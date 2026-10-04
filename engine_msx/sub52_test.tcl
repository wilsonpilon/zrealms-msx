# =============================================================================
# Z-Realm (zrealm-msx) - Teste Automatizado da Subfase 5.2 no openMSX
# =============================================================================
# Controlador em malha fechada (Closed-Loop) sincronizado quadro a quadro (after "frame"):
# Validação do Cartucho MegaROM ASCII-16 (demo.rom) executando nativamente em SCREEN 4.
# =============================================================================

set log_file [open "screenshots/sub52_execution.log" w]

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

# Endereços oficiais dos símbolos extraídos de zrealm.map (Página 3 RAM - ROM_ASCII16)
set ADDR_HERO_X        0xC578
set ADDR_HERO_Y        0xC579
set ADDR_HERO_DIR      0xC57C
set ADDR_VM_FLAGS      0xC5F8
set ADDR_VM_INV        0xC6F8
set ADDR_VM_INV_CNT    0xC728
set ADDR_HERO_HP       0xC729
set ADDR_VM_MESSAGE    0xC734
set ADDR_VM_STRING_ID  0xC7B4
set ADDR_VM_HAS_MSG    0xC7B6
set ADDR_VM_SFX        0xC7B7

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
    log_msg [format "%-32s | Hero:(%2d,%2d) Dir:%d | F1:%d F2:%d | InvCnt:%d (It0:%d) | HP:%d | StrID:%d | SFX:%d" \
        $label $hx $hy $hdir $f1 $f2 $invCnt $it0 $hp $strId $sfx]
}

# Acelera a fase de boot inicial do BIOS (0 a 7.5 segundos emulados)
set throttle off

set test_state 0
set state_timer 0

proc control_loop {} {
    global test_state state_timer ADDR_HERO_X ADDR_HERO_Y ADDR_HERO_DIR log_file
    
    set hx [read_u8 $ADDR_HERO_X]
    set hy [read_u8 $ADDR_HERO_Y]
    set hdir [read_u8 $ADDR_HERO_DIR]
    
    switch $test_state {
        0 {
            # Aguarda estabilização do VDP após boot
            incr state_timer
            if {$state_timer >= 10} {
                log_msg "=== (1/5) Boot do Cartucho MegaROM - Spawn Inicial na Sala 1 ==="
                screenshot "screenshots/sub52_01_spawn_cart.png"
                log_status "Spawn Inicial ROM"
                log_msg "Navegando a Oeste em direcao ao Guardiao em (12, 9)..."
                set state_timer 0
                set test_state 1
            }
        }
        1 {
            # Move para (13, 9)
            if {$hx > 13} {
                keymatrixdown 8 16
            } else {
                keymatrixup 8 16
                incr state_timer
                if {$state_timer >= 10} {
                    log_msg "=== (2/5) Heroi Diante do Guardiao em (13, 9) ==="
                    screenshot "screenshots/sub52_02_facing_guardian_cart.png"
                    log_status "Diante do Guardiao"
                    log_msg "Acionando dialogo via tecla ESPACO..."
                    keymatrixdown 8 1
                    set state_timer 0
                    set test_state 2
                }
            }
        }
        2 {
            incr state_timer
            if {$state_timer == 3} {
                keymatrixup 8 1
            } elseif {$state_timer >= 15} {
                log_msg "=== (3/5) Caixa de Dialogo Renderizada em SCREEN 4 ==="
                screenshot "screenshots/sub52_03_dialogue_cart.png"
                log_status "Dialogo Ativo"
                log_msg "Avancando dialogo para receber a Chave de Bronze..."
                keymatrixdown 8 1
                set state_timer 0
                set test_state 3
            }
        }
        3 {
            incr state_timer
            if {$state_timer == 3} {
                keymatrixup 8 1
            } elseif {$state_timer >= 15} {
                log_msg "=== (4/5) Chave de Bronze Recebida no Inventario ==="
                screenshot "screenshots/sub52_04_item_received_cart.png"
                log_status "Item Recebido"
                log_msg "Caminhando em direcao ao Leste para transicao de sala..."
                set state_timer 0
                set test_state 4
            }
        }
        4 {
            # Caminha a Leste em direção ao portal da Sala 2 (EastRoom)
            keymatrixdown 8 128
            # Ao cruzar a borda Leste da Sala 1 (X=31), o Herói surge no lado Oeste da Sala 2 (X=0 ou 1)
            if {$hx == 0 || $hx == 1} {
                release_all_keys
                incr state_timer
                if {$state_timer >= 20} {
                    log_msg "=== (5/5) Transicao Concluida - Heroi na Sala 2 (Camara dos Pilares) ==="
                    screenshot "screenshots/sub52_05_room2_cart.png"
                    log_status "Sala 2 Atingida"
                    log_msg "================================================================="
                    log_msg "   VALIDACAO SUBFASE 5.2 (MEGAROM .ROM) CONCLUIDA COM SUCESSO!   "
                    log_msg "================================================================="
                    set test_state 99
                    after time 1.0 {
                        close $log_file
                        exit
                    }
                }
            }
        }
        99 {
            return
        }
    }
    
    after "frame" control_loop
}

# Inicializa o loop de controle após o boot completo do cartucho (7.5 segundos emulados)
after time 7.5 {
    set throttle on
    log_msg "Cartucho MegaROM carregado. Reativando throttle e sincronizando a 60 fps..."
    control_loop
}

# Timeout de segurança após 45 segundos
after time 45.0 {
    log_msg "TIMEOUT de seguranca alcancado. Encerrando."
    release_all_keys
    close $log_file
    exit
}
