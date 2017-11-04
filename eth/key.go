// Package eth the eth crypto library
package eth

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/pborman/uuid"

	"github.com/inwecrypto/cryptox/keystore"
)

// const variables
var (
	StandardScryptN = 1 << 18
	StandardScryptP = 1
	LightScryptN    = 1 << 12
	LightScryptP    = 6
)

// Key eth wallet key
type Key struct {
	ID         uuid.UUID // Key ID
	Address    string    // address
	PrivateKey *ecdsa.PrivateKey
}

func keystoreKeyToEthKey(key *keystore.Key) (*Key, error) {

	ecdsaKey, err := toECDSA(key.PrivateKey, false)

	if err != nil {
		return nil, err
	}

	return &Key{
		ID:         uuid.UUID(key.ID),
		Address:    key.Address,
		PrivateKey: ecdsaKey,
	}, nil
}

func ethKeyToKeyStoreKey(key *Key) (*keystore.Key, error) {
	bytes := key.PrivateKey.D.Bytes()

	return &keystore.Key{
		ID:         key.ID,
		Address:    key.Address,
		PrivateKey: bytes,
	}, nil
}

// WriteScryptKeyStore write keystore with Scrypt format
func WriteScryptKeyStore(key *Key, password string) ([]byte, error) {
	keyStoreKey, err := ethKeyToKeyStoreKey(key)

	if err != nil {
		return nil, err
	}

	attrs := map[string]interface{}{
		"ScryptN": StandardScryptN,
		"ScryptP": StandardScryptP,
	}

	return keystore.Encrypt(keyStoreKey, password, attrs)
}

// WriteLightScryptKeyStore write keystore with Scrypt format
func WriteLightScryptKeyStore(key *Key, password string) ([]byte, error) {
	keyStoreKey, err := ethKeyToKeyStoreKey(key)

	if err != nil {
		return nil, err
	}

	attrs := map[string]interface{}{
		"ScryptN": LightScryptN,
		"ScryptP": LightScryptP,
	}

	return keystore.Encrypt(keyStoreKey, password, attrs)
}

// ReadKeyStore read key from keystore
func ReadKeyStore(data []byte, password string) (*Key, error) {
	keystore, err := keystore.Decrypt(data, password)

	if err != nil {
		return nil, err
	}

	return keystoreKeyToEthKey(keystore)
}

func toECDSA(d []byte, strict bool) (*ecdsa.PrivateKey, error) {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = secp256k1.S256()

	if strict && 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}

	priv.D = new(big.Int).SetBytes(d)

	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)

	return priv, nil
}
