package goinsta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

// StoryReelMention represent story reel mention
type StoryReelMention struct {
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
	Z        int     `json:"z"`
	Width    float64 `json:"width"`
	Height   float64 `json:"height"`
	Rotation float64 `json:"rotation"`
	IsPinned int     `json:"is_pinned"`
	IsHidden int     `json:"is_hidden"`
	User     User
}

// StoryCTA represent story cta
type StoryCTA struct {
	Links []struct {
		LinkType                                int         `json:"linkType"`
		WebURI                                  string      `json:"webUri"`
		AndroidClass                            string      `json:"androidClass"`
		Package                                 string      `json:"package"`
		DeeplinkURI                             string      `json:"deeplinkUri"`
		CallToActionTitle                       string      `json:"callToActionTitle"`
		RedirectURI                             interface{} `json:"redirectUri"`
		LeadGenFormID                           string      `json:"leadGenFormId"`
		IgUserID                                string      `json:"igUserId"`
		AppInstallObjectiveInvalidationBehavior interface{} `json:"appInstallObjectiveInvalidationBehavior"`
	} `json:"links"`
}

//StoryMedia is the struct that handles the information from the methods to get info about Stories.
type StoryMedia struct {
	inst     *Instagram
	endpoint string
	uid      int64

	err error

	Pk              interface{} `json:"id"`
	LatestReelMedia int64       `json:"latest_reel_media"`
	ExpiringAt      float64     `json:"expiring_at"`
	HaveBeenSeen    float64     `json:"seen"`
	CanReply        bool        `json:"can_reply"`
	Title           string      `json:"title"`
	CanReshare      bool        `json:"can_reshare"`
	ReelType        string      `json:"reel_type"`
	User            User        `json:"user"`
	Items           []Item      `json:"items"`
	ReelMentions    []string    `json:"reel_mentions"`
	PrefetchCount   int         `json:"prefetch_count"`
	// this field can be int or bool
	HasBestiesMedia      interface{} `json:"has_besties_media"`
	StoryRankingToken    string      `json:"story_ranking_token"`
	Broadcasts           []Broadcast `json:"broadcasts"`
	FaceFilterNuxVersion int         `json:"face_filter_nux_version"`
	HasNewNuxStory       bool        `json:"has_new_nux_story"`
	Status               string      `json:"status"`
}

// Delete removes instragram story.
//
// See example: examples/media/deleteStories.go
func (media *StoryMedia) Delete() error {
	insta := media.inst

	_, err := insta.HttpRequest(
		&reqOptions{
			ApiPath: fmt.Sprintf(urlMediaDelete, media.ID()),
			Query: map[string]interface{}{
				"media_id": media.ID(),
			},
			IsPost: true,
		},
	)
	return err
}

// ID returns Story id
func (media *StoryMedia) ID() string {
	switch id := media.Pk.(type) {
	case int64:
		return strconv.FormatInt(id, 10)
	case string:
		return id
	}
	return ""
}

func (media *StoryMedia) instagram() *Instagram {
	return media.inst
}

func (media *StoryMedia) setValues() {
	for i := range media.Items {
		setToItem(&media.Items[i], media)
	}
}

// Error returns error happened any error
func (media StoryMedia) Error() error {
	return media.err
}

// Seen marks story as seen.
/*
func (media *StoryMedia) Seen() error {
	insta := media.inst
	data, err := insta.prepareData(
		map[string]interface{}{
			"container_module":   "feed_timeline",
			"live_vods_skipped":  "",
			"nuxes_skipped":      "",
			"nuxes":              "",
			"reels":              "", // TODO xd
			"live_vods":          "",
			"reel_media_skipped": "",
		},
	)
	if err == nil {
		_, err = insta.sendRequest(
			&reqOptions{
				ApiPath: urlMediaSeen, // reel=1&live_vod=0
				Query:    generateSignature(data),
				IsPost:   true,
				UseV2:    true,
			},
		)
	}
	return err
}
*/

