package neo

import (
	"bytes"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"math"
	"sort"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/goany/slf4go"
	"github.com/inwecrypto/neogo"
)

var logger = slf4go.Get("neo")

// Asserts .
const (
	GasAssert = "602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7"
	NEOAssert = "c56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b"
)

// Err
var (
	ErrNoUTXO = errors.New("no enough utxo")
)

// Transaction types
const (
	MinerTransaction      byte = 0x00
	IssueTransaction      byte = 0x01
	ClaimTransaction      byte = 0x02
	EnrollmentTransaction byte = 0x20
	RegisterTransaction   byte = 0x40
	ContractTransaction   byte = 0x80
	PublishTransaction    byte = 0xd0
	InvocationTransaction byte = 0xd1
)

// Attr Usage
const (
	ContractHash   = byte(0x00)
	ECDH02         = byte(0x02)
	ECDH03         = byte(0x03)
	Script         = byte(0x20)
	Vote           = byte(0x30)
	CertURL        = byte(0x80)
	DescriptionURL = byte(0x81)
	Description    = byte(90)
	Hash1          = byte(0xa1)
	Hash2          = byte(0xa2)
	Hash3          = byte(0xa3)
	Hash4          = byte(0xa4)
	Hash5          = byte(0xa5)
	Hash6          = byte(0xa6)
	Hash7          = byte(0xa7)
	Hash8          = byte(0xa8)
	Hash9          = byte(0xa9)
	Hash10         = byte(0xaa)
	Hash11         = byte(0xab)
	Hash12         = byte(0xac)
	Hash13         = byte(0xad)
	Hash14         = byte(0xae)
	Hash15         = byte(0xaf)
	Remark         = byte(0xf0)
	Remark1        = byte(0xf1)
	Remark2        = byte(0xf2)
	Remark3        = byte(0xf3)
	Remark4        = byte(0xf4)
	Remark5        = byte(0xf5)
	Remark6        = byte(0xf6)
	Remark7        = byte(0xf7)
	Remark8        = byte(0xf8)
	Remark9        = byte(0xf9)
	Remark10       = byte(0xfa)
	Remark11       = byte(0xfb)
	Remark12       = byte(0xfc)
	Remark13       = byte(0xfd)
	Remark14       = byte(0xfe)
	Remark15       = byte(0xff)
)

// RawTxSerializable .
type RawTxSerializable func(writer io.Writer) error

// RawTx raw transaction object
type RawTx struct {
	Type       byte              // transaction type
	Version    byte              // tx version
	XData      RawTxSerializable // special tx data
	Attributes []*RawTxAttr      // tx attribute
	Inputs     []*RawTxInput     // tx inputs
	Outputs    []*RawTxOutput    // tx output
	Scripts    []*RawTxScript    // tx scripts
}

// NewRawTx create new raw tx
func NewRawTx(txType byte) *RawTx {
	return &RawTx{
		Type:    txType,
		Version: 0,
	}
}

// GenerateWithSign generate raw tx with sign data
func (tx *RawTx) GenerateWithSign(key *Key) ([]byte, string, error) {

	var buff bytes.Buffer

	if err := tx.writeSignData(&buff); err != nil {
		return nil, "", err
	}

	txid := sha256.Sum256(buff.Bytes())

	txid = sha256.Sum256(txid[:])

	sign, err := key.PrivateKey.Sign(buff.Bytes(), elliptic.P256())

	if err != nil {
		return nil, "", err
	}

	redeemScript := key.PrivateKey.PublicKey.ToBytes()

	tx.Scripts = []*RawTxScript{
		&RawTxScript{
			StackScript:  sign,
			RedeemScript: redeemScript,
		},
	}

	buff.Reset()

	err = tx.WriteBytes(&buff)

	if err != nil {
		return nil, "", err
	}

	return buff.Bytes(), hex.EncodeToString(reverseBytes(txid[:])), nil
}

