import json
import struct

import zstandard
import pathlib
import shutil

# https://i.instagram.com/challenge/?next=/api/v1/users/50333860614/info/%253Ffrom_module%253Dself_profile

# path = "C:/Users/Administrator/Desktop/project/github/instagram_project/抓包/"
path = r"C:/Users/Administrator/Desktop/"
file_name = "无标题1"
input_file = pathlib.Path(path + file_name)

file_object = open(input_file, 'rb')
data = file_object.read()
decodedata = None

for i in range(0, len(data)):
    try:
        decomp = zstandard.ZstdDecompressor()
        # decodedata = decomp.decompress((data[:len(data) - i]))
        decodedata = decomp.decompress((data[i:]))
        print(i, decodedata)
        # print(hpack.Decoder().decode(data[i:]))
    except Exception as e:
        # print(i, e)
        pass

jso = json.loads(decodedata)
data = json.dumps(jso, sort_keys=True, indent=4)

f = open(path + file_name + ".json", "w")
f.write(data)
f.close()
print(data)