type trayRequest struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Sync function is used when Highlight must be sync.
// Highlight must be sync when User.Highlights does not return any object inside StoryMedia slice.
//
// This function does NOT update Stories items.
//
// This function updates StoryMedia.Items
//func (media *StoryMedia) Sync() error {
//	insta := media.inst
//	query := []trayRequest{
//		{"SUPPORTED_SDK_VERSIONS", "9.0,10.0,11.0,12.0,13.0,14.0,15.0,16.0,17.0,18.0,19.0,20.0,21.0,22.0,23.0,24.0"},
//		{"FACE_TRACKER_VERSION", "10"},
//		{"segmentation", "segmentation_enabled"},
//		{"COMPRESSION", "ETC2_COMPRESSION"},
//	}
//	qjson, err := json.Marshal(query)
//	if err != nil {
//		return err
//	}
//
//	id := media.Pk.(string)
//
//
//	body, err := insta.HttpRequest(
//		&reqOptions{
//			ApiPath: urlReelMedia,
//			Query:  	map[string]interface{}{
//				"user_ids":                   []string{id},
//				"supported_capabilities_new": tools.B2s(qjson),
//			},
//			IsPost:   true,
//		},
//	)
//	if err == nil {
//		resp := trayResp{}
//		err = json.Unmarshal(body, &resp)
//		if err == nil {
//			m, ok := resp.Reels[id]
//			if ok {
//				media.Items = m.Items
//				media.setValues()
//				return nil
//			}
//			err = fmt.Errorf("cannot find %s structure in response", id)
//		}
//	}
//	return err
//}

// Next allows pagination after calling:
// User.Stories
//
//
// returns false when list reach the end
// if StoryMedia.Error() is ErrNoMore no problem have been occurred.
func (media *StoryMedia) Next(params ...interface{}) bool {
	if media.err != nil {
		return false
	}

	insta := media.inst
	endpoint := media.endpoint
	if media.uid != 0 {
		endpoint = fmt.Sprintf(endpoint, media.uid)
	}

	body, err := insta.sendSimpleRequest(endpoint)
	if err == nil {
		m := StoryMedia{}
		err = json.Unmarshal(body, &m)
		if err == nil {
			// TODO check NextID media
			*media = m
			media.inst = insta
			media.endpoint = endpoint
			media.err = ErrNoMore // TODO: See if stories has pagination
			media.setValues()
			return true
		}
	}
	media.err = err
	return false
}

// FeedMedia represent a set of media items
type FeedMedia struct {
	inst *Instagram

	err error

	uid       int64
	endpoint  string
	timestamp string

	Items               []Item `json:"items"`
	NumResults          int    `json:"num_results"`
	MoreAvailable       bool   `json:"more_available"`
	AutoLoadMoreEnabled bool   `json:"auto_load_more_enabled"`
	Status              string `json:"status"`
	// Can be int64 and string
	// this is why we recommend Next() usage :')
	NextID interface{} `json:"next_max_id"`
}

// Delete deletes all items in media. Take care...
//
// See example: examples/media/mediaDelete.go
func (media *FeedMedia) Delete() error {
	for i := range media.Items {
		media.Items[i].Delete()
	}
	return nil
}

func (media *FeedMedia) instagram() *Instagram {
	return media.inst
}

// SetInstagram set instagram
func (media *FeedMedia) SetInstagram(inst *Instagram) {
	media.inst = inst
}

// SetID sets media ID
// this value can be int64 or string
func (media *FeedMedia) SetID(id interface{}) {
	media.NextID = id
}

// Sync updates media values.
func (media *FeedMedia) Sync() error {
	id := media.ID()
	insta := media.inst

	body, err := insta.HttpRequest(
		&reqOptions{
			ApiPath: fmt.Sprintf(urlMediaInfo, id),
			Query: map[string]interface{}{
				"media_id": id,
			},
			IsPost: false,
		},
	)
	if err != nil {
		return err
	}

	m := FeedMedia{}
	err = json.Unmarshal(body, &m)
	*media = m
	media.endpoint = urlMediaInfo
	media.inst = insta
	media.NextID = id
	media.setValues()
	return err
}

func (media *FeedMedia) setValues() {
	for i := range media.Items {
		setToItem(&media.Items[i], media)
	}
}

func (media FeedMedia) Error() error {
	return media.err
}

// ID returns media id.
func (media *FeedMedia) ID() string {
	switch s := media.NextID.(type) {
	case string:
		return s
	case int64:
		return strconv.FormatInt(s, 10)
	case json.Number:
		return string(s)
	}
	return ""
}

