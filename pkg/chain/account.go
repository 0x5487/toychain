package chain

import (
	"crypto/sha256"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/shengdoushi/base58"
)

type Address struct {
	Address    string
	PrivateKey string
}

func CreateAddress() (Address, error) {
	addr := Address{}

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return addr, err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	addr.PrivateKey = hexutil.Encode(privateKeyBytes)[2:]

	publicAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	addressTron := make([]byte, 0)
	addressTron = append(addressTron, byte(0x41))
	addressTron = append(addressTron, publicAddress.Bytes()...)

	if addressTron[0] == 0 {
		addr.Address = new(big.Int).SetBytes(addressTron).String()
	}
	addr.Address = EncodeCheck(addressTron)

	return addr, nil
}

func EncodeCheck(input []byte) string {
	h256h0 := sha256.New()
	h256h0.Write(input)
	h0 := h256h0.Sum(nil)

	h256h1 := sha256.New()
	h256h1.Write(h0)
	h1 := h256h1.Sum(nil)

	inputCheck := input
	inputCheck = append(inputCheck, h1[:4]...)

	return Encode(inputCheck)
}

func Encode(input []byte) string {
	return base58.Encode(input, base58.BitcoinAlphabet)
}
