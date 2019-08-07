package ns

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"
)

// _sha512 returns a SHA-512 hex string.
func _sha512(data []byte) string {
	hash := sha512.New()
	hash.Write(data)
	bytes := hash.Sum(nil)
	return hex.EncodeToString(bytes)
}

// Sha256 returns a SHA-256 hex string.
func _sha256(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	bytes := hash.Sum(nil)
	return hex.EncodeToString(bytes)
}

var (
	emptyHash = _sha256([]byte(""))[:16]
	//addrRegx  = regexp.MustCompile(`[0-9a-f]{70}`)
)

// Namespace sawtooth namespace
type Namespace interface {
	fmt.Stringer
	MakeAddress(string) string
	//Validate(string) bool
}

// GeneralNS prefix 64 character with sha512
type GeneralNS struct {
	prefix string
}

// MakeAddress ...
func (ns *GeneralNS) MakeAddress(addr string) string {
	return ns.prefix + _sha512([]byte(addr))[:64]
}

/*
// Validate ...
func (ns *GeneralNS) Validate(addr string) bool {
	return IsAddress(addr) && strings.HasPrefix(addr, ns.prefix)
}
*/
func (ns *GeneralNS) String() string {
	return ns.prefix
}

// SawtoothNS is implements namespace rules of sawtooth.
type SawtoothNS struct {
	GeneralNS
}

// MakeAddress ...
func (sns *SawtoothNS) MakeAddress(addr string) string {
	tmp := strings.SplitN(addr, ".", 4)
	b := strings.Builder{}

	for _, x := range tmp {
		b.WriteString(_sha256([]byte(x))[:16])
	}

	if len(tmp) < 4 {
		b.WriteString(strings.Repeat(emptyHash, 4-len(tmp)))
	}

	return sns.prefix + b.String()
}

var (
	settings = &SawtoothNS{GeneralNS{"000000"}} // build-in settings family of sawtooth.
)

// New return sawtooth namespace
func New(name string) Namespace {
	switch name {
	case "000000":
		return settings
	default:
		return &GeneralNS{prefix: _sha512([]byte(name))[:6]}
	}
}

// NewSawtoothNS ...
func NewSawtoothNS(name string) Namespace {
	return &SawtoothNS{GeneralNS{prefix: _sha512([]byte(name))[:6]}}
}

// Settings return settings-tp namespace
func Settings() Namespace {
	return settings
}

// EmptyHash returns empty hash code
func EmptyHash() string {
	return emptyHash
}

/*
// IsAddress returns address length is 70 or not.
func IsAddress(addr string) bool {
	if len(addr) != 70 {
		return false
	}

	return sawtk.IsHexString(addr)
}
*/
