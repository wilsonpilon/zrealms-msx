# As Catacumbas de Cristal: O Desafio do Rei Esquecido
## Tutorial Passo a Passo para Criar o Jogo de Referência Completo no Z-Realm

> **Projeto de Referência Oficial:** Z-Realm (`zrealm-msx`) v0.5.3  
> **Gênero:** cRPG de Exploração em Labirinto por Salas (Grid Walk 8x8)  
> **Plataforma:** MSX 2 / MSX 2+ (SCREEN 4, 16 Cores, Sprites Modo 2, Som PSG AY-3-8910)  
> **Formatos de Saída:** Cartucho MegaROM ASCII-16 (`.ROM`) ou Disquete MSX-DOS 2 (`.DSK` / `.COM`)  

---

## 1. Visão Geral e Game Design Document (GDD)

### 1.1 A Lore e a Premissa
Nas profundezas esquecidas do reino de Valdor jazem **As Catacumbas de Cristal**, uma cripta ancestral construída sob a antiga cidadela. Durante séculos, o local serviu de repouso aos reis caídos e seus tesouros místicos.

O herói desperta na entrada das catacumbas com apenas 75 de HP. Para selar a fenda mágica e vencer o jogo, ele deve explorar **20 salas interconectadas em um labirinto 4x5**, obter chaves de diferentes hierarquias (Bronze, Ferro e Real Dourada), desviar de armadilhas de espinhos, resgatar um prisioneiro, consultar um sábio eremita, derrotar em combate de tempo real um esqueleto guerreiro guardião e erguer o lendário **Cristal Primordial** no Santuário do Rei Esquecido.

### 1.2 Mapa Estrutural do Mundo (Grid Espacial 4x5)
O mundo é modelado como uma malha retangular de 4 linhas por 5 colunas:

```
          Coluna 0            Coluna 1            Coluna 2            Coluna 3            Coluna 4
Linha 0: [04: Armaria]    <-> [03: Corredor]  <-> [01: Entrada]   <-> [02: Pilares]   <-> [05: Salão Reis]
              ^                   ^                   ^                   ^                   ^
              v                   v                   v                   v                   v
Linha 1: [08: Galeria]    <-> [07: Cripta]    <-> [06: Vale]      <-> [09: Labirinto] <-> [10: Fosso]
              ^                   ^                   ^                   ^                   ^
              v                   v                   v                   v                   v
Linha 2: [12: Fonte]      <-> [11: Eremita]   <-> [13: Celas]     <-> [14: Prisão]    <-> [15: Tortura]
              ^                   ^                   ^                   ^                   ^
              v                   v                   v                   v                   v
Linha 3: [16: Esgotos]    <-> [17: Pórtico]   <-> [18: Cristais]  <-> [19: Antecâmara]<-> [20: Santuário]
```

### 1.3 Cadeia de Progressão e Resolução de Enigmas
1. **Sala 01 (Entrada):** O herói fala com o **Guardião Sentinela** e recebe a **Chave de Bronze** (Flag 1 = 1, SFX 1).
2. **Sala 04 (Armaria):** O herói abre o baú guardado e obtém a **Chave de Ferro** e uma **Poção de Cura** (+25 HP, Flag 2 = 1, SFX 1).
3. **Sala 14 (Cárcere):** Usando a Chave de Ferro, o herói liberta o **Prisioneiro Trancafiado**, que lhe entrega o **Amuleto de Cristal** (Flag 3 = 1, SFX 1).
4. **Sala 11 (Refúgio do Eremita):** O herói apresenta o Amuleto ao **Eremita Sábio**, que decifra a profecia e revela que o **Esqueleto Guerreiro** na Câmara de Tortura porta a Chave Real (Flag 4 = 1, SFX 1).
5. **Sala 12 (Fonte Sagrada):** O herói bebe da água benta curativa, restaurando seu HP para 100/100 (SFX 1).
6. **Sala 15 (Câmara de Tortura):** O herói enfrenta o **Esqueleto Guerreiro** em combate corpo-a-corpo em tempo real. O golpe de espada destrói o monstro e dropa a **Chave Real Dourada** (Flag 5 = 1, SFX 1).
7. **Sala 19 (Antecâmara Real):** O herói usa a Chave Real para destrancar a grande **Porta Real** (Flag 6 = 1, SFX 3).
8. **Sala 20 (Santuário do Rei Esquecido):** O herói ergue o **Cristal Primordial**, alcançando a **Vitória Suprema** (Flag 7 = 1, SFX 1)!

