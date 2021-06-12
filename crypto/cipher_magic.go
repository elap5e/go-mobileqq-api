package crypto

import (
	_ "crypto"
	_ "unsafe"
)

//go:linkname xorBytes crypto/cipher.xorBytes
func xorBytes(dst, a, b []byte) int
