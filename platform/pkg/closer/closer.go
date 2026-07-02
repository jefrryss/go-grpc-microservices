package closer

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type closeFunc struct {
	name string
	fn   func(context.Context) error
}

type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	funcs []closeFunc
	err   error
}

func New() *Closer {
	return &Closer{}
}

func (c *Closer) Add(name string, fn func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, closeFunc{name: name, fn: fn})
}

func (c *Closer) Close(ctx context.Context) error {
	c.once.Do(func() {
		c.mu.Lock()
		funcs := append([]closeFunc(nil), c.funcs...)
		c.funcs = nil
		c.mu.Unlock()

		var result error
		for i := len(funcs) - 1; i >= 0; i-- {
			if err := funcs[i].fn(ctx); err != nil {
				result = errors.Join(result, fmt.Errorf("close %s: %w", funcs[i].name, err))
			}
		}
		c.err = result
	})

	return c.err
}
