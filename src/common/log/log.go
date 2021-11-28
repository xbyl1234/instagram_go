package log

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	colorRed    = 31
	colorYellow = 33
	colorBlue   = 34

	levelT = "[Info] "
	levelE = "[Error] "
	levelW = "[Warn] "

	defaultFileSize = 60 * 1024 * 1024
	minFileSize     = 1 * 1024 * 1024
	defaultLogDir   = "log"
	defaultLogName  = "default.log"

	logTypeStd logType = iota + 1
	logTypeFile
)

type (
	logType int
	myLog   struct {
		sync.Once
		sync.Mutex
		outs     map[logType]io.Writer //writer集合
		file     *os.File              //文件句柄
		fileName string                //日志名
		dir      string                //日志存放路径
		size     int64                 //单个日志文件的大小限制
	}
)

var (
	defaultLogger = &myLog{}
)

func (m *myLog) init() {
	if m.dir == "" {
		m.dir = defaultLogDir
	}
	if m.fileName == "" {
		m.fileName = defaultLogName
	}
	if m.size == 0 {
		m.size = defaultFileSize
	} else {
		if m.size < minFileSize {
			panic(fmt.Sprintf("invalid size: %d", m.size))
		}
	}

	if m.outs == nil {
		m.outs = make(map[logType]io.Writer)
	}
	if !isExist(m.dir) {
		if err := os.Mkdir(m.dir, 0777); err != nil {
			panic(err)
		}
	}
	name := path.Join(m.dir, m.fileName)
	file, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}

	m.file = file
	m.outs[logTypeStd] = os.Stdout
	m.outs[logTypeFile] = file
}

func (m *myLog) checkLogSize() {
	if m.file == nil {
		return
	}
	m.Lock()
	defer m.Unlock()
	fileInfo, err := m.file.Stat()
	if err != nil {
		panic(err)
	}
	if m.size > fileInfo.Size() {
		return
	}
	//需要分割
	newName := path.Join(m.dir, time.Now().Format("2006_01_02_15:04:03")+".log")
	name := path.Join(m.dir, m.fileName)

	err = os.Rename(name, newName)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}

	m.file.Close()
	m.file = file
	m.outs[logTypeFile] = file
	return
}

func (m *myLog) write(level string, content string) {
	m.checkLogSize()
	var colorText int
	switch level {
	case levelT:
		colorText = colorBlue
	case levelW:
		colorText = colorYellow
	case levelE:
		colorText = colorRed
	}

	for k, wr := range m.outs {
		if k == logTypeStd {
			fmt.Fprintf(wr, setColor(content, colorText))
		} else {
			fmt.Fprintf(wr, content)
		}
	}
}

func InitLogger() {
	defaultLogger.init()
}

//Info
func Info(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	timeStr := time.Now().Format("2006-01-02 15:04:05.0000") + " "
	codeLine := "[" + timeStr + shortFileName(file) + ":" + strconv.Itoa(line) + "]"
	content := levelT + codeLine + fmt.Sprintf(format, v...) + "\n"
	defaultLogger.write(levelT, content)
}

//Error
func Error(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	timeStr := time.Now().Format("2006-01-02 15:04:05.0000") + " "
	codeLine := "[" + timeStr + shortFileName(file) + ":" + strconv.Itoa(line) + "]"
	content := levelE + codeLine + fmt.Sprintf(format, v...) + "\n"
	defaultLogger.write(levelE, content)
}

//Warn
func Warn(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	timeStr := time.Now().Format("2006-01-02 15:04:05.0000") + " "
	codeLine := "[" + timeStr + shortFileName(file) + ":" + strconv.Itoa(line) + "]"
	content := levelW + codeLine + fmt.Sprintf(format, v...) + "\n"
	defaultLogger.write(levelW, content)
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func shortFileName(file string) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return short
}

func setColor(msg string, text int) string {
	return fmt.Sprintf("%c[%dm%s%c[0m", 0x1B, text, msg, 0x1B)
}
