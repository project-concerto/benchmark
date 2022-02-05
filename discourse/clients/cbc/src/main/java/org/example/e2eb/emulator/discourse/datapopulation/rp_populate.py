import create_topics
import create_user
import create_posts
import requests

if __name__ == '__main__':
    num = 256
    # create_user.populate_users(num)
    #
    # # All users create a topic first
    # for i in range(1, num + 1):
    #     s = requests.session()
    #     create_topics.create_topic(s, str(i)+"qqcom")
    #     s.close()

    # All users create a post with image in its own topic
    for i in range(1, num + 1):
        topic_id = 9 + i
        s = requests.session()
        create_posts.create_post_with_image(s, topic_id, i)
        s.close()
