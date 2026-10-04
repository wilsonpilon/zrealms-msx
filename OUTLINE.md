# OUTLINE.md — Acompanhamento do Projeto Z-Realm (`zrealm-msx`)

Este documento serve como mapa de bordo vivo do projeto, refletindo o roadmap estabelecido em [SPEC.md](SPEC.md), as fases concluídas, o trabalho em andamento e os próximos marcos a serem executados.

---

## Roadmap Global & Status Atual

| Fase | Descrição | Status | Marco Principal |
| :--- | :--- | :---: | :--- |
| **Fase 1** | **Estruturação de Dados & Abstrações do Editor (Go + SQLite)** | 🟢 **100% Concluída** | Banco `.rpgproj`, Modelos de Hardware, Repositórios & Exporter |
| **Fase 2** | **Prototipagem de Baixo Nível no MSX (C + MSXgl + SDCC)** | 🟢 **100% Concluída** | MSX-DOS 2, Paging Mapper (Page 2), SCREEN 4 e Binary Loader de Disco |
| **Fase 3** | **Desenvolvimento da GUI Desktop com Fyne (O Editor Visual)** | 🟢 **100% Concluída** | Editores de Tilesets, Sprites, Salas 32x18, Roteiros e Regras |
| **Fase 4** | **Gameplay Engine & Máquina de Eventos (MSX)** | 🟢 **100% Concluída** | Controle no Grid, VM de Eventos, HUD e Caixas de Diálogo |
| **Fase 5** | **Pipeline Integrado de Build & Jogo de Referência** | 🟡 **Em Andamento (5.1 OK)** | "One-Click Run", Teste em openMSX e Jogo Demonstrador |

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
- [x] **Subfase 3.1 — Shell da Aplicação & Navegação:**
  - [x] Janela desktop com tema retrô escuro customizado inspirado no MSX 2 e V9938 (`RetroDarkTheme`).
  - [x] Gerenciador de estado reativo (`ProjectState`) com ciclo de vida, observadores e histórico de arquivos recentes.
  - [x] Menu principal (Arquivo, Ferramentas com validação SQLite, Ajuda com especificações MSX 2).
  - [x] Navegação por abas (Salas 32x18, Tilesets 8x8, Sprites 16x16, Diálogos, Regras de RPG e Exportador MSX 2).
  - [x] Barra de status inferior com monitoramento em tempo real de arquivo e estimativa de memória do Memory Mapper.
  - [x] Tela de boas-vindas inicial com atalhos para Novo, Abrir e carregar Masmorra Demo.
  - [x] Integração unificada no CLI (`zrealm`, `zrealm <arquivo.rpgproj>`, retrocompatível com flags `-new`, `-export`, etc.).
- [x] **Subfase 3.2 — Editor de Tiles (8x8):**
  - [x] Paleta oficial de 16 cores do MSX 2 V9938 (`MSXPalette`).
  - [x] Grid interativo de desenho pixel-a-pixel (`TileCanvas`) com suporte a clique e arrasto para alternar, desenhar e apagar.
  - [x] Seletor de cores de primeiro plano (Foreground) e fundo (Background) individual por linha de varredura (8 scanlines).
  - [x] Preview ampliado (64x64) em tempo real com escala nítida pixel-art.
  - [x] Seletor de propriedades físicas e colisão (Passável, Sólido, Água, Dano, Gatilho).
  - [x] Transformações matemáticas de padrão: Girar 90°, Espelhar Horizontal, Espelhar Vertical, Limpar e Preencher.
  - [x] Adição dinâmica de novos tiles até o limite físico de 256 padrões do VDP.
  - [x] Persistência imediata no SQLite com sincronização instantânea de estado.
