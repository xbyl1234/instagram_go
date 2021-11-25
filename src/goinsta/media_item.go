package goinsta

import (
	"encoding/json"
	"fmt"
	"io"
	"makemoney/common/log"
	neturl "net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

type Media interface {
	// Next allows pagination
	Next(...interface{}) bool
	// Error returns error (in case it have been occurred)
	Error() error
	// ID returns media id
	ID() string
	// Delete removes media
	Delete() error

	instagram() *Instagram
}

type MediaType int

var (
	MediaType_Photo    MediaType = 1
	MediaType_Video    MediaType = 2
	MediaType_Carousel MediaType = 8
)

type Item struct {
	media    Media
	Comments *Comments `json:"-"`

	CanSeeInsightsAsBrand      bool   `json:"can_see_insights_as_brand"`
	CanViewMorePreviewComments bool   `json:"can_view_more_preview_comments"`
	CommercialityStatus        string `json:"commerciality_status"`
	DeletedReason              int    `json:"deleted_reason"`
	FundraiserTag              struct {
		HasStandaloneFundraiser bool `json:"has_standalone_fundraiser"`
	} `json:"fundraiser_tag"`
	HideViewAllCommentEntrypoint bool   `json:"hide_view_all_comment_entrypoint"`
	IntegrityReviewDecision      string `json:"integrity_review_decision"`
	IsCommercial                 bool   `json:"is_commercial"`
	IsInProfileGrid              bool   `json:"is_in_profile_grid"`
	IsPaidPartnership            bool   `json:"is_paid_partnership"`
	IsUnifiedVideo               bool   `json:"is_unified_video"`
	LikeAndViewCountsDisabled    bool   `json:"like_and_view_counts_disabled"`
	NextMaxId                    int64  `json:"next_max_id"`
	ProductType                  string `json:"product_type"`
	ProfileGridControlEnabled    bool   `json:"profile_grid_control_enabled"`

	TakenAt          int64   `json:"taken_at"`
	Pk               int64   `json:"pk"`
	ID               string  `json:"id"`
	CommentsDisabled bool    `json:"comments_disabled"`
	DeviceTimestamp  int64   `json:"device_timestamp"`
	MediaType        int     `json:"media_type"`
	Code             string  `json:"code"`
	ClientCacheKey   string  `json:"client_cache_key"`
	FilterType       int     `json:"filter_type"`
	CarouselParentID string  `json:"carousel_parent_id"`
	CarouselMedia    []Item  `json:"carousel_media,omitempty"`
	User             User    `json:"user"`
	CanViewerReshare bool    `json:"can_viewer_reshare"`
	Caption          Caption `json:"caption"`
	CaptionIsEdited  bool    `json:"caption_is_edited"`
	Likes            int     `json:"like_count"`
	HasLiked         bool    `json:"has_liked"`
	// Toplikers can be `string` or `[]string`.
	// Use TopLikers function instead of getting it directly.
	Toplikers                    interface{} `json:"top_likers"`
	Likers                       []User      `json:"likers"`
	CommentLikesEnabled          bool        `json:"comment_likes_enabled"`
	CommentThreadingEnabled      bool        `json:"comment_threading_enabled"`
	HasMoreComments              bool        `json:"has_more_comments"`
	MaxNumVisiblePreviewComments int         `json:"max_num_visible_preview_comments"`
	// Previewcomments can be `string` or `[]string` or `[]Comment`.
	// Use PreviewComments function instead of getting it directly.
	Previewcomments interface{} `json:"preview_comments,omitempty"`
	CommentCount    int         `json:"comment_count"`
	PhotoOfYou      bool        `json:"photo_of_you"`
	// Tags are tagged people in photo
	Tags struct {
		In []Tag `json:"in"`
	} `json:"usertags,omitempty"`
	FbUserTags           Tag    `json:"fb_user_tags"`
	CanViewerSave        bool   `json:"can_viewer_save"`
	OrganicTrackingToken string `json:"organic_tracking_token"`
	// Images contains URL images in different versions.
	// Version = quality.
	Images          Images   `json:"image_versions2,omitempty"`
	OriginalWidth   int      `json:"original_width,omitempty"`
	OriginalHeight  int      `json:"original_height,omitempty"`
	ImportedTakenAt int64    `json:"imported_taken_at,omitempty"`
	Location        Location `json:"location,omitempty"`
	Lat             float64  `json:"lat,omitempty"`
	Lng             float64  `json:"lng,omitempty"`

	// Videos
	Videos            []Video `json:"video_versions,omitempty"`
	HasAudio          bool    `json:"has_audio,omitempty"`
	VideoDuration     float64 `json:"video_duration,omitempty"`
	ViewCount         float64 `json:"view_count,omitempty"`
	IsDashEligible    int     `json:"is_dash_eligible,omitempty"`
	VideoDashManifest string  `json:"video_dash_manifest,omitempty"`
	NumberOfQualities int     `json:"number_of_qualities,omitempty"`

	// Only for stories
	StoryEvents              []interface{}      `json:"story_events"`
	StoryHashtags            []interface{}      `json:"story_hashtags"`
	StoryPolls               []interface{}      `json:"story_polls"`
	StoryFeedMedia           []interface{}      `json:"story_feed_media"`
	StorySoundOn             []interface{}      `json:"story_sound_on"`
	CreativeConfig           interface{}        `json:"creative_config"`
	StoryLocations           []interface{}      `json:"story_locations"`
	StorySliders             []interface{}      `json:"story_sliders"`
	StoryQuestions           []interface{}      `json:"story_questions"`
	StoryProductItems        []interface{}      `json:"story_product_items"`
	StoryCTA                 []StoryCTA         `json:"story_cta"`
	ReelMentions             []StoryReelMention `json:"reel_mentions"`
	SupportsReelReactions    bool               `json:"supports_reel_reactions"`
	ShowOneTapFbShareTooltip bool               `json:"show_one_tap_fb_share_tooltip"`
	HasSharedToFb            int64              `json:"has_shared_to_fb"`
	Mentions                 []Mentions
	Audience                 string `json:"audience,omitempty"`
	StoryMusicStickers       []struct {
		X              float64 `json:"x"`
		Y              float64 `json:"y"`
		Z              int     `json:"z"`
		Width          float64 `json:"width"`
		Height         float64 `json:"height"`
		Rotation       float64 `json:"rotation"`
		IsPinned       int     `json:"is_pinned"`
		IsHidden       int     `json:"is_hidden"`
		IsSticker      int     `json:"is_sticker"`
		MusicAssetInfo struct {
			ID                       string `json:"id"`
			Title                    string `json:"title"`
			Subtitle                 string `json:"subtitle"`
			DisplayArtist            string `json:"display_artist"`
			CoverArtworkURI          string `json:"cover_artwork_uri"`
			CoverArtworkThumbnailURI string `json:"cover_artwork_thumbnail_uri"`
			ProgressiveDownloadURL   string `json:"progressive_download_url"`
			HighlightStartTimesInMs  []int  `json:"highlight_start_times_in_ms"`
			IsExplicit               bool   `json:"is_explicit"`
			DashManifest             string `json:"dash_manifest"`
			HasLyrics                bool   `json:"has_lyrics"`
			AudioAssetID             string `json:"audio_asset_id"`
			IgArtist                 struct {
				Pk            int    `json:"pk"`
				Username      string `json:"username"`
				FullName      string `json:"full_name"`
				IsPrivate     bool   `json:"is_private"`
				ProfilePicURL string `json:"profile_pic_url"`
				ProfilePicID  string `json:"profile_pic_id"`
				IsVerified    bool   `json:"is_verified"`
			} `json:"ig_artist"`
			PlaceholderProfilePicURL string `json:"placeholder_profile_pic_url"`
			ShouldMuteAudio          bool   `json:"should_mute_audio"`
			ShouldMuteAudioReason    string `json:"should_mute_audio_reason"`
			OverlapDurationInMs      int    `json:"overlap_duration_in_ms"`
			AudioAssetStartTimeInMs  int    `json:"audio_asset_start_time_in_ms"`
		} `json:"music_asset_info"`
	} `json:"story_music_stickers,omitempty"`
}

// Comment pushes a text comment to media item.
//
// If parent media is a Story this function will send a private message
// replying the Instagram story.
func (item *Item) Comment(text string) error {
	var opt *reqOptions
	insta := item.media.instagram()

	switch item.media.(type) {
	case *StoryMedia:
		//to, err := prepareRecipients(item)
		//if err != nil {
		//	return err
		//}
		//
		//query := insta.prepareDataQuery(
		//	map[string]interface{}{
		//		"recipient_users": to,
		//		"action":          "send_item",
		//		"media_id":        item.ID,
		//		"client_context":  generateUUID(),
		//		"text":            text,
		//		"entry":           "reel",
		//		"reel_id":         item.User.ID,
		//	},
		//)
		//opt = &reqOptions{
		//	Connection: "keep-alive",
		//	ApiPath:   fmt.Sprintf("%s?media_type=%s", urlReplyStory, item.GetMediaType()),
		//	Query:      query,
		//	IsPost:     true,
		//}
	case *FeedMedia: // normal media
		opt = &reqOptions{
			ApiPath: fmt.Sprintf(urlCommentAdd, item.Pk),
			Query: map[string]interface{}{
				"comment_text": text,
			},
			IsPost: true,
		}
	}

	// ignoring response
	_, err := insta.HttpRequest(opt)
	return err
}

func (item *Item) GetMediaType() MediaType {
	switch item.MediaType {
	case 1:
		return MediaType_Photo
	case 2:
		return MediaType_Video
	case 8:
		return MediaType_Carousel
	}
	log.Error("GetMediaType error: %v", item.MediaType)
	return MediaType_Photo
}

func setToItem(item *Item, media Media) {
	item.media = media
	item.User.inst = media.instagram()
	item.Comments = newComments(item)
	for i := range item.CarouselMedia {
		item.CarouselMedia[i].User = item.User
		setToItem(&item.CarouselMedia[i], media)
	}
}

// setToMediaItem is a utility function that
// mimics the setToItem but for the SavedMedia items
func setToMediaItem(item *MediaItem, media Media) {
	item.Media.media = media
	item.Media.User.inst = media.instagram()

	item.Media.Comments = newComments(&item.Media)

	for i := range item.Media.CarouselMedia {
		item.Media.CarouselMedia[i].User = item.Media.User
		setToItem(&item.Media.CarouselMedia[i], media)
	}
}

func getname(name string) string {
	nname := name
	i := 1
	for {
		ext := path.Ext(name)

		_, err := os.Stat(name)
		if err != nil {
			break
		}
		if ext != "" {
			nname = strings.Replace(nname, ext, "", -1)
		}
		name = fmt.Sprintf("%s.%d%s", nname, i, ext)
		i++
	}
	return name
}

func download(inst *Instagram, url, dst string) (string, error) {
	file, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer file.Close()

	resp, err := inst.c.Get(url)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, resp.Body)
	return dst, err
}

