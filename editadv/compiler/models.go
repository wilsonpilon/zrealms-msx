package compiler

import (
	"fmt"
	"strings"
)

// =============================================================================
// Clean-Room Reverse-Engineered Adventure Engine Models for PC (Go)
// Based on "Editor de Adventures" by Renato Degiovani
// =============================================================================

// Constantes de Limites Arquiteturais
const (
	MaxPositions   = 99  // Cap. 6: Posições numeradas de 1 a 99
	MaxObjects     = 99  // Cap. 5: Objetos numerados de 1 a 99
	MaxRegisters   = 256 // Cap. 3: Tabela de 256 registradores
	MaxVerbs       = 200 // Cap. 4: Verbos de 1 a 200 (1..34 reservados)
	MaxFunctions   = 200 // Cap. 10: Funções de 1 a 200 (1..20 especiais)
	MaxMessages    = 200 // Cap. 9: Mensagens de 1 a 200

	DefaultMaxCarried     = 5 // Cap. 12: Quantidade máxima de objetos na mão
	DefaultMaxInContainer = 3 // Cap. 12: Quantidade máxima de objetos no objeto 3

	PositionCondOffset = 100 // Cap. 6: Saída > 100 executa FUNÇÃO (saída - 100)
	ObjectOcultoOffset = 100 // Cap. 5: Situação (sala + 100) é objeto oculto
	RegisterObjOffset  = 100 // Cap. 3: Situação do objeto X acessível em REG[100 + X]
)

// =============================================================================
// 1. REGISTRADORES PRÉ-DEFINIDOS (Capítulo 3)
// =============================================================================
const (
	RegPosicao           = 1  // Posição atual do jogador (1..99)
	RegJogadasL          = 2  // Contador de jogadas LSB (incrementado a cada ENTER)
	RegJogadasH          = 3  // Contador de jogadas MSB
	RegContador4         = 4  // Contador especial incrementado a cada frase se != 0
	RegContador5         = 5  // Contador especial / bomba decrementado a cada frase.
	                          // Quando zera, executa a FUNÇÃO 2 automaticamente.
	RegPassosEscuro      = 6  // Passos no escuro: incrementado a cada frase no escuro.
	                          // Ao atingir 5, executa a FUNÇÃO 3 automaticamente.
	RegObjetosNoObj3     = 7  // Quantidade de objetos dentro do Objeto 3
	RegObjetosCarregados = 8  // Quantidade de objetos carregados pelo jogador
	RegIluminacao        = 9  // Flag iluminação do local: 0 = claro, 1 = escuro
	RegEstadoObj2        = 10 // Estado do Objeto 2 (fonte de luz): 0 = apagado, 1 = aceso
	RegRelogioMin        = 11 // Minutos do relógio de jogo
	RegRelogioHora       = 12 // Horas do relógio de jogo
	RegRelogioDia        = 13 // Dias do relógio de jogo
	RegRelogioAjuste     = 14 // Cadência do relógio de jogo

	RegUserStart = 15 // Início dos registradores de uso livre pelo autor (15..99)
	RegUserEnd   = 99
)

// =============================================================================
// 2. VERBOS PRÉ-DEFINIDOS (Capítulo 4)
// =============================================================================
const (
	VerboNorte        = 1
	VerboSul          = 2
	VerboLeste        = 3
	VerboOeste        = 4
	VerboGrave        = 5
	VerboRecupere     = 6
	VerboEntre        = 7
	VerboSuba         = 8
	VerboSaia         = 9
	VerboDesca        = 10
	VerboHoras        = 11
	VerboQuanto       = 12
	VerboTemos        = 13
	VerboRecomece     = 14
	VerboHa           = 15
	VerboGarimpe      = 16
	VerboPense        = 17
	VerboGrite        = 18
	VerboCorra        = 19
	VerboPegue        = 20
	VerboColoque      = 21
	VerboTroque       = 22
	VerboCompre       = 23
	VerboRoube        = 24
	VerboTire         = 25
	VerboQuebre       = 26
	VerboSolte        = 27
	VerboExamine      = 28
	VerboProcure      = 29
	VerboOfereca      = 30
	VerboFaca         = 31
	VerboJogue        = 32
	VerboConserte     = 33
	VerboVenda        = 34
	VerboCustomStart  = 35
)

