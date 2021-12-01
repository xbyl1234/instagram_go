import json
import struct

import zstandard
import pathlib
import shutil

path = "C:/Users/Administrator/Desktop/project/github/instagram_project/抓包/"
file_name = "s1"
input_file = pathlib.Path(path + file_name)

file_object = open(input_file, 'rb')
decomp = zstandard.ZstdDecompressor()
jso = json.loads(decomp.decompress(file_object.read()))
data = json.dumps(jso, sort_keys=True, indent=4)

f = open(path + file_name + ".json", "w")
f.write(data)
f.close()
print(data)
