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


all_header = set()
header_values = set()


def get_header_maps(file_name, xml_data):
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
                    h = head(line[:line.find(":")], line[line.find(":") + 1:])
                    if h.k == "X-Ig-Abr-Connection-Speed-Kbps":
                        header_values.add(h.v)
                    all_header.add(h.k)
                    headers.append(h)
            if not iinstagram_header_list.get(purl.path + " - " + file_name):
                iinstagram_header_list[purl.path + " - " + file_name] = []
            iinstagram_header_list[purl.path + " - " + file_name].append(headers)
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
        s += head.k + ","
    m = hashlib.md5()
    m.update(s.encode())
    return m.hexdigest()


def str_md5(s):
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


def md5_in_list(md5, l):
    for item in l:
        if md5 == item["md5"]:
            return True
    return False


def remove_cookies(ppath, pmd5):
    tags = ["X-Ig-Extended-Cdn-Thumbnail-Sizes,", "Cookie,"]
    for item in pmd5:
        if tags[0] in item["header"] and tags[1] in item["header"]:
            nocmd5 = str_md5(item["header"].replace(tags[0], "").replace(tags[1], ""))
            if md5_in_list(nocmd5, pmd5):
                for pitem in ppath:
                    if nocmd5 in pitem["md5"]:
                        print("find no ", tags[0], tags[1], nocmd5)
                        pitem["md5"] = "no-" + tags[0] + tags[1] + item["md5"]

    for tag in tags:
        for item in pmd5:
            if tag in item["header"]:
                nocmd5 = str_md5(item["header"].replace(tag, ""))
                if md5_in_list(nocmd5, pmd5):
                    for pitem in ppath:
                        if nocmd5 in pitem["md5"]:
                            print("find no ", tag, nocmd5)
                            pitem["md5"] = "no-" + tag + item["md5"]
    pmd5set = set()
    pset = set()
    ret = []
    for item in ppath:
        sp = item["md5"].split(",")
        if sp[len(sp) - 1] not in pmd5set:
            pmd5set.add(sp[len(sp) - 1])
        if item["md5"] + item["path"] not in pset:
            pset.add(item["md5"] + item["path"])
            ret.append(item)
    return ret, pmd5set


# file_name = ["2021-12-31-一些操作burp.xml",
#              "2022-01-03-获取评论等.xml",
#              "第一次打开",
#              "登录.xml",
#              "邮箱失败.xml",
#              "邮箱失败2.xml",
#              "注册成功.xml",
#              "第二次安装app.xml"]

file_name = [
    # '2022-01-12-第一次安装第一次打开-注册.xml',
    # '2022-01-12-第一次安装第二次打开-注册.xml',
    # '2022-01-12-第一次安装第二次打开-登录.xml'
    # "2022-01-12-第一次安装第一次打开-手机注册-注册后.xml"
    # "2021-12-31-一些操作burp.xml",
    # "2022-01-03-获取评论等.xml"
    # "190聊天1.xml"
    # "语音.xml"
    # "私信发图视频关注.xml"
    # "聊天链接等.xml"
    # "发本地短视频2.xml",
    # "2022-02-22-置顶快拍.xml",
    # "发图片帖子拍照.xml",
    # "2022-02-22-发多个帖子.xml",
    # "发story拍照.xml",
    # "快拍设置不私.xml"
    # "注册第二次风控.xml"
    # "bio.xml"
    # "刷视频主页.xml",
    # "评论.xml"
    "分享视频链接.xml"
]

maps = []
for file in file_name:
    file_path = "./抓包/" + file
    header_maps1 = get_header_maps(file, get_xml(file_path))
    maps.append(header_maps1)

old = None
for i in maps:
    if old is None:
        old = i
        continue

    old = merge_header_list_map(i, old)

ppath, pmd5 = make_map(old)
ppath, md5s = remove_cookies(ppath, pmd5)

print(str(sorted(ppath, key=lambda x: x["path"], reverse=False)).replace("'", '"'))
print(str(sorted(pmd5, key=lambda x: x["header"], reverse=False)).replace("'", '"'))
print(str(md5s).replace("'", '"'))
print(str(header_values).replace("'", '"'))
