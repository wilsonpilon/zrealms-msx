# Tutorial: Criando a Aventura "A Mansão Misteriosa" no Editor TUI

Este guia detalha o passo a passo completo para recriar a aventura de demonstração **"A Mansão Misteriosa"** do zero utilizando a ferramenta TUI (Text User Interface) do compilador do projeto.

---

## 🚀 1. Como Iniciar o Editor TUI

No terminal, execute o editor informando o arquivo YAML onde o jogo será criado/editado:

```powershell
go run ./compiler/cmd/main.go games/minha_aventura.yaml
```

> **Atalhos Globais do Editor:**
> - Teclas **`1` a `5`**: Alternam entre as abas do editor:
>   - **`1`**: Mapa e Salas
>   - **`2`**: Objetos e Byte de Consistência
>   - **`3`**: Comandos e Funções (Bytecode)
>   - **`4`**: Mensagens do Jogo
>   - **`5`**: Build e Geração da ROM MSX
> - **`Tab` / `Shift+Tab`**: Navega entre campos, listas e botões da tela
> - **`Enter`**: Seleciona ou ativa um botão/campo
> - **`F2`**: Salva o arquivo YAML
> - **`F5`**: Compila a aventura para código C (`game_data.h` / `game_data.c`)
> - **`F9`**: Compila a ROM MSX 32K (`advent.rom`)
> - **`Esc`**: Volta o foco para a lista principal da aba

---

## 🗺️ 2. Aba 1: Criando as Salas do Mapa (Tecla `1`)

Na aba **1. Mapa**, defina as 4 salas da mansão com suas conexões cardeais. Para cada sala, selecione na lista da esquerda (ou adicione) e preencha o formulário à direita:

### Sala 1: Hall de Entrada
- **ID da Sala:** `1`
- **Nome da Sala:** `Hall de Entrada`
- **Descrição:**
  `Voce esta no grande hall de entrada de uma antiga mansao. Ao norte ha uma sala de estar. A saida ao sul esta trancada.`
- **Saídas Cardeais:**
  - **Norte:** `2` *(Leva à Sala de Estar)*
  - **Sul:** `0` *(Trancada inicialmente)*
  - **Leste:** `0`
  - **Oeste:** `0`
- Pressione o botão **[Salvar Alterações]**.

---

### Sala 2: Sala de Estar
- **ID da Sala:** `2`
- **Nome da Sala:** `Sala de Estar`
- **Descrição:**
  `Voce esta em uma confortavel sala de estar com lareira apagada. Ha passagens para o sul, leste e oeste.`
- **Saídas Cardeais:**
  - **Norte:** `0`
  - **Sul:** `1` *(Retorna ao Hall)*
  - **Leste:** `3` *(Leva à Biblioteca)*
  - **Oeste:** `4` *(Leva ao Porão)*
- Pressione o botão **[Salvar Alterações]**.

---

### Sala 3: Biblioteca
- **ID da Sala:** `3`
- **Nome da Sala:** `Biblioteca`
- **Descrição:**
  `Uma biblioteca silenciosa com estantes repletas de livros empoeirados. A saida fica a oeste.`
- **Saídas Cardeais:**
  - **Norte:** `0`
  - **Sul:** `0`
  - **Leste:** `0`
  - **Oeste:** `2` *(Retorna à Sala de Estar)*
- Pressione o botão **[Salvar Alterações]**.

---

### Sala 4: Porão
- **ID da Sala:** `4`
- **Nome da Sala:** `Porao`
- **Descrição:**
  `Um porao escuro e gelado sob a mansao. A escada para a sala de estar sobe a leste.`
- **Saídas Cardeais:**
  - **Norte:** `0`
  - **Sul:** `0`
  - **Leste:** `2` *(Sobe para a Sala de Estar)*
  - **Oeste:** `0`
- Pressione o botão **[Salvar Alterações]**.

---

## 📦 3. Aba 2: Criando os Objetos e Consistências (Tecla `2`)

Na aba **2. Objetos**, criamos os 5 objetos do jogo, configurando seus nomes com sinônimos (separados por barra `/`), localização inicial e as permissões de interação (Byte de Consistência):

