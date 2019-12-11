package wallet

import "testing"

func TestWallet(t *testing.T) {
	pubkey := `03511c83916ac338835b07f6b9f7c0aa10b7b427b48e16b5e91360c919c9cf60cb`
	w, err := MakeFromHex(pubkey)

	if err != nil {
		t.Error(err)
	}

	ans := `1Lpgbz8o24ENRsZD3Rr5fVzJr2Ln4BBi5F`
	if w != ans {
		t.Errorf("wallet want %q, but %q", ans, w)
	}
}
