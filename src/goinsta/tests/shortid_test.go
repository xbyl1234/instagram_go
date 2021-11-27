package tests

import (
	"encoding/json"
	"testing"
)

type name struct {
	CachedCommentsCursor string `json:"cached_comments_cursor"`
}

func TestMediaIDFromShortID(t *testing.T) {
	ret := &name{}
	print(json.Unmarshal(([]byte("{ \"cached_comments_cursor\": \"17938678366697942\"}"))[:], ret))
	print(ret)
	//mediaID, err := goinsta.MediaIDFromShortID("BR_repxhx4O")
	//if err != nil {
	//	t.Fatal(err)
	//	return
	//}
	//if mediaID != "1477090425239445006" {
	//	t.Fatal("Invalid mediaID")
	//}
}
