# 📜 Resumo Completo da Sessão & Especificação de Continuidade
## Projeto: Editor de Adventures (MSX Clean-Room Edition)
> **Data:** 25 de Setembro de 2026  
> **Referência Histórica:** Renato Degiovani (1986), *Sistema Editor de Adventures Versão 3.4*  
> **Plataforma:** MSX (Screen 0, 40×24) | PC (Go, Compilador & TUI)

---

## 🎯 1. Resumo Executivo das Implementações Realizadas

Nesta sessão, avançamos substancialmente no ecossistema do **Editor de Adventures**, resolvendo problemas de jogabilidade, introduzindo documentação histórica completa, modernizando a interface gráfica clássica e criando um inovador sistema de acentuação ergonômico para o MSX.

### Principais Entregas:

1. **OCR Integral do Manual Histórico (1986):**
   - Transcrição completa e conversão para Markdown do manual original de Renato Degiovani: [`docs/editor_adventure.md`](docs/editor_adventure.md).
   - Preservação de todos os 11 capítulos, registradores, tabela dos 45 bytecodes, regras do byte de consistência e o jogo de exemplo *MANSÃO*.

2. **Correção do Game Loop da Engine C (MSX):**
   - **Exibição de Saídas:** Corrigida a renderização das saídas visíveis (`N, S, L, O...`) na descrição da sala no Campo Central.
   - **Telas Introdutórias:** Implementação da exibição da tela de introdução, instruções e contextualização histórica antes de iniciar o jogo, com mensagens padronizadas de *"Pressione qualquer tecla..."* e pausa ativa de teclado para permitir a leitura.
   - **Comandos Canônicos do Sistema:**
     - `VERBO` / `VERBOS`: Lista todos os verbos conhecidos pelo vocabulário da aventura.
     - `INSTRUCAO` / `INSTRUCOES`: Reexibe a tela com as instruções e regras gerais.
     - `DICA` / `DICAS`: Dispara as dicas contextuais preparadas pelo autor do adventure.

3. **Novo Sistema de Acentuação e Teclado MSX:**
   - **13 Letras Acentuadas:** Suporte completo a `À, Á, Â, Ã, Ç, É, Ê, Í, Ó, Ô, Õ, Ú, Ü`.
   - **Glifo Customizado para `Ü`:** Desenhado e carregado dinamicamente na VRAM do MSX no slot `0x9F` (substituindo o símbolo `£`), completando o alfabeto acentuado.
   - **Normalização Sintática:** O parser converte as letras acentuadas para seus equivalentes sem acento internamente, preservando a escrita exata para o jogador sem quebrar correspondência com o dicionário.

4. **Janela [TAB] — Inserção Rápida de Acentos:**
   - Caixa de diálogo centralizada na tela (largura 32, colunas 4 a 35).
   - Moldura gráfica dos anos 1990 com cantos e divisórias contínuas.
   - Grade 5×3 perfeitamente alinhada em ordem alfabética estrita:
     - Linha 1: `À`, `Á`, `Â`, `Ã`, `Ç`
     - Linha 2: `É`, `Ê`, `Í`, `Ó`, `Ô`
     - Linha 3: `Õ`, `Ú`, `Ü`
   - Navegação interativa com as **setas do teclado** ($\leftarrow, \rightarrow, \uparrow, \downarrow$), cursor com colchetes `[Á]` e inserção com **ENTER**.

5. **Janela [SELECT] — Configuração dos Atalhos `Shift+1`..`Shift+0`:**
   - Grade de dois níveis (topo: 13 acentos, base: 10 atalhos).
   - **Fluxo em 2 Fases:**
     - **Fase 1 (Atalhos):** O cursor percorre os atalhos `1:Á` a `0:Ç`. Pressiona-se ENTER ou o número correspondente.
     - **Fase 2 (Acentos):** O atalho escolhido é marcado (`>1:Á<`), o cursor sobe para a grade de acentos, o usuário navega com as setas e confirma com ENTER.
     - O novo atalho entra em vigor imediatamente para uso com `Shift+1`..`Shift+0`.
   - **Otimização no Compilador Go:** O compilador analisa a frequência das palavras acentuadas de cada jogo e define os 10 atalhos iniciais ideais automaticamente.

6. **Acabamento Visual IBM-PC dos Anos 1990 & Alinhamento Rígido:**
   - Identificação dos caracteres semigráficos na ROM do MSX via análise de `images/editadv-02.png` e dump de VRAM (`vram.dat`):
     - Cantos: `0x81` (Sup-Esq), `0x9A` (Sup-Dir), `0xA6` (Inf-Esq), `0xA7` (Inf-Dir)
     - Barras: `0x5F` (Horizontal), `0x5E` (Vertical)
   - Resolução definitiva do efeito torto/desalinhado: cada célula das janelas possui largura exata de 6 caracteres (`col * 6`), sincronizando milimetricamente o topo e a base da janela.

7. **Build & Deploy:**
   - Compilação dos binários do clássico *Amazônia* (`amazonia.dsk`, `amazonia.com`, `advent.rom`) com SDCC e MSXgl.

---

## 📐 2. Especificação Técnica de Detalhes Internos

