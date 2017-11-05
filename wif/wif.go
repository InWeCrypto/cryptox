package wif

import (
	"bytes"
	"crypto/sha256"
	"fmt"

	b58 "github.com/inwecrypto/cryptox/base58"
)

// FromWIF get privatekey from wif format
func FromWIF(wif string) ([]byte, error) {
	base58 := b58.NewBase58()

	decodedWIF, err := base58.Decode(wif)
	if err != nil {
		return nil, err
	}

	if len(decodedWIF) != 38 {
		return nil, fmt.Errorf(
			"Expected length of decoded WIF to be 38, got: %d", len(decodedWIF),
		)
	}

	if decodedWIF[0] != 0x80 {
		return nil, fmt.Errorf(
			"Expected first byte of decoded WIF to be '0x80', got: %x", decodedWIF[0],
		)
	}

	if decodedWIF[33] != 0x01 {
		return nil, fmt.Errorf(
			"Expected 34th byte of decoded WIF to be '0x01', got: %x", decodedWIF[33],
		)
	}

	subString := decodedWIF[:len(decodedWIF)-4]

	rawFirstSHA := sha256.Sum256([]byte(subString))
	firstSHA := rawFirstSHA[:]

	rawSecondSHA := sha256.Sum256(firstSHA)
	secondSHA := rawSecondSHA[:]

	firstFourBytes := secondSHA[:4]
	lastFourBytes := decodedWIF[len(decodedWIF)-4 : len(decodedWIF)]

	if !bytes.Equal(firstFourBytes, lastFourBytes) {
		return nil, fmt.Errorf("WIF failed checksum validation")
	}

	return decodedWIF[1:33], nil
}
