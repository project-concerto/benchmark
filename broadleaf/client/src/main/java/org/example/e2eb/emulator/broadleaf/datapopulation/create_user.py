import requests
import re
import warnings
warnings.filterwarnings('ignore')

header = {"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36", "Connection": "keep-alive"}


def register(session, num):
    home_url = "https://localhost:8443/"
    response = session.get(url=home_url, headers=header, verify=False)
    token = re.search('\"csrfToken\":\"(.*?)\"', response.text).group(1)

    register_url = "https://localhost:8443/register"
    data = { "redirectUrl": "",
             "customer.emailAddress": str(num)+"@qq.com",
             "customer.firstName": str(num),
             "customer.lastName": "qq",
             "password": "zxd123",
             "passwordConfirm": "zxd123",
             "csrfToken": token
            }
    response = session.post(url=register_url, data=data, headers=header, verify=False)
    print(response.text)

for i in range(1,1025):
    session = requests.Session()
    register(session, i)
    session.close()
