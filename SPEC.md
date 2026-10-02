# SPEC.md — Arquitetura e Especificação Técnica
## Z-Realm (`zrealm-msx`) — O ZZT dos cRPGs para MSX
### (Editor Desktop em Go/Fyne/SQLite + Engine MSX2/MSX-DOS 2 em C/MSXgl)

**Projeto Oficial:** Z-Realm (`zrealm-msx`)  
**Versão:** 0.3.2  
**Data:** Outubro de 2026  
**Status:** Fases 1 e 2 Concluídas — Rumo à Fase 3 (Editor Desktop Go + Fyne)  
**Autor:** Equipe de Arquitetura de Software & Retrocomputação  

---

## 1. Visão Geral e Escopo

### 1.1. O Conceito: "O ZZT dos cRPGs para MSX"
O objetivo deste projeto é conceber e implementar um ecossistema completo de desenvolvimento de jogos de interpretação (cRPGs) e aventuras interativas para a plataforma MSX.

Inspirado na facilidade criativa e modularidade de **ZZT** (Tim Sweeney, 1991) e na profundidade dos clássicos de exploração por salas e matrizes de tiles (como *Ultima*, *Dragon Quest*, *Rogue* e *The Magic Candle*), o sistema divide-se em duas metades perfeitamente acopladas:

1. **O Editor (PC/Desktop):** Uma ferramenta de criação visual, moderna, ergonômica e multiplataforma desenvolvida em **Go + Fyne**, que armazena a totalidade do projeto (tiles, sprites, mapas, roteiros, regras e tabelas de atributos) em um arquivo único **SQLite**.
2. **A Engine e o Compilador (MSX):** Um runtime ultracompacto e de alta performance desenvolvido em **C puro (SDCC + biblioteca MSXgl)**, voltado para computadores **MSX 2 e MSX 2+** sob o ambiente de disco **MSX-DOS 2**.

### 1.2. Filosofia de Design: Compilação vs. Interpretação Lenta
Ao contrário de engines genéricas com interpretadores pesados que sofrem no Z80 a 3.58 MHz, a nossa arquitetura opera por **compilação e empacotamento estático**:
- O Editor no PC valida, otimiza e empacota o banco SQLite em estruturas binárias nativas alinhadas a blocos de memória Z80.
- A Engine embutida executa um loop de processamento determinístico orientado a eventos discretos (Grid-based / Room-based), minimizando overhead de CPU.
- O resultado é um executável nativo do MSX-DOS 2 (`.COM`) que roda com taxa de atualização estável (50/60 Hz VBLANK), com transições instantâneas entre salas.

---

## 2. Pilha Tecnológica e Arquitetura do Sistema

```
+-----------------------------------------------------------------------+
|                         PC / DESKTOP (Editor)                        |
|                                                                       |
|  [ Fyne GUI (Go) ]  <--->  [ Project Repository (Go) ]                |
|         |                            |                                |
|         v                            v                                |
|  [ Visual Editors ]          [ SQLite Database ]                      |
|  (Tiles, Maps, Scripts)      (Entities, Dialogues, Flags, Maps)       |
|                                      |                                |
|                                      v                                |
|                           [ Asset Packer / Exporter ]                 |
+--------------------------------------|--------------------------------+
                                       | Chunks binários / .BNK
                                       v
+-----------------------------------------------------------------------+
|                    MSX TOOLCHAIN (SDCC + MSXgl)                       |
|                                                                       |
|  [ Engine Core (C) ] + [ Generated Assets ] ---> [ SDCC Linker ]      |
+--------------------------------------|--------------------------------+
                                       | Output: GAME.COM + GAME.DAT
                                       v
+-----------------------------------------------------------------------+
|                     MSX 2 / MSX 2+ (Target Runtime)                   |
|                                                                       |
|  [ MSX-DOS 2 Kernel ]                                                 |
|  [ V9938 VDP ] --------> SCREEN 4 (256x192, 3 Banks, Sprite Mode 2)   |
|  [ Memory Mapper ] ----> Paging Window (Page 2: 8000h-BFFFh)          |
|                          (>= 256KB até 2MB+ de Assets em RAM)         |
+-----------------------------------------------------------------------+
```

### 2.1. Ambiente do Editor (PC)
- **Linguagem:** Go (Go 1.22+)
- **Interface Gráfica (GUI):** Fyne v2 (renderização via OpenGL acelerada por GPU, layouts responsivos, suporte nativo a temas escuros e estética retro-moderna).
- **Armazenamento de Projeto:** **SQLite 3** (arquivo `.rpgproj` ou `.db`). Todas as tabelas são normalizadas, permitindo auditoria direta por SQL, histórico de revisões e consistência transacional ACID.
- **Pipeline de Exportação:** Módulos em Go dedicados à conversão de dados do SQLite para formatos binários nativos do Z80/VDP (ex.: planos de bits de caracteres, tabelas de cores, mapas de índices compactados e bytecode de eventos).

### 2.2. Ambiente da Engine (MSX)
- **Linguagem:** C padrão (compatível com C99 / SDCC).
- **Compilador:** SDCC (Small Device C Compiler) devidamente ajustado com otimizações para Z80.
- **Framework & Hardware Abstraction:** **MSXgl** (gerenciamento de VDP V9938/V9958, PSG/AY-3-8910, entradas de teclado/joystick e rotinas de bootstrap).
- **Plataforma Alvo Mínima:**
  - MSX 2 (V9938, CPU Z80A @ 3.579545 MHz).
  - Mínimo de **256 KB de Memory Mapper** instalada.
  - Sistema Operacional **MSX-DOS 2** (versões 2.20 ou superior).
  - 128 KB de VRAM.
- **Formato Primário de Saída:** Executável `GAME.COM` para MSX-DOS 2 acompanhado de arquivos de dados empacotados (`GAME.DAT` / `.BNK`).
- **Formato Secundário (Portabilidade):** Imagem de Cartucho MegaROM (ex: Konami SCC ou ASCII 16KB / ASCII 8KB), usando a mesma camada lógica de abstração de bancos.

---

## 3. Modelo de Vídeo e Apresentação (MSX 2 / V9938)

### 3.1. Escolha do Modo de Vídeo: SCREEN 4 (Graphic 3)
Embora as SCREEN 1 e SCREEN 2 ofereçam facilidades herdadas do MSX 1, o alvo prioritário da nossa engine em MSX 2 é a **SCREEN 4 (Graphic 3 do V9938)**, combinada com alternativas para SCREEN 2 caso necessário:
- **Resolução:** 256 x 192 pixels.
- **Grade Matricial:** 32 x 24 tiles (cada tile com 8 x 8 pixels).
- **Organização da Tabela de Padrões:**
  - Dividida em 3 seções verticais de 8 linhas de caracteres cada (Superior: linhas 0-7, Central: 8-15, Inferior: 16-23).
  - Cada seção pode referenciar até 256 caracteres distintos (totalizando até 768 padrões de 8x8 simultâneos em tela).
- **Atributos de Cor:** 2 cores por linha de 8 pixels de cada caractere (herança do modo Graphic 2, mas com acesso mais flexível via V9938).
- **Vantagem Crítica sobre o MSX 1 (Sprites Modo 2):**
  - Sprites de 16 x 16 pixels com resolução de atributos de cor por linha de varredura.
  - Até 8 sprites por scanline (em vez do limite limitador de 4 do TMS9918), permitindo que heróis, NPCs e projéteis coexistam na mesma linha sem o *flicker* agressivo que destrói RPGs do MSX 1.
  - Não ocorrência de desativação total de sprites na tela.

### 3.2. Layout de Tela Típico (View & HUD)
A tela de 32 x 24 tiles é dividida funcionalmente em zonas estáticas:
```
+-----------------------------------+
|  [ÁREA DE JOGO / SALA ATUAL]      |  Linhas 0 a 17 (32 x 18 tiles)
|  256 x 144 pixels                 |  Viewport de exploração contínua
|  Matriz de salas estilo ZZT/cRPG  |
+-----------------------------------+
|  [HUD / STATUS BAR]               |  Linhas 18 a 19 (32 x 2 tiles)
|  HP: 120/120  MP: 45/45  LV: 03   |  Informações vitais instantâneas
+-----------------------------------+
|  [CAIXA DE DIÁLOGO / MENSAGENS]   |  Linhas 20 a 23 (32 x 4 tiles)
|  "O guarda sussurra: cuidado..."  |  Janela de texto com rolagem/DTE
+-----------------------------------+
```

---

## 4. Design de Memória e Paginação (Memory Mapper & MegaROM)

Este é o diferencial técnico estruturante do projeto. O jogo não cabe em 64 KB de RAM e não deve depender de leituras mecânicas e lentas de disco para cada passo ou diálogo do jogador.

### 4.1. O Desafio dos 64 KB do Z80 e a Arquitetura de Slots do MSX
O Z80 possui um barramento de endereços de 16 bits (0000h a FFFFh), dividido no MSX em quatro páginas de 16 KB:
- **Página 0:** `0000h - 3FFFh`
- **Página 1:** `4000h - 7FFFh`
- **Página 2:** `8000h - BFFFh`
- **Página 3:** `C000h - FFFFh`

No ambiente **MSX-DOS 2**:
- **Página 0:** Ocupada pelos tratadores de interrupção, hooks da BIOS e variáveis internas do DOS 2.
- **Página 3:** Ocupada pela área de trabalho do MSX-DOS 2, buffers de disco, pilha do sistema operacional e pela **Tabela de Vetores do Gerenciador de Mapper** (`EXTBIO` / Jump Table em `F380h+`).
- **Páginas 1 e 2:** Área de Programa Transiente (**TPA** - *Transient Program Area*), onde o binário `.COM` é carregado e executado (iniciando em `0100h`).

### 4.2. Estratégia de Paginação da Engine: "A Janela Aperture em Página 2"
Para manter a engine simples, rápida e livre de corrupção do DOS, definimos uma partição fixa de memória de execução:

```
0000h +---------------------------------------------------------------+
      | Página 0 (0000h - 3FFFh): Kernel MSX-DOS 2 / Rotinas de Sistema|
4000h +---------------------------------------------------------------+
      | Página 1 (4000h - 7FFFh): NÚCLEO DA ENGINE FIXO               |
      |   - Game Loop principal                                       |
      |   - Máquina de Estados / Interpretador de Eventos             |
      |   - Renderizador VDP e HAL                                    |
      |   - Driver de Som / Interrupções VBLANK                       |
8000h +---------------------------------------------------------------+
      | Página 2 (8000h - BFFFh): JANELA DINÂMICA DE PAGINAÇÃO (16KB) |
      |   Segmento X mapeado sob demanda:                             |
      |   [ Banco de Salas ] OU [ Banco de Diálogos ] OU              |
      |   [ Banco de Patterns/Tiles ] OU [ Banco de Áudio ]           |
C000h +---------------------------------------------------------------+
      | Página 3 (C000h - FFFFh): RAM DO SISTEMA & ESTADO GLOBAL      |
      |   - Pilha da Engine (SP)                                      |
      |   - Estruturas de Dados do Jogador (Stats, Inventário, Flags)  |
      |   - Estado das Entidades da Sala Atual                        |
      |   - Vetores de Interrupção e Jump Table do Mapper             |
FFFFh +---------------------------------------------------------------+
```

### 4.3. Interface com o Gerenciador de Memory Mapper do MSX-DOS 2
No MSX-DOS 2, é terminantemente proibido acessar diretamente os registradores de I/O de mapper (`0xFC`, `0xFD`, `0xFE`, `0xFF`) sem sincronização com o kernel, pois o MSX-DOS 2 pode gerenciar múltiplos mappers físicos simultaneamente e reservar segmentos para buffers de disco e RAM disk.

A nossa Engine se integrará obrigatoriamente através da **Tabela de Rotinas do Mapper** disponibilizada pelo MSX-DOS 2 (`EXTBIO` ID 0 / Rotinas em RAM da Página 3):
- `ALL_SEG`: Aloca segmentos livres de 16 KB no início da execução.
- `FRE_SEG`: Libera segmentos ao sair do jogo de volta ao prompt do DOS.
- `PUT_P2`: Mapeia um segmento físico de 16 KB diretamente na **Página 2 (`8000h - BFFFh`)**.
- `GET_P2`: Lê qual segmento está atualmente ativo na Página 2.

### 4.4. Escalonamento Dinâmico de Memória (256 KB a 2 MB+)
No arranque, a Engine interroga o DOS 2 para determinar quantos segmentos livres existem:
1. **Perfil Mínimo (256 KB = 16 segmentos de 16 KB):**
   - ~4 a 6 segmentos reservados pelo MSX-DOS 2 e TPA base.
   - **10 a 12 segmentos livres (~160 KB - 192 KB)** alocados para cache de jogo.
   - A engine carrega os bancos essenciais (mundo inicial, tilesets principais). Se o jogador viaja para uma masmorra distante que não está em RAM, a engine lê o arquivo `GAME.DAT` via MSX-DOS 2 com I/O de alta velocidade e substitui o segmento menos utilizado (**LRU Cache**).
2. **Perfil Confortável (512 KB = 32 segmentos):**
   - 26+ segmentos livres (~416 KB). Praticamente 100% dos dados de um RPG de tamanho médio permanecem residentes em RAM, eliminando leituras de disco após o boot.
3. **Perfil Expandido (1 MB a 2 MB+):**
   - Todos os mapas, todas as variações de gráficos de tiles, sons e a totalidade das árvores de diálogos são pré-carregados durante a tela de abertura. O jogo executa com latência de carregamento zero.

### 4.5. Portabilidade para MegaROM (Cartucho)
A mesma arquitetura de dados alinhada a 16 KB viabiliza a compilação como **MegaROM** (mappers ASCII 16K ou Konami 8K/16K):
- No formato Cartucho, a rotina `MAPPER_SetPage2(segment_id)` é compilada condicionalmente: em vez de chamar o vetor `PUT_P2` do MSX-DOS 2, ela faz uma escrita simples no registrador do mapper de cartucho (ex.: `*(volatile u8*)0x8000 = segment_id;`).
- Toda a lógica da engine, ponteiros e structs permanecem 100% inalterados.

### 4.6. Subsistema de Carregamento de Disco (Binary Disk Loader para DOS 2)
Implementado no módulo `engine_msx/loader.c`, o carregador estabelece o protocolo de streaming de alta performance:
1. **`HEADER.BIN` (Tabela Mestra):**
   - Cabeçalho de 32 bytes validando Magic `ZR01`, versão binária (`1`), total de segmentos de 16 KB no arquivo de dados e parâmetros iniciais do herói (sala, posição X/Y e tileset inicial).
   - Diretório de recursos de 8 bytes por registro mapeando `[Tipo (1B)] -> [Segmento (1B)] -> [ID Lógico (2B)] -> [Offset Página 2 (2B)] -> [Tamanho (2B)]`.
2. **Streaming de `GAME.DAT`:**
   - Utiliza os descritores nativos do MSX-DOS 2 (`DOS_OpenHandle`, `DOS_ReadHandle`, `DOS_CloseHandle`).
   - Para cada segmento de 16 KB indicado no cabeçalho mestre, a engine chaveia o segmento lógico para a Página 2 (`0x8000`) e lê 16.384 bytes diretamente do disco (em 2 blocos seguros de 8.192 bytes).
3. **Resolução de Recursos em Tempo de Execução:**
   - `LOADER_GetTileset(id)` e `LOADER_GetRoom(id)` localizam o recurso no diretório em O(1), chaveiam automaticamente a Página 2 para o segmento correspondente via `MAPPER_SetPage2` e retornam o ponteiro direto em memória (`0x8000 + Offset`).
4. **Finalização Graciosa e Proteção de Sistema:**
   - Ao encerrar a execução, a engine chama `MAPPER_Cleanup()`, devolvendo 100% dos segmentos alocados ao MSX-DOS 2 (`FRE_SEG`).
   - O modo de vídeo é restaurado com chamada de interslot BIOS para `INITXT` (`0x006C`), recarregando o gerador de caracteres e retornando ao prompt `A:\>` sem congelamento ou corrupção de tela.

---

## 5. Arquitetura de Dados: O Banco SQLite do Projeto

O projeto inteiro residirá em um arquivo SQLite estruturado com integridade referencial:

### 5.1. Esquema Relacional Proposto (Resumo das Tabelas Core)

```sql
-- Configurações globais do projeto
CREATE TABLE project_settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

-- Conjunto de padrões de tiles (8x8 pixels)
CREATE TABLE tilesets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT
);

CREATE TABLE tiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tileset_id INTEGER NOT NULL REFERENCES tilesets(id) ON DELETE CASCADE,
    tile_index INTEGER NOT NULL, -- 0 a 255 no banco local
    pattern_bytes BLOB NOT NULL, -- 8 bytes (1 bit por pixel)
    color_bytes BLOB NOT NULL,   -- 8 bytes (atributos de cor por linha)
    collision_type INTEGER DEFAULT 0, -- 0: Passável, 1: Sólido, 2: Água, 3: Dano, etc.
    animation_next_tile_id INTEGER REFERENCES tiles(id)
);

-- Sprites (16x16 pixels - Modo 2 do V9938)
CREATE TABLE sprites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    pattern_bytes BLOB NOT NULL, -- 32 bytes (4 blocos de 8 bytes para 16x16)
    color_bytes BLOB NOT NULL    -- 16 bytes (uma cor por linha de varredura)
);

-- Salas / Telas Estáticas (Grid 32x18 tiles)
CREATE TABLE rooms (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    world_x INTEGER NOT NULL,
    world_y INTEGER NOT NULL,
    name TEXT NOT NULL,
    tileset_id INTEGER NOT NULL REFERENCES tilesets(id),
    tile_matrix BLOB NOT NULL,   -- 576 bytes (32 * 18 indices de tiles)
    north_room_id INTEGER REFERENCES rooms(id),
    south_room_id INTEGER REFERENCES rooms(id),
    east_room_id  INTEGER REFERENCES rooms(id),
    west_room_id  INTEGER REFERENCES rooms(id)
);

-- Atores, Entidades e Triggers posicionados nas salas
CREATE TABLE entities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    room_id INTEGER NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    pos_x INTEGER NOT NULL,      -- Coordenada X no grid (0-31)
    pos_y INTEGER NOT NULL,      -- Coordenada Y no grid (0-17)
    sprite_id INTEGER REFERENCES sprites(id),
    behavior_type INTEGER NOT NULL, -- 0: NPC Estático, 1: NPC Errante, 2: Baú, 3: Gatilho/Trigger, etc.
    event_script_id INTEGER REFERENCES scripts(id)
);

-- Scripts e Roteiros de Eventos
CREATE TABLE scripts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    source_code TEXT NOT NULL,    -- Código fonte em linguagem de alto nível do Editor
    bytecode BLOB                 -- Bytecode compilado para o Z80
);

-- Diálogos e Textos
CREATE TABLE string_table (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    context_tag TEXT NOT NULL,
    text_content TEXT NOT NULL
);

-- Regras e Atributos de RPG
CREATE TABLE hero_classes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    base_hp INTEGER NOT NULL,
    base_mp INTEGER NOT NULL,
    base_atk INTEGER NOT NULL,
    base_def INTEGER NOT NULL
);

CREATE TABLE items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    item_type INTEGER NOT NULL,  -- Arma, Armadura, Poção, Chave, Quest
    modifier_stat INTEGER,
    modifier_value INTEGER,
    price INTEGER DEFAULT 0
);
```

### 5.2. Pipeline de Serialização Binária para o MSX
Durante o processo de Build:
1. O exportador em Go itera sobre as salas e agrupa conjuntos de 10 a 20 salas por **Segmento de 16 KB**.
2. Os textos passam por compressão leve (DTE - *Dual Tile Encoding* ou dicionário de pares de caracteres), otimizando o consumo de RAM.
3. É gerado um arquivo de tabela mestra de cabeçalho (`HEADER.BIN`) contendo o diretório de alocação de recursos:
   `[ID do Recurso] -> [Número do Segmento Lógico] -> [Offset no Segmento (0x0000 - 0x3FFF)]`.

---

## 6. Roteiro de Desenvolvimento (Roadmap Incremental)

O projeto será construído de forma iterativa, validando cada camada antes de avançar.

```
+-------------------------------------------------------------------------+
| FASE 1: Estruturação dos Dados & Abstrações do Editor (Go + SQLite)     |
+-------------------------------------------------------------------------+
                                    |
                                    v
+-------------------------------------------------------------------------+
| FASE 2: Prototipagem de Baixo Nível no MSX (C + MSXgl + SDCC)           |
+-------------------------------------------------------------------------+
                                    |
                                    v
+-------------------------------------------------------------------------+
| FASE 3: Desenvolvimento da GUI Desktop com Fyne (O Editor Visual)       |
+-------------------------------------------------------------------------+
                                    |
                                    v
+-------------------------------------------------------------------------+
| FASE 4: Gameplay Engine & Máquina de Eventos (MSX)                      |
+-------------------------------------------------------------------------+
                                    |
                                    v
+-------------------------------------------------------------------------+
| FASE 5: Pipeline Integrado de Build & Jogo de Referência                |
+-------------------------------------------------------------------------+
```

### Fase 1: Estruturação dos Dados & Abstrações do Editor (Go + SQLite)
- **Subfase 1.1 — Schema & Database Foundation:**
  - Criação do pacote Go de gerenciamento de banco de dados (`pkg/project`).
  - Implementação das migrações SQLite embutidas (`embed.FS`).
  - Criação das rotinas de abertura, criação e validação de integridade de arquivos `.rpgproj`.
- **Subfase 1.2 — Repositórios & Modelos em Go:**
  - Estruturas de dados em Go para `Tileset`, `Tile`, `Sprite`, `Room`, `Entity`, `Script` e `Item`.
  - Operações CRUD completas com testes unitários automatizados.
- **Subfase 1.3 — Engine de Serialização Binária (The Exporter):**
  - Empacotador de matriz de tiles (conversão de estruturas do SQLite para blocos planos de 16 KB).
  - Codificador de padrões VDP (formato binário nativo do V9938).
  - Tabela mestra de índices de segmentos para o runtime MSX.

### Fase 2: Prototipagem de Baixo Nível no MSX (C + MSXgl + SDCC)
- **Subfase 2.1 — Bootstrap MSX-DOS 2 & Abstração de Mapper:**
  - Criação do esqueleto de aplicação MSX-DOS 2 (`.COM`) via MSXgl.
  - Implementação do módulo `dos2_mapper.c`: interrogação de segmentos via `EXTBIO`, alocação segura (`ALL_SEG`), liberação graciosa ao sair (`FRE_SEG`).
  - Rotina de paginação segura na Página 2 (`MAPPER_SetPage2`).
- **Subfase 2.2 — Configuração do VDP V9938 (SCREEN 4):**
  - Inicialização de SCREEN 4 com 3 bancos de padrões.
  - Carregamento de paleta de 16 cores personalizada.
  - Suporte a Sprites Modo 2 (posicionamento, atributos e cores).
- **Subfase 2.3 — Teste de Carga de Sala a partir do Mapper:**
  - Programa de teste no MSX que aloca 256 KB, preenche 4 segmentos com dados simulados de salas e alterna instantaneamente entre eles desenhando os tiles na VRAM.

### Fase 3: Desenvolvimento da GUI Desktop com Fyne (O Editor Visual)
- **Subfase 3.1 — Shell da Aplicação & Navegação:**
  - Janela principal com tema retrô profissional escuro.
  - Gerenciador de projetos (Novo, Abrir, Salvar, Histórico).
  - Painéis de navegação: Tileset Editor, Sprite Editor, Room Editor, Rules & Stats, Script Editor.
- **Subfase 3.2 — Editor de Tiles (8x8):**
  - Grid interativo de desenho pixel-a-pixel.
  - Seletor de cores da paleta MSX (16 cores V9938).
  - Atribuição de propriedades de física/colisão por tile.
- **Subfase 3.3 — Editor de Sprites (16x16):**
  - Edição de sprites Modo 2 com pré-visualização de cores por linha.
  - Configuração de quadros de animação (Walk cycle, Idle, Attack).
- **Subfase 3.4 — Editor de Salas (Room Matrix View):**
  - Área de pintura matricial (32 x 18 tiles) com carimbo de tiles.
  - Ferramentas de pincel, balde de tinta e seleção retangular.
  - Conexão visual entre salas adjacentes (Norte, Sul, Leste, Oeste).
  - Inserção e posicionamento visual de entidades (NPCs, portas, baús).
- **Subfase 3.5 — Editor de Regras, Tabelas de RPG e Roteiros:**
  - Formulários para criação de heróis, monstros, itens e feitiços.
  - Editor textual para scripts simples de gatilho/evento.

### Fase 4: Gameplay Engine & Máquina de Eventos (MSX)
- **Subfase 4.1 — Controle do Herói & Movimentação por Grid:**
  - Leitura de controles (D-Pad/Teclado cursores).
  - Movimentação discreta de 8 em 8 pixels com verificação de colisão contra a tabela de tiles da sala atual.
  - Transição de borda de sala: detecção de saída e acionamento instantâneo do chaveamento de segmento de mapper.
- **Subfase 4.2 — Sistema de Entidades e Atores da Sala:**
  - Spawn dinâmico de até 8 entidades ativas na sala.
  - Atualização de posição e estados de IA simples (NPCs estáticos, patrulheiros).
  - Interação por tecla de ação (barra de espaço / botão 1 do joystick).
- **Subfase 4.3 — Máquina Virtual de Eventos (Bytecode VM):**
  - Interpretador compacto de instruções de aventura:
    - `OP_MSG [str_id]`: Exibe caixa de diálogo na zona inferior.
    - `OP_GIVE_ITEM [item_id]`: Adiciona item ao inventário.
    - `OP_CHECK_FLAG [flag_id] [jump_offset]`: Desvio condicional de fluxo.
    - `OP_SET_FLAG [flag_id] [val]`: Altera estado global do mundo.
    - `OP_TELEPORT [room_id] [x] [y]`: Teletransporte forçado.
- **Subfase 4.4 — Caixa de Diálogo & Overlay de HUD:**
  - Desenho e rolagem de texto na área de diálogo (linhas 20 a 23).
  - Impressão formatada de números de atributos (HP, MP) no HUD.

### Fase 5: Pipeline Integrado de Build & Jogo de Referência
- **Subfase 5.1 — Automação de Build ("One-Click Run"):**
  - O editor aciona o exportador binário, empacota os arquivos de dados, invoca o SDCC/MSXgl e gera o `.COM`.
  - Integração para disparo automático em emulador (ex.: **openMSX** com máquina MSX2 configurada para 256KB/512KB de Mapper).
- **Subfase 5.2 — Backend MegaROM:**
  - Opção no exportador para gerar uma imagem `.ROM` unificada para uso em gravadores de EPROM ou flashcards.
- **Subfase 5.3 — Criação do Jogo Demonstrador de Referência:**
  - Mini-RPG completo de teste contendo masmorra com 20 salas, 5 NPCs com diálogos ramificados, 3 quebra-cabeças com chaves/alavancas e sistema de combate simples.

---

## 7. Próximos Passos Imediatos
Com a **Fase 1 (Estruturação de Dados & Abstrações do Editor)** e a **Fase 2 (Prototipagem de Baixo Nível no MSX)** 100% concluídas e verificadas no openMSX, a sequência de trabalho avança para a **Fase 3: Desenvolvimento da GUI Desktop com Fyne (O Editor Visual)**:
1. **Subfase 3.1 — Shell da Aplicação & Navegação:** Construir a janela mestra da aplicação com tema retrô escuro, gerenciamento de arquivos `.rpgproj` (Criar, Abrir, Salvar) e barra de ferramentas de navegação por abas.
2. **Subfase 3.2 — Editor de Tiles (8x8):** Canvas interativo de desenho pixel-a-pixel com paleta V9938 e atributos de colisão física.
3. **Subfase 3.3 — Editor de Sprites (16x16):** Canvas de edição para Modo 2 com preview de cores por scanline.
4. **Subfase 3.4 — Editor de Salas:** Matriz 32x18 com carimbador de tiles, interligações cardeais e posicionamento visual de entidades.