- [x] **Subfase 3.3 — Editor de Sprites (16x16 Modo 2):**
  - [x] Grid interativo de desenho 16x16 pixels (`SpriteCanvas`) com destaque visual nos 4 quadrantes de 8x8 e fundo xadrez sutil para transparência.
  - [x] Suporte a pintura fluida por clique e arrasto do mouse (ferramentas Alternar, Lápis e Borracha).
  - [x] Seletor individual de cor para cada uma das 16 scanlines do Modo 2 do V9938 com atalho "Aplicar em Todas".
  - [x] Pré-visualizações dinâmicas em escala real 1x (16x16) e 4x (64x64) com pixels nítidos.
  - [x] Transformações completas de sprite: Espelhar Horizontal, Espelhar Vertical (com inversão de scanline de cor), Girar 90°, Deslocamento direcional (Shift Up/Down/Left/Right), Limpar e Preencher.
  - [x] Duplicação inteligente de quadros ("Duplicar Quadro") para criação ágil de ciclos de animação (Walk cycle, Idle, Attack).
  - [x] CRUD completo com criação, edição, persistência no SQLite e exclusão de sprites.
- [x] **Subfase 3.4 — Editor de Salas (Room Matrix View):**
  - [x] Matriz interativa de pintura 32x18 tiles (256x144 pixels) correspondente à geometria da viewport MSX 2 SCREEN 4 no chip V9938.
  - [x] Conjunto completo de ferramentas de edição: Pincel (Carimbo), Balde de Tinta (`FloodFillRoom`), Borracha (Tile 0) e Conta-Gotas (`ToolEyedropper`).
  - [x] Utilitários rápidos de preenchimento: "Limpar (Tile 0)", "Preencher Tudo" e "Preencher Bordas" (`FillBorderRoom`).
  - [x] Paleta de carimbo com seleção visual de tiles do tileset ativo, preview em tempo real e seletor numérico direto (0..255).
  - [x] Painel de Conexões Cardeais (Norte, Sul, Leste, Oeste) com salto rápido ("Ir para Sala"), vinculação manual e algoritmo inteligente de "Auto-Conectar por Coordenadas" (`AutoConnectRooms`).
  - [x] Gestor completo de entidades da sala: inserção, edição, exclusão e posicionamento de atores/NPCs, baús, portas, gatilhos e inimigos, com marcadores visuais sobrepostos diretamente na matriz da sala.
  - [x] Persistência direta no banco SQLite via `RoomEditorWidget` integrado ao `RoomView` com testes unitários cobrindo todos os algoritmos.
- [x] **Subfase 3.5 — Editor de Regras, Tabelas de RPG e Roteiros:**
  - [x] Editor completo de classes de herói com parâmetros base de combate (HP, MP, Ataque, Defesa) e cálculo de tamanho de registro Z80.
  - [x] Gestor do catálogo de itens e equipamentos com categorias (Arma, Armadura, Consumível, Chave, Quest), modificadores de atributos e preços.
  - [x] Simulador de caixa de diálogo com proporção nativa MSX 2 (SCREEN 4: 32 colunas x 4 linhas nas scanlines 20 a 23) e quebra automática de texto (`wrapText`).
  - [x] Compilador e montador de scripts de eventos (`CompileScript`) para a Bytecode VM com opcodes compactos (`OP_MSG`, `OP_GIVE_ITEM`, `OP_TAKE_ITEM`, `OP_SET_FLAG`, `OP_CHECK_FLAG`, `OP_TELEPORT`, `OP_HEAL`, `OP_DAMAGE`, `OP_PLAY_SFX`, `OP_END`).
  - [x] Desassemblador de bytecode (`DisassembleScript`) para engenharia reversa e exibição de bytecode hexadecimal formatado.
  - [x] Integração de CRUD com persistência SQLite reativa em `RulesView` e `ScriptView`.
  - [x] **Conclusão da Fase 3 (Editor Visual Desktop em Fyne):** 100% dos editores visuais operacionais, reativos e cobertos por testes unitários automatizados.

---

