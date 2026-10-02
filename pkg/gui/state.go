package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/zrealm-msx/zrealm/pkg/exporter"
	"github.com/zrealm-msx/zrealm/pkg/project"
)

// ProjectState gerencia o estado reativo global do projeto carregado na interface Fyne.
type ProjectState struct {
	mu          sync.RWMutex
	project     *project.Project
	filePath    string
	isModified  bool
	recentFiles []string

	onProjectLoaded   []func(proj *project.Project, path string)
	onProjectClosed   []func()
	onModifiedChanged []func(modified bool)
	onDataChanged     []func()
}

// NewProjectState cria uma nova instância de gerenciamento de estado da GUI.
func NewProjectState() *ProjectState {
	return &ProjectState{
		recentFiles: make([]string, 0),
	}
}

// Current retorna o ponteiro para o projeto atualmente aberto, se houver.
func (s *ProjectState) Current() *project.Project {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.project
}

// FilePath retorna o caminho do arquivo do projeto ativo no disco.
func (s *ProjectState) FilePath() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.filePath
}

// IsOpen informa se existe um projeto ativo aberto na interface.
func (s *ProjectState) IsOpen() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.project != nil
}

// IsModified informa se há alterações pendentes de salvamento.
func (s *ProjectState) IsModified() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isModified
}

// SetModified atualiza a flag de modificado e notifica ouvintes.
func (s *ProjectState) SetModified(mod bool) {
	s.mu.Lock()
	if s.isModified == mod {
		s.mu.Unlock()
		return
	}
	s.isModified = mod
	listeners := append([]func(bool){}, s.onModifiedChanged...)
	s.mu.Unlock()

	for _, l := range listeners {
		l(mod)
	}
}

// NotifyDataChanged dispara notificações de que os dados do projeto sofreram alterações.
func (s *ProjectState) NotifyDataChanged() {
	s.SetModified(true)
	s.mu.RLock()
	listeners := append([]func(){}, s.onDataChanged...)
	s.mu.RUnlock()

	for _, l := range listeners {
		l()
	}
}

// OpenProject abre um arquivo .rpgproj existente.
func (s *ProjectState) OpenProject(path string) error {
	s.mu.Lock()
	if s.project != nil {
		_ = s.project.Close()
		s.project = nil
	}

	proj, err := project.Open(path)
	if err != nil {
		s.mu.Unlock()
		return fmt.Errorf("erro ao abrir projeto: %w", err)
	}

	s.project = proj
	s.filePath = path
	s.isModified = false
	s.addRecentLocked(path)

	loadedListeners := append([]func(*project.Project, string){}, s.onProjectLoaded...)
	s.mu.Unlock()

	for _, l := range loadedListeners {
		l(proj, path)
	}
	return nil
}

// NewProject cria um novo projeto .rpgproj no caminho especificado.
func (s *ProjectState) NewProject(path string, name string) error {
	s.mu.Lock()
	if s.project != nil {
		_ = s.project.Close()
		s.project = nil
	}

	proj, err := project.Create(path, name)
	if err != nil {
		s.mu.Unlock()
		return fmt.Errorf("erro ao criar projeto: %w", err)
	}

	s.project = proj
	s.filePath = path
	s.isModified = false
	s.addRecentLocked(path)

	loadedListeners := append([]func(*project.Project, string){}, s.onProjectLoaded...)
	s.mu.Unlock()

	for _, l := range loadedListeners {
		l(proj, path)
	}
	return nil
}

// OpenDemoProject gera um projeto de demonstração temporário ou no caminho especificado e o abre.
func (s *ProjectState) OpenDemoProject(targetPath string) error {
	if targetPath == "" {
		tmpDir := os.TempDir()
		targetPath = filepath.Join(tmpDir, "demo_zrealm.rpgproj")
	}

	proj, err := project.CreateDemoProject(targetPath)
	if err != nil {
		return fmt.Errorf("erro ao criar demo: %w", err)
	}
	_ = proj.Close()

	return s.OpenProject(targetPath)
}

// CloseProject fecha o projeto atualmente ativo.
func (s *ProjectState) CloseProject() {
	s.mu.Lock()
	if s.project != nil {
		_ = s.project.Close()
		s.project = nil
	}
	s.filePath = ""
	s.isModified = false

	closedListeners := append([]func(){}, s.onProjectClosed...)
	s.mu.Unlock()

	for _, l := range closedListeners {
		l()
	}
}

// Export executa o exportador binário para MSX 2.
func (s *ProjectState) Export(outDir string) (*exporter.ExportResult, error) {
	s.mu.RLock()
	proj := s.project
	filePath := s.filePath
	s.mu.RUnlock()

	if proj == nil {
		return nil, fmt.Errorf("nenhum projeto aberto para exportar")
	}

	if outDir == "" {
		outDir = filepath.Join(filepath.Dir(filePath), "build_msx")
	}

	return exporter.Export(proj, exporter.ExportOptions{
		OutputDir:       outDir,
		ExportBankFiles: true,
	})
}

// AddRecentFile adiciona um caminho à lista de recentes.
func (s *ProjectState) AddRecentFile(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.addRecentLocked(path)
}

func (s *ProjectState) addRecentLocked(path string) {
	filtered := make([]string, 0, len(s.recentFiles)+1)
	filtered = append(filtered, path)
	for _, f := range s.recentFiles {
		if f != path {
			filtered = append(filtered, f)
		}
		if len(filtered) >= 10 {
			break
		}
	}
	s.recentFiles = filtered
}

// RecentFiles retorna os arquivos recentes.
func (s *ProjectState) RecentFiles() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]string, len(s.recentFiles))
	copy(res, s.recentFiles)
	return res
}

// OnProjectLoaded registra um callback invocado ao carregar um projeto.
func (s *ProjectState) OnProjectLoaded(fn func(proj *project.Project, path string)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onProjectLoaded = append(s.onProjectLoaded, fn)
}

// OnProjectClosed registra um callback invocado ao fechar um projeto.
func (s *ProjectState) OnProjectClosed(fn func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onProjectClosed = append(s.onProjectClosed, fn)
}

// OnModifiedChanged registra um callback invocado quando o status de modificado muda.
func (s *ProjectState) OnModifiedChanged(fn func(modified bool)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onModifiedChanged = append(s.onModifiedChanged, fn)
}

// OnDataChanged registra um callback invocado quando dados são modificados.
func (s *ProjectState) OnDataChanged(fn func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onDataChanged = append(s.onDataChanged, fn)
}
