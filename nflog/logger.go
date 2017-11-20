package nflog

import (
	"io"
	"log"
)

var (
	info       *log.Logger
	err        *log.Logger
	debug      *log.Logger
	fatal      *log.Logger
	moduleName string
)

func Init(infoHandle io.Writer, errHandle io.Writer, debugHandle io.Writer, fatalHandle io.Writer, name string) {
	info = log.New(infoHandle, "[Info]"+"["+name+"]", log.Ldate|log.Ltime)
	err = log.New(errHandle, "[Error]"+"["+name+"]", log.Ldate|log.Ltime)
	debug = log.New(debugHandle, "[Debug]"+"["+name+"]", log.Ldate|log.Ltime|log.Lshortfile)
	fatal = log.New(fatalHandle, "[Fatal]"+"["+name+"]", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(format string, v ...interface{}) {
	info.Printf(format, v...)
}

func Err(format string, v ...interface{}) {
	err.Printf(format, v...)
}

func Debug(format string, v ...interface{}) {
	debug.Printf(format, v...)
}

func Fatal(format string, v ...interface{}) {
	fatal.Fatalf(format, v...)
}
