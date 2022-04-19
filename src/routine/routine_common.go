package routine

import (
	"makemoney/common"
	"makemoney/common/log"
	"os"
	"strings"
)

func LoadKeyWord(KeywordPath string) ([]string, error) {
	var KeyWord []string
	data, err := os.ReadFile(KeywordPath)
	if err != nil {
		log.Error("read key word file error: %v", err)
		return nil, err
	}
	sp := strings.Split(common.B2s(data), "\n")
	KeyWord = make([]string, len(sp))
	index := 0
	for _, item := range sp {
		item = strings.Trim(item, " ")
		item = strings.ReplaceAll(item, "\r", "")
		item = strings.ReplaceAll(item, "\n", "")
		if item != "" {
			KeyWord[index] = item
			index++
		}
	}
	KeyWord = KeyWord[:index]
	return KeyWord, err
}
