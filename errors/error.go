package errors

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

type Error struct {
	err string
}

func New(text string) *Error {
	_, file, line, _ := runtime.Caller(1)
	idx := strings.LastIndexByte(file, '/')
	return &Error{err: fmt.Sprintf("[%s,%s]", file[idx+1:]+":"+strconv.Itoa(line), text)}
}

func (e *Error) Append(text string) *Error {
	_, file, line, _ := runtime.Caller(1)
	idx := strings.LastIndexByte(file, '/')
	e.err = fmt.Sprintf("[%s,%s,%s]", file[idx+1:]+":"+strconv.Itoa(line), text, e.err)
	return e
}

func (e *Error) Error() string {
	return e.err
}
