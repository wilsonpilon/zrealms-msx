// ____________________________
// Z-Realm (zrealm-msx) - Driver de Áudio PSG (AY-3-8910)
// Subfase 5.3: SFX Driver Nativo com Controle de Registradores
//─────────────────────────────────────────────────────────────────────────────
#include "audio.h"

// Portas I/O do PSG no padrão MSX (AY-3-8910 / YM2149)
__sfr __at(0xA0) s_PSG_Reg;
__sfr __at(0xA1) s_PSG_Data;

static void PSG_Write(u8 reg, u8 val)
{
	s_PSG_Reg = reg;
	s_PSG_Data = val;
}

static u8 s_SFXType = 0;
static u8 s_SFXStep = 0;
static u8 s_SFXTimer = 0;

void AUDIO_Init(void)
{
	// Silencia canais A, B, C
	PSG_Write(8, 0);
	PSG_Write(9, 0);
	PSG_Write(10, 0);

	// Mixer: Desativa todos os canais de tom e ruído (Bits 0..5 = 1), Port A input, Port B output
	PSG_Write(7, 0xBF);

	s_SFXType = 0;
	s_SFXStep = 0;
	s_SFXTimer = 0;
}

void AUDIO_PlaySFX(u8 sfxID)
{
	s_SFXType = sfxID;
	s_SFXStep = 0;
	s_SFXTimer = 0;

	// Aplica o primeiro quadro do efeito imediatamente
	AUDIO_Update();
}

void AUDIO_Update(void)
{
	if (s_SFXType == 0)
	{
		return;
	}

	if (s_SFXTimer > 0)
	{
		s_SFXTimer--;
		return;
	}

	switch (s_SFXType)
	{
		case 1: // Item / Chime / Vitória (Arpejo ascendente C5, E5, G5, C6)
		{
			if (s_SFXStep == 0)
			{
				PSG_Write(7, 0xBE); // Canal A Tom ON, ruído OFF
				PSG_Write(0, 0xFE); // C5 fine
				PSG_Write(1, 0x00); // C5 coarse
				PSG_Write(8, 14);   // Vol A
				s_SFXTimer = 4;
				s_SFXStep = 1;
			}
			else if (s_SFXStep == 1)
			{
				PSG_Write(0, 0xCA); // E5 fine
				PSG_Write(1, 0x00);
				PSG_Write(8, 14);
				s_SFXTimer = 4;
				s_SFXStep = 2;
			}
			else if (s_SFXStep == 2)
			{
				PSG_Write(0, 0xAA); // G5 fine
				PSG_Write(1, 0x00);
				PSG_Write(8, 14);
				s_SFXTimer = 4;
				s_SFXStep = 3;
			}
			else if (s_SFXStep == 3)
			{
				PSG_Write(0, 0x7F); // C6 fine
				PSG_Write(1, 0x00);
				PSG_Write(8, 15);
				s_SFXTimer = 10;
				s_SFXStep = 4;
			}
			else if (s_SFXStep == 4)
			{
				PSG_Write(8, 8); // Decay
				s_SFXTimer = 6;
				s_SFXStep = 5;
			}
			else
			{
				// Fim do efeito
				PSG_Write(8, 0);
				PSG_Write(7, 0xBF);
				s_SFXType = 0;
			}
			break;
		}

		case 2: // Dano / Combate / Perigo (Ruído + Tom grave com decaimento rápido)
		{
			if (s_SFXStep == 0)
			{
				PSG_Write(7, 0xB6); // Canal A Tom ON + Ruído ON
				PSG_Write(6, 0x18); // Ruído grosso
				PSG_Write(0, 0x57); // Tom grave fine (0x0357)
				PSG_Write(1, 0x03);
				PSG_Write(8, 15);   // Vol máx
				s_SFXTimer = 3;
				s_SFXStep = 1;
			}
			else if (s_SFXStep == 1)
			{
				PSG_Write(6, 0x10);
				PSG_Write(0, 0x00);
				PSG_Write(1, 0x04);
				PSG_Write(8, 12);
				s_SFXTimer = 3;
				s_SFXStep = 2;
			}
			else if (s_SFXStep == 2)
			{
				PSG_Write(6, 0x08);
				PSG_Write(0, 0x00);
				PSG_Write(1, 0x05);
				PSG_Write(8, 7);
				s_SFXTimer = 4;
				s_SFXStep = 3;
			}
			else if (s_SFXStep == 3)
			{
				PSG_Write(6, 0x04);
				PSG_Write(8, 3);
				s_SFXTimer = 3;
				s_SFXStep = 4;
			}
			else
			{
				PSG_Write(8, 0);
				PSG_Write(7, 0xBF);
				s_SFXType = 0;
			}
			break;
		}

		case 3: // Porta / Teletransporte (Duplo tom com intervalo mágico)
		{
			if (s_SFXStep == 0)
			{
				PSG_Write(7, 0xBC); // Canais A e B Tom ON
				PSG_Write(0, 0x50); // Tom A (0x0150)
				PSG_Write(1, 0x01);
				PSG_Write(2, 0x00); // Tom B (0x0200)
				PSG_Write(3, 0x02);
				PSG_Write(8, 13);
				PSG_Write(9, 13);
				s_SFXTimer = 5;
				s_SFXStep = 1;
			}
			else if (s_SFXStep == 1)
			{
				PSG_Write(0, 0xA0); // Tom A (0x01A0)
				PSG_Write(1, 0x01);
				PSG_Write(2, 0x80); // Tom B (0x0280)
				PSG_Write(3, 0x02);
				PSG_Write(8, 10);
				PSG_Write(9, 10);
				s_SFXTimer = 5;
				s_SFXStep = 2;
			}
			else if (s_SFXStep == 2)
			{
				PSG_Write(0, 0x20); // Tom A (0x0220)
				PSG_Write(1, 0x02);
				PSG_Write(2, 0x20); // Tom B (0x0320)
				PSG_Write(3, 0x03);
				PSG_Write(8, 6);
				PSG_Write(9, 6);
				s_SFXTimer = 6;
				s_SFXStep = 3;
			}
			else
			{
				PSG_Write(8, 0);
				PSG_Write(9, 0);
				PSG_Write(7, 0xBF);
				s_SFXType = 0;
			}
			break;
		}

		case 4: // Diálogo / Blip (Tom agudo curto)
		{
			if (s_SFXStep == 0)
			{
				PSG_Write(7, 0xBE); // Canal A Tom ON
				PSG_Write(0, 0xB4); // Tom agudo
				PSG_Write(1, 0x00);
				PSG_Write(8, 12);
				s_SFXTimer = 2;
				s_SFXStep = 1;
			}
			else if (s_SFXStep == 1)
			{
				PSG_Write(8, 6);
				s_SFXTimer = 2;
				s_SFXStep = 2;
			}
			else
			{
				PSG_Write(8, 0);
				PSG_Write(7, 0xBF);
				s_SFXType = 0;
			}
			break;
		}

		default:
			AUDIO_Init();
			break;
	}
}
