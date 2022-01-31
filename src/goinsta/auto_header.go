package goinsta

import (
	"makemoney/common"
	"makemoney/common/log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	IGHeader_EncryptionId           = "Password-Encryption-Key-Id"
	IGHeader_EncryptionKey          = "Password-Encryption-Pub-Key"
	IGHeader_Authorization          = "Authorization"
	IGHeader_udsUserID              = "Ig-U-Ds-User-Id"
	IGHeader_iguiggDirectRegionHint = "Ig-U-Ig-Direct-Region-Hint"
	IGHeader_iguShbid               = "Ig-U-Shbid"
	IGHeader_iguShbts               = "Ig-U-Shbts"
	IGHeader_iguRur                 = "Ig-U-Rur"
	IGHeader_UseAuthHeaderForSso    = "Use-Auth-Header-For-Sso"
	IGHeader_XMid                   = "X-Mid"
	IGHeader_igwwwClaim             = "X-Ig-Www-Claim"
)

func GetAutoHeaderFunc(header []string) []AutoSetHeaderFun {
	ret := make([]AutoSetHeaderFun, len(header))
	index := 0
	var serverHeader = func(key string) {
		ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
			value := inst.GetHeader(key)
			req.Header.Set(key, value)
			if value == "" {
				log.Warn("user: %s ignore header %s", inst.User, key)
			}
		}
		index++
	}
	for _, item := range header {
		switch item {
		case "Ig-U-Rur":
			serverHeader(item)
			break
		case "Ig-U-Ds-User-Id":
			serverHeader(item)
			break
		case "X-Mid":
			serverHeader(item)
			break
		case "Authorization":
			serverHeader(item)
			break
		case "Ig-U-Ig-Direct-Region-Hint":
			serverHeader(item)
			break
		case "X-Ig-Www-Claim":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				claim := inst.GetHeader("X-Ig-Www-Claim")
				if claim == "" {
					claim = "0"
				}
				req.Header.Set("X-Ig-Www-Claim", claim)
			}
			index++
			break
		case "Authorization-Others":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("Authorization-Others", "")
			}
			index++
			break

		case "X-Entity-Length":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
			}
			index++
			break
		case "X-Idfa":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Idfa", inst.Device.IDFA)
			}
			index++
			break
		case "X-Ig-App-Locale":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-App-Locale", inst.Device.AppLocale)
			}
			index++
			break

		case "X-Ig-Eu-Configure-Disabled":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Eu-Configure-Disabled", "true")
			}
			index++
			break
		case "X-Ig-Timezone-Offset":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Timezone-Offset", inst.Device.TimezoneOffset)
			}
			index++
			break

		case "X-Pigeon-Session-Id":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Pigeon-Session-Id", inst.sessionID)
			}
			index++
			break
		case "X-Fb-Http-Engine":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Fb-Http-Engine", "Liger")
			}
			index++
			break

		case "X-Ig-Connection-Speed":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Connection-Speed", InstagramReqSpeed[common.GenNumber(0, len(InstagramReqSpeed))])
			}
			index++
			break
		case "X-Ig-Bandwidth-Speed-Kbps":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Bandwidth-Speed-Kbps", "0.000")
			}
			index++
			break
		case "Connection":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("Connection", "close")
			}
			index++
			break
		case "X-Fb-Server-Cluster":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Fb-Server-Cluster", "True")
			}
			index++
			break

		case "X-Ig-App-Startup-Country":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-App-Startup-Country", inst.Device.StartupCountry)
			}
			index++
			break
		case "X-Pigeon-Rawclienttime":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Pigeon-Rawclienttime", strconv.FormatInt(time.Now().Unix(), 10)+".000000")
			}
			index++
			break
		case "Content-Type":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
			}
			index++
			break
		case "X-Bloks-Is-Panorama-Enabled":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Bloks-Is-Panorama-Enabled", "true")
			}
			index++
			break
		case "X-Fb-Client-Ip":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Fb-Client-Ip", "True")
			}
			index++
			break
		case "X-Ig-Device-Id":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Device-Id", inst.Device.DeviceID)
			}
			index++
			break
		case "Ig-Intended-User-Id":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("Ig-Intended-User-Id", strconv.FormatInt(inst.ID, 10))
			}
			index++
			break
		case "X-Ig-Capabilities":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Capabilities", "36r/Fx8=")
			}
			index++
			break
		case "Accept-Language":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("Accept-Language", inst.Device.AcceptLanguage)
			}
			index++
			break
		case "X-Ig-Connection-Type":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Connection-Type", inst.Device.NetWorkType)
			}
			index++
			break
		case "X-Bloks-Version-Id":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Bloks-Version-Id", inst.Device.BloksVersionID)
			}
			index++
			break

		case "X-Ig-Device-Locale":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Device-Locale", inst.Device.AppLocale)
			}
			index++
			break

		case "X-Device-Id":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Device-Id", inst.Device.DeviceID)
			}
			index++
			break
		case "X-Ig-Family-Device-Id":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Family-Device-Id", inst.Device.FamilyID)
			}
			index++
			break
		case "X-Ig-App-Id":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-App-Id", InstagramAppID)
			}
			index++
			break
		case "X-Ig-Mapped-Locale":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Mapped-Locale", strings.Replace(inst.Device.AppLocale, "-", "_", -1))
			}
			index++
			break
		case "Accept-Encoding":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("Accept-Encoding", "gzip, deflate")
			}
			index++
			break
		case "X-Ig-Abr-Connection-Speed-Kbps":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Ig-Abr-Connection-Speed-Kbps", "35")
			}
			index++
			break
		case "Family_device_id":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("Family_device_id", inst.Device.FamilyID)
			}
			index++
			break
		case "User-Agent":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("User-Agent", inst.Device.UserAgent)
			}
			index++
			break
		case "X-Tigon-Is-Retry":
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
				req.Header.Set("X-Tigon-Is-Retry", "False")
			}
			index++
			break

		//case "Media_hash":
		//	ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
		//	}
		//	index++
		//	break
		//case "X-Ig-Prefetch-Request":
		//	ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
		//	}
		//	index++
		//	break
		//case "X-Fb":
		//	ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
		//	}
		//	index++
		//	break
		//case "Content-Length":
		//	ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
		//	}
		//	index++
		//	break
		//case "X-Entity-Type":
		//	ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
		//	}
		//	index++
		//	break
		//case "X-Ads-Opt-Out":
		//	ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
		//	}
		//	index++
		//	break
		//case "X-Entity-Name":
		//	ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
		//	}
		//	index++
		//	break
		//case "Offset":
		//	ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
		//	}
		//	index++
		//	break
		//case "X_fb_photo_waterfall_id":
		//	ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
		//	}
		//	index++
		//	break
		//case "X-Instagram-Rupload-Params":
		//	ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
		//	}
		//	index++
		//	break
		default:
			if item != "Content-Length" {
				log.Info("header: %s fun is default!", item)
			}
			ret[index] = func(inst *Instagram, opt *reqOptions, req *http.Request) {
			}
			index++
		}
	}
	return ret
}
