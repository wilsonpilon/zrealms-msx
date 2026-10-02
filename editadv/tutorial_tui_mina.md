# Tutorial TUI: Criando a Aventura "A Mina do Abismo" Passo a Passo

Este tutorial detalha o processo completo para recriar o adventure **"A Mina do Abismo"** através da ferramenta TUI (Text User Interface) do compilador do projeto.

---

## 🚀 1. Abrindo o Editor TUI

No terminal, a partir do diretório raiz do projeto, acesse a pasta `compiler` e inicie a ferramenta TUI com o arquivo de destino:

```powershell
cd e:\editadv\compiler
go run ./cmd/main.go ../games/mina_do_abismo.yaml
```

> **Comandos e Navegação Global da TUI:**
> - Teclas **`1` a `5`**: Alternam diretamente entre as abas:
>   - **`1`**: Mapa e Salas (Posições)
>   - **`2`**: Objetos e Consistência
>   - **`3`**: Comandos e Funções (Bytecode)
>   - **`4`**: Mensagens do Jogo
>   - **`5`**: Build e Geração da ROM
> - **`Tab` / `Shift+Tab`**: Move o foco entre listas, campos de texto e botões
> - **`F2`**: Salva o arquivo YAML
> - **`F5`**: Compila os arquivos C (`game_data.h` e `game_data.c`)
> - **`F9`**: Compila a ROM MSX de 32 KB (`advent.rom`)
> - **`Esc`**: Retorna o foco para a lista principal da aba

---

## 🗺️ 2. Aba 1: Salas do Mapa (Tecla `1`)

Cadastre as 8 salas com suas respectivas descrições e saídas cardeais:

### Sala 1: Entrada da Mina
- **ID:** `1`
- **Nome:** `Entrada da Mina`
- **Descrição:**
  `Você está na entrada desmoronada de uma antiga mina de carvão. O calor é sufocante e o ar é seco. Rochas bloqueiam a saída para a superfície ao norte. Ao sul há uma fenda sombria, e a leste uma galeria de mineração.`
- **Saídas:** Norte: `0` | Sul: `6` | Leste: `2` | Oeste: `0`

### Sala 2: Galeria de Mineração
- **ID:** `2`
- **Nome:** `Galeria de Mineração`
- **Descrição:**
  `Uma ampla galeria escavada na rocha sólida. Vigas antigas escoram o teto. Ao norte ouve-se o gotejar constante de água. Ao leste há uma passagem para o almoxarifado, e a oeste fica a entrada da mina.`
- **Saídas:** Norte: `4` | Sul: `0` | Leste: `3` | Oeste: `1`

### Sala 3: Almoxarifado Abandonado
- **ID:** `3`
- **Nome:** `Almoxarifado Abandonado`
- **Descrição:**
  `Um depósito empoeirado com prateleiras quebradas e ferramentas enferrujadas. Há uma passagem de volta para a galeria a oeste e uma porta ao sul que leva ao escritório do feitor.`
- **Saídas:** Norte: `0` | Sul: `5` | Leste: `0` | Oeste: `2`

### Sala 4: Lago Subterrâneo
- **ID:** `4`
- **Nome:** `Lago Subterrâneo`
- **Descrição:**
  `Uma gruta serena iluminada por musgo fosforescente no teto. Uma vertente de água cristalina brota da rocha, formando uma lagoa límpida antes de sumir no chão. A saída fica ao sul.`
- **Saídas:** Norte: `0` | Sul: `2` | Leste: `0` | Oeste: `0`

### Sala 5: Escritório do Feitor
- **ID:** `5`
- **Nome:** `Escritório do Feitor`
- **Descrição:**
  `O antigo escritório da mineradora. Há uma mesa de madeira apodrecida e um baú de ferro pesado no canto. A única saída fica ao norte de volta para o almoxarifado.`
- **Saídas:** Norte: `3` | Sul: `0` | Leste: `0` | Oeste: `0`

### Sala 6: Fenda dos Cristais (Sala Escura)
- **ID:** `6`
- **Nome:** `Fenda dos Cristais`
- **Descrição:**
  `Uma fenda estreita cortando a rocha. As paredes refletem pequenos brilhos se houver luz, mas o ambiente é dominado por uma escuridão profunda. Há passagens ao norte e a leste.`
- **Saídas:** Norte: `1` | Sul: `0` | Leste: `7` | Oeste: `0`

