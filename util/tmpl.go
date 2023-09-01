package util

import (
	"fmt"
	"runtime"

	"gorm.io/gorm"
)

func FirstError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func AllowPanic(fn func() error) (panicked bool, err error) {
	if fn == nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			err = fmt.Errorf("panic: %v\n%s", r, buf[:n])
		}
	}()
	err = fn()
	return
}

func ExecWithinTransaction(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.Debug().Begin()
	if err := tx.Error; err != nil {
		return err
	}

	if _, err := AllowPanic(func() error {
		return fn(tx)
	}); err != nil {
		return FirstError(tx.Rollback().Error, err)
	}

	if err := tx.Commit().Error; err != nil {
		return FirstError(tx.Rollback().Error, err)
	}

	return nil
}
