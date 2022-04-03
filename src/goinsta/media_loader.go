package goinsta

import (
	"bytes"
	"encoding/json"
	"github.com/edgeware/mp4ff/mp4"
	"image"
	"image/jpeg"
	"makemoney/common"
	"os"
	"os/exec"
	"strconv"
	"strings"
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
	Loc       *LocationSearch
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
	VideoCodec  string
	ImageData   []byte
	ImagePath   string
	ImageMD5    string
	From        string
	Waterfall   string
	UploadId    string
	MDatStart   int
	MDatLen     int
}

const FromCamera = "camera"
const FromLibrary = "library"

type FfStream struct {
	Index            int    `json:"index"`
	CodecName        string `json:"codec_name"`
	CodecLongName    string `json:"codec_long_name"`
	Profile          string `json:"profile"`
	CodecType        string `json:"codec_type"`
	CodecTagString   string `json:"codec_tag_string"`
	CodecTag         string `json:"codec_tag"`
	Width            int    `json:"width,omitempty"`
	Height           int    `json:"height,omitempty"`
	CodedWidth       int    `json:"coded_width,omitempty"`
	CodedHeight      int    `json:"coded_height,omitempty"`
	ClosedCaptions   int    `json:"closed_captions,omitempty"`
	FilmGrain        int    `json:"film_grain,omitempty"`
	HasBFrames       int    `json:"has_b_frames,omitempty"`
	PixFmt           string `json:"pix_fmt,omitempty"`
	Level            int    `json:"level,omitempty"`
	ColorRange       string `json:"color_range,omitempty"`
	ColorSpace       string `json:"color_space,omitempty"`
	ColorTransfer    string `json:"color_transfer,omitempty"`
	ColorPrimaries   string `json:"color_primaries,omitempty"`
	ChromaLocation   string `json:"chroma_location,omitempty"`
	FieldOrder       string `json:"field_order,omitempty"`
	Refs             int    `json:"refs,omitempty"`
	IsAvc            string `json:"is_avc,omitempty"`
	NalLengthSize    string `json:"nal_length_size,omitempty"`
	Id               string `json:"id"`
	RFrameRate       string `json:"r_frame_rate"`
	AvgFrameRate     string `json:"avg_frame_rate"`
	TimeBase         string `json:"time_base"`
	StartPts         int    `json:"start_pts"`
	StartTime        string `json:"start_time"`
	DurationTs       int    `json:"duration_ts"`
	Duration         string `json:"duration"`
	BitRate          string `json:"bit_rate"`
	BitsPerRawSample string `json:"bits_per_raw_sample,omitempty"`
	NbFrames         string `json:"nb_frames"`
	ExtradataSize    int    `json:"extradata_size"`
	Disposition      struct {
		Default         int `json:"default"`
		Dub             int `json:"dub"`
		Original        int `json:"original"`
		Comment         int `json:"comment"`
		Lyrics          int `json:"lyrics"`
		Karaoke         int `json:"karaoke"`
		Forced          int `json:"forced"`
		HearingImpaired int `json:"hearing_impaired"`
		VisualImpaired  int `json:"visual_impaired"`
		CleanEffects    int `json:"clean_effects"`
		AttachedPic     int `json:"attached_pic"`
		TimedThumbnails int `json:"timed_thumbnails"`
		Captions        int `json:"captions"`
		Descriptions    int `json:"descriptions"`
		Metadata        int `json:"metadata"`
		Dependent       int `json:"dependent"`
		StillImage      int `json:"still_image"`
	} `json:"disposition"`
	Tags struct {
		Language    string `json:"language"`
		HandlerName string `json:"handler_name"`
		VendorId    string `json:"vendor_id"`
	} `json:"tags"`
	SampleFmt     string `json:"sample_fmt,omitempty"`
	SampleRate    string `json:"sample_rate,omitempty"`
	Channels      int    `json:"channels,omitempty"`
	ChannelLayout string `json:"channel_layout,omitempty"`
	BitsPerSample int    `json:"bits_per_sample,omitempty"`
}

type Ffprobe struct {
	Streams []*FfStream `json:"streams"`
}

func (this *RawVideoMedia) LoadVideo(videoPath string, imgPath string) error {
	this.ImagePath = imgPath
	this.Path = videoPath

	cmd := exec.Command(".\\tools\\video\\ffprobe", "-v", "quiet", "-print_format", "json", "-show_streams", videoPath)
	result, err := cmd.Output()

	ffprobe := &Ffprobe{}
	err = json.Unmarshal(result, ffprobe)
	if err != nil {
		return err
	}
	var video *FfStream
	for _, item := range ffprobe.Streams {
		if item.CodecType == "video" {
			video = item
		}
	}
	if video == nil {
		return &common.MakeMoneyError{ErrStr: "not find video stream!"}
	}

	this.High = video.Height
	this.Width = video.Width
	this.VideoCodec = video.CodecTagString
	this.BitRate, _ = strconv.ParseFloat(video.BitRate, 64)

	sp := strings.Split(video.AvgFrameRate, "/")
	f1, _ := strconv.ParseFloat(sp[0], 64)
	f2, _ := strconv.ParseFloat(sp[1], 64)
	this.FrameRate = f1 / f2

	times, _ := strconv.ParseFloat(video.Duration, 64)
	this.Duration = times * 1000

	this.YcbcrMatrix = "ITU_R_709_2"

	this.VideoData, _ = os.ReadFile(videoPath)
	this.ImageData, _ = os.ReadFile(imgPath)

	mp4File, err := mp4.ReadMP4File(videoPath)
	if err != nil {
		return err
	}

	this.MDatStart = int(mp4File.Mdat.StartPos)
	this.MDatLen = int(mp4File.Mdat.DataLength())

	return nil
}

func (this *RawVideoMedia) CopyAndModifyMd5() *RawVideoMedia {
	n := &RawVideoMedia{
		RawMediaBase: RawMediaBase{},
		High:         this.High,
		Width:        this.Width,
		Duration:     this.Duration,
		Path:         this.Path,
		MD5:          this.MD5,
		Caption:      this.Caption,
		AudioTitle:   this.AudioTitle,
		FrameRate:    this.FrameRate,
		BitRate:      this.BitRate,
		YcbcrMatrix:  this.YcbcrMatrix,
		VideoCodec:   this.VideoCodec,
		ImagePath:    this.ImagePath,
		ImageMD5:     this.ImageMD5,
		From:         this.From,
		Waterfall:    this.Waterfall,
		UploadId:     this.UploadId,
		MDatStart:    this.MDatStart,
		MDatLen:      this.MDatLen,
	}

	copy(n.VideoData, this.VideoData)
	copy(n.ImageData, this.ImageData)

	for count := 0; count < 5; count++ {
		n.VideoData[n.MDatStart+count] = byte(common.GenNumber(0, 255))
	}
	return n
}
