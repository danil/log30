package log30_test

import (
	"runtime"
	"testing"
)

func line() int { _, _, l, _ := runtime.Caller(1); return l }

type testprinter struct {
	t    *testing.T
	link string
}

func (p testprinter) Errorf(msg string, args ...interface{}) {
	p.t.Errorf(p.link+"\n"+msg, args...)
}

type Struct struct {
	Name string
	Age  int
}
