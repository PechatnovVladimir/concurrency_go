package concurrency

type Promise struct {
	result   chan error
	promised bool
}

func NewPromise() Promise {
	return Promise{
		result: make(chan error, 1),
	}
}

func (p *Promise) Set(value error) {
	if p.promised {
		return
	}

	p.promised = true
	p.result <- value
	close(p.result)
}

type Future struct {
	result <-chan error
}

func NewFuture(result <-chan error) Future {
	return Future{
		result: result,
	}
}

func (f *Future) Get() error {
	return <-f.result
}

func (p *Promise) GetFuture() Future {
	return NewFuture(p.result)
}
