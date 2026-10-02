# RELEASE.md — Detalhamento do Release Oficial

**Versão Atual:** 0.4.0  
**Data:** 02 de Outubro de 2026  
**Status do Release:** Phase 3 Complete (Editor Visual Desktop em Fyne 100% Concluído)  
**Alvo:** Windows (x64) para o Toolkit & GUI / MSX 2 & MSX-DOS 2 para a Engine  

---

## 1. Resumo Executivo do Release

A **Fase 3 (Desenvolvimento da GUI Desktop com Fyne - O Editor Visual)** está **100% concluída**. Este release entrega a aplicação de desktop completa para criação de RPGs no MSX 2, integrando banco de dados relacional SQLite, renderização gráfica com fidelidade aos modos do VDP V9938, ferramentas de desenho em tempo real, gerador de conexões de mundo, gestor de entidades, catálogo de RPG e compilador de scripts da máquina virtual de eventos.

O ecossistema agora oferece a experiência completa de criação:
1. **Tilesets & Tiles 8x8:** Desenho pixel-a-pixel com 16 cores do MSX 2, atributos de colisão e transformações completas.
2. **Sprites 16x16 (Modo 2):** Grid com quadrantes guias, scanlines de cores independentes, previews 1x e 4x e duplicação de frames para ciclos de animação.
3. **Salas 32x18 (SCREEN 4):** Matriz de 576 bytes contíguos com carimbo contínuo, flood fill, borracha, conta-gotas, conexões cardeais automáticas por coordenadas e posicionamento de atores.
4. **Regras, Classes e Itens:** Ficha de heróis (HP/MP/ATK/DEF), inventário e catálogo de equipamentos.
5. **Roteiros & Scripts:** Compilador e desassemblador de scripts de eventos para a Bytecode VM do MSX Z80 e simulador de caixa de diálogo com proporção nativa 32x4 caracteres.

---

## 2. O que foi Criado (Novos Componentes)

### 2.1. Shell da Aplicação & Navegação (Subfase 3.1)
* **`theme.go`:** Tema retrô escuro `RetroDarkTheme` inspirado na paleta V9938 (fundos azuis/grafite escuros, destaques em ciano MSX e dourado, tipografia limpa).
* **`state.go`:** Gerenciador `ProjectState` seguro para concorrência com callbacks reativos, dirty tracking e proteção contra perda de dados.
* **`menu.go`:** Menu mestre da aplicação com verificação de integridade física (`PRAGMA quick_check`) e referencial (`PRAGMA foreign_key_check`) do SQLite.
* **`statusbar.go`:** Barra de rodapé com telemetria em tempo real e estimativa de consumo de blocos de 16 KB no Memory Mapper MSX.

### 2.2. Editor de Tiles 8x8 (Subfase 3.2)
* **`tile_editor.go` & `tile_view.go`:**
  * Grid pixel-a-pixel interativo com paleta V9938 de 16 cores.
  * Seletores de cores individuais para Foreground e Background por linha de 8 pixels.
  * Seletor de física/colisão (Passável, Sólido, Água, Dano, Gatilho).
  * Transformações completas: Rotação 90°, Espelhamento H/V, Deslocamento direcional (Shift), Inversão, Limpar e Preencher.

### 2.3. Editor de Sprites 16x16 Modo 2 (Subfase 3.3)
* **`sprite_editor.go` & `sprite_view.go`:**
  * Grid interativo de 16x16 pixels com divisores de quadrantes 8x8.
  * Seletor de cor individual para cada uma das 16 scanlines do Modo 2 do V9938 com atalho "Aplicar em Todas".
  * Pré-visualizações dinâmicas em escala real 1x (16x16) e ampliada 4x (64x64) com pixels nítidos.
  * Recurso "Duplicar Quadro" para prototipagem ágil de ciclos de animação (Walk/Idle/Attack).

### 2.4. Editor de Salas 32x18 SCREEN 4 (Subfase 3.4)
* **`room_canvas.go` & `room_editor.go` & `room_view.go`:**
  * Viewport interativo de 32x18 tiles (256x144 pixels) correspondente à geometria da viewport MSX SCREEN 4.
  * Conjunto completo de ferramentas: Pincel/Carimbo, Balde de Tinta (`FloodFillRoom`), Borracha (Tile 0) e Conta-Gotas (`ToolEyedropper`).
  * Utilitários rápidos: Limpeza total, Preenchimento total e Preenchimento automático de bordas (`FillBorderRoom`).
  * Paleta de carimbo dinâmica com tiles do tileset ativo, preview do tile selecionado e seletor numérico direto (0..255).
  * Painel de Conexões Cardeais (Norte, Sul, Leste, Oeste) com navegação rápida ("Ir para Sala") e algoritmo de "Auto-Conectar por Coordenadas" (`AutoConnectRooms`).
  * Gestor de Entidades com marcadores visuais sobrepostos diretamente na matriz da sala.

