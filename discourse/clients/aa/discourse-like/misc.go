package discourse

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"go.uber.org/zap"
)

var _topic_posts map[int][]int

const topicOwnersFileName = "discourse-like/topic_posts"

func GetStoredLargeTopcPosts() map[int][]int {
	if _topic_posts == nil {
		_topic_posts = make(map[int][]int)
		for _, topicId := range TargetTopics {
			postIds := GetTopicPosts(topicId)
			_topic_posts[topicId] = postIds
		}
	}
	return _topic_posts
}

func GetRealTimeTopicPosts(topicIds []int) map[int][]int {
	resp := make(map[int][]int)
	for _, topicId := range topicIds {
		postIds := GetTopicPosts(topicId)
		resp[topicId] = postIds
	}
	return resp
}

func Get(url string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return bodyBytes
}

func Request(method string, username string, payload []byte, location string) ([]byte, bool) {
	body := bytes.NewReader(payload)

	url := fmt.Sprintf("http://%v/%v", BaseUrl, location)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Api-Key", ApiKey)
	req.Header.Set("Api-Username", username)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode/100 != 2 {
		zap.L().Debug("Request Finished",
			zap.String("method", method),
			zap.String("url", url),
			zap.String("payload", string(payload)),
			zap.String("username", username),
			zap.Int("StatusCode", resp.StatusCode),
			zap.String("body", string(bodyBytes)),
		)
		return nil, false
	}

	return bodyBytes, true
}

func GetCategoryIds() []int {
	type SiteResp struct {
		CategoryList struct {
			Categories []struct {
				ID int `json:"id"`
			} `json:"categories"`
		} `json:"category_list"`
	}

	resp := Get("http://%s/categories.json")
	var siteResp SiteResp
	err := json.Unmarshal(resp, &siteResp)
	if err != nil {
		log.Fatal(err)
	}

	// Extract the category ids
	ids := make([]int, len(siteResp.CategoryList.Categories))
	for _, category := range siteResp.CategoryList.Categories {
		ids = append(ids, category.ID)
	}
	return ids
}

func GetTopicPosts(topicId int) []int {
	type GetTopicResp struct {
		PostStream struct {
			Stream []int `json:"stream"`
		} `json:"post_stream"`
	}

	url := fmt.Sprintf("http://%v/t/%v.json", BaseUrl, topicId)
	resp := Get(url)
	var getTopicResp GetTopicResp
	err := json.Unmarshal(resp, &getTopicResp)
	if err != nil {
		log.Fatal(err)
	}

	return getTopicResp.PostStream.Stream
}

func getPostOwner(postId int) string {
	type GetPostResp struct {
		Username string `json:"username"`
	}

	url := fmt.Sprintf("http://%v/posts/%v.json", BaseUrl, postId)
	resp := Get(url)

	var getPostResp GetPostResp
	err := json.Unmarshal(resp, &getPostResp)
	if err != nil {
		log.Fatal(err)
	}

	return getPostResp.Username
}

func PreparePostsOwners() {
	l := zap.L()

	_topic_posts = make(map[int][]int)
	for _, topicId := range TargetTopics {
		postIds := GetTopicPosts(topicId)
		l.Info("Got topic posts",
			zap.Int("topicId", topicId),
			zap.Int("postIds Length", len(postIds)),
		)
		_topic_posts[topicId] = postIds
	}

	f, err := os.Create(topicOwnersFileName)
	defer f.Close()
	if err != nil {
		l.Fatal("Failed to create spree_order_number", zap.Error(err))
	}

	for topicId, postIds := range _topic_posts {
		for _, postId := range postIds {
			fmt.Fprintf(f, "%v %v\n", topicId, postId)
		}
	}
}

// Preparation to fill some post_actions
func LikeEveryPostWithOneUser(userId int) {
	taskChan := make(chan int)
	stopChan := make(chan bool)

	l := zap.L()

	go func() {
		for _, topicId := range AllTopics {
			postIds := GetTopicPosts(topicId)
			l.Info("Got topic posts",
				zap.Int("topicId", topicId),
				zap.Int("postIds Length", len(postIds)),
			)
			for _, postId := range postIds {
				taskChan <- postId
			}
		}
		close(stopChan)
	}()

	for i := 0; i < 12; i++ {
		go func() {
			for {
				select {
				case postId := <-taskChan:
					// fmt.Println("Like post", postId)
					dp := DiscoursePost{
						postId:   postId,
						username: fmt.Sprintf("admin%v", userId),
						doLike:   true,
					}
					dp.Run(context.Background())
				case <-stopChan:
					return
				}
			}
		}()
	}

	<-stopChan
}