type bestMedia struct {
	w, h int
	url  string
}

// GetBest returns best quality image or video.
//
// Arguments can be []Video or []Candidate
func GetBest(obj interface{}) string {
	m := bestMedia{}

	switch t := obj.(type) {
	// getting best video
	case []Video:
		for _, video := range t {
			if m.w < video.Width && video.Height > m.h && video.URL != "" {
				m.w = video.Width
				m.h = video.Height
				m.url = video.URL
			}
		}
		// getting best image
	case []Candidate:
		for _, image := range t {
			if m.w < image.Width && image.Height > m.h && image.URL != "" {
				m.w = image.Width
				m.h = image.Height
				m.url = image.URL
			}
		}
	}
	return m.url
}

var rxpTags = regexp.MustCompile(`#\w+`)

// Hashtags returns caption hashtags.
//
// Item media parent must be FeedMedia.
//
// See example: examples/media/hashtags.go
func (item *Item) Hashtags() []Tags {
	tags := rxpTags.FindAllString(item.Caption.Text, -1)

	hsh := make([]Tags, len(tags))

	i := 0
	for _, tag := range tags {
		hsh[i].name = tag[1:]
		i++
	}

	for _, comment := range item.PreviewComments() {
		tags := rxpTags.FindAllString(comment.Text, -1)

		for _, tag := range tags {
			hsh = append(hsh, Tags{name: tag[1:]})
		}
	}

	return hsh
}

