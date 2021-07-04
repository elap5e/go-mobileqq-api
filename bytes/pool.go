package bytes

import (
	"sync"
)

type Pool struct {
	pool sync.Pool
}

func NewPool(length int) *Pool {
	return &Pool{
		pool: sync.Pool{
			New: func() interface{} {
				var r []byte
				if length > 0 {
					r = make([]byte, 0, length)
				}
				return &Buffer{buf: r}
			},
		},
	}
}

func (b *Pool) Put(buf *Buffer) {
	b.pool.Put(buf)
}

func (b *Pool) Get() *Buffer {
	buf := b.pool.Get().(*Buffer)
	buf.Reset()
	return buf
}
