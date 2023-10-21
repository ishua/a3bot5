import yaml
import sys

class Conf():
    def __init__(self):
        self.conf = {}
        with open("conf/ytdf_config.yaml") as f:
            self.conf = yaml.load(f, Loader=yaml.FullLoader)

        self.users =  self.conf.get("users")
        if len(self.users) < 1:
            print("need users in config")
            sys.exit(1) 
    
    @property
    def redis( self ) ->str: 
        return self.conf.get("redis", "redis:6379")
    
    @property
    def channel( self ) ->str: 
        return self.conf.get("channel", "ytd2feed")
    
    @property
    def path2content( self ) ->str: 
        return self.conf.get("path2content", "temp")
    
    @property
    def url2content( self ) ->str: 
        return self.conf.get("url2content", "")
    
    def user( self, userName: str ) -> dict:
        _user = {}
        for u in self.users:
            if u["user"] == userName:
                _user["name"] = u["user"]
                _user["feedName"] = u["feed_name"]
                _user["feedDescription"] = u["feed_description"]
                _user["format"] = u["format"]

        return _user
            
    