---

## 2. Passo 1: Inicializando o Projeto no Editor

1. Abra o **Z-Realm**:
   ```powershell
   .\zrealm.exe
   ```
2. No menu superior ou no *Painel de Início Rápido*:
   - Clique em **`Criar Novo Projeto (.rpgproj)`**.
   - Defina o nome do arquivo como `catacumbas.rpgproj` e o título como `"As Catacumbas de Cristal"`.
3. O editor criará automaticamente o banco SQLite relacional com todas as tabelas migradas e o schema pronto.

---

## 3. Passo 2: Desenhando os 13 Tiles na Aba "Tilesets (8x8)"

No padrão **MSX 2 (SCREEN 4)**, cada tile de 8x8 pixels possui 8 bytes de padrão de bits (1 = frente, 0 = fundo) e 8 bytes de atributos de cores de 4 bits (Nibble alto = cor de frente, Nibble baixo = cor de fundo).

Acesse a aba **🧱 Tilesets (8x8)** e crie/configure os seguintes tiles:

| ID | Nome do Tile | Tipo de Colisão | Cor Predominante (Fg/Bg) | Função no Cenário |
|:--:|:-------------|:----------------|:-------------------------|:------------------|
| **0** | **Chão de Pedra** | `Passável (0)` | Cinza Claro (9) / Preto (1) | Piso comum navegável em todas as salas |
| **1** | **Parede de Tijolos** | `Sólido (1)` | Cinza (14) / Azul Escuro (4) | Delimitação externa e paredes internas |
| **2** | **Portal de Arco** | `Gatilho (4)` | Amarelo Claro (10) / Preto (1) | Portas mágicas e acessos ornamentados |
| **3** | **Água Cristalina** | `Água (2)` | Azul Claro (5) / Azul Escuro (4) | Canais de água nos esgotos e vales |
| **4** | **Espinhos Ardentes** | `Dano (3)` | Vermelho Escuro (6) / Preto (1) | Armadilha de chão: drena 5 HP por toque |
| **5** | **Tocha na Parede** | `Sólido (1)` | Amarelo (10) / Preto (1) | Iluminação atmosférica das catacumbas |
| **6** | **Estátua Antiga** | `Sólido (1)` | Cinza (14) / Preto (1) | Esculturas fúnebres decorativas |
| **7** | **Altar Místico** | `Sólido (1)` | Amarelo Brilhante (11) / Preto (1) | Suporte ritual do templo e santuário |
| **8** | **Grades da Cela** | `Sólido (1)` | Cinza (14) / Preto (1) | Jaula de retenção do prisioneiro |
| **9** | **Lajotas Antigas** | `Passável (0)` | Cinza Claro (9) / Preto (1) | Chão desgastado por séculos |
| **10**| **Cristal Místico** | `Sólido (1)` | Ciano Brilhante (7) / Preto (1) | Formações de cristal azul brilhante |
| **11**| **Fonte de Água** | `Sólido (1)` | Ciano (7) / Azul Escuro (4) | Bacia da fonte sagrada de cura |
| **12**| **Alavanca Bronze** | `Sólido (1)` | Amarelo (10) / Preto (1) | Mecanismo de abertura de portas |

> **Dica de Design:** Configure as cores de cada linha na lateral direita e selecione a propriedade física correta (ex: marque `Dano (3)` para o Tile 4 de espinhos).

---

## 4. Passo 3: Criando os 8 Sprites na Aba "Sprites (16x16)"

No V9938 (Sprites Modo 2), cada sprite de 16x16 pixels é composto por 32 bytes de máscara de pixels e 16 bytes de atributos de cores (1 cor independente por scanline vertical).

Acesse a aba **👾 Sprites (16x16)** e desenhe os 8 sprites do jogo:

1. **Sprite 1 — Herói Guerreiro (`Heroi`):**
   - Cores: Branco brilhante (`0x0F`) nos elmos e lâmina.
   - Animação: 16x16 pixels com espada na mão e pernas separadas.
2. **Sprite 2 — Guardião Sentinela (`Guardiao`):**
   - Cores: Amarelo (`0x0A`) nas linhas 0-3 (elmo de ouro), Ciano (`0x07`) nas linhas 4-7 (armadura peitoral), Azul escuro (`0x04`) nas linhas 8-11 (cota de malha) e Cinza (`0x0E`) nas linhas 12-15.
