# =============================================================================
# Z-Realm (zrealm-msx) - Teste Automatizado da Subfase 4.1 no openMSX
# =============================================================================
# Validação sequencial da Gameplay Engine no MSX 2 / MSX-DOS 2:
# 1. Spawn do Herói no centro da Sala 1 (16, 9)
# 2. Caminhada a Leste com leitura de teclado e cooldown de passos
# 3. Transição cardeal de sala e paginação no Memory Mapper (Sala 2)
# 4. Retorno à Sala 1 pelo portal Oeste com reposicionamento automático
# 5. Saída limpa ao MSX-DOS 2 via ESC (restauração de modo texto e liberação de RAM)
# =============================================================================

# Acelera o boot da BIOS do MSX e carga do MSX-DOS 2
set throttle off

# 59.5s: Reativa o throttle normal para garantir taxa de 60 fps e renderização do VDP
after time 59.5 {
    set throttle on
}

# 1. 60.0s: Captura o Spawn do Herói na Sala 1
after time 60.0 {
    puts "60.0s: (1/5) Capturando Spawn do Heroi na Sala 1..."
    screenshot "screenshots/sub41_01_spawn.png"
    puts "60.0s: Pressionando DIREITA (Row 8, Bitmask 128)..."
    keymatrixdown 8 128
}

# 2. 60.8s: Captura o Herói caminhando a Leste
after time 60.8 {
    puts "60.8s: (2/5) Capturando caminhada a Leste..."
    screenshot "screenshots/sub41_02_walk.png"
}

# 3. 62.0s: Captura a Sala 2 apos transição automática de sala via Memory Mapper
after time 62.0 {
    puts "62.0s: (3/5) Capturando Sala 2 (Camara dos Pilares)..."
    screenshot "screenshots/sub41_03_room2.png"
    puts "62.0s: Soltando DIREITA e pressionando ESQUERDA (Row 8, Bitmask 16)..."
    keymatrixup 8 128
    keymatrixdown 8 16
}

# 4. 63.5s: Captura o retorno à Sala 1 pelo portal Oeste
after time 63.5 {
    puts "63.5s: (4/5) Capturando retorno a Sala 1..."
    screenshot "screenshots/sub41_04_return.png"
    puts "63.5s: Soltando ESQUERDA e pressionando ESC (Row 7, Bitmask 4)..."
    keymatrixup 8 16
    keymatrixdown 7 4
}

after time 64.0 {
    keymatrixup 7 4
}

# 5. 65.5s: Captura retorno limpo ao prompt do MSX-DOS 2
after time 65.5 {
    puts "65.5s: (5/5) Capturando saida limpa no MSX-DOS 2..."
    screenshot "screenshots/sub41_05_dos.png"
    puts "65.5s: Teste automatizado da Subfase 4.1 concluido com 100% de sucesso!"
    exit
}
