package goinsta

import "fmt"

type UserOperate struct {
	inst *Instagram
}

func newUserOperate(inst *Instagram) *UserOperate {
	return &UserOperate{inst: inst}
}

type RespLikeUser struct {
	BaseApiResp
	PreviousFollowing bool       `json:"previous_following,omitempty"`
	FriendshipStatus  Friendship `json:"friendship_status"`
}

func (this *UserOperate) LikeUser(userID int64) error {
	params := map[string]interface{}{
		"_uuid":            this.inst.Device.DeviceID,
		"_uid":             this.inst.ID,
		"user_id":          userID,
		"device_id":        this.inst.Device.DeviceID,
		"container_module": "profile",
	}
	resp := &RespLikeUser{}
	err := this.inst.HttpRequestJson(&reqOptions{
		ApiPath: fmt.Sprintf(urlUserFollow, userID),
		IsPost:  true,
		Signed:  true,
		Query:   params,
	}, resp)

	err = resp.CheckError(err)
	return err
}
