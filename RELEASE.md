# RELEASE.md — Detalhamento do Release Oficial

**Versão Atual:** 0.3.0  
**Data:** 02 de Outubro de 2026  
**Status do Release:** Phase 2 Complete (Prototipagem de Baixo Nível no MSX 2 Concluída)  
**Alvo:** Windows (x64) para o Toolkit / MSX 2 & MSX-DOS 2 para a Engine  

---

## 1. Resumo Executivo do Release

A **Fase 2 (Prototipagem de Baixo Nível no MSX)** está 100% concluída. Este release entrega a **Engine Nativa do MSX 2 (`engine_msx/`)** em C (SDCC 4.6.0 + MSXgl), o subsistema de **Memory Mapper na Página 2 (`0x8000 - 0xBFFF`)** integrado ao MSX-DOS 2, o controlador gráfico de **SCREEN 4 (Graphic 3 do V9938)** e o **Binary Disk Loader (`loader.c`)**.

O pipeline está totalmente conectado de ponta a ponta:
1. O comando `zrealm -demo demo.rpgproj` gera um projeto SQLite com tileset de masmorra e 2 salas conectadas.
2. O comando `zrealm -export demo.rpgproj -out engine_msx/emul/dos2` exporta os dados binários para `HEADER.BIN` e `GAME.DAT`.
3. O build do MSX (`build.bat`) gera o executável `zrealm.com` e a imagem de disquete `DOS2_zrealm.dsk` com todos os dados.
4. O emulador openMSX carrega o jogo no DOS 2, pagina os segmentos na Página 2, renderiza as duas salas e retorna com texto 100% limpo ao prompt do DOS, liberando toda a RAM alocada.

---

## 2. O que foi Criado (Novos Componentes)

### 2.1. Engine MSX 2 & MSX-DOS 2 Core (`engine_msx/`)
* **`zrealm.c`:** Núcleo de inicialização e game loop em C para MSX-DOS 2:
  * Inicialização via `crt0_dos.asm` (início em `0x0100`).
  * Tratamento de teclado sem conflito com chamadas BDOS (leitura direta via `Keyboard_IsKeyPressed`).
  * Finalização limpa restaurando o modo de texto BIOS (Screen 0) via `DOS_InterSlotCall(g_EXPTBL[0], R_INITXT)`.
* **`mapper.c` / `mapper.h`:** Gerenciador de Memory Mapper para MSX-DOS 2:
  * Localização dinâmica da Jump Table do DOS 2 (`EXTBIOS`).
  * Alocação de segmentos via `ALL_SEG` e liberação via `FRE_SEG`.
  * Paginação de segmentos de 16 KB na Página 2 (`0x8000 - 0xBFFF`) via `PUT_P2`.
  * Zero vazamento de memória (todas as estruturas alocadas são devolvidas ao DOS 2).
* **`vdp_screen4.c` / `vdp_screen4.h`:** Controlador de Vídeo V9938 SCREEN 4:
  * Ativação de Graphic 3 com 3 bancos verticais de padrões (2048 bytes cada) e cores (2048 bytes cada).
  * Renderização direta de Viewport de 32x18 tiles (576 bytes contíguos na Name Table `0x1800`).
  * Limpeza dedicada de áreas de HUD (linhas 18-19) e Diálogo (linhas 20-23) com tile preto.
  * Ocultamento de sprites Modo 2 (Y=216).
* **`loader.c` / `loader.h`:** Carregador de Disco Binário para MSX-DOS 2:
  * Abertura e leitura com descritores nativos (`DOS_OpenHandle`, `DOS_ReadHandle`, `DOS_CloseHandle`).
  * Validação do cabeçalho mestre `HEADER.BIN` (Magic `ZR01`, versão 1).
  * Streaming contínuo de `GAME.DAT` em blocos de 8 KB diretamente para o segmento mapeado na Página 2.
  * Resolução de recursos com chaveamento automático de segmento (`LOADER_GetTileset`, `LOADER_GetRoom`).

### 2.2. Extensões no Toolkit Go (`cmd/zrealm` e `pkg/project`)
* **`pkg/project/demo.go`:** Gerador do projeto demonstrativo `CreateDemoProject`:
  * Tileset com Chão de masmorra (Tile 0), Parede de tijolos (Tile 1) e Portal místico dourado (Tile 2).
  * Sala 1 ("Entrada das Catacumbas"): Paredes perimetrais com saída no Leste (Tile 2).
  * Sala 2 ("Câmara dos Pilares"): Conexão no Oeste, 4 pilares maciços e altar cruciforme central.
  * Sprite do Herói (16x16) e Diálogo de boas-vindas.
