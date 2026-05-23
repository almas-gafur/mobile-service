package hash

import (
	"crypto/rand"
	"math/big"
)

const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789"

func Generate(length int) (string, error) {
	result := make([]byte, length)
	max := big.NewInt(int64(len(alphabet)))

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		result[i] = alphabet[n.Int64()]
	}

	return string(result), nil
}
