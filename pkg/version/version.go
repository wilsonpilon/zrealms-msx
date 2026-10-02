package version

import (
	_ "embed"
	"strings"
)

//go:embed VERSION
var rawVersion string

func clean(s string) string {
	s = strings.TrimPrefix(s, "\ufeff")
	return strings.TrimSpace(s)
}

// Version retorna a versão atual do sistema Z-Realm no formato X.Y.Z.
var Version = clean(rawVersion)

// String retorna a versão formatada com prefixo 'v'.
func String() string {
	v := clean(Version)
	if v == "" {
		return "v0.1.0-dev"
	}
	if !strings.HasPrefix(v, "v") {
		return "v" + v
	}
	return v
}
