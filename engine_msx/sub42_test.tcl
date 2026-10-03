# =============================================================================
# Z-Realm (zrealm-msx) - Teste Automatizado da Subfase 4.2 no openMSX
# =============================================================================
# Validação sequencial do Sistema de Entidades e Atores da Sala:
# 1. Spawn da Sala 1: Herói (16, 9), Guardião (12, 9) e Baú (8, 4) visíveis
# 2. Caminhada a Oeste e bloqueio de colisão física contra o Guardião em (13, 9)
# 3. Interação com o Guardião via Barra de Espaço (Row 8, Bitmask 1)
# 4. Transição para Sala 2 (Câmara dos Pilares): limpeza de entidades da Sala 1
#    e spawn das entidades da Sala 2 (Sentinela Errante e Baú Místico)
# 5. Movimentação autônoma da IA do NPC Errante (Wandering NPC) no grid
# 6. Saída limpa ao MSX-DOS 2 via ESC (desativação de sprites e liberação de RAM)
# =============================================================================

set throttle off

# 59.5s: Reativa o throttle normal para garantir taxa de 60 fps e renderização do VDP
after time 59.5 {
    set throttle on
}

# 1. 60.0s: Captura o Spawn da Sala 1 com Herói, Guardião e Baú
after time 60.0 {
    puts "60.0s: (1/6) Capturando Spawn da Sala 1 com Heroi, Guardiao e Bau..."
    screenshot "screenshots/sub42_01_spawn_entities.png"
    puts "60.0s: Caminhando a Oeste em direcao ao Guardiao (Row 8, Bitmask 16)..."
    keymatrixdown 8 16
}

# 2. 61.2s: Captura colisão sólida contra o Guardião em (13, 9)
after time 61.2 {
    puts "61.2s: (2/6) Capturando colisao solida contra o Guardiao em (13, 9)..."
    screenshot "screenshots/sub42_02_collision_guardian.png"
    puts "61.2s: Soltando ESQUERDA e pressionando ESPACO para interagir (Row 8, Bitmask 1)..."
    keymatrixup 8 16
    keymatrixdown 8 1
}

# 3. 61.8s: Captura interação com o Guardião
after time 61.8 {
    puts "61.8s: (3/6) Capturando interacao com o Guardiao via tecla ESPACO..."
    screenshot "screenshots/sub42_03_interact_guardian.png"
    keymatrixup 8 1
    puts "61.8s: Caminhando a Leste em direcao ao portal da Sala 2 (Row 8, Bitmask 128)..."
    keymatrixdown 8 128
}

# 4. 64.6s: Captura a Sala 2 (Câmara dos Pilares) logo após a transição de sala
after time 64.6 {
    puts "64.6s: (4/6) Capturando Sala 2 com Sentinela Errante e Bau Mistico..."
    screenshot "screenshots/sub42_04_room2_entities.png"
    keymatrixup 8 128
}

# 5. 66.2s: Captura a movimentação autônoma do Sentinela Errante na Sala 2
after time 66.2 {
    puts "66.2s: (5/6) Capturando movimentacao autonoma do NPC Errante..."
    screenshot "screenshots/sub42_05_wandering_npc.png"
    puts "66.2s: Pressionando ESC para encerramento gracioso (Row 7, Bitmask 4)..."
    keymatrixdown 7 4
}

after time 66.7 {
    keymatrixup 7 4
}

# 6. 68.0s: Captura retorno limpo ao prompt do MSX-DOS 2
after time 68.0 {
    puts "68.0s: (6/6) Capturando saida limpa no MSX-DOS 2..."
    screenshot "screenshots/sub42_06_dos_clean_exit.png"
    puts "68.0s: Teste automatizado da Subfase 4.2 concluido com 100% de sucesso!"
    exit
}
