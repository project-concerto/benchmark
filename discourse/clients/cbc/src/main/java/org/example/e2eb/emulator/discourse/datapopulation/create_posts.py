import requests
import warnings
import json
import common

warnings.filterwarnings('ignore')


def create_post(session, topic_id, username):
    add_url = common.host + "/posts.json"
    data = {
        "raw": common.random_string(50),
        "topic_id": str(topic_id)
    }
    header = common.header
    header["Api_Username"] = username
    response = session.post(url=add_url, headers=common.header, data=data, verify=False)
    if response.status_code != 200:
        exit(-1)


def get_url_by_id(user_id):
    image_id = (user_id + 7) / 8
    image_id = int(image_id)
    file = open("./images.json", "r")
    images = json.load(file)
    return "![{}|{}x{}]({})".format(image_id, images[str(image_id)]["thumbnail_width"],
                                    images[str(image_id)]["thumbnail_height"],
                                    images[str(image_id)]["short_url"])


def create_post_with_image(session, topic_id, user_id):
    add_url = common.host + "/posts.json"
    image_url = get_url_by_id(user_id)
    data = {
        "raw": common.random_string(50) + image_url,
        "topic_id": str(topic_id)
    }
    header = common.header
    header["Api_Username"] = str(user_id) + "qqcom"
    response = session.post(url=add_url, headers=common.header, data=data, verify=False)
    if response.status_code != 200:
        exit(-1)


if __name__ == '__main__':
    session = requests.Session()
    create_post(session, 10, "1qqcom")
    session.close()
