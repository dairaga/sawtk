/**
 * Copyright 2017 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * ------------------------------------------------------------------------------
 */

/*
Package signing is from https://github.com/hyperledger/sawtooth-sdk-go
*/
package signing

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"math/big"

	ellcurv "github.com/btcsuite/btcd/btcec"
)

var cachedCurve = ellcurv.S256()

// -- Private Key --

// Secp256k1PrivateKey represents a secp256k1 private key.
type Secp256k1PrivateKey struct {
	privateKey []byte
}

// NewSecp256k1PrivateKey creates a PrivateKey instance from private key bytes.
func NewSecp256k1PrivateKey(privateKey []byte) PrivateKey {
	return &Secp256k1PrivateKey{privateKey}
}

// GetAlgorithmName returns the string "secp256k1".
func (key *Secp256k1PrivateKey) GetAlgorithmName() string {
	return "secp256k1"
}

// AsHex returns the private key as a hex-encoded string.
func (key *Secp256k1PrivateKey) AsHex() string {
	return hex.EncodeToString(key.privateKey)
}

// AsBytes returns the bytes of the private key.
func (key *Secp256k1PrivateKey) AsBytes() []byte {
	return key.privateKey
}

// -- Public Key --

// Secp256k1PublicKey represents secp256k1 public key.
type Secp256k1PublicKey struct {
	publicKey []byte
}

// NewSecp256k1PublicKey creates a PublicKey instance from public key bytes.
func NewSecp256k1PublicKey(publicKey []byte) PublicKey {
	return &Secp256k1PublicKey{publicKey}
}

// GetAlgorithmName returns the string "secp256k1".
func (key *Secp256k1PublicKey) GetAlgorithmName() string {
	return "secp256k1"
}

// AsHex returns the public key as a hex-encoded string.
func (key *Secp256k1PublicKey) AsHex() string {
	return hex.EncodeToString(key.publicKey)
}

// AsBytes returns the bytes of the public key.
func (key *Secp256k1PublicKey) AsBytes() []byte {
	return key.publicKey
}

// -- Context --

// Secp256k1Context implements signing.Context with KoblitzCurve.
type Secp256k1Context struct {
	curve *ellcurv.KoblitzCurve
}

// NewSecp256k1Context returns a new secp256k1 context.
func NewSecp256k1Context() Context {
	return &Secp256k1Context{ellcurv.S256()}
}

// GetAlgorithmName returns the string "secp256k1".
func (ctx *Secp256k1Context) GetAlgorithmName() string {
	return "secp256k1"
}

// NewRandomPrivateKey generates a new random secp256k1 private key.
func (ctx *Secp256k1Context) NewRandomPrivateKey() PrivateKey {
	priv, _ := ellcurv.NewPrivateKey(cachedCurve)

	return &Secp256k1PrivateKey{priv.Serialize()}
}

// GetPublicKey produces a public key for the given private key.
func (ctx *Secp256k1Context) GetPublicKey(privateKey PrivateKey) PublicKey {
	_, publicKey := ellcurv.PrivKeyFromBytes(
		cachedCurve,
		privateKey.AsBytes())

	return NewSecp256k1PublicKey(publicKey.SerializeCompressed())
}

// Sign uses the given private key to calculate a signature for the
// given data. A sha256 hash of the data is first calculated and this
// is what is actually signed. Returns the signature as bytes using
// the compact serialization (which is just (r, s)).
func (ctx *Secp256k1Context) Sign(message []byte, privateKey PrivateKey) []byte {
	priv, _ := ellcurv.PrivKeyFromBytes(
		ctx.curve,
		privateKey.AsBytes())

	hash := doSHA256(message)

	sig, err := priv.Sign(hash)
	if err != nil {
		panic("Signing failed")
	}

	return serializeCompact(sig)
}

// Verify uses the given public key to verify that the given signature
// was created from the given data using the associated private key. A
// sha256 hash of the data is calculated first and this is what is
// actually used to verify the signature.
func (ctx *Secp256k1Context) Verify(signature []byte, message []byte, publicKey PublicKey) bool {
	sig := deserializeCompact(signature)
	hash := doSHA256(message)

	pub, err := ellcurv.ParsePubKey(
		publicKey.AsBytes(),
		ctx.curve)
	if err != nil {
		panic(err.Error())
	}

	return sig.Verify(hash, pub)
}

// -- SHA --

func doSHA512(input []byte) []byte {
	hash := sha512.New()
	hash.Write(input)
	return hash.Sum(nil)
}

func doSHA256(input []byte) []byte {
	hash := sha256.New()
	hash.Write(input)
	return hash.Sum(nil)
}

// ---

func serializeCompact(sig *ellcurv.Signature) []byte {
	b := make([]byte, 0, 64)
	// TODO: Padding
	rbytes := pad(sig.R.Bytes(), 32)
	sbytes := pad(sig.S.Bytes(), 32)
	b = append(b, rbytes...)
	b = append(b, sbytes...)
	if len(b) != 64 {
		panic("Invalid signature length")
	}
	return b
}

func deserializeCompact(b []byte) *ellcurv.Signature {
	return &ellcurv.Signature{
		R: new(big.Int).SetBytes(b[:32]),
		S: new(big.Int).SetBytes(b[32:]),
	}
}

func pad(buf []byte, size int) []byte {
	newbuf := make([]byte, 0, size)
	padLength := size - len(buf)
	for i := 0; i < padLength; i++ {
		newbuf = append(newbuf, 0)
	}
	return append(newbuf, buf...)
}
