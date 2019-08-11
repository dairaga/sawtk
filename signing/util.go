package signing

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"io/ioutil"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"gitlab.com/dataforceme/sawtk/signing"
)

var (
	secp256k1Ctx = NewSecp256k1Context()
	factory      = NewCryptoFactory(secp256k1Ctx)
)

// ----------------------------------------------------------------------------

// RandomPrivateKey returns a random private key.
func RandomPrivateKey() PrivateKey {
	return secp256k1Ctx.NewRandomPrivateKey()
}

// NewSigner returns a signer from private key.
func NewSigner(key PrivateKey) *Signer {
	return factory.NewSigner(key)
}

// GenerateSigner returns a signer with a random private key.
func GenerateSigner() *Signer {
	return NewSigner(secp256k1Ctx.NewRandomPrivateKey())
}

// LoadSignerFromBytes returns a signer loading from bytes of private key.
func LoadSignerFromBytes(raw []byte) *Signer {
	return NewSigner(NewSecp256k1PrivateKey(raw))
}

// LoadSignerFromHex returns a signer from hex string.
func LoadSignerFromHex(hexstr string) (*Signer, error) {
	raw, err := hex.DecodeString(hexstr)
	if err != nil {
		return nil, err
	}

	return LoadSignerFromBytes(raw), nil
}

// LoadSignerFromFile returns signer from private key file.
func LoadSignerFromFile(file string) (*Signer, error) {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return LoadSignerFromBytes(raw), nil
}

// MustSigner ...
func MustSigner(s *Signer, err error) *Signer {
	if err != nil || s == nil {
		panic("signer failure")
	}

	return s
}

// ----------------------------------------------------------------------------

// Verify returns true if signature is legal, or false.
func Verify(signature, message, public []byte) bool {
	pubKey := NewSecp256k1PublicKey(public)

	return secp256k1Ctx.Verify(signature, message, pubKey)
}

// VerifyMsg verifies signature is legal or not.
// message will be cast to bytes.
func VerifyMsg(signature, message, publicKey string) (bool, error) {

	pubKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return false, err
	}

	signBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, err
	}

	return Verify(signBytes, []byte(message), pubKeyBytes), nil
}

// VerifyHex verifies signature is legal or not.
// all inputs must be hex string.
func VerifyHex(signature, message, publicKey string) (bool, error) {
	pubKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return false, err
	}

	signBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, err
	}

	msgBytes, err := hex.DecodeString(message)
	if err != nil {
		return false, err
	}

	return Verify(signBytes, msgBytes, pubKeyBytes), nil
}

// ----------------------------------------------------------------------------

var one = new(big.Int).SetInt64(1)
var curv = btcec.S256()

func fill(dest, src []byte) {
	dlen := len(dest)
	slen := len(src)

	if dlen <= slen {
		copy(dest, src)
		return
	}

	count := 0
	for {
		if count >= dlen {
			return
		}

		copy(dest[count:], src)
		count += slen
	}
}

func codeFieldElement(c elliptic.Curve, code string) *big.Int {
	params := c.Params()
	b := make([]byte, params.BitSize/8+8)
	fill(b, []byte(code))
	k := new(big.Int).SetBytes(b)
	n := new(big.Int).Sub(params.N, one)
	k.Mod(k, n)
	k.Add(k, one)
	return k
}

func generateKeyByCode(c elliptic.Curve, code string) *ecdsa.PrivateKey {
	k := codeFieldElement(c, code)

	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = c
	priv.D = k
	priv.PublicKey.X, priv.PublicKey.Y = c.ScalarBaseMult(k.Bytes())
	return priv
}

// GeneratePrivKeyFromCode returns a private key from some code.
func GeneratePrivKeyFromCode(code string) PrivateKey {
	p := (*btcec.PrivateKey)(generateKeyByCode(curv, code))
	return signing.NewSecp256k1PrivateKey(p.Serialize())
}

// GenerateSignerFromCode returns a signer from some code.
func GenerateSignerFromCode(code string) *Signer {
	return factory.NewSigner(GeneratePrivKeyFromCode(code))
}
