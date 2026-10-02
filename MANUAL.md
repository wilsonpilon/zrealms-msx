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
* **💾 Exportar:** Compilação dos dados para os binários nativos da Engine (`HEADER.BIN`, `GAME.DAT`).

---

### 3.2. Linha de Comando (CLI Administrativa)
O utilitário `zrealm.exe` também opera em modo de linha de comando para automação em pipelines CI/CD:

#### 3.2.1. Consultar a Versão do Sistema
```bash
./bin/zrealm.exe -version
```
*Saída esperada:*
```text
Z-Realm (zrealm-msx) - v0.1.0
O ZZT dos cRPGs para MSX (MSX 2 / MSX-DOS 2)
```

### 3.2. Criar um Novo Projeto
Para inicializar um novo projeto com o schema SQLite completo, metadados de configuração e o tileset padrão `Overworld`:
```bash
./bin/zrealm.exe -new meu_jogo.rpgproj -name "A Lenda de Valdor"
```

### 3.3. Validar a Integridade de um Projeto
Para checar a integridade estrutural física do SQLite e a integridade de todas as chaves estrangeiras:
```bash
./bin/zrealm.exe -check meu_jogo.rpgproj
```
*Saída esperada:*
```text
Verificando integridade do projeto: meu_jogo.rpgproj...
Projeto: A Lenda de Valdor (v0.1.0)
Plataforma Alvo: MSX2_MSXDOS2
Status de Integridade: OK (Integridade física e Foreign Keys válidas)
```

### 3.4. Exportar Projeto para Binários do MSX 2
Para compilar o banco SQLite para as estruturas binárias nativas alinhadas à Página 2 do MSX-DOS 2 (`HEADER.BIN`, `GAME.DAT` e `SEGxx.BNK`):
```bash
./bin/zrealm.exe -export meu_jogo.rpgproj -out ./build_msx
```
*Saída esperada:*
```text
Exportando projeto para formato nativo MSX 2: meu_jogo.rpgproj...
Exportação concluída com sucesso!
  Tabela Mestra:    ./build_msx/HEADER.BIN
  Dados (GAME.DAT): ./build_msx/GAME.DAT (16384 bytes em 1 segmentos de 16KB)
  Total Recursos:   5 catalogados
```

### 3.5. Gerar Projeto Demonstrativo Completo
Para gerar automaticamente um projeto pronto com masmorra, tileset customizado e 2 salas conectadas:
```bash
./bin/zrealm.exe -demo demo.rpgproj
```

### 3.6. Execução e Teste no Emulador openMSX
A imagem de disquete gerada pelo build contém o sistema de boot completo sob MSX-DOS 2. Para executar no openMSX com perfil de hardware oficial (MSX 2 com 512 KB Mapper e MSX-DOS 2):
```bash
openmsx -machine Philips_NMS_8250 -ext msxdos2 -ext ram512k -diska engine_msx/emul/dsk/DOS2_zrealm.dsk
```
* **Controles no Jogo:**
  * `[ESPAÇO]`: Pula antecipadamente para a próxima sala no loop demonstrativo.
  * `[ESC]`: Encerra a execução e retorna com o modo de texto e RAM restaurados para o prompt `A:\>`.

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

O Z80 enxerga apenas 64 KB de endereçamento por vez. Sob o **MSX-DOS 2**, o Z-Realm utiliza a seguinte divisão fixa:

```
0000h - 3FFFh (Página 0): Kernel MSX-DOS 2 e tratadores de interrupção
4000h - 7FFFh (Página 1): Núcleo da Engine FIXO (Game Loop, VM de eventos, áudio)
8000h - BFFFh (Página 2): JANELA DINÂMICA DE PAGINAÇÃO (16 KB)
C000h - FFFFh (Página 3): Estado global do jogador, RAM de variáveis e EXTBIO
```

### O Chaveamento de Segmentos de 16 KB
1. Os dados do jogo são empacotados em **Segmentos lógicos de 16 KB** (16.384 bytes).
2. Um único segmento de 16 KB comporta confortavelmente **28 salas completas** com seus cabeçalhos.
3. Ao mover o herói para uma nova sala, a engine chama a rotina `PUT_P2` do MSX-DOS 2, chaveando o segmento desejado para o endereço `8000h` em questão de microsegundos, com **zero leitura mecânica de disco**.
4. Em máquinas com 256 KB de Memory Mapper, um cache LRU gerencia os 10 a 12 segmentos livres. Em máquinas com 512 KB ou mais, 100% do mundo reside em RAM.

---

## 6. Solução de Problemas Comuns

* **Erro "o arquivo de projeto já existe":** O Z-Realm protege arquivos existentes contra sobrescrita acidental. Para recriar, forneça um novo nome ou apague o arquivo anterior manualmente.
* **Erro de violação de chave estrangeira ao manipular banco manualmente:** O SQLite do Z-Realm opera com `PRAGMA foreign_keys = ON;`. Certifique-se de que ao inserir um tile ou entidade, o `tileset_id` ou `room_id` correspondente já exista na tabela pai.
* **O script build.ps1 é bloqueado pela política de execução do PowerShell:** Execute uma vez no terminal:
  ```powershell
  Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass
  ```