### Fase 4: Gameplay Engine & Máquina de Eventos (MSX)
- [x] **Subfase 4.1 — Movimentação do Herói & Colisão no Grid:**
  - [x] Movimentação discreta de 8x8 pixels alinhada ao grid 32x18 (SCREEN 4) com leitura combinada de Joystick (Porta 1) e Teclado (Setas do MSX).
  - [x] Cooldown de repetição de passos (`HERO_STEP_COOLDOWN_FRAMES = 6`) garantindo responsividade imediata no toque inicial e cadência suave ao manter pressionado.
  - [x] Verificação de física e colisão de tiles em tempo real (`COLLISION_SOLID`, `COLLISION_WATER`, `COLLISION_DAMAGE`, `COLLISION_TRIGGER`, `COLLISION_PASSABLE`).
  - [x] Transição cardeal de salas nas bordas (Norte, Sul, Leste, Oeste) com paginação dinâmica automática na Página 2 (`0x8000 - 0xBFFF`) via Memory Mapper (EXTBIOS).
  - [x] Renderização de Sprites Modo 2 (16x16 pixels) com `VDP_LoadSprite`, `VDP_SetSpritePos` e ocultação limpa com `VDP_Screen4_HideSprite`.
  - [x] Saída limpa ao MSX-DOS 2 via tecla `ESC` restaurando modo texto BIOS (`R_INITXT`), desativando sprites e liberando todos os segmentos de RAM com `MAPPER_Cleanup()`.
  - [x] Validação automatizada em emulador openMSX comprovando com 5 screenshots sequenciais: spawn na Sala 1, caminhada a Leste, travessia do portal e carregamento da Sala 2 ("Câmara dos Pilares"), retorno pelo portal Oeste e encerramento limpo ao DOS 2.
- [x] **Subfase 4.2 — Sistema de Entidades e Atores da Sala:**
  - [x] Gerenciador dinâmico de entidades com até 8 instâncias ativas na sala (`entity.h`, `entity.c`).
  - [x] Alocação automática de Sprites Modo 2 (16x16) nos slots 1 a 8 do VDP V9938 com ciclo de vida, carregamento e ocultação limpa.
  - [x] Inteligências Artificiais e comportamentos: NPCs estáticos, NPCs errantes com PRNG Z80 e checagem de colisão mútua, baús sólidos e gatilhos de piso.
  - [x] Interação contextual por tecla de ação unificada (Barra de Espaço no teclado e Gatilho A no Joystick) com debounce e verificação direcional.
  - [x] Sistema de colisão física de entidades integrado ao herói impedindo travessia de NPCs e baús.
  - [x] Validação automatizada no openMSX com 6 screenshots sequenciais comprovando spawn de entidades, colisão sólida contra o Guardião, interação via ESPAÇO, transição para Sala 2, IA do NPC errante e encerramento limpo no MSX-DOS 2.
- [x] **Subfase 4.3 — Máquina Virtual de Eventos (Bytecode VM):**
  - [x] Interpretador de instruções compactas (`OP_NOP`, `OP_MSG`, `OP_GIVE_ITEM`, `OP_TAKE_ITEM`, `OP_SET_FLAG`, `OP_CHECK_FLAG`, `OP_TELEPORT`, `OP_HEAL`, `OP_DAMAGE`, `OP_PLAY_SFX`, `OP_END`).
  - [x] Buffer seguro de execução na RAM da Página 1 (`s_VMScriptBuffer[256]`) imune a chaveamentos de segmento e teleportes.
  - [x] Gerenciamento de estado global: 256 flags de evento, inventário de 16 slots com controle de quantidade e estatísticas de RPG do herói.
  - [x] Validação automatizada em malha fechada no openMSX com 8 screenshots comprovando quest do Guardião, obtenção de itens, desvios condicionais de flags (diálogo alternativo), abertura e esvaziamento do baú com cura (+25 HP) e retorno limpo ao MSX-DOS 2.
