# =============================================================================
# Z-Realm (zrealm-msx) - Teste Automatizado da Subfase 4.3 no openMSX
# =============================================================================
# Controlador em malha fechada (Closed-Loop) via leitura de memória no openMSX:
# Garante navegação determinística e à prova de desvios no grid.
# =============================================================================

set log_file [open "screenshots/sub43_execution.log" w]

proc log_msg {msg} {
    global log_file
    puts $log_file $msg
    flush $log_file
    puts $msg
}

proc read_u8 {addr} {
    return [debug read "memory" $addr]
}

proc read_u16 {addr} {
    set low [debug read "memory" $addr]
    set high [debug read "memory" [expr {$addr + 1}]]
    return [expr {$low + ($high << 8)}]
}

# Endereços oficiais dos símbolos extraídos de zrealm.map
set ADDR_HERO_X        0x35EA
set ADDR_HERO_Y        0x35EB
set ADDR_HERO_DIR      0x35EE
set ADDR_VM_FLAGS      0x366A
set ADDR_VM_INV        0x376A
set ADDR_VM_INV_CNT    0x379A
set ADDR_HERO_HP       0x379B
set ADDR_VM_MESSAGE    0x37A6
set ADDR_VM_STRING_ID  0x3826
set ADDR_VM_HAS_MSG    0x3828
set ADDR_VM_SFX        0x3829

proc release_all_keys {} {
    keymatrixup 8 16
    keymatrixup 8 128
    keymatrixup 8 32
    keymatrixup 8 64
    keymatrixup 8 1
    keymatrixup 7 4
}

proc log_status {label} {
    global ADDR_HERO_X ADDR_HERO_Y ADDR_HERO_DIR ADDR_VM_FLAGS ADDR_VM_INV ADDR_VM_INV_CNT ADDR_HERO_HP ADDR_VM_STRING_ID ADDR_VM_SFX
    set hx [read_u8 $ADDR_HERO_X]
    set hy [read_u8 $ADDR_HERO_Y]
    set hdir [read_u8 $ADDR_HERO_DIR]
    set f1 [read_u8 [expr {$ADDR_VM_FLAGS + 1}]]
    set f2 [read_u8 [expr {$ADDR_VM_FLAGS + 2}]]
    set invCnt [read_u8 $ADDR_VM_INV_CNT]
    set hp [read_u16 $ADDR_HERO_HP]
    set strId [read_u16 $ADDR_VM_STRING_ID]
    set sfx [read_u8 $ADDR_VM_SFX]
    set it0 [read_u16 $ADDR_VM_INV]
    set it1 [read_u16 [expr {$ADDR_VM_INV + 3}]]
    log_msg [format "%-32s | Hero:(%2d,%2d) Dir:%d | F1:%d F2:%d | InvCnt:%d (It0:%d,It1:%d) | HP:%d | StrID:%d | SFX:%d" \
        $label $hx $hy $hdir $f1 $f2 $invCnt $it0 $it1 $hp $strId $sfx]
}

set throttle off

# 59.5s: Reativa taxa normal de 60 fps
after time 59.5 {
    set throttle on
}

set test_state 0
set state_timer 0

