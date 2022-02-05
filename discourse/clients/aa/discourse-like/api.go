package discourse

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"associated-access/utils"
	"log"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

var TopicsNum = flag.Int("topics", 1, "Number of topics to like")
var threadPerTopic = flag.Int("tpt", 12, "Thread per topic. Only for old.")
var docker = flag.Bool("docker", false, "Prepare Discourse")

type DiscoursePost struct {
	postId   int
	username string
	doLike   bool
}

func (d *DiscoursePost) like(postId int, username string) bool {
	type Payload struct {
		ID               int  `json:"id"`
		PostActionTypeID int  `json:"post_action_type_id"`
		FlagTopic        bool `json:"flag_topic"`
	}

	data := Payload{
		// fill struct
		ID:               postId,
		PostActionTypeID: 2,
		FlagTopic:        false,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	resp, ok := Request("POST", username, payloadBytes, "post_actions.json")
	if !ok {
		zap.L().Debug("Failed to like post", zap.String("response", string(resp)))
	}

	return ok
}

func (d *DiscoursePost) dislike(postId int, username string) bool {
	type Payload struct {
		PostActionTypeID int `json:"post_action_type_id"`
	}
	payload := Payload{
		PostActionTypeID: 2,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	_, ok := Request("DELETE", d.username, payloadBytes, fmt.Sprintf("post_actions/%v.json", d.postId))

	return ok
}

func (d *DiscoursePost) Run(ctx context.Context) utils.APIResult {
	result := utils.APIResult{
		SuccessLatency: nil,
		FailLatencies:  make([]int64, 0),
	}
	start := time.Now()
	var ok bool
	if d.doLike {
		ok = d.like(d.postId, d.username)
	} else {
		ok = d.dislike(d.postId, d.username)
	}
	elapsed := time.Since(start).Microseconds()
	if ok {
		result.SuccessLatency = &elapsed
	} else {
		result.FailLatencies = append(result.FailLatencies, elapsed)
	}
	return result
}

type DiscourseIncresingTopicFactory struct {
	workForThread [][]DiscoursePost
}

func NewDiscourseFactory(threads int) *DiscourseIncresingTopicFactory {
	l := zap.L()
	l.Info("NewDiscourseFactory",
		zap.Int("topicsNum", *TopicsNum),
		zap.Int("threadPerTopic", *threadPerTopic))

	l.Info("Getting topics and its posts from Discourse API")
	topicPosts := GetStoredLargeTopcPosts()
	topics := make([]int, 0)
	for topicId := range topicPosts {
		topics = append(topics, topicId)
	}
	topicWork := make(map[int][]DiscoursePost)
	topicWorkTotal := make(map[int]int)

	l.Info("Making work for each thread")
	workForThread := make([][]DiscoursePost, 0)
	for i := 0; i < threads; i++ {
		targetTopic := topics[i/(*threadPerTopic)]
		if _, ok := topicWork[targetTopic]; !ok {
			l.Info("Making work for topic", zap.Int("topicId", targetTopic))
			work := make([]DiscoursePost, 0)
			for _, post := range topicPosts[targetTopic] {
				for user := 0; user < UserCount; user++ {
					work = append(work, DiscoursePost{
						postId:   post,
						username: fmt.Sprintf("admin%v", user),
						doLike:   true,
					})
				}
			}
			rand.Shuffle(len(work), func(i int, j int) {
				work[i], work[j] = work[j], work[i]
			})
			topicWork[targetTopic] = work
			topicWorkTotal[targetTopic] = len(work)
		}
		numWork := topicWorkTotal[targetTopic] / Min(*threadPerTopic, threads-i/(*threadPerTopic)*(*threadPerTopic))
		workForThread = append(workForThread, topicWork[targetTopic][0:numWork])
		topicWork[targetTopic] = topicWork[targetTopic][numWork:]
		l.Info("Assigning work for this thread",
			zap.Int("NumberOfWorkForThread", len(workForThread[i])),
			zap.Int("NumberOfWorkForTopic", topicWorkTotal[targetTopic]),
		)
	}

	return &DiscourseIncresingTopicFactory{
		workForThread: workForThread,
	}
}

func (f *DiscourseIncresingTopicFactory) Prepare() {
	PrepareDiscourseDocker()
}

func (f *DiscourseIncresingTopicFactory) Make(threadId int) utils.API {
	work := f.workForThread[threadId][0]
	f.workForThread[threadId] = f.workForThread[threadId][1:]
	zap.L().Debug("Making work",
		zap.Int("threadId", threadId),
		zap.Int("Remaining", len(f.workForThread[threadId])),
		zap.Int("postId", work.postId),
		zap.String("username", work.username),
	)
	return &work
}

func (f *DiscourseIncresingTopicFactory) Stop() {}

type DiscourseConstantTopicFactory struct {
	cancel   context.CancelFunc
	workChan chan DiscoursePost
}

func NewDiscourseConstantTopicFactory() *DiscourseConstantTopicFactory {
	l := zap.L()
	l.Info("NewDiscourseConstantTopicFactory", zap.Int("topicsNum", *TopicsNum))

	work := make([]DiscoursePost, 0)

	topicPosts := GetStoredLargeTopcPosts()
	cnt := 0
	for topic, posts := range topicPosts {
		cnt++
		if cnt > *TopicsNum {
			break
		}
		l.Info("Adding work", zap.Int("TopicId", topic))
		for _, post := range posts {
			for user := 0; user < UserCount; user++ {
				work = append(work, DiscoursePost{
					postId:   post,
					username: fmt.Sprintf("admin%v", user),
					doLike:   true,
				})
			}
		}
	}

	rand.Shuffle(len(work), func(i, j int) {
		work[i], work[j] = work[j], work[i]
	})

	ctx, cancel := context.WithCancel(context.Background())
	workChan := make(chan DiscoursePost, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case workChan <- work[0]:
				work = work[1:]
			}
		}
	}()

	return &DiscourseConstantTopicFactory{
		cancel:   cancel,
		workChan: workChan,
	}
}

