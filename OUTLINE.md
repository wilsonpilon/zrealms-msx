# OUTLINE.md — Acompanhamento do Projeto Z-Realm (`zrealm-msx`)

Este documento serve como mapa de bordo vivo do projeto, refletindo o roadmap estabelecido em [SPEC.md](SPEC.md), as fases concluídas, o trabalho em andamento e os próximos marcos a serem executados.

---

## Roadmap Global & Status Atual

| Fase | Descrição | Status | Marco Principal |
| :--- | :--- | :---: | :--- |
| **Fase 1** | **Estruturação de Dados & Abstrações do Editor (Go + SQLite)** | 🟢 **100% Concluída** | Banco `.rpgproj`, Modelos de Hardware, Repositórios & Exporter |
| **Fase 2** | **Prototipagem de Baixo Nível no MSX (C + MSXgl + SDCC)** | 🟢 **100% Concluída** | MSX-DOS 2, Paging Mapper (Page 2), SCREEN 4 e Binary Loader de Disco |
| **Fase 3** | **Desenvolvimento da GUI Desktop com Fyne (O Editor Visual)** | ⚪ *Planejado* | Editores de Tilesets, Sprites, Salas 32x18, Roteiros e Regras |
| **Fase 4** | **Gameplay Engine & Máquina de Eventos (MSX)** | ⚪ *Planejado* | Controle no Grid, VM de Eventos, HUD e Caixas de Diálogo |
| **Fase 5** | **Pipeline Integrado de Build & Jogo de Referência** | ⚪ *Planejado* | "One-Click Run", Teste em openMSX e Jogo Demonstrador |

---

## Detalhamento das Subfases

### Fase 1: Estruturação de Dados & Abstrações do Editor (Go + SQLite)
- [x] **Subfase 1.1 — Schema & Database Foundation em Go:**
  - [x] Inicialização do módulo Go (`github.com/zrealm-msx/zrealm`) e driver SQLite Pure-Go (`modernc.org/sqlite`).
  - [x] Schema DDL relacional versionado com suporte a migrações embutidas (`embed.FS`).
  - [x] Gerenciador `pkg/project` com criação, abertura, validação de integridade física e Foreign Keys.
  - [x] Testes unitários cobrindo transações, WAL e integridade referencial.
- [x] **Subfase 1.2 — Repositórios & Modelos de Domínio em Go:**
  - [x] Constantes de hardware do V9938 (tiles 8x8, sprites 16x16, matriz 32x18, áreas de tela).
  - [x] Structs de domínio (`Tileset`, `Tile`, `Sprite`, `Room`, `Entity`, `Script`, `HeroClass`, `Item`).
  - [x] Métodos de cálculo de bits e scanlines de hardware V9938 (ordenação de blocos verticais do Modo 2).
  - [x] Camada de repositórios `pkg/storage` (CRUD de todas as entidades com constraints e limites).
  - [x] Testes unitários e de integração completos com 100% de aprovação.
- [x] **Subfase 1.3 — Engine de Serialização Binária (The Exporter):**
  - [x] Empacotador de salas em blocos lógicos alinhados a 16 KB (Página 2 / 8000h-BFFFh).
  - [x] Codificador de tabelas de padrões e cores para SCREEN 4 e Sprites Modo 2.
  - [x] Gerador da tabela mestra de cabeçalho (`HEADER.BIN` / Manifest de recursos) com Magic `ZR01`.
  - [x] Geração contígua de `GAME.DAT` e modular de `SEGxx.BNK`.
  - [x] Testes automatizados do exportador validando alinhamento e integridade binária.

---

### Fase 2: Prototipagem de Baixo Nível no MSX (C + MSXgl + SDCC)
- [x] **Subfase 2.1 — Bootstrap MSX-DOS 2 & Abstração de Mapper:**
  - [x] Esqueleto de aplicação `.COM` com inicialização MSXgl sob MSX-DOS 2.
  - [x] Módulo `mapper.c`: interrogação via `EXTBIOS` (Jump Table do DOS 2), alocação de segmentos (`ALL_SEG`) e liberação garantida (`FRE_SEG`).
  - [x] Rotina de paginação segura na Página 2 (`PUT_P2` / `MAPPER_SetPage2`), mantendo Página 1 fixa (Engine) e Página 3 (DOS 2 / Stack).
  - [x] Restauração perfeita do modo de texto BIOS (Screen 0) via `INITXT` (0x006C) ao sair, com zero vazamento de memória.
