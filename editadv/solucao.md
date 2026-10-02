# Solução Completa e Mapa do Jogo de Demonstração
## *A Mansão Misteriosa*

Guia de teste, mapa de conexões e passo a passo de resolução da aventura de exemplo (`games/demo.yaml`) para teste no MSX e no openMSX.

---

## 🗺️ 1. Mapa das Salas (Conexões Geográficas)

```
                     +--------------------+
                     |    [3] BIBLIOTECA   |
                     |   (Chave de bronze)|
                     +---------+----------+
                               |
                               | (Oeste / Leste)
                               |
+--------------------+   (Leste)+--------------------+
|     [4] PORÃO      |<-------->| [2] SALA DE ESTAR  |
|  (Local escuro)    |   (Oeste)|       (Mala)       |
+--------------------+          +---------+----------+
                                          |
                                          | (Sul / Norte)
                                          v
                                +--------------------+
                                | [1] HALL DE ENTRADA| <--- Início do Jogo
                                |  (Vela e Fósforo)  |
                                +---------+----------+
                                          |
                                          | [PORTÃO SUL - TRANCA]
                                          v (ENTRE c/ Chave)
                                    [ VITÓRIA! ]
```

### Detalhes das Salas:
| ID | Sala | Descrição | Saídas Disponíveis | Itens no Início |
|---|---|---|---|---|
| **1** | **Hall de Entrada** | Grande hall da mansão. Saída ao sul trancada. | **N** $\rightarrow$ Sala de Estar (2) | `VELA` (ID 2), `FOSFORO` (ID 5) |
| **2** | **Sala de Estar** | Sala de estar com lareira apagada. | **S** $\rightarrow$ Hall (1), **L** $\rightarrow$ Biblioteca (3), **O** $\rightarrow$ Porão (4) | `MALA` (ID 3) |
| **3** | **Biblioteca** | Sala silenciosa com estantes empoeiradas. | **O** $\rightarrow$ Sala de Estar (2) | `CHAVE` (ID 4) |
| **4** | **Porão** | Porão escuro e gelado sob a mansão. | **L** $\rightarrow$ Sala de Estar (2) | Nenhum |

---

## 🎒 2. Objetos do Jogo

| ID | Nome / Sinônimo | Início | Propriedades | Função |
|---|---|---|---|---|
| **1** | `LOCAL` | — | Fixo (não pegável) | Referência ao ambiente |
| **2** | `VELA` / `VELAS` | Sala 1 | Pegável, Guardável | Pode ser acesa com fósforos |
| **3** | `MALA` / `BOLSA` | Sala 2 | Pegável (Recipiente) | Armazena até 3 objetos dentro dela |
| **4** | `CHAVE` / `CHAVES` | Sala 3 | Pegável, Guardável | Necessária para abrir o portão no Hall (1) |
| **5** | `FOSFORO` / `FOSFOROS` | Sala 1 | Pegável, Guardável | Usado para acender a vela |

---

## 🏆 3. Passo a Passo Mais Rápido para Vencer (Speedrun / Teste Rápido)

Para testar rapidamente a condição de vitória e a mecânica de comandos:

1. **Início** no *Hall de Entrada* (Sala 1).
2. Digite: `N` *(Move-se para a Sala de Estar)*
3. Digite: `L` *(Move-se para a Biblioteca)*
4. Digite: `PEGUE CHAVE` *(Coleta a chave de bronze)*
5. Digite: `O` *(Retorna à Sala de Estar)*
6. Digite: `S` *(Retorna ao Hall de Entrada)*
7. Digite: `ENTRE` *(Usa a chave para destrancar o portão)*

🎉 **Resultado:**
> *"Voce usou a chave e destrancou o portao sul! Parabens, voce venceu!"*  
> O jogo finaliza a partida com sucesso.

---

## 🧪 4. Roteiro de Teste Completo (Mecânicas da Engine)

Para validar todas as funcionalidades da engine MSX (inventário, recipientes, iluminação, parser e exibição):

### Teste A: Coleta e Exame de Objetos
```text
> PEGUE VELA
Pegou.

> PEGUE FOSFORO
Pegou.

> EXAMINE VELA
E apenas uma vela.

> INVENTARIO (ou TEMOS)
(Lista os objetos que você está carregando)
```