3. **Sprite 3 — Baú de Tesouro (`Bau`):**
   - Cores: Amarelo ouro (`0x0A`) na tampa e fechadura, Vermelho escuro (`0x06`) na madeira do corpo.
4. **Sprite 4 — Esqueleto Guerreiro (`Esqueleto`):**
   - Cores: Branco (`0x0F`) no crânio, Vermelho médio (`0x08`) nas órbitas oculares e coração, Cinza claro (`0x0E`) nas costelas e pernas.
   - Silhueta: Empunha uma espada longa vertical pronta para o ataque.
5. **Sprite 5 — Eremita Sábio (`Eremita`):**
   - Cores: Cinza claro (`0x0E`) na longa barba e cabelos, Verde floresta (`0x0C`) na túnica de eremita.
   - Silhueta: Apoia-se em um longo cajado de madeira.
6. **Sprite 6 — Prisioneiro Trancafiado (`Prisioneiro`):**
   - Cores: Amarelo pálido (`0x0B`) no rosto desnutrido, Marrom escuro (`0x06`) nos trapos de prisioneiro. Correntes desenhadas nos punhos.
7. **Sprite 7 — Cristal Primordial (`Cristal`):**
   - Cores: Ciano reluzente (`0x07`) na parte superior e Azul místico (`0x05`) na base. Formato de diamante lapidado cintilante.
8. **Sprite 8 — Goblin Ladino (`Goblin`):**
   - Cores: Verde goblin (`0x03`) na pele e orelhas pontudas, Amarelo couro (`0x0A`) no colete e pernas.

---

## 5. Passo 4: Cadastrando Itens na Aba "Regras & RPG Stats"

Acesse a aba **⚔️ Regras & RPG Stats** e cadastre os 7 itens na seção **Itens & Equipamentos**:

| ID | Nome do Item | Tipo | Modificador | Valor | Preço (GP) | Descrição |
|:--:|:-------------|:-----|:------------|:-----:|:----------:|:----------|
| **1** | `Chave de Bronze` | `Chave (Key)` | Nenhum | 0 | 0 | Entregue pelo Guardião na entrada |
| **2** | `Pocao de Vida` | `Consumível` | `HP Máximo (1)` | +25 | 10 | Encontrada no baú da Armaria |
| **3** | `Chave de Ferro` | `Chave (Key)` | Nenhum | 0 | 0 | Chave pesada que abre celas |
| **4** | `Chave Real Dourada` | `Chave (Key)` | Nenhum | 0 | 0 | Dropada pelo Esqueleto Guardião |
| **5** | `Amuleto de Cristal` | `Quest` | Nenhum | 0 | 0 | Presente de gratidão do prisioneiro |
| **6** | `Elixir Magico` | `Consumível` | `MP Máximo (2)` | +30 | 20 | Tônico revigorante de magia |
| **7** | `Cristal Primordial` | `Quest` | Nenhum | 0 | 0 | O artefato sagrado da vitória final |

---

## 6. Passo 5: Configurando Diálogos na Aba "Diálogos & Roteiros"

Na sub-aba **Tabela de Diálogos**, registre as 15 mensagens que darão voz à narrativa:

* **MSG 1 (`MSG_WELCOME`):** `"Bem-vindo as Catacumbas de Cristal! Pressione [ESPACO] para interagir."`
* **MSG 2 (`MSG_GUARDIAN_1`):** `"Guardiao: As profundezas sao perigosas. Tome esta Chave de Bronze!"`
* **MSG 3 (`MSG_GUARDIAN_2`):** `"Guardiao: Que a luz guie seus passos nas profundezas."`
* **MSG 4 (`MSG_CHEST_ARMORY`):** `"Voce abriu o bau e pegou a Chave de Ferro e uma Pocao de Vida!"`
* **MSG 5 (`MSG_CHEST_EMPTY`):** `"O bau de tesouro esta vazio."`
* **MSG 6 (`MSG_SKELETON_DEFEATED`):** `"Esqueleto destruido! Entre os ossos reluz a Chave Real Dourada!"`
* **MSG 7 (`MSG_PRISONER_LOCKED`):** `"Prisioneiro: Socorro! Estou preso! Preciso da Chave de Ferro da Armaria!"`
* **MSG 8 (`MSG_PRISONER_FREE`):** `"Prisioneiro: Livre! Muito obrigado! Leve este Amuleto de Cristal!"`
* **MSG 9 (`MSG_PRISONER_THANKS`):** `"Prisioneiro: Va em frente! Encontre o Eremita e o Santuario!"`
* **MSG 10 (`MSG_HERMIT_RIDDLE`):** `"Eremita: Traga o Amuleto de Cristal do prisioneiro para decifrar a profecia."`
* **MSG 11 (`MSG_HERMIT_SOLVED`):** `"Eremita: A profecia revela: o Esqueleto na Tortura guarda a Chave Real!"`
* **MSG 12 (`MSG_FOUNTAIN_HEAL`):** `"Voce bebe da Fonte Sagrada. Sua vida foi completamente restaurada! (+100 HP)"`
* **MSG 13 (`MSG_DOOR_LOCKED`):** `"Porta Real: Trancada com selo magico de ouro. Requer a Chave Real!"`
* **MSG 14 (`MSG_DOOR_UNLOCKED`):** `"O selo dourado se rompe e a Porta do Santuario se abre!"`
* **MSG 15 (`MSG_VICTORY`):** `"VITORIA! Voce ergueu o Cristal Primordial e salvou o Reino de Z-Realm!"`

