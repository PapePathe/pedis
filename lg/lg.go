package lg

import (
	"fmt"
	"log"
)

func NewLg(id uint64) *Lg {
	return &Lg{
		warn:  log.New(log.Writer(), fmt.Sprintf("%d WARN ", id), log.LstdFlags|log.LUTC),
		info:  log.New(log.Writer(), fmt.Sprintf("%d INFO ", id), log.LstdFlags|log.LUTC|log.Lmsgprefix),
		debug: log.New(log.Writer(), fmt.Sprintf("%d DEBUG ", id), log.Lshortfile),
		err:   log.New(log.Writer(), fmt.Sprintf("%d ERROR ", id), log.Lshortfile),
		fatal: log.New(log.Writer(), fmt.Sprintf("%d FATAL ", id), log.Lshortfile),
	}
}

type Lg struct {
	warn  *log.Logger
	info  *log.Logger
	debug *log.Logger
	err   *log.Logger
	fatal *log.Logger
}

func (rl Lg) Warning(messages ...interface{}) {
	rl.warn.Println(messages...)
}

func (rl Lg) Warningf(_ string, messages ...interface{}) {
	rl.warn.Println(messages...)
}

func (rl Lg) Panic(messages ...interface{}) {
	rl.fatal.Println(messages...)
}

func (rl Lg) Panicf(string, ...interface{}) {}

func (rl Lg) Info(messages ...interface{}) {
	rl.info.Println(messages...)
}

func (rl Lg) Infof(pattern string, messages ...interface{}) {
	rl.info.Println(fmt.Sprintf(pattern, messages...))
}

func (rl Lg) Debug(messages ...interface{}) {
	rl.debug.Println(messages...)
}

func (rl Lg) Fatal(...interface{}) {}

func (rl Lg) Fatalf(string, ...interface{}) {}

func (rl Lg) Error(messages ...interface{}) {
	rl.err.Println(messages...)
}

func (rl Lg) Errorf(string, ...interface{}) {}
func (rl Lg) Debugf(string, ...interface{}) {}
