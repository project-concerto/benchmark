import requests
import warnings

import common

warnings.filterwarnings('ignore')


def register(session, num):
    register_url = common.host + "/users.json"
    data = {
        "name": str(num) + "qqcom",
        "email": str(num) + "@qq.com",
        "password": "zxd1234567",
        "username": str(num) + "qqcom",
        "active": "true",
        "approved": "true"
    }
    response = session.post(url=register_url, data=data, headers=common.header, verify=False)
    if response.status_code != 200:
        exit(-1)


def populate_users(num):
    for i in range(1, num + 1):
        s = requests.Session()
        register(s, i)
        s.close()


if __name__ == '__main__':
    session = requests.Session()
    register(session, 0)
    session.close()
