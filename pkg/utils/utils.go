package utils

import (
	"github.com/google/uuid"
	"math/big"
)

// Define the Base62 character set
const base62Charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenUniqueValue() string {
	uuidObj := uuid.New()

	// Convert the UUID to a big.Int
	number := new(big.Int)
	number.SetBytes(uuidObj[:])

	return Base62EncodeBigInt(number)
}

func Base62EncodeBigInt(number *big.Int) string {
	// compare number with 0
	zero := big.NewInt(0)
	if number.Cmp(zero) == 0 {
		return string(base62Charset[0])
	}

	var result string
	base := big.NewInt(int64(len(base62Charset)))
	for number.Cmp(zero) > 0 {
		// mod is an empty big.Int, which holds the remainder of the division operation (number % base)
		mod := new(big.Int)
		number.DivMod(number, base, mod)
		result = string(base62Charset[mod.Int64()]) + result
	}
	return result
}