- [x] **Subfase 4.4 — Caixa de Diálogos & HUD:**
  - [x] Definição de layout no Banco 2 da SCREEN 4 (V9938): linhas 18-19 reservadas ao HUD e linhas 20-23 reservadas à Caixa de Diálogos e painel de repouso.
  - [x] Carregador de fonte ASCII 8x8 (Texas Instruments TMS9900) e glifos de UI (coração, estrela, chave, cantoneiras e molduras) no Banco 2 da VRAM (`0x1400..0x17FF`) com tabela de cores dedicada (`0x3400..0x37FF`).
  - [x] Renderizador de HUD em tempo real com estatísticas vitais do herói: HP (`♥ xxx/xxx`), MP (`★ xxx/xxx`), Nível (`LV:xx`) e contagem de chaves (`🗝:x`) com rotina de formatação decimal ultrarrápida por subtrações sucessivas sem divisão Z80.
  - [x] Sistema de Diálogos avançado com moldura retangular (`┌─┐│└─┘`), quebra automática de linha (word-wrapping inteligente de até 30 caracteres por linha em 3 linhas por página), indicador de prompt (`[ESPACO] ▼`) e paginação suave para textos longos.
  - [x] Pausa automática de movimentação do herói e comportamentos de entidades durante a exibição de diálogos, com controle de debounce de entrada do teclado e joystick.
  - [x] Painel de repouso ("Standby Panel") exibindo título da área e dica de contexto ao fechar caixas de diálogo.
  - [x] **Conclusão da Fase 4 (Gameplay Engine & Máquina de Eventos no MSX):** 100% das subfases de gameplay (Física/Grid, Entidades/IAs, VM de Eventos e UI/HUD) plenamente operacionais e integradas.

---

### Fase 5: Pipeline Integrado de Build & Jogo de Referência
- [x] **Subfase 5.1 — Automação "One-Click Run":**
  - [x] Criação do pacote `pkg/runner` com abstração completa do ciclo de vida de execução.
  - [x] Empacotador de disquetes MSX-DOS 2 (`.DSK` de 720 KB FAT12) via utilitário `msxtar`.
  - [x] Injeção e extração automática dos binários essenciais de sistema embutidos (`autoexec.bat`, `zrealm.com`, `COMMAND2.COM`, `MSXDOS2.SYS`, `HEADER.BIN` e `GAME.DAT`).
  - [x] Detector e orquestrador dinâmico do emulador `openMSX` (`FindOpenMSX`) com busca resiliente no PATH, Scoop e diretórios convencionais do Windows.
  - [x] Integração completa na GUI desktop Fyne (`ExportView` com botão One-Click Run `[F5]`, seletor de perfil de hardware e console de logs em tempo real).
  - [x] Adição do comando CLI `zrealm -run <arquivo.rpgproj>` (com flags opcionais `-machine`, `-emu` e `-script`).
  - [x] Validação automatizada em malha fechada no openMSX (`sub51_test.tcl`) comprovando com 3 screenshots sequenciais: boot do disquete gerado, gameplay interativo com HUD e saída limpa ao MSX-DOS 2 via ESC com liberação total de RAM.