---

## 7. Passo 6: Compilando os 8 Scripts na Bytecode VM

Na sub-aba **Scripts de Eventos (VM)**, crie os 8 scripts usando os mnemônicos da máquina virtual do Z80 e clique em **`Compilar Bytecode`**:

### Script 1: Guardião na Entrada (`Script Guardiao`)
```assembly
CHECK_FLAG 1 ja_falou
MSG 2
GIVE_ITEM 1
SET_FLAG 1 1
PLAY_SFX 1
END
ja_falou:
MSG 3
END
```

### Script 2: Baú da Armaria (`Script Bau Armaria`)
```assembly
CHECK_FLAG 2 bau_aberto
MSG 4
GIVE_ITEM 3
GIVE_ITEM 2
HEAL 25
SET_FLAG 2 1
PLAY_SFX 1
END
bau_aberto:
MSG 5
END
```

### Script 3: Esqueleto Guerreiro (`Script Esqueleto`)
```assembly
CHECK_FLAG 5 ja_morto
MSG 6
GIVE_ITEM 4
SET_FLAG 5 1
PLAY_SFX 1
END
ja_morto:
END
```

### Script 4: Prisioneiro nas Celas (`Script Prisioneiro`)
```assembly
CHECK_FLAG 3 livre
CHECK_FLAG 2 tem_chave
MSG 7
END
tem_chave:
MSG 8
GIVE_ITEM 5
SET_FLAG 3 1
PLAY_SFX 1
END
livre:
MSG 9
END
```

### Script 5: Eremita Sábio (`Script Eremita`)
```assembly
CHECK_FLAG 4 decifrado
CHECK_FLAG 3 tem_amuleto
MSG 10
END
tem_amuleto:
MSG 11
SET_FLAG 4 1
PLAY_SFX 1
END
decifrado:
MSG 11
END
```

### Script 6: Fonte Sagrada (`Script Fonte Sagrada`)
```assembly
MSG 12
HEAL 100
PLAY_SFX 1
END
```

### Script 7: Porta Real (`Script Porta Real`)
```assembly
CHECK_FLAG 6 porta_aberta
CHECK_FLAG 5 tem_chave_real
MSG 13
END
tem_chave_real:
MSG 14
SET_FLAG 6 1
PLAY_SFX 3
END
porta_aberta:
END
```

### Script 8: Cristal da Vitória (`Script Cristal Vitoria`)
```assembly
CHECK_FLAG 7 ja_venceu
MSG 15
GIVE_ITEM 7
SET_FLAG 7 1
PLAY_SFX 1
END
ja_venceu:
MSG 15
END
```

---

## 8. Passo 7: Criando as 20 Salas na Aba "Salas (32x18)"

Cada sala possui **32 colunas x 18 linhas** (576 tiles).

### 8.1 Regra de Ouro do Alinhamento das Portas Cardeais
Para que o herói transite suavemente entre salas sem travar nas paredes:
- **Portas Norte/Sul:** Devem estar centralizadas nas colunas **X = 15 e X = 16** (linhas Y = 0 no Norte e Y = 17 no Sul).
- **Portas Leste/Oeste:** Devem estar centralizadas nas linhas **Y = 8 e Y = 9** (coluna X = 0 no Oeste e X = 31 no Leste).
- **Nunca coloque estátuas, paredes, água ou obstáculos sólidos** nestes eixos de travessia!

### 8.2 Tabela das 20 Salas e Coordenadas
Crie as salas na aba **Salas** preenchendo as conexões cardeais na barra lateral direita:

