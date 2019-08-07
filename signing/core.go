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

package signing

// -- Keys --

// PrivateKey A private key instance. The underlying content is dependent on
// implementation.
type PrivateKey interface {
	// Returns the algorithm name used for this private key.
	GetAlgorithmName() string

	// Returns the private key encoded as a hex string.
	AsHex() string

	// Returns the private key bytes.
	AsBytes() []byte
}

// PublicKey A public key instance. The underlying content is dependent on
// implementation.
type PublicKey interface {
	// Returns the algorithm name used for this public key.
	GetAlgorithmName() string

	// Returns the public key encoded as a hex string.
	AsHex() string

	// Returns the public key bytes.
	AsBytes() []byte
}

// -- Context --

// Context A context for a cryptographic signing algorithm.
type Context interface {
	// Returns the algorithm name used for this context.
	GetAlgorithmName() string

	// Generates a new random private key.
	NewRandomPrivateKey() PrivateKey

	// Produces a public key for the given private key.
	GetPublicKey(privateKey PrivateKey) PublicKey

	// Sign uses the given private key to calculate a signature for
	// the given data.
	Sign(message []byte, privateKey PrivateKey) []byte

	// Verify uses the given public key to verify that the given
	// signature was created from the given data using the associated
	// private key.
	Verify(signature []byte, message []byte, publicKey PublicKey) bool
}

// CreateContext Returns a Context instance by name.
func CreateContext(algorithmName string) Context {
	if algorithmName == "secp256k1" {
		return NewSecp256k1Context()
	}

	panic("No such algorithm")
}

// -- Signer --

// Signer A convenient wrapper of Context and PrivateKey.
type Signer struct {
	context    Context
	privateKey PrivateKey
}

// Sign Signs the given message.
func (s *Signer) Sign(message []byte) []byte {
	return s.context.Sign(message, s.privateKey)
}

// GetPublicKey returns the public key for this Signer instance.
func (s *Signer) GetPublicKey() PublicKey {
	return s.context.GetPublicKey(s.privateKey)
}

// -- CryptoFactory --

// CryptoFactory A factory for generating Signers.
type CryptoFactory struct {
	context Context
}

// NewCryptoFactory Creates a factory for generating Signers.
func NewCryptoFactory(context Context) *CryptoFactory {
	return &CryptoFactory{context: context}
}

// GetContext returns the context that backs this Factory instance.
func (factory *CryptoFactory) GetContext() Context {
	return factory.context
}

// NewSigner creates a new Signer for the given private key.
func (factory *CryptoFactory) NewSigner(privateKey PrivateKey) *Signer {
	return &Signer{factory.context, privateKey}
}
