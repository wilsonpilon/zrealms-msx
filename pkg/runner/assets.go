package runner

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed assets/*
var assetsFS embed.FS

// ExtractSystemFiles extrai os binários essenciais do runtime MSX-DOS 2 para o diretório de destino.
func ExtractSystemFiles(destDir string) error {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("falha ao criar pasta de destino %s: %w", destDir, err)
	}

	files := []string{"autoexec.bat", "COMMAND2.COM", "MSXDOS2.SYS", "zrealm.com"}
	for _, f := range files {
		data, err := assetsFS.ReadFile("assets/" + f)
		if err != nil {
			return fmt.Errorf("erro ao ler asset embutido %s: %w", f, err)
		}
		targetPath := filepath.Join(destDir, f)
		if err := os.WriteFile(targetPath, data, 0644); err != nil {
			return fmt.Errorf("erro ao extrair asset %s: %w", targetPath, err)
		}
	}
	return nil
}

// GetAssetBytes retorna os bytes de um asset do sistema pelo nome do arquivo.
func GetAssetBytes(name string) ([]byte, error) {
	return assetsFS.ReadFile("assets/" + name)
}