// Next allows pagination after calling:
// User.Feed
// Params: ranked_content is set to "true" by default, you can set it to false by either passing "false" or false as parameter.
// returns false when list reach the end.
// if FeedMedia.Error() is ErrNoMore no problem have been occurred.
func (media *FeedMedia) Next(params ...interface{}) bool {
	if media.err != nil {
		return false
	}

	insta := media.inst
	endpoint := media.endpoint
	next := media.ID()
	ranked := "true"

	if media.uid != 0 {
		endpoint = fmt.Sprintf(endpoint, media.uid)
	}

	for _, param := range params {
		switch s := param.(type) {
		case string:
			if _, err := strconv.ParseBool(s); err == nil {
				ranked = s
			}
		case bool:
			if !s {
				ranked = "false"
			}
		}
	}
	body, err := insta.HttpRequest(
		&reqOptions{
			ApiPath: endpoint,
			Query: map[string]interface{}{
				"max_id": next,
				//"rank_token":     insta.rankToken,
				"min_timestamp":  media.timestamp,
				"ranked_content": ranked,
			},
		},
	)
	if err == nil {
		m := FeedMedia{}
		d := json.NewDecoder(bytes.NewReader(body))
		d.UseNumber()
		err = d.Decode(&m)
		if err == nil {
			*media = m
			media.inst = insta
			media.endpoint = endpoint
			if m.NextID == 0 || !m.MoreAvailable {
				media.err = ErrNoMore
			}
			media.setValues()
			return true
		}
	}
	return false
}

// MediaItem defines a item media for the
// SavedMedia struct
type MediaItem struct {
	Media Item `json:"media"`
}

// SavedMedia stores the information about media being saved before in my account.
type SavedMedia struct {
	inst     *Instagram
	endpoint string

	err error

	Items []MediaItem `json:"items"`

	NumResults          int    `json:"num_results"`
	MoreAvailable       bool   `json:"more_available"`
	AutoLoadMoreEnabled bool   `json:"auto_load_more_enabled"`
	Status              string `json:"status"`

	NextID interface{} `json:"next_max_id"`
}

// Next allows pagination
func (media *SavedMedia) Next(params ...interface{}) bool {
	// Inital error check
	// if last pagination had errors
	if media.err != nil {
		return false
	}

	insta := media.inst
	endpoint := media.endpoint
	next := media.ID()

	opts := &reqOptions{
		ApiPath: endpoint,
		Query: map[string]interface{}{
			"max_id": next,
		},
	}

	body, err := insta.HttpRequest(opts)
	if err != nil {
		media.err = err
		return false
	}

	m := SavedMedia{}

	if err := json.Unmarshal(body, &m); err != nil {
		media.err = err
		return false
	}

	*media = m

	media.inst = insta
	media.endpoint = endpoint
	media.err = nil

	if m.NextID == 0 || !m.MoreAvailable {
		media.err = ErrNoMore
	}

	media.setValues()

	return true
}

// Error returns the SavedMedia error
func (media *SavedMedia) Error() error {
	return media.err
}

// ID returns the SavedMedia next id
func (media *SavedMedia) ID() string {
	switch id := media.NextID.(type) {
	case int64:
		return strconv.FormatInt(id, 10)
	case string:
		return id
	}
	return ""
}

// Delete method TODO
//
// I think this method should use the
// Unsave method, instead of the Delete.
func (media *SavedMedia) Delete() error {
	return nil
}

// instagram returns the media instagram
func (media *SavedMedia) instagram() *Instagram {
	return media.inst
}

// setValues set the SavedMedia items values
func (media *SavedMedia) setValues() {
	for i := range media.Items {
		setToMediaItem(&media.Items[i], media)
	}
}

