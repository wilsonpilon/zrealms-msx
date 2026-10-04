package project

import (
	"context"
	"fmt"

	"github.com/zrealm-msx/zrealm/pkg/models"
	"github.com/zrealm-msx/zrealm/pkg/script"
)

// CreateDemoProject cria o jogo de referência completo da Fase 5.3:
// "As Catacumbas de Cristal: O Desafio do Rei Esquecido"
// Contendo 20 salas interligadas em grid 4x5, 13 tiles, 8 sprites, 7 itens,
// IAs de NPCs estáticos/errantes/patrulhas, combate em tempo real com monstro hostil,
// quebra-cabeças progressivos de chaves e cristal primordial de vitória.
func CreateDemoProject(filePath string) (*Project, error) {
	proj, err := Create(filePath, "As Catacumbas de Cristal")
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	st := proj.Storage()

	// 1. Obter ou atualizar o Tileset 1
	tilesets, err := st.Tilesets.ListTilesets(ctx)
	if err != nil || len(tilesets) == 0 {
		return nil, fmt.Errorf("tileset inicial não encontrado: %w", err)
	}
	tilesetID := tilesets[0].ID
	tilesets[0].Name = "Catacumbas Reais"
	tilesets[0].Description = "Masmorra antiga dos reis esquecidos com canais de água e cristais"
	if err := st.Tilesets.UpdateTileset(ctx, tilesets[0]); err != nil {
		return nil, err
	}

	// -------------------------------------------------------------
	// 2. CRIAÇÃO DOS 13 TILES
	// -------------------------------------------------------------
	type tileDef struct {
		index     int
		pattern   []byte
		color     byte
		collision models.CollisionType
	}

	tiles := []tileDef{
		// Tile 0: Chão de pedra limpo (Passable)
		{
			index:     0,
			pattern:   []byte{0x00, 0x00, 0x00, 0x18, 0x18, 0x00, 0x00, 0x00},
			color:     0x91, // Fg=9 (Cinza claro), Bg=1 (Preto)
			collision: models.CollisionPassable,
		},
		// Tile 1: Parede de tijolos de alvenaria (Solid)
		{
			index:     1,
			pattern:   []byte{0xFF, 0x88, 0x88, 0xFF, 0x22, 0x22, 0xFF, 0x88},
			color:     0xE4, // Fg=14 (Cinza), Bg=4 (Azul escuro)
			collision: models.CollisionSolid,
		},
		// Tile 2: Portal de arco / Porta mágica (Trigger)
		{
			index:     2,
			pattern:   []byte{0x7E, 0x81, 0xBD, 0xA5, 0xA5, 0xBD, 0x81, 0xFF},
			color:     0xA1, // Fg=10 (Amarelo claro), Bg=1 (Preto)
			collision: models.CollisionTrigger,
		},
		// Tile 3: Água cristalina / Canaleta subterrânea (Water)
		{
			index:     3,
			pattern:   []byte{0x00, 0x24, 0x42, 0x00, 0x00, 0x18, 0x24, 0x00},
			color:     0x54, // Fg=5 (Azul claro), Bg=4 (Azul escuro)
			collision: models.CollisionWater,
		},
		// Tile 4: Espinhos no chão (Damage)
		{
			index:     4,
			pattern:   []byte{0x00, 0x10, 0x38, 0x54, 0x92, 0x10, 0x38, 0xFE},
			color:     0x61, // Fg=6 (Vermelho escuro), Bg=1 (Preto)
			collision: models.CollisionDamage,
		},
		// Tile 5: Tocha acesa na parede (Solid)
		{
			index:     5,
			pattern:   []byte{0x10, 0x38, 0x38, 0x10, 0x10, 0x10, 0x28, 0x00},
			color:     0xA1, // Fg=10 (Amarelo), Bg=1 (Preto)
			collision: models.CollisionSolid,
		},
		// Tile 6: Estátua antiga esculpida (Solid)
		{
			index:     6,
			pattern:   []byte{0x3C, 0x42, 0x3C, 0x18, 0x7E, 0x18, 0x24, 0x42},
			color:     0xE1, // Fg=14 (Cinza), Bg=1 (Preto)
			collision: models.CollisionSolid,
		},
		// Tile 7: Altar místico real (Solid)
		{
			index:     7,
			pattern:   []byte{0x00, 0x7E, 0x42, 0x5A, 0x5A, 0x42, 0x7E, 0xFF},
			color:     0xB1, // Fg=11 (Amarelo brilhante), Bg=1 (Preto)
			collision: models.CollisionSolid,
		},
		// Tile 8: Grade de ferro da cela (Solid)
		{
			index:     8,
			pattern:   []byte{0xFF, 0x24, 0x24, 0xFF, 0x24, 0x24, 0xFF, 0x24},
			color:     0xE1, // Fg=14 (Cinza), Bg=1 (Preto)
			collision: models.CollisionSolid,
		},
		// Tile 9: Chão com fissuras / lajotas antigas (Passable)
		{
			index:     9,
			pattern:   []byte{0x80, 0x80, 0x80, 0xFF, 0x08, 0x08, 0x08, 0xFF},
			color:     0x91, // Fg=9 (Cinza claro), Bg=1 (Preto)
			collision: models.CollisionPassable,
		},
		// Tile 10: Cristal místico azul brilhante (Solid)
		{
			index:     10,
			pattern:   []byte{0x10, 0x38, 0x7C, 0xFE, 0x7C, 0x38, 0x10, 0x00},
			color:     0x71, // Fg=7 (Ciano), Bg=1 (Preto)
			collision: models.CollisionSolid,
		},
		// Tile 11: Fonte de água benta (Solid)
		{
			index:     11,
			pattern:   []byte{0x00, 0x3C, 0x42, 0x99, 0x99, 0x42, 0x3C, 0x00},
			color:     0x74, // Fg=7 (Ciano), Bg=4 (Azul escuro)
			collision: models.CollisionSolid,
		},
		// Tile 12: Alavanca de bronze na parede (Solid)
		{
			index:     12,
			pattern:   []byte{0x02, 0x04, 0x08, 0x10, 0x20, 0x7E, 0x7E, 0x00},
			color:     0xA1, // Fg=10 (Amarelo), Bg=1 (Preto)
			collision: models.CollisionSolid,
		},
	}

	for _, td := range tiles {
		colors := make([]byte, 8)
		for j := range colors {
			colors[j] = td.color
		}
		t := &models.Tile{
			TilesetID:     tilesetID,
			TileIndex:     td.index,
			PatternBytes:  td.pattern,
			ColorBytes:    colors,
			CollisionType: td.collision,
		}
		if err := st.Tilesets.SaveTile(ctx, t); err != nil {
			return nil, err
		}
	}

	// -------------------------------------------------------------
	// 3. CRIAÇÃO DOS 8 SPRITES
	// -------------------------------------------------------------
	// Sprite 1: Herói Guerreiro
	sprHero := models.NewSprite("Heroi")
	for i := range sprHero.ColorBytes {
		sprHero.ColorBytes[i] = 0x0F
	}
	_ = sprHero.SetPixel(7, 2, true); _ = sprHero.SetPixel(8, 2, true)
	_ = sprHero.SetPixel(6, 3, true); _ = sprHero.SetPixel(7, 3, true); _ = sprHero.SetPixel(8, 3, true); _ = sprHero.SetPixel(9, 3, true)
	_ = sprHero.SetPixel(7, 4, true); _ = sprHero.SetPixel(8, 4, true)
	for x := 5; x <= 10; x++ { _ = sprHero.SetPixel(x, 5, true) }
	_ = sprHero.SetPixel(7, 6, true); _ = sprHero.SetPixel(8, 6, true)
	_ = sprHero.SetPixel(6, 7, true); _ = sprHero.SetPixel(9, 7, true)
	_ = sprHero.SetPixel(6, 8, true); _ = sprHero.SetPixel(9, 8, true)
	if err := st.Sprites.CreateSprite(ctx, sprHero); err != nil {
		return nil, err
	}

	// Sprite 2: Guardião Sentinela
	sprGuardian := models.NewSprite("Guardiao")
	for i := 0; i <= 3; i++ { sprGuardian.ColorBytes[i] = 0x0A }
	for i := 4; i <= 7; i++ { sprGuardian.ColorBytes[i] = 0x07 }
	for i := 8; i <= 11; i++ { sprGuardian.ColorBytes[i] = 0x04 }
	for i := 12; i <= 15; i++ { sprGuardian.ColorBytes[i] = 0x0E }
	_ = sprGuardian.SetPixel(7, 1, true); _ = sprGuardian.SetPixel(8, 1, true)
	for x := 6; x <= 9; x++ { _ = sprGuardian.SetPixel(x, 2, true); _ = sprGuardian.SetPixel(x, 3, true) }
	for x := 4; x <= 11; x++ { _ = sprGuardian.SetPixel(x, 4, true) }
	_ = sprGuardian.SetPixel(3, 5, true); _ = sprGuardian.SetPixel(12, 5, true)
	for x := 5; x <= 10; x++ { _ = sprGuardian.SetPixel(x, 6, true); _ = sprGuardian.SetPixel(x, 7, true) }
	for x := 6; x <= 9; x++ { _ = sprGuardian.SetPixel(x, 8, true); _ = sprGuardian.SetPixel(x, 9, true) }
	if err := st.Sprites.CreateSprite(ctx, sprGuardian); err != nil {
		return nil, err
	}

	// Sprite 3: Baú de Tesouro
	sprChest := models.NewSprite("Bau")
	for i := 0; i <= 7; i++ { sprChest.ColorBytes[i] = 0x0A }
	for i := 8; i <= 15; i++ { sprChest.ColorBytes[i] = 0x06 }
	for x := 4; x <= 11; x++ { _ = sprChest.SetPixel(x, 5, true) }
	for y := 6; y <= 11; y++ {
		for x := 3; x <= 12; x++ { _ = sprChest.SetPixel(x, y, true) }
	}
	_ = sprChest.SetPixel(3, 12, true); _ = sprChest.SetPixel(4, 12, true); _ = sprChest.SetPixel(11, 12, true); _ = sprChest.SetPixel(12, 12, true)
	if err := st.Sprites.CreateSprite(ctx, sprChest); err != nil {
		return nil, err
	}

	// Sprite 4: Esqueleto Guerreiro (Monstro Hostil)
	sprSkeleton := models.NewSprite("Esqueleto")
	for i := 0; i <= 5; i++ { sprSkeleton.ColorBytes[i] = 0x0F } // Crânio branco
	for i := 6; i <= 7; i++ { sprSkeleton.ColorBytes[i] = 0x08 } // Olhos/tórax avermelhado
	for i := 8; i <= 15; i++ { sprSkeleton.ColorBytes[i] = 0x0E } // Ossos cinza claro
	// Crânio
	for x := 6; x <= 9; x++ { _ = sprSkeleton.SetPixel(x, 2, true); _ = sprSkeleton.SetPixel(x, 3, true) }
	_ = sprSkeleton.SetPixel(7, 4, true); _ = sprSkeleton.SetPixel(8, 4, true)
	// Costelas e espinha
	for x := 5; x <= 10; x++ { _ = sprSkeleton.SetPixel(x, 6, true); _ = sprSkeleton.SetPixel(x, 8, true) }
	_ = sprSkeleton.SetPixel(7, 7, true); _ = sprSkeleton.SetPixel(8, 7, true)
	// Espada na mão direita
	for y := 3; y <= 9; y++ { _ = sprSkeleton.SetPixel(12, y, true) }
	_ = sprSkeleton.SetPixel(11, 8, true); _ = sprSkeleton.SetPixel(13, 8, true)
	// Pernas
	_ = sprSkeleton.SetPixel(6, 10, true); _ = sprSkeleton.SetPixel(9, 10, true)
	_ = sprSkeleton.SetPixel(6, 11, true); _ = sprSkeleton.SetPixel(9, 11, true)
	_ = sprSkeleton.SetPixel(5, 12, true); _ = sprSkeleton.SetPixel(10, 12, true)
	if err := st.Sprites.CreateSprite(ctx, sprSkeleton); err != nil {
		return nil, err
	}

	// Sprite 5: Eremita Sábio
	sprHermit := models.NewSprite("Eremita")
	for i := 0; i <= 5; i++ { sprHermit.ColorBytes[i] = 0x0E }  // Barba/cabelo branco
	for i := 6; i <= 15; i++ { sprHermit.ColorBytes[i] = 0x0C } // Manto verde floresta
	for x := 6; x <= 9; x++ { _ = sprHermit.SetPixel(x, 2, true); _ = sprHermit.SetPixel(x, 3, true) }
	for x := 5; x <= 10; x++ { _ = sprHermit.SetPixel(x, 4, true); _ = sprHermit.SetPixel(x, 5, true) }
	// Manto longo e cajado
	for y := 6; y <= 12; y++ {
		for x := 5; x <= 10; x++ { _ = sprHermit.SetPixel(x, y, true) }
		_ = sprHermit.SetPixel(3, y, true) // Cajado
	}
	_ = sprHermit.SetPixel(2, 5, true); _ = sprHermit.SetPixel(4, 5, true)
	if err := st.Sprites.CreateSprite(ctx, sprHermit); err != nil {
		return nil, err
	}

	// Sprite 6: Prisioneiro
	sprPrisoner := models.NewSprite("Prisioneiro")
	for i := 0; i <= 5; i++ { sprPrisoner.ColorBytes[i] = 0x0B } // Rosto pálido
	for i := 6; i <= 15; i++ { sprPrisoner.ColorBytes[i] = 0x06 } // Trapos castanhos
	for x := 6; x <= 9; x++ { _ = sprPrisoner.SetPixel(x, 3, true); _ = sprPrisoner.SetPixel(x, 4, true) }
	for x := 5; x <= 10; x++ { _ = sprPrisoner.SetPixel(x, 6, true); _ = sprPrisoner.SetPixel(x, 7, true) }
	// Correntes nos punhos
	_ = sprPrisoner.SetPixel(3, 7, true); _ = sprPrisoner.SetPixel(4, 7, true)
	_ = sprPrisoner.SetPixel(11, 7, true); _ = sprPrisoner.SetPixel(12, 7, true)
	for y := 8; y <= 12; y++ {
		_ = sprPrisoner.SetPixel(6, y, true); _ = sprPrisoner.SetPixel(9, y, true)
	}
	if err := st.Sprites.CreateSprite(ctx, sprPrisoner); err != nil {
		return nil, err
	}

	// Sprite 7: Cristal Primordial (Vitória)
	sprCrystal := models.NewSprite("Cristal")
	for i := 0; i <= 7; i++ { sprCrystal.ColorBytes[i] = 0x07 }  // Ciano brilhante
	for i := 8; i <= 15; i++ { sprCrystal.ColorBytes[i] = 0x05 } // Azul místico
	_ = sprCrystal.SetPixel(7, 2, true); _ = sprCrystal.SetPixel(8, 2, true)
	_ = sprCrystal.SetPixel(6, 3, true); _ = sprCrystal.SetPixel(9, 3, true)
	for x := 5; x <= 10; x++ { _ = sprCrystal.SetPixel(x, 4, true); _ = sprCrystal.SetPixel(x, 5, true) }
	for x := 4; x <= 11; x++ { _ = sprCrystal.SetPixel(x, 6, true); _ = sprCrystal.SetPixel(x, 7, true); _ = sprCrystal.SetPixel(x, 8, true) }
	for x := 5; x <= 10; x++ { _ = sprCrystal.SetPixel(x, 9, true); _ = sprCrystal.SetPixel(x, 10, true) }
	_ = sprCrystal.SetPixel(6, 11, true); _ = sprCrystal.SetPixel(9, 11, true)
	_ = sprCrystal.SetPixel(7, 12, true); _ = sprCrystal.SetPixel(8, 12, true)
	if err := st.Sprites.CreateSprite(ctx, sprCrystal); err != nil {
		return nil, err
	}

	// Sprite 8: Goblin Ladino
	sprGoblin := models.NewSprite("Goblin")
	for i := 0; i <= 5; i++ { sprGoblin.ColorBytes[i] = 0x03 }  // Verde Goblin
	for i := 6; i <= 15; i++ { sprGoblin.ColorBytes[i] = 0x0A } // Roupas de couro amarelo
	for x := 6; x <= 9; x++ { _ = sprGoblin.SetPixel(x, 3, true); _ = sprGoblin.SetPixel(x, 4, true) }
	_ = sprGoblin.SetPixel(4, 3, true); _ = sprGoblin.SetPixel(11, 3, true) // Orelhas pontudas
	for x := 5; x <= 10; x++ { _ = sprGoblin.SetPixel(x, 6, true); _ = sprGoblin.SetPixel(x, 7, true) }
	_ = sprGoblin.SetPixel(3, 7, true) // Adaga
	_ = sprGoblin.SetPixel(6, 9, true); _ = sprGoblin.SetPixel(9, 9, true)
	_ = sprGoblin.SetPixel(6, 10, true); _ = sprGoblin.SetPixel(9, 10, true)
	if err := st.Sprites.CreateSprite(ctx, sprGoblin); err != nil {
		return nil, err
	}

	// -------------------------------------------------------------
	// 4. CATÁLOGO DE ITENS (IDs 1 a 7)
	// -------------------------------------------------------------
	items := []*models.Item{
		// 1: Chave de Bronze
		{Name: "Chave de Bronze", ItemType: models.ItemKey, Price: 0},
		// 2: Poção de Vida
		{Name: "Pocao de Vida", ItemType: models.ItemConsumable, ModifierStat: 1, ModifierValue: 25, Price: 10},
		// 3: Chave de Ferro
		{Name: "Chave de Ferro", ItemType: models.ItemKey, Price: 0},
		// 4: Chave Real Dourada
		{Name: "Chave Real Dourada", ItemType: models.ItemKey, Price: 0},
		// 5: Amuleto de Cristal
		{Name: "Amuleto de Cristal", ItemType: models.ItemQuest, Price: 0},
		// 6: Elixir Magico
		{Name: "Elixir Magico", ItemType: models.ItemConsumable, ModifierStat: 2, ModifierValue: 30, Price: 20},
		// 7: Cristal Primordial do Rei
		{Name: "Cristal Primordial", ItemType: models.ItemQuest, Price: 0},
	}
	for _, it := range items {
		if err := st.GameData.CreateItem(ctx, it); err != nil {
			return nil, err
		}
	}

	// -------------------------------------------------------------
	// 5. TABELA DE STRINGS DE DIÁLOGO (IDs 1 a 15)
	// -------------------------------------------------------------
	stringsList := []*models.StringEntry{
		{ContextTag: "MSG_WELCOME", TextContent: "Bem-vindo as Catacumbas de Cristal! Pressione [ESPACO] para interagir."},
		{ContextTag: "MSG_GUARDIAN_1", TextContent: "Guardiao: As profundezas sao perigosas. Tome esta Chave de Bronze!"},
		{ContextTag: "MSG_GUARDIAN_2", TextContent: "Guardiao: Que a luz guie seus passos nas profundezas."},
		{ContextTag: "MSG_CHEST_ARMORY", TextContent: "Voce abriu o bau e pegou a Chave de Ferro e uma Pocao de Vida!"},
		{ContextTag: "MSG_CHEST_EMPTY", TextContent: "O bau de tesouro esta vazio."},
		{ContextTag: "MSG_SKELETON_DEFEATED", TextContent: "Esqueleto destruido! Entre os ossos reluz a Chave Real Dourada!"},
		{ContextTag: "MSG_PRISONER_LOCKED", TextContent: "Prisioneiro: Socorro! Estou preso! Preciso da Chave de Ferro da Armaria!"},
		{ContextTag: "MSG_PRISONER_FREE", TextContent: "Prisioneiro: Livre! Muito obrigado! Leve este Amuleto de Cristal!"},
		{ContextTag: "MSG_PRISONER_THANKS", TextContent: "Prisioneiro: Va em frente! Encontre o Eremita e o Santuario!"},
		{ContextTag: "MSG_HERMIT_RIDDLE", TextContent: "Eremita: Traga o Amuleto de Cristal do prisioneiro para decifrar a profecia."},
		{ContextTag: "MSG_HERMIT_SOLVED", TextContent: "Eremita: A profecia revela: o Esqueleto na Tortura guarda a Chave Real!"},
		{ContextTag: "MSG_FOUNTAIN_HEAL", TextContent: "Voce bebe da Fonte Sagrada. Sua vida foi completamente restaurada! (+100 HP)"},
		{ContextTag: "MSG_DOOR_LOCKED", TextContent: "Porta Real: Trancada com selo magico de ouro. Requer a Chave Real!"},
		{ContextTag: "MSG_DOOR_UNLOCKED", TextContent: "O selo dourado se rompe e a Porta do Santuario se abre!"},
		{ContextTag: "MSG_VICTORY", TextContent: "VITORIA! Voce ergueu o Cristal Primordial e salvou o Reino de Z-Realm!"},
	}
	for _, s := range stringsList {
		if err := st.GameData.CreateString(ctx, s); err != nil {
			return nil, err
		}
	}

	// -------------------------------------------------------------
	// 6. SCRIPTS DE EVENTOS DA BYTECODE VM (IDs 1 a 8)
	// -------------------------------------------------------------
	type scriptDef struct {
		name string
		src  string
	}
	scriptsDef := []scriptDef{
		// Script 1: Guardião na Entrada (Dá Chave de Bronze)
		{
			name: "Script Guardiao",
			src:  "CHECK_FLAG 1 ja_falou\nMSG 2\nGIVE_ITEM 1\nSET_FLAG 1 1\nPLAY_SFX 1\nEND\nja_falou:\nMSG 3\nEND",
		},
		// Script 2: Baú da Armaria (Dá Chave de Ferro + Poção de Vida)
		{
			name: "Script Bau Armaria",
			src:  "CHECK_FLAG 2 bau_aberto\nMSG 4\nGIVE_ITEM 3\nGIVE_ITEM 2\nHEAL 25\nSET_FLAG 2 1\nPLAY_SFX 1\nEND\nbau_aberto:\nMSG 5\nEND",
		},
		// Script 3: Esqueleto Guerreiro na Câmara de Tortura (Dá Chave Real Dourada)
		{
			name: "Script Esqueleto",
			src:  "CHECK_FLAG 5 ja_morto\nMSG 6\nGIVE_ITEM 4\nSET_FLAG 5 1\nPLAY_SFX 1\nEND\nja_morto:\nEND",
		},
		// Script 4: Prisioneiro nas Celas (Requer Chave de Ferro, dá Amuleto)
		{
			name: "Script Prisioneiro",
			src:  "CHECK_FLAG 3 livre\nCHECK_FLAG 2 tem_chave\nMSG 7\nEND\ntem_chave:\nMSG 8\nGIVE_ITEM 5\nSET_FLAG 3 1\nPLAY_SFX 1\nEND\nlivre:\nMSG 9\nEND",
		},
		// Script 5: Eremita Sábio (Enigma do Amuleto)
		{
			name: "Script Eremita",
			src:  "CHECK_FLAG 4 decifrado\nCHECK_FLAG 3 tem_amuleto\nMSG 10\nEND\ntem_amuleto:\nMSG 11\nSET_FLAG 4 1\nPLAY_SFX 1\nEND\ndecifrado:\nMSG 11\nEND",
		},
		// Script 6: Fonte Sagrada (Cura Completa +100 HP)
		{
			name: "Script Fonte Sagrada",
			src:  "MSG 12\nHEAL 100\nPLAY_SFX 1\nEND",
		},
		// Script 7: Porta Real da Antecâmara (Requer Chave Real)
		{
			name: "Script Porta Real",
			src:  "CHECK_FLAG 6 porta_aberta\nCHECK_FLAG 5 tem_chave_real\nMSG 13\nEND\ntem_chave_real:\nMSG 14\nSET_FLAG 6 1\nPLAY_SFX 3\nEND\nporta_aberta:\nEND",
		},
		// Script 8: Cristal Primordial no Santuário (Vitória!)
		{
			name: "Script Cristal Vitoria",
			src:  "CHECK_FLAG 7 ja_venceu\nMSG 15\nGIVE_ITEM 7\nSET_FLAG 7 1\nPLAY_SFX 1\nEND\nja_venceu:\nMSG 15\nEND",
		},
	}

	createdScripts := make([]*models.Script, len(scriptsDef))
	for idx, sd := range scriptsDef {
		bc, err := script.CompileScript(sd.src)
		if err != nil {
			return nil, fmt.Errorf("falha ao compilar %s: %w", sd.name, err)
		}
		scr := &models.Script{
			Name:       sd.name,
			SourceCode: sd.src,
			Bytecode:   bc,
		}
		if err := st.GameData.CreateScript(ctx, scr); err != nil {
			return nil, err
		}
		createdScripts[idx] = scr
	}

	// -------------------------------------------------------------
	// 7. CRIAÇÃO DAS 20 SALAS (GRID 4x5)
	// -------------------------------------------------------------
	type roomMeta struct {
		id     int // 1 a 20
		row    int // 0 a 3
		col    int // 0 a 4
		name   string
		custom func(m []byte)
	}

	roomMetas := []roomMeta{
		// 1: Entrada das Catacumbas (r=0, c=2)
		{id: 1, row: 0, col: 2, name: "Entrada das Catacumbas", custom: func(m []byte) {
			m[4*32+12] = 5; m[4*32+19] = 5 // Tochas na entrada
		}},
		// 2: Câmara dos Pilares (r=0, c=3)
		{id: 2, row: 0, col: 3, name: "Camara dos Pilares", custom: func(m []byte) {
			m[5*32+10] = 1; m[5*32+21] = 1; m[12*32+10] = 1; m[12*32+21] = 1 // 4 Pilares
		}},
		// 3: Corredor das Sombras (r=0, c=1)
		{id: 3, row: 0, col: 1, name: "Corredor das Sombras", custom: func(m []byte) {
			m[4*32+10] = 1; m[4*32+21] = 1; m[13*32+10] = 1; m[13*32+21] = 1
		}},
		// 4: Armaria dos Antigos (r=0, c=0)
		{id: 4, row: 0, col: 0, name: "Armaria dos Antigos", custom: func(m []byte) {
			m[5*32+8] = 5; m[5*32+23] = 5 // Tochas
			m[5*32+16] = 6                // Estátua decorativa no norte
		}},
		// 5: Salão dos Reis (r=0, c=4)
		{id: 5, row: 0, col: 4, name: "Salao dos Reis", custom: func(m []byte) {
			m[7*32+16] = 7; m[6*32+15] = 6; m[6*32+17] = 6 // Altar e estátuas
		}},
		// 6: Vale das Almas (r=1, c=2)
		{id: 6, row: 1, col: 2, name: "Vale das Almas", custom: func(m []byte) {
			for y := 4; y <= 13; y++ {
				m[y*32+7] = 3
				m[y*32+24] = 3
			}
		}},
		// 7: Cripta dos Heróis (r=1, c=1)
		{id: 7, row: 1, col: 1, name: "Cripta dos Herois", custom: func(m []byte) {
			m[5*32+8] = 6; m[5*32+23] = 6; m[12*32+8] = 6; m[12*32+23] = 6 // 4 Estátuas tumulares
		}},
		// 8: Galeria Subterrânea (r=1, c=0)
		{id: 8, row: 1, col: 0, name: "Galeria Subterranea", custom: func(m []byte) {
			for x := 10; x <= 21; x++ { m[8*32+x] = 9 } // Lajotas antigas no centro
		}},
		// 9: Labirinto de Pedra (r=1, c=3)
		{id: 9, row: 1, col: 3, name: "Labirinto de Pedra", custom: func(m []byte) {
			for x := 8; x <= 14; x++ { m[6*32+x] = 1 }
			for x := 17; x <= 23; x++ { m[11*32+x] = 1 }
		}},
		// 10: Fosso de Espinhos (r=1, c=4)
		{id: 10, row: 1, col: 4, name: "Fosso de Espinhos", custom: func(m []byte) {
			for x := 8; x <= 23; x++ {
				m[7*32+x] = 4
				m[10*32+x] = 4
			}
		}},
		// 11: Refúgio do Eremita (r=2, c=1)
		{id: 11, row: 2, col: 1, name: "Refugio do Eremita", custom: func(m []byte) {
			m[4*32+12] = 5; m[13*32+12] = 7 // Tocha e Altar fora do corredor central
		}},
		// 12: Fonte Sagrada (r=2, c=0)
		{id: 12, row: 2, col: 0, name: "Fonte Sagrada", custom: func(m []byte) {
			m[12*32+14] = 3; m[12*32+15] = 11; m[12*32+16] = 11; m[12*32+17] = 3
			m[13*32+15] = 3; m[13*32+16] = 3
		}},
		// 13: Celas Subterrâneas (r=2, c=2)
		{id: 13, row: 2, col: 2, name: "Celas Subterraneas", custom: func(m []byte) {
			for x := 6; x <= 12; x++ { m[5*32+x] = 8 }
			for x := 19; x <= 25; x++ { m[5*32+x] = 8 }
		}},
		// 14: Cárcere do Prisioneiro (r=2, c=3)
		{id: 14, row: 2, col: 3, name: "Carcere do Prisioneiro", custom: func(m []byte) {
			for x := 10; x <= 22; x++ { m[5*32+x] = 8 } // Grades no topo, corredor Y=8..9 desobstruído
			m[6*32+10] = 8; m[6*32+22] = 8
		}},
		// 15: Câmara de Tortura (r=2, c=4)
		{id: 15, row: 2, col: 4, name: "Camara de Tortura", custom: func(m []byte) {
			m[4*32+8] = 5; m[4*32+23] = 5 // Tochas sinistras
			m[13*32+10] = 4                // Espinhos fora do corredor central
		}},
		// 16: Esgotos da Cidadela (r=3, c=0)
		{id: 16, row: 3, col: 0, name: "Esgotos da Cidadela", custom: func(m []byte) {
			for x := 4; x <= 27; x++ {
				if x != 15 && x != 16 {
					m[8*32+x] = 3; m[9*32+x] = 3
				}
			}
		}},
		// 17: Pórtico Antigo (r=3, c=1)
		{id: 17, row: 3, col: 1, name: "Portico Antigo", custom: func(m []byte) {
			m[5*32+10] = 6; m[5*32+21] = 6 // Estátuas guardiãs
		}},
		// 18: Caverna de Cristais (r=3, c=2)
		{id: 18, row: 3, col: 2, name: "Caverna de Cristais", custom: func(m []byte) {
			m[4*32+8] = 10; m[4*32+23] = 10; m[13*32+8] = 10; m[13*32+23] = 10 // Cristais azuis
		}},
		// 19: Antecâmara Real (r=3, c=3)
		{id: 19, row: 3, col: 3, name: "Antecamara Real", custom: func(m []byte) {
			m[8*32+30] = 2; m[9*32+30] = 2 // Portal ornamentado para o Santuário
		}},
		// 20: Santuário do Rei Esquecido (r=3, c=4)
		{id: 20, row: 3, col: 4, name: "Santuario do Rei Esquecido", custom: func(m []byte) {
			m[7*32+13] = 10; m[7*32+18] = 10
			m[8*32+14] = 5; m[8*32+15] = 7; m[8*32+16] = 7; m[8*32+17] = 5
			m[9*32+13] = 10; m[9*32+18] = 10
		}},
	}

	// Função geradora de matriz com aberturas de portas cardeais
	buildMatrix := func(hasN, hasS, hasW, hasE bool, custom func(m []byte)) []byte {
		mat := make([]byte, models.RoomMatrixSize)
		for y := 0; y < models.RoomHeight; y++ {
			for x := 0; x < models.RoomWidth; x++ {
				idx := y*models.RoomWidth + x
				isWall := false
				if y == 0 {
					if !(hasN && (x == 15 || x == 16)) {
						isWall = true
					}
				} else if y == models.RoomHeight-1 {
					if !(hasS && (x == 15 || x == 16)) {
						isWall = true
					}
				} else if x == 0 {
					if !(hasW && (y == 8 || y == 9)) {
						isWall = true
					}
				} else if x == models.RoomWidth-1 {
					if !(hasE && (y == 8 || y == 9)) {
						isWall = true
					}
				}

				if isWall {
					mat[idx] = 1 // Parede sólida
				} else {
					mat[idx] = 0 // Chão livre
				}
			}
		}
		if custom != nil {
			custom(mat)
		}
		return mat
	}

	// Mapa em memória para associar ID às instâncias criadas
	roomMap := make(map[int]*models.Room)

	// Primeiro passo: Criar as salas no banco
	for _, rm := range roomMetas {
		hasN := (rm.row > 0)
		hasS := (rm.row < 3)
		hasW := (rm.col > 0)
		hasE := (rm.col < 4)

		mat := buildMatrix(hasN, hasS, hasW, hasE, rm.custom)
		r := &models.Room{
			WorldX:     rm.col,
			WorldY:     rm.row,
			Name:       rm.name,
			TilesetID:  tilesetID,
			TileMatrix: mat,
		}
		if err := st.Rooms.CreateRoom(ctx, r); err != nil {
			return nil, err
		}
		roomMap[rm.id] = r
	}

	// Segundo passo: Configurar as conexões cardeais entre as 20 salas no grid 4x5
	// Matriz do grid: grid[row][col] -> roomID
	grid := [4][5]int{
		{4, 3, 1, 2, 5},
		{8, 7, 6, 9, 10},
		{12, 11, 13, 14, 15},
		{16, 17, 18, 19, 20},
	}

	for row := 0; row < 4; row++ {
		for col := 0; col < 5; col++ {
			rID := grid[row][col]
			r := roomMap[rID]

			if row > 0 {
				nID := grid[row-1][col]
				nRoom := roomMap[nID]
				r.NorthRoomID = &nRoom.ID
			}
			if row < 3 {
				sID := grid[row+1][col]
				sRoom := roomMap[sID]
				r.SouthRoomID = &sRoom.ID
			}
			if col > 0 {
				wID := grid[row][col-1]
				wRoom := roomMap[wID]
				r.WestRoomID = &wRoom.ID
			}
			if col < 4 {
				eID := grid[row][col+1]
				eRoom := roomMap[eID]
				r.EastRoomID = &eRoom.ID
			}

			if err := st.Rooms.UpdateRoom(ctx, r); err != nil {
				return nil, err
			}
		}
	}

	// -------------------------------------------------------------
	// 8. CRIAÇÃO DAS ENTIDADES DINÂMICAS NAS SALAS
	// -------------------------------------------------------------
	// Sala 1: Entrada das Catacumbas
	// - Guardião Sentinela em (14, 7) com Script 1
	_ = st.Rooms.CreateEntity(ctx, &models.Entity{
		RoomID:        roomMap[1].ID,
		Name:          "Guardiao Real",
		PosX:          14,
		PosY:          7,
		SpriteID:      &sprGuardian.ID,
		BehaviorType:  models.BehaviorStaticNPC,
		EventScriptID: &createdScripts[0].ID,
	})

	// Sala 4: Armaria dos Antigos
	// - Baú de Armaria em (8, 5) com Script 2
	_ = st.Rooms.CreateEntity(ctx, &models.Entity{
		RoomID:        roomMap[4].ID,
		Name:          "Bau da Armaria",
		PosX:          8,
		PosY:          5,
		SpriteID:      &sprChest.ID,
		BehaviorType:  models.BehaviorChest,
		EventScriptID: &createdScripts[1].ID,
	})
	// - Sentinela Errante em (24, 5)
	_ = st.Rooms.CreateEntity(ctx, &models.Entity{
		RoomID:       roomMap[4].ID,
		Name:         "Sentinela da Armaria",
		PosX:         24,
		PosY:         5,
		SpriteID:     &sprGuardian.ID,
		BehaviorType: models.BehaviorWanderingNPC,
	})

	// Sala 2: Câmara dos Pilares
	// - Sentinela em Patrulha em (6, 6)
	_ = st.Rooms.CreateEntity(ctx, &models.Entity{
		RoomID:       roomMap[2].ID,
		Name:         "Sentinela dos Pilares",
		PosX:         6,
		PosY:         6,
		SpriteID:     &sprGuardian.ID,
		BehaviorType: models.BehaviorPatrolNPC,
	})

	// Sala 10: Fosso de Espinhos
	// - Goblin Ladino Patrulhando em (16, 6)
	_ = st.Rooms.CreateEntity(ctx, &models.Entity{
		RoomID:       roomMap[10].ID,
		Name:         "Goblin Ladino",
		PosX:         16,
		PosY:         6,
		SpriteID:     &sprGoblin.ID,
		BehaviorType: models.BehaviorPatrolNPC,
	})

	// Sala 11: Refúgio do Eremita
	// - Eremita Sábio em (15, 8) com Script 5
	_ = st.Rooms.CreateEntity(ctx, &models.Entity{
		RoomID:        roomMap[11].ID,
		Name:          "Eremita Sabio",
		PosX:          15,
		PosY:          8,
		SpriteID:      &sprHermit.ID,
		BehaviorType:  models.BehaviorStaticNPC,
		EventScriptID: &createdScripts[4].ID,
	})

	// Sala 12: Fonte Sagrada
	// - Fonte Interativa em (15, 12) com Script 6 (Cura 100 HP)
	_ = st.Rooms.CreateEntity(ctx, &models.Entity{
		RoomID:        roomMap[12].ID,
		Name:          "Fonte Sagrada",
		PosX:          15,
		PosY:          12,
		SpriteID:      &sprChest.ID,
		BehaviorType:  models.BehaviorChest,
		EventScriptID: &createdScripts[5].ID,
	})

	// Sala 14: Cárcere do Prisioneiro
	// - Prisioneiro em (16, 6) com Script 4
	_ = st.Rooms.CreateEntity(ctx, &models.Entity{
		RoomID:        roomMap[14].ID,
		Name:          "Prisioneiro Trancafiado",
		PosX:          16,
		PosY:          6,
		SpriteID:      &sprPrisoner.ID,
		BehaviorType:  models.BehaviorStaticNPC,
		EventScriptID: &createdScripts[3].ID,
	})

	// Sala 15: Câmara de Tortura
	// - Esqueleto Guerreiro Hostil em (16, 9) com BEHAVIOR_HOSTILE e Script 3
	_ = st.Rooms.CreateEntity(ctx, &models.Entity{
		RoomID:        roomMap[15].ID,
		Name:          "Esqueleto Guerreiro",
		PosX:          16,
		PosY:          9,
		SpriteID:      &sprSkeleton.ID,
		BehaviorType:  models.BehaviorHostile,
		EventScriptID: &createdScripts[2].ID,
	})

	// Sala 19: Antecâmara Real
	// - Porta Real em (30, 8) com Script 7
	_ = st.Rooms.CreateEntity(ctx, &models.Entity{
		RoomID:        roomMap[19].ID,
		Name:          "Porta Real",
		PosX:          30,
		PosY:          8,
		SpriteID:      &sprChest.ID,
		BehaviorType:  models.BehaviorDoor,
		EventScriptID: &createdScripts[6].ID,
	})

	// Sala 20: Santuário do Rei Esquecido
	// - Cristal Primordial em (15, 8) com Script 8 (Vitória)
	_ = st.Rooms.CreateEntity(ctx, &models.Entity{
		RoomID:        roomMap[20].ID,
		Name:          "Cristal Primordial",
		PosX:          15,
		PosY:          8,
		SpriteID:      &sprCrystal.ID,
		BehaviorType:  models.BehaviorChest,
		EventScriptID: &createdScripts[7].ID,
	})

	// -------------------------------------------------------------
	// 9. CONFIGURAÇÃO INICIAL DO JOGO
	// -------------------------------------------------------------
	initialRoom := roomMap[1] // Sala 1: Entrada das Catacumbas
	_ = proj.SetSetting("initial_room_id", fmt.Sprintf("%d", initialRoom.ID))
	_ = proj.SetSetting("initial_hero_x", "16")
	_ = proj.SetSetting("initial_hero_y", "9")
	_ = proj.SetSetting("initial_tileset", fmt.Sprintf("%d", tilesetID))

	return proj, nil
}
