package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"makemoney/goinsta"
	"os"
)

type RawMediaBase struct {
	Latitude  float32
	Longitude float32
}

type RawImgMedia struct {
	RawMediaBase
	Image     image.Image
	ImageData []byte
	High      int
	Width     int
	Path      string `json:"path"`
	MD5       string
	Caption   string
	Loc       *goinsta.LocationSearch
}

func (this *RawImgMedia) LoadImage(path string) error {
	var err error
	this.ImageData, err = os.ReadFile(path)
	if err != nil {
		return err
	}

	imgWirte := bytes.NewBuffer(this.ImageData)
	img, err := jpeg.Decode(imgWirte)
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
	return this.ImageData
}

type RawVideoMedia struct {
	RawMediaBase
	VideoData   []byte
	ImageData   []byte
	High        int
	Width       int
	Duration    float64
	Path        string `json:"path"`
	MD5         string
	Caption     string
	AudioTitle  string
	FrameRate   float64
	BitRate     float64
	YcbcrMatrix string
}
