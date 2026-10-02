package migrations

import (
	"testing"
)

func TestLoadMigrations(t *testing.T) {
	list, err := LoadMigrations()
	if err != nil {
		t.Fatalf("falha ao carregar migrações: %v", err)
	}

	if len(list) == 0 {
		t.Fatal("esperava pelo menos 1 migração embutida, obteve 0")
	}

	if list[0].Version != 1 {
		t.Errorf("primeira migração esperada com versão 1, obteve %d", list[0].Version)
	}

	if list[0].Name != "0001_initial_schema.sql" {
		t.Errorf("nome esperado '0001_initial_schema.sql', obteve '%s'", list[0].Name)
	}

	if len(list[0].SQL) == 0 {
		t.Error("conteúdo SQL da migração 0001 está vazio")
	}
}
