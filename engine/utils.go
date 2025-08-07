package engine

import (
	"crypto/rand"
	"math/big"
)

type ID string

func generateID() ID {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	id := make([]byte, 8)
	for i := range id {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		id[i] = charset[num.Int64()]
	}
	return ID(id)
}
