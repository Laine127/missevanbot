import redis

server = {
    "host": "localhost",
    "port": 6379,
    "db": 0,
    "password": ""
}

r = redis.Redis(host=server["host"], port=server["port"], db=server["db"], password=server["password"])


def get_rid(key):
    return str.split(decode(key), ":")[2]


def decode(val: bytes):
    if val is None:
        return "NONE"
    s = val.decode("utf-8")
    if s == "":
        return "NONE"
    return s


def get_info(rid):
    key = "missevan:{}:info".format(rid)
    alias = r.hget(key, "alias")
    bot = r.hget(key, "bot")
    ret = "BOT={}   ALIAS={}".format(decode(bot), decode(alias).ljust(15))
    return ret


def get_detail(key):
    info = get_info(get_rid(key))

    online = r.hget(key, "online")
    count = r.hget(key, "count")
    game = r.hget(key, "game")

    ret = "\n\tONLINE={}   COUNT={}   GAME={}".format(decode(online), decode(count), decode(game))
    return info + ret


def current_alive(alive_rooms):
    print("当前运行：")
    if len(alive_rooms) == 0:
        print("NONE")
        return
    for idx, key in enumerate(alive_rooms):
        print("[RUNNING #{}] ".format(idx + 1), get_info(get_rid(key)))


def current_online(online_rooms):
    print("当前在线：")
    if len(online_rooms) == 0:
        print("NONE")
        return
    for idx, key in enumerate(online_rooms):
        print("[ONLINE #{}]".format(idx + 1), get_detail(key))


if __name__ == "__main__":
    current_alive(r.keys("missevan:alive:*"))
    current_online(r.keys("missevan:online:*"))
