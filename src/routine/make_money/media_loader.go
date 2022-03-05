package main

import (
	"image"
	"image/jpeg"
	"makemoney/goinsta"
	"os"
)

type RawMediaBase struct {
}

type RawImgMedia struct {
	Image   image.Image
	High    int
	Width   int
	Path    string `json:"path"`
	MD5     string
	Caption string
	Loc     *goinsta.LocationSearch
}

func (this *RawImgMedia) LoadImage(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}

	this.Image = img
	this.Path = path

	bound := img.Bounds()
	this.Width = bound.Max.X
	this.High = bound.Max.Y
	return nil
}

func (this *RawImgMedia) GetImage() []byte {
	file, err := os.ReadFile(this.Path)
	if err != nil {
		return nil
	}
	return file

	//var buff bytes.Buffer
	//err := jpeg.Encode(&buff, this.Image, &jpeg.Options{Quality: 90})
	//if err != nil {
	//	return nil
	//}
	//b := buff.Bytes()
	//return buff.Bytes()
}

////测试读取视频
//window := gocv.NewWindow("Hello")
//window.ResizeWindow(960,540)
//img := gocv.NewMat()
//cap,_ := gocv.VideoCaptureFile("stopTest.avi")     <font >
//for {
//	cap.Read(&img)
//	window.IMShow(img)
//	window.WaitKey(1)
//}