| ID | Nome da Sala | Linha (Y) | Coluna (X) | Conexão Norte | Conexão Sul | Conexão Leste | Conexão Oeste |
|:--:|:-------------|:---------:|:----------:|:-------------:|:-----------:|:-------------:|:-------------:|
| **1** | `Entrada das Catacumbas` | 0 | 2 | *(Nenhuma)* | Sala 6 | Sala 2 | Sala 3 |
| **2** | `Camara dos Pilares` | 0 | 3 | *(Nenhuma)* | Sala 9 | Sala 5 | Sala 1 |
| **3** | `Corredor das Sombras` | 0 | 1 | *(Nenhuma)* | Sala 7 | Sala 1 | Sala 4 |
| **4** | `Armaria dos Antigos` | 0 | 0 | *(Nenhuma)* | Sala 8 | Sala 3 | *(Nenhuma)* |
| **5** | `Salao dos Reis` | 0 | 4 | *(Nenhuma)* | Sala 10 | *(Nenhuma)* | Sala 2 |
| **6** | `Vale das Almas` | 1 | 2 | Sala 1 | Sala 13 | Sala 9 | Sala 7 |
| **7** | `Cripta dos Herois` | 1 | 1 | Sala 3 | Sala 11 | Sala 6 | Sala 8 |
| **8** | `Galeria Subterranea` | 1 | 0 | Sala 4 | Sala 12 | Sala 7 | *(Nenhuma)* |
| **9** | `Labirinto de Pedra` | 1 | 3 | Sala 2 | Sala 14 | Sala 10 | Sala 6 |
| **10**| `Fosso de Espinhos` | 1 | 4 | Sala 5 | Sala 15 | *(Nenhuma)* | Sala 9 |
| **11**| `Refugio do Eremita` | 2 | 1 | Sala 7 | Sala 17 | Sala 13 | Sala 12 |
| **12**| `Fonte Sagrada` | 2 | 0 | Sala 8 | Sala 16 | Sala 11 | *(Nenhuma)* |
| **13**| `Celas Subterraneas` | 2 | 2 | Sala 6 | Sala 18 | Sala 14 | Sala 11 |
| **14**| `Carcere do Prisioneiro`| 2 | 3 | Sala 9 | Sala 19 | Sala 15 | Sala 13 |
| **15**| `Camara de Tortura` | 2 | 4 | Sala 10 | Sala 20 | *(Nenhuma)* | Sala 14 |
| **16**| `Esgotos da Cidadela` | 3 | 0 | Sala 12 | *(Nenhuma)* | Sala 17 | *(Nenhuma)* |
| **17**| `Portico Antigo` | 3 | 1 | Sala 11 | *(Nenhuma)* | Sala 18 | Sala 16 |
| **18**| `Caverna de Cristais` | 3 | 2 | Sala 13 | *(Nenhuma)* | Sala 19 | Sala 17 |
| **19**| `Antecamara Real` | 3 | 3 | Sala 14 | *(Nenhuma)* | Sala 20 | Sala 18 |
| **20**| `Santuario do Rei Esquecido` | 3 | 4 | Sala 15 | *(Nenhuma)* | *(Nenhuma)* | Sala 19 |

---

## 9. Passo 8: Inserindo Entidades e Atores nas Salas

Na barra lateral direita da aba **Salas**, no painel **Entidades e Atores**, adicione os seguintes personagens e objetos interativos:

1. **Na Sala 1 (Entrada):**
   - Nome: `Guardiao Real` | Posição: `X=14, Y=7` | Sprite: `Sprite 2 (Guardiao)` | Comportamento: `NPC Estático (1)` | Script: `Script 1 (Guardiao)`.
2. **Na Sala 4 (Armaria):**
   - Nome: `Bau da Armaria` | Posição: `X=8, Y=5` | Sprite: `Sprite 3 (Bau)` | Comportamento: `Baú (3)` | Script: `Script 2 (Bau Armaria)`.
   - Nome: `Sentinela da Armaria` | Posição: `X=24, Y=5` | Sprite: `Sprite 2 (Guardiao)` | Comportamento: `NPC Andarilho (4)`.
3. **Na Sala 2 (Pilares):**
   - Nome: `Sentinela dos Pilares` | Posição: `X=6, Y=6` | Sprite: `Sprite 2 (Guardiao)` | Comportamento: `NPC Patrulheiro (5)`.
