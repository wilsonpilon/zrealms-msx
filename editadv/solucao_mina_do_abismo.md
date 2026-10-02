# Solução Completa e Mapa do Jogo
## *A Mina do Abismo*

Guia de teste, mapa de conexões geográficas, tabela de objetos, variáveis e passo a passo de resolução do adventure **"A Mina do Abismo"** (`games/mina_do_abismo.yaml`).

---

## 🗺️ 1. Mapa Completo das Salas (Conexões Geográficas)

```
                     +---------------------------+
                     |    [4] LAGO SUBTERRÂNEO   |
                     |    (Fonte de Água Pura)   |
                     +-------------+-------------+
                                   |
                                   | (Norte / Sul)
                                   v
+--------------------+       +-----+---------------------+       +---------------------------+
|                    |       |  [2] GALERIA DE MINERAÇÃO |       | [3] ALMOXARIFADO          |
|                    |<----->|  (Pederneira)             |<----->|     (Mochila de Lona)     |
|                    | (Oeste|                           | (Leste+-------------+-------------+
|                    | /Leste+-------------+-------------+ /Oeste)             |
| [1] ENTRADA        |                     |                                   | (Sul / Norte)
|     DA MINA        |<--------------------+                                   v
| (Vela e Cantil)    |                                           +-------------+-------------+
| [INÍCIO DO JOGO]   |                                           | [5] ESCRITÓRIO DO FEITOR  |
+---------+----------+                                           |     (Baú de Ferro Trancado|
          |                                                      |      -> Contém a CORDA)   |
          | (Sul / Norte)                                        +---------------------------+
          v
+---------+----------+       +---------------------------+       +---------------------------+
| [6] FENDA DOS      |       |  [7] TÚNEL DOS TRILHOS    |       |  [8] ABISMO DO TÚNEL      |
|     CRISTAIS       |<----->|      (Antigos Trilhos)    |<----->|      (Pilar de Pedra /    |
| (Chave de Ferro)   | (Leste|                           | (Leste|       Vento da Saída)     |
| [SALA ESCURA!]     | /Oeste)      [SALA ESCURA!]       | /Oeste)     [AMARRAR CORDA E      |
+--------------------+       +---------------------------+       |      DESCER = VITÓRIA!]   |
                                                                 +---------------------------+
```

### Detalhes das 8 Salas:
| ID | Sala | Iluminação | Conexões Disponíveis | Itens Iniciais |
|---|---|---|---|---|
| **1** | **Entrada da Mina** | Clara | **L** $\rightarrow$ Galeria (2), **S** $\rightarrow$ Fenda (6) | `VELA` (ID 2), `CANTIL` (ID 6) |
| **2** | **Galeria de Mineração** | Clara | **O** $\rightarrow$ Entrada (1), **N** $\rightarrow$ Lago (4), **L** $\rightarrow$ Almoxarifado (3) | `PEDERNEIRA` (ID 5) |
| **3** | **Almoxarifado** | Clara | **O** $\rightarrow$ Galeria (2), **S** $\rightarrow$ Escritório (5) | `MOCHILA` (ID 3, Recipiente) |
| **4** | **Lago Subterrâneo** | Clara | **S** $\rightarrow$ Galeria (2) | `AGUA` (ID 8, Fonte pura) |
| **5** | **Escritório do Feitor** | Clara | **N** $\rightarrow$ Almoxarifado (3) | `BAU` (ID 9, Baú trancado) |
| **6** | **Fenda dos Cristais** | **Escura!** | **N** $\rightarrow$ Entrada (1), **L** $\rightarrow$ Trilhos (7) | `CHAVE` (ID 4, Chave de ferro) |
| **7** | **Túnel dos Trilhos** | **Escura!** | **O** $\rightarrow$ Fenda (6), **L** $\rightarrow$ Abismo (8) | Nenhum |
| **8** | **Abismo do Túnel** | Clara | **O** $\rightarrow$ Trilhos (7) | `PILAR` (ID 10, Fixo no abismo) |

---

## 🎒 2. Tabela de Objetos

O jogo possui limite de **3 itens carregados na mão** (`max_carregados = 3`), permitindo carregar a Vela acesa, a Mochila de transporte e ainda ter uma mão livre para manipular suprimentos e ferramentas:

