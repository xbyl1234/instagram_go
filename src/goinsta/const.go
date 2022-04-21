package goinsta

import (
	"net/http"
)

var (
	InstagramHost       = "https://i.instagram.com/"
	InstagramHost_B     = "https://b.i.instagram.com/"
	InstagramHost_Graph = "https://graph.instagram.com/"
	//InstagramUserAgent2  = "Instagram %s (iPhone7,2; iOS 12_5_5; en_US; en-US; scale=2.00; 750x1334; %s) AppleWebKit/420+"
	InstagramUserAgent = "Instagram %s (%s; iOS %s; en_US; en-US; %s; %s; %s) AppleWebKit/%s"
	//InstagramLocation    = "en_US"
	InstagramVersions    []*InstDeviceInfo
	InstagramAppID       = "124024574287414"
	InstagramAccessToken = "124024574287414|84a456d620314b6e92a16d8ff1c792dc"

	InstagramDeviceList = []string{"iPhone7,2 12_5_5 scale=2.00 750x1334"}

	//InstagramDeviceList = []string{
	//	"iPhone7,1 12_5_5 scale=2.61 1080x1920",
	//	//"iPhone7,2 12_5_5 scale=2.00 750x1334",
	//	//"iPhone7,2 11_2_2 scale=2.00 750x1334",
	//	"iPhone8,1 13_7 scale=2.00 750x1334",
	//	"iPhone8,1 14_6 scale=2.00 750x1334",
	//	"iPhone8,1 12_1_4 scale=2.34 750x1331",
	//	"iPhone8,1 14_7_1 scale=2.00 750x1334",
	//	"iPhone9,1 15_1 scale=2.00 750x1334",
	//	"iPhone9,3 15_1 scale=2.00 750x1334",
	//	"iPhone9,3 14_8_1 scale=2.00 750x1334",
	//	"iPhone9,3 14_1 scale=2.00 750x1334",
	//	"iPhone9,4 14_6 scale=2.61 1080x1920",
	//	"iPhone9,4 15_1 scale=2.61 1080x1920",
	//	"iPhone9,4 14_8_1 scale=2.61 1080x1920",
	//	"iPhone10,1 15_2 scale=2.00 750x1334",
	//	"iPhone10,1 14_8_1 scale=2.00 750x1334",
	//	//"iPhone10,3 15_1 scale=3.00 1125x2436",
	//	//"iPhone10,4 15_0 scale=2.00 750x1334",
	//	//"iPhone10,4 14_8_1 scale=2.00 750x1334",
	//	//"iPhone10,4 15_1 scale=2.00 750x1334",
	//	//"iPhone10,4 14_7_1 scale=2.00 750x1334",
	//	//"iPhone10,4 15_0_2 scale=2.00 750x1334",
	//	//"iPhone10,5 15_1 scale=2.88 1080x1920",
	//	//"iPhone10,5 15_0_1 scale=2.61 1080x1920",
	//	//"iPhone10,5 13_5_1 scale=2.61 1080x1920",
	//	//"iPhone10,6 15_2 scale=3.00 1125x2436",
	//	//"iPhone10,6 15_1 scale=3.00 1125x2436",
	//	"iPhone11,2 14_7_1 scale=3.00 1125x2436",
	//	"iPhone11,2 15_1 scale=3.00 1125x2436",
	//	"iPhone11,2 14_6 scale=3.00 1125x2436",
	//	//"iPhone11,6 15_1 scale=3.00 1242x2688",
	//	//"iPhone11,6 14_8_1 scale=3.00 1242x2688",
	//	//"iPhone11,8 14_8_1 scale=2.00 828x1792",
	//	//"iPhone11,8 14_6 scale=2.00 828x1792",
	//	//"iPhone11,8 15_0_2 scale=2.00 828x1792",
	//	//"iPhone11,8 15_2 scale=2.00 828x1792",
	//	//"iPhone11,8 15_1 scale=2.00 828x1792",
	//	"iPhone12,1 14_7_1 scale=2.00 828x1792",
	//	"iPhone12,1 14_8_1 scale=2.00 828x1792",
	//	"iPhone12,1 15_2 scale=2.00 828x1792",
	//	"iPhone12,1 15_0 scale=2.00 828x1792",
	//	"iPhone12,1 13_7 scale=2.00 828x1792",
	//	"iPhone12,1 15_0_2 scale=2.00 828x1792",
	//	"iPhone12,3 15_0_2 scale=3.00 1125x2436",
	//	"iPhone12,3 14_6 scale=3.00 1125x2436",
	//	"iPhone12,3 14_2 scale=3.00 1125x2436",
	//	"iPhone12,3 15_1 scale=3.00 1125x2436",
	//	//"iPhone12,5 14_3 scale=3.00 1242x2688",
	//	//"iPhone12,5 15_1 scale=3.31 1242x2689",
	//	//"iPhone12,5 14_6 scale=3.00 1242x2688",
	//	//"iPhone12,5 13_3_1 scale=3.31 1242x2689",
	//	//"iPhone12,5 14_7_1 scale=3.00 1242x2688",
	//	//"iPhone12,5 14_8_1 scale=3.00 1242x2688",
	//	"iPhone12,8 15_2 scale=2.00 750x1334",
	//	"iPhone12,8 15_1 scale=2.00 750x1334",
	//	"iPhone12,8 14_8_1 scale=2.00 750x1334",
	//	"iPhone12,8 14_7_1 scale=2.00 750x1334",
	//	"iPhone13,1 14_7_1 scale=2.88 1080x2338",
	//	"iPhone13,2 15_2 scale=3.00 1170x2532",
	//	"iPhone13,2 15_0_2 scale=3.00 1170x2532",
	//	"iPhone13,2 15_1_1 scale=3.00 1170x2532",
	//	"iPhone13,2 14_6 scale=3.00 1170x2532",
	//	"iPhone13,2 14_7_1 scale=3.00 1170x2532",
	//	//"iPhone13,3 15_1_1 scale=3.00 1170x2532",
	//	//"iPhone13,3 15_2 scale=3.00 1170x2532",
	//	//"iPhone13,3 14_8 scale=3.00 1170x2532",
	//	//"iPhone13,3 14_8_1 scale=3.00 1170x2532",
	//	"iPhone13,4 15_1_1 scale=3.00 1284x2778",
	//	"iPhone13,4 15_2 scale=3.00 1284x2778",
	//	//"iPhone14,2 15_1_1 scale=3.00 1170x2532",
	//	//"iPhone14,3 15_2 scale=3.00 1284x2778",
	//	//"iPhone14,5 15_0 scale=3.00 1170x2532",
	//	//"iPhone14,5 15_2 scale=3.00 1170x2532",
	//	//"iPhone14,5 15_0 scale=3.66 1170x2533",
	//	//"iPhone14,5 15_1_1 scale=3.00 1170x2532",
	//}

	//InstagramDeviceList = []string{
	//	"iPhone7,2 12_5_5 scale=2.00 750x1334",
	//}
	LensModel = map[string]string{
		"7,2": "iPhone,6,back camera,4.15mm,f/2.2",
		//"7,1":  "iPhone,6,back camera,4.15mm,f/2.2",
		//"8,1":  "iPhone,6s,back camera,4.15mm,f/2.2",
		//"9,1":  "iPhone,7,back camera,3.99mm,f/1.8",
		//"9,2":  "iPhone,7,Plus back dual camera,3.99mm,f/1.8",
		//"9,3":  "iPhone,7,back camera,3.99mm,f/1.8",
		//"9,4":  "iPhone,7,Plus back dual camera,3.99mm,f/1.8",
		//"10,1": "iPhone,8,back camera,3.99mm,f/1.8",
		//"11,2": "iPhone,XS,Max back dual camera,4.25mm,f/1.8",
		//"12,1": "iPhone,11,back dual wide camera,4.25mm,f/1.8",
		//"12,3": "iPhone,11,Pro back triple camera,4.25mm,f/1.8",
		//"12,8": "iPhone,SE,(2nd generation) back camera,3.99mm,f/1.8",
		//"13,1": "iPhone,12,mini back dual wide camera,4.2mm,f/1.6",
		//"13,2": "iPhone,12,back dual wide camera,4.2mm,f/1.6",
		//"13,4": "iPhone,12,Pro Max back triple camera,5.1mm,f/1.6",
	}

	InstagramVersionData = []string{
		"190.0.0.26.119 294609445 5538d18f11cad2fa88efa530f8a717c5d5339d1d53fc5140af9125216d1f7a89",
		//"191.0.0.25.122 296543649 bf3e79f2304601044c85a6f9c44dab59a72558ca9f9a821b96882a4a54ca3c3a",
		//"192.0.0.37.119 298025452 9fd8ac08308424f3385019b6c63fc3eb52f3d9d1314f33b78b5db21716d3bf7a",
		//"193.0.0.29.121 299401192 02872c8277b5f0ccc5275f61be6e600aeda09ad6927f75ed89e4177ba297bd4f",
		//"194.0.0.33.171 301014278 f487d54d4da25bc844e5abc7af030f30a3f83d27bf24d8765274c62268d17dab",
		//"195.0.0.24.119 302211069 64f0fd5651d93707c9c3f98da6e5af03e46140770801f4d6eb3c2799911ec21f",
		//"196.0.0.21.120 303649428 d5507f7c0ad817ba90a28666727e7cfbe294b33203ab934be3e68e8823c27605",
		//"197.0.0.20.119 305020938 29d3248efc1cfc10a0dbafd84eb58fd7eebc6b99e626c5bfa0fe615b8ff784d9",
		//"198.0.0.27.119 306495430 4e3e06f8c5ab8a9ab19a536615e0cd79a8eb4742068f2f63af62a4511a69187c",
		//"199.0.0.27.120 307960803 4591eb23633140d3162279972d8981f471b81a93fe8dcc4f28d7c24014470eee",
		//"200.0.0.19.118 309467717 c8855ee853644de5034628e4716bd7c915e25b17ba524f999d904099b41e910d",
		//"200.1.0.20.118 310109037 c8855ee853644de5034628e4716bd7c915e25b17ba524f999d904099b41e910d",
		//"202.0.0.23.119 312612729 212380a0da87f76cbe923f35a5e31cd753ec95b4d6d584729e8054420f7822ca",
		//"203.0.0.26.117 314122909 0d09d64dfdf69cb3c47fd6f30f0948fbdb48a267aff5aaf8e450a5122dd0a68f",
		//"204.0.0.16.119 315460786 f646f3986f57a81c3eef6b8fde357c324737ea272a2f59045ad73fd2275519e5",
		//"205.0.0.20.115 317250287 d28e6e26fb0da1c17e9f33a6165eb94acc93bea287ea880be847d1a6796acbb9",
		//"206.0.0.30.118 318760365 965488f430b5bb53716f1a026db174f1c0a069c657b3ff9cc11afa555acd2797",
		//"206.1.0.30.118 320107871 965488f430b5bb53716f1a026db174f1c0a069c657b3ff9cc11afa555acd2797",
		//"207.0.0.28.118 320361397 3ced1a8d80642ea42beaafe4d2769cd5afd3b763d091b0eddb60f30850562184",
		//"208.0.0.26.131 322184766 c9e961cf88c88c69f97f01c5d4342d299bdb880c11894424c622d5b97e8a6a76",
		//"209.0.0.23.113 324180638 a94a49452ea2139f4109a16413cf473fe4bbe021b27e5f6e86625923ff3969d9",
		//"209.1.0.25.113 324511477 a94a49452ea2139f4109a16413cf473fe4bbe021b27e5f6e86625923ff3969d9",
		//"210.0.0.16.67 325544617 857f81be49bc66e64a152220da016951113eec808e531a9e235d55e342cf43c4",
		//"211.0.0.21.118 327311214 93202078f561251d153a675632166c73377a4c41e127ace6f25ee9484a77c7ce",
		//"212.0.0.22.118 328988229 0d38efe9f67cf51962782e8aae19001881099884d8d86c683d374fc1b89ffad1",
		//"212.1.0.25.118 329643252 0d38efe9f67cf51962782e8aae19001881099884d8d86c683d374fc1b89ffad1",
		//"213.0.0.19.117 330663239 5bf3152c14a8e8651b2ec5a689994b294f4e0a74b86b5652da331aa7035d1c62",
		//"213.1.0.22.117 332048479 5bf3152c14a8e8651b2ec5a689994b294f4e0a74b86b5652da331aa7035d1c62",
	}

	InstagramReqSpeed = []string{
		"35kbps",
		"29kbps",
		"42kbps",
		"45kbps",
		"58kbps",
		"78kbps",
		"9kbps",
		"46kbps",
	}

	CoordMap = map[string]InstLocationInfo{
		//"纽约": {
		//	Country:        "美国",
		//	City:           "纽约",
		//	Lon:            -73.87141290786828,
		//	Lat:            40.8385895611293,
		//	Timezone:       "-18000",
		//	AppLocale:      "en-US",
		//	StartupCountry: "US",
		//	MappedLocale:   "en_US",
		//	AcceptLanguage: "en-US;q=1.0",
		//},
		"美国英文": {
			Country:        "美国",
			City:           "纽约",
			Lon:            -73.87141290786828,
			Lat:            40.8385895611293,
			Timezone:       "28800",
			AppLocale:      "en-US",
			StartupCountry: "US",
			MappedLocale:   "en_US",
			AcceptLanguage: "en-US;q=1.0",
		},
		//"纽约中文2": {
		//	Country:        "美国",
		//	City:           "纽约",
		//	Lon:            -73.87141290786828,
		//	Lat:            40.8385895611293,
		//	Timezone:       "-18000",
		//	AppLocale:      "zh-Hans-US",
		//	StartupCountry: "US",
		//	MappedLocale:   "en_US",
		//	AcceptLanguage: "zh-CN;q=1.0",
		//},
		//"纽约中文3": {
		//	Country:        "美国",
		//	City:           "纽约",
		//	Lon:            -73.87141290786828,
		//	Lat:            40.8385895611293,
		//	Timezone:       "-18000",
		//	AppLocale:      "zh-Hans-US",
		//	StartupCountry: "US",
		//	MappedLocale:   "zh_CN",
		//	AcceptLanguage: "zh-CN;q=1.0",
		//},
		//"纽约中文4": {
		//	Country:        "美国",
		//	City:           "纽约",
		//	Lon:            -73.87141290786828,
		//	Lat:            40.8385895611293,
		//	Timezone:       "-18000",
		//	AppLocale:      "zh-Hans-US",
		//	StartupCountry: "JP",
		//	MappedLocale:   "zh_CN",
		//	AcceptLanguage: "zh-CN;q=1.0",
		//},
	}

	NoLoginHeaderMap map[string]*HeaderSequence
	LoginHeaderMap   map[string]*HeaderSequence
	HeaderMD5Map     map[string]*HeaderSequence
	ReqHeaderJson    reqHeaderJson
)

type AutoSetHeaderFun func(inst *Instagram, opt *reqOptions, req *http.Request)

type HeaderSequence struct {
	HeaderFun []AutoSetHeaderFun
	HeaderSeq []string
}

type pathsMap struct {
	Path string `json:"path"`
	Md5  string `json:"md5"`
}

type reqHeaderJson struct {
	PathsNoLogin []pathsMap `json:"paths_no_login"`
	PathsLogin   []pathsMap `json:"paths_login"`
	Md5S         []*struct {
		Desp      string `json:"desp,omitempty"`
		Md5       string `json:"md5"`
		Header    string `json:"header"`
		headerSeq HeaderSequence
	} `json:"md5s"`
}

type muteOption string

const (
	MuteAll   muteOption = "all"
	MuteStory muteOption = "story"
	MuteFeed  muteOption = "feed"
)

// Endpoints (with format vars)
