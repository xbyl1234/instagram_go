import base64
import xml.dom.minidom
import urllib.parse


class head:
    def __init__(self, k, v):
        self.k = k
        self.v = v


file_name = "2021-12-31-一些操作burp.xml"
file_path = "./抓包/" + file_name
burp = xml.dom.minidom.parse(file_path)
root = burp.documentElement
items = root.getElementsByTagName('item')

iinstagram_header_list = []
for item in items:
    url = item.getElementsByTagName("url")[0].firstChild.data
    purl = urllib.parse.urlparse(url)
    if purl.netloc == "i.instagram.com":
        req_body = base64.b64decode(item.getElementsByTagName("request")[0].firstChild.data)
        request = req_body[:req_body.find(b"\r\n\r\n")].decode()

        sp = request.split("\r\n")
        headers = []
        for line in sp:
            if line == "":
                break

            if line.startswith("POST") or line.startswith("GET"):
                headers.append(head([line[:line.find(" ")]], line[line.find(" ") + 1:]))
            else:
                headers.append(head([line[:line.find(":")]], line[line.find(":") + 1:]))

        print(headers)
