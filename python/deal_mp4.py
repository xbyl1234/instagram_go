import json
import os

import threading
import time

ffmpeg = "C:/Users/Administrator/Desktop/project/github/instagram_project/tools/video/ffmpeg.exe"


def conv_mp4_format(dir_path, file_name, ico_path, text):
    cmd = ffmpeg + " -i " + dir_path + "raw_video/" + file_name + \
          " -color_primaries bt709 " + \
          " -y " + dir_path + "video/" + file_name
    os.system(cmd)

    cmd = ffmpeg + " -i " + dir_path + "video/" + file_name + \
          " -i " + ico_path + " -filter_complex \"overlay=10:10\" " + \
          " -y " + dir_path + "video_ico/" + file_name
    print(cmd)
    os.system(cmd)

    cmd = ffmpeg + " -i " + dir_path + "video_ico/" + file_name + \
          " -vf drawtext='fontfile=comic.ttf:text=\"" + text + "\":x=(w-t*300):y=(h-th)/5:fontcolor=green:fontsize=100'" + \
          " -y " + dir_path + "video_ico_text/" + file_name
    print(cmd)
    os.system(cmd)


video_dir = "C:/Users/Administrator/Desktop/mn/刘二/"
ico_path = r"C:\Users\Administrator\Desktop\ico.png"
text = "insemail.work open it for more sexy"

try:
    os.makedirs(video_dir + "video")
except Exception as e:
    pass
try:
    os.makedirs(video_dir + "video_ico")
except Exception as e:
    pass

try:
    os.makedirs(video_dir + "video_ico_text")
except Exception as e:
    pass

mp4_infos = []
files = os.listdir(video_dir + "raw_video/")

index = 1
for file in files:
    print("---------", index)
    index += 1
    try:
        mp4_info = {}
        name = file.replace(".mp4", "")
        mp4_info["caption"] = name[:name.rfind("-")]
        mp4_info["file_name"] = name[name.rfind("-") + 1:]
        os.rename(video_dir + "raw_video/" + file,
                  video_dir + "raw_video/" + mp4_info["file_name"] + ".mp4")

        os.rename(video_dir + "cover/" + name + ".jpg",
                  video_dir + "cover/" + mp4_info["file_name"] + ".jpg")

        mp4_info["file_name"] = mp4_info["file_name"] + ".mp4"
        mp4_infos.append(mp4_info)
        conv_mp4_format(video_dir, mp4_info["file_name"], ico_path, text)
    except Exception as e:
        print(e)

print(json.dumps(mp4_infos))
