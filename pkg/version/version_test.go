package version

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	v := String()
	if !strings.HasPrefix(v, "v0.") {
		t.Errorf("versão esperada começando com 'v0.', obtido: %s", v)
	}
}