func (tx *RawTx) writeSignData(writer io.Writer) error {
	_, err := writer.Write([]byte{tx.Type, 0x00})

	if err != nil {
		return err
	}

	if tx.XData != nil {
		if err := tx.XData(writer); err != nil {
			return err
		}
	}

	_, err = writer.Write([]byte{byte(len(tx.Attributes))})

	if err != nil {
		return err
	}

	for _, attr := range tx.Attributes {
		if err := attr.WriteBytes(writer); err != nil {
			return err
		}
	}

	_, err = writer.Write([]byte{byte(len(tx.Inputs))})

	if err != nil {
		return err
	}

	for _, input := range tx.Inputs {
		if err := input.WriteBytes(writer); err != nil {
			return err
		}
	}

	_, err = writer.Write([]byte{byte(len(tx.Outputs))})

	if err != nil {
		return err
	}

	for _, output := range tx.Outputs {
		if err := output.WriteBytes(writer); err != nil {
			return err
		}
	}

	return nil
}

// WriteBytes .
func (tx *RawTx) WriteBytes(writer io.Writer) error {

	if err := tx.writeSignData(writer); err != nil {
		return err
	}

	_, err := writer.Write([]byte{byte(len(tx.Scripts))})

	if err != nil {
		return err
	}

	for _, script := range tx.Scripts {
		if err := script.WriteBytes(writer); err != nil {
			return err
		}
	}

	return nil
}

// RawTxAttr raw transaction attribute
type RawTxAttr struct {
	Usage byte
	Data  []byte
}

// WriteBytes .
func (attr *RawTxAttr) WriteBytes(writer io.Writer) error {

	_, err := writer.Write([]byte{attr.Usage})

	if err != nil {
		return err
	}

	if !(attr.Usage <= ECDH03 || attr.Usage == Vote || (attr.Usage <= Hash15 && attr.Usage >= Hash1)) {
		_, err := writer.Write([]byte{byte(len(attr.Data))})

		if err != nil {
			return err
		}
	}

	_, err = writer.Write(attr.Data)

	if err != nil {
		return err
	}

	return nil
}

// RawTxInput raw tx input parameter
type RawTxInput struct {
	TxID string
	Vout uint16
}

// WriteBytes .
func (input *RawTxInput) WriteBytes(writer io.Writer) error {

	data, err := hex.DecodeString(strings.TrimPrefix(input.TxID, "0x"))

	if err != nil {
		return err
	}

	_, err = writer.Write(reverseBytes(data))

	if err != nil {
		return err
	}

	data = make([]byte, 2)

	binary.LittleEndian.PutUint16(data, input.Vout)

	_, err = writer.Write(data)

	if err != nil {
		return err
	}

	return nil
}

// RawTxOutput raw tx output utxo
type RawTxOutput struct {
	AssertID string
	Value    float64
	Address  string
}

// WriteBytes .
func (output *RawTxOutput) WriteBytes(writer io.Writer) error {

	data, err := hex.DecodeString(strings.TrimPrefix(output.AssertID, "0x"))

	if err != nil {
		return err
	}

	_, err = writer.Write(reverseBytes(data))

	if err != nil {
		return err
	}

	value := uint64(math.Floor(output.Value * 100000000))

	data = make([]byte, 8)

	binary.LittleEndian.PutUint64(data, value)

	_, err = writer.Write(data)

	if err != nil {
		return err
	}

	data, err = decodeAddress(output.Address)

	if err != nil {
		return err
	}

	_, err = writer.Write(data)

	if err != nil {
		return err
	}

	return nil
}

// RawTxScript .
type RawTxScript struct {
	StackScript  []byte
	RedeemScript []byte
}

// WriteBytes .
func (script *RawTxScript) WriteBytes(writer io.Writer) error {

	length := byte(len(script.StackScript))

	_, err := writer.Write([]byte{length + 1, length})

	if err != nil {
		return err
	}

	_, err = writer.Write(script.StackScript)

	if err != nil {
		return err
	}

	length = byte(len(script.RedeemScript))

	_, err = writer.Write([]byte{length + 2, length})

	if err != nil {
		return err
	}

	_, err = writer.Write(script.RedeemScript)

	if err != nil {
		return err
	}

	_, err = writer.Write([]byte{0xac})

	if err != nil {
		return err
	}

	return nil
}

