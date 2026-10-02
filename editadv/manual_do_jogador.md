# Manual do Explorador & Guia de Sobrevivência
## *A Mina do Abismo*

Bem-vindo a **A Mina do Abismo**, uma aventura de texto interativa e sobrevivência desenvolvida para os computadores padrão MSX.

---

## 📖 1. A Ambientação e o Desafio

Você liderava uma expedição arqueológica aos confins de uma cordilheira desértica quando o solo sob seus pés cedeu subitamente. Você despencou nas galerias abandonadas de uma antiga mina de carvão subterrânea.

A queda bloqueou a entrada principal com toneladas de rochas e cascalho. O ar nas profundezas é quente, sufocante e seco. O seu tempo é limitado: a desidratação avança a cada passo, a escuridão esconde poços sem fundo e a única saída conhecida fica além de um abismo intransponível.

Seu objetivo é sobreviver, encontrar água potável, iluminar o labirinto e encontrar um meio seguro de descer o abismo vertical para escapar com vida.

---

## 🧭 2. Regras e Sistemas de Jogo

### ⏳ A Sede Implacável (Gestão de Hidratação)
- O esforço físico e o calor subterrâneo consomem sua água corporal continuamente.
- A cada comando digitado, seu nível de hidratação diminui.
- **Aviso:** Quando suas forças estiverem se esgotando (*"Sua boca está seca como poeira..."*), você terá apenas alguns turnos para beber água antes de desmaiar fatalmente.
- **Como se hidratar:**
  - Encontre a nascente de água cristalina nas profundezas e use o comando `BEBA AGUA`.
  - Use um recipiente para armazenar água e levá-la consigo para beber nos túneis mais distantes (`ENCHA CANTIL` e `BEBA CANTIL`).

### 🕯️ As Trevas Subterrâneas (Sistema de Iluminação)
- As galerias mais profundas da mina são desprovidas de luz natural e musgo fosforescente.
- **Atenção:** Andar no escuro sem uma fonte de luz acesa é extremamente perigoso. Se você der 5 passos às cegas no breu total, despencará em um poço abandonado!
- Para produzir luz, encontre uma fonte de combustível (como uma vela de sebo) e uma ferramenta capaz de tirar faíscas (pederneira). Use `ACENDA VELA` para iluminar os arredores.

### 🎒 Limite de Carga e Backtracking
- Nas condições precárias da mina, você consegue carregar até **3 itens soltos nas mãos** por vez (ex: a vela de iluminação, a mochila de transporte e uma ferramenta ou suprimento em uso).
- Se tentar pegar um 4º item, o jogo avisará que suas mãos estão cheias (*"Não dá para carregar mais nada."*).
- **Estratégia:** Utilize a mochila para armazenar até 2 itens extras (`GUARDE <item>`), liberando suas mãos para novas explorações.

---

## ⌨️ 3. Resumo de Comandos Úteis

| Comando | Descrição |
|---|---|
| `N`, `S`, `L` (ou `E`), `O` (ou `W`) | Move-se pelas 4 direções cardeais: Norte, Sul, Leste e Oeste. |
| `NE`, `NO` (ou `NW`), `SE`, `SO` (ou `SW`) | Move-se pelas 4 direções colaterais: Nordeste, Noroeste, Sudeste e Sudoeste. |
| `PEGUE <item>` | Recolhe um objeto do chão da sala atual ou de dentro da mochila. |
| `SOLTE <item>` | Deixa um objeto no chão do ambiente atual. |
| `GUARDE <item>` | Armazena um objeto dentro da mochila (se você estiver carregando a mochila). |
| `TIRE <item>` | Retira um objeto guardado dentro da mochila e coloca na mão. |
| `EXAMINE MOCHILA` | Lista todos os itens que estão guardados dentro da mochila. |
| `EXAMINE <item>` | Olha detalhadamente um item ou elemento do cenário. |
| `TEMOS` ou `INV` | Lista todos os objetos que você está carregando no momento. |
| `ENTER` (sem texto) | Repete a descrição completa da sala em que você se encontra. |
| `ACENDA VELA` | Risca a pederneira e acende a vela (requer ambos os itens). |
| `BEBA AGUA` | Bebe da fonte cristalina no lago subterrâneo. |
| `ENCHA CANTIL` | Enche o cantil com água da fonte. |
| `BEBA CANTIL` | Bebe a água armazenada no cantil em qualquer sala. |
| `ABRA BAU` | Destranca o baú de ferro (requer a chave correspondente). |
| `AMARRE CORDA` | Fixa a corda de escalada no pilar de sustentação à beira do abismo. |
| `DESCA` | Desce pelo abismo (requer que a corda já esteja amarrada com segurança). |
| `VERBO` ou `VERBOS` | Lista todos os verbos aceitos pelo analisador do jogo. |
| `INSTRUCAO` | Exibe a tela com as orientações e história inicial do jogo. |
| `DICA` ou `DICAS` | Solicita uma pista ou orientação quando estiver perdido. |

---

## 🔤 4. Acentuação e Teclado MSX

O jogo reconhece e preserva a acentuação correta da língua portuguesa:
$$\text{À, Á, Â, Ã, Ç, É, Ê, Í, Ó, Ô, Õ, Ú, Ü}$$

### ⚡ Atalhos Rápidos (Shift + 1 a 0)
- Ao segurar **Shift** e pressionar os números **1 a 9** ou **0**, o caractere acentuado correspondente é inserido instantaneamente no texto.
- Os atalhos já vêm pré-configurados de fábrica com as letras mais usadas na aventura.

### 📑 Tecla [TAB] — Tabela de Acentos Interativa
- Pressione **[TAB]** durante a digitação para abrir uma elegante janela com moldura gráfica clássica.
- Use as **setas do teclado** ($\leftarrow, \rightarrow, \uparrow, \downarrow$) para navegar entre os acentos em ordem alfabética.
- Pressione **ENTER** sobre a letra escolhida para inseri-la de volta no comando.
- Pressione **ESC** para cancelar.

### ⚙️ Tecla [SELECT] — Personalizar os Atalhos
- Deseja trocar os atalhos `Shift+1..0` pelas suas letras favoritas? Pressione **[SELECT]**.
- Com as setas, escolha o número do atalho que quer alterar (ex: `1:..`) e aperte **ENTER**.
- O cursor subirá para a tabela de acentos: escolha a nova letra e aperte **ENTER**.
- A nova configuração passa a valer imediatamente! Pressione **ESC** para sair.

---

## 💡 5. Dicas Estratégicas para o Explorador (Sem Spoilers Críticos)

1. **Priorize a água cedo:** Não se aventure pelas fendas profundas antes de localizar uma fonte estável de água potável. Saber onde se hidratar é a diferença entre a vida e a morte.
2. **Crie um depósito seguro:** A Galeria Central (Sala 2) e o Almoxarifado (Sala 3) são salas centrais e claras. Elas são ótimos pontos para deixar temporariamente itens pesados que você não vai usar naquele momento.
3. **Não entre no breu despreparado:** Nunca se aventure pela Fenda dos Cristais (ao sul da Entrada) sem antes ter acendido sua vela.
4. **Examine o que encontrar:** Baús pesados de ferro e gavetas costumam proteger os equipamentos mais valiosos de expedições passadas. Procure chaves nas reentrâncias da rocha.
5. **A gravidade não perdoa:** Nunca tente descer um precipício ou abismo em um adventure de texto sem antes ter um meio seguro de sustentação fixado no cenário!

Boa sorte em sua expedição pelas profundezas da Mina do Abismo!
