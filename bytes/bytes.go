package bytes

type Encoder interface {
	Encode(b *Buffer)
}

type Decoder interface {
	Decode(b *Buffer) error
}