- [x] **Subfase 2.2 — Configuração do VDP V9938 (SCREEN 4):**
  - [x] Ativação de Graphic 3 com 3 bancos verticais de padrões (2048 bytes cada) e cores (2048 bytes cada).
  - [x] Viewport de jogo em 32 x 18 tiles (576 bytes contíguos copiados diretamente para VRAM).
  - [x] Área reservada e limpa para HUD (linhas 18-19) e Diálogo (linhas 20-23).
  - [x] Desativação e ocultamento de sprites via coordenadas Y fora da tela (Modo 2).
- [x] **Subfase 2.3 — Binary Disk Loader & Integração com Go Exporter:**
  - [x] Módulo `loader.c` para leitura direta de `HEADER.BIN` e streaming de `GAME.DAT` em blocos de 8 KB para o Memory Mapper via handles do DOS 2 (`DOS_OpenHandle`, `DOS_ReadHandle`, `DOS_CloseHandle`).
  - [x] Lookup O(1) de recursos (`LOADER_GetTileset`, `LOADER_GetRoom`) que chaveiam a Página 2 automaticamente para o segmento do recurso.
  - [x] Comando `-demo` adicionado ao CLI Go `zrealm` gerando projeto demonstrativo com 2 salas conectadas e tileset customizado.
  - [x] Teste ponta a ponta verificado no openMSX (`Philips_NMS_8250` + `ram512k` + `msxdos2`), capturando screenshots automáticos da Sala 1 (Entrada), Sala 2 (Câmara dos Pilares) e saída limpa de volta ao prompt `A:\>`.

---

### Fase 3: Desenvolvimento da GUI Desktop com Fyne (O Editor Visual)
- [ ] **Subfase 3.1 — Shell da Aplicação & Navegação:**
  - [ ] Janela principal com tema retrô escuro, gerenciamento de projetos e abas de ferramentas.
- [ ] **Subfase 3.2 — Editor de Tiles (8x8):**
  - [ ] Grid de desenho pixel-a-pixel com paleta MSX2 e colisão física.
- [ ] **Subfase 3.3 — Editor de Sprites (16x16):**
  - [ ] Edição em Modo 2 com atribuição de cores por scanline e quadros de animação.
- [ ] **Subfase 3.4 — Editor de Salas (Room Matrix View):**
  - [ ] Matriz de pintura 32x18 com carimbo de tiles, conexões cardeais e inserção visual de entidades.
- [ ] **Subfase 3.5 — Editor de Regras, Tabelas de RPG e Roteiros:**
  - [ ] Formulários para classes, itens e editor textual de scripts.

---

### Fase 4: Gameplay Engine & Máquina de Eventos (MSX)
- [ ] **Subfase 4.1 — Movimentação do Herói & Colisão no Grid:**
  - [ ] Entrada de teclado/joystick com movimentação em passos de 8 pixels e troca automática de sala nas bordas.
- [ ] **Subfase 4.2 — Sistema de Entidades e Atores:**
  - [ ] Até 8 entidades ativas na sala, IAs simples e interação por tecla de ação.
- [ ] **Subfase 4.3 — Máquina Virtual de Eventos (Bytecode VM):**
  - [ ] Interpretador de instruções compactas (`OP_MSG`, `OP_GIVE_ITEM`, `OP_CHECK_FLAG`, etc.).
- [ ] **Subfase 4.4 — Caixa de Diálogos & HUD:**
  - [ ] Renderizador de texto nas linhas 20-23 e mostrador de HP/MP nas linhas 18-19.

---

### Fase 5: Pipeline Integrado de Build & Jogo de Referência
- [ ] **Subfase 5.1 — Automação "One-Click Run":**
  - [ ] Disparo automático de exportação, compilação e execução no openMSX.
- [ ] **Subfase 5.2 — Backend MegaROM (Cartucho):**
  - [ ] Geração de arquivo `.ROM` unificado com chaveador de banco em cartucho.
- [ ] **Subfase 5.3 — Jogo de Referência:**
  - [ ] Mini-cRPG demonstrador com masmorra de 20 salas, enigmas, NPCs e combate simples.
