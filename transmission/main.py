#!/usr/bin/python3
import app
import os
import redis
from multiprocessing import Process
import time

def go_command(a: app.BotMsg, trhost: str, trport: int, tdownloaddir: str):
    print(trhost, trport, a.msg, a.redis, a.channel)
    t = app.Tr(trhost, trport, tdownloaddir)
    a.runBotCommand(t)   
    # 
    # reply = t.add_torrent(a.getTorrentUrl())
    # a.reply(reply)

    

if __name__ == '__main__':
    print("start app")
    cfg = app.Conf()
    print("redis host: {}, listern channel: {}".format(cfg.redis, cfg.channel))

    r = redis.Redis(host=cfg.redis)
    p = r.pubsub()
    p.subscribe(cfg.channel)

    print("Start to lisen")
    while True:
        try:
            message = p.get_message()
        except redis.ConnectionError:
            # Do reconnection attempts here such as sleeping and retrying
            print("reconnect to redis after 3 sec")
            time.sleep(3)
            p = r.pubsub()
            p.subscribe(cfg.channel)
        if message:
            if message["type"] == "message":
                payload = message["data"]
                a = app.BotMsg(payload, cfg.redis, cfg.tbotchannel)
                # start in background
                proccess = Process(target=go_command, 
                                   args=(a, cfg.trhost, cfg.trport, cfg.tdownloaddir))
                proccess.start()
        time.sleep(1)  # be nice to the system :)