* **Flag `-demo` no CLI `zrealm`:** Comando direto para gerar projetos de teste.

### 2.3. Automação de Testes e Empacotamento
* **`project_config.js`:** Adição de `HEADER.BIN` e `GAME.DAT` em `DiskFiles` para geração automática do disquete `DOS2_zrealm.dsk` via `msxtar`.
* **`boot_test.tcl`:** Script TCL do openMSX configurado para capturar automaticamente:
  * Sala 1 (Entrada com portal dourado).
  * Sala 2 (Câmara dos Pilares via troca dinâmica de segmento de Mapper).
  * Saída limpa ao prompt do MSX-DOS 2.
* **`build.ps1`:** Empacotamento dos binários MSX (`DOS2_zrealm.dsk`, `zrealm.com`, `HEADER.BIN`, `GAME.DAT`) no diretório `dist/msx/` e no pacote ZIP de release.

---

## 3. O que foi Corrigido / Otimizado

* **Interferência BDOS 0x06:** Substituição de `DOS_PollKey()` por leitura direta de teclado de hardware (`Keyboard_IsKeyPressed`), evitando que chamadas de polling BDOS corrompam registradores Z80 no meio de modos gráficos do VDP.
* **Restauração de Fonte no MSX-DOS 2:** Uso de `DOS_InterSlotCall(g_EXPTBL[0], R_INITXT)` para recarregar a fonte BIOS do gerador de caracteres e restaurar completamente a Screen 0 ao sair para o DOS 2.
* **Limpeza do HUD e Diálogo:** Ajuste de `VDP_ClearHUDAndDialogue` para usar Tile 255 (preto/vazio), evitando que o chão pontilhado do jogo invada as áreas reservadas de status e texto.
* **Configuração de Boot openMSX:** Inclusão obrigatória de `-ext msxdos2` em conjunto com `-ext ram512k` para a máquina `Philips_NMS_8250`.

---

## 4. O que foi Testado & Resultados

### 4.1. Testes Unitários Go (100% de Aprovação)
```text
ok  	github.com/zrealm-msx/zrealm/pkg/exporter	2.53s
ok  	github.com/zrealm-msx/zrealm/pkg/models	0.87s
ok  	github.com/zrealm-msx/zrealm/pkg/project	4.30s
ok  	github.com/zrealm-msx/zrealm/pkg/project/migrations	0.84s
ok  	github.com/zrealm-msx/zrealm/pkg/storage	2.40s
ok  	github.com/zrealm-msx/zrealm/pkg/version	0.83s
```

### 4.2. Testes de Execução no openMSX
* **Máquina:** `Philips_NMS_8250` + `ram512k` + `msxdos2`.
* **Boot:** Limpo via `autoexec.bat` chamando `zrealm.com`.
* **Leitura de Disco:** `HEADER.BIN` (72 bytes) e `GAME.DAT` (16.384 bytes) lidos sem erros.
* **Alocação de Mapper:** 16 segmentos alocados com sucesso.
* **Renderização:** Sala 1 e Sala 2 renderizadas com 100% de fidelidade ao schema SQLite exportado.
* **Retorno:** Prompt `A:\>` restaurado sem travamento, congelamento ou vazamento de memória.

---

## 5. Próximos Passos (Fase 3)

Com a fundação de dados (Fase 1) e o runtime de baixo nível do MSX 2 (Fase 2) concluídos e integrados, o projeto avança para a **Fase 3: Desenvolvimento da GUI Desktop com Fyne (O Editor Visual)**:
1. **Subfase 3.1 — Shell da Aplicação & Navegação:** Criação da janela mestre, seletor de arquivos `.rpgproj` e layout de abas temáticas.
2. **Subfase 3.2 — Editor de Tiles (8x8):** Canvas de edição com paleta do MSX 2 e marcação de colisão.
3. **Subfase 3.3 — Editor de Sprites (16x16):** Visualizador de camadas Modo 2.
4. **Subfase 3.4 — Editor de Salas:** Matriz 32x18 com carimbo de tiles e conexão cardeal visual.