// Delete deletes your media item. StoryMedia or FeedMedia
//
// See example: examples/media/mediaDelete.go
func (item *Item) Delete() error {
	insta := item.media.instagram()

	_, err := insta.HttpRequest(
		&reqOptions{
			ApiPath: fmt.Sprintf(urlMediaDelete, item.ID),
			Query: map[string]interface{}{
				"media_id": item.ID,
			},
			IsPost: true,
		},
	)
	return err
}

// SyncLikers fetch new likers of a media
//
// This function updates Item.Likers value
func (item *Item) SyncLikers() error {
	resp := respLikers{}
	insta := item.media.instagram()
	body, err := insta.sendSimpleRequest(urlMediaLikers, item.ID)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &resp)
	if err == nil {
		item.Likers = resp.Users
	}
	return err
}

// Unlike mark media item as unliked.
//
// See example: examples/media/unlike.go
func (item *Item) Unlike() error {
	insta := item.media.instagram()

	_, err := insta.HttpRequest(
		&reqOptions{
			ApiPath: fmt.Sprintf(urlMediaUnlike, item.ID),
			Query: map[string]interface{}{
				"media_id": item.ID,
			},
			IsPost: true,
		},
	)
	return err
}

// Like mark media item as liked.
//
// See example: examples/media/like.go
func (item *Item) Like() error {
	insta := item.media.instagram()

	_, err := insta.HttpRequest(
		&reqOptions{
			ApiPath: fmt.Sprintf(urlMediaLike, item.ID),
			Query: map[string]interface{}{
				"media_id": item.ID,
			},
			IsPost: true,
		},
	)
	return err
}

