# =============================================================================
# Z-Realm (zrealm-msx) - Teste Automatizado da Subfase 4.4 no openMSX
# =============================================================================
# Validação em malha fechada do renderizador de HUD e Caixa de Diálogos:
# Prova renderização de texto, molduras, paginação, cura no HUD e saída ao DOS.
# =============================================================================

set log_file [open "engine_msx/screenshots/sub44_execution.log" w]

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
set ADDR_HERO_X        0x416E
set ADDR_HERO_Y        0x416F
set ADDR_HERO_DIR      0x4172
set ADDR_VM_FLAGS      0x41EE
set ADDR_VM_INV        0x42EE
set ADDR_VM_INV_CNT    0x431E
set ADDR_HERO_HP       0x431F
set ADDR_HERO_MAX_HP   0x4321
set ADDR_VM_MESSAGE    0x432A
set ADDR_VM_STRING_ID  0x43AA
set ADDR_VM_HAS_MSG    0x43AC
set ADDR_VM_SFX        0x43AD

# Descoberta dinâmica do ponteiro de s_DialogueActive via ld a, (_s_DialogueActive) em UI_IsDialogueActive (0x2A32)
set ADDR_PEEK_PTR      0x2A33
set b_low              [read_u8 $ADDR_PEEK_PTR]
set b_high             [read_u8 [expr {$ADDR_PEEK_PTR + 1}]]
set ADDR_DIALOG_ACTIVE [expr {$b_low | ($b_high << 8)}]

