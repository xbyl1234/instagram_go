package common

import (
	"io/ioutil"
	math_rand "math/rand"
	"os"
	"strings"
	"time"
)

type Resourcet struct {
	ico      []string
	username []string
}

var Resource Resourcet

func InitResource(icoPath string, usernamePath string) error {
	if usernamePath != "" {
		data, err := ioutil.ReadFile(usernamePath)
		if err != nil {
			return err
		}
		sp := strings.Split(string(data), "\n")
		for idx := range sp {
			username := sp[idx]
			username = strings.ReplaceAll(username, " ", "")
			username = strings.ReplaceAll(username, "\n", "")
			username = strings.ReplaceAll(username, "\r", "")
			if len(username) > 3 {
				Resource.username = append(Resource.username, username)
			}
		}
	}

	if icoPath != "" {
		//ico
		dir, err := ioutil.ReadDir(icoPath)
		if err != nil {
			return err
		}
		PthSep := string(os.PathSeparator)
		for _, fi := range dir {
			if fi.IsDir() {
				continue
			}
			Resource.ico = append(Resource.ico, icoPath+PthSep+fi.Name())
		}
	}
	return nil
}

func (this *Resourcet) ChoiceUsername() string {
	math_rand.Seed(time.Now().UnixNano())
	return this.username[math_rand.Intn(len(this.username))]
}

func (this *Resourcet) ChoiceIco() string {
	math_rand.Seed(time.Now().UnixNano())
	return this.ico[math_rand.Intn(len(this.ico))]
}
