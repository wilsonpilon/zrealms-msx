# RELEASE.md — Detalhamento do Release Oficial

**Versão Atual:** 0.1.0  
**Data:** 02 de Outubro de 2026  
**Status do Release:** Alpha Foundation (Fase 1: Subfases 1.1 e 1.2 Concluídas)  
**Alvo:** Windows (x64) para o Editor / MSX 2 & MSX-DOS 2 para a Engine  

---

## 1. Resumo Executivo do Release

Este release estabelece a fundação de dados, modelos de hardware do V9938 e repositórios de persistência do ecossistema **Z-Realm (`zrealm-msx`)**. A arquitetura foi validada para permitir o armazenamento robusto de projetos em arquivos únicos SQLite (`.rpgproj`), com garantia de integridade referencial por chaves estrangeiras ativas e layout alinhado aos limites do MSX 2.

---

## 2. O que foi Criado (Novos Componentes)

### 2.1. Núcleo de Banco de Dados (`pkg/project`)
* **`0001_initial_schema.sql`:** DDL contendo 10 tabelas relacionais normalizadas:
  * `schema_migrations`: Versionamento das tabelas.
  * `project_settings`: Configurações chave-valor do projeto.
  * `tilesets` e `tiles`: Padrões de 8x8 pixels para SCREEN 4 (Graphic 3).
  * `sprites`: Padrões de 16x16 pixels para Modo 2 (V9938).
  * `rooms`: Matriz matricial de 576 bytes (32 colunas x 18 linhas de viewport) e links cardeais.
  * `entities`: Atores, baús, portas e gatilhos com constraints estritas de coordenadas.
  * `scripts`, `string_table`, `hero_classes` e `items`.
* **`migrations.go`:** Sistema de auto-descoberta e aplicação transacional de migrações SQL embutidas via `embed.FS`.
* **`project.go`:** Gerenciador do ciclo de vida de projetos (`Create`, `Open`, `ValidateIntegrity`, `Close`, `GetSetting`, `SetSetting`, `Storage`).

### 2.2. Modelos de Domínio e Hardware MSX2 (`pkg/models`)
* **`constants.go`:** Constantes de hardware do V9938 (dimensões de tiles, sprites, viewport de 32x18, divisões de tela) e enums (`CollisionType`, `BehaviorType`, `ItemType`).
* **`models.go`:** Estruturas de dados fortemente tipadas com funções auxiliares de hardware:
  * Extração e injeção de pixels em padrões de 8x8 (`GetPixel`/`SetPixel`).
  * Atributos de cor de linha de varredura (nibble alto Fg / nibble baixo Bg).
  * Mapeamento de sprites 16x16 na ordem de memória nativa do V9938 (coluna esquerda nos primeiros 16 bytes, coluna direita nos 16 bytes subsequentes).
  * Acesso bidimensional à matriz de 576 bytes de salas (`GetTile`/`SetTile`).

### 2.3. Camada de Repositórios (`pkg/storage`)
* **`tileset_repo.go`:** Repositório de tilesets com suporte a upsert transacional de tiles.
* **`sprite_repo.go`:** Repositório de sprites 16x16 com validação de dimensões de buffers.
* **`room_repo.go`:** Repositório de salas (com busca por coordenadas `world_x, world_y`) e de entidades da sala.
* **`game_data_repo.go`:** Repositório de scripts, diálogos/textos, classes e catálogo de itens.
* **`storage.go`:** Agregador integrado diretamente na instância do projeto (`proj.Storage()`).

### 2.4. Utilitários e Governança
* **`cmd/zrealm/main.go`:** Utilitário CLI para verificação de versão, criação e validação de arquivos `.rpgproj`.
* **`pkg/version/version.go`:** Módulo dinâmico de controle de versão.
* **`build.ps1`:** Script PowerShell completo que auto-incrementa o número de compilação `Z`, executa os testes, compila os binários e gera o pacote `.zip` em `dist/`.
* **Documentação:** `LICENSE` (GPL 3), `README.md`, `OUTLINE.md`, `MANUAL.md`, `CHANGELOG.md` e `RELEASE.md`.

---

## 3. O que foi Corrigido / Refatorado

* **Resolução de Dependência CGO:** Adoção de `modernc.org/sqlite` em vez de `mattn/go-sqlite3`, permitindo compilações instantâneas no Windows sem dependência de MinGW ou bibliotecas C de terceiros no host de compilação.
* **Desacoplamento de Testes:** Eliminação de ciclo de importação detectado pelo compilador Go entre `pkg/storage` e `pkg/project` em arquivos de teste, isolando a fixture de teste de storage.
* **Integridade Referencial Ativa:** Configuração explícita do DSN do SQLite e validação via `PRAGMA foreign_keys = ON;`, prevenindo que tiles ou entidades órfãs sejam persistidos.

---

## 4. O que foi Testado & Resultados

Todos os testes foram executados e validados no ambiente Windows x64:

```text
=== RUN   TestTilePixelAndColors          --- PASS (0.00s)
=== RUN   TestSpriteModo2Layout           --- PASS (0.00s)
=== RUN   TestRoomMatrixAccess            --- PASS (0.00s)
=== RUN   TestCreateProject               --- PASS (0.15s)
=== RUN   TestCreateExistingFileFails     --- PASS (0.15s)
=== RUN   TestOpenExistingProject         --- PASS (1.18s)
=== RUN   TestForeignKeyEnforcement       --- PASS (0.21s)
=== RUN   TestCascadeDeleteTileset        --- PASS (0.16s)
=== RUN   TestSettingsCRUD                --- PASS (0.16s)
=== RUN   TestValidateIntegrity           --- PASS (0.16s)
=== RUN   TestProjectStorageIntegration   --- PASS (0.17s)
=== RUN   TestLoadMigrations              --- PASS (0.00s)
=== RUN   TestTilesetAndTileCRUD          --- PASS (0.70s)
=== RUN   TestSpriteCRUD                  --- PASS (0.16s)
=== RUN   TestRoomAndEntities             --- PASS (0.19s)
=== RUN   TestGameDataCRUD                --- PASS (0.16s)
=== RUN   TestVersion                     --- PASS (0.00s)
PASS: 100% de sucesso em todos os testes.
```

---

## 5. Próximos Passos para o Próximo Incremento (`0.1.x` -> `0.2.0`)

* **Subfase 1.3:** Implementação do módulo `pkg/exporter` (serialização binária dos segmentos de 16 KB para a Página 2 do MSX-DOS 2 e geração do manifesto `HEADER.BIN`).
* **Incremento Semântico:** Ao concluir a Subfase 1.3 (fechamento da Fase 1 de dados), o componente `Y` será incrementado.
