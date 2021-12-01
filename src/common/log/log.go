package log

import (
	"fmt"
	"github.com/utahta/go-cronowriter"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	errorRed   = 31
	warnYellow = 33
	infoBlue   = 34
)

type Log struct {
	w           *cronowriter.CronoWriter
	preFit      string
	write2File  bool
	write2Front bool
}

var defaultLog *Log

func setColor(msg string, text int) string {
	return fmt.Sprintf("%c[%dm%s%c[0m", 0x1B, text, msg, 0x1B)
}

func (this *Log) printLog(lev int, data string) {
	_, file, line, ok := runtime.Caller(2)

	var log = "["
	if this.preFit != "" {
		log = this.preFit + " "
	}
	log += time.Now().Format("2006-01-02 15:04:05") + " "
	if ok {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		log += short
		log += ":"
		log += strconv.Itoa(line)
	}
	log += "]\t"
	switch lev {
	case errorRed:
		log += "Error"
		break
	case warnYellow:
		log += "Warn"
		break
	case infoBlue:
		log += "Info"
		break
	}
	log += ": "
	log += data
	log += "\n"

	log = setColor(log, lev)
	if this.write2File {
		this.w.Write([]byte(log))
	}
	if this.write2Front {
		print(log)
	}
}

func Info(format string, v ...interface{}) {
	defaultLog.printLog(infoBlue, fmt.Sprintf(format, v...))
}

func Error(format string, v ...interface{}) {
	defaultLog.printLog(errorRed, fmt.Sprintf(format, v...))
}

func Warn(format string, v ...interface{}) {
	defaultLog.printLog(warnYellow, fmt.Sprintf(format, v...))
}

func ListDir(dirPth string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		files = append(files, fi.Name())
	}
	return files, nil
}

func getTimeArr(start, end string) int64 {
	timeLayout := "20060102"
	loc, _ := time.LoadLocation("Local")
	startUnix, _ := time.ParseInLocation(timeLayout, start, loc)
	endUnix, _ := time.ParseInLocation(timeLayout, end, loc)
	startTime := startUnix.Unix()
	endTime := endUnix.Unix()
	date := (endTime - startTime) / 86400
	return date
}

func cleanLog(t1 interface{}) {
	for {
		select {
		case <-t1.(*time.Ticker).C:
			path, _ := os.Getwd()
			today := time.Now().Format("20060102")
			logs, _ := ListDir(path + "/log")
			for _, logFile := range logs {
				oldFile := strings.Replace(logFile, "log_", "", -1)
				oldFile = strings.Replace(oldFile, ".txt", "", -1)
				arr := getTimeArr(oldFile, today)
				if arr >= 30 {
					Info("delete old log file %s", path+"/log/"+logFile)
					os.Remove(path + "/log/" + logFile)
				}
			}
		}
	}
}

func InitDefaultLog(logName string, write2Front bool, write2File bool) {
	path, _ := os.Getwd()
	defaultLog = &Log{}
	defaultLog.preFit = ""
	defaultLog.write2Front = write2Front
	defaultLog.write2File = write2File
	defaultLog.w = cronowriter.MustNew(path + "/log/log_" + logName + "_%Y%m%d.txt")

	timeClean := time.NewTicker(time.Hour * 24)
	go cleanLog(timeClean)
}
