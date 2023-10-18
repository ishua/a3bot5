import yaml

class Conf():
    def __init__(self):
        self.conf = {}
        with open("conf/ytdf_config.yaml") as f:
            self.conf = yaml.load(f, Loader=yaml.FullLoader)
    
    @property
    def redis( self ) ->str: 
        return self.conf["redis"] 