// RawClaimTx .
type RawClaimTx struct {
	*RawTx
	Claims []*RawTxInput
}

// NewRawClaimTx .
func NewRawClaimTx() *RawClaimTx {
	tx := &RawClaimTx{
		RawTx: NewRawTx(ClaimTransaction),
	}

	tx.RawTx.XData = func(writer io.Writer) error {

		logger.DebugF("======%x", len(tx.Claims))

		_, err := writer.Write([]byte{byte(len(tx.Claims))})

		if err != nil {
			return err
		}

		for _, clamin := range tx.Claims {
			if err := clamin.WriteBytes(writer); err != nil {
				return err
			}
		}

		return nil
	}

	return tx
}

func reverseBytes(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

func decodeAddress(address string) ([]byte, error) {

	result, _, err := base58.CheckDecode(address)

	if err != nil {
		logger.DebugF("decode address :%s -- failed\n\t%s", address, err)
		return nil, err
	}

	return result[0:20], nil
}

type utxoSorter []*neogo.UTXO

func (s utxoSorter) Len() int      { return len(s) }
func (s utxoSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s utxoSorter) Less(i, j int) bool {

	ival, _ := s[i].Value()
	jval, _ := s[j].Value()

	return ival < jval
}

// CalcTxInput .
func CalcTxInput(amount float64, unspent []*neogo.UTXO) ([]*neogo.UTXO, float64, error) {
	sort.Sort(utxoSorter(unspent))

	selected := make([]*neogo.UTXO, 0)
	vinvalue := float64(0)

	for _, utxo := range unspent {
		var err error
		selected = append(selected, utxo)
		vinvalue, err = utxo.Value()

		if err != nil {
			return nil, 0, err
		}

		if vinvalue > amount {
			return selected, vinvalue, nil
		}
	}

	return selected, vinvalue, nil
}

// CreateSendAssertTx create send assert tx object
func CreateSendAssertTx(assert, from, to string, amount float64, unspent []*neogo.UTXO) (*RawTx, error) {

	sendUTXOs, totalAmount, err := CalcTxInput(amount, unspent)

	if err != nil {
		return nil, err
	}

	if totalAmount < amount {
		return nil, ErrNoUTXO
	}

	tx := NewRawTx(ContractTransaction)

	for _, utxo := range sendUTXOs {
		tx.Inputs = append(tx.Inputs, &RawTxInput{
			TxID: utxo.TransactionID,
			Vout: uint16(utxo.Vout.N),
		})
	}

	tx.Outputs = append(tx.Outputs, &RawTxOutput{
		AssertID: assert,
		Value:    amount,
		Address:  to,
	})

	if totalAmount > amount {
		tx.Outputs = append(tx.Outputs, &RawTxOutput{
			AssertID: assert,
			Value:    totalAmount - amount,
			Address:  from,
		})
	}

	return tx, nil
}

type claimSorter []*neogo.UTXO

func (s claimSorter) Len() int      { return len(s) }
func (s claimSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s claimSorter) Less(i, j int) bool {
	return s[i].SpentBlock < s[j].SpentBlock
}

// CreateClaimTx .
func CreateClaimTx(val float64, address string, unspent []*neogo.UTXO) (*RawTx, error) {
	tx := NewRawClaimTx()

	sort.Sort(claimSorter(unspent))

	for _, utxo := range unspent {
		tx.Claims = append(tx.Claims, &RawTxInput{
			TxID: utxo.TransactionID,
			Vout: uint16(utxo.Vout.N),
		})
	}

	tx.Outputs = append(tx.Outputs, &RawTxOutput{
		AssertID: GasAssert,
		Value:    val,
		Address:  address,
	})

	return tx.RawTx, nil
}