### 2.1 Mapeamento de Códigos de Teclado MSX
| Tecla | Código BIOS / ASCII | Ação |
|---|---|---|
| `TAB` | `9` (`0x09`) | Abre a janela de inserção rápida de acentos |
| `SELECT` | `24` (`0x18`) | Abre a janela de configuração dos atalhos Shift |
| `ESC` | `27` (`0x1B`) | Cancela/fecha janelas de diálogo |
| `ENTER` | `13` (`0x0D`) | Confirma seleção de caractere ou atalho |
| `Seta Esquerda` | `29` (`0x1D`) | Cursor para a esquerda |
| `Seta Direita` | `28` (`0x1C`) | Cursor para a direita |
| `Seta Cima` | `30` (`0x1E`) | Cursor para cima |
| `Seta Baixo` | `31` (`0x1F`) | Cursor para baixo |
| `Shift + 1` | `33` (`!`) | Atalho slot 0 |
| `Shift + 2` | `64` (`@`) | Atalho slot 1 |
| `Shift + 3` | `35` (`#`) | Atalho slot 2 |
| `Shift + 4` | `36` (`$`) | Atalho slot 3 |
| `Shift + 5` | `37` (`%`) | Atalho slot 4 |
| `Shift + 6` | `94` (`^`) | Atalho slot 5 |
| `Shift + 7` | `38` (`&`) | Atalho slot 6 |
| `Shift + 8` | `42` (`*`) | Atalho slot 7 |
| `Shift + 9` | `40` (`(`) | Atalho slot 8 |
| `Shift + 0` | `41` (`)`) | Atalho slot 9 |

### 2.2 Caracteres Gráficos da Moldura MSX
- `0x81`: Canto Superior Esquerdo
- `0x9A`: Canto Superior Direito
- `0xA6`: Canto Inferior Esquerdo
- `0xA7`: Canto Inferior Direito
- `0x5F`: Linha Horizontal
- `0x5E`: Linha Vertical
- `0x9F`: Glifo personalizado de VRAM para `Ü` maiúsculo

---

## 🗂️ 3. Mapeamento dos Arquivos Modificados / Criados

| Arquivo | Descrição das Modificações |
|---|---|
| [`engine/src/ui.c`](engine/src/ui.c) | Implementação de `UI_DrawWindowFrame`, `UI_DrawWindowDivider`, `UI_DrawAccentCell`, `UI_DrawShortcutCell`, navegação por setas e novos diálogos `UI_SelectAccentedChar` e `UI_ConfigShiftShortcuts`. |
| [`engine/src/ui.h`](engine/src/ui.h) | Declarações dos novos métodos de interface e atalhos. |
| [`engine/src/font_custom.h`](engine/src/font_custom.h) | Definição do bitmap customizado do `Ü` (`0x9F`). |
| [`engine/src/parser.c`](engine/src/parser.c) | Adicionada normalização de `Ü` para `U` no parser sintático. |
| [`engine/src/interpreter.c`](engine/src/interpreter.c) | Implementação dos comandos canônicos `VERBOS`, `INSTRUCAO` e `DICA`. |
| [`engine/src/game_loop.c`](engine/src/game_loop.c) | Listagem correta de saídas na sala e fluxo de introdução com pausa. |
| [`compiler/compiler.go`](compiler/compiler.go) | Otimização estatística de atalhos de acentos baseada em frequência no texto e suporte ao `Ü`. |
| [`MSXgl/projects/advent/*`](MSXgl/projects/advent/) | Sincronização dos códigos C e script de compilação SDCC. |
| [`amazonia/amazonia.dsk`](amazonia/amazonia.dsk) | Imagem de disco DOS1 com o jogo Amazônia pronto para execução. |
| [`amazonia/amazonia.com`](amazonia/amazonia.com) | Executável MSX-DOS do jogo Amazônia. |
| [`docs/editor_adventure.md`](docs/editor_adventure.md) | OCR integral do manual oficial de 1986 de Renato Degiovani. |
| [`manual.md`](manual.md) | Atualizado com Seções 13 (Acentuação/Atalhos), 14 (Comandos de Sistema) e 15 (Manual Histórico). |
| [`manual_do_jogador.md`](manual_do_jogador.md) | Atualizado com instruções para o jogador das teclas TAB, SELECT e atalhos Shift. |
| [`README.md`](README.md) | Atualizado com novidades, créditos e documentação dos novos recursos. |

---

## 🚀 4. O Que Falta Fazer / Próximos Passos Recomendados

Ao retomar os trabalhos em outro computador ou em uma nova sessão, as seguintes atividades podem ser priorizadas:

1. **Apresentação Gráfica de Telas SCR (Opcional):**
   - O manual original prevê telas gráficas opcionais de abertura (`.SCR`). Atualmente a engine trabalha 100% em modo texto Screen 0. Uma rotina de carregamento de imagem de abertura em Screen 2 antes de alternar para Screen 0 pode ser adicionada se desejado.
2. **Exportação Direta de DSK pelo TUI:**
   - Adicionar ao menu F9 da ferramenta de terminal (`edadv.exe`) a opção de gerar diretamente o arquivo `.dsk` do jogo além da `.rom`.
3. **Novas Histórias e Testes de Comunidade:**
   - Testar emuladores adicionais (WebMSX, BlueMSX, fMSX) e cartuchos reais (Carnivore2, MegaFlashROM SCC+).
   - Validar com novos adventures criados pela comunidade no formato YAML.

---

*Documento gerado automaticamente para registro de sessão e continuidade de desenvolvimento.*
