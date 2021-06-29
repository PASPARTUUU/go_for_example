package lang

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Procedure func()
type ErrorFunction func() error

func Recover(proceed Procedure) Procedure {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Infof("recover from panic: %v", r)
			}
		}()

		proceed()
	}
}

func RecoverWithError(proceed ErrorFunction) (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			logrus.Infof("recover from panic with error: %v", err)
		}
	}()

	err = proceed()
	return
}
