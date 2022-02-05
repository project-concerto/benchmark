import random

# Server Info
header = {
    "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 "
                  "Safari/537.36",
    "Connection": "keep-alive",
    "Api_Key": "1f863193afbbe719ac251c29ee73b749a562c4395d37e5c7595986d72e4c87e6",
    "Api_Username": "zxd"
}

host = "http://localhost:3000"

names = ["We ", "I ", "They ", "He ", "She ", "Jack ", "Jim ", "Bxb "]
verbs = ["was ", "is ", "were ", "are ", "do ", "does ", "doing ", "done ", "did "]
nouns = ["playing a game ", "watching television ", "talking ", "dancing ", "speaking ", "like ", "love ", "fuck "]


def random_string(length):
    res = ""
    while len(res) <= length:
        res += random.choice(names)
        res += random.choice(verbs)
        res += random.choice(nouns)
    return res[0:length]
