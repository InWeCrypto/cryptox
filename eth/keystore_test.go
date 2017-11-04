package eth

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScryptKeyStore(t *testing.T) {
	neo, err := ioutil.ReadFile("testdata/scrypt.json")

	if err != nil {
		t.Fatalf("%s", err)
	}

	key, err := ReadKeyStore(neo, "test")

	if err != nil {
		t.Fatalf("%s", err)
	}

	json, err := WriteLightScryptKeyStore(key, "test2")

	if err != nil {
		t.Fatalf("%s", err)
	}

	key2, err := ReadKeyStore(json, "test2")

	if err != nil {
		t.Fatalf("%s", err)
	}

	assert.Equal(t, key, key2)

}