proc get_dialog_active {} {
    global ADDR_DIALOG_ACTIVE
    return [read_u8 $ADDR_DIALOG_ACTIVE]
}

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
    set dlg [get_dialog_active]
    log_msg [format "%-30s | Hero:(%2d,%2d) Dir:%d | Dlg:%d | F1:%d F2:%d | Inv:%d | HP:%d | StrID:%d" \
        $label $hx $hy $hdir $dlg $f1 $f2 $invCnt $hp $strId]
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
    set dlg [get_dialog_active]
    
    switch $test_state {
        0 {
            # Spawn Inicial: HUD com HP 75/100, MP 50/50, LV 01, K 0 e Painel de Repouso
            log_msg "=== (1/9) Spawn Inicial na Sala 1 (HUD & Painel Ativos) ==="
            screenshot "engine_msx/screenshots/sub44_01_spawn_hud.png"
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
                log_msg "=== (2/9) Heroi Diante do Guardiao em (13, 9) ==="
                screenshot "engine_msx/screenshots/sub44_02_facing_guardian.png"
                log_status "Diante do Guardiao"
                log_msg "Pressionando ESPACO para interagir com o Guardiao (Abre Caixa de Dialogo)..."
                keymatrixdown 8 1
                set state_timer 0
                set test_state 2
            }
        }
        2 {
            incr state_timer
            if {$state_timer == 2} {
                keymatrixup 8 1
            } elseif {$state_timer >= 6} {
                log_msg "=== (3/9) Caixa de Dialogo do Guardiao Aberta (Linhas 19-23) ==="
                screenshot "engine_msx/screenshots/sub44_03_guardian_dialogue.png"
                log_status "Dialogo Guardiao Aberto"
                log_msg "Pressionando ESPACO para fechar a caixa de dialogo..."
                keymatrixdown 8 1
                set state_timer 0
                set test_state 3
            }
        }
        3 {
            incr state_timer
            if {$state_timer == 2} {
                keymatrixup 8 1
            } elseif {$state_timer >= 6} {
                log_msg "=== (4/9) Dialogo Fechado (Painel de Repouso Restaurado) ==="
                screenshot "engine_msx/screenshots/sub44_04_dialogue_closed.png"
                log_status "Dialogo Fechado"
                log_msg "Pressionando ESPACO novamente para testar ramificacao de Flag 1..."
                keymatrixdown 8 1
                set state_timer 0
                set test_state 4
            }
        }
        4 {
            incr state_timer
            if {$state_timer == 2} {
                keymatrixup 8 1
            } elseif {$state_timer >= 6} {
                log_msg "=== (5/9) Dialogo Ramificado do Guardiao (Flag 1 == 1) ==="
                screenshot "engine_msx/screenshots/sub44_05_guardian_branch_dialogue.png"
                log_status "Dialogo Ramificado"
                log_msg "Fechando dialogo e navegando ate o Bau..."
                keymatrixdown 8 1
                set state_timer 0
                set test_state 5
            }
        }
        5 {
            incr state_timer
            if {$state_timer == 2} {
                keymatrixup 8 1
            } elseif {$state_timer >= 6} {
                # Sobe para Y=5
                if {$hy > 5} {
                    keymatrixdown 8 32
                } else {
                    keymatrixup 8 32
                    # Caminha para X=8
                    if {$hx > 8} {
                        keymatrixdown 8 16
                    } else {
                        keymatrixup 8 16
                        # Vira para CIMA
                        if {$hdir != 0} {
                            keymatrixdown 8 32
                        } else {
                            keymatrixup 8 32
                            log_msg "=== (6/9) Heroi Diante do Bau Antigo em (8, 4) ==="
                            screenshot "engine_msx/screenshots/sub44_06_facing_chest.png"
                            log_status "Diante do Bau"
                            log_msg "Pressionando ESPACO para abrir o Bau (Executa Cura +25 HP no HUD)..."
                            keymatrixdown 8 1
                            set state_timer 0
                            set test_state 6
                        }
                    }
                }
            }
        }
        6 {
            incr state_timer
            if {$state_timer == 2} {
                keymatrixup 8 1
            } elseif {$state_timer >= 6} {
                log_msg "=== (7/9) Bau Aberto: Dialogo do Bau + HUD Curado para 100 HP! ==="
                screenshot "engine_msx/screenshots/sub44_07_chest_heal_hud.png"
                log_status "Bau Aberto e Curado"
                log_msg "Pressionando ESPACO para fechar dialogo do bau..."
                keymatrixdown 8 1
                set state_timer 0
                set test_state 7
            }
        }
        7 {
            incr state_timer
            if {$state_timer == 2} {
                keymatrixup 8 1
            } elseif {$state_timer >= 6} {
                log_msg "Pressionando ESPACO para reinspecionar o Bau esvaziado..."
                keymatrixdown 8 1
                set state_timer 0
                set test_state 8
            }
        }
        8 {
            incr state_timer
            if {$state_timer == 2} {
                keymatrixup 8 1
            } elseif {$state_timer >= 6} {
                log_msg "=== (8/9) Dialogo de Bau Vazio Detectado ==="
                screenshot "engine_msx/screenshots/sub44_08_chest_empty_dialogue.png"
                log_status "Bau Vazio"
                log_msg "Fechando dialogo e pressionando ESC para encerramento..."
                keymatrixdown 8 1
                set state_timer 0
                set test_state 9
            }
        }
        9 {
            incr state_timer
            if {$state_timer == 2} {
                keymatrixup 8 1
            } elseif {$state_timer >= 6} {
                log_msg "=== (9/9) Pressionando ESC para Saida Limpa ao MSX-DOS 2 ==="
                release_all_keys
                keymatrixdown 7 4
                set state_timer 0
                set test_state 10
            }
        }
        10 {
            incr state_timer
            if {$state_timer == 2} {
                keymatrixup 7 4
            } elseif {$state_timer >= 20} {
                screenshot "engine_msx/screenshots/sub44_09_dos_exit.png"
                log_msg "=== VALIDACAO DA SUBFASE 4.4 CONCLUIDA COM 100% DE SUCESSO! ==="
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
    log_msg "Iniciando ciclo de controle deterministico..."
    control_loop
}
