package reedam

import (
	"errors"
	"fmt"
	"log"
	"path"
	"runtime"
)

type Reedam struct {
	Err     error
	ErrName string
	Message string
	Status  int
}

func New() *Reedam {
	return &Reedam{}
}

func (r *Reedam) WithError(err error) *Reedam {
	r.Err = err
	return r
}

func (r *Reedam) WithErrName(errName string) *Reedam {
	r.ErrName = errName
	return r
}

func (r *Reedam) WithMessage(message string) *Reedam {
	r.Message = message
	return r
}

func (r *Reedam) WithStatus(status int) *Reedam {
	r.Status = status
	return r
}

func (r *Reedam) WithLog(data ...any) *Reedam {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return r
	}
	funcName := path.Base(runtime.FuncForPC(pc).Name())
	fileName := path.Base(file)

	fmt.Println("-------------------------------------------------------------------------")
	log.Println("/ERROR/", r.Err)
	log.Println("/SOURCE/", fmt.Sprintf("file: %s, func: %s, line: %d", fileName, funcName, line))
	log.Println("/DATA/", data)
	fmt.Println("-------------------------------------------------------------------------")

	return r
}

func (r *Reedam) Error() string {
	if r.Err == nil {
		return ""
	}
	return r.Err.Error()
}

func (r *Reedam) GetError() error {
	return r.Err
}

func (r *Reedam) GetMessage() string {
	return r.Message
}

func (r *Reedam) GetStatus() int {
	return r.Status
}

func As(err error) *Reedam {
	if err == nil {
		return nil
	}
	var r *Reedam
	if errors.As(err, &r) {
		return r
	}
	return nil
}