### 2.5. Editor de Regras, Tabelas de RPG e Roteiros (Subfase 3.5)
* **`rules_view.go`:**
  * Gestor de Classes de Heróis: cadastro, edição e exclusão de classes com parâmetros base de combate (HP, MP, Ataque, Defesa).
  * Catálogo de Itens: gerenciamento de armas, armaduras, consumíveis, chaves e quests com modificadores de stat e preços.
* **`script_compiler.go` & `script_view.go`:**
  * Compilador de scripts (`CompileScript`) para a Bytecode VM do MSX Z80 com opcodes compactos (`OP_MSG`, `OP_GIVE_ITEM`, `OP_TAKE_ITEM`, `OP_SET_FLAG`, `OP_CHECK_FLAG`, `OP_TELEPORT`, `OP_HEAL`, `OP_DAMAGE`, `OP_PLAY_SFX`, `OP_END`).
  * Resolução automática de labels e cálculo de offsets relativos de salto.
  * Desassemblador (`DisassembleScript`) e exibição de bytecode hexadecimal formatado.
  * Simulador de caixa de diálogo com proporção nativa MSX (32 colunas x 4 linhas) e quebra automática de texto (`wrapText`).

---

## 3. O que foi Testado & Resultados

### 3.1. Testes Unitários Go (100% de Aprovação)
```text
ok  	github.com/zrealm-msx/zrealm/pkg/exporter           2.61s
ok  	github.com/zrealm-msx/zrealm/pkg/gui                4.62s
ok  	github.com/zrealm-msx/zrealm/pkg/models             0.91s
ok  	github.com/zrealm-msx/zrealm/pkg/project            4.83s
ok  	github.com/zrealm-msx/zrealm/pkg/project/migrations 0.89s
ok  	github.com/zrealm-msx/zrealm/pkg/storage            3.00s
ok  	github.com/zrealm-msx/zrealm/pkg/version            0.82s
```

* **Testes de GUI & Compilador (`pkg/gui`):**
  - `TestRetroDarkTheme`: conformidade da paleta V9938.
  - `TestProjectStateLifecycle` & `TestProjectStateDemo`: reatividade e integridade.
  - `TestFloodFillRoom`, `TestFillBorderRoom`, `TestClearRoomMatrix`: algoritmos matriciais de sala 32x18.
  - `TestAutoConnectRooms`: topologia de conexão automática em grade de salas.
  - `TestRoomCanvasPainting`: ferramentas de pintura e marcadores visuais.
  - `TestCompileScriptBasic`, `TestCompileScriptWithLabelsAndBranches`, `TestCompileScriptErrors`: montador de bytecode da VM.
  - `TestDisassembleAndFormatHex`: desassemblagem e hex formatting.
  - `TestWrapText`: quebra de texto na caixa de diálogo de 32 colunas.
  - `TestSpriteCanvasPixelPainting`, `TestSpriteTransformations`: Modo 2 do V9938.
  - `TestTileCanvasPixelPainting`, `TestPatternTransformations`: padrões 8x8 do V9938.

---

## 4. Próximos Passos (Fase 4: Gameplay Engine & Máquina de Eventos no MSX)

Com as Fases 1, 2 e 3 100% concluídas, o desenvolvimento avança para a **Fase 4**:
1. **Subfase 4.1 — Movimentação do Herói & Colisão no Grid:** Controle direcional, colisão contra tiles sólidos/água e transição de sala nas bordas.
2. **Subfase 4.2 — Sistema de Entidades e Atores:** Spawning dinâmico de até 8 entidades ativas, IAs simples e interação por tecla de ação.
3. **Subfase 4.3 — Máquina Virtual de Eventos (Bytecode VM no Z80):** Execução do interpretador de instruções compactas compiladas pelo editor.
4. **Subfase 4.4 — Caixa de Diálogos & HUD:** Renderização de mensagens nas linhas 20-23 e mostrador de HP/MP nas linhas 18-19.