### Sala 7: Túnel dos Trilhos (Sala Escura)
- **ID:** `7`
- **Nome:** `Túnel dos Trilhos`
- **Descrição:**
  `Um longo túnel onde antigos trilhos de vagonetes enferrujam pelo chão. O ar aqui é rarefeito e escuro como breu. O túnel segue a leste para o abismo e a oeste para a fenda.`
- **Saídas:** Norte: `0` | Sul: `0` | Leste: `8` | Oeste: `6`

### Sala 8: Abismo do Túnel
- **ID:** `8`
- **Nome:** `Abismo do Túnel`
- **Descrição:**
  `O túnel termina na borda de um imenso abismo vertical. Um vento fresco sopra do fundo, indicando a saída da mina sob a luz do dia lá embaixo. Ao lado da borda há um sólido pilar de pedra.`
- **Saídas:** Norte: `0` | Sul: `0` | Leste: `0` | Oeste: `7`

---

## 🎒 3. Aba 2: Objetos e Consistência (Tecla `2`)

Defina os 10 objetos e configure o byte de consistência:

| ID | Nome | Sinônimos | Sala Inicial | Pode Pegar | Pode Guardar | Pode Tirar |
|---|---|---|---|---|---|---|
| **1** | `LOCAL` | — | `0` | ❌ Não | ❌ Não | ❌ Não |
| **2** | `VELA` | `VELAS` | `1` (Entrada) | ✅ Sim | ✅ Sim | ✅ Sim |
| **3** | `MOCHILA` | `SACO, BOLSA` | `3` (Almoxarifado) | ✅ Sim | ❌ Não (Recipiente) | ❌ Não |
| **4** | `CHAVE` | `CHAVES` | `6` (Fenda) | ✅ Sim | ✅ Sim | ✅ Sim |
| **5** | `PEDERNEIRA` | `FOSFORO` | `2` (Galeria) | ✅ Sim | ✅ Sim | ✅ Sim |
| **6** | `CANTIL` | `GARRAFA` | `1` (Entrada) | ✅ Sim | ✅ Sim | ✅ Sim |
| **7** | `CORDA` | `CORDAS` | `0` (Oculta no Baú) | ✅ Sim | ✅ Sim | ✅ Sim |
| **8** | `AGUA` | `FONTE, LAGO` | `4` (Lago) | ❌ Não | ❌ Não | ❌ Não |
| **9** | `BAU` | `COFRE` | `5` (Escritório) | ❌ Não | ❌ Não | ❌ Não |
| **10**| `PILAR` | `GANCHO, BORDA`| `8` (Abismo) | ❌ Não | ❌ Não | ❌ Não |

---

## 💬 4. Aba 4: Mensagens do Jogo (Tecla `4`)

Cadastre as mensagens do enredo:
- **11**: `A Mina do Abismo - Uma expedição perigosa em uma mina abandonada.`
- **30**: `Você se ajoelhou e bebeu da água cristalina da fonte. A sede foi saciada!`
- **31**: `Você encheu o cantil até a boca com a água fresca da lagoa.`
- **32**: `Você abriu o cantil e bebeu goles generosos de água. Sua sede desapareceu!`
- **33**: `Com a pederneira você tirou faíscas e acendeu a vela! A chama dissipa a escuridão.`
- **34**: `Você usou a chave e destrancou o baú de ferro. Dentro dele há uma corda de escalada!`
- **35**: `Você amarrou firmemente a corda no pilar de pedra. Ela desce pelo abismo até o solo!`
- **36**: `Segurando na corda, você desce pelo abismo até alcançar a saída da mina sob a luz do dia! Parabéns, você sobreviveu!`
- **37**: `O abismo tem dezenas de metros de queda livre! Descer sem uma corda seria fatal.`
- **38**: `O baú já está destrancado.`
- **39**: `Seu cantil está completamente vazio!`
- **40**: `A corda já está amarrada no pilar.`
- **41**: `Você precisa de algo para tirar faíscas e acender a vela.`
- **42**: `Não há água para encher o cantil aqui.`
- **50**: `Sua garganta secou e você sucumbiu à sede implacável da mina... Fim da partida.`
- **51**: `Caminhando às cegas no escuro, você despencou em um poço profundo! Fim da partida.`
- **52**: `Sua boca está seca como poeira. A sede consome suas forças! Encontre água logo!`
- **53**: `Pegou.`
- **54**: `Soltou no chão.`
- **55**: `Guardou dentro da mochila.`

---

## ⚙️ 5. Aba 3: Comandos e Funções Bytecode (Tecla `3`)