// =============================================================================
// 3. OBJETOS, SITUAÇÃO E BYTE DE CONSISTÊNCIA (Capítulo 5)
// =============================================================================
const (
	ObjIdLocal     = 1 // Palavra LOCAL para checagens de visibilidade
	ObjIdLuz       = 2 // Objeto de iluminação (vela, tocha, lanterna)
	ObjIdContainer = 3 // Objeto recipiente (mala, saco, bolsa)

	ObjSitInexistente    = 0   // Inexistente ou palavra sintática
	ObjSitCarregado      = 250 // Na mão do jogador (inventário)
	ObjSitEmObj3Aberto   = 251 // Dentro do objeto 3 liberado para pegar (aberto)
	ObjSitEmObj3Fechado  = 253 // Dentro do objeto 3 trancado/fechado

	// Bits do Byte de Consistência
	ObjBitPegar       = 1 << 0 // Bit 0: Pode ser pego -> Função 6
	ObjBitColocarObj3 = 1 << 1 // Bit 1: Pode ser colocado no Obj 3 -> Função 7
	ObjBitTrocar      = 1 << 2 // Bit 2: Pode ser trocado -> Função 8
	ObjBitComprar     = 1 << 3 // Bit 3: Pode ser comprado -> Função 9
	ObjBitRoubar      = 1 << 4 // Bit 4: Pode ser roubado -> Função 10
	ObjBitTirar       = 1 << 5 // Bit 5: Pode ser tirado de outro -> Função 11
	ObjBitQuebrar     = 1 << 6 // Bit 6: Pode ser quebrado -> Função 12
)

// ConsistencyConfig representa as predisposições de ação de um objeto de forma declarativa.
type ConsistencyConfig struct {
	CanTake       bool `json:"pode_pegar" yaml:"pode_pegar"`               // Bit 0 -> Função 6
	CanPutInObj3  bool `json:"pode_guardar" yaml:"pode_guardar"`           // Bit 1 -> Função 7
	CanTrade      bool `json:"pode_trocar" yaml:"pode_trocar"`             // Bit 2 -> Função 8
	CanBuy        bool `json:"pode_comprar" yaml:"pode_comprar"`           // Bit 3 -> Função 9
	CanSteal      bool `json:"pode_roubar" yaml:"pode_roubar"`             // Bit 4 -> Função 10
	CanRemove     bool `json:"pode_tirar" yaml:"pode_tirar"`               // Bit 5 -> Função 11
	CanBreak      bool `json:"pode_quebrar" yaml:"pode_quebrar"`           // Bit 6 -> Função 12
}

// ToByte converte a configuração em um byte de consistência conforme o Capítulo 5.3.
func (c ConsistencyConfig) ToByte() uint8 {
	var b uint8
	if c.CanTake {
		b |= ObjBitPegar
	}
	if c.CanPutInObj3 {
		b |= ObjBitColocarObj3
	}
	if c.CanTrade {
		b |= ObjBitTrocar
	}
	if c.CanBuy {
		b |= ObjBitComprar
	}
	if c.CanSteal {
		b |= ObjBitRoubar
	}
	if c.CanRemove {
		b |= ObjBitTirar
	}
	if c.CanBreak {
		b |= ObjBitQuebrar
	}
	return b
}

// FromByte inicializa a configuração a partir de um byte.
func (c *ConsistencyConfig) FromByte(b uint8) {
	c.CanTake = (b & ObjBitPegar) != 0
	c.CanPutInObj3 = (b & ObjBitColocarObj3) != 0
	c.CanTrade = (b & ObjBitTrocar) != 0
	c.CanBuy = (b & ObjBitComprar) != 0
	c.CanSteal = (b & ObjBitRoubar) != 0
	c.CanRemove = (b & ObjBitTirar) != 0
	c.CanBreak = (b & ObjBitQuebrar) != 0
}

// Object representa um item ou palavra reconhecível pelo sistema.
type Object struct {
	ID               int               `json:"id" yaml:"id"`
	Name             string            `json:"nome" yaml:"nome"`                             // Nome principal
	Synonyms         []string          `json:"sinonimos,omitempty" yaml:"sinonimos,omitempty"` // Sinônimos aceitos
	InitialSituation uint8             `json:"situacao" yaml:"situacao"`                     // Situação inicial
	Consistency      ConsistencyConfig `json:"consistencia" yaml:"consistencia"`             // Configuração do byte de consistência
	Description      string            `json:"descricao" yaml:"descricao"`                   // Mensagem que descreve o objeto
}

// FullName retorna a string no formato MSX do Editor: "NOME/SIN1/SIN2"
func (o Object) FullName() string {
	if len(o.Synonyms) == 0 {
		return o.Name
	}
	parts := append([]string{o.Name}, o.Synonyms...)
	return strings.Join(parts, "/")
}

