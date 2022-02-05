import requests
import warnings
import common

warnings.filterwarnings('ignore')


def create_topic(session, username):
    add_url = common.host + "/posts.json"
    data = {
        "raw": common.random_string(50),
        "title": common.random_string(50),
        "category": "3",
    }
    header = common.header
    header["Api_Username"] = username
    response = session.post(url=add_url, headers=common.header, data=data, verify=False)
    if response.status_code != 200:
        exit(-1)


if __name__ == '__main__':
    session = requests.Session()
    create_topic(session, "1qqcom")
    session.close()
