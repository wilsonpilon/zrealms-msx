# CHANGELOG.md — Registro de Mudanças do Z-Realm

Todas as alterações notáveis neste projeto serão documentadas neste arquivo.

O formato baseia-se em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/) e este projeto adota as seguintes regras para versionamento **X.Y.Z**:
* **Z (Build / Patch):** Incrementado automaticamente a cada compilação executada pelo script `build.ps1`.
* **Y (Minor / Feature):** Incrementado a cada nova feature concluída e integrada ao projeto.
* **X (Major):** Incrementado a cada transição estrutural ou conclusão de uma grande fase (ex.: finalização da Camada de Dados, conclusão do Editor Gráfico, etc.).

## [0.4.0] - 2026-10-02

### Adicionado (Added)
- **Fase 3 Completa: Editor Desktop Visual em Fyne (GUI 100% Funcional):**
  - **Subfase 3.1 — Shell da Aplicação & Navegação:**
    - Tema retro escuro `RetroDarkTheme` inspirado nas cores e fontes do V9938.
    - Gerenciador de estado reativo e concorrente `ProjectState` com detecção de alterações e dirty tracking.
    - Menu mestre da aplicação com verificação de integridade física (`quick_check`) e referencial (`foreign_key_check`) do SQLite.
    - Barra de status ao vivo com indicador de recursos e estimativa de consumo do Memory Mapper MSX.
  - **Subfase 3.2 — Editor de Tiles (8x8 SCREEN 4):**
    - Grid pixel-a-pixel interativo com paleta oficial V9938 de 16 cores.
    - Seletores de cor individuais de Foreground e Background por scanline.
    - Configuração de tipo de colisão física (Passável, Sólido, Água, Dano, Gatilho).
    - Transformações completas de matriz: Rotação 90°, Flip H/V, Deslocamento (Shift direcional), Inversão, Limpeza e Preenchimento.
  - **Subfase 3.3 — Editor de Sprites (16x16 Modo 2):**
    - Grid de 16x16 pixels com divisores de quadrantes 8x8 e transparência sutil.
    - Atribuição independente de cor para cada uma das 16 scanlines do Modo 2 do V9938.
    - Pré-visualização com pixels nítidos em escala 1x (16x16) e 4x (64x64).
    - Recurso "Duplicar Quadro" para prototipagem rápida de ciclos de animação (Walk/Idle).
  - **Subfase 3.4 — Editor de Salas (Room Matrix View):**
    - Viewport interativo de 32x18 tiles (256x144 pixels) correspondente à área jogável SCREEN 4.
    - Ferramentas: Pincel/Carimbo contínuo, Balde de Tinta (`FloodFillRoom`), Borracha e Conta-Gotas (`ToolEyedropper`).
    - Utilitários: Limpeza total, Preenchimento total e Preenchimento automático de bordas (`FillBorderRoom`).
    - Paleta de carimbo com preview de tiles e seletor numérico direto (0..255).
    - Conexões cardeais (Norte, Sul, Leste, Oeste) com navegação direta e "Auto-Conectar por Coordenadas" (`AutoConnectRooms`).
    - Gestor visual de entidades com inserção de NPCs, baús, portas e marcadores gráficos sobrepostos na sala.
  - **Subfase 3.5 — Editor de Regras, Tabelas de RPG e Roteiros:**
    - Ficha de classes de herói com parâmetros de combate (HP, MP, Ataque, Defesa).
    - Catálogo de itens e equipamentos com categorias, preços e modificadores de atributos.
    - Simulador de caixa de diálogo com proporção nativa MSX (32x4 caracteres nas linhas 20-23).
    - Compilador de scripts de eventos (`CompileScript`) para a Bytecode VM com opcodes compactos (`OP_MSG`, `OP_GIVE_ITEM`, `OP_TAKE_ITEM`, `OP_SET_FLAG`, `OP_CHECK_FLAG`, `OP_TELEPORT`, `OP_HEAL`, `OP_DAMAGE`, `OP_PLAY_SFX`, `OP_END`).
    - Desassemblador de bytecode (`DisassembleScript`) e exibição de hexadecimal formatado.

---

## [0.3.0] - 2026-10-02

