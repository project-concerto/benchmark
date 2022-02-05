import requests
import warnings

import common

warnings.filterwarnings('ignore')


def set_answer(session, post_id):
    data = {
        "id": str(post_id)
    }
    add_url = common.host + "/solution/accept"
    response = session.post(url=add_url, headers=common.header, data=data, verify=False)
    assert (response.status_code == 200)

def unset_answer(session, post_id):
    add_url = common.host + "/solution/unaccept"
    data = {
        "id": str(post_id)
    }
    response = session.post(url=add_url, headers=common.header, data=data, verify=False)
    assert (response.status_code == 200)


if __name__ == '__main__':
    session = requests.Session()
    set_answer(session, 399)
    set_answer(session, 400)
    session.close()