// Save saves media item.
//
// You can get saved media using Account.Saved()
func (item *Item) Save() error {
	insta := item.media.instagram()

	_, err := insta.HttpRequest(
		&reqOptions{
			ApiPath: fmt.Sprintf(urlMediaSave, item.ID),
			Query: map[string]interface{}{
				"media_id": item.ID,
			},
			IsPost: true,
		},
	)
	return err
}

// Unsave unsaves media item.
func (item *Item) Unsave() error {
	insta := item.media.instagram()

	_, err := insta.HttpRequest(
		&reqOptions{
			ApiPath: fmt.Sprintf(urlMediaUnsave, item.ID),
			Query: map[string]interface{}{
				"media_id": item.ID,
			},
			IsPost: true,
		},
	)
	return err
}

// Download downloads media item (video or image) with the best quality.
//
// Input parameters are folder and filename. If filename is "" will be saved with
// the default value name.
//
// If file exists it will be saved
// This function makes folder automatically
//
// This function returns an slice of location of downloaded items
// The returned values are the output path of images and videos.
//
// This function does not download CarouselMedia.
//
// See example: examples/media/itemDownload.go
func (item *Item) Download(folder, name string) (imgs, vds string, err error) {
	var u *neturl.URL
	var nname string
	imgFolder := path.Join(folder, "images")
	vidFolder := path.Join(folder, "videos")
	inst := item.media.instagram()

	os.MkdirAll(folder, 0777)
	os.MkdirAll(imgFolder, 0777)
	os.MkdirAll(vidFolder, 0777)

	vds = GetBest(item.Videos)
	if vds != "" {
		if name == "" {
			u, err = neturl.Parse(vds)
			if err != nil {
				return
			}

			nname = path.Join(vidFolder, path.Base(u.Path))
		} else {
			nname = path.Join(vidFolder, name)
		}
		nname = getname(nname)

		vds, err = download(inst, vds, nname)
		return "", vds, err
	}

	imgs = GetBest(item.Images.Versions)
	if imgs != "" {
		if name == "" {
			u, err = neturl.Parse(imgs)
			if err != nil {
				return
			}

			nname = path.Join(imgFolder, path.Base(u.Path))
		} else {
			nname = path.Join(imgFolder, name)
		}
		nname = getname(nname)

		imgs, err = download(inst, imgs, nname)
		return imgs, "", err
	}

	return imgs, vds, fmt.Errorf("cannot find any image or video")
}

// TopLikers returns string slice or single string (inside string slice)
// Depending on TopLikers parameter.
func (item *Item) TopLikers() []string {
	switch s := item.Toplikers.(type) {
	case string:
		return []string{s}
	case []string:
		return s
	}
	return nil
}

// PreviewComments returns string slice or single string (inside Comment slice)
// Depending on PreviewComments parameter.
// If PreviewComments are string or []string only the Text field will be filled.
func (item *Item) PreviewComments() []Comment {
	switch s := item.Previewcomments.(type) {
	case []interface{}:
		if len(s) == 0 {
			return nil
		}

		switch s[0].(type) {
		case interface{}:
			comments := make([]Comment, 0)
			for i := range s {
				if buf, err := json.Marshal(s[i]); err != nil {
					return nil
				} else {
					comment := &Comment{}

					if err = json.Unmarshal(buf, comment); err != nil {
						return nil
					} else {
						comments = append(comments, *comment)
					}
				}
			}
			return comments
		case string:
			comments := make([]Comment, 0)
			for i := range s {
				comments = append(comments, Comment{
					Text: s[i].(string),
				})
			}
			return comments
		}
	case string:
		comments := []Comment{
			{
				Text: s,
			},
		}
		return comments
	}
	return nil
}

// StoryIsCloseFriends returns a bool
// If the returned value is true the story was published only for close friends
func (item *Item) StoryIsCloseFriends() bool {
	return item.Audience == "besties"
}
