package main

import (
	"fmt"
	"makemoney/common"
	"time"
)

type Extra2 struct {
	WaterfallId string  `json:"waterfall_id"`
	StartTime   float32 `json:"start_time"`
	ElapsedTime float32 `json:"elapsed_time"`
	Step        string  `json:"step"`
	Flow        string  `json:"flow"`
}

type name struct {
	Extra2
}

func GenUploadID() string {
	upId := fmt.Sprintf("%d", time.Now().UnixMicro())
	return upId[:len(upId)-1]
}

func main() {
	var err error
	encodePasswd := common.InstagramQueryEscape("{\"quality_info\":\"{\\\"original_video_codec\\\":\\\"hvc1\\\",\\\"encoded_video_codec\\\":\\\"avc1\\\",\\\"original_width\\\":720,\\\"original_frame_rate\\\":30,\\\"encoded_bit_rate\\\":6270412,\\\"encoded_height\\\":1280,\\\"original_bit_rate\\\":690021.6875,\\\"encoded_color_primaries\\\":\\\"ITU_R_709_2\\\",\\\"original_height\\\":1280,\\\"encoded_frame_rate\\\":30,\\\"encoded_ycbcr_matrix\\\":\\\"ITU_R_709_2\\\",\\\"original_ycbcr_matrix\\\":\\\"ITU_R_709_2\\\",\\\"original_bits_per_component\\\":8,\\\"encoded_width\\\":720,\\\"encoded_transfer_function\\\":\\\"ITU_R_709_2\\\",\\\"measured_frames\\\":[{\\\"ssim\\\":0.98715811967849731,\\\"timestamp\\\":0},{\\\"ssim\\\":0.9876900315284729,\\\"timestamp\\\":0.83333333333333337},{\\\"ssim\\\":0.98564529418945312,\\\"timestamp\\\":1.6666666666666667},{\\\"ssim\\\":0.98104578256607056,\\\"timestamp\\\":2.5},{\\\"ssim\\\":0.97924661636352539,\\\"timestamp\\\":3.3333333333333335},{\\\"ssim\\\":0.98242074251174927,\\\"timestamp\\\":4.166666666666667},{\\\"ssim\\\":0.98820382356643677,\\\"timestamp\\\":5},{\\\"ssim\\\":0.98724061250686646,\\\"timestamp\\\":5.833333333333333},{\\\"ssim\\\":0.98652297258377075,\\\"timestamp\\\":6.666666666666667},{\\\"ssim\\\":0.98639136552810669,\\\"timestamp\\\":7.5},{\\\"ssim\\\":0.98449891805648804,\\\"timestamp\\\":8.3333333333333339},{\\\"ssim\\\":0.9821438193321228,\\\"timestamp\\\":9.1666666666666661},{\\\"ssim\\\":0.98663371801376343,\\\"timestamp\\\":10}]}\",\"_uuid\":\"DB51C2A4-B9F7-4F02-B0B1-9E5FE9CC0752\",\"_uid\":\"52027913737\",\"upload_id\":\"1647106610259307\"}")
	print(encodePasswd)
	encodePasswd = common.InstagramQueryEscape(encodePasswd)
	if err != nil {

	}
	print(encodePasswd)
	print(encodePasswd)
}