### Teste B: Ação Especial (Acender a Vela)
```text
> ACENDA VELA (ou FACA VELA)
Voce riscou um fosforo e acendeu a vela. Uma luz suave ilumina o ambiente.
```

### Teste C: Recipiente (Mala / Objeto 3)
```text
> N
(Chega à Sala de Estar)

> PEGUE MALA
Pegou.

> GUARDE VELA
Guardou dentro da mala.

> TEMOS
(Exibe que a vela está guardada dentro da mala)
```

### Teste D: Exploração do Porão e Biblioteca
```text
> O
(Entra no Porão)

> L
(Volta à Sala de Estar)

> L
(Entra na Biblioteca)

> PEGUE CHAVE
Pegou.

> O
(Volta à Sala de Estar)

> S
(Volta ao Hall de Entrada)

> ENTRE
Voce usou a chave e destrancou o portao sul! Parabens, voce venceu!
```

---

## ⌨️ 5. Resumo de Comandos Aceitos

- **Direções:** `N` (Norte), `S` (Sul), `L` (Leste), `O` (Oeste).
- **Ações:** `PEGUE <obj>`, `SOLTE <obj>`, `GUARDE <obj>`, `EXAMINE <obj>`, `ACENDA <obj>`, `ENTRE`.
- **Informações:** `TEMOS` (ou `INVENTARIO`), `ENTER` (repete descrição do local).

---

## 🔤 6. Acentuação e Teclas de Atalho (Suporte Completo ao Português)

A engine utiliza o banco gráfico customizado do adventure (`vram.dat`), carregado na VRAM da SCREEN 0 (`0x0800..0x0FFF`). Para garantir compatibilidade com todos os teclados MSX e emuladores:

### Teclas de Atalho com SHIFT (Padrão Original do Editor de Adventures):
Conforme documentado no manual original (Item 1-2 e 1-8), para teclados sem caracteres acentuados nativos (como o Gradiente Expert 1.0 ou MSX internacional):
- **`SHIFT + 0`**: Digita **`Ç`**
- **`SHIFT + 1`**: Digita **`Á`**
- **`SHIFT + 2`**: Digita **`É`**
- **`SHIFT + 3`**: Digita **`Í`**
- **`SHIFT + 4`**: Digita **`Ó`**
- **`SHIFT + 5`**: Digita **`Ú`**
- **`SHIFT + 6`**: Digita **`Ã`**
- **`SHIFT + 7`**: Digita **`Õ`**
- **`SHIFT + 8`**: Digita **`Â`**
- **`SHIFT + 9`**: Digita **`Ô`**

### Outras Formas de Digitar o Ç e Acentos:
- **`CTRL + C`**: Digita diretamente **`Ç`**.
- **Teclas Mortas (Dead Keys para PC / emulador openMSX):**
  - Digite `'` seguido de `C` para obter **`Ç`**.
  - Digite `'` seguido de vogal para obter **`Á`**, **`É`**, **`Í`**, **`Ó`**, **`Ú`**.
  - Digite `~` seguido de `A` ou `O` para obter **`Ã`**, **`Õ`**.
  - Digite `^` seguido de vogal para obter **`Â`**, **`Ê`**, **`Ô`**.
- **Teclados Nacionais MSX (Sharp HotBit / Expert com BIOS BR):**
  - A tecla dedicada `Ç` e os caracteres acentuados nativos da BIOS são reconhecidos diretamente.

### Apresentação Adequada na Tela (Decodificação UTF-8):
- Todas as mensagens do sistema (como *"É impossível ir nesta direção."*, *"Perdão, não entendi..."*, etc.) e textos do autor são decodificados dinamicamente para os códigos de caractere únicos da fonte do MSX.
- Cada acento é exibido como **1 único caractere correto** desenhado na tela, eliminando caracteres estranhos duplos.
- O parser do jogo reconhece comandos digitados com ou sem acentos (por exemplo, aceita tanto `DESCA` quanto `DESÇA`, `FACA` quanto `FAÇA`).

