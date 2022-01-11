import base64
import hashlib
import xml.dom.minidom
import urllib.parse


class head:
    def __init__(self, k, v):
        self.k = k
        self.v = v


def get_xml(file_path):
    burp = xml.dom.minidom.parse(file_path)
    root = burp.documentElement
    items = root.getElementsByTagName('item')
    return items


def get_header_maps(xml_data):
    iinstagram_header_list = {}
    for item in xml_data:
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
                    pass
                    # headers.append(head(line[:line.find(" ")], line[line.find(" ") + 1:]))
                else:
                    headers.append(head(line[:line.find(":")], line[line.find(":") + 1:]))
            if not iinstagram_header_list.get(purl.path):
                iinstagram_header_list[purl.path] = []
            iinstagram_header_list[purl.path].append(headers)
    return iinstagram_header_list


def merge_header_list_map(l1, l2):
    l3 = {}
    for path in l1:
        if not l3.get(path):
            l3[path] = []
        for item in l1[path]:
            l3[path].append(item)

    for path in l2:
        if not l3.get(path):
            l3[path] = []
        for item in l2[path]:
            l3[path].append(item)
    return l3


def list_md5(l):
    s = ""
    for head in l:
        s += head.k
    m = hashlib.md5()
    m.update(s.encode())
    return m.hexdigest()


def print_header(header):
    s = ""
    for item in header:
        s += item.k + ","
        print(item.k + "\t" + item.v)
        if item.v == "":
            print("k is null ", item.k)
    return s


def make_map(iinstagram_header_list):
    header_map = {}
    for path in iinstagram_header_list:
        _m = ""
        for header in iinstagram_header_list[path]:
            m5 = list_md5(header)
            if _m == "":
                _m = m5
            if _m != m5:
                print("not same:", path)
            header_map[m5] = header

    spath = set()
    ppath = []
    for path in iinstagram_header_list:
        for header in iinstagram_header_list[path]:
            key = path + list_md5(header)
            if key not in spath:
                spath.add(key)
                ppath.append({
                    "path": path,
                    "md5": list_md5(header)
                })

    pmd5 = []
    for md5 in header_map:
        pmd5.append({
            "md5": md5,
            "header": print_header(header_map[md5])
        })
    return ppath, pmd5


file_name = ["2021-12-31-一些操作burp.xml",
             "2022-01-03-获取评论等.xml",
             "第一次打开",
             "登录.xml",
             "邮箱失败.xml",
             "邮箱失败2.xml",
             "注册成功.xml",
             "第二次安装app.xml"]
maps = []
for file in file_name:
    file_path = "./抓包/" + file
    header_maps1 = get_header_maps(get_xml(file_path))
    maps.append(header_maps1)

old = None
for i in maps:
    if old is None:
        old = i
        continue

    old = merge_header_list_map(i, old)

ppath, pmd5 = make_map(old)
print(ppath)
print(pmd5)
