# import pathlib
#
# import hpack
#
# path = r"C:/Users/Administrator/Desktop/"
# file_name = "无标题1"
# input_file = pathlib.Path(path + file_name)
# file_object = open(input_file, 'rb')
# data = file_object.read()
#
# for j in range(0, len(data)):
#     for i in range(0, len(data)):
#         try:
#             h = hpack.Decoder().decode(data[j:len(data) - i])
#             if len(h) != 0:
#                 print(j, i, h)
#             # print(hpack.Decoder().decode(data[i:]))
#         except Exception as e:
#             # print(i, e)
#             pass



# from urllib.parse import quote
#
# print(quote("{\n  \"Version3.1\" : {\n    \"iad-attribution\" : \"false\"\n  }\n}", 'utf-8'))


import gzip
import io
import os

read_file_name = r'C:\Users\Administrator\Desktop\1.gzip'
file_mode = 'rb'

with gzip.open(read_file_name, file_mode) as input_file:
    with io.TextIOWrapper(input_file, encoding='utf-8') as dec:
        print(dec.read())


