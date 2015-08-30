package test // import "moetang.info/go/test"

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
	"time"
)

var (
	g_NO_ERROR = errors.New("no error.")
)

type Defer interface {
	MoetangGoTest() bool
}

type deferInfo struct {
	ch <-chan error
}

func (*deferInfo) MoetangGoTest() bool {
	return true
}

type Assertion struct {
	tb testing.TB
}

type errorContainer struct {
	err error
}

func newErrorContainer(a ...interface{}) errorContainer {
	err := errors.New(fmt.Sprintln(a...))
	return errorContainer{
		err: err,
	}
}

func newAssertion(t testing.TB) *Assertion {
	return &Assertion{
		tb: t,
	}
}

func RunTest(t testing.TB, f func()) {
	defer func() {
		if i := recover(); i != nil {
			switch c := i.(type) {
			case errorContainer:
				if c.err != nil {
					t.Fatal(c.err.Error())
				}
			default:
				panic(i)
			}
		}
	}()
	f()
}

func AssertError(err error) {
	if err != nil {
		panic(newErrorContainer(err))
	}
}

func AssertErrorWithMsg(msg string, err error) {
	if err != nil {
		panic(newErrorContainer(msg, ":", err))
	}
}

func AssertNotNull(e interface{}) {
	if e == nil {
		panic(newErrorContainer("param is null"))
	}
}

func AssertNotNullWithMsg(msg string, e interface{}) {
	if e == nil {
		panic(newErrorContainer(msg, ": is null"))
	}
}

func AssertStringEquals(result, expect string) {
	if result != expect {
		panic(newErrorContainer("result:", result, "expect:", expect))
	}
}

func AssertStringEqualsWithMsg(msg, result, expect string) {
	if result != expect {
		panic(newErrorContainer(msg, "result:", result, "expect:", expect))
	}
}

func AssertTrue(b bool) {
	if !b {
		panic(newErrorContainer("param is not true"))
	}
}

func AssertTrueWithMsg(msg string, b bool) {
	if !b {
		panic(newErrorContainer(msg, ": is not true"))
	}
}

func AssertFalse(b bool) {
	if b {
		panic(newErrorContainer("param is not false"))
	}
}

func AssertFalseWithMsg(msg string, b bool) {
	if b {
		panic(newErrorContainer(msg, ": is not false"))
	}
}

func AssertByteSliceEqualsWithMsg(msg string, expect, result []byte) {
	if !bytes.Equal(expect, result) {
		panic(newErrorContainer(msg, ": not equals", "expect:", expect, "result:", result))
	}
}

func AssertByteSliceEquals(expect, result []byte) {
	if !bytes.Equal(expect, result) {
		panic(newErrorContainer("not equals", "expect:", expect, "result:", result))
	}
}

func WaitFor(sec int, d Defer) {
	timeout := time.After(time.Duration(time.Duration(sec) * time.Second))
	ch := d.(*deferInfo).ch
	select {
	case <-timeout:
		panic(newErrorContainer("time out."))
	case e := <-ch:
		if e != nil {
			if e != g_NO_ERROR {
				panic(newErrorContainer(e))
			}
		}
	}
}

func WaitForWithTimeoutMsg(timeoutMsg string, sec int, d Defer) {
	timeout := time.After(time.Duration(time.Duration(sec) * time.Second))
	ch := d.(*deferInfo).ch
	select {
	case <-timeout:
		panic(newErrorContainer(timeoutMsg, ":", "time out."))
	case e := <-ch:
		if e != nil {
			if e != g_NO_ERROR {
				panic(newErrorContainer(e))
			}
		}
	}
}

func DeferTestTask(f func()) Defer {
	ch := make(chan error, 1)
	go func() {
		defer func() {
			if i := recover(); i != nil {
				switch c := i.(type) {
				case errorContainer:
					if c.err != nil {
						ch <- errors.New(c.err.Error())
					}
				default:
					panic(i)
				}
			}
		}()
		f()
		ch <- nil
	}()
	return &deferInfo{
		ch: ch,
	}
}
