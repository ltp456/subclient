package subclient

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/ChainSafe/gossamer/lib/crypto"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/blake2b"
	"hash"
	"math/big"
	"strings"
	"subclient/types"
	"subclient/xxhash"
)

func NewStorageKey(name types.ModuleName, key types.StorageKey, networkId []byte, opts ...types.Option) (string, error) {
	moduleHash := xxhash.New128([]byte(name)).Sum(nil)
	storageHash := xxhash.New128([]byte(key)).Sum(nil)
	keyBytes := append(moduleHash, storageHash...)
	for _, opt := range opts {
		value := opt.Value
		switch opt.Type {
		case types.Twox64Concat:
			hashSum := xxhash.New64Concat(nil)
			keyHash, err := getHashValue(value, hashSum, networkId)
			if err != nil {
				return "", err
			}
			keyBytes = append(keyBytes, keyHash...)
		case types.Blake2_128Concat:
			hashSum, err := xxhash.NewBlake2b128Concat(nil)
			if err != nil {
				return "", err
			}
			keyHash, err := getHashValue(value, hashSum, networkId)
			if err != nil {
				return "", err
			}
			keyBytes = append(keyBytes, keyHash...)
		default:
			return "", fmt.Errorf("unSupport type: %v", opt.Type)
		}

	}
	return fmt.Sprintf("%x", keyBytes), nil

}

func getHashValue(value interface{}, hash hash.Hash, networkId []byte) ([]byte, error) {
	switch value.(type) {
	case string: // only addr, todo alias
		address := value.(string)
		public, err := AddrToPublic(address, networkId)
		if err != nil {
			return nil, err
		}
		hash.Write(public)
		return hash.Sum(nil), nil
	case uint32:
		index := value.(uint32)
		eraIndexSerialized := make([]byte, 4)
		binary.LittleEndian.PutUint32(eraIndexSerialized, index)
		hash.Write(eraIndexSerialized)
		paramHash := hash.Sum(nil)
		return paramHash, nil
	default:
		return nil, fmt.Errorf("unsupport: %v", value)
	}

}

func AddrToPublic(addr string, networkId []byte) ([]byte, error) {
	bytes := base58.Decode(addr)
	enc := bytes[:len(bytes)-2]
	pub := enc[len(networkId):]
	return pub, nil
}

var ss58Prefix = []byte("SS58PRE")

func PublicKeyToAddress(pub crypto.PublicKey, networkId []byte) (string, error) {
	enc := append(networkId, pub.Encode()...)
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

func TrimSlash(value string) string {
	tValue := strings.TrimLeft(value, "\"")
	result := strings.TrimRight(tValue, "\"")
	return result
}

func TrimQuotes(value string) string {
	key := fmt.Sprintf("%s", strings.ReplaceAll(value, `"`, ""))
	return key
}

func Hash256(data []byte) ([]byte, error) {
	checksum, _ := blake2b.New(32, []byte{})
	_, _ = checksum.Write(data)
	h := checksum.Sum(nil)
	return h, nil
}

func HashHex(data string) (string, error) {
	checksum, _ := blake2b.New(32, []byte{})
	_, _ = checksum.Write([]byte(data))
	h := checksum.Sum(nil)
	return hex.EncodeToString(h), nil
}

func RatToStr(value *big.Rat, decimal int) string {
	//todo
	if decimal == 0 {
		return value.FloatString(0)
	}
	precision := "1" + strings.Repeat("0", decimal)
	precisionRat, _ := NewRatFromStr(precision)
	ratDiv := RatDiv(value, precisionRat)
	return ratDiv.FloatString(decimal)

}

func StrToRat(value string, decimal int) (*big.Rat, bool) {
	//todo
	rat, ok := NewRatFromStr(value)
	if !ok {
		return nil, false
	}
	if decimal == 0 {
		return rat, true
	}
	precision := "1" + strings.Repeat("0", decimal)
	precisionRat, ok := NewRatFromStr(precision)
	if !ok {
		return nil, false
	}
	ratMul := RatMul(rat, precisionRat)
	return ratMul, true
}

func RatDivDecimal(value *big.Rat, decimal int) *big.Rat {
	if decimal == 0 {
		return new(big.Rat).Set(value)
	}
	precision := "1" + strings.Repeat("0", int(decimal))
	precisionRat, _ := NewRatFromStr(precision)
	ratDiv := RatDiv(value, precisionRat)
	return ratDiv

}

func NewRateFromBigInt(value *big.Int) *big.Rat {
	rat := new(big.Rat)
	newRate := rat.SetInt(value)
	return newRate
}

func NewRatFromStr(value string) (*big.Rat, bool) {
	rat := new(big.Rat)
	newRat, ok := rat.SetString(value)
	return newRat, ok
}

func NewRatFromBigInt(value *big.Int) (*big.Rat, bool) {
	return NewRatFromStr(value.String())
}

func NewRatFromFloat(value float64) *big.Rat {
	rat := new(big.Rat)
	newRat := rat.SetFloat64(value)
	return newRat
}

func NewRatFromInt(value int64) *big.Rat {
	rat := new(big.Rat)
	newRat := rat.SetInt64(value)
	return newRat
}

func NewRatFromUint(value uint64) *big.Rat {
	rat := new(big.Rat)
	newRat := rat.SetUint64(value)
	return newRat
}

func RatAdd(a, b *big.Rat) *big.Rat {
	rat := new(big.Rat)
	newRat := rat.Add(a, b)
	return newRat
}

func RatSub(a, b *big.Rat) *big.Rat {
	rat := new(big.Rat)
	newRat := rat.Sub(a, b)
	return newRat
}

func RatMul(a, b *big.Rat) *big.Rat {
	rat := new(big.Rat)
	newRat := rat.Mul(a, b)
	return newRat
}

func NewRat(a *big.Rat) *big.Rat {
	rat := new(big.Rat)
	rat.Set(a)
	return rat
}

func RatDiv(a, b *big.Rat) *big.Rat {
	rat := new(big.Rat)
	invRat := rat.Inv(b)
	tmpRat := new(big.Rat)
	result := tmpRat.Mul(a, invRat)
	return result
}
func Big2Str(x *big.Int, decimal int) string {
	if decimal == 0 {
		return x.String()
	}

	// Integral part
	i := x.String()
	if len(i) <= decimal {
		i = "0"
	} else {
		i = i[0 : len(i)-decimal]
	}

	// Decimal part
	d := x.String()
	if len(d) < decimal {
		d = strings.Repeat("0", decimal-len(d)) + d
	} else {
		d = d[len(d)-decimal:]
	}

	return i + "." + d
}

// ParseBigInt parse hex string value to big.Int
func ParseBigInt(value string) (*big.Int, error) {
	i := big.NewInt(0)
	_, err := fmt.Sscan(value, i)

	return i, err
}
