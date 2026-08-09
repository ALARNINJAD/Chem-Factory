package reedam

import (
	"chem-factory/pkg/lang"

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

func (r *Reedam) WithLog() *Reedam {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return r
	}
	funcName := path.Base(runtime.FuncForPC(pc).Name())
	fileName := path.Base(file)

	fmt.Println("-------------------------------------------------------------------------")
	log.Println("/ERROR/", r.Err)
	log.Println("/SOURCE/", fmt.Sprintf("file: %s, func: %s, line: %d", fileName, funcName, line))
	fmt.Println("-------------------------------------------------------------------------")

	return r
}

func (r *Reedam) Error() string {
	return r.Err.Error()
}

func InternalError(err error) error {
	return New().WithError(err).WithMessage(lang.ErrorUnexpected).WithStatus(StatusInternalServerError).WithLog()
}
