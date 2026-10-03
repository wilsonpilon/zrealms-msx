// ____________________________
// Z-Realm (zrealm-msx) - Controlador e Entidade do Herói (MSX 2)
//─────────────────────────────────────────────────────────────────────────────
#include "hero.h"

Hero g_Hero;

// Padrão de fallback 16x16 caso nenhum sprite esteja no banco binário
static const u8 s_DefaultHeroPattern[32] = {
	0x00, 0x00, 0x00, 0x18, 0x3C, 0x7E, 0x7E, 0x3C,
	0x18, 0x3C, 0x7E, 0x7E, 0x3C, 0x18, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x18, 0x3C, 0x7E, 0x7E, 0x3C,
	0x18, 0x3C, 0x7E, 0x7E, 0x3C, 0x18, 0x00, 0x00
};

static const u8 s_DefaultHeroColor[16] = {
	0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F,
	0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F, 0x0F
};

void HERO_Init(u8 startTileX, u8 startTileY, u8 spriteResourceID)
{
	BinarySprite* spr;

	g_Hero.TileX = startTileX;
	g_Hero.TileY = startTileY;
	g_Hero.PixelX = startTileX * 8;
	g_Hero.PixelY = startTileY * 8;
	g_Hero.Direction = HERO_DIR_DOWN;
	g_Hero.StepCooldown = 0;
	g_Hero.Moved = FALSE;

	// Carrega o sprite do herói no slot 0 do VDP
	spr = LOADER_GetSprite(spriteResourceID);
	if (spr)
	{
		VDP_LoadSprite(0, spr->Pattern, spr->Color);
	}
	else
	{
		VDP_LoadSprite(0, s_DefaultHeroPattern, s_DefaultHeroColor);
	}

	HERO_Draw();
}

void HERO_SetPosition(u8 tileX, u8 tileY)
{
	g_Hero.TileX = tileX;
	g_Hero.TileY = tileY;
	g_Hero.PixelX = tileX * 8;
	g_Hero.PixelY = tileY * 8;
	HERO_Draw();
}

void HERO_Draw(void)
{
	VDP_SetSpritePos(0, g_Hero.PixelX, g_Hero.PixelY);
}

void HERO_Update(void)
{
	u8 joy;
	bool joyUp, joyDown, joyLeft, joyRight;
	bool kbUp, kbDown, kbLeft, kbRight;
	bool moveUp, moveDown, moveLeft, moveRight;
	i8 dx = 0;
	i8 dy = 0;

	g_Hero.Moved = FALSE;

	// Cooldown de repetição de passos
	if (g_Hero.StepCooldown > 0)
	{
		g_Hero.StepCooldown--;
	}

	// 1. Leitura de Joystick (Porta 1) com sanitização de hardware desconectado
	joy = Joystick_Read(JOY_PORT_1);
	joyUp    = IS_JOY_PRESSED(joy, JOY_INPUT_DIR_UP);
	joyDown  = IS_JOY_PRESSED(joy, JOY_INPUT_DIR_DOWN);
	joyLeft  = IS_JOY_PRESSED(joy, JOY_INPUT_DIR_LEFT);
	joyRight = IS_JOY_PRESSED(joy, JOY_INPUT_DIR_RIGHT);

	// Se direções opostas forem acionadas simultaneamente, trata como porta aberta/aterrada
	if ((joyUp && joyDown) || (joyLeft && joyRight))
	{
		joyUp = joyDown = joyLeft = joyRight = FALSE;
	}

	// 2. Leitura de Teclado (Setas do MSX)
	kbUp    = Keyboard_IsKeyPressed(KEY_UP);
	kbDown  = Keyboard_IsKeyPressed(KEY_DOWN);
	kbLeft  = Keyboard_IsKeyPressed(KEY_LEFT);
	kbRight = Keyboard_IsKeyPressed(KEY_RIGHT);

	// Combinação de teclado e joystick
	moveUp    = kbUp || joyUp;
	moveDown  = kbDown || joyDown;
	moveLeft  = kbLeft || joyLeft;
	moveRight = kbRight || joyRight;

	// Se nenhuma direção for pressionada, reseta o cooldown imediatamente
	// garantindo resposta instantânea ao primeiro toque da tecla
	if (!moveUp && !moveDown && !moveLeft && !moveRight)
	{
		g_Hero.StepCooldown = 0;
		return;
	}

	// Se ainda estiver em cooldown de repetição ao segurar tecla, aguarda
	if (g_Hero.StepCooldown > 0)
	{
		return;
	}

	// Determina direção e delta de movimento
	if (moveUp)
	{
		dy = -1;
		g_Hero.Direction = HERO_DIR_UP;
	}
	else if (moveDown)
	{
		dy = 1;
		g_Hero.Direction = HERO_DIR_DOWN;
	}
	else if (moveLeft)
	{
		dx = -1;
		g_Hero.Direction = HERO_DIR_LEFT;
	}
	else if (moveRight)
	{
		dx = 1;
		g_Hero.Direction = HERO_DIR_RIGHT;
	}

	// Calcula coordenada de destino no grid
	{
		i8 targetX = (i8)g_Hero.TileX + dx;
		i8 targetY = (i8)g_Hero.TileY + dy;
		u8 newX, newY;

		// Checa se o destino ultrapassa as bordas da sala atual (32x18)
		if (targetX < 0 || targetX >= VIEWPORT_WIDTH || targetY < 0 || targetY >= VIEWPORT_HEIGHT)
		{
			// Tenta transição de sala pelas conexões cardeais
			if (WORLD_CheckRoomTransition(targetX, targetY, &newX, &newY))
			{
				g_Hero.TileX = newX;
				g_Hero.TileY = newY;
				g_Hero.PixelX = newX * 8;
				g_Hero.PixelY = newY * 8;
				g_Hero.StepCooldown = HERO_STEP_COOLDOWN_FRAMES + 2; // Folga suave na troca de tela
				g_Hero.Moved = TRUE;
				HERO_Draw();
			}
			else
			{
				// Borda sem conexão de sala: bloqueio
				g_Hero.StepCooldown = HERO_STEP_COOLDOWN_FRAMES;
			}
			return;
		}

		// Destino dentro da sala: checa colisão física do tile
		{
			u8 col = WORLD_GetCollision((u8)targetX, (u8)targetY);
			if (col == COLLISION_SOLID || col == COLLISION_WATER)
			{
				// Bloqueado por obstáculo sólido ou água
				g_Hero.StepCooldown = HERO_STEP_COOLDOWN_FRAMES;
				return;
			}

			// Terreno livre para travessia (Passável, Dano ou Gatilho)
			g_Hero.TileX = (u8)targetX;
			g_Hero.TileY = (u8)targetY;
			g_Hero.PixelX = g_Hero.TileX * 8;
			g_Hero.PixelY = g_Hero.TileY * 8;
			g_Hero.StepCooldown = HERO_STEP_COOLDOWN_FRAMES;
			g_Hero.Moved = TRUE;
			HERO_Draw();
		}
	}
}
