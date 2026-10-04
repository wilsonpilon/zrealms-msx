# MANUAL.md — Manual do Desenvolvedor & Usuário do Z-Realm

Bem-vindo ao manual oficial do **Z-Realm (`zrealm-msx`)** — o ecossistema de criação de cRPGs para computadores MSX 2 e MSX 2+.

---

## 1. Pré-Requisitos do Sistema

Para desenvolver, compilar e executar o ecossistema Z-Realm em seu computador, você precisará de:

### Ambiente PC (Desenvolvimento do Editor & Tooling)
* **Sistema Operacional:** Windows 10/11 (64-bit), Linux ou macOS.
* **Go Compiler:** Go versão 1.22 ou superior ([golang.org](https://go.dev/)).
* **PowerShell:** PowerShell 5.1 ou PowerShell 7+ (para execução do script de build automatizado no Windows).
* **Git:** Para controle de versão ([git-scm.com](https://git-scm.com/)).

### Ambiente de Teste MSX (Opcional para Emulação)
* **Emulador MSX:** **openMSX** ([openmsx.org](https://openmsx.org/)) configurado com uma máquina MSX2 ou MSX2+ brasileira, europeia ou japonesa (ex.: Philips NMS 8250, Gradiente Expert 2.0, Sony HB-F700) com no mínimo 256 KB de Memory Mapper.
* **Sistema de Disco:** Disco virtual com **MSX-DOS 2.20+** (`MSXDOS2.SYS` e `COMMAND2.COM`).

---

## 2. Como Baixar e Compilar

### 2.1. Clonando o Repositório
Abra o seu terminal e execute:
```bash
git clone https://github.com/zrealm-msx/zrealm.git
cd zrealm
```

### 2.2. Compilação Automatizada (Recomendado)
No Windows, execute o script PowerShell na raiz do projeto:
```powershell
./build.ps1
```

O que o script faz automaticamente:
1. **Incrementa o número de build (`Z`)** da versão semântica `X.Y.Z` (ex.: de `0.1.0` para `0.1.1`).
2. Executa a suite completa de testes unitários (`go test -v ./...`).
3. Compila os binários executáveis otimizados no diretório `dist/bin/`.
4. Copia a documentação de distribuição (`README.md`, `LICENSE`, `MANUAL.md`, `RELEASE.md`) para `dist/`.
5. Compacta tudo em um arquivo `.zip` na pasta `dist/` (ex.: `dist/zrealm-msx-v0.1.1-windows-amd64.zip`), pronto para distribuição ou publicação de release no GitHub.

#### Parâmetros do Script de Build
* **Incremento de Feature (`Y`):**
  ```powershell
  ./build.ps1 -BumpFeature
  ```
* **Incremento de Grande Fase (`X`):**
  ```powershell
  ./build.ps1 -BumpMajor
  ```
* **Compilação rápida sem executar os testes unitários:**
  ```powershell
  ./build.ps1 -SkipTests
  ```

### 2.3. Compilação Manual via Go
Caso prefira compilar os utilitários manualmente pelo Go:
```bash
go build -o bin/zrealm.exe ./cmd/zrealm
```

Para rodar todos os testes unitários com saída verbosa:
```bash
go test -v ./...
```

---

## 3. Guia de Uso do Editor Visual Desktop & CLI

### 3.1. Iniciar o Editor Gráfico Desktop (GUI Fyne)
Para abrir a interface gráfica completa do Z-Realm, basta executar o binário compilado sem argumentos ou dar um duplo-clique no executável:
```powershell
./dist/bin/zrealm.exe
# ou na raiz após o build:
./zrealm.exe
```
O editor abrirá com o tema retrô escuro do MSX 2 e abas integradas:
* **🧱 Tilesets (8x8):** Desenho pixel-a-pixel com paleta V9938 de 16 cores, cores de frente e fundo por scanline, atributos de colisão física e transformações completas.
* **👾 Sprites (16x16 Modo 2):** Editor com 16 scanlines de cores independentes, visualização em escala real 1x e 4x, e recurso de duplicação de quadros para ciclos de animação.
* **🗺️ Salas (32x18 SCREEN 4):** Matriz contígua de 576 bytes com carimbo contínuo, flood fill, borracha, conta-gotas, conexões cardeais automáticas e posicionamento de entidades/atores.
* **⚔️ Regras & RPG:** Gestão de classes de personagens (HP, MP, Ataque, Defesa) e catálogo de itens/equipamentos.
* **📜 Scripts & Diálogos:** Editor com compilador e desassemblador de scripts para a Bytecode VM do Z80 e simulador de caixa de diálogo com proporção nativa MSX (32x4 caracteres).
* **💾 Exportar & One-Click Run:** Compilação dos dados para binários nativos (`HEADER.BIN`, `GAME.DAT`), cartuchos MegaROM (`.ROM` ASCII-16), empacotamento em disco MSX-DOS 2 (`.dsk` de 720 KB) e execução com 1 clique no openMSX (`[F5]` ou menu *Ferramentas -> 🚀 Executar no openMSX (F5)*).

---

### 3.2. Linha de Comando (CLI Administrativa)
O utilitário `zrealm.exe` também opera em modo de linha de comando para automação em pipelines CI/CD e testes rápidos:

#### 3.2.1. Consultar a Versão do Sistema
```bash
./bin/zrealm.exe -version
```
*Saída esperada:*
```text
Z-Realm (zrealm-msx) - v0.5.2
O ZZT dos cRPGs para MSX (MSX 2 / MSX-DOS 2 / MegaROM)
```

#### 3.2.2. Criar um Novo Projeto
Para inicializar um novo projeto com o schema SQLite completo, metadados de configuração e o tileset padrão `Overworld`:
```bash
./bin/zrealm.exe -new meu_jogo.rpgproj -name "A Lenda de Valdor"
```

#### 3.2.3. Validar a Integridade de um Projeto
Para checar a integridade estrutural física do SQLite e a integridade de todas as chaves estrangeiras:
```bash
./bin/zrealm.exe -check meu_jogo.rpgproj
```
*Saída esperada:*
```text
Verificando integridade do projeto: meu_jogo.rpgproj...
Projeto: A Lenda de Valdor (v0.5.2)
Plataforma Alvo: MSX2_MSXDOS2
Status de Integridade: OK (Integridade física e Foreign Keys válidas)
```

#### 3.2.4. Exportar Projeto para Binários do MSX 2 (MSX-DOS 2)
Para compilar o banco SQLite para as estruturas binárias nativas alinhadas à Página 2 do MSX-DOS 2 (`HEADER.BIN`, `GAME.DAT` e `SEGxx.BNK`):
```bash
./bin/zrealm.exe -export meu_jogo.rpgproj -out ./build_msx
```

#### 3.2.5. Exportar Projeto para Cartucho MegaROM (.ROM ASCII-16)
Para compilar o banco SQLite e montar um cartucho MegaROM `.ROM` pronto para rodar em emuladores ou gravar em FlashROM/EPROM física:
```bash
./bin/zrealm.exe -export-rom meu_jogo.rpgproj -out ./build_msx -pad-size 128
```
*Saída esperada:*
```text
Exportando projeto para Cartucho MegaROM ASCII-16 (.ROM): meu_jogo.rpgproj...
Exportação MegaROM concluída com sucesso!
  Arquivo .ROM:     ./build_msx/meu_jogo.rom (131072 bytes / 128 KB)
  Tipo de Mapper:   ASCII16 (8 bancos de 16KB)
  Segmentos Jogo:   1 banco(s)
  Total Recursos:   9 catalogados
```

#### 3.2.6. Automação "One-Click Run" via Disquete (MSX-DOS 2)
Para compilar, gerar o disco virtual `.dsk` inicializável de 720 KB e disparar o openMSX em um único comando:
```bash
./bin/zrealm.exe -run meu_jogo.rpgproj
```

#### 3.2.7. Automação "One-Click Run" via Cartucho MegaROM (.ROM)
Para compilar, montar o cartucho `.ROM` e disparar o openMSX diretamente com `-cart`:
```bash
./bin/zrealm.exe -run-rom meu_jogo.rpgproj -machine Philips_NMS_8250
```

Opções adicionais de emulação via linha de comando:
* `-machine <perfil>`: Seleciona o perfil de hardware (`Philips_NMS_8250` [padrão] ou `Panasonic_FS-A1GT`).
* `-emu <caminho>`: Especifica o executável do emulador (caso não esteja no `PATH`).
* `-script <arquivo.tcl>`: Executa um script de automação TCL do openMSX (ideal para testes de regressão e CI/CD).
* `-pad-size <KB>`: Fixa o tamanho do cartucho ROM gerado (128, 256, 512, 1024, etc.; padrão 0 = automático).

Exemplo avançado com script de teste:
```bash
./bin/zrealm.exe -run-rom demo.rpgproj -script engine_msx/sub52_test.tcl
```

#### 3.2.8. Gerar Projeto Demonstrativo Completo
Para gerar automaticamente um projeto pronto com masmorra, tileset customizado e 2 salas conectadas:
```bash
./bin/zrealm.exe -demo demo.rpgproj
```

#### 3.2.9. Execução Manual no openMSX
Caso deseje montar e executar manualmente no openMSX:
* **Modo Disquete DOS2:**
  ```bash
  openmsx -machine Philips_NMS_8250 -ext msxdos2 -ext ram512k -diska build_msx/zrealm.dsk
  ```
* **Modo Cartucho MegaROM:**
  ```bash
  openmsx -machine Philips_NMS_8250 -cart build_msx/meu_jogo.rom
  ```
* **Controles na Engine:**
  * `[Setas]` / `[WASD]`: Movimentação do herói pelo grid de 32x18 tiles (SCREEN 4).
  * `[ESPAÇO]`: Interação com NPCs, baús e avanço de diálogos.
  * `[ESC]`: Encerra a sessão da engine (retorna ao DOS 2 em disquete ou reinicia na sala inicial em cartucho).

---

## 4. Formato dos Arquivos & Modelo de Dados

O Z-Realm armazena a totalidade de um jogo em um arquivo único com extensão **`.rpgproj`**, que é um banco de dados relacional **SQLite 3** formatado com transações atômicas e foreign keys ativadas. Você pode inspecionar e editar este arquivo com ferramentas gráficas como o **DB Browser for SQLite** ([sqlitebrowser.org](https://sqlitebrowser.org/)).

### 4.1. Estrutura das Tabelas Principais

```
+-------------------------------------------------------------+
|                     ARQUIVO .rpgproj                        |
+-------------------------------------------------------------+
|  schema_migrations  : Histórico de migrações DDL aplicadas  |
|  project_settings   : Parâmetros globais (chave -> valor)   |
|  tilesets           : Catálogos de ambientes (8x8 tiles)    |
|  tiles              : Padrão binário V9938 (8B) + Cores (8B)|
|  sprites            : Modo 2 V9938 (32B padrão + 16B cores) |
|  rooms              : Matriz 32x18 (576B) + Conexões cardeais|
|  entities           : NPCs, Baús, Portas e Gatilhos na sala |
|  scripts            : Código de alto nível + Bytecode Z80   |
|  string_table       : Textos e diálogos do jogo             |
|  hero_classes       : Classes de personagens e atributos base|
|  items              : Catálogo de armas, armaduras e poções |
+-------------------------------------------------------------+
```

### 4.2. Formato dos Gráficos no Hardware MSX 2 (V9938)

#### Tiles (8 x 8 pixels — SCREEN 4 / Graphic 3)
Cada tile é composto por dois buffers binários de **8 bytes** cada:
* **`pattern_bytes` (8 bytes):** 1 bit por pixel. O bit 7 representa o pixel mais à esquerda (X=0) e o bit 0 o pixel mais à direita (X=7). Cada byte representa uma linha vertical (Y=0 a Y=7).
* **`color_bytes` (8 bytes):** Atributos de cor para cada uma das 8 linhas de varredura. Cada byte armazena duas cores de 4 bits (paleta de 16 cores):
  * Nibble alto (bits 7..4): Cor de primeiro plano (*Foreground*).
  * Nibble baixo (bits 3..0): Cor de fundo (*Background*).

#### Sprites (16 x 16 pixels — Modo 2 do V9938)
O Modo 2 do V9938 supera o limitador clássico do MSX 1 permitindo até 8 sprites por scanline e cores individuais por scanline:
* **`pattern_bytes` (32 bytes):** Organizado na ordem de memória nativa do VDP V9938 em duas colunas verticais de 16 bytes:
  * Bytes `0..15`: Coluna esquerda do sprite (X: 0..7, Y: 0..15).
  * Bytes `16..31`: Coluna direita do sprite (X: 8..15, Y: 0..15).
* **`color_bytes` (16 bytes):** 1 byte por scanline (Y: 0..15), definindo a cor de cada uma das 16 linhas do sprite.

### 4.3. Formato das Salas (Rooms)
Cada sala possui uma matriz contígua de exatamente **576 bytes** (`tile_matrix`):
* Grade de **32 colunas x 18 linhas**.
* Cada byte contém o índice do tile (`0` a `255`) correspondente ao tileset ativo da sala.
* Conexões cardeais automáticas: Campos `north_room_id`, `south_room_id`, `east_room_id` e `west_room_id` interligam salas vizinhas, permitindo transições contínuas de tela quando o herói atravessa as bordas do grid.

---

## 5. Estratégia de Paginação na Memória do MSX (Página 2)

O Z80 enxerga apenas 64 KB de endereçamento por vez. O Z-Realm implementa uma estratégia simétrica e elegante tanto no ambiente de disco (MSX-DOS 2) quanto em cartucho (MegaROM):

### 5.1. Ambiente MSX-DOS 2 (Memory Mapper)
```
0000h - 3FFFh (Página 0): Kernel MSX-DOS 2 e tratadores de interrupção
4000h - 7FFFh (Página 1): Núcleo da Engine FIXO (Game Loop, VM de eventos, áudio)
8000h - BFFFh (Página 2): JANELA DINÂMICA DE PAGINAÇÃO (16 KB)
C000h - FFFFh (Página 3): Estado global do jogador, RAM de variáveis e EXTBIO
```
1. Os dados do jogo são empacotados em **Segmentos lógicos de 16 KB** (16.384 bytes).
2. Um único segmento de 16 KB comporta confortavelmente **28 salas completas** com seus cabeçalhos.
3. Ao mover o herói para uma nova sala, a engine chama a rotina `PUT_P2` do MSX-DOS 2 via `MAPPER_SetPage2`, chaveando o segmento desejado para o endereço `8000h` em questão de microsegundos, com **zero leitura mecânica de disco**.
4. Em máquinas com 256 KB de Memory Mapper, um cache LRU gerencia os 10 a 12 segmentos livres. Em máquinas com 512 KB ou mais, 100% do mundo reside em RAM.

### 5.2. Ambiente Cartucho MegaROM (ASCII-16)
No formato de cartucho MegaROM padrão **ASCII-16**, a organização de memória é otimizada para execução em hardware puro:
```
0000h - 3FFFh (Página 0): MSX BIOS / Sistema Operacional
4000h - 7FFFh (Página 1): Banco 0 FIXO do Cartucho (Código da Engine + Tabela Mestra em 7C00h)
8000h - BFFFh (Página 2): JANELA DINÂMICA DE BANCOS DO CARTUCHO (Bancos 1 a N)
C000h - FFFFh (Página 3): Memória RAM do MSX (Estado global, variáveis de jogo e stack)
```
1. **Banco 0 Fixo (`4000h..7FFFh`):** Contém o ponto de entrada do cartucho (`4000h` assinatura `AB`), o código compilado da engine e a **Tabela Mestra** (`MasterHeader`) gravada estaticamente no endereço `7C00h` (offset `0x3C00` do arquivo ROM).
2. **Chaveamento por Hardware ASCII-16:** A comutação de banco na Página 2 é realizada escrevendo o número do banco desejado no registrador de controle em `0x77FF`:
   ```c
   Poke(0x77FF, logicalSegment + 1);
   ```
3. **Desempenho Instantâneo:** Ao cruzar um portal de sala, a nova sala e seus recursos tornam-se imediatamente visíveis na CPU em um único ciclo de escrita I/O (`OUT`/`POKE`), sem custo de transferência ou alocação de RAM.
4. **Padronização Comercial:** O exportador Go (`pkg/exporter/megarom.go`) alinha o cartucho com padding `0xFF` para os tamanhos físicos comerciais: 128 KB, 256 KB, 512 KB, 1 MB, 2 MB ou 4 MB.

---

## 6. Solução de Problemas Comuns

* **Erro "o arquivo de projeto já existe":** O Z-Realm protege arquivos existentes contra sobrescrita acidental. Para recriar, forneça um novo nome ou apague o arquivo anterior manualmente.
* **Erro de violação de chave estrangeira ao manipular banco manualmente:** O SQLite do Z-Realm opera com `PRAGMA foreign_keys = ON;`. Certifique-se de que ao inserir um tile ou entidade, o `tileset_id` ou `room_id` correspondente já exista na tabela pai.
* **O script build.ps1 é bloqueado pela política de execução do PowerShell:** Execute uma vez no terminal:
  ```powershell
  Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass
  ```
