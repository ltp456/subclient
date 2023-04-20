package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/blake2b"
	"math/big"
	"reflect"
	"strings"
)

var ss58Prefix = []byte("SS58PRE")

func PublicKeyToAddress(pub []byte, networkId []byte) (string, error) {
	enc := append(networkId, pub...)
	hasher, err := blake2b.New(64, nil)
	if err != nil {
		return "", err
	}
	_, err = hasher.Write(append(ss58Prefix, enc...))
	if err != nil {
		return "", err
	}
	checksum := hasher.Sum(nil)
	return base58.Encode(append(enc, checksum[:2]...)), nil
}

func AddrToPublic(addr string, networkId []byte) ([]byte, error) {
	bytes := base58.Decode(addr)
	enc := bytes[:len(bytes)-2]
	pub := enc[len(networkId):]
	return pub, nil
}

func ObjToBytes(obj interface{}) ([][]byte, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	var value [][]byte
	err = Unmarshal(bytes, &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}
func ObjToBigInt(obj interface{}) (*big.Int, error) {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	value := big.NewInt(0)
	err = Unmarshal(bytes, &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}
func ObjToObj(obj interface{}, value interface{}) error {
	typeOf := reflect.ValueOf(value)
	if typeOf.Kind() != reflect.Ptr {
		return fmt.Errorf("value is mutst pointer")
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return Unmarshal(data, value)
}

func Marshal(obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encode := json.NewEncoder(&buf)
	err := encode.Encode(obj)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func Unmarshal(data []byte, value interface{}) error {
	typeOf := reflect.ValueOf(value)
	if typeOf.Kind() != reflect.Ptr {
		return fmt.Errorf("value is mutst pointer")
	}
	decode := json.NewDecoder(strings.NewReader(string(data)))
	decode.UseNumber()
	err := decode.Decode(value)
	if err != nil {
		return err
	}
	return err
}
