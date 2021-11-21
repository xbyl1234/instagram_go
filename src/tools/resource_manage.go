package tools

import (
	"io/ioutil"
	math_rand "math/rand"
	"os"
	"strings"
	"time"
)

type Picture struct {
	filename string
}

type Resourcet struct {
	ico      []Picture
	username []string
}

var Resource Resourcet

func InitResource(icoPath string, usernamePath string) error {
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
		Resource.ico = append(Resource.ico, Picture{icoPath + PthSep + fi.Name()})
	}
	return nil
}

func (this *Resourcet) ChoiceUsername() string {
	math_rand.Seed(time.Now().UnixNano())
	return this.username[math_rand.Intn(len(this.username))]
}

func (this *Resourcet) ChoiceIco() Picture {
	math_rand.Seed(time.Now().UnixNano())
	return this.ico[math_rand.Intn(len(this.ico))]
}
