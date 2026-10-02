# Z-Realm (`zrealm-msx`)
### *O ZZT dos cRPGs para MSX*

[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](LICENSE)
[![Target: MSX2 / MSX-DOS 2](https://img.shields.io/badge/Platform-MSX2%20%7C%20MSX--DOS%202-red.svg)](SPEC.md)
[![Language: Go](https://img.shields.io/badge/Editor-Go%20%2B%20Fyne%20%2B%20SQLite-00ADD8.svg)](pkg/)
[![Engine: C / MSXgl](https://img.shields.io/badge/Engine-C%20%2B%20SDCC%20%2B%20MSXgl-green.svg)](MSXgl/)

---

## 1. Visão Geral e Conceito

O **Z-Realm** (`zrealm-msx`) é um ecossistema completo de desenvolvimento e execução de jogos de interpretação (cRPGs) e aventuras interativas para computadores da linha **MSX 2** e **MSX 2+** sob o sistema operacional **MSX-DOS 2**.

### Inspirações
* **ZZT (Tim Sweeney, 1991):** Facilidade de criação, exploração orientada a salas contínuas, gatilhos e eventos interativos.
* **Os Clássicos dos cRPGs de 8 bits:** A profundidade tática e a atmosfera envolvente de *Ultima*, *Dragon Quest*, *Rogue* e *The Magic Candle*.

---

## 2. A Filosofia de Design: Compilação Estática vs. Interpretação Lenta

Diferente de engines concebidas para PCs velozes que tentam interpretar scripts pesados e matrizes gigantescas no modesto Z80 a 3.58 MHz, o Z-Realm adota uma abordagem de **compilação e empacotamento estático**:

1. **O Editor Visual Desktop (PC):** Desenvolvido em **Go + Fyne + SQLite**, permite modelar o mundo, desenhar tiles e sprites, traçar roteiros de eventos e definir regras de RPG. O editor valida, otimiza e serializa o banco relacional diretamente para **estruturas binárias compactas alinhadas a blocos de 16 KB**.
2. **A Engine MSX (Target Runtime):** Desenvolvida em **C puro (SDCC + MSXgl)**, opera um loop determinístico orientado a eventos discretos no grid (32x18 tiles), com chaveamento instantâneo de segmentos de memória e sem latência mecânica de disco.

---

## 3. Pilha Tecnológica & Ferramentas

```
+-----------------------------------------------------------------------+
|                         PC / DESKTOP (Editor)                         |
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

### Backend do Editor & Ferramental PC
* **Linguagem:** Go (Go 1.22+)
* **Interface Gráfica:** Fyne v2 (renderização via OpenGL acelerada por GPU)
* **Banco de Dados:** SQLite 3 (Pure Go via `modernc.org/sqlite` - zero dependência de compiladores C no Windows)
* **Controle Transacional:** Foreign Keys estritas (`PRAGMA foreign_keys = ON;`), WAL mode e migrações embutidas (`embed.FS`).

### Engine de Execução no MSX
* **Hardware Alvo:** MSX 2 / MSX 2+ (Z80A @ 3.58 MHz, VDP V9938, Mínimo de 256 KB Memory Mapper).
* **Sistema Operacional:** MSX-DOS 2 (versão 2.20 ou superior).
* **Modo de Vídeo:** **SCREEN 4 (Graphic 3)**:
  * Resolução: 256 x 192 pixels (Grade de 32 x 24 tiles de 8x8).
  * 3 bancos verticais de caracteres (até 768 padrões simultâneos na tela).
  * **Sprites Modo 2:** 16x16 pixels com resolução de atributos de cor por scanline, permitindo até 8 sprites na mesma linha sem desativação nem flicker agressivo de MSX 1.
* **Layout de Apresentação:**
  * **Área Jogável / Viewport:** Linhas 0 a 17 (32 x 18 tiles = matriz fixa de 576 bytes).
  * **HUD / Status Bar:** Linhas 18 a 19 (32 x 2 tiles: HP, MP, nível, status).
  * **Caixa de Diálogo:** Linhas 20 a 23 (32 x 4 tiles para texto narrativo).
* **Gerenciamento de Memória:** Janela de Paginação na **Página 2 (`8000h - BFFFh`)** gerenciada através das rotinas oficiais do Mapper do MSX-DOS 2 (`PUT_P2`, `GET_P2`, `ALL_SEG`, `FRE_SEG`).

---

## 4. Estrutura do Repositório

```text
zrealm-msx/
├── cmd/
│   └── zrealm/            # CLI e utilitário de administração e exportação
├── pkg/
│   ├── exporter/          # Motor de serialização binária para o MSX 2 (V9938/Z80)
│   ├── models/            # Modelos de domínio e cálculos de hardware MSX 2 (V9938)
│   ├── project/           # Gerenciador de projetos (.rpgproj), migrações e demo seeder
│   ├── storage/           # Repositórios de persistência CRUD com SQLite
│   └── version/           # Controle dinâmico de versão semântica (X.Y.Z)
├── engine_msx/            # Runtime C do MSX 2 (V9938 SCREEN 4, Mapper Page 2, Disk Loader)
├── MSXgl/                 # MSX Game Library (biblioteca C para SDCC)
├── dist/                  # Diretório de distribuição gerado pelo build.ps1
├── build.ps1              # Script de automação de build, testes e empacotamento ZIP
├── SPEC.md                # Arquitetura detalhada e especificação técnica mestra
├── OUTLINE.md             # Visão do roadmap, fases e progresso atual
├── MANUAL.md              # Guia completo de instalação, uso e formatos de dados
├── CHANGELOG.md           # Histórico de alterações por versão
├── RELEASE.md             # Notas detalhadas da versão atual
└── LICENSE                # Licença GNU General Public License v3 (GPL 3)
```

---

## 5. Como Compilar e Executar

Consulte o documento [MANUAL.md](MANUAL.md) para o passo a passo completo.

Para compilação rápida com teste automatizado e geração do pacote de distribuição no Windows:
```powershell
./build.ps1
```

O script:
1. Incrementa automaticamente o contador de build `Z` da versão `X.Y.Z`.
2. Executa a suite completa de testes unitários (`go test -v ./...`).
3. Compila os executáveis no diretório `dist/`.
4. Copia os manuais e licenças para `dist/`.
5. Gera o arquivo `.zip` pronto para publicação de release no GitHub.

---

## 6. Licença

Este projeto é software livre distribuído sob os termos da **GNU General Public License v3.0 (GPL 3)**. Consulte o arquivo [LICENSE](LICENSE) para detalhes.