| ID | Nome e Sinônimos | Sala Inicial | Pode Pegar? | Pode Guardar no Recipiente? | Descrição |
|:---:|:---|:---:|:---:|:---:|:---|
| **1** | `LOCAL` | `0` (Fixo) | [ ] Não | [ ] Não | `O local onde voce se encontra.` |
| **2** | `VELA/VELAS` | `1` (Hall) | [x] Sim | [x] Sim | `Uma vela de cera amarelada.` |
| **3** | `MALA/BOLSA` | `2` (Sala) | [x] Sim | [ ] Não | `Uma mala de couro aberta para guardar objetos.` |
| **4** | `CHAVE/CHAVES` | `3` (Biblio) | [x] Sim | [x] Sim | `Uma chave de bronze pesada.` |
| **5** | `FOSFORO/FOSFOROS` | `1` (Hall) | [x] Sim | [x] Sim | `Uma caixinha com fosforos.` |

> **Observações sobre os Objetos:**
> - O **Objeto 1 (`LOCAL`)** é a palavra-chave de sistema usada em comandos como `EXAMINE LOCAL`.
> - O **Objeto 2 (`VELA`)** atua como objeto de iluminação (associado ao registrador de luz).
> - O **Objeto 3 (`MALA`)** atua como o recipiente padrão da engine (armazena até 3 itens).
> - Para cada objeto editado, clique no botão **[Salvar Alterações]**.

---

## 💬 4. Aba 4: Criando as Mensagens do Jogo (Tecla `4`)

Antes de programar os comandos e funções, cadastre as mensagens de texto personalizadas na aba **4. Mensagens**:

| ID | Texto da Mensagem | Uso / Contexto |
|:---:|:---|:---|
| **11** | `A Mansao Misteriosa - Uma aventura em texto para MSX.` | Mensagem de Introdução ao iniciar |
| **30** | `Voce riscou um fosforo e acendeu a vela. Uma luz suave ilumina o ambiente.` | Resposta ao acender a vela |
| **31** | `Voce assoprou e apagou a vela.` | Resposta ao apagar a vela |
| **32** | `Pegou.` | Feedback da ação de pegar |
| **33** | `Soltou no chao.` | Feedback da ação de soltar |
| **34** | `Guardou dentro da mala.` | Feedback ao guardar na mala |
| **35** | `Voce usou a chave e destrancou o portao sul! Parabens, voce venceu!` | Mensagem de vitória do jogo |

---

## ⚙️ 5. Aba 3: Comandos e Funções em Bytecode (Tecla `3`)

Na aba **3. Comandos**, criamos a lógica de jogo dividida em:
1. **Comandos Personalizados** (Combinações específicas de Verbo + Objeto digitados pelo jogador).
2. **Funções do Sistema** (Reset inicial e rotinas padrão de ações como Pegar, Guardar, Soltar e Examinar).

Pressione **`Tab`** para alternar o botão **[Alternar Comandos / Funções]**.

---

### 5.1 Comandos Personalizados

#### Comando 1: Acender a Vela
- **Verbo:** `31` (`FACA` / `ACENDA`)
- **Objeto 1:** `2` (`VELA`)
- **Objeto 2:** `0`
- **Sequência de Instruções Bytecode:**
  1. `TEMOS 2, 3` *(Testa se o jogador possui a vela no inventário; se não tiver, salta para o passo 3)*
  2. `MSG 17` *(Exibe mensagem padrão 17: "Você não tem isso")*
  3. `NVC` *(Interrompe o fluxo e volta ao prompt)*
  4. `LDR 10, 1` *(Passo 3: Define Registrador 10 = 1, indicando luz acesa)*
  5. `MSG 30` *(Exibe mensagem 30: "Você riscou um fósforo e acendeu a vela...")*
  6. `NVC` *(Fim da execução do comando)*

---

#### Comando 2: Destrancar o Portão com a Chave (Condição de Vitória)
- **Verbo:** `7` (`ENTRE`)
- **Objeto 1:** `0`
- **Objeto 2:** `0`
- **Sequência de Instruções Bytecode:**
  1. `LOCAL 1, 3` *(Testa se o jogador está na Sala 1 - Hall; se não estiver, salta para o passo 3)*
  2. `MSG 16` *(Exibe mensagem padrão 16: "Não há como entrar por aqui")*
  3. `NVC`
  4. `TEMOS 4, 6` *(Passo 3: Testa se o jogador carrega a Chave ID 4; se tiver, salta para o passo 6)*
  5. `MSG 15` *(Exibe mensagem padrão 15: "A porta está trancada")*
  6. `NVC`
  7. `MSG 35` *(Passo 6: Exibe mensagem 35 de vitória: "Você usou a chave e destrancou...")*
  8. `FIM` *(Finaliza a partida com encerramento formal do jogo)*

