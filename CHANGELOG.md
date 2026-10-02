# CHANGELOG.md — Registro de Mudanças do Z-Realm

Todas as alterações notáveis neste projeto serão documentadas neste arquivo.

O formato baseia-se em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/) e este projeto adota as seguintes regras para versionamento **X.Y.Z**:
* **Z (Build / Patch):** Incrementado automaticamente a cada compilação executada pelo script `build.ps1`.
* **Y (Minor / Feature):** Incrementado a cada nova feature concluída e integrada ao projeto.
* **X (Major):** Incrementado a cada transição estrutural ou conclusão de uma grande fase (ex.: finalização da Camada de Dados, conclusão do Editor Gráfico, etc.).

---

## [0.1.0] - 2026-10-02

### Adicionado (Added)
- **Database Foundation (Subfase 1.1):**
  - Módulo Go inicializado como `github.com/zrealm-msx/zrealm`.
  - Driver SQLite 3 Pure-Go `modernc.org/sqlite` integrado sem dependência de toolchain CGO no Windows.
  - Esquema DDL inicial em `pkg/project/migrations/0001_initial_schema.sql` com 10 tabelas relacionais (`schema_migrations`, `project_settings`, `tilesets`, `tiles`, `sprites`, `scripts`, `rooms`, `entities`, `string_table`, `hero_classes`, `items`).
  - Carregador de migrações SQL embutidas com `embed.FS` e execução transacional ordenada.
  - Gerenciador `pkg/project` com suporte a `Create`, `Open`, `ValidateIntegrity` (verificação física com `PRAGMA quick_check` e referencial com `PRAGMA foreign_key_check`), e controle de settings.
- **Modelos de Domínio & Repositórios (Subfase 1.2):**
  - Definições e constantes de hardware do **MSX 2 (V9938)** em `pkg/models/constants.go`: SCREEN 4 (tiles 8x8), Sprites Modo 2 (16x16), áreas de tela e grid de 32x18 (576 bytes por sala).
  - Enums de colisão (`CollisionType`), comportamentos de IA (`BehaviorType`) e tipos de itens (`ItemType`).
  - Structs e métodos de hardware em `pkg/models/models.go`:
    - `Tile`: Manipulação de pixels e scanlines de cor (Fg/Bg).
    - `Sprite`: Ordem nativa de blocos verticais do Modo 2 do V9938 (coluna esquerda bytes 0..15, coluna direita bytes 16..31) e cores por scanline.
    - `Room`: Matriz matricial de 576 bytes com `GetTile` e `SetTile` validados.
  - Camada de persistência `pkg/storage`:
    - `TilesetRepository`: CRUD de tilesets e upsert atômico de tiles.
    - `SpriteRepository`: CRUD de sprites de 16x16 com validação de bytes.
    - `RoomRepository`: CRUD de salas (com busca por coordenadas do mundo) e entidades no grid 32x18.
    - `GameDataRepository`: CRUD para scripts, diálogos, classes de heróis e catálogo de itens.
    - Agregador `Storage` integrado a `Project.Storage()`.
- **Governança & Automação:**
  - Licença GNU General Public License v3.0 (`LICENSE`).
  - CLI administrativa inicial em `cmd/zrealm/main.go` para criação e checagem de projetos `.rpgproj`.
  - Controle de versão dinâmico em `VERSION` e `pkg/version/version.go`.
  - Script PowerShell `build.ps1` com auto-incremento de `Z`, execução de testes e geração de pacote `dist/*.zip`.
  - Documentação mestra: `README.md`, `SPEC.md`, `OUTLINE.md`, `MANUAL.md`, `CHANGELOG.md` e `RELEASE.md`.

### Testes (Tested)
- Testes unitários de integridade de banco de dados, migrações e chaves estrangeiras (`pkg/project/project_test.go`).
- Testes de validação dos cálculos de hardware V9938 em tiles e sprites (`pkg/models/models_test.go`).
- Testes completos de CRUD e constraints dos repositórios (`pkg/storage/storage_test.go`).
- Teste de integração ponta a ponta `TestProjectStorageIntegration`.
