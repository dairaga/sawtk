package tp

import (
	"fmt"
	"strings"
)

// Family is for TransactionProcessor.
type Family struct {
	name       string
	versions   []string
	namespaces []string
}

func (f *Family) String() string {
	v := strings.Join(f.versions, `","`)
	n := strings.Join(f.namespaces, `","`)

	return fmt.Sprintf(`{"name": "%s", "versions": ["%s"], "namespaces": ["%s"]}`, f.name, v, n)
}

// FamilyName returns family name defined in processor.TransactionHandler.
func (f *Family) FamilyName() string {
	return f.name
}

// FamilyVersions returns family support versions defined in processor.TransactionHandler.
func (f *Family) FamilyVersions() []string {
	return f.versions
}

// Namespaces returns namespaces of family defined in processor.TransactionHandler.
func (f *Family) Namespaces() []string {
	return f.namespaces
}

// NewFamily returns a Sawtooth TP family.
func NewFamily(name string, versions, namespaces []string) *Family {
	if len(versions) <= 0 || len(namespaces) <= 0 {
		panic("a sawtooth family must have support versions and namespaces")
	}

	return &Family{
		name:       name,
		versions:   versions,
		namespaces: namespaces,
	}
}
