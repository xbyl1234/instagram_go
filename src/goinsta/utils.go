package goinsta

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"image"
	"makemoney/common"
	"time"

	// Required for getImageDimensionFromReader in jpg and png format
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strconv"
)

func toString(i interface{}) string {
	switch s := i.(type) {
	case string:
		return s
	case bool:
		return strconv.FormatBool(s)
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32)
	case int:
		return strconv.Itoa(s)
	case int64:
		return strconv.FormatInt(s, 10)
	case int32:
		return strconv.Itoa(int(s))
	case int16:
		return strconv.FormatInt(int64(s), 10)
	case int8:
		return strconv.FormatInt(int64(s), 10)
	case uint:
		return strconv.FormatInt(int64(s), 10)
	case uint64:
		return strconv.FormatInt(int64(s), 10)
	case uint32:
		return strconv.FormatInt(int64(s), 10)
	case uint16:
		return strconv.FormatInt(int64(s), 10)
	case uint8:
		return strconv.FormatInt(int64(s), 10)
	case []byte:
		return common.B2s(s)
	case error:
		return s.Error()
	}
	return ""
}

//func prepareRecipients(cc interface{}) (bb string, err error) {
//	var b []byte
//	ids := make([][]int64, 0)
//	switch c := cc.(type) {
//	case *Conversation:
//		for i := range c.Users {
//			ids = append(ids, []int64{c.Users[i].ID})
//		}
//	case *Item:
//		ids = append(ids, []int64{c.User.ID})
//	case int64:
//		ids = append(ids, []int64{c})
//	}
//	b, err = json.Marshal(ids)
//	bb = tools.B2s(b)
//	return
//}

// getImageDimensionFromReader return image dimension , types is .jpg and .png
func getImageDimensionFromReader(rdr io.Reader) (int, int, error) {
	image, _, err := image.DecodeConfig(rdr)
	if err != nil {
		return 0, 0, err
	}
	return image.Width, image.Height, nil
}

func RSAEncrypt(pubKey []byte, plainText []byte) ([]byte, error) {
	block, _ := pem.Decode(pubKey)
	//x509解码
	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	//类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	//对明文进行加密
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
	if err != nil {
		return nil, err
	}
	//返回密文
	return cipherText, nil
}

func AesGcmEncrypt(key []byte, iv []byte, plainText []byte, add []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, iv, plainText, add)
	return ciphertext, nil
}

func encryptPassword(password string, encId string, encPubKey string) (string, error) {
	//byte[] rand_key = new byte[32], iv = new byte[12];
	_time := strconv.FormatInt(time.Now().Unix(), 10)
	randKey := common.GenString(common.CharSet_All, 32)
	iv := common.GenString(common.CharSet_All, 12)
	decodedPubKey, err := base64.StdEncoding.DecodeString(encPubKey)
	if err != nil {
		return "", err
	}

	randKeyEncrypted, err := RSAEncrypt(decodedPubKey, []byte(randKey))
	if err != nil {
		return "", err
	}
	passwordEncrypted, err := AesGcmEncrypt([]byte(randKey), []byte(iv), []byte(password), []byte(_time))
	if err != nil {
		return "", err
	}

	buff := bytes.Buffer{}
	buff.WriteByte(1)
	encid, _ := strconv.Atoi(encId)
	buff.WriteByte(byte(encid))
	buff.WriteString(iv)
	lenByte := make([]byte, 2)
	binary.LittleEndian.PutUint16(lenByte, uint16(len(randKeyEncrypted)))
	buff.Write(lenByte)
	buff.Write(randKeyEncrypted)
	buff.Write(passwordEncrypted[len(passwordEncrypted)-16:])
	buff.Write(passwordEncrypted[:len(passwordEncrypted)-16])

	return fmt.Sprintf("#PWD_INSTAGRAM:%s:%s:%s", "4", _time,
		base64.StdEncoding.EncodeToString(buff.Bytes())), nil
}
