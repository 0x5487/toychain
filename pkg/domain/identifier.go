package domain

import (
	"encoding/hex"
)

type Identifier [32]byte

func (i Identifier) ToByte() []byte {
	return i[:]
}

func (i Identifier) String() string {
	return hex.EncodeToString(i.ToByte())
}
