from PIL import Image
import piexif, os, time


def setDate(photoAddress):
    mytime = str(time.strftime("%Y:%m:%d %H:%M:%S", time.localtime()))
    # 固定格式 mytime = "2013:12:19 10:10:10" ,不能用time.strftime("%Y/%m/%d %H:%M:%S",time.localtime())
    exif_ifd = {
        piexif.ExifIFD.DateTimeOriginal: mytime,
        piexif.ExifIFD.DateTimeDigitized: mytime,
        piexif.ExifIFD.CameraOwnerName: "ken"
    }
    exif_dict = {"Exif": exif_ifd}
    exif_bytes = piexif.dump(exif_dict)
    im = Image.open(photoAddress)
    im.save(photoAddress, exif=exif_bytes)


def setMyInfo(photoAddress):
    im = Image.open(photoAddress)
    exif_dict = piexif.load(im.info["exif"])
    exif_dict["0th"][piexif.ImageIFD.Artist] = "kenwanmao"
    exif_bytes = piexif.dump(exif_dict)
    # newfile = photoAddress
    im.save(photoAddress, exif=exif_bytes)


def clearExif(path):
    startTime = time.time()
    countNums = 0
    # os.walk() 方法用于通过在目录树中游走输出在目录中的文件名
    for root, dirs, files in os.walk(path):
        for name in files:
            if name.endswith(".JPG") or name.endswith(".jpg"):
                photoAddress = os.path.join(root, name)
                print("{},已经被抹去exif信息。".format(photoAddress))
                # 调用piexif库的remove函数直接去除exif信息。
                piexif.remove(photoAddress)
                countNums += 1
                # setDate(photoAddress)
                # setMyInfo(photoAddress)

    print("本次程序共清除{}张JPG照片的Exif信息，耗时{:.2f} s".format(countNums, time.time() - startTime))
    input("\n照片信息处理完毕，按任意键退出...")


if __name__ == "__main__":
    print("欢迎使用EXIF信息清除程序！\n使用规则如下：")
    print("1.可以将本程序放在图片目录下点开使用")
    print("2.将照片目录手动输入\n")
    # 获取当前程序所在的目录
    # nowDir = str(os.getcwd())
    # photoDir = input("手动输入照片目录:") or nowDir
    # 启动清除Exif信息函数

    clearExif(r"C:\Users\Administrator\Desktop\project\github\instagram_project\data\girl_picture")