| ID | Nome / Sinônimo | Início | Propriedades | Função no Enredo |
|---|---|---|---|---|
| **1** | `LOCAL` | — | Fixo (Situação 0) | Palavra sintática de ambiente |
| **2** | `VELA` / `VELAS` | Sala 1 | Pegável, Guardável | Objeto de iluminação (acesa com pederneira) |
| **3** | `MOCHILA` / `SACO` | Sala 3 | Pegável (Recipiente) | Armazena até 2 itens extras dentro dela |
| **4** | `CHAVE` / `CHAVES` | Sala 6 | Pegável, Guardável | Destranca o baú de ferro no Escritório (5) |
| **5** | `PEDERNEIRA` / `FOSFORO` | Sala 2 | Pegável, Guardável | Tira faíscas para acender a vela (`ACENDA VELA`) |
| **6** | `CANTIL` / `GARRAFA` | Sala 1 | Pegável, Guardável | Armazena água da fonte para matar a sede no caminho |
| **7** | `CORDA` / `CORDAS` | Baú (Sala 5) | Pegável, Guardável | Amarrada no pilar da Sala 8 para descer o abismo |
| **8** | `AGUA` / `FONTE` | Sala 4 | Fixo (Situação 4) | Beber água (`BEBA AGUA`) ou encher cantil |
| **9** | `BAU` / `COFRE` | Sala 5 | Fixo (Situação 5) | Contém a corda de escalada |
| **10**| `PILAR` / `GANCHO` | Sala 8 | Fixo (Situação 8) | Ponto de ancoragem para fixar a corda |

---

## 🧭 3. Suporte aos 8 Pontos Cardeais

A engine aceita navegação pelos **8 pontos cardeais e colaterais**, com atalhos de 1 e 2 letras:
- **Norte**: `N` ou `NORTE`
- **Sul**: `S` ou `SUL`
- **Leste**: `L`, `LESTE`, `E` ou `ESTE`
- **Oeste**: `O`, `OESTE`, `W` ou `WEST`
- **Nordeste**: `NE` ou `NORDESTE`
- **Noroeste**: `NO`, `NOROESTE` ou `NW`
- **Sudeste**: `SE` ou `SUDESTE`
- **Sudoeste**: `SO`, `SUDOESTE` ou `SW`

---

## ⚙️ 4. Mecânicas e Registradores da Engine

| Registrador | Nome | Função |
|---|---|---|
| **REG 1** | `REG_POSICAO` | Sala atual onde o jogador se encontra (1 a 8). |
| **REG 5** | `REG_CONTADOR_5` | **Temporizador de Sede:** Inicia em 30 na Função 1. Decrementado automaticamente a cada turno. Se atingir 0, aciona a **Função 2** (Morte por desidratação). Quando atinge 5, emite aviso na tela. |
| **REG 6** | `REG_PASSOS_ESCURO`| **Passos no Escuro:** Se o jogador caminhar nas Salas 6 ou 7 sem a vela acesa, o contador sobe. Ao atingir 5 passos, aciona a **Função 3** (Morte por queda no fosso). |
| **REG 9** | `REG_ILUMINACAO` | Indicador de sala escura: `0` (claro) e `1` (escuro, Salas 6 e 7). |
| **REG 10**| `REG_ESTADO_OBJ2`| Estado da Vela: `0` (apagada) e `1` (acesa). |
| **REG 15**| Flag do Cantil | `0` (vazio) e `1` (cheio de água potável). |
| **REG 16**| Flag do Baú | `0` (fechado/trancado) e `1` (destrancado com a chave). |
| **REG 17**| Flag da Corda | `0` (desamarrada) e `1` (fixada no pilar do abismo). |

---

## 🏆 5. Passo a Passo Mais Rápido para Vencer (Speedrun com Backtracking)

Como você pode carregar **3 objetos** na mão e até **2 dentro da mochila**, siga a rota abaixo para nunca morrer de sede, iluminar a escuridão e obter a corda:

### Etapa 1: Coletar os Itens Iniciais e a Pederneira
1. Você inicia na **Sala 1 (Entrada da Mina)**.
2. `PEGUE VELA` *(Item 1/3 na mão)*
3. `PEGUE CANTIL` *(Item 2/3 na mão)*
4. `L` *(Move-se para a Sala 2 - Galeria)*
5. `SOLTE CANTIL` *(Deixa o cantil no chão da Galeria por enquanto)*
6. `PEGUE PEDERNEIRA` *(Pega a pederneira; agora você carrega Vela + Pederneira - 2/3)*

