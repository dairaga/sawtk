package signing

import (
	"encoding/hex"
	"testing"
)

func TestCode(t *testing.T) {
	code := "123456"
	data := "1234567890"

	signer := GenerateSignerFromCode(code)
	signature := signer.Sign([]byte(data))

	if !Verify(signature, []byte(data), signer.GetPublicKey().AsBytes()) {
		t.Error("verify must be true but false")
	}

	signer2 := GenerateSignerFromCode("otherone")
	if Verify(signature, []byte(data), signer2.GetPublicKey().AsBytes()) {
		t.Error("verify with wrong private key must be false, but true")
	}

	signer3 := GenerateSignerFromCode(code)
	if !Verify(signature, []byte(data), signer3.GetPublicKey().AsBytes()) {
		t.Error("verify with the same code must be true but false")
	}
}

func TestGetSigner(t *testing.T) {
	singer, err := LoadSignerFromFile("test.priv")
	if err != nil {
		t.Fatal("get signer: ", err)
	}

	if singer.GetPublicKey().AsHex() != "03d73e65987f716a33fb2cf1bec01711c6bee200b90143ee656f1282fdd1276a9c" {
		t.Fatal("public key failure")
	}

	s := hex.EncodeToString(singer.Sign([]byte("1234567890")))

	if s != "24a5bfdd510b3fe5ce689353be0936bad5d07a1eb6cd35c50881c7215c55e87b2461778dc9fc66f0063b5b4d0cd18faa273f4744cc333a7df4ccc6e86af14add" {
		t.Fatal("signature failure")
	}
}