### Comandos:
1. **Beber Água da Fonte** (`BEBA AGUA`): Verbo `35`, Objeto1 `8`
   - `LOCAL 4 3` $\rightarrow$ `MSG 19` $\rightarrow$ `NVC` $\rightarrow$ `LDR 5 14` $\rightarrow$ `MSG 30` $\rightarrow$ `NVC`
2. **Encher o Cantil** (`ENCHA CANTIL`): Verbo `21`, Objeto1 `6`
   - `TEMOS 6 3` $\rightarrow$ `MSG 17` $\rightarrow$ `NVC` $\rightarrow$ `LOCAL 4 6` $\rightarrow$ `MSG 42` $\rightarrow$ `NVC` $\rightarrow$ `LDR 15 1` $\rightarrow$ `MSG 31` $\rightarrow$ `NVC`
3. **Beber do Cantil** (`BEBA CANTIL`): Verbo `35`, Objeto1 `6`
   - `TEMOS 6 3` $\rightarrow$ `MSG 17` $\rightarrow$ `NVC` $\rightarrow$ `REG= 15 1 6` $\rightarrow$ `MSG 39` $\rightarrow$ `NVC` $\rightarrow$ `LDR 15 0` $\rightarrow$ `LDR 5 14` $\rightarrow$ `MSG 32` $\rightarrow$ `NVC`
4. **Acender a Vela** (`ACENDA VELA`): Verbo `31`, Objeto1 `2`
   - `TEMOS 2 3` $\rightarrow$ `MSG 17` $\rightarrow$ `NVC` $\rightarrow$ `TEMOS 5 6` $\rightarrow$ `MSG 41` $\rightarrow$ `NVC` $\rightarrow$ `LDR 10 1` $\rightarrow$ `MSG 33` $\rightarrow$ `NVC`
5. **Abrir o Baú** (`ABRA BAU`): Verbo `7`, Objeto1 `9`
   - `LOCAL 5 3` $\rightarrow$ `MSG 16` $\rightarrow$ `NVC` $\rightarrow$ `REG= 16 1 7` $\rightarrow$ `TEMOS 4 9` $\rightarrow$ `MSG 16` $\rightarrow$ `NVC` $\rightarrow$ `MSG 38` $\rightarrow$ `NVC` $\rightarrow$ `LDR 16 1` $\rightarrow$ `CRIA 7` $\rightarrow$ `MSG 34` $\rightarrow$ `NVC`
6. **Amarrar a Corda** (`AMARRE CORDA`): Verbo `21`, Objeto1 `7`
   - `LOCAL 8 3` $\rightarrow$ `MSG 16` $\rightarrow$ `NVC` $\rightarrow$ `REG= 17 1 7` $\rightarrow$ `TEMOS 7 9` $\rightarrow$ `MSG 17` $\rightarrow$ `NVC` $\rightarrow$ `MSG 40` $\rightarrow$ `NVC` $\rightarrow$ `LDR 17 1` $\rightarrow$ `SOLTA 7` $\rightarrow$ `MSG 35` $\rightarrow$ `NVC`
7. **Descer pelo Abismo** (`DESCA`): Verbo `10`, Objeto1 `0`
   - `LOCAL 8 3` $\rightarrow$ `MSG 15` $\rightarrow$ `NVC` $\rightarrow$ `REG= 17 1 6` $\rightarrow$ `MSG 37` $\rightarrow$ `NVC` $\rightarrow$ `MSG 36` $\rightarrow$ `FIM`

### Funções:
- **Função 1 (Reset):** Inicializa `REG 1 = 1`, `REG 10 = 0`, `REG 5 = 14`, `REG 15 = 0`, `REG 16 = 0`, `REG 17 = 0`, `RET`.
- **Função 2 (Morte por Sede):** `MSG 50`, `FIM`.
- **Função 3 (Morte no Escuro):** `MSG 51`, `FIM`.
- **Função 5 (Turn Loop):** Se sala 6 ou 7 define `REG 9 = 1` (escuro), senão `REG 9 = 0` (claro). Se `REG 5 == 5`, exibe `MSG 52` (aviso de sede).
- **Funções 6, 7, 13, 14:** Funções padrão com verificação de limite de 3 itens carregados na mão (`REG> 8 2 7`).

---

## 💾 6. Salvando e Gerando a ROM (Tecla `5`)

1. Pressione **`F2`** para salvar o arquivo YAML completo.
2. Pressione **`F5`** para exportar o código C para `engine/src/`.
3. Pressione **`F9`** para iniciar o build da ROM MSX (`advent.rom`).
4. Abra o openMSX para jogar e testar a sobrevivência na mina!