// UploadPhoto post image from io.Reader to instagram.
func (this *Instagram) UploadPhoto(photo io.Reader, photoCaption string, quality int, filterType int) (Item, error) {
	out := Item{}

	config, err := this.postPhoto(photo, photoCaption, quality, filterType, false)
	if err != nil {
		return out, err
	}

	body, err := this.HttpRequest(&reqOptions{
		ApiPath: "media/configure/?",
		Query:   config,
		IsPost:  true,
	})
	if err != nil {
		return out, err
	}
	var uploadResult struct {
		Media    Item   `json:"media"`
		UploadID string `json:"upload_id"`
		Status   string `json:"status"`
	}
	err = json.Unmarshal(body, &uploadResult)
	if err != nil {
		return out, err
	}

	if uploadResult.Status != "ok" {
		return out, fmt.Errorf("invalid status, result: %s", uploadResult.Status)
	}

	return uploadResult.Media, nil
}

func (this *Instagram) postPhoto(photo io.Reader, photoCaption string, quality int, filterType int, isSidecar bool) (map[string]interface{}, error) {
	uploadID := time.Now().Unix()
	photoName := fmt.Sprintf("pending_media_%d.jpg", uploadID)
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	w.WriteField("upload_id", strconv.FormatInt(uploadID, 10))
	w.WriteField("_uuid", this.uuid)
	w.WriteField("_csrftoken", this.token)
	var compression = map[string]interface{}{
		"lib_name":    "jt",
		"lib_version": "1.3.0",
		"quality":     quality,
	}
	cBytes, _ := json.Marshal(compression)
	w.WriteField("image_compression", toString(cBytes))
	if isSidecar {
		w.WriteField("is_sidecar", toString(1))
	}
	fw, err := w.CreateFormFile("photo", photoName)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	rdr := io.TeeReader(photo, &buf)
	if _, err = io.Copy(fw, rdr); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", goInstaAPIUrl+"upload/photo/", &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-IG-Capabilities", "3Q4=")
	req.Header.Set("X-IG-Connection-Type", "WIFI")
	req.Header.Set("Cookie2", "$Version=1")
	req.Header.Set("Accept-Language", "en-US")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Content-type", w.FormDataContentType())
	req.Header.Set("Connection", "close")
	req.Header.Set("User-Agent", goInstaUserAgent)

	resp, err := this.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code, result: %s", resp.Status)
	}
	var result struct {
		UploadID       string      `json:"upload_id"`
		XsharingNonces interface{} `json:"xsharing_nonces"`
		Status         string      `json:"status"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if result.Status != "ok" {
		return nil, fmt.Errorf("unknown error, status: %s", result.Status)
	}
	width, height, err := getImageDimensionFromReader(&buf)
	if err != nil {
		return nil, err
	}
	config := map[string]interface{}{
		"media_folder": "Instagram",
		"source_type":  4,
		"caption":      photoCaption,
		"upload_id":    strconv.FormatInt(uploadID, 10),
		"device":       goInstaDeviceSettings,
		"edits": map[string]interface{}{
			"crop_original_size": []int{width * 1.0, height * 1.0},
			"crop_center":        []float32{0.0, 0.0},
			"crop_zoom":          1.0,
			"filter_type":        filterType,
		},
		"extra": map[string]interface{}{
			"source_width":  width,
			"source_height": height,
		},
	}
	return config, nil
}

// UploadAlbum post image from io.Reader to instagram.
func (this *Instagram) UploadAlbum(photos []io.Reader, photoCaption string, quality int, filterType int) (Item, error) {
	out := Item{}

	var childrenMetadata []map[string]interface{}
	for _, photo := range photos {
		config, err := this.postPhoto(photo, photoCaption, quality, filterType, true)
		if err != nil {
			return out, err
		}

		childrenMetadata = append(childrenMetadata, config)
	}
	albumUploadID := time.Now().Unix()

	config := map[string]interface{}{
		"caption":           photoCaption,
		"client_sidecar_id": albumUploadID,
		"children_metadata": childrenMetadata,
	}

	body, err := this.HttpRequest(&reqOptions{
		ApiPath: "media/configure_sidecar/?",
		Query:   config,
		IsPost:  true,
	})
	if err != nil {
		return out, err
	}

	var uploadResult struct {
		Media           Item   `json:"media"`
		ClientSideCarID int64  `json:"client_sidecar_id"`
		Status          string `json:"status"`
	}
	err = json.Unmarshal(body, &uploadResult)
	if err != nil {
		return out, err
	}

	if uploadResult.Status != "ok" {
		return out, fmt.Errorf("invalid status, result: %s", uploadResult.Status)
	}

	return uploadResult.Media, nil
}
