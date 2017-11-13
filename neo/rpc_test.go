package neo

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRPCAccountSate(t *testing.T) {
	client := NewClient("http://47.52.173.179:20332")

	accoutState, err := client.GetAccountState("AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr")

	assert.NoError(t, err)

	printResult(accoutState)
}

func TestGetBalance(t *testing.T) {
	client := NewClient("http://47.52.173.179:20332")

	balance, err := client.GetBalance("0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b")

	assert.NoError(t, err)

	printResult(balance)
}

func TestConnectionCount(t *testing.T) {
	client := NewClient("http://47.52.173.179:20332")

	count, err := client.GetConnectionCount()

	assert.NoError(t, err)

	fmt.Printf("connection count :%d\n", count)
}

func TestBestBlockHash(t *testing.T) {
	client := NewClient("http://47.52.173.179:20332")

	hash, err := client.GetBestBlockHash()

	assert.NoError(t, err)

	block, err := client.GetBlock(hash)

	assert.NoError(t, err)

	blockjson, _ := json.MarshalIndent(block, "", "\t")

	fmt.Printf("the best block :\n\t%s\n", string(blockjson))
}

func TestBlockCount(t *testing.T) {
	client := NewClient("http://47.52.173.179:20332")

	count, err := client.GetBlockCount()

	assert.NoError(t, err)

	fmt.Printf("the block count :%d\n", count)
}

func TestBlockByIndex(t *testing.T) {
	client := NewClient("http://47.52.173.179:10332")

	block, err := client.GetBlockByIndex(1546852)

	assert.NoError(t, err)

	blockjson, _ := json.MarshalIndent(block, "", "\t")

	fmt.Printf("the best block :\n\t%s\n", string(blockjson))
}

func TestGetRawTransaction(t *testing.T) {
	client := NewClient("http://47.52.173.179:10332")

	block, err := client.GetRawTransaction("0x8e977f49006bf768dc80f1938b0bf9536478b1d9c206686728de2020ea1ab259")

	assert.NoError(t, err)

	blockjson, _ := json.MarshalIndent(block, "", "\t")

	fmt.Printf("trans:\n\t%s\n", string(blockjson))
}

func TestGetTxOut(t *testing.T) {
	client := NewClient("http://47.52.173.179:10332")

	block, err := client.GetTxOut("0x0ae13c1ba01d30a8238a0ec89019171fcf9eee61802dd468cc797a02ac48798d", 0)

	assert.NoError(t, err)

	blockjson, _ := json.MarshalIndent(block, "", "\t")

	fmt.Printf("trans:\n\t%s\n", string(blockjson))
}

func TestGetPeers(t *testing.T) {
	client := NewClient("http://47.52.173.179:10332")

	block, err := client.GetPeers()

	assert.NoError(t, err)

	blockjson, _ := json.MarshalIndent(block, "", "\t")

	fmt.Printf("peers:\n\t%s\n", string(blockjson))
}

func TestSendRawTransaction(t *testing.T) {

	wallet, err := KeyFromWIF("L4Ns4Uh4WegsHxgDG49hohAYxuhj41hhxG6owjjTWg95GSrRRbLL")

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, wallet.Address, "AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr")

	client := NewClient("http://47.52.173.179:20332")

	client.GetAccountState("AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr")
}

func printResult(result interface{}) {

	data, _ := json.MarshalIndent(result, "", "\t")

	fmt.Println(string(data))
}
