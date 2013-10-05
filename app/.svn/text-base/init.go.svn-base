package app

//-------------------------------------------------------------------------
//-
//- This package provides a centralized place for all the other packages
//- in this app to pass around configs, logging, security, and database
//- connections.
//-
//-------------------------------------------------------------------------
import (
	"database/sql"
	"fmt"
	_ "github.com/ziutek/mymysql/godrv"
	"log/syslog"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// const (
// 	LOG_LEVEL = 3 // 1-Fatal  2-Error  3-Warn  4-Info  5-Debug
// )

var (
	Cfg         = readConfigFile("config.ini")
	logger, _   = syslog.New(syslog.LOG_LOCAL2, "pid")
	Db          = OpenDB()
	logLevel, _ = strconv.Atoi(Cfg["log.level"])
	AppID       = 0
	LoggerName  = ""
	StartTime   time.Time
)

func Access(message string) {
	myLogger("http", message)
}

func Fatal(message string) {
	if logLevel >= 1 {
		myLogger("FATAL", message)
	}
}

func Error(message string) {
	if logLevel >= 2 {
		myLogger("ERROR", message)
	}
}

func Warn(message string) {
	if logLevel >= 3 {
		myLogger("WARN", message)
	}
}

func Info(message string) {
	if logLevel >= 4 {
		myLogger("INFO", message)
	}
}

func Debug(message string) {
	if logLevel >= 5 {
		myLogger("DEBUG", message)
	}
}

func myLogger(level string, message string) {

	_, file, line, _ := runtime.Caller(2)

	// if Cfg["log.showFileName"] == "On" {
	// 	logger.Printf("[%s][%s] %s - %s:%d", ApiKey, level, message, strings.TrimRight(filepath.Base(file), ".go"), line)
	// } else {
	// 	logger.Printf("[%s][%s] %s", ApiKey, level, message)
	// }

	currentTime := time.Now()
	duration := currentTime.Sub(StartTime)
	if Cfg["log.showFileName"] == "On" {
		logger.Notice(fmt.Sprintf("%04d [%-5s]  %.2fms  %s : %s - %s:%d", AppID, level, duration.Seconds()*1000, LoggerName, message, strings.TrimRight(filepath.Base(file), ".go"), line))
	} else {
		logger.Notice(fmt.Sprintf("%04d [%-5s]  %.2fms  %s : %s", AppID, level, duration.Seconds()*1000, LoggerName, message))
	}
}

func OpenDB() *sql.DB {

	conStr := fmt.Sprintf("%s:%s:%s*%s/%s/%s", Cfg["db.proto"], Cfg["db.host"], Cfg["db.port"], Cfg["db.name"], Cfg["db.user"], Cfg["db.pass"])
	Db, err := sql.Open("mymysql", conStr)
	if err != nil {
		Error(fmt.Sprintf("problem open a connection to %s - %s", Cfg["db.name"], err.Error()))
		return nil
	}

	return Db
}

func fileWatch() {

	for {
		time.Sleep(time.Second * 10)
		// Info("reloading config.ini")
		Cfg = readConfigFile("config.ini")
		logLevel, _ = strconv.Atoi(Cfg["log.level"])
	}

}

func Init() bool {

	//TODO - need to learn a little more and get the public variables declared at the top of the file
	//     - down into this function

	go fileWatch()

	return true

}
