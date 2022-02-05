import create_topics
import create_user
import create_posts
import requests

if __name__ == '__main__':
    num = 256
    create_user.populate_users(num)

    # make odd user create topic
    # new created topic_id start from 140 (inclusive)
    # related posts id start from 143 (inclusive)
    for i in range(1, num, 2):
        s = requests.session()
        create_topics.create_topic(s, str(i) + "qqcom")
        s.close()

    # make even create posts for correspond topics
    # new created posts start from 143 + num/2
    for i in range(2, num + 1, 2):
        topic_id = i/2 + 15
        s = requests.session()
        create_posts.create_post(s, topic_id, str(i) + "qqcom")
        create_posts.create_post(s, topic_id, str(i) + "qqcom")
        s.close()
