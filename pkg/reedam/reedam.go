package reedam

import (
	"fmt"
	"log"
	"path"
	"runtime"
)

type Reedam struct {
	Err     error
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
		return r.Message
	}
	return r.Err.Error()
}