### Etapa 2: Saciar a Sede e Acender a Vela
7. `N` *(Move-se para a Sala 4 - Lago Subterrâneo)*
8. `BEBA AGUA` *(Sua sede é completamente saciada, reiniciando o timer de 30 turnos!)*
9. `ACENDA VELA` *(Usa a pederneira na vela: a chama acende e dissipa as trevas!)*
10. `S` *(Retorna à Sala 2 - Galeria)*

### Etapa 3: Coletar a Mochila e Guardar o Cantil Cheio
11. `L` *(Move-se para a Sala 3 - Almoxarifado)*
12. `SOLTE PEDERNEIRA` *(A vela já está acesa; deixa a pederneira aqui)*
13. `PEGUE MOCHILA` *(Coleta a mochila; você carrega Vela + Mochila - 2/3)*
14. `O` *(Retorna à Sala 2 - Galeria)*
15. `PEGUE CANTIL` *(Pega o cantil do chão; agora você carrega Vela + Mochila + Cantil - 3/3)*
16. `N` *(Vai ao Lago Subterrâneo)*
17. `ENCHA CANTIL` *(Enche o cantil até a boca!)*
18. `GUARDE CANTIL` *(Guarda o cantil cheio dentro da mochila; volta a ter 2/3 itens na mão!)*
19. `S` *(Retorna à Sala 2)*
20. `O` *(Retorna à Sala 1 - Entrada)*

### Etapa 4: Explorar a Sala Escura e Obter a Chave
21. `S` *(Entra na Sala 6 - Fenda dos Cristais. Como a vela está acesa, você enxerga tudo perfeitamente!)*
22. `PEGUE CHAVE` *(Coleta a chave de ferro)*
23. `N` *(Retorna à Sala 1 - Entrada)*
24. `L` *(Move-se para a Sala 2)*
25. `L` *(Move-se para a Sala 3 - Almoxarifado)*
26. `S` *(Entra na Sala 5 - Escritório do Feitor)*

### Etapa 5: Abrir o Baú e Coletar a Corda
27. `ABRA BAU` *(Destranca o baú de ferro com a chave; a corda de escalada surge na sala!)*
28. `SOLTE CHAVE` *(A chave não é mais necessária)*
29. `PEGUE CORDA` *(Pega a corda de escalada)*

### Etapa 6: A Marcha Final até o Abismo
30. `N` *(Retorna à Sala 3)*
31. `O` *(Retorna à Sala 2)*
32. `O` *(Retorna à Sala 1)*
33. `BEBA CANTIL` *(Bebe do cantil guardado na mochila para garantir que não sentirá sede no abismo)*
34. `S` *(Entra na Sala 6 - Fenda dos Cristais)*
35. `L` *(Entra na Sala 7 - Túnel dos Trilhos)*
36. `L` *(Chega à Sala 8 - Abismo do Túnel)*
37. `AMARRE CORDA` *(Amarra a corda firme no pilar de pedra à beira do precipício!)*
38. `DESCA` *(Desce pela corda até o fundo do abismo)*

🎉 **Resultado:**
> *"Segurando na corda, voce desce pelo abismo ate alcancar a saida da mina sob a luz do dia! Parabens, voce sobreviveu!"*  
> **Vitória da partida com sucesso!**

---

## 🧪 5. Roteiro de Teste Completo das Mecânicas

### Teste A: Morte por Sede
- Inicie a partida e execute comandos repetidos (como `L` e `O` ou examine objetos) por 30 turnos sem beber água.
- No turno 5 restante, a mensagem de aviso `MSG 52` surge: *"Sua boca esta seca como poeira..."*
- Ao zerar o timer, `MSG 50` é disparada e o jogo finaliza com Game Over por desidratação (com tela pausada para leitura, sem reiniciar o hardware).

### Teste B: Morte no Escuro
- Inicie a partida e entre diretamente na Sala 6 (`S`) sem acender a vela.
- Dê 5 passos no escuro (`L`, `O`, `L`, `O`, `L`).
- `MSG 51` é disparada e o jogo finaliza por queda no fosso.

### Teste C: Limite de Inventário
- Carregue 3 itens na mão (`VELA`, `CANTIL`, `PEDERNEIRA`) e tente pegar um 4º item:
  - O sistema exibe: *"Nao da para carregar mais nada."*
  - Requer o uso da `MOCHILA` (`GUARDE <item>`) ou deixar itens no chão (`SOLTE <item>`).
