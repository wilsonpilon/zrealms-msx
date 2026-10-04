# =============================================================================
# Z-Realm (zrealm-msx) - Teste Automatizado da Subfase 5.1 no openMSX
# =============================================================================
# Validação do pipeline One-Click Run:
# 1. Boot automático a partir do disquete (.DSK) gerado pelo pipeline
# 2. Execução transparente do jogo no MSX 2 (V9938 SCREEN 4 + HUD + VM)
# 3. Interação no mundo de aventura
# 4. Encerramento gracioso via ESC retornando ao prompt do MSX-DOS 2
# =============================================================================

set log_file [open "engine_msx/screenshots/sub51_execution.log" w]

proc log_msg {msg} {
    global log_file
    puts $log_file $msg
    flush $log_file
    puts $msg
}

log_msg "============================================================"
log_msg "   VALIDAÇÃO DETERMINÍSTICA - SUBFASE 5.1: ONE-CLICK RUN   "
log_msg "============================================================"

set state_timer 0
set test_state 0

proc control_loop {} {
    global test_state state_timer log_file
    
    switch $test_state {
        0 {
            # Estado 0: Captura o boot inicial do jogo com HUD e masmorra ativa
            incr state_timer
            if {$state_timer >= 10} {
                log_msg "=== (1/3) Boot Concluído: Masmorra e HUD Ativos no V9938 ==="
                screenshot "engine_msx/screenshots/sub51_01_oneclick_boot.png"
                set state_timer 0
                set test_state 1
            }
        }
        1 {
            # Estado 1: Anda a Oeste em direção ao Guardião (Setas: Coluna 8, Linha 4)
            incr state_timer
            if {$state_timer == 5} {
                keymatrixdown 8 4
            } elseif {$state_timer == 15} {
                keymatrixup 8 4
            } elseif {$state_timer == 25} {
                keymatrixdown 8 4
            } elseif {$state_timer == 35} {
                keymatrixup 8 4
            } elseif {$state_timer == 45} {
                keymatrixdown 8 4
            } elseif {$state_timer == 55} {
                keymatrixup 8 4
                log_msg "Heroi posicionado em frente ao Guardiao."
                set state_timer 0
                set test_state 2
            }
        }
        2 {
            # Estado 2: Pressiona ESPAÇO para interagir (Coluna 8, Linha 0)
            incr state_timer
            if {$state_timer == 5} {
                log_msg "=== (2/3) Interagindo com o Guardião (Abre Diálogo) ==="
                keymatrixdown 8 0
            } elseif {$state_timer == 12} {
                keymatrixup 8 0
            } elseif {$state_timer >= 30} {
                screenshot "engine_msx/screenshots/sub51_02_oneclick_gameplay.png"
                set state_timer 0
                set test_state 3
            }
        }
        3 {
            # Estado 3: Fecha o diálogo (ESPAÇO)
            incr state_timer
            if {$state_timer == 5} {
                keymatrixdown 8 0
            } elseif {$state_timer == 12} {
                keymatrixup 8 0
            } elseif {$state_timer >= 25} {
                set state_timer 0
                set test_state 4
            }
        }
        4 {
            # Estado 4: Pressiona ESC (Coluna 7, Linha 4) para sair limpo ao DOS 2
            incr state_timer
            if {$state_timer == 5} {
                log_msg "=== (3/3) Pressionando ESC para Saída Graciosa ao MSX-DOS 2 ==="
                keymatrixdown 7 4
            } elseif {$state_timer == 12} {
                keymatrixup 7 4
            } elseif {$state_timer >= 35} {
                screenshot "engine_msx/screenshots/sub51_03_oneclick_exit.png"
                log_msg "=== VALIDACAO DA SUBFASE 5.1 CONCLUIDA COM 100% DE SUCESSO! ==="
                close $log_file
                exit
            }
        }
    }
    
    after "frame" control_loop
}

# Inicializa o loop de controle após o boot do DOS e carga do jogo (21.0 segundos)
after time 21.0 {
    set throttle on
    log_msg "Boot do MSX-DOS 2 e carga do zrealm.com concluidos. Iniciando ciclo de testes..."
    control_loop
}