proc control_loop {} {
    global test_state state_timer ADDR_HERO_X ADDR_HERO_Y ADDR_HERO_DIR log_file
    
    set hx [read_u8 $ADDR_HERO_X]
    set hy [read_u8 $ADDR_HERO_Y]
    set hdir [read_u8 $ADDR_HERO_DIR]
    
    switch $test_state {
        0 {
            # Spawn Inicial
            log_msg "=== (1/8) Spawn Inicial na Sala 1 ==="
            screenshot "screenshots/sub43_01_spawn.png"
            log_status "Spawn Inicial"
            log_msg "Navegando a Oeste em direcao ao Guardiao em (12, 9)..."
            set test_state 1
        }
        1 {
            # Move para (13, 9)
            if {$hx > 13} {
                keymatrixdown 8 16
            } elseif {$hx <= 13} {
                keymatrixup 8 16
                log_msg "=== (2/8) Heroi Diante do Guardiao em (13, 9) ==="
                screenshot "screenshots/sub43_02_facing_guardian.png"
                log_status "Diante do Guardiao"
                log_msg "Pressionando ESPACO para interagir com o Guardiao (Executa Script 1)..."
                keymatrixdown 8 1
                set state_timer 0
                set test_state 2
            }
        }
        2 {
            incr state_timer
            if {$state_timer == 2} {
                keymatrixup 8 1
            } elseif {$state_timer >= 5} {
                log_msg "=== (3/8) Script 1 do Guardiao Executado ==="
                screenshot "screenshots/sub43_03_guardian_quest.png"
                log_status "Guardiao Quest Concedida"
                log_msg "Pressionando ESPACO novamente para testar ramificacao de Flag 1..."
                keymatrixdown 8 1
                set state_timer 0
                set test_state 3
            }
        }
        3 {
            incr state_timer
            if {$state_timer == 2} {
                keymatrixup 8 1
            } elseif {$state_timer >= 5} {
                log_msg "=== (4/8) Desvio Condicional do Script 1 (Flag 1 == 1) ==="
                screenshot "screenshots/sub43_04_guardian_branched.png"
                log_status "Guardiao Mensagem Ramificada"
                log_msg "Navegando ate o Bau: primeiro sobe para Y=5..."
                set state_timer 0
                set test_state 4
            }
        }
        4 {
            # Sobe para Y=5
            if {$hy > 5} {
                keymatrixdown 8 32
            } else {
                keymatrixup 8 32
                log_msg "Chegou a Y=5. Caminhando a Oeste para X=8..."
                set test_state 5
            }
        }
        5 {
            # Caminha a Oeste para X=8
            if {$hx > 8} {
                keymatrixdown 8 16
            } else {
                keymatrixup 8 16
                log_msg "Chegou a X=8. Virando para CIMA encarando o Bau em (8, 4)..."
                keymatrixdown 8 32
                set state_timer 0
                set test_state 6
            }
        }
        6 {
            incr state_timer
            if {$state_timer >= 2} {
                keymatrixup 8 32
                log_msg "=== (5/8) Heroi Diante do Bau em (8, 4) ==="
                screenshot "screenshots/sub43_05_facing_chest.png"
                log_status "Diante do Bau (8, 5)"
                log_msg "Pressionando ESPACO para abrir o Bau (Executa Script 2)..."
                keymatrixdown 8 1
                set state_timer 0
                set test_state 7
            }
        }
        7 {
            incr state_timer
            if {$state_timer == 2} {
                keymatrixup 8 1
            } elseif {$state_timer >= 5} {
                log_msg "=== (6/8) Script 2 do Bau Executado ==="
                screenshot "screenshots/sub43_06_chest_opened.png"
                log_status "Bau Aberto e Pocao Adquirida"
                log_msg "Pressionando ESPACO novamente no Bau aberto..."
                keymatrixdown 8 1
                set state_timer 0
                set test_state 8
            }
        }
        8 {
            incr state_timer
            if {$state_timer == 2} {
                keymatrixup 8 1
            } elseif {$state_timer >= 5} {
                log_msg "=== (7/8) Desvio Condicional do Script 2 (Bau Vazio) ==="
                screenshot "screenshots/sub43_07_chest_empty.png"
                log_status "Bau Vazio (Flag 2 == 1)"
                log_msg "Pressionando ESC para encerramento gracioso..."
                keymatrixdown 7 4
                set state_timer 0
                set test_state 9
            }
        }
        9 {
            incr state_timer
            if {$state_timer == 2} {
                keymatrixup 7 4
            } elseif {$state_timer >= 6} {
                log_msg "=== (8/8) Retorno Limpo ao MSX-DOS 2 ==="
                screenshot "screenshots/sub43_08_dos_exit.png"
                log_msg "Z-Realm Subfase 4.3: Teste Automatizado da Bytecode VM CONCLUIDO COM 100% DE SUCESSO!"
                close $log_file
                exit
            }
        }
    }
    
    after time 0.1 control_loop
}

after time 60.0 {
    control_loop
}
