package neo

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/dynamicgo/config"
	"github.com/inwecrypto/neogo"
	"github.com/stretchr/testify/assert"
)

var cnf *config.Config

func init() {
	cnf, _ = config.NewFromFile("./test.json")
}

func TestType(t *testing.T) {
	assert.Equal(t, (ClaimTransaction), byte(0x02))
	assert.Equal(t, (ContractTransaction), byte(0x80))
}

func TestSign(t *testing.T) {

	client := neogo.NewClient(cnf.GetString("testnode", "xxxxx") + "/extend")

	key, err := KeyFromWIF(cnf.GetString("wallet", "xxxxx"))

	assert.NoError(t, err)

	utxos, err := client.GetBalance(key.Address, "0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b")

	assert.NoError(t, err)

	printResult(utxos)

	tx, err := CreateSendAssertTx(
		"0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b",
		key.Address,
		key.Address,
		1, utxos)

	if !assert.NoError(t, err) {

	}

	rawtx, txid, err := tx.GenerateWithSign(key)

	assert.NoError(t, err)

	logger.DebugF("txid %s with raw data:\n%s", txid, hex.EncodeToString(rawtx))

	client = neogo.NewClient(cnf.GetString("testnode", "xxxxx"))

	status, err := client.SendRawTransaction(rawtx)

	assert.NoError(t, err)

	println(status)
}

func printResult(result interface{}) {

	data, _ := json.MarshalIndent(result, "", "\t")

	fmt.Println(string(data))
}

func TestDecodeAddress(t *testing.T) {
	address, err := decodeAddress("AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr")

	assert.NoError(t, err)

	logger.Debug(hex.EncodeToString(address))
}
