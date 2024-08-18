package utils

import (
	"github.com/google/uuid"
	"math/big"
)

// Define the Base62 character set
const base62Charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// GenUniqueValue generates a unique value
func GenUniqueValue() string {
	// Generate a new UUID
	uuidObj := uuid.New()

	// Convert the UUID to a big.Int
	number := new(big.Int)
	number.SetBytes(uuidObj[:])

	return Base62EncodeBigInt(number)
}

// Base62EncodeBigInt encodes a big.Int to a Base62 string
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

		// number = number / base, and mod = number % base
		number.DivMod(number, base, mod)

		// prepend the character at the index of the remainder to the result
		result = string(base62Charset[mod.Int64()]) + result
	}
	return result
}
