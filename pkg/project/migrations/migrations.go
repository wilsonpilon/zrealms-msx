package migrations

import (
	"embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//go:embed *.sql
var migrationFiles embed.FS

// Migration representa um script de migração versionado.
type Migration struct {
	Version int
	Name    string
	SQL     string
}

// LoadMigrations lê todos os scripts SQL embutidos e os retorna ordenados por versão crescente.
func LoadMigrations() ([]Migration, error) {
	entries, err := migrationFiles.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("falha ao ler diretório de migrações: %w", err)
	}

	var migrations []Migration
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		parts := strings.SplitN(entry.Name(), "_", 2)
		if len(parts) < 2 {
			return nil, fmt.Errorf("nome de migração inválido (deve seguir NNNN_nome.sql): %s", entry.Name())
		}

		ver, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("versão inválida na migração %s: %w", entry.Name(), err)
		}

		content, err := migrationFiles.ReadFile(entry.Name())
		if err != nil {
			return nil, fmt.Errorf("falha ao ler arquivo de migração %s: %w", entry.Name(), err)
		}

		migrations = append(migrations, Migration{
			Version: ver,
			Name:    entry.Name(),
			SQL:     string(content),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}
