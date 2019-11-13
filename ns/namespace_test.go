package ns_test

import (
	"testing"

	"github.com/dairaga/sawtk/ns"
)

func TestNamespace(t *testing.T) {
	ns1 := ns.New("intkey")

	if ns1.Prefix() != "1cf126" {
		t.Fatal("namespace not match:", ns1.String(), "1cf126")
	}

	ns1 = ns.New("000000")
	tmp := ns1.MakeAddress("sawtooth.config.vote.proposals")
	if tmp != "000000a87cb5eafdcca6a8b79606fb3afea5bdab274474a6aa82c1c0cbf0fbcaf64c0b" {
		t.Fatal("settings namespace fail", tmp, "000000a87cb5eafdcca6a8b79606fb3afea5bdab274474a6aa82c1c0cbf0fbcaf64c0b")
	}

	tmp = ns1.MakeAddress("mykey")

	if tmp != "0000005e50f405ace6cbdfe3b0c44298fc1c14e3b0c44298fc1c14e3b0c44298fc1c14" {
		t.Fatal("settings namespace fail", tmp, "0000005e50f405ace6cbdfe3b0c44298fc1c14e3b0c44298fc1c14e3b0c44298fc1c14")
	}

	tmp = ns1.MakeAddress("diviner.exchange")
	if tmp != "0000008923f4638a4a5030ab27b729d9cc4cb1e3b0c44298fc1c14e3b0c44298fc1c14" {
		t.Fatal("settings namespace fail", tmp, "0000008923f4638a4a5030ab27b729d9cc4cb1e3b0c44298fc1c14e3b0c44298fc1c14")
	}

	//fmt.Println(sawtk.GetEmptyHash())
}

/*
func TestValidateAddress(t *testing.T) {
	myns := ns.New("mytest")

	addr := myns.MakeAddress("testtest")

	if !ns.IsAddress(addr) {
		t.Errorf("%s is address", addr)
	}

	addr = strings.Repeat("0", 70)
	if !ns.IsAddress(addr) {
		t.Errorf("%s is address", addr)
	}

	addr = strings.Repeat("f", 70)
	if !ns.IsAddress(addr) {
		t.Errorf("%s is address", addr)
	}

	addr = strings.Repeat("g", 70)
	if ns.IsAddress(addr) {
		t.Errorf("%s is not address", addr)
	}

	addr = strings.Repeat("a", 69)
	if ns.IsAddress(addr) {
		t.Errorf("%s is not address", addr)
	}
}
*/
