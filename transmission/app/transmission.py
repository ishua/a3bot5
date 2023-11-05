import requests

from transmission_rpc import Client

class Tr:
    def __init__( self, host: str, port: int, download_dir: str):
        self.client = Client(host=host, port=port)
        self.download_dir = download_dir
       
    def addTorrent(self, type_file: str, torrent_url: str) -> str:
        _download_dir = self.download_dir
        if type_file == "m" or type_file == "movie":
            _download_dir += "movie"
        if type_file == "s" or type_file == "shows":
            _download_dir += "shows"
        if type_file == "c" or type_file == "cartoon":
            _download_dir += "cartoon"
        if type_file == "a" or type_file == "audiobook":
            _download_dir += "audiobook"
        if type_file == "ap" or type_file == "audiobook_p":
            _download_dir += "audiobook_p"
        if type_file == "c" or type_file == "cartoon_s":
            _download_dir += "cartoon_s"
        if _download_dir == self.download_dir:
            return "wrong category"
        
        try:
            t = self.client.add_torrent(torrent=torrent_url, download_dir=_download_dir)
        except Exception as e:
            return e.__str__()
            
        return str(t.id) + " " + t.name
    
    def listTorrents(self) -> str:
        try:
            tlist = self.client.get_torrents()
        except Exception as e:
            return e.__str__()
        if len(tlist) == 0:
            return "no torrents"
        ret = "torrent list \n"
        for t in tlist:
            ret += str(t.id) 
            ret +=  "-" + t.name 
            ret += "-" + t.status
            ret += "-" + str(t.progress) 
            ret += "-" + '{:.3f}'.format(t.rate_download/1024/1024)
            ret += "\n" 
        return ret

    def statusTorrent(self, torrent_id: int) -> str:
        try:
            t = self.client.get_torrent(torrent_id)
        except Exception as e:
            return e.__str__()
        return str(t.id) + " " + t.name

    def delTorrent(self, torrent_id: int) -> str:
        try:
            t = self.client.remove_torrent(torrent_id)
        except Exception as e:
            return e.__str__()
        return "ok"