// =============================================================================
// 4. POSIÇÕES / SALAS (Capítulo 6 - Suporte a 8 Pontos Cardeais)
// =============================================================================
type Exits struct {
	North     int `json:"norte" yaml:"norte"`
	South     int `json:"sul" yaml:"sul"`
	East      int `json:"leste" yaml:"leste"`
	West      int `json:"oeste" yaml:"oeste"`
	Northeast int `json:"nordeste,omitempty" yaml:"nordeste,omitempty"`
	Northwest int `json:"noroeste,omitempty" yaml:"noroeste,omitempty"`
	Southeast int `json:"sudeste,omitempty" yaml:"sudeste,omitempty"`
	Southwest int `json:"sudoeste,omitempty" yaml:"sudoeste,omitempty"`
}

// ToArray converte as saídas para a ordem fixa da engine [Norte, Sul, Leste, Oeste, NE, NO, SE, SO].
func (e Exits) ToArray() [8]uint8 {
	return [8]uint8{
		uint8(e.North), uint8(e.South), uint8(e.East), uint8(e.West),
		uint8(e.Northeast), uint8(e.Northwest), uint8(e.Southeast), uint8(e.Southwest),
	}
}

// Position representa um local onde o jogador pode estar fisicamente.
type Position struct {
	ID          int    `json:"id" yaml:"id"`
	Name        string `json:"nome" yaml:"nome"`
	Description string `json:"descricao" yaml:"descricao"`
	Exits       Exits  `json:"saidas" yaml:"saidas"`
}

// =============================================================================
// 5. INSTRUÇÕES E OPCODES (Capítulo 8)
// =============================================================================
type Opcode uint8

const (
	OpNOP    Opcode = 0
	OpMSG    Opcode = 1
	OpNVC    Opcode = 2
	OpLLIST  Opcode = 3
	OpCLIST  Opcode = 4
	OpDLIST  Opcode = 5
	OpOBJ    Opcode = 6
	OpINC    Opcode = 7
	OpDEC    Opcode = 8
	OpLDR    Opcode = 9
	OpSOMA   Opcode = 10
	OpRND    Opcode = 11
	OpREGEq  Opcode = 12 // REG=
	OpREGGt  Opcode = 13 // REG>
	OpREGLt  Opcode = 14 // REG<
	OpAQUI   Opcode = 15
	OpLOCAL  Opcode = 16
	OpTEMOS  Opcode = 17
	OpSOLTA  Opcode = 18
	OpPEGA   Opcode = 19
	OpCRIA   Opcode = 20
	OpAPAG   Opcode = 21
	OpGOSUB  Opcode = 22
	OpLIBR   Opcode = 23
	OpTRC    Opcode = 24
	OpPOE    Opcode = 25
	OpESV    Opcode = 26
	OpOK     Opcode = 27
	OpREGN   Opcode = 28
	OpNVF    Opcode = 29
	OpREF    Opcode = 30
	OpFIM    Opcode = 31
	OpNEU    Opcode = 32
	OpDESC   Opcode = 33
	OpRET    Opcode = 34
	OpGOTO   Opcode = 35
	OpPAUSA  Opcode = 36
	OpFLAG   Opcode = 37
	OpEVID   Opcode = 38
	OpCLS    Opcode = 39
	OpEVDEq  Opcode = 40 // EVD=
	OpCHRS   Opcode = 41
	OpPRT    Opcode = 42
	OpDNT    Opcode = 43
	OpCMD    Opcode = 44
	OpCount  Opcode = 45
)

// OpcodeNames mapeia o mnemônico textual para o código de operação numérico.
var OpcodeMap = map[string]Opcode{
	"NOP":   OpNOP,
	"MSG":   OpMSG,
	"NVC":   OpNVC,
	"LLIST": OpLLIST,
	"CLIST": OpCLIST,
	"DLIST": OpDLIST,
	"OBJ":   OpOBJ,
	"INC":   OpINC,
	"DEC":   OpDEC,
	"LDR":   OpLDR,
	"SOMA":  OpSOMA,
	"RND":   OpRND,
	"REG=":  OpREGEq,
	"REG>":  OpREGGt,
	"REG<":  OpREGLt,
	"AQUI":  OpAQUI,
	"LOCAL": OpLOCAL,
	"TEMOS": OpTEMOS,
	"SOLTA": OpSOLTA,
	"PEGA":  OpPEGA,
	"CRIA":  OpCRIA,
	"APAG":  OpAPAG,
	"GOSUB": OpGOSUB,
	"LIBR":  OpLIBR,
	"TRC":   OpTRC,
	"POE":   OpPOE,
	"ESV":   OpESV,
	"OK":    OpOK,
	"REGN":  OpREGN,
	"NVF":   OpNVF,
	"REF":   OpREF,
	"FIM":   OpFIM,
	"NEU":   OpNEU,
	"DESC":  OpDESC,
	"RET":   OpRET,
	"GOTO":  OpGOTO,
	"PAUSA": OpPAUSA,
	"FLAG":  OpFLAG,
	"EVID":  OpEVID,
	"CLS":   OpCLS,
	"EVD=":  OpEVDEq,
	"CHRS":  OpCHRS,
	"PRT":   OpPRT,
	"DNT":   OpDNT,
	"CMD":   OpCMD,
}

