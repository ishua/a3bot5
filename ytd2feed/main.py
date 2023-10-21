#!/usr/bin/python3
import app
import os
import redis
from multiprocessing import Process
import time
import json

# def start_download(data: dict):
#     print("try to download")

#     c = app.Conf()
#     fc = FeedCreater(data, l, c)
#     download(data, Conf(), l, fc)
#     print("Download complited")



if __name__ == '__main__':
    print("start app")
    cfg = app.Conf()
    print("redis host: {}, listern channel: {}".format(cfg.redis, cfg.channel))
    # init content path
    if not os.path.isdir(cfg.path2content):
         os.makedirs(cfg.path2content)

    fo = cfg.user("AlekseyIm")
    print(fo)
    fc = app.FeedCreater(cfg.path2content, cfg.url2content, fo["feedName"], fo["feedDescription"])

    print("start")
    app.download(
        url = 'https://www.youtube.com/watch?v=O11dfJVJusk',
        format = fo["format"],
        fc=fc
    )

    # init redis
    # r = redis.Redis(host=cfg.redis)
    # p = r.pubsub()
    # p.subscribe(cfg.channel)

    # print("Start to lisen")
    # while True:
    #     try:
    #         message = p.get_message()
    #     except redis.ConnectionError:
    #         # Do reconnection attempts here such as sleeping and retrying
    #         print("reconnect to redis after 3 sec")
    #         time.sleep(3)
    #         p = r.pubsub()
    #         p.subscribe(cfg.channel)
    #     if message:
    #         if message["type"] == "message":
    #             m = json.loads(message["data"])
    #             # start download in background
    #             proccess = Process(target=start_download, args=(m,))
    #             proccess.start()
    #     time.sleep(1)  # be nice to the system :)
