package stream

type ErrorPasser struct {
	errCh chan error
}

func NewErrorPasser() *ErrorPasser {
	return &ErrorPasser{
		// default size is 2
		errCh: make(chan error, 2),
	}
}

func NewErrorPasserWithCap(maxErrCnt int) *ErrorPasser {
	if maxErrCnt < 0 {
		maxErrCnt = 0
	}
	return &ErrorPasser{
		errCh: make(chan error, maxErrCnt),
	}
}

// Check try get err in a non-block way.
// NOTE: if done, err is nil.
func (e *ErrorPasser) Check() (err error, done bool) {
	select {
	case err, ok := <-e.errCh:
		return err, !ok
	default:
		return nil, false
	}
}

func (e *ErrorPasser) Get() error {
	for err := range e.errCh {
		return err
	}
	return nil
}

func (e *ErrorPasser) Put(err error) {
	e.errCh <- err
}

func (e *ErrorPasser) Close() {
	close(e.errCh)
}

func (e *ErrorPasser) Cap() int {
	return cap(e.errCh)
}