// Instruction define uma instrução na história do autor (em YAML/JSON).
type Instruction struct {
	Op    string `json:"op" yaml:"op"`                                     // Mnemônico (ex: "MSG", "PEGA", "REG=")
	P1    int    `json:"p1,omitempty" yaml:"p1,omitempty"`                 // Parâmetro 1
	P2    int    `json:"p2,omitempty" yaml:"p2,omitempty"`                 // Parâmetro 2
	P3    int    `json:"p3,omitempty" yaml:"p3,omitempty"`                 // Parâmetro 3 (label em REG=, REG>, REG<)
	Label string `json:"label,omitempty" yaml:"label,omitempty"`           // Rótulo textual opcional (ex: "A", "LOOP")
}

// CompiledInstruction representa a instrução de 4 bytes compilada para SDCC/MSX.
type CompiledInstruction struct {
	Op uint8
	P1 uint8
	P2 uint8
	P3 uint8
}

// ToCompiled resolve o mnemônico para a estrutura binária compatível com Z80.
func (inst Instruction) ToCompiled() (CompiledInstruction, error) {
	opCode, found := OpcodeMap[strings.ToUpper(strings.TrimSpace(inst.Op))]
	if !found {
		return CompiledInstruction{}, fmt.Errorf("mnemônico desconhecido: %s", inst.Op)
	}
	return CompiledInstruction{
		Op: uint8(opCode),
		P1: uint8(inst.P1),
		P2: uint8(inst.P2),
		P3: uint8(inst.P3),
	}, nil
}

// =============================================================================
// 6. VERBOS (Capítulo 4)
// =============================================================================
type Verb struct {
	ID       int      `json:"id" yaml:"id"`
	Name     string   `json:"nome" yaml:"nome"`
	Synonyms []string `json:"sinonimos,omitempty" yaml:"sinonimos,omitempty"`
}

// FullName retorna a string no formato MSX: "VERBO/SINONIMO1/SINONIMO2"
func (v Verb) FullName() string {
	if len(v.Synonyms) == 0 {
		return v.Name
	}
	parts := append([]string{v.Name}, v.Synonyms...)
	return strings.Join(parts, "/")
}

// =============================================================================
// 7. COMANDOS (Capítulo 7)
// =============================================================================
// Cabeçalho estruturado: Verbo + Objeto1 + Objeto2.
type Command struct {
	Verb         int           `json:"verbo" yaml:"verbo"`                   // ID do verbo
	Object1      int           `json:"objeto1,omitempty" yaml:"objeto1,omitempty"` // ID do primeiro objeto (0 se ausente)
	Object2      int           `json:"objeto2,omitempty" yaml:"objeto2,omitempty"` // ID do segundo objeto (0 se ausente)
	Comment      string        `json:"comentario,omitempty" yaml:"comentario,omitempty"`
	Instructions []Instruction `json:"instrucoes" yaml:"instrucoes"`         // Procedimento a ser executado
}

// =============================================================================
// 8. FUNÇÕES (Capítulo 10)
// =============================================================================
const (
	FuncReset          = 1  // Reset inicial
	FuncTimerBomba     = 2  // Disparada quando Reg 5 atinge 0
	FuncEscuro         = 3  // Disparada quando Reg 6 atinge 5
	FuncRealtime       = 4  // Executada continuamente no loop de tempo real
	FuncPreComando     = 5  // Executada no início de cada turno antes do comando
	FuncPadraoPegar    = 6  // Ação padrão de pegar (bit 0)
	FuncPadraoColocar  = 7  // Ação padrão de colocar no obj 3 (bit 1)
	FuncPadraoTrocar   = 8  // Ação padrão de trocar (bit 2)
	FuncPadraoComprar  = 9  // Ação padrão de comprar (bit 3)
	FuncPadraoRoubar   = 10 // Ação padrão de roubar (bit 4)
	FuncPadraoTirar    = 11 // Ação padrão de tirar de objeto (bit 5)
	FuncPadraoQuebrar  = 12 // Ação padrão de quebrar (bit 6)
	FuncPadraoSoltar   = 13 // Ação padrão de soltar
	FuncPadraoExaminar = 14 // Ação padrão de examinar
	FuncPadraoProcurar = 15 // Ação padrão de procurar
	FuncPadraoOferecer = 16 // Ação padrão de oferecer
	FuncPadraoFazer    = 17 // Ação padrão de fazer
	FuncPadraoJogar    = 18 // Ação padrão de jogar
	FuncPadraoConsertar= 19 // Ação padrão de consertar
	FuncPadraoVender   = 20 // Ação padrão de vender
	FuncUserStart      = 21 // Funções livres do autor (21..200)
)

