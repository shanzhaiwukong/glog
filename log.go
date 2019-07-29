package glog

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

type (
	// Environment 环境
	Environment int
	// Level 级别
	Level int
	// DataBase 数据库接口
	DataBase interface {
		Exec(sql string, params ...interface{}) error
	}
	// Log 结构体
	Log struct {
		lock     sync.Mutex
		duration time.Duration
		list     []string
		filePath string
		isPause  bool
	}
)

const (
	EnvConsole Environment = iota
	EnvDevelop
	EnvProduct
)
const (
	LevelInfo Level = iota
	LevelWarn
	LevelDebug
	LevelError
	LevelPanic
)

var (
	env        = EnvConsole
	timeFormat = "2006/01/02 15:04:05.000"
	mdb        DataBase
	fileSize   int64 = 1024 * 1024 * 20

	reg, _  = regexp.Compile(`_\d{1,}\.?[^\\/]*?$`)
	reg2, _ = regexp.Compile(`\d{1,}`)
	prefix  = map[Level]string{LevelInfo: "INFO", LevelWarn: "WARN", LevelDebug: "DEBUG", LevelError: "ERROR", LevelPanic: "PANIC"}
)

// SetEnvironment 设置当前环境
func SetEnvironment(ev Environment) {
	env = ev
}

// SetTimeFormat 设置日期格式
func SetTimeFormat(tf string) {
	timeFormat = tf
}

// SetDataBase 设置数据库
func SetDataBase(db DataBase) {
	mdb = db
}

// SetFileSize 设置单个文件大小
func SetFileSize(size int64) {
	fileSize = size
}

// New 实例化log
func New(duration time.Duration, filePath string) *Log {
	ext := path.Ext(filePath)
	if ext == "" {
		filePath += time.Now().Format("-20060102_0.log")
	} else {
		filePath = strings.Replace(filePath, ext, time.Now().Format("-20060102_0")+ext, 1)
	}
	log := &Log{
		duration: duration,
		list:     make([]string, 0),
		filePath: filePath,
	}
	go log.interval()
	return log
}

//定时器
func (log *Log) interval() {
	for {
		if log.isPause {
			continue
		}
		time.Sleep(log.duration)
		log.Flush()
	}
}

// Pasue 暂停
func (log *Log) Pasue() {
	log.isPause = true
}

// Continue 继续
func (log *Log) Continue() {
	log.isPause = false
}

// Log 记录信息
func (log *Log) Log(level Level, a ...interface{}) {
	switch level {
	case LevelInfo:
		log.Info(a...)
		break
	case LevelWarn:
		log.Warn(a...)
		break
	case LevelError:
		log.Error(a...)
		break
	case LevelDebug:
		log._debug(2, a...)
		break
	case LevelPanic:
		log.Panic(a...)
		break
	}
}

// Logf 记录信息
func (log *Log) Logf(level Level, format string, a ...interface{}) {
	switch level {
	case LevelInfo:
		log.Infof(format, a...)
		break
	case LevelWarn:
		log.Warnf(format, a...)
		break
	case LevelError:
		log.Errorf(format, a...)
		break
	case LevelDebug:
		log._debugf(2, format, a...)
		break
	case LevelPanic:
		log.Panicf(format, a...)
		break
	}
}

// Info 记录Info信息
func (log *Log) Info(a ...interface{}) {
	log.append(fmt.Sprintf("[%s: %s] ", time.Now().Format(timeFormat), prefix[LevelInfo]) + fmt.Sprintln(a...))
}

// Error 记录Error信息
func (log *Log) Error(a ...interface{}) {
	log.append(fmt.Sprintf("[%s: %s] ", time.Now().Format(timeFormat), prefix[LevelError]) + fmt.Sprintln(a...))
}

// Debug 记录Debug信息 附加调用文件行数信息（FL）
func (log *Log) Debug(a ...interface{}) {
	log._debug(2, a...)
}
func (log *Log) _debug(depth int, a ...interface{}) {
	_, file, line, _ := runtime.Caller(depth)
	log.append(fmt.Sprintf("[%s: %s][FL:%s,%d] ", time.Now().Format(timeFormat), prefix[LevelDebug], file, line) + fmt.Sprintln(a...))
}

