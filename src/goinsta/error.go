package goinsta

import "fmt"

// ErrorN is general instagram error
type ErrorN struct {
	Message   string `json:"message"`
	Status    string `json:"status"`
	ErrorType string `json:"error_type"`
}

// Error503 is instagram API error
type Error503 struct {
	Message string
}

func (e Error503) Error() string {
	return e.Message
}

func (e ErrorN) Error() string {
	return fmt.Sprintf("%s: %s (%s)", e.Status, e.Message, e.ErrorType)
}

// Error400 is error returned by HTTP 400 status code.
type Error400 struct {
	ChallengeError
	Action     string `json:"action"`
	StatusCode string `json:"status_code"`
	Payload    struct {
		ClientContext string `json:"client_context"`
		Message       string `json:"message"`
	} `json:"payload"`
	Status string `json:"status"`
}

func (e Error400) Error() string {
	return fmt.Sprintf("%s: %s", e.Status, e.Payload.Message)
}

// ChallengeError is error returned by HTTP 400 status code.
type ChallengeError struct {
	Message   string `json:"message"`
	Challenge struct {
		URL               string `json:"url"`
		APIPath           string `json:"api_path"`
		HideWebviewHeader bool   `json:"hide_webview_header"`
		Lock              bool   `json:"lock"`
		Logout            bool   `json:"logout"`
		NativeFlow        bool   `json:"native_flow"`
	} `json:"challenge"`
	Status    string `json:"status"`
	ErrorType string `json:"error_type"`
}

func (e ChallengeError) Error() string {
	return fmt.Sprintf("%s: %s", e.Status, e.Message)
}

type ApiError struct {
	ErrorType string `json:"error_type"`
}

func (e ApiError) Error() string {
	return e.ErrorType
}

func CheckApiError(resp interface{}, err error) error {
	apiResp := (resp).(BaseApiResp)
	if apiResp.isError() {
		return &ApiError{apiResp.ErrorType}
	}
	return err
}