type Function struct {
	ID           int           `json:"id" yaml:"id"`
	Name         string        `json:"nome,omitempty" yaml:"nome,omitempty"`
	Instructions []Instruction `json:"instrucoes" yaml:"instrucoes"`
}

// =============================================================================
// 9. MENSAGENS (Capítulo 9)
// =============================================================================
const (
	MsgIntro             = 11 // Apresentação inicial
	MsgAchei             = 12 // "Achei o que você queria."
	MsgEscuro            = 13 // "Está muito escuro aqui..."
	MsgNaoEntendi        = 14 // "Perdão, não entendi..."
	MsgMovimentoInvalido = 15 // "É impossível ir nesta direção."
	MsgNaoPossivel       = 16 // "Isto não é possível."
	MsgNaoTemos          = 17 // "Nós não temos..."
	MsgJaTemos           = 18 // "Nós já temos..."
	MsgNaoEstouVendo     = 19 // "Eu não estou vendo..."
	MsgObjetoComum       = 20 // "É apenas..."
	MsgCargaMaxima       = 21 // "Não dá para carregar mais nada."
	MsgObj3Lotado        = 22 // "Não cabe mais nada."
)

type Message struct {
	ID   int    `json:"id" yaml:"id"`
	Text string `json:"texto" yaml:"texto"`
}

// =============================================================================
// 10. ESTRUTURA GERAL DO ARQUIVO DA AVENTURA (Documento YAML / JSON)
// =============================================================================
type AdventureMetadata struct {
	Title       string `json:"titulo" yaml:"titulo"`
	Author      string `json:"autor" yaml:"autor"`
	Version     string `json:"versao" yaml:"versao"`
	Description string `json:"descricao,omitempty" yaml:"descricao,omitempty"`
}

type AdventureConfig struct {
	InitialPosition  int  `json:"posicao_inicial" yaml:"posicao_inicial"`
	MaxCarried       int  `json:"max_carregados,omitempty" yaml:"max_carregados,omitempty"`           // Padrão: 5
	MaxInContainer   int  `json:"max_em_recipiente,omitempty" yaml:"max_em_recipiente,omitempty"`     // Padrão: 3
	EnableRealTime   bool `json:"tempo_real,omitempty" yaml:"tempo_real,omitempty"`
}

// AdventureGame representa a estrutura completa de uma aventura.
type AdventureGame struct {
	Meta             AdventureMetadata `json:"meta" yaml:"meta"`
	Config           AdventureConfig   `json:"config" yaml:"config"`
	InitialRegisters map[int]uint8     `json:"registradores_iniciais,omitempty" yaml:"registradores_iniciais,omitempty"`
	Positions        []Position        `json:"posicoes" yaml:"posicoes"`
	Objects          []Object          `json:"objetos" yaml:"objetos"`
	Verbs            []Verb            `json:"verbos,omitempty" yaml:"verbos,omitempty"`
	Commands         []Command         `json:"comandos" yaml:"comandos"`
	Functions        []Function        `json:"funcoes" yaml:"funcoes"`
	Messages         []Message         `json:"mensagens" yaml:"mensagens"`
}

// Validate realiza validações estruturais essenciais conforme o manual.
func (g *AdventureGame) Validate() error {
	if g.Config.InitialPosition <= 0 || g.Config.InitialPosition > MaxPositions {
		return fmt.Errorf("posição inicial inválida: %d (deve ser entre 1 e %d)", g.Config.InitialPosition, MaxPositions)
	}
	if len(g.Positions) > MaxPositions {
		return fmt.Errorf("excesso de posições: %d (máximo permitido: %d)", len(g.Positions), MaxPositions)
	}
	if len(g.Objects) > MaxObjects {
		return fmt.Errorf("excesso de objetos: %d (máximo permitido: %d)", len(g.Objects), MaxObjects)
	}
	return nil
}