// Warn 记录Warn信息
func (log *Log) Warn(a ...interface{}) {
	log.append(fmt.Sprintf("[%s: %s] ", time.Now().Format(timeFormat), prefix[LevelWarn]) + fmt.Sprintln(a...))
}

// Panic 记录Panic信息 附加堆栈信息
func (log *Log) Panic(a ...interface{}) {
	log.append(fmt.Sprintf("[%s: %s][Stack:%s] ", time.Now().Format(timeFormat), prefix[LevelPanic], string(debug.Stack())) + fmt.Sprintln(a...))
}

// Infof 记录Info信息
func (log *Log) Infof(format string, a ...interface{}) {
	log.append(fmt.Sprintf("[%s: %s] ", time.Now().Format(timeFormat), prefix[LevelInfo]) + fmt.Sprintf(format+"\n", a...))
}

// Errorf 记录Error信息
func (log *Log) Errorf(format string, a ...interface{}) {
	log.append(fmt.Sprintf("[%s: %s] ", time.Now().Format(timeFormat), prefix[LevelError]) + fmt.Sprintf(format+"\n", a...))
}

// Debugf 记录Debug信息 附加调用文件行数信息（FL）
func (log *Log) Debugf(format string, a ...interface{}) {
	log._debugf(2, format, a...)
}

func (log *Log) _debugf(depth int, format string, a ...interface{}) {
	_, file, line, _ := runtime.Caller(depth)
	log.append(fmt.Sprintf("[%s: %s][FL:%s,%d] ", time.Now().Format(timeFormat), prefix[LevelDebug], file, line) + fmt.Sprintf(format+"\n", a...))
}

// Warnf 记录Warn信息
func (log *Log) Warnf(format string, a ...interface{}) {
	log.append(fmt.Sprintf("[%s: %s] ", time.Now().Format(timeFormat), prefix[LevelWarn]) + fmt.Sprintf(format+"\n", a...))
}

// Panicf 记录Panic信息 附加堆栈信息
func (log *Log) Panicf(format string, a ...interface{}) {
	log.append(fmt.Sprintf("[%s: %s][Stack:%s] ", time.Now().Format(timeFormat), prefix[LevelPanic], string(debug.Stack())) + fmt.Sprintf(format+"\n", a...))
}

// Db 把信息记录到数据库中
func (log *Log) Db(sql string, params ...interface{}) error {
	log.append(fmt.Sprintf("[%s: %s][%s params:%s]\n", time.Now().Format(timeFormat), "SQL", sql, fmt.Sprint(params...)))
	return mdb.Exec(sql, params...)
}

// Flush 清空缓存区
func (log *Log) Flush() {
	log.lock.Lock()
	list := log.list[0:]
	log.list = make([]string, 0)
	log.lock.Unlock()
	if len(list) == 0 {
		return
	}
	content := strings.Join(list, "")
	if env == EnvConsole || log.filePath == "" {
		fmt.Println(content)
	} else {
		log.writeToFile(content)
	}
}
func (log *Log) append(item string) {
	if log.isPause {
		return
	}
	var count int
	log.lock.Lock()
	defer func() {
		log.lock.Unlock()
		if int64(count*1024) > fileSize {
			log.Flush()
		}
	}()
	log.list = append(log.list, item)
	count = len(log.list)
}

func (log *Log) writeToFile(content string) {
	log.lock.Lock()
	defer log.lock.Unlock()
check:
	file, err := os.Stat(log.filePath)
	if os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(log.filePath), 0333)
		os.Create(log.filePath)
	} else {
		if file.Size() > fileSize {
			log.filePath = string(reg.ReplaceAllFunc([]byte(log.filePath), func(in []byte) []byte {
				return reg2.ReplaceAllFunc(in, func(_in []byte) []byte {
					n, _ := strconv.Atoi(string(_in))
					return []byte(strconv.Itoa(n + 1))
				})
			}))
			goto check
		}
	}
	if _file, err := os.OpenFile(log.filePath, os.O_APPEND|os.O_WRONLY, 0333); err == nil {
		defer _file.Close()
		_file.WriteString(content)
	}
}
