import json
import os

# from androguard.misc import APK
#
dir = "C:/Users/Administrator/Desktop/apk/"
# for file in os.listdir(dir):
#     if ".apk" in file:
#         path = dir + file
#         os.system("adb install " + path)
#         input()
#         os.system("adb install " + path)

#
# for file in os.listdir(dir):
#     if ".apk" in file:
#         path = dir + file
#         try:
#             apk_obj = APK(path)
#             print(file, apk_obj.get_package(), "----------", apk_obj.get_app_name())
#         except Exception as e:
#             print(file, "error", "----------", e)
import shutil

#
# apks = []
# f = open(r"C:\Users\Administrator\Desktop\apk\apk包名.csv", "r", encoding="utf8")
# d = f.read()
# sp = d.split("\n")
# for s in sp:
#     ssp = s.split(",")
#     if len(ssp) == 3:
#         apk = {}
#         apk["file"] = ssp[0]
#         apk["pkg"] = ssp[1]
#         apk["name"] = ssp[2]
#         apks.append(apk)
#     else:
#         print(s)
#
# f = open(r"C:\Users\Administrator\Desktop\apk\雷电已有安装包.txt", "r", encoding="utf8")
# d = f.read()
# sp = d.split("\n")
# preset = set()
# for s in sp:
#     preset.add(s)
#
# f = open(r"C:\Users\Administrator\Desktop\apk\雷电安装后.txt", "r", encoding="utf8")
# d = f.read()
# sp = d.split("\n")
# aftset = list()
# for s in sp:
#     if s not in preset:
#         aftset.append(s)
#         # print(s)
#
# # w = open(r"C:\Users\Administrator\Desktop\apk\统计2.csv", "w", encoding="utf8")
#
# for apk in apks:
#     # w.write(apk["file"])
#     # w.write(",")
#     # w.write(apk["name"])
#     # w.write(",")
#     # w.write(apk["pkg"])
#     # w.write(",")
#     if apk["pkg"] not in aftset:
#         apk["install"] = "false"
#         # w.write("false")
#         # w.write(",")
#         print(apk["file"])
#         try:
#             shutil.move(dir + apk["file"], dir + "雷电安装失败/" + apk["file"])
#         except Exception as e:
#             pass
#
#     # else:
#     #     w.write(",")
#     # w.write("\n")
# print(json.dumps(apks, ensure_ascii=False))


app = """62267d5355944d21857dd79f.apk
6226928355944d21857dde60.apk
6226a50a55944d21857de425.apk
6226ad8155944d21857de6aa.apk
6226c56255944d21857de9fd.apk
62274ee655944d21857df1f6.apk
622751a455944d21857df246.apk
6227ab1d55944d21857df7c0.apk
6227cf0455944d21857dfa53.apk
6228930355944d21857e0373.apk
62293acc55944d21857e1132.apk
622bde502614c9550a60e49d.apk"""
sp = app.split("\n")
for i in sp:
    try:
        shutil.move(dir + "雷电安装失败/" + i, dir + "雷电安装失败/3/" + i)
    except Exception as e:
        pass