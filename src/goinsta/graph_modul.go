package goinsta

import (
	"fmt"
	"reflect"
)

type ExtraPk struct {
	Pk string `json:"pk"`
}

type ExtraRadioType struct {
	RadioType string `json:"radio_type"`
}

type ExtraContainerModule struct {
	ContainerModule string `json:"containermodule"`
}

type Extra2 struct {
	WaterfallId string  `json:"waterfall_id"`
	StartTime   float32 `json:"start_time"`
	ElapsedTime float32 `json:"elapsed_time"`
	Step        string  `json:"step"`
	Flow        string  `json:"flow"`
}

type ExtraContactPointType struct {
	ContactPointType string `json:"contact_point_type"`
}

type ExtraSomeId struct {
	FbFamilyDeviceId string `json:"fb_family_device_id"`
	AppDeviceId      string `json:"app_device_id"`
}

func initBaseExtra(inst *Instagram, extra interface{}) {
	extraValue := reflect.ValueOf(extra)

	field := extraValue.FieldByName("Pk")
	if field.CanSet() {
		field.SetString(fmt.Sprintf("%s", inst.ID))
	}

	field = extraValue.FieldByName("radio_type")
	if field.CanSet() {
		field.SetString("wifi-none")
	}

	field = extraValue.FieldByName("containermodule")
	if field.CanSet() {
		field.SetString("waterfall_log_in")
	}

	field = extraValue.FieldByName("waterfall_id")
	if field.CanSet() {
		field.SetString("")
	}
	field = extraValue.FieldByName("start_time")
	if field.CanSet() {
		field.SetString("")
	}

	field = extraValue.FieldByName("elapsed_time")
	if field.CanSet() {
		field.SetString("")
	}

	field = extraValue.FieldByName("step")
	if field.CanSet() {
		field.SetString("")
	}

	field = extraValue.FieldByName("flow")
	if field.CanSet() {
		field.SetString("phone")
	}

	field = extraValue.FieldByName("contact_point_type")
	if field.CanSet() {
		field.SetString("phone")
	}

	field = extraValue.FieldByName("fb_family_device_id")
	if field.CanSet() {
		field.SetString("")
	}

	field = extraValue.FieldByName("app_device_id")
	if field.CanSet() {
		field.SetString("")
	}
}

type nextButtonTapped struct {
	ExtraPk
	ExtraRadioType
	ExtraContainerModule
	Extra2
	ExtraContactPointType
	PhoneNumber string `json:"phone_number"`
}

func MakeNextButtonTapped(inst *Instagram, params map[string]interface{}) *GraphData {
	var ret = &GraphData{}
	ret.Module = "app"
	ret.Name = "next_button_tapped"
	ret.Extra = &nextButtonTapped{
		PhoneNumber: inst.RegisterPhoneArea + inst.RegisterPhoneNumber,
	}

	initBaseExtra(inst, ret.Extra)
	return ret
}

type regFieldInteracted struct {
	ExtraPk
	ExtraRadioType
	Extra2
	ExtraContainerModule
	FieldName       string  `json:"field_name"`
	InteractionType string  `json:"interaction_type"`
	CurrentTime     float32 `json:"current_time"`
	Guid            string  `json:"guid"`
}

func MakeRegFieldInteracted(inst *Instagram, params map[string]interface{}) *GraphData {
	var ret = &GraphData{}
	ret.Module = "waterfall_log_in"
	ret.Name = "reg_field_interacted"
	ret.Extra = &regFieldInteracted{}
	initBaseExtra(inst, ret.Extra)
	return ret
}

type fxSsoLibrary struct {
	ExtraRadioType
	ExtraPk
	VersionId                      string `json:"version_id"`
	TargetAccountType              int    `json:"target_account_type"`
	FxSsoLibraryCredentialSource   string `json:"fx_sso_library_credential_source"`
	FxSsoLibraryFlowUsingAuthToken string `json:"fx_sso_library_flow_using_auth_token"`
	FxSsoLibraryEvent              string `json:"fx_sso_library_event"`
	LogLocation                    string `json:"log_location"`
}

func MakeFxSsoLibrary(inst *Instagram, params map[string]interface{}) *GraphData {
	var ret = &GraphData{}
	ret.Module = "app"
	ret.Name = "fx_sso_library"
	ret.Extra = &fxSsoLibrary{}
	initBaseExtra(inst, ret.Extra)
	return ret
}

type analyticsFileDeleted struct {
	ExtraRadioType
	ExtraPk
	Channel string `json:"channel"`
}

func MakeAnalyticsFileDeleted(inst *Instagram, params map[string]interface{}) *GraphData {
	var ret = &GraphData{}
	ret.Module = "app"
	ret.Name = "analytics_file_deleted"
	ret.Extra = &analyticsFileDeleted{}
	initBaseExtra(inst, ret.Extra)
	return ret
}

type nextBlocked struct {
	ExtraPk
	ExtraRadioType
	Extra2
	ExtraSomeId
	ExtraContactPointType
	Reason string `json:"reason"`
}

func MakeNextBlocked(inst *Instagram, params map[string]interface{}) *GraphData {
	var ret = &GraphData{}
	ret.Module = "waterfall_log_in"
	ret.Name = "next_blocked"
	ret.Extra = &nextBlocked{}
	initBaseExtra(inst, ret.Extra)
	return ret
}

type timeSpentBitArray struct {
	TosPersistentUptime int     `json:"tos_persistent_uptime"`
	TosArray            []int   `json:"tos_array"`
	TosSeq              int     `json:"tos_seq"`
	TosLen              int     `json:"tos_len"`
	TosCum              int     `json:"tos_cum"`
	TosUptime           int     `json:"tos_uptime"`
	TosTime             float32 `json:"tos_time"`
	StartSessionId      string  `json:"start_session_id"`
	ExtraPk
}

func MakeTimeSpentBitArray(inst *Instagram, params map[string]interface{}) *GraphData {
	var ret = &GraphData{}
	ret.Module = "app"
	ret.Name = "time_spent_bit_array"
	ret.Extra = &timeSpentBitArray{}
	initBaseExtra(inst, ret.Extra)
	return ret
}

type igEmergencyPushDidSetInitialVersion struct {
	ExtraPk
	ExtraRadioType
	CurrentVersion int `json:"current_version"`
}

func MakeIgEmergencyPushDidSetInitialVersion(inst *Instagram, params map[string]interface{}) *GraphData {
	var ret = &GraphData{}
	ret.Module = "ig_emergency_push"
	ret.Name = "ig_emergency_push_did_set_initial_version"
	ret.Extra = &igEmergencyPushDidSetInitialVersion{}
	initBaseExtra(inst, ret.Extra)
	return ret
}

type proceedWithPhoneNumber struct {
	ExtraPk
	ExtraRadioType
	Extra2
	ExtraSomeId
	ExtraContactPointType
}

func MakeProceedWithPhoneNumber(inst *Instagram, params map[string]interface{}) *GraphData {
	var ret = &GraphData{}
	ret.Module = "waterfall_log_in"
	ret.Name = "proceed_with_phone_number"
	ret.Extra = &proceedWithPhoneNumber{}
	initBaseExtra(inst, ret.Extra)
	return ret
}

type stepViewLoaded struct {
	Extra2
	ExtraPk
	ExtraRadioType
	ExtraContainerModule
}

func MakeStepViewLoaded(inst *Instagram, params map[string]interface{}) *GraphData {
	var ret = &GraphData{}
	ret.Module = "app"
	ret.Name = "step_view_loaded"
	ret.Extra = &stepViewLoaded{}
	initBaseExtra(inst, ret.Extra)
	return ret
}