- [x] **Subfase 5.2 — Backend MegaROM (Cartucho .ROM):**
  - [x] Compilação multi-target da engine C para ambos os formatos: MSX-DOS 2 (`.COM`) e MegaROM ASCII-16 (`.ROM`).
  - [x] Implementação de chaveamento de bancos MegaROM ASCII-16 via escrita em `0x77FF` (`MAPPER_SetPage2`), mapeando a Página 2 (`0x8000..0xBFFF`) diretamente para os bancos de 16 KB do cartucho.
  - [x] Injeção e leitura da Tabela Mestra (`MasterHeader`) e Diretório de Recursos diretamente no Banco 0 do cartucho em `ROM_HEADER_ADDR = 0x7C00` (offset `0x3C00` do arquivo ROM), aproveitando o espaço restante de 1 KB do Banco 0.
  - [x] Criação do pacote Go `pkg/exporter/megarom.go` com embutimento da base ROM (`zrealm_base.rom`), injeção do cabeçalho e append sequencial dos segmentos de dados a partir do Banco 1 (`0x4000`), com preenchimento (padding `0xFF`) para tamanhos padronizados de cartucho (128 KB, 256 KB, 512 KB, 1024 KB, 2048 KB, 4096 KB).
  - [x] Integração completa do pipeline de One-Click Run para cartuchos (`OneClickRunROM`, `BuildOpenMSXROMArgs`, `LaunchOpenMSXCart`) no pacote `pkg/runner`.
  - [x] Atualização da interface gráfica desktop Fyne (`ExportView` e `ProjectState`) com botão dedicado "🕹️ Testar Cartucho MegaROM (.ROM)" e exportação direta com log detalhado de bancos.
  - [x] Adição das flags CLI no binário `zrealm`: `-export-rom <arquivo.rpgproj>`, `-run-rom <arquivo.rpgproj>` e `-pad-size <KB>`.
  - [x] Validação automatizada em malha fechada no emulador openMSX (`sub52_test.tcl`) comprovando com 5 screenshots sequenciais: spawn inicial no cartucho com HUD e SCREEN 4, aproximação do NPC Guardião, diálogo interativo na tela, recebimento da Chave de Bronze no inventário com atualização do HUD, e transição cardeal para a Sala 2 (Câmara dos Pilares) carregada e renderizada a partir do banco MegaROM.
- [x] **Subfase 5.3 — Jogo de Referência Completo ("As Catacumbas de Cristal: O Desafio do Rei Esquecido"):**
  - [x] Criação do jogo de referência completo em `pkg/project/demo.go` com masmorra de 20 salas interligadas em grid 4x5 com travessias cardeais contínuas (Norte, Sul, Leste, Oeste).
  - [x] Catálogo rico de assets: 13 tipos de tiles (Passable, Solid, Trigger, Water, Damage), 8 sprites Modo 2 (Herói, Guardião, Baú, Eremita, Esqueleto, Prisioneiro, Cristal, Goblin), 7 itens e 15 strings de diálogo.
  - [x] Driver de Áudio PSG Nativo (`engine_msx/audio.h` e `audio.c`) para o processador de som AY-3-8910 utilizando acesso direto via portas SDCC `__sfr __at(0xA0)` e `__sfr __at(0xA1)`, com 4 canais de SFX (Item/Chime, Combate/Dano, Porta/Teleporte, Diálogo/Blip).
  - [x] Mecânicas de Combate em Tempo Real (`BEHAVIOR_HOSTILE` tipo 6) em `entity.c`: IA com detecção do herói dentro de raio de 7 tiles, perseguição com pathing Manhattan, ataque corpo a corpo contra o herói (-8 HP com SFX 2 e atualização imediata do HUD), e golpe de espada do herói desferido pela tecla de Ação (abate do inimigo, ocultação do sprite, SFX 1 de vitória e execução do script de drop).
  - [x] Dano de piso em tiles perigosos (`COLLISION_DAMAGE` tipo 3) em `hero.c`: penalidade de -5 HP, atualização do HUD e SFX de perigo.
  - [x] Cadeia de quebra-cabeças progressiva e eventos da Bytecode VM: Guardião Sentinela -> Baú da Armaria (Chave de Ferro e Poção) -> Resgate do Prisioneiro -> Enigma do Eremita -> Fonte Sagrada de Cura (+100 HP) -> Combate na Câmara de Tortura (Chave Real Dourada) -> Santuário do Rei Esquecido (Cristal Primordial de Vitória).
  - [x] Validação automatizada em malha fechada no emulador openMSX (`sub53_test.tcl`) comprovando com 6 capturas de tela em alta fidelidade: spawn na Sala 1, diálogo e chave do Guardião, abertura do baú na Armaria (Sala 4), restauração na Fonte Sagrada (Sala 12), combate em tempo real com o Esqueleto Guerreiro (Sala 15) e conquista do Cristal Primordial no Santuário (Sala 20).
  - [x] **Conclusão da Fase 5 (Pipeline Integrado de Build & Jogo de Referência):** 100% das subfases concluídas com excelência e validação empírica no MSX 2.

