// ____________________________
// Z-Realm (zrealm-msx) - Driver de Áudio PSG (AY-3-8910)
// Subfase 5.3: SFX Driver Nativo com Controle de Registradores
//─────────────────────────────────────────────────────────────────────────────
#pragma once

#include "core.h"

// Inicializa o subsistema de áudio PSG e zera os canais
void AUDIO_Init(void);

// Dispara um efeito sonoro (SFX)
// 1: Chime / Item coletado / Vitória
// 2: Dano / Combate / Perigo
// 3: Porta / Teletransporte
// 4: Diálogo / Blip
void AUDIO_PlaySFX(u8 sfxID);

// Atualiza o estado dos envelopes de áudio e decay de volume a cada quadro V-Blank (50/60 Hz)
void AUDIO_Update(void);
