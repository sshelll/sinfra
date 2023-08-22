package promise

func All(promises ...*Promise) *Promise {
	return New(func(resolve, reject func(v any)) {
		var values []any
		for _, p := range promises {
			p.Await()
			values = append(values, p.Result())
		}
		resolve(values)
	})
}

func Race(promises ...*Promise) *Promise {
	safeClose := func(ch chan struct{}) {
		defer func() {
			_ = recover()
		}()
		close(ch)
	}
	return New(func(resolve, reject func(v any)) {
		ch := make(chan struct{})
		for _, p := range promises {
			go func(p *Promise) {
				defer safeClose(ch)
				defer func() {
					if r := recover(); r != nil {
						reject(r)
					}
				}()
				p.Await()
				resolve(p.Result())
			}(p)
		}
		<-ch
	})
}