---

### 5.2 Funções de Sistema (Pressione [Alternar Comandos / Funções])

Cadastre as funções padrão da arquitetura do livro:

#### Função 1: Reset do Jogo (Chamada automaticamente no boot)
- **ID da Função:** `1`
- **Instruções:**
  1. `LDR 1, 1` *(Registrador 1 = 1: Sala inicial = Hall de Entrada)*
  2. `LDR 10, 0` *(Registrador 10 = 0: Luz apagada)*
  3. `RET` *(Retorna ao fluxo principal)*

---

#### Função 6: Rotina Padrão de Pegar Objeto (Ativada pelo bit 0 da consistência)
- **ID da Função:** `6`
- **Instruções:**
  1. `AQUI 0, 3` *(Verifica se o objeto em evidência está no local atual; se não estiver, salta para 3)*
  2. `MSG 19` *(Exibe: "Não vejo isso aqui")*
  3. `NVC`
  4. `REG> 8, 4, 7` *(Passo 3: Testa se Registrador 8 [Itens carregados] > 4; se sim, salta para 7)*
  5. `PEGA 0` *(Move o objeto para o inventário do jogador)*
  6. `MSG 32` *(Exibe mensagem 32: "Pegou.")*
  7. `NVC`
  8. `MSG 21` *(Passo 7: Exibe mensagem de erro de carga máxima: "Você não consegue carregar mais nada")*
  9. `NVC`

---

#### Função 7: Rotina Padrão de Guardar Objeto no Recipiente (Bit 1 da consistência)
- **ID da Função:** `7`
- **Instruções:**
  1. `TEMOS 0, 3` *(Testa se o jogador tem o objeto na mão; se sim, salta para 3)*
  2. `MSG 17` *(Exibe: "Você não tem isso")*
  3. `NVC`
  4. `POE 0` *(Passo 3: Move o objeto para dentro da mala / recipiente)*
  5. `MSG 34` *(Exibe mensagem 34: "Guardou dentro da mala.")*
  6. `NVC`

---

#### Função 13: Rotina Padrão de Soltar Objeto
- **ID da Função:** `13`
- **Instruções:**
  1. `TEMOS 0, 3` *(Testa se o jogador carrega o item; se sim, salta para 3)*
  2. `MSG 17` *(Exibe: "Você não tem isso")*
  3. `NVC`
  4. `SOLTA 0` *(Passo 3: Coloca o item no chão da sala atual)*
  5. `MSG 33` *(Exibe mensagem 33: "Soltou no chão.")*
  6. `NVC`

---

#### Função 14: Rotina Padrão de Examinar Objeto
- **ID da Função:** `14`
- **Instruções:**
  1. `AQUI 0, 4` *(Testa se está no chão; se estiver, salta para o passo 4)*
  2. `TEMOS 0, 4` *(Testa se está no inventário; se estiver, salta para o passo 4)*
  3. `MSG 19` *(Não está aqui)*
  4. `NVC`
  5. `MSG 20` *(Passo 4: Exibe mensagem padrão 20: "É apenas...")*
  6. `OBJ 0` *(Imprime a descrição detalhada do objeto)*
  7. `CHRS 46` *(Imprime o caractere ponto final '.')*
  8. `NVC`

---

## 💾 6. Salvando e Gerando a ROM do MSX

1. Pressione **`F2`** em qualquer momento para salvar seu projeto em disco no formato YAML.
2. Mude para a aba **5. Build** (pressione a tecla **`5`**):
   - Pressione o botão **`[1. Compilar para C (game_data.h/c) [F5]]`**:  
     O compilador irá validar todos os ponteiros, saídas cardeais, consistências e montar as tabelas C estáticas e constantes para o Z80.
   - Pressione o botão **`[2. Construir ROM MSX 32K (build.bat) [F9]]`**:  
     O MSXgl e o SDCC irão compilar o código C e montar o arquivo binário `advent.rom`.

---

## 🕹️ 7. Testando no openMSX

Com a ROM gerada, execute no emulador:

```powershell
openmsx -machine C-BIOS_MSX1_EU -cart MSXgl\projects\advent\out\advent.rom
```

A aventura estará 100% pronta para ser jogada no MSX!
Para conferir a solução e o mapa completo, consulte o documento [solucao.md](file:///e:/editadv/solucao.md).
