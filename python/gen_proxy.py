import json
import uuid

filename = "zone3_ips_数据中心.csv"
location = "us"
f = open("../data/" + filename)
ips = {}
for line in f.readlines():
    sp = line.strip().split(",")
    if len(sp) != 2:
        print("err:", line)
        continue

    if sp[1] == location:
        ip = {}
        ip['rip'] = sp[0]
        ip['country'] = location
        ip['ip'] = 'zproxy.lum-superproxy.io'
        ip['port'] = '22225'
        ip['username'] = 'lum-customer-hl_28871e6d-zone-zone3-ip-' + sp[0]
        ip['passwd'] = '13vps2jrphhn'
        ip['proxy_type'] = 0
        ip['need_auth'] = True
        ip['is_used'] = False
        ip['is_conn_error'] = False
        ip['is_risk'] = False
        ip['id'] = str(uuid.uuid4())
        if ips.get(ip['id']):
            print("你妈卖批!!!!")
        ips[ip['id']] = ip

print("all count", str(len(ips)))

f = open("../data/" + filename.replace(".csv", "") + "_" + location + ".json", "w")
d = json.dumps(ips)
json.dump(ips, f)