func (f *DiscourseConstantTopicFactory) Prepare() {
	PrepareDiscourseDocker()
}

func (f *DiscourseConstantTopicFactory) Make(threadId int) utils.API {
	work := <-f.workChan
	zap.L().Debug("DiscourseConstantTopicFactory Make",
		zap.Int("postId", work.postId),
		zap.String("username", work.username),
	)
	return &work
}

func (f *DiscourseConstantTopicFactory) Stop() {
	f.cancel()
}

// Every thread will use one topic exclusively.
// So there is less contention.
type DiscourseNoContentionFactory struct {
	workForThread [][]DiscoursePost
}

func NewDiscourseNoContentionFactory(threads int, low_contention bool) *DiscourseNoContentionFactory {
	l := zap.L()

	// Find the topic that each thread should reply to.
	topicsForThread := make(map[int][]int, 0)
	for i := 0; i < threads; i++ {
		topicsForThread[i] = make([]int, 0)
	}
	usedTopicIds := make([]int, 0)
	if low_contention {
		dist := len(AllTopics) / threads
		for i := 0; i < threads; i++ {
			topicsForThread[i] = append(topicsForThread[i], AllTopics[i*dist])
			usedTopicIds = append(usedTopicIds, AllTopics[i*dist])
		}
	} else {
		for i := 0; i < *TopicsNum; i++ {
			topicId := TargetTopics[i]
			usedTopicIds = append(usedTopicIds, topicId)
			for t := 0; t < threads; t++ {
				topicsForThread[t] = append(topicsForThread[t], topicId)
			}
		}
	}

	topicPosts := GetRealTimeTopicPosts(usedTopicIds)

	workForThread := make([][]DiscoursePost, 0)
	userPerThread := UserCount / threads
	for i := 0; i < threads; i++ {
		work := make([]DiscoursePost, 0)
		for _, topicId := range topicsForThread[i] {
			for _, post := range topicPosts[topicId] {
				for user := i * userPerThread; user < (i+1)*userPerThread; user++ {
					work = append(work, DiscoursePost{
						postId:   post,
						username: fmt.Sprintf("admin%v", user),
						doLike:   true,
					})
				}
			}
		}
		rand.Shuffle(len(work), func(i, j int) {
			work[i], work[j] = work[j], work[i]
		})
		workForThread = append(workForThread, work)
		l.Debug("Assigning work for this thread",
			zap.Int("Number of Topics", len(topicsForThread[i])),
			zap.Int("Number of WorkForThread", len(workForThread[i])),
		)
	}

	return &DiscourseNoContentionFactory{
		workForThread: workForThread,
	}
}

func (f *DiscourseNoContentionFactory) Prepare() {
	PrepareDiscourseDocker()
}

func (f *DiscourseNoContentionFactory) Make(threadId int) utils.API {
	work := f.workForThread[threadId][0]
	f.workForThread[threadId] = f.workForThread[threadId][1:]
	zap.L().Debug("DiscourseNoContentionFactory Make",
		zap.Int("postId", work.postId),
		zap.String("username", work.username),
	)
	return &work
}

func (f *DiscourseNoContentionFactory) Stop() {}

func PrepareDiscourseDocker() {
	l := zap.L()
	if !*docker {
		l.Info("Skip preparing discourse docker")
		return
	}

	l.Info("Preparing Discourse API")
	output, err := exec.Command(RefreshScript).CombinedOutput()
	if err != nil {
		l.Fatal("Failed to refresh Discourse",
			zap.Error(err),
			zap.String("output", string(output)),
		)
	}

	l.Info("Waiting for Discourse to be ready")
	for {
		output, err = exec.Command("bash", "-c", "docker stats --no-stream --format \"{{.CPUPerc}}\" app").CombinedOutput()
		if err != nil {
			l.Fatal("Failed to get Discourse CPU usage",
				zap.Error(err),
				zap.String("output", string(output)),
			)
		}

		percent, err := strconv.ParseFloat(strings.TrimSuffix(string(output), "%\n"), 64)
		if err != nil {
			l.Fatal("Failed to parse Discourse CPU usage",
				zap.Error(err),
				zap.String("output", string(output)),
			)
		}
		l.Info("Discourse CPU usage", zap.Float64("percent", percent))
		if percent < 10.0 {
			break
		}

		time.Sleep(time.Second)
	}
	l.Info("Refreshed Discourse")
}

func Min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
