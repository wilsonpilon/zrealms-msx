set throttle off

# 1. 45.0s: SCREEN 4 ativa exibindo a Sala 0 (Segmento 0 na Página 2, borda azul de tijolos)
after time 45.0 {
    puts "45.0s: Capturando Room 0..."
    screenshot "screenshots/zrealm_room0.png"
}

# 2. 53.0s: SCREEN 4 com Sala 1 (Segmento 1 na Página 2, cruz central dividindo a sala)
after time 53.0 {
    puts "53.0s: Capturando Room 1 (apos troca automatica para Segmento 1)..."
    screenshot "screenshots/zrealm_room1.png"
}

# 3. 64.0s: Captura retorno ao prompt do MSX-DOS 2 (TEXT1 restaurado, RAM 100% liberada)
after time 64.0 {
    puts "64.0s: Capturando saida limpa no MSX-DOS 2..."
    screenshot "screenshots/zrealm_exit.png"
}

# 4. 66s: Finaliza teste
after time 66 {
    puts "66s: Teste concluido com sucesso!"
    exit
}