4. **Na Sala 10 (Fosso de Espinhos):**
   - Nome: `Goblin Ladino` | Posição: `X=16, Y=6` | Sprite: `Sprite 8 (Goblin)` | Comportamento: `NPC Patrulheiro (5)`.
5. **Na Sala 11 (Refúgio do Eremita):**
   - Nome: `Eremita Sabio` | Posição: `X=15, Y=8` | Sprite: `Sprite 5 (Eremita)` | Comportamento: `NPC Estático (1)` | Script: `Script 5 (Eremita)`.
6. **Na Sala 12 (Fonte Sagrada):**
   - Nome: `Fonte Sagrada` | Posição: `X=15, Y=12` | Sprite: `Sprite 3 (Bau)` | Comportamento: `Baú (3)` | Script: `Script 6 (Fonte Sagrada)`.
7. **Na Sala 14 (Cárcere do Prisioneiro):**
   - Nome: `Prisioneiro Trancafiado` | Posição: `X=16, Y=6` | Sprite: `Sprite 6 (Prisioneiro)` | Comportamento: `NPC Estático (1)` | Script: `Script 4 (Prisioneiro)`.
8. **Na Sala 15 (Câmara de Tortura):**
   - Nome: `Esqueleto Guerreiro` | Posição: `X=16, Y=9` | Sprite: `Sprite 4 (Esqueleto)` | Comportamento: `Hostil / Inimigo (6)` | Script: `Script 3 (Esqueleto)`.
9. **Na Sala 19 (Antecâmara Real):**
   - Nome: `Porta Real` | Posição: `X=30, Y=8` | Sprite: `Sprite 3 (Bau)` | Comportamento: `Porta (2)` | Script: `Script 7 (Porta Real)`.
10. **Na Sala 20 (Santuário do Rei Esquecido):**
    - Nome: `Cristal Primordial` | Posição: `X=15, Y=8` | Sprite: `Sprite 7 (Cristal)` | Comportamento: `Baú (3)` | Script: `Script 8 (Cristal Vitoria)`.

---

## 10. Passo 9: Configurando a Posição Inicial de Spawn

No menu **Arquivo** -> **Configurações do Projeto**:
- **Sala Inicial (`initial_room_id`):** `1` (Entrada das Catacumbas)
- **Spawn do Herói X (`initial_hero_x`):** `16`
- **Spawn do Herói Y (`initial_hero_y`):** `9`
- **Tileset Inicial (`initial_tileset`):** `1`

---

## 11. Passo 10: Compilação, Exportação e One-Click Run!

Tudo pronto! Seu cRPG para MSX 2 está completamente configurado e modelado. Agora você pode testá-lo imediatamente:

### 11.1 Pela Interface Gráfica
1. Clique na aba **🚀 Exportador MSX 2**.
2. Clique no botão azul **`[Testar Cartucho MegaROM (.ROM)]`** ou no botão ciano **`[Testar Disquete DOS2 (.DSK) [F5]]`** (ou pressione a tecla de atalho **`F5`**).
3. O Z-Realm empacota todos os dados SQLite, compila o executável/ROM, invoca o emulador **openMSX** e você estará jogando em menos de 2 segundos!

### 11.2 Pela Linha de Comando (PowerShell / Terminal)
Se preferir automatizar por scripts:
```powershell
# Execução direta com 1 clique via Cartucho MegaROM:
.\zrealm.exe -run-rom catacumbas.rpgproj

# Ou via Disquete MSX-DOS 2:
.\zrealm.exe -run catacumbas.rpgproj
```

---

## 12. Como Jogar na Engine do MSX 2

| Comando | Teclado do PC / MSX | Função no Jogo |
|:--------|:--------------------|:---------------|
| **Mover Herói** | `[Setas]` ou `[WASD]` | Caminha pelo grid de 32x18 tiles (SCREEN 4) |
| **Ação / Interagir** | `[ESPAÇO]` ou `[ENTER]` | Fala com NPCs, abre baús, investiga portas |
| **Atacar (Combate)** | `[ESPAÇO]` contra inimigo | Desfere golpe de espada no monstro hostil |
| **Avançar Diálogo** | `[ESPAÇO]` | Passa para a próxima página do diálogo |
| **Sair da Sessão** | `[ESC]` | Retorna ao MSX-DOS 2 (disquete) ou reseta na sala 1 (cartucho) |

Parabéns! Você acaba de construir um cRPG completo e nativo para a lendária arquitetura MSX 2 usando o ecossistema **Z-Realm**! 🎮🏰🗡️
