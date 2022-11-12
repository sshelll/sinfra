package stream

type ErrorPasser struct {
	errCh chan error
}

func NewErrorPasser() *ErrorPasser {
	return &ErrorPasser{
		errCh: make(chan error, 2),
	}
}

func (e *ErrorPasser) Check() (err error, done bool) {
	select {
	case err, ok := <-e.errCh:
		if err != nil {
			return err, true
		}
		return nil, !ok
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
