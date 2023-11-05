import yaml
import sys

class Conf():
    def __init__(self):
        self.conf = {}
        with open("conf/transmission_config.yaml") as f:
            self.conf = yaml.load(f, Loader=yaml.FullLoader)

    
    @property
    def redis( self ) ->str: 
        return self.conf.get("redis", "redis")
    
    @property
    def channel( self ) ->str: 
        return self.conf.get("channel", "transmission")
    
    @property
    def tbotchannel( self ) ->str: 
        return self.conf.get("tbotchannel", "tbot")
    
    @property
    def trhost( self ) ->str: 
        return self.conf.get("trhost", "transmission")
    
    @property
    def trport( self ) ->int: 
        return self.conf.get("trport", 9091)
    
    @property
    def tdownloaddir( self ) ->int: 
        return self.conf.get("tdownload_dir", "/downloads/complete/")

            
    