### Adicionado (Added)
- **Engine MSX 2 & Prototipagem de Baixo Nível (Fase 2):**
  - Bootstrap completo no **MSX-DOS 2** (`zrealm.com`) utilizando compilador SDCC 4.6.0 e framework MSXgl.
  - Gerenciador de **Memory Mapper** (`mapper.c`) em conformidade com as regras do MSX-DOS 2 (`EXTBIOS`):
    - Alocação dinâmica de segmentos de usuário via `ALL_SEG`.
    - Paginação segura na **Página 2 (`0x8000 - 0xBFFF`)** via `PUT_P2`, mantendo a Engine fixa na Página 1 (`0x4000 - 0x7FFF`) e o sistema operacional/pilha na Página 3 (`0xC000 - 0xFFFF`).
    - Liberação garantida de 100% dos segmentos alocados via `FRE_SEG` ao finalizar, sem vazamento de RAM.
  - Driver de vídeo para **V9938 SCREEN 4 (Graphic 3)** (`vdp_screen4.c`):
    - Configuração de 3 bancos verticais com 2048 bytes de padrões e 2048 bytes de atributos de cores.
    - Viewport de exploração de **32 x 18 tiles (256x144 pixels / 576 bytes contíguos)**.
    - Área de HUD (linhas 18-19) e Diálogo (linhas 20-23) limpas com tile preto/vazio (Tile 255).
    - Desativação limpa de sprites do Modo 2 por coordenadas de scanline fora da tela.
    - Restauração perfeita do modo de texto BIOS (Screen 0) e gerador de caracteres via chamada de interslot `INITXT` (0x006C).
  - **Binary Disk Loader para MSX-DOS 2** (`loader.c`):
    - Carregador de disco de alta velocidade usando descritores de arquivo nativos do DOS 2 (`DOS_OpenHandle`, `DOS_ReadHandle`, `DOS_CloseHandle`).
    - Leitura e validação do cabeçalho mestre `HEADER.BIN` (Magic `ZR01`, versão, alocação de segmentos).
    - Streaming direto de `GAME.DAT` em blocos de 8 KB para a janela da Página 2 do Memory Mapper.
    - Funções de acesso instantâneo a recursos com paginação transparente de segmentos (`LOADER_GetTileset`, `LOADER_GetRoom`).
  - **Comando `-demo` no Toolkit Go (`cmd/zrealm`):**
    - Geração de projeto demonstrativo (`demo.rpgproj`) com tileset de masmorra, 2 salas interligadas (Entrada com portal Leste e Câmara dos Pilares com cruz central), sprite do herói e diálogos.
    - Integração de exportação direta do SQLite para o disco de boot MSX (`HEADER.BIN` e `GAME.DAT`).
  - **Empacotamento de Disco & Testes Automatizados no openMSX:**
    - Atualização do `project_config.js` para inclusão automática de `HEADER.BIN` e `GAME.DAT` na imagem de disco `DOS2_zrealm.dsk` via `msxtar`.
    - Script de teste automatizado `boot_test.tcl` executado com `-machine Philips_NMS_8250 -ext msxdos2 -ext ram512k`, validando visualmente a Sala 1, a Sala 2 (por troca de segmento) e o retorno com texto limpo ao prompt `A:\>`.

---

## [0.2.0] - 2026-10-02

### Adicionado (Added)
- **Engine de Serialização Binária para MSX 2 (Subfase 1.3):**
  - Módulo `pkg/exporter` responsável por compilar dados do SQLite em estruturas nativas Z80/VDP.
  - Empacotador de memória alinhado a blocos de 16 KB (Página 2 / `8000h - BFFFh`) com detecção automática de overflow de segmento.
  - Tabela mestra de cabeçalho `HEADER.BIN` (Magic `ZR01`, versão, contagem de segmentos, dados da sala inicial e diretório de 8 bytes por recurso).
  - Serializador de salas em blocos compactos de 656 bytes (cabeçalho de 16 bytes com conexões cardeais + matriz de 576 bytes + 8 entidades com posições e scripts).
  - Serializador de tilesets em formato nativo V9938 SCREEN 4 (4.608 bytes: pattern table, color table, tabela de física/colisão e tabela de animações cíclicas).
  - Serializador de sprites de 16x16 Modo 2 (48 bytes: 32 bytes pattern + 16 bytes atributos de cor por scanline).
  - Serializador de strings e diálogos com índice relativo de offsets de 16 bits para acesso aleatório sem latência.
  - Geração automática de `GAME.DAT` contíguo e arquivos modulares de banco `SEGxx.BNK`.
  - Integração da flag `-export` no utilitário de linha de comando `zrealm.exe`.
- **Conclusão da Fase 1:** Toda a fundação de dados, modelos e exportador binário para MSX 2 finalizada e validada.

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
