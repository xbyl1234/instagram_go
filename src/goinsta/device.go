package goinsta

import (
	"fmt"
	"makemoney/common"
	"makemoney/common/log"
	"strconv"
	"strings"
)

type InstRegisterInfo struct {
	RegisterEmail       string `json:"register_email" bson:"register_email"`
	RegisterPhoneNumber string `json:"register_phone_number" bson:"register_phone_number"`
	RegisterPhoneArea   string `json:"register_phone_area" bson:"register_phone_area"`
	RegisterIpCountry   string `json:"register_ip_country" bson:"register_ip_country"`
	RegisterTime        int64  `json:"register_time" bson:"register_time"`
}

type InstLocationInfo struct {
	Country        string  `json:"country" bson:"country"`
	City           string  `json:"city" bson:"city"`
	Lat            float64 `json:"lat" bson:"lat"`
	Lon            float64 `json:"lon" bson:"lon"`
	Timezone       string  `json:"timezone" bson:"timezone"`
	AppLocale      string  `json:"app_locale" bson:"app_locale"`
	StartupCountry string  `json:"startup_country" bson:"startup_country"`
	AcceptLanguage string  `json:"accept_language" bson:"accept_language"`
}

type InstDeviceInfo struct {
	Version        string  `json:"version" bson:"version"`
	VersionCode    string  `json:"version_code" bson:"version_code"`
	BloksVersionID string  `json:"bloks_version_id" bson:"bloks_version_id"`
	UserAgent      string  `json:"user_agent" bson:"user_agent"`
	IDFA           string  `json:"idfa" bson:"idfa"`
	NetWorkType    string  `json:"net_work_type" bson:"net_work_type"`
	DeviceID       string  `json:"device_id" bson:"device_id"`
	FamilyID       string  `json:"family_id" bson:"family_id"`
	WaterID        string  `json:"water_id" bson:"water_id"`
	DeviceToken    string  `json:"device_token" bson:"device_token"`
	SystemVersion  string  `json:"system_version" bson:"system_version"`
	LensModel      string  `json:"lens_model" bson:"lens_model"`
	FocalLength    float64 `json:"focal_length" bson:"focal_length"`
	Aperture       float64 `json:"aperture" bson:"aperture"`
}

type InstAccountInfo struct {
	Device   InstDeviceInfo   `json:"device" bson:"device"`
	Location InstLocationInfo `json:"location" bson:"location"`
	Register InstRegisterInfo `json:"register" bson:"register"`
}

func InitInstagramConst() error {
	InstagramVersions = make([]*InstDeviceInfo, len(InstagramVersionData))
	for index, item := range InstagramVersionData {
		sp := strings.Split(item, " ")
		InstagramVersions[index] = &InstDeviceInfo{
			Version:        sp[0],
			VersionCode:    sp[1],
			BloksVersionID: sp[2],
		}
	}

	err := common.LoadJsonFile("config/http_header_sequence.json", &ReqHeaderJson)
	if err != nil {
		return err
	}
	HeaderMD5Map = make(map[string]*HeaderSequence)
	for _, md5 := range ReqHeaderJson.Md5S {
		sp := strings.Split(md5.Header, ",")
		if len(sp) == 0 {
			return &common.MakeMoneyError{
				ErrStr: "header is null,md5: " + md5.Md5,
			}
		}
		sp = sp[:len(sp)-1]
		md5.headerSeq.HeaderFun = GetAutoHeaderFunc(sp)
		md5.headerSeq.HeaderSeq = sp
		HeaderMD5Map[md5.Md5] = &md5.headerSeq
	}

	var makeSeqMap = func(headerMap *map[string]*HeaderSequence, pahts []pathsMap) {
		*headerMap = make(map[string]*HeaderSequence)
		for _, path := range pahts {
			find := false
			for _, md5 := range ReqHeaderJson.Md5S {
				if path.Md5 == md5.Md5 {
					(*headerMap)[path.Path] = &md5.headerSeq
					find = true
					break
				}
			}
			if !find {
				log.Warn("not find header md5: %s", path.Md5)
			}
		}
	}

	makeSeqMap(&NoLoginHeaderMap, ReqHeaderJson.PathsNoLogin)
	makeSeqMap(&LoginHeaderMap, ReqHeaderJson.PathsLogin)
	return err
}

func GenInstDeviceInfo() *InstAccountInfo {
	version := InstagramVersions[common.GenNumber(0, len(InstagramVersions))]
	device := InstagramDeviceList[common.GenNumber(0, len(InstagramDeviceList))]
	sp := strings.Split(device, " ")
	SystemVersion := strings.ReplaceAll(sp[1], "_", ".")

	spLens := strings.Split(LensModel[sp[0][6:]], ",")
	Lens := ""
	for _, item := range spLens {
		Lens += item + " "
	}
	Lens = Lens[:len(Lens)-1]
	FocalLength, _ := strconv.ParseFloat(strings.ReplaceAll(spLens[4], "f/", ""), 64)
	Aperture, _ := strconv.ParseFloat(strings.ReplaceAll(spLens[3], "mm", ""), 64)

	coord := CoordMap["纽约"]

	instVersion := &InstAccountInfo{
		Device: InstDeviceInfo{
			IDFA:           strings.ToUpper(common.GenUUID()),
			Version:        version.Version,
			VersionCode:    version.VersionCode,
			BloksVersionID: version.BloksVersionID,
			UserAgent:      fmt.Sprintf(InstagramUserAgent, version.Version, sp[0], sp[1], sp[2], sp[3], version.VersionCode, "420+"),
			NetWorkType:    "WiFi",
			DeviceID:       strings.ToUpper(common.GenUUID()),
			WaterID:        common.GenString(common.CharSet_16_Num, 32),
			SystemVersion:  SystemVersion,
			LensModel:      Lens,
			FocalLength:    FocalLength,
			Aperture:       Aperture,
		},
		Location: coord,
		Register: InstRegisterInfo{},
	}

	instVersion.Device.FamilyID = instVersion.Device.DeviceID
	return instVersion
}
