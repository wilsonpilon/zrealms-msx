# SISTEMA EDITOR DE ADVENTURES

**Versão 3.4 — 1986**  
**Programação:** Renato Degiovani  
**Produção:** PRO KIT Software / PRO Software  
**Agradecimentos:** Cláudio Costa e José Eduardo Neves  
*TODOS OS DIREITOS RESERVADOS*

---

## ESTE MANUAL

Para uma perfeita compreensão da funcionalidade do sistema **EDITOR**, é necessário que se faça uma leitura completa deste manual. Em seguida, o usuário deverá implementar e analisar o exemplo que acompanha este sistema (*MANSÃO*), a fim de que os conceitos de construção de um adventure fiquem mais claros.

A versão para MSX disco do sistema **EDITOR** é compatível com o padrão Microsol de interface de drive.

---

## ÍNDICE

1. **[1. INTRODUÇÃO](#1-introdução)**
   - 1.1 A funcionalidade do jogo
   - 1.2 Procedimento de carga
   - 1.3 Inicialização do EDITOR
   - 1.4 Organização do sistema
   - 1.5 Elementos estruturais do jogo
   - 1.6 A funcionalidade do EDITOR
   - 1.7 Tabela dos comandos do EDITOR
   - 1.8 Edição de alfabetos
2. **[2. VÍDEO E EDITORAÇÃO](#2-vídeo-e-editoração)**
   - 2.1 Campos funcionais
   - 2.2 Editor de tela
   - 2.3 Comandos de edição de tela
3. **[3. REGISTRADORES](#3-registradores)**
   - 3.1 Registradores especiais
   - 3.2 Registradores pré-definidos
4. **[4. VERBOS](#4-verbos)**
   - 4.1 Divisão funcional dos verbos
5. **[5. OBJETOS](#5-objetos)**
   - 5.1 Objetos pré-definidos
   - 5.2 Situação de um objeto
   - 5.3 Byte de consistência
6. **[6. POSIÇÕES](#6-posições)**
   - 6.1 Movimento entre posições
7. **[7. COMANDOS](#7-comandos)**
   - 7.1 Comandos já definidos
8. **[8. INSTRUÇÕES DO SISTEMA](#8-instruções-do-sistema)**
   - 8.1 Edição das instruções
   - 8.2 Tabela de instruções
9. **[9. MENSAGENS](#9-mensagens)**
   - 9.1 Mensagens especiais
   - 9.2 Mensagens pré-definidas
   - 9.3 A estrutura de uma mensagem
10. **[10. FUNÇÕES](#10-funções)**
    - 10.1 Funções especiais
11. **[11. UM EXEMPLO (O JOGO "MANSÃO")](#11-um-exemplo-o-jogo-mansão)**
    - 11.1 Criando um mapa
    - 11.2 A inicialização
    - 11.3 As passagens do jogo
    - 11.4 Criando Objetos
    - 11.5 Projetando a saída
    - 11.6 Conclusão
12. **[12. INFORMAÇÕES COMPLEMENTARES](#12-informações-complementares)**

---

## 1. INTRODUÇÃO

O sistema **EDITOR** foi concebido para monitorar a criação e edição de programas adventures. A sua potencialidade está diretamente vinculada ao conhecimento dos diversos elementos que o compõem.

Todo o trabalho de programação, ou de uso de uma linguagem de computação, foi minimizado a fim de que o autor possa usar seu esforço unicamente na elaboração funcional do jogo. No entanto, as práticas e noções elementares de programação são úteis quando da criação e depuração dos **COMANDOS** e **FUNÇÕES** do sistema.

### 1.1 A Funcionalidade do Jogo

O adventure é um jogo no qual o jogador é sempre o agente criador e idealizador das ações. O computador é passivo e apenas responde a um enredo previamente estruturado pelo autor do jogo.

A forma de interagir com o jogo é através de frases de comando que expressam o desejo do jogador; ou seja, se o jogador deseja pegar um objeto — uma caneta, por exemplo —, ele simplesmente induz o computador a essa ação, na forma:
```
PEGUE A CANETA
```
O computador se restringe a executar a ordem e fornecer um relatório do seu desempenho.

O computador também se encarrega de orientar o jogador em relação ao local onde ele se encontra, descrevendo-o e simulando os movimentos através de comandos do tipo `VÁ PARA O NORTE`, por exemplo.

Os adventures criados pelo **EDITOR** permitem não só uma diversidade de ações, mas também de estruturas de construção das frases. Em função das letras fornecidas pelo jogador, o sistema procura no seu banco um conjunto que satisfaça a intenção da frase. Por exemplo, basta digitar `EX` que o sistema compreende a palavra `EXAMINE`, ou apenas `LOC`, que é interpretado como `LOCAL`.

O sistema possui também formas de simplificação para a abreviação dos comandos mais usados, principalmente as movimentações, com a substituição das frases do tipo `VÁ PARA O NORTE` por um comando mais simples: `NORTE` ou apenas `N`.

A tecla **ENTER**, pressionada sem nenhuma frase para interpretação, assume a função de fazer uma descrição sumária dos objetos que estão no local onde o jogador se encontra.

O jogo permite também a referência indireta ao objeto, ou seja, a construção `PEGUE A CANETA` pode ser estruturada também como `PEGUE-A`, desde que o objeto "caneta" tenha sido referenciado anteriormente.

A versão 3.4 trabalha exclusivamente com textos e produz adventures no estilo de **AMAZÔNIA** e **SERRA PELADA**. O uso da acentuação obedece aos procedimentos normais de cada versão de computador. Uma vez que o modelo Gradiente Expert 1.0 não possui todas as letras acentuadas e nem permite a acentuação correta via teclado, é aconselhável definir as teclas `SHIFT + NÚMERO` como letras acentuadas (veja o item 1.8 sobre edição de alfabetos).

### 1.2 Procedimento de Carga

O procedimento de carga (*LOAD*) do sistema **EDITOR** é bastante simples, porém deve ser executado exatamente como descrito abaixo:
1. Desligue o computador e aguarde alguns segundos para que os capacitores internos descarreguem.
2. Coloque o disco **EDITOR** no drive A.
3. Ligue o micro e feche a portinhola do drive.

### 1.3 Inicialização do EDITOR

O sistema **EDITOR** é carregado automaticamente quando se reseta o computador. A primeira intervenção do operador diz respeito à pergunta:
```
Editar (S/N)?
```
Neste momento o usuário deve optar entre editar um jogo que já existe ou criar um novo jogo.

O disco **EDITOR** possui um jogo exemplo chamado **MANSÃO**, que pode ser acessado neste momento pressionando-se a tecla **"S"** e a tecla **ENTER** a seguir.

Caso se deseje editar um jogo, o sistema solicitará a presença, no drive, de um disco que contenha tal jogo. Caso a tecla pressionada seja **"N"**, o sistema solicitará um nome para o jogo a ser criado e, logo após, solicitará a presença de um disco previamente formatado.

Em ambos os casos o sistema trabalhará com uma organização própria para o disco. Tal organização não responde a comandos do tipo `DIR`, `FILES`, `COPY`, etc.

Quando em edição, é dispensável a presença do disco **EDITOR** no drive, com exceção da ocasião em que for solicitada a carga de um conjunto novo de letras (veja tópico sobre edição de alfabetos). Não é aconselhável, durante a criação de um jogo, usar o disco **EDITOR** como arquivo de trabalho. Tanto o **EDITOR** quanto um jogo por ele criado utilizam apenas 1 drive.

Após a operação de carga ou criação de um jogo, será apresentada ao operador uma tela com um balanço da situação dos elementos das tabelas, bem como o espaço de memória utilizado e o espaço disponível.

### 1.4 Organização do Sistema

O sistema **EDITOR** compõe-se de:
- Um **módulo de jogo** (o player propriamente dito);
- Um **módulo de edição**;
- As **tabelas dos elementos funcionais** do jogo.

Quando em edição, estando no jogo, o retorno ao **EDITOR** é feito mediante o acionamento da tecla **ESC**.

A cópia do jogo gerada pelo **EDITOR** pode ser reeditada apenas se ambos possuírem uma versão compatível. A gravação do jogo final, efetuada pelo comando `CJ`, não carrega consigo o módulo **EDITOR**.

### 1.5 Elementos Estruturais do Jogo

Um jogo criado pelo **EDITOR** é composto por uma série de elementos organizados em tabelas:

* **REGISTRADORES:** São variáveis de 1 byte (0 a 255) cuja função é permitir o registro de estados diversos. Contêm informações sobre as condições de determinados objetos ou estão associados a alguma circunstância relativa a uma posição (por exemplo, se o local onde o jogador se encontra está claro ou escuro). Podem servir também como contadores ou temporizadores e se assemelham, na sua aplicação, às variáveis de um programa em BASIC.
* **VERBOS:** São os verbos reconhecidos pelo sistema e compõem sempre a primeira palavra da frase de comando do jogador. São responsáveis pela definição da ação pretendida durante o jogo.
* **OBJETOS:** São os elementos, ou palavras, que podem ser manipulados ou referenciados pelo jogador (por exemplo: vela, mala, chão, árvore, etc.).
* **POSIÇÕES:** São os locais onde é possível a presença física do jogador no mapa do jogo.
* **COMANDOS:** São séries de procedimentos (instruções de programação) que configuram a sistemática operacional de uma determinada frase do jogador. Por exemplo, para a frase `PEGUE A CANETA`, existe um comando correspondente que permitirá ao sistema efetivar essa ação.
* **MENSAGENS:** São as mensagens emitidas pelo sistema para se comunicar com o jogador, como resultado de uma determinada ação ou descrição de ambiente.
* **FUNÇÕES:** São séries de procedimentos semelhantes aos comandos, cuja requisição é determinada por um acontecimento ou situação especial (passagem de um local para outro, temporizador, início de jogo ou chamada explícita por um comando). Funcionam também como sub-rotinas.

### 1.6 A Funcionalidade do EDITOR

O sistema operacional do **EDITOR** baseia-se no uso direto de comandos, acompanhados ou não de parâmetros, constituídos na forma abreviada da ação:
- Para listar todas as mensagens: `LM` + `ENTER`
- Para listar a mensagem número 5: `LM 5` + `ENTER`

Quando o **EDITOR** está em processo de criação ou edição de um determinado elemento (**MENSAGENS**, **POSIÇÕES** ou **VERBOS**), a tecla **ESC** fica desativada para não cancelar inadvertidamente a operação. Essa restrição não se aplica às **FUNÇÕES** e nem aos **COMANDOS**.

Toda a funcionalidade do sistema do jogo está baseada na indexação dos elementos:
- `OBJETO 3`
- `POSIÇÃO 5`
- `MENSAGEM 10`
- `REGISTRADOR 17`
- `VERBO 9`

Os **COMANDOS** não são indexados por números: são acessados pela sua construção efetiva (por exemplo: `PEGUE A CANETA`).

Uma vez criado um elemento, ele não poderá mais ser deletado do banco de dados do jogo. No caso de não ser mais precisa a sua presença, basta deixá-lo com um valor nulo ou qualquer.

### 1.7 Tabela dos Comandos do EDITOR

Os comandos de listagem funcionam mediante o uso das setas de navegação:
- **SETA PARA BAIXO:** Avança um elemento.
- **SETA PARA CIMA:** Retrocede um elemento.
- **SETA PARA DIREITA:** Edita o elemento em evidência.
- **ESC:** Encerra o comando de listagem e retorna ao prompt do editor.

| Mnemônico | Parâmetros | Função |
|---|---|---|
| **LR** | `[X]` | **Listar Registrador:** Lista o registrador X. Se X for omitido, a listagem será integral (um a um). |
| **MR** | `X , Y` | **Mudar Registrador:** Modifica o valor do registrador X para o valor Y (ambos de 0 a 255). |
| **LV** | `[X]` | **Listar Verbo:** Lista o verbo X. Se X for omitido, lista todos os verbos sequencialmente. |
| **CV** | | **Criar Verbo:** Cria mais um verbo na tabela e atribui-lhe o próximo índice sequencial. |
| **EV** | `X` | **Editar Verbo:** Permite a edição do verbo X. |
| **LP** | `[X]` | **Listar Posição:** Lista a posição X ou todas as posições sequencialmente. |
| **CP** | | **Criar Posição:** Cria uma nova posição geográfica e atribui-lhe o próximo índice. |
| **EP** | `X` | **Editar Posição:** Permite editar a posição X (descrição e saídas Norte/Sul/Leste/Oeste). |
| **CC** | | **Criar Comando:** Cria um novo comando na tabela. A identificação sintática será solicitada a seguir. |
| **EC** | `[frase]` | **Editar Comando:** Permite editar a lista de instruções de um comando existente. |
| **LM** | `[X]` | **Listar Mensagem:** Lista a mensagem X ou todas as mensagens. |
| **CM** | | **Criar Mensagem:** Cria uma nova mensagem e atribui-lhe o próximo índice. |
| **EM** | `X` | **Editar Mensagem:** Permite a edição de texto da mensagem X. |
| **CF** | | **Criar Função:** Cria uma nova função e atribui-lhe o próximo índice. |
| **EF** | `X` | **Editar Função:** Permite a edição das instruções da função X. |
| **LO** | `[X]` | **Listar Objeto:** Lista o objeto X ou todos os objetos cadastrados. |
| **CO** | | **Criar Objeto:** Cria um novo objeto com seu nome, situação e byte de consistência. |
| **EO** | `X` | **Editar Objeto:** Permite editar as propriedades do objeto X. |
| **CJ** | | **Copiar Jogo:** Grava uma cópia de produção do jogo em disco. |
| **TJ** | | **Testar Jogo:** Inicia o teste imediato do jogo em edição desde a sua fase inicial. |
| **RJ** | | **Retornar ao Jogo:** Retorna ao teste no ponto exato onde houve a interrupção via ESC. |
| **POK** | `X , Y` | **POKear:** Modifica o endereço de memória X com o valor Y. |
| **EA** | | **Editar Alfabeto:** Permite a edição do conjunto de fontes ou carga de fontes do disco. |

### 1.8 Edição de Alfabetos

A edição de alfabetos permite utilizar um desenho de caracteres mais adequado ao tema do adventure. Ao ser executado o comando `EA`, é apresentado todo o conjunto de caracteres em uso.

Um quadro em destaque ampliado indica o desenho do caractere sob o cursor:
- **SETAS:** Movimentam o cursor pelo conjunto de caracteres.
- **CTRL + S:** Salva o conjunto de caracteres no disco EDITOR.
- **CTRL + L:** Carrega um conjunto de caracteres do disco EDITOR.
- **CTRL + R:** Anula as alterações efetuadas e recupera o alfabeto original.
- **ESC:** Cancela o comando `EA` e retorna ao editor sem modificações.
- **SELECT:** Aceita as alterações e retorna ao editor.
- **ENTER:** Abre o modo de edição do caractere selecionado.

**Comandos na edição de caractere:**
- **Qualquer tecla:** Copia o desenho correspondente àquela tecla.
- **ESC:** Cancela a edição do caractere sem alterações.
- **CTRL + Q:** Inverte os pixels do caractere.
- **CTRL + L:** Limpa o caractere (apaga todos os pixels).
- **CTRL + I:** Espelha o caractere horizontalmente.
- **SETAS:** Movem o cursor pixel a pixel dentro da matriz 8x8.
- **BARRA DE ESPAÇO:** Inverte o estado do pixel sob o cursor (liga/desliga).
- **ENTER:** Confirma as modificações e retorna à tabela de caracteres.

O disco **EDITOR** possui três buffers de alfabetos pré-definidos. O aplicativo gráfico **GRAPHOS III** também pode fornecer alfabetos e telas para jogos criados pelo EDITOR.

---

## 2. VÍDEO E EDITORAÇÃO

O sistema **EDITOR** baseia-se na utilização do vídeo como processo interativo permanente entre o jogador e a narrativa. A tela é sempre dividida em três campos principais: **superior**, **central** e **inferior**, delimitados por barras divisórias horizontais.

### 2.1 Campos Funcionais

```
Linha  0: ┌────────────────────────────────────────┐  CAMPO SUPERIOR:
Linha  1: ├────────────────────────────────────────┤  Eco da frase comando interpretada
Linhas 2  │ EXAMINE O LOCAL                        │
  a   20: │                                        │  CAMPO CENTRAL:
          │ Neste local tem:                       │  Descrição do ambiente, saídas,
          │ uma vela                               │  respostas e mensagens emitidas
          │ uma mala                               │
Linha 21: ├────────────────────────────────────────┤
Linha 22: │ > PEGUE A VELA                         │  CAMPO INFERIOR:
Linha 23: └────────────────────────────────────────┘  Prompt de entrada de comandos
```

* **CAMPO SUPERIOR:** Apresenta a frase do jogador após ter sido interpretada pelo parser. Todas as palavras irrelevantes (artigos, preposições não cadastradas) são eliminadas, exibindo a forma canônica da ação.
* **CAMPO CENTRAL:** Apresenta o resultado das ações, descrições de salas e mensagens.
* **CAMPO INFERIOR:** Recebe a digitação do jogador em tempo real. O cursor pode se deslocar até a última coluna da linha.

Durante a **edição**, os campos assumem funções específicas:
* **CAMPO SUPERIOR:** Apresenta o nome e a versão do sistema.
* **CAMPO CENTRAL:** Apresenta o balanço de memória, listagens e permite a edição direta de textos, comandos e funções.
* **CAMPO INFERIOR:** Linha de comando do editor para entrada dos mnemônicos (`LM`, `EP`, etc.).

### 2.2 Editor de Tela

Todas as operações de criação de mensagens ou descrição de salas utilizam o editor de tela embutido. O sistema apresenta um cursor não-destrutivo e a posição inicial padrão de impressão indicada pelo caractere `)`.

O texto digitado após o caractere `)` será incorporado ao jogo exatamente na posição indicada.

### 2.3 Comandos de Edição de Tela

| Tecla / Comando | Função |
|---|---|
| **INS** | Insere um espaço na posição do cursor, deslocando o texto à direita. |
| **DEL** | Deleta o caractere na posição do cursor, puxando o texto à direita. |
| **ENTER** | Finaliza a edição e compila a mensagem para a tabela correspondente. |
| **HOME** | Coloca o cursor no canto superior esquerdo do campo central. |
| **CAPS** | Alterna entre maiúsculas e minúsculas. |
| **CTRL + D** | Insere uma quebra de linha forçada (*Line Feed*). |
| **CTRL + I** | Inverte o vídeo do caractere. |
| **CTRL + J** | Repõe o caractere `)` indicativo da posição inicial padrão de impressão. |
| **CTRL + O** | Insere uma chamada de mensagem especial (`#1` a `#9`). |

---

## 3. REGISTRADORES

Os registradores são variáveis numéricas de 1 byte (faixa de `0` a `255`). Existem 99 registradores no sistema (`1` a `99`), sendo que os 14 primeiros possuem atribuições pré-definidas pela engine.

### 3.1 Registradores Especiais (Acesso aos Objetos)

Durante a partida, a condição e localização de cada objeto do jogo é mantida em uma tabela de bytes internos. O acesso direto ao estado de um objeto pode ser feito tratando-o como um registrador através da fórmula:
$$	ext{Índice} = 	ext{Objeto} + 100$$

*Exemplo:* Para verificar ou alterar o estado do **Objeto 7**, referencia-se o **Registrador 107**.

### 3.2 Registradores Pré-definidos

| Registrador | Nome / Função | Descrição |
|---|---|---|
| **Reg 1** | Posição Corrente | Número da sala/posição onde o jogador se encontra atualmente. |
| **Reg 2** | Contador de Jogadas (Low) | Byte menos significativo do total de comandos/ENTERs executados. |
| **Reg 3** | Contador de Jogadas (High) | Byte mais significativo do contador de jogadas. |
| **Reg 4** | Contador Auxiliar | Incrementado a cada frase de comando sempre que for diferente de zero. |
| **Reg 5** | Temporizador (Bomba-relógio) | Decrementado a cada frase de comando se for $> 0$. Ao atingir zero, dispara automaticamente a **Função 2**. |
| **Reg 6** | Passos no Escuro | Incrementado a cada comando dado em local escuro. Ao atingir o valor 5, dispara a **Função 3** (morte no escuro). |
| **Reg 7** | Objetos no Recipiente | Quantidade de objetos guardados dentro do Objeto 3 (normalmente limite de 3). |
| **Reg 8** | Inventário do Jogador | Quantidade de objetos sendo carregados pelo jogador (normalmente limite de 5). |
| **Reg 9** | Iluminação do Local | `0` = Local claro; `1` = Local escuro. |
| **Reg 10** | Estado da Fonte de Luz | Estado do Objeto 2: `0` = Apagado; `1` = Aceso. |
| **Reg 11** | Minutos do Relógio | Contador de minutos virtuais do jogo. |
| **Reg 12** | Horas do Relógio | Contador de horas virtuais do jogo. |
| **Reg 13** | Dias do Relógio | Contador de dias virtuais do jogo. |
| **Reg 14** | Fator de Cadência | Fator de ajuste de tempo real para o relógio do sistema. |

---

## 4. VERBOS

Os verbos compõem a primeira palavra da frase de comando e definem a ação pretendida. Podem ser criados até 209 verbos ou conjuntos de sinônimos.

Ao criar sinônimos para um verbo, utiliza-se a barra `/` como delimitador.  
*Exemplo:* `PEGUE/TOME/APANHE/AGARRE`

### 4.1 Divisão Funcional dos Verbos (Tabela Obrigatória 1 a 34)

Os primeiros 34 verbos possuem distribuição fixa reservada pela engine:

| ID | Verbo Canônico | Ação Padrão Associada |
|---|---|---|
| **1** | `NORTE` | Movimento cardeal Norte |
| **2** | `SUL` | Movimento cardeal Sul |
| **3** | `LESTE` | Movimento cardeal Leste |
| **4** | `OESTE` | Movimento cardeal Oeste |
| **5** | `GRAVE` | Salva o estado do jogo |
| **6** | `RECUPERE` | Restaura partida salva |
| **7** | `ENTRE` | Entrar em local/objeto |
| **8** | `SUBA` | Movimento para cima / subir |
| **9** | `SAIA` | Sair de local/objeto |
| **10** | `DESCA` | Movimento para baixo / descer |
| **11** | `HORAS` | Consulta o relógio do jogo |
| **12** | `QUANTO` | Consulta pontuação ou contadores |
| **13** | `TEMOS` | Inventário de objetos carregados |
| **14** | `RECOMECE` | Reinicia a partida |
| **15** | `GARIMPE` | Ação de mineração/busca |
| **16** | `PENSE` | Reflexão ou dica |
| **17** | `GRITE` | Emitir som alto |
| **18** | `CORRA` | Movimento rápido |
| **19** | `PEGUE` | Apanhar objeto (Função 6) |
| **20** | `COLOQUE` | Guardar objeto no Objeto 3 (Função 7) |
| **21** | `TROQUE` | Trocar objetos (Função 8) |
| **22** | `COMPRE` | Comprar item (Função 9) |
| **23** | `ROUBE` | Subtrair item (Função 10) |
| **24** | `TIRE` | Retirar de recipiente (Função 11) |
| **25** | `QUEBRE` | Danificar/destruir item (Função 12) |
| **26** | `SOLTE` | Largar objeto no chão (Função 13) |
| **27** | `EXAMINE` | Inspecionar item ou sala (Função 14) |
| **28** | `PROCURE` | Buscar itens ocultos (Função 15) |
| **29** | `OFERECA` | Dar presente/item (Função 16) |
| **30** | `FACA` | Construir/manufaturar (Função 17) |
| **31** | `JOGUE` | Arremessar objeto (Função 18) |
| **32** | `CONSERTE` | Reparar item (Função 19) |
| **33** | `VENDA` | Vender item comercial (Função 20) |
| **34** | `ACENDA` / Livre | Ação personalizada de iluminação ou livre |

---

## 5. OBJETOS

Todas as palavras substantivas passíveis de referência pelo jogador devem ser cadastradas como **OBJETOS**, mesmo que sejam apenas elementos cenográficos ou decorativos. O sistema suporta até 99 objetos (indexados de 1 a 99).

### 5.1 Objetos Pré-definidos

* **Objeto 1 (`LOCAL`):** Palavra reservada para permitir a inspeção do ambiente (`EXAMINE O LOCAL`).
* **Objeto 2 (Fonte de Luz):** Objeto reservado para iluminação (ex: lampião, tocha ou vela). O controle de luminosidade e posse é automatizado pelo sistema.
* **Objeto 3 (Recipiente):** Objeto reservado com capacidade de conter outros itens dentro de si (ex: saco, mochila, baú ou mala).

### 5.2 Situação de um Objeto

O estado de localização de um objeto é determinado pelos seguintes códigos numéricos:

| Valor | Situação do Objeto |
|---|---|
| **0** | Objeto inexistente ou apenas palavra decorativa sem presença física. |
| **1 a 99** | Objeto presente no chão da posição/sala correspondente a este número. |
| **250** | Objeto sendo carregado pelo jogador (no inventário). |
| **251** | Objeto guardado dentro do Objeto 3 quando este está aberto. |
| **253** | Objeto guardado dentro do Objeto 3 quando este está fechado. |

### 5.3 Byte de Consistência

Cada objeto possui um byte de propriedades cujos bits indicam predisposição para ações padronizadas sem necessidade de criar comandos manuais:

| Bit | Valor | Ação Automática | Função Chamada |
|---|---|---|---|
| **Bit 0** | 1 | Objeto pode ser pego livremente (`PEGUE`) | **Função 6** |
| **Bit 1** | 1 | Pode ser colocado no Objeto 3 (`COLOQUE`) | **Função 7** |
| **Bit 2** | 1 | Pode ser trocado (`TROQUE`) | **Função 8** |
| **Bit 3** | 1 | Pode ser comprado (`COMPRE`) | **Função 9** |
| **Bit 4** | 1 | Pode ser roubado (`ROUBE`) | **Função 10** |
| **Bit 5** | 1 | Pode ser retirado do Objeto 3 (`TIRE`) | **Função 11** |
| **Bit 6** | 1 | Pode ser quebrado (`QUEBRE`) | **Função 12** |

Se um bit for `0`, a ação padrão é bloqueada para aquele objeto, a menos que o autor crie um **COMANDO** explícito para tratar a situação.

---

## 6. POSIÇÕES

Posições são as salas ou locais do mapa onde o jogador pode estar. O sistema suporta até 99 posições (numeradas de 1 a 99).

Cada posição possui:
1. Um **texto descritivo** exibido ao entrar ou ao pressionar ENTER;
2. Uma matriz com as **quatro conexões cardeais**: `Norte / Sul / Leste / Oeste`.

### 6.1 Movimento Entre Posições

- Se a saída apontar para um valor de `1 a 99`, o jogador é movido imediatamente para essa posição e ela é descrita.
- Se a saída for `0`, o movimento é bloqueado e o sistema emite a **Mensagem 15** (*"É impossível ir nesta direção."*).
- Se a saída for um valor **superior a 100**, trata-se de um **Movimento Condicional**. O sistema desvia a execução para a **Função (`Valor - 100`)**, permitindo verificar portas trancadas, armadilhas ou desmoronamentos antes de autorizar a passagem.

---

## 7. COMANDOS

Os comandos definem as ações específicas do jogo e possuem prioridade máxima de execução sobre os comportamentos padrão.

A estrutura sintática de um comando é formada por:
$$	ext{VERBO} + 	ext{OBJETO 1} + 	ext{OBJETO 2}$$

*Exemplo:* `COLOQUE A VELA NA MALA` $ightarrow$ `(COLOQUE, VELA, MALA)`.

### 7.1 Comandos Pré-definidos

A engine já inclui nativamente três comandos operacionais:
* `EXAMINE O LOCAL`: Descreve a sala e lista todos os objetos presentes no chão.
* `TEMOS`: Lista todos os objetos no inventário do jogador.
* `REINICIE`: Reinicia a partida, restaurando todas as variáveis e posições iniciais.

---

## 8. INSTRUÇÕES DO SISTEMA (BYTECODES)

Os comandos e as funções são programados através de uma linguagem de instruções estruturadas (45 opcodes no total).

### 8.1 Edição das Instruções

No editor, o cursor `)` aponta para a instrução em evidência:
- **SETAS:** Navegam para cima e para baixo pelas instruções da listagem.
- **INS:** Insere uma nova linha de instrução.
- **DEL:** Remove a instrução sob o cursor.
- **SELECT:** Compila a listagem.
- **Labels (Rótulos de Salto):** Letras de `A` a `T` para desvios condicionais.

### 8.2 Tabela Completa de Instruções

| Mnemônico | Parâmetros | Descrição da Ação |
|---|---|---|
| **NOP** | | Sem efeito (No Operation). |
| **MSG** | `X` | Imprime o texto da Mensagem X no campo central. |
| **NVC** | | Aguarda novo comando do jogador (encerra o processamento do turno). |
| **LLIST** | | Lista os objetos presentes no local e aguarda novo comando. |
| **CLIST** | | Lista os objetos carregados pelo jogador e aguarda novo comando. |
| **DLIST** | | Lista os objetos guardados dentro do Objeto 3 e aguarda novo comando. |
| **OBJ** | `X` | Imprime o nome do Objeto X (se X=0, usa o objeto em evidência). |
| **INC** | `X` | Incrementa em 1 o valor do Registrador X. |
| **DEC** | `X` | Decrementa em 1 o valor do Registrador X. |
| **REG=** | `X , Y` | Atribui o valor Y ao Registrador X ($X \leftarrow Y$). |
| **SOMA** | `X , Y` | Soma Y ao Registrador X ($X \leftarrow X + Y$). |
| **RAND** | `X , Y` | Gera um número aleatório entre 0 e Y e soma ao Registrador X. |
| **REG==** | `X , Y , L` | Se Registrador $X = Y$, salta para o Label L; senão continua. |
| **REG>** | `X , Y , L` | Se Registrador $X > Y$, salta para o Label L; senão continua. |
| **REG<** | `X , Y , L` | Se Registrador $X < Y$, salta para o Label L; senão continua. |
| **AQUI** | `X , L` | Se o Objeto X estiver no local corrente, salta para o Label L. |
| **LOCAL** | `X , L` | Se a posição atual do jogador for X, salta para o Label L. |
| **TEMOS** | `X , L` | Se o Objeto X estiver sendo carregado, salta para o Label L. |
| **SOLTA** | `X` | Solta o Objeto X no chão (se X=100, solta todos os carregados). |
| **PEGA** | `X` | Coloca o Objeto X no inventário do jogador. |
| **COL** | `X` | Coloca o Objeto X na posição atual do jogador. |
| **APAGA** | `X` | Apaga o Objeto X do jogo (define sua situação como 0). |
| **GOSUB** | `X` | Executa a Função X como sub-rotina e retorna no RET. |
| **LIBR** | | Destranca e abre o Objeto 3 (situação dos itens vira 251). |
| **TRC** | | Tranca o Objeto 3 (situação dos itens vira 253). |
| **POE** | `X` | Coloca o Objeto X dentro do Objeto 3. |
| **ESV** | | Esvazia o Objeto 3, jogando todo o seu conteúdo no chão. |
| **OK** | | Imprime `"Ok."` e aguarda novo comando (equivale a `MSG + NVC`). |
| **REGN** | `X` | Imprime o valor numérico decimal do Registrador X. |
| **NVF** | | Aguarda novo comando sem testar a Função 5. |
| **REF** | | Retorna da execução da Função 4 (*Real Time*). |
| **FIM** | | Encerra a partida com vitória ou fim de jogo. |
| **NEU** | | Reinicia completamente a partida (Reset geral). |
| **DESC** | | Redescreve o local corrente e aguarda novo comando. |
| **RET** | | Retorna de uma chamada `GOSUB`. |
| **GOTO** | `L` | Salto incondicional para o Label L. |
| **PAUSA** | `X` | Pausa a execução durante X segundos. |
| **FLAG** | `X , Y` | Atribui o valor Y à Flag interna X. |
| **EVID** | `X` | Coloca em evidência o Objeto 1 ou Objeto 2 da frase do jogador. |
| **CLS** | | Limpa o campo central do vídeo. |
| **EVD=** | `X , L` | Se o objeto em evidência for igual a X, salta para o Label L. |
| **CHR$** | `X` | Imprime diretamente o caractere de código ASCII X. |
| **DNT** | `X , L` | Se o Objeto X estiver dentro do Objeto 3, salta para o Label L. |
| **CMD** | `X` | Redireciona a execução para o Comando X. |

---

## 9. MENSAGENS

As mensagens são o canal principal de comunicação e atmosfera com o jogador. São indexadas de 1 a 209.

### 9.1 Mensagens Especiais (Otimização de Espaço 1 a 9)

As mensagens 1 a 9 podem ser inseridas no início de outras mensagens ou descrições através do caractere especial `#N`, economizando espaço de memória para prefixos repetitivos:
- `#1` $ightarrow$ `"Sinto muito, mas "`
- `#2` $ightarrow$ `"Nós estamos "`

### 9.2 Mensagens Pré-definidas do Sistema

| ID | Uso Padrão pela Engine | Texto Típico Recomendado |
|---|---|---|
| **10** | Apresentação / Prólogo | Introdução do enredo exibida no início da partida. |
| **11** | Sucesso em busca | `"Achei o que você queria."` |
| **12** | Local escuro sem luz | `"Está muito escuro aqui. É melhor arranjar alguma luz..."` |
| **13** | Comando não compreendido | `"Perdão, não entendi..."` |
| **14** | Direção bloqueada | `"É impossível ir nesta direção."` |
| **15** | Ação impossível / erro | `"Isto não é possível."` |
| **16** | Objeto não possuído | `"Nós não temos isso."` |
| **17** | Já possui o objeto | `"Nós já temos isso."` |
| **18** | Objeto ausente no local | `"Eu não estou vendo isso aqui."` |
| **19** | Exame de objeto comum | `"É apenas um objeto comum."` |
| **20** | Limite de inventário cheio | `"Não dá para carregar mais nada."` |
| **21** | Recipiente cheio | `"Não cabe mais nada dentro."` |

---

## 10. FUNÇÕES

Funções são rotinas acionadas por eventos, passagens de mapa, temporizadores ou chamadas explícitas via `GOSUB`.

### 10.1 Funções Especiais Pré-definidas

* **Função 1:** Rotina de inicialização e reset da partida (executada no início do jogo).
* **Função 2:** Disparada automaticamente quando o temporizador (**Reg 5**) zera (bomba-relógio ou sede).
* **Função 3:** Disparada após 5 comandos consecutivos no escuro (**Reg 6** $\ge 5$, morte no escuro).
* **Função 4:** Rotina de tempo real (*Real Time*), executada ciclicamente se o primeiro opcode não for NOP.
* **Função 5:** Rotina de pré-comando, executada a cada turno antes de processar a frase do jogador.
* **Funções 6 a 20:** Implementações padrão para os verbos do byte de consistência:
  - **Função 6:** Pegar objeto comum.
  - **Função 7:** Colocar objeto no Objeto 3.
  - **Função 8:** Trocar objeto.
  - **Função 9:** Comprar objeto.
  - **Função 10:** Roubar objeto.
  - **Função 11:** Tirar objeto do Objeto 3.
  - **Função 12:** Quebrar objeto.
  - **Função 13:** Soltar objeto no chão.
  - **Função 14:** Examinar objeto.
  - **Função 15:** Procurar objeto.
  - **Função 16:** Dar/oferecer objeto.
  - **Função 17:** Fabricar/fazer objeto.
  - **Função 18:** Jogar/arremessar objeto.
  - **Função 19:** Consertar objeto.
  - **Função 20:** Vender objeto.

---

## 11. UM EXEMPLO PRÁTICO: O JOGO "MANSÃO"

Para consolidar os conceitos expostos, este capítulo descreve a construção passo a passo do jogo exemplo **MANSÃO**.

### 11.1 O Mapa da Mansão

O jogo se passa em uma mansão misteriosa composta por 13 posições:

```
                  [ 11. Sótão ]
                       │
                       │ (alçapão)
                       │
[ 10. Quarto ] ── [ 9. Corredor ] ── [ 8. Banheiro ]
                       │
                       │ (escada)
                       │
 [ 4. Sala ]   ── [ 1. Hall ]    ── [ 2. Copa ] ── [ 3. Cozinha ]
                       │                               │
                  [ 12. Entrada ]                 [ 5. Adega ]
                       │
                 [ 13. Portão ] (Saída / Vitória)
```

**Tabela de Saídas das Posições (Norte / Sul / Leste / Oeste):**
1. **Hall:** `9 / 12 / 2 / 4`
2. **Copa:** `0 / 0 / 3 / 1`
3. **Cozinha:** `0 / 105 / 0 / 2` *(Sul = 105: Função 5 condicional para descer à adega)*
4. **Sala de Estar:** `0 / 0 / 1 / 0`
5. **Adega:** `103 / 0 / 0 / 0` *(Norte = 103: Função 3 para subir à cozinha)*
8. **Banheiro:** `0 / 0 / 0 / 9`
9. **Corredor:** `111 / 1 / 8 / 10` *(Norte = 111: Função 11 para subir ao sótão)*
10. **Quarto:** `0 / 0 / 9 / 0`
11. **Sótão:** `0 / 109 / 0 / 0` *(Sul = 109: Função 9 para descer pelo alçapão)*
12. **Entrada da Mansão:** `1 / 113 / 0 / 0` *(Sul = 113: Porta da frente trancada)*
13. **Portão Externo:** Local de vitória e término da partida.

### 11.2 A Inicialização (Função 1)

A Função 1 configura a partida:
- Posiciona o jogador no **Hall** (`Reg 1 = 1`);
- Define as lâmpadas apagadas (`Reg 10 = 0`);
- Inicializa a mensagem de apresentação (Mensagem 10).

### 11.3 Objetos e Quebra-cabeças

* **A Vela e os Fósforos:**
  - A vela (Objeto 2) permite explorar a Adega (Posição 5, escura).
  - Enquanto acesa, um contador consome o pavio; se queimar por completo, apaga-se.
* **O Rádio e as Pilhas:**
  - O rádio necessita das pilhas colocadas em seu interior para sintonizar a rádio FM e fornecer as horas do relógio.
* **A Geladeira (Objeto 3):**
  - Funciona como recipiente para armazenar a comida e itens gelados.
* **A Chave da Frente:**
  - Oculta no sótão, necessária para destrancar a porta da posição 12 e alcançar a posição 13 (vitória).

---

## 12. INFORMAÇÕES COMPLEMENTARES

1. **Loop Principal do Sistema:**
   - 1. Executa a **Função 5** (pré-comando).
   - 2. Lê a frase de comando do jogador no campo inferior.
   - 3. Incrementa o contador de jogadas (**Reg 4**).
   - 4. Decrementa o temporizador (**Reg 5**) e testa se chegou a zero.
   - 5. Analisa a sintaxe da frase através do parser.
   - 6. Dispara o **Comando** específico ou a **Ação Padrão**.
   - 7. Retorna ao passo 1.
2. **Endereços Notáveis em Memória:**
   - `0x2693` (`0x6593`): Quantidade máxima de objetos no inventário pessoal.
   - `0x2694` (`0x6594`): Quantidade máxima de objetos suportados no recipiente (Objeto 3).
3. **Compatibilidade de Teclado no MSX:**
   - Nas máquinas da linha Gradiente Expert 1.0, o conjunto estendido de caracteres acentuados é mapeado para as combinações de `SHIFT` com os numerais.
