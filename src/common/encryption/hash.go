package encryption

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/twmb/murmur3"
)

func Hash(data []byte) (string, error) {
	return Mmh3Hash32(StandBase64(data))
}

func Mmh3Hash32(raw []byte) (string, error) {
	var h32 = murmur3.New32()
	_, err := h32.Write(raw)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", int32(h32.Sum32())), nil
}

func StandBase64(braw []byte) []byte {
	bckd := base64.StdEncoding.EncodeToString(braw)
	var buffer bytes.Buffer
	for i := 0; i < len(bckd); i++ {
		ch := bckd[i]
		buffer.WriteByte(ch)
		if (i+1)%76 == 0 {
			buffer.WriteByte('\n')
		}
	}
	buffer.WriteByte('\n')
	return buffer.Bytes()
}
