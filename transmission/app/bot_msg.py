import json
from urllib.parse import urlparse
import redis
from app import Tr as client

# type ChannelMsg struct {
# 	Command          string `json:"command"`
# 	UserName         string `json:"userName"`
# 	MsgId            int    `json:"msgId"`
# 	ReplyToMessageID int    `json:"replyToMessageID"`
# 	ChatId           int64  `json:"chatId"`
# 	Text             string `json:"text"`
# 	ReplyText        string `json:"replyText"`
# 	Caption          string `json:"caption"`
# 	FileUrl          string `json:"fileUrl"`
# }

class BotMsg:
    def __init__( self, payload: str, redis: str, channel: str):
        self.msg = json.loads(payload)
        self.redis = redis
        self.channel = channel
        
    def setReply( self, reply: str):
        self.msg["ReplyText"] = reply

    def getReply( self ) -> str:
        return self.msg.get("ReplyText", "")
    
    def getUserName( self) -> str:
        return self.msg.get("userName", "")
    
    def getTorrentUrl (self ) -> str:
        return self.msg.get("fileUrl")
    
    def validateError( self ) -> bool:
        if self.msg.get("caption") is None and self.msg.get("text") is None:
            self.reply("no command")
            return True
        
        commands = ""
        if self.msg.get("text") != "":
            commands = self.msg["text"].split(" ")
        else:
            commands = self.msg["caption"].split(" ")
            _uri = urlparse(self.msg.get("fileUrl"))
            if _uri.hostname is None:
                self.reply("no link to torrent file")
                return True

        if len(commands) < 2:
            self.reply("command is too small")
            return True
        
        return False


    def reply( self, replyText: str ):
        self.setReply(replyText)
        r = redis.Redis(self.redis)
        r.publish(self.channel, json.dumps(self.msg))

    def runBotCommand(self, c: client):
        if self.validateError():
            return
        commands = ""
        if self.msg.get("text") != "":
            commands = self.msg["text"].split(" ")
        else:
            commands = self.msg["caption"].split(" ")

        if commands[1] == "add":
            if len(commands) < 3:
                self.reply("Can't add torrent. Need a label")
                return
            t_url = self.getTorrentUrl()
            if t_url == "":
                self.reply("Can't add torrent. Need a torrent file")
                return 
            self.reply(c.addTorrent(commands[2], t_url))
            return

        if commands[1] == "list":
            self.reply(c.listTorrents())
            return
        
        if commands[1] == "status":
            if len(commands) < 3:
                self.reply("Can't check torrent. Need torrent id")
                return
            self.reply(c.statusTorrent(int(commands[2])))
            return
        
        if commands[1] == "del":
            if len(commands) < 3:
                self.reply("Can't del torrent. Need torrent id")
                return
            self.reply(c.delTorrent(int(commands[2])))
            return
        
        if commands[1] == "help":
            _reply = "This is help for /torrent command \n"
            _reply += "start words: /torrent, t, T\n"
            _reply += "next words:\n"
            _reply += '- "add" + category: "movie/m","shows/s", "cartoon/c" - add torrents in category\n'
            _reply += '- "list" - list torrens in actions\n'
            _reply += '- "status #" - where # it is id torrent\n'
            _reply += '- "del #" - where # it is id torrent\n'
            self.reply(_reply)
            return
        
        self.reply("don't know command: " + commands[1])

