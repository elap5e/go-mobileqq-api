package bytes

type Buffer struct {
	buf []byte
	idx int
}

// NewBuffer allocates a new Buffer initialized with buf,
// where the contents of buf are considered the unread portion of the buffer.
func NewBuffer(buf []byte) *Buffer {
	return &Buffer{buf: buf}
}

func (b *Buffer) Bytes() []byte {
	return b.buf
}