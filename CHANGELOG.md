# CHANGELOG.md — Registro de Mudanças do Z-Realm

Todas as alterações notáveis neste projeto serão documentadas neste arquivo.

O formato baseia-se em [Keep a Changelog](https://keepachangelog.com/pt-BR/1.0.0/) e este projeto adota as seguintes regras para versionamento **X.Y.Z**:
* **Z (Build / Patch):** Incrementado automaticamente a cada compilação executada pelo script `build.ps1`.
* **Y (Minor / Feature):** Incrementado a cada nova feature concluída e integrada ao projeto.
* **X (Major):** Incrementado a cada transição estrutural ou conclusão de uma grande fase (ex.: finalização da Camada de Dados, conclusão do Editor Gráfico, etc.).

## [0.5.1] - 2026-10-04

### Adicionado (Added)
- **Subfase 5.1: Automação "One-Click Run" (Exportação + Montagem DSK + Boot no openMSX):**
  - **Pacote `pkg/runner` (`runner.go`, `disk.go`, `assets.go`):**
    - Abstração do ciclo completo de execução com `OneClickRun` e `RunOptions`.
    - Empacotador automático de disquetes MSX-DOS 2 de 720 KB FAT12 (`PackDsk`) utilizando o utilitário nativo `msxtar`.
    - Embutimento dos binários essenciais de sistema (`autoexec.bat`, `zrealm.com`, `COMMAND2.COM`, `MSXDOS2.SYS`) via `embed.FS`, garantindo que o toolkit seja 100% autônomo.
    - Detector dinâmico e resiliente do executável `openmsx.exe` (`FindOpenMSX`) com busca em variáveis de ambiente, PATH, Scoop (`%USERPROFILE%\scoop\shims\openmsx.exe`) e diretórios padrão do Windows.
    - Inspetor de conteúdo de disquete `VerifyDskContents` validando catálogo completo de arquivos do disco.
  - **Interface Gráfica Fyne Desktop (`ExportView`, `menu.go`, `state.go`):**
    - Nova seção de destaque **"🚀 Automação 'One-Click Run' (openMSX)"** com botão de ação rápida `[F5]` e feedback visual imediato.
    - Seletor de hardware emulado permitindo alternar entre o perfil padrão `Philips_NMS_8250 (MSX 2 + 512KB Mapper)` e `Panasonic_FS-A1GT (MSX 2+ / MSX turbo R)`.
    - Campo opcional para especificar caminho customizado do emulador.
    - Item de menu "Testar no openMSX (One-Click Run)..." adicionado ao menu Ferramentas.
    - Console de logs em tempo real na interface exibindo cada etapa do pipeline (exportação, empacotamento DSK e disparo do emulador).
  - **Interface de Linha de Comando (CLI - `cmd/zrealm/main.go`):**
    - Flag `-run <arquivo.rpgproj>` permitindo disparo do One-Click Run diretamente pelo terminal.
    - Flags complementares `-machine <nome>`, `-emu <caminho>` e `-script <arquivo.tcl>` para automação e testes headless.
  - **Validação Automatizada em Malha Fechada no openMSX (`sub51_test.tcl`):**
    - Execução do script emulado no modelo `Panasonic_FS-A1GT` capturando 3 screenshots sequenciais:
      1. `sub51_01_oneclick_boot.png`: Boot direto a partir do disquete `.dsk` recém-gerado pelo pipeline, com masmorra, sprites e HUD ativos.
      2. `sub51_02_oneclick_gameplay.png`: Movimentação do herói e abertura de diálogo com o Guardião.
      3. `sub51_03_oneclick_exit.png`: Saída limpa ao MSX-DOS 2 via tecla ESC com restauração de modo de texto e liberação total de RAM.

## [0.4.4] - 2026-10-03

### Adicionado (Added)
- **Subfase 4.4: Caixa de Diálogos & HUD no VDP V9938 (Conclusão da Fase 4):**
  - **Módulo de Interface e HUD (`ui.h`, `ui.c`):**
    - **Layout VRAM no Banco 2 da SCREEN 4:**
      - Linhas 18-19 reservadas ao HUD fixo em tempo real.
      - Linhas 20-23 reservadas à Caixa de Diálogos interativa e ao Painel de Repouso ("Standby").
      - Fonte ASCII 8x8 completa baseada no padrão Texas Instruments TMS9900 carregada nos padrões 128..223 da VRAM (`0x1400..0x16FF`) com atributos de cor correspondentes (`0x3400..0x36FF`).
      - Glifos customizados de UI mapeados nos padrões 224..239: molduras e cantoneiras (`┌ ─ ┐ │ └ ┘`), ícones coloridos de RPG (Coração Rosa `♥`, Estrela Azul `★`, Chave Dourada `🗝`) e seta indicadora de paginação (`▼`).
    - **HUD Dinâmico em Tempo Real:**
      - Renderização contínua das estatísticas vitais do herói: Pontos de Vida (`♥ xxx/xxx`), Mana (`★ xxx/xxx`), Nível (`LV:xx`) e contagem dinâmica de chaves (`🗝:x`).
      - Algoritmo ultrarrápido de conversão decimal sem operações de divisão por hardware (`__divuint`/`__moduint`), utilizando loops de subtração de alta performance para Z80.
      - Disparo de atualização imediata do HUD sempre que a Máquina Virtual executa `VM_OP_GIVE_ITEM`, `VM_OP_TAKE_ITEM`, `VM_OP_HEAL` ou `VM_OP_DAMAGE`.
    - **Caixa de Diálogos com Word-Wrapping & Paginação:**
      - Janela textual com moldura de alta definição em ciano sobre fundo preto (colunas 1 a 30, linhas 19 a 23).
      - Algoritmo de quebra automática de palavras (*word-wrapping*) limitando o texto a 30 caracteres por linha e 3 linhas úteis por página.
      - Paginação automática para mensagens que excedem 3 linhas, exibindo prompt com seta pulsante (`[ESPACO] ▼`) para avanço e fechamento com debounce seguro.
      - Pausa total da movimentação do herói e IAs de entidades durante a exibição de caixas de diálogo.
    - **Painel de Repouso ("Standby Panel"):**
      - Exibição de contexto automático após o término de diálogos, apresentando o título do reino ("Z-REALM: CATACUMBAS") e instruções de ação ("`[ESPACO] INTERAGIR / ACAO`").
  - **Aprimoramentos de Toolchain & Compilador:**
    - Correção do pipeline de compilação em `MSXgl/engine/script/js/compiler.js`: sanitização de paths com remoção de trailing slashes nos parâmetros `-I` do `sdasz80.exe` e remoção automática de diretivas espúrias `!extern` geradas pelo SDCC 4.6.0 em modo `c1mode` que causavam estouro de pilha (`0xC00000FD`).
  - **Validação Automatizada em Malha Fechada no openMSX (`sub44_test.tcl`):**
    - Execução do script emulado no modelo `Panasonic_FS-A1GT` capturando 9 screenshots sequenciais:
      1. `sub44_01_spawn_hud.png`: Spawn inicial com HUD ativo (`♥ 075/100`, `★ 030/030`, `LV:01`, `🗝:0`) e painel de repouso.
      2. `sub44_02_facing_guardian.png`: Aproximação ao Guardião no grid em (13, 9).
      3. `sub44_03_guardian_dialogue.png`: Abertura da caixa de diálogo do Guardião com moldura ciano e texto quebrado; atualização imediata do contador de chaves no HUD (`🗝:1`).
      4. `sub44_04_dialogue_closed.png`: Fechamento do diálogo e restauração do painel de repouso.
      5. `sub44_05_guardian_branch_dialogue.png`: Reinteração com o Guardião disparando diálogo ramificado por Flag 1.
      6. `sub44_06_facing_chest.png`: Posicionamento em (8, 5) em frente ao Baú.
      7. `sub44_07_chest_heal_hud.png`: Abertura do baú com diálogo descritivo e cura em tempo real de +25 HP no HUD (`♥ 100/100`).
      8. `sub44_08_chest_empty_dialogue.png`: Reinspeção do baú confirmando diálogo de baú vazio.
      9. `sub44_09_dos_exit.png`: Saída limpa ao MSX-DOS 2 via tecla ESC sem resíduos de vídeo ou memória.

## [0.4.3] - 2026-10-03

### Adicionado (Added)
- **Subfase 4.3: Máquina Virtual de Eventos (Bytecode VM no Z80):**
  - **Módulo da Máquina Virtual (`vm.h`, `vm.c`):**
    - Interpretador Z80 de alta eficiência para os 11 opcodes canônicos da especificação:
      - `OP_NOP (0x00)`: Nenhuma operação.
      - `OP_MSG (0x01)`: Carrega e exibe string pelo ID lógico via Memory Mapper.
      - `OP_GIVE_ITEM (0x02)`: Adiciona item ao inventário com empilhamento de quantidade.
      - `OP_TAKE_ITEM (0x03)`: Remove item do inventário e compacta slots contíguos.
      - `OP_SET_FLAG (0x04)`: Altera estado global de uma das 256 flags.
      - `OP_CHECK_FLAG (0x05)`: Desvio condicional relativo (`int16 offset`) caso a flag esteja ativa (`!= 0`).
      - `OP_TELEPORT (0x06)`: Teletransporte forçado de sala e coordenadas `(X, Y)`.
      - `OP_HEAL (0x07)`: Restauração de HP do herói respeitando o teto de `MaxHP`.
      - `OP_DAMAGE (0x08)`: Subtração de HP do herói com proteção de piso em 0.
      - `OP_PLAY_SFX (0x09)`: Disparo de efeitos sonoros com `DOS_Beep()` imune a chamadas diretas de BIOS não mapeada.
      - `OP_END (0xFF)`: Finalização de execução de script.
    - Buffer seguro de execução na RAM residente da Página 1 (`s_VMScriptBuffer[256]`), copiando o bytecode antes de executar para garantir imunidade contra trocas de segmento no Memory Mapper (Página 2) causadas por `OP_TELEPORT` ou `WORLD_LoadRoom`.
    - Gerenciamento de estado global:
      - Tabela de 256 flags de evento globais (`g_VMFlags[256]`).
      - Inventário do jogador com 16 slots (`g_VMInventory[16]`), armazenando `ItemID` (16-bit) e `Quantity` (8-bit) com rotinas de busca, contagem e remoção sem uso de `memcpy`.
      - Estatísticas do herói (`g_HeroStats`): HP, MaxHP, MP, MaxMP, Nível, Ataque e Defesa.
  - **Integração no Loader e Exporter:**
    - Definição de `ResTypeScript = 6` em [pkg/exporter/types.go](file:///e:/zrealms-msx/pkg/exporter/types.go).
    - Suporte a indexação direta O(1) de Strings em [pkg/exporter/exporter.go](file:///e:/zrealms-msx/pkg/exporter/exporter.go).
    - Rotinas de cópia segura com preservação de paginação no mapper: `LOADER_CopyScript()` e `LOADER_CopyString()`.
  - **Integração nas Entidades e Gameplay:**
    - [engine_msx/entity.c](file:///e:/zrealms-msx/engine_msx/entity.c): Acionamento automático de `VM_ExecuteScript(ent->EventScriptID)` ao interagir (`ENTITY_InteractAt`) ou pisar em gatilhos (`ENTITY_CheckStepTrigger`).
    - [engine_msx/zrealm.c](file:///e:/zrealms-msx/engine_msx/zrealm.c): Inicialização da VM de eventos no boot do jogo (`VM_Init()`).
  - **Validação Automatizada em Malha Fechada no openMSX (`sub43_test.tcl`):**
    - Controlador em malha fechada via leitura direta de memória RAM no openMSX inspecionando coordenadas, flags, inventário, HP e IDs de strings.
    - 8 screenshots sequenciais comprovando:
      1. `sub43_01_spawn.png`: Spawn inicial com herói em (16, 9), flags zeradas e inventário vazio.
      2. `sub43_02_facing_guardian.png`: Aproximação a Oeste até (13, 9) encarando o Guardião.
      3. `sub43_03_guardian_quest.png`: Disparo do Script 1: obtenção da Chave de Ferro, Flag 1 definida para 1, carregamento da Mensagem 2 e execução de SFX.
      4. `sub43_04_guardian_branched.png`: Segunda interação com o Guardião: detecção de Flag 1 ativa e desvio condicional para a Mensagem 3 sem duplicar itens.
      5. `sub43_05_facing_chest.png`: Navegação no grid até (8, 5) encarando o Baú em (8, 4).
      6. `sub43_06_chest_opened.png`: Abertura do baú: obtenção da Poção de Vida, cura de 25 HP (75 -> 100 HP), Flag 2 definida para 1 e Mensagem 4.
      7. `sub43_07_chest_empty.png`: Segunda interação no baú aberto: desvio condicional para a Mensagem 5 ("O baú está vazio").
      8. `sub43_08_dos_exit.png`: Encerramento limpo via tecla ESC com restauração de vídeo e liberação integral de memória RAM no MSX-DOS 2.

---

## [0.4.2] - 2026-10-03

### Adicionado (Added)
- **Subfase 4.2: Sistema de Entidades, Atores e IAs da Sala (MSX 2 Gameplay Engine):**
  - **Módulo de Entidades da Sala (`entity.h`, `entity.c`):**
    - Gerenciador com capacidade de até 8 instâncias ativas simultâneas (`MAX_ACTIVE_ENTITIES = 8`).
    - Atribuição dinâmica de hardware no VDP V9938: slots de sprites 1 a 8 em Modo 2 (16x16 pixels com cor por scanline), mantendo slot 0 exclusivo para o Herói.
    - Ciclo de vida integrado: rotina `ENTITY_LoadRoomEntities` que limpa entidades da sala anterior e carrega dinamicamente novas entidades a partir de `BinaryRoom.Entities` via Memory Mapper (Página 2: `0x8000-0xBFFF`).
    - Finalização graciosa com `ENTITY_Cleanup()`, ocultando os 8 sprites de hardware do VDP.
  - **Inteligências Artificiais e Comportamentos (Behavior Types):**
    - `BEHAVIOR_STATIC_NPC`: NPCs estacionários com colisão sólida e resposta a interação por proximidade.
    - `BEHAVIOR_WANDERING_NPC`: NPCs errantes com passos autônomos no grid, acionados por PRNG Z80 leve, verificação de limites do mapa, passabilidade de tiles (`COLLISION_PASSABLE`), colisão contra o herói e contra outras entidades ativas.
    - `BEHAVIOR_PATROL_NPC`: NPCs patrulheiros em rota contínua de vaivém com reversão ao encontrar obstáculos ou entidades.
    - `BEHAVIOR_CHEST`: Baú de tesouro com colisão sólida intransponível e alternância de estado (aberto/fechado) ao interagir.
    - `BEHAVIOR_DOOR`: Porta ou passagem com bloqueio físico transitável conforme estado.
    - `BEHAVIOR_TRIGGER`: Gatilhos invisíveis de piso acionados instantaneamente ao pisar (`ENTITY_CheckStepTrigger`).
  - **Colisão Física de Entidades no Grid:**
    - Verificação de colisão `ENTITY_IsSolidAt` integrada a `HERO_Update()`: o herói não atravessa NPCs ou baús.
    - Entidades ativas em movimento verificam colisão mútua e colisão contra a posição do herói.
  - **Interação por Tecla de Ação (Barra de Espaço / Gatilho do Joystick):**
    - Suporte unificado à Barra de Espaço (`KEY_SPACE`) no teclado e Botão 1 (`JOY_INPUT_TRIGGER_A`) no Joystick da Porta 1.
    - Temporizador suave de debounce (`ActionCooldown = 15 frames`, ~0.25s).
    - Cálculo direcional de abordagem à frente do herói com fallback contextual sob os pés para acionar `ENTITY_InteractAt()`.
  - **Assets e Projeto Demonstrador (`demo.go`):**
    - Adicionado Sprite 2 ("Guardião") com armadura ciano, elmo dourado e lança.
    - Adicionado Sprite 3 ("Baú de Tesouro") com ferragens douradas e corpo castanho.
    - Sala 1: Guardião em `(12, 9)` e Baú em `(8, 4)`.
    - Sala 2: Sentinela Errante em `(6, 6)` e Baú Místico em `(21, 9)`.
  - **Validação Automatizada no openMSX (`sub42_test.tcl`):**
    - Sequência automatizada com 6 screenshots gravados em disco:
      1. `sub42_01_spawn_entities.png`: Spawn na Sala 1 com Herói, Guardião e Baú renderizados simultaneamente.
      2. `sub42_02_collision_guardian.png`: Herói caminhando a Oeste e colisão física sólida bloqueando avanço em `(13, 9)`.
      3. `sub42_03_interact_guardian.png`: Acionamento da tecla ESPACO e interação contextual com o Guardião.
      4. `sub42_04_room2_entities.png`: Travessia para a Sala 2 via Memory Mapper com limpeza dos sprites da Sala 1 e spawn das entidades da Sala 2.
      5. `sub42_05_wandering_npc.png`: Sentinela Errante movendo-se autonomamente pelo grid SCREEN 4.
      6. `sub42_06_dos_clean_exit.png`: Retorno limpo ao MSX-DOS 2 via ESC, desativação de sprites e 100% da RAM liberada.

---

## [0.4.1] - 2026-10-03

### Adicionado (Added)
- **Subfase 4.1: Movimentação do Herói & Colisão no Grid (MSX 2 Gameplay Engine):**
  - **Módulo do Herói (`hero.h`, `hero.c`):**
    - Entidade do Herói com coordenadas discretas de grid (`TileX`, `TileY`), coordenadas de tela (`PixelX`, `PixelY`), orientação direcional (`HERO_DIR_UP`, `HERO_DIR_DOWN`, `HERO_DIR_LEFT`, `HERO_DIR_RIGHT`) e estado de movimento (`Moved`).
    - Cooldown de repetição de passos configurável (`HERO_STEP_COOLDOWN_FRAMES = 6`), proporcionando latência zero no primeiro acionamento e taxa contínua de ~10 passos por segundo.
    - Suporte simultâneo e higienizado a Joystick na Porta 1 (`JOY_PORT_1`) e Teclado (Setas direcionais do MSX via `Keyboard_IsKeyPressed`).
  - **Módulo de Mundo e Transição de Salas (`world.h`, `world.c`):**
    - Gerenciador global `g_World` com buffers locais para os 576 tiles da sala e 256 bytes de colisão física.
    - Verificação de colisão física em tempo real via `WORLD_GetCollision`: suporte a `COLLISION_SOLID`, `COLLISION_WATER`, `COLLISION_DAMAGE`, `COLLISION_TRIGGER` e `COLLISION_PASSABLE`.
    - Transição de bordas e portais cardeais (Norte, Sul, Leste, Oeste) via `WORLD_CheckRoomTransition`: paginação automática do segmento no Memory Mapper (Página 2: `0x8000 - 0xBFFF`), recarga dos dados da sala e reposicionamento automático na borda oposta.
  - **Sprites V9938 Modo 2 (`vdp_screen4.h`, `vdp_screen4.c`):**
    - Configuração de tabelas de padrões (`0x3800`), atributos (`0x1E00`) e cores (`0x1C00`) para sprites 16x16.
    - Funções de renderização: `VDP_LoadSprite` (32 bytes de pattern + 16 bytes de color per-scanline), `VDP_SetSpritePos` e ocultação via `VDP_Screen4_HideSprite`.
  - **Loop Principal e Encerramento Limpo (`zrealm.c`):**
    - Sincronização a 50/60 Hz no V-Blank via interrupções (`Halt()`).
    - Tecla `ESC` para encerramento gracioso: ocultação de sprites, restauração do modo texto BIOS (`R_INITXT`) e liberação de 100% dos segmentos alocados no Memory Mapper (`MAPPER_Cleanup()`).
  - **Validação Automatizada com openMSX:**
    - Script de teste `sync_trace.tcl` sincronizado com o renderizador VDP gerando 5 screenshots comprovando:
      1. Spawn do herói no centro da Sala 1 `(16, 9)`.
      2. Caminhada a Leste rumo ao portal.
      3. Chegada à Sala 2 ("Câmara dos Pilares") via chaveamento do Memory Mapper.
      4. Retorno pelo portal Oeste de volta à Sala 1.
      5. Retorno limpo ao prompt do MSX-DOS 2 com mensagem de sucesso e memória liberada.

---

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
