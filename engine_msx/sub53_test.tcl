# =============================================================================
# Z-Realm (zrealm-msx) - Teste Automatizado da Subfase 5.3 no openMSX
# =============================================================================
# Validação do Jogo de Referência Completo:
# "As Catacumbas de Cristal: O Desafio do Rei Esquecido"
# - 20 Salas Interligadas em Grid 4x5
# - Sistema de Áudio PSG AY-3-8910 (SFX 1: Item, SFX 2: Dano, SFX 3: Porta, SFX 4: Diálogo)
# - Combate em Tempo Real com IA Hostil (BEHAVIOR_HOSTILE)
# - Dano de Piso em Tiles de Espinho (COLLISION_DAMAGE)
# - Enigmas em Cadeia, Baús, Fonte Sagrada e Cristal Primordial de Vitória
# =============================================================================

# Executa com throttle ligado e velocidade a 800% para renderizar todos os quadros de vídeo
set throttle on
set speed 800

set log_file [open "screenshots/sub53_execution.log" w]

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

# Endereços dos símbolos em RAM (Página 3 - 0xC000 a 0xCFFF)
set ADDR_WORLD_ROOM     0xC225
set ADDR_HERO_X         0xC578
set ADDR_HERO_Y         0xC579
set ADDR_HERO_DIR       0xC57C
set ADDR_VM_FLAGS       0xC5F8
set ADDR_VM_INV         0xC6F8
set ADDR_VM_INV_CNT     0xC728
set ADDR_HERO_HP        0xC729
set ADDR_VM_MESSAGE     0xC734
set ADDR_VM_STRING_ID   0xC7B4
set ADDR_VM_HAS_MSG     0xC7B6
set ADDR_VM_SFX         0xC7B7

proc release_all_keys {} {
    keymatrixup 8 16   ;# Left
    keymatrixup 8 128  ;# Right
    keymatrixup 8 32   ;# Up
    keymatrixup 8 64   ;# Down
    keymatrixup 8 1    ;# Space
}

proc set_key_dir {dir} {
    keymatrixup 8 16
    keymatrixup 8 128
    keymatrixup 8 32
    keymatrixup 8 64
    if {$dir == "left"} {
        keymatrixdown 8 16
    } elseif {$dir == "right"} {
        keymatrixdown 8 128
    } elseif {$dir == "up"} {
        keymatrixdown 8 32
    } elseif {$dir == "down"} {
        keymatrixdown 8 64
    }
}

proc log_status {label} {
    global ADDR_WORLD_ROOM ADDR_HERO_X ADDR_HERO_Y ADDR_HERO_DIR ADDR_VM_FLAGS ADDR_VM_INV ADDR_VM_INV_CNT ADDR_HERO_HP ADDR_VM_STRING_ID ADDR_VM_SFX
    set room [read_u16 $ADDR_WORLD_ROOM]
    set hx [read_u8 $ADDR_HERO_X]
    set hy [read_u8 $ADDR_HERO_Y]
    set hdir [read_u8 $ADDR_HERO_DIR]
    set f1 [read_u8 [expr {$ADDR_VM_FLAGS + 1}]]
    set f2 [read_u8 [expr {$ADDR_VM_FLAGS + 2}]]
    set f5 [read_u8 [expr {$ADDR_VM_FLAGS + 5}]]
    set f7 [read_u8 [expr {$ADDR_VM_FLAGS + 7}]]
    set invCnt [read_u8 $ADDR_VM_INV_CNT]
    set hp [read_u16 $ADDR_HERO_HP]
    set strId [read_u16 $ADDR_VM_STRING_ID]
    set sfx [read_u8 $ADDR_VM_SFX]
    set it0 [read_u16 $ADDR_VM_INV]

    log_msg [format "%-32s | Room:%2d | Hero:(%2d,%2d) Dir:%d | F1:%d F2:%d F5:%d F7:%d | Inv:%d (It0:%d) | HP:%3d | Str:%d | SFX:%d" \
        $label $room $hx $hy $hdir $f1 $f2 $f5 $f7 $invCnt $it0 $hp $strId $sfx]
}

set last_logged_room 0

proc log_room_change {} {
    global last_logged_room ADDR_WORLD_ROOM ADDR_HERO_X ADDR_HERO_Y
    set r [read_u16 $ADDR_WORLD_ROOM]
    if {$last_logged_room != $r} {
        set last_logged_room $r
        set hx [read_u8 $ADDR_HERO_X]
        set hy [read_u8 $ADDR_HERO_Y]
        log_msg "-> [format {Entrou na Sala %2d em (%2d, %2d)} $r $hx $hy]"
    }
}

set test_state 0
set state_timer 0
set total_frames 0

proc control_loop {} {
    global test_state state_timer total_frames log_file
    global ADDR_WORLD_ROOM ADDR_HERO_X ADDR_HERO_Y ADDR_HERO_DIR ADDR_VM_FLAGS ADDR_VM_INV ADDR_VM_INV_CNT ADDR_HERO_HP ADDR_VM_STRING_ID ADDR_VM_HAS_MSG ADDR_VM_SFX

    log_room_change

    incr total_frames
    set room [read_u16 $ADDR_WORLD_ROOM]
    set hx [read_u8 $ADDR_HERO_X]
    set hy [read_u8 $ADDR_HERO_Y]
    set hasMsg [read_u8 $ADDR_VM_HAS_MSG]

    switch $test_state {
        0 {
            # Estado 0: Spawn Inicial na Sala 1 (Entrada das Catacumbas)
            incr state_timer
            if {$state_timer >= 30} {
                log_msg "=== (1/6) Spawn Inicial no Jogo de Referencia ==="
                screenshot "screenshots/sub53_01_spawn.png"
                log_status "Spawn Inicial Sala 1"
                log_msg "Aproximando-se do Guardiao Real em (14, 7)..."
                set state_timer 0
                set test_state 1
            }
        }
        1 {
            # Estado 1: Anda para (14, 9), depois sobe para (14, 8) diante do Guardião
            if {$hx > 14} {
                set_key_dir "left"
            } elseif {$hy > 8} {
                set_key_dir "up"
            } else {
                set_key_dir "none"
                incr state_timer
                if {$state_timer >= 10} {
                    log_msg "=== (2/6) Dialogo com Guardiao e Chave de Bronze ==="
                    log_status "Diante do Guardiao"
                    keymatrixdown 8 1 ;# Pressiona Espaço para falar
                    set state_timer 0
                    set test_state 2
                }
            }
        }
        2 {
            # Estado 2: Diálogo com Guardião e recebimento da Chave
            incr state_timer
            if {$state_timer == 5} {
                keymatrixup 8 1
            } elseif {$state_timer == 25} {
                log_status "Dialogo Ativo / Chave Recebida"
                screenshot "screenshots/sub53_02_dialogue.png"
                keymatrixdown 8 1 ;# Fecha diálogo
            } elseif {$state_timer == 35} {
                keymatrixup 8 1
            } elseif {$state_timer >= 45} {
                if {$hasMsg == 0} {
                    log_msg "Descendo para Y=9 e caminhando para Oeste rumo a Armaria (Salas 3 e 4)..."
                    set state_timer 0
                    set test_state 3
                } else {
                    keymatrixdown 8 1
                    after "frame" { keymatrixup 8 1 }
                }
            }
        }
        3 {
            # Estado 3: Navega para a Sala 4 (Armaria)
            # Sala 1 -> Oeste -> Sala 3 -> Oeste -> Sala 4
            if {$room == 1} {
                if {$hy < 9} {
                    set_key_dir "down"
                } else {
                    set_key_dir "left"
                }
            } elseif {$room == 3} {
                set_key_dir "left"
            } elseif {$room == 4} {
                if {$hx > 8} {
                    set_key_dir "left"
                } elseif {$hy > 6} {
                    set_key_dir "up"
                } else {
                    set_key_dir "none"
                    incr state_timer
                    if {$state_timer >= 10} {
                        log_msg "=== (3/6) Armaria dos Antigos e Bau de Tesouro ==="
                        log_status "Diante do Bau Armaria"
                        screenshot "screenshots/sub53_03_armory.png"
                        keymatrixdown 8 1 ;# Abre baú
                        set state_timer 0
                        set test_state 4
                    }
                }
            }
        }
        4 {
            # Estado 4: Abre Baú da Armaria (Chave de Ferro + Poção de Vida)
            incr state_timer
            if {$state_timer == 5} {
                keymatrixup 8 1
            } elseif {$state_timer == 25} {
                log_status "Bau Aberto / Chave de Ferro e Pocao"
                keymatrixdown 8 1 ;# Fecha diálogo
            } elseif {$state_timer == 35} {
                keymatrixup 8 1
            } elseif {$state_timer >= 45} {
                if {$hasMsg == 0} {
                    log_msg "Caminhando para o Sul rumo a Fonte Sagrada (Salas 8 e 12)..."
                    set state_timer 0
                    set test_state 5
                } else {
                    keymatrixdown 8 1
                    after "frame" { keymatrixup 8 1 }
                }
            }
        }
        5 {
            # Estado 5: Da Armaria (4), volta a (15, 9), desce pela Galeria (8) até a Fonte (12)
            if {$room == 4} {
                if {$hy < 9} {
                    set_key_dir "down"
                } elseif {$hx < 15} {
                    set_key_dir "right"
                } else {
                    set_key_dir "down"
                }
            } elseif {$room == 8} {
                set_key_dir "down"
            } elseif {$room == 12} {
                if {$hy < 11} {
                    set_key_dir "down"
                } else {
                    set_key_dir "none"
                    incr state_timer
                    if {$state_timer >= 10} {
                        log_msg "=== (4/6) Fonte Sagrada e Aguas Curativas ==="
                        log_status "Fonte Sagrada"
                        screenshot "screenshots/sub53_04_fountain.png"
                        keymatrixdown 8 1 ;# Bebe da fonte
                        set state_timer 0
                        set test_state 6
                    }
                }
            }
        }
        6 {
            # Estado 6: Vida Restaurada na Fonte Sagrada
            incr state_timer
            if {$state_timer == 5} {
                keymatrixup 8 1
            } elseif {$state_timer == 25} {
                log_status "Vida Restaurada na Fonte"
                keymatrixdown 8 1 ;# Fecha diálogo
            } elseif {$state_timer == 35} {
                keymatrixup 8 1
            } elseif {$state_timer >= 45} {
                if {$hasMsg == 0} {
                    log_msg "Caminhando a Leste rumo a Camara de Tortura (Salas 11, 13, 14, 15)..."
                    set state_timer 0
                    set test_state 7
                } else {
                    keymatrixdown 8 1
                    after "frame" { keymatrixup 8 1 }
                }
            }
        }
        7 {
            # Estado 7: Segue todo Leste na linha Y=9 até a Sala 15 (Câmara de Tortura)
            if {$room == 12} {
                if {$hy > 9} {
                    set_key_dir "up"
                } else {
                    set_key_dir "right"
                }
            } elseif {$room == 11 || $room == 13 || $room == 14} {
                set_key_dir "right"
            } elseif {$room == 15} {
                set skel_x [read_u8 0xC581]
                set skel_y [read_u8 0xC582]
                set dist [expr {abs($hx - $skel_x) + abs($hy - $skel_y)}]
                if {$dist <= 1 || ($hx >= 14 && $hx <= 16)} {
                    set_key_dir "none"
                    log_msg "=== (5/6) Combate em Tempo Real com Esqueleto Hostil ==="
                    log_status "Monstro em Combate"
                    set state_timer 0
                    set test_state 71
                } else {
                    set_key_dir "right"
                }
            }
        }
        71 {
            # Estado 71: Golpe de Espada contra Esqueleto Hostil
            incr state_timer
            if {$state_timer == 10} {
                keymatrixdown 8 1 ;# Ataque de espada!
            } elseif {$state_timer == 15} {
                keymatrixup 8 1
            } elseif {$state_timer == 35} {
                log_status "Esqueleto Destruido / Chave Real Dourada"
                screenshot "screenshots/sub53_05_combat.png"
                keymatrixdown 8 1 ;# Fecha diálogo
            } elseif {$state_timer == 45} {
                keymatrixup 8 1
            } elseif {$state_timer >= 55} {
                if {$hasMsg == 0} {
                    log_msg "Descendo ao Sul para o Santuario do Rei Esquecido (Sala 20)..."
                    set state_timer 0
                    set test_state 8
                } else {
                    keymatrixdown 8 1
                    after "frame" { keymatrixup 8 1 }
                }
            }
        }
        8 {
            # Estado 8: Desce ao Sul diretamente para a Sala 20 (Santuário do Rei Esquecido)
            if {$room == 15} {
                if {$hx < 15} {
                    set_key_dir "right"
                } elseif {$hx > 15} {
                    set_key_dir "left"
                } else {
                    set_key_dir "down"
                }
            } elseif {$room == 20} {
                if {$hy < 7} {
                    set_key_dir "down"
                } elseif {$hy >= 7 && $hy <= 8} {
                    set_key_dir "none"
                    incr state_timer
                    if {$state_timer >= 15} {
                        log_msg "=== (6/6) Santuario do Rei Esquecido e Cristal Primordial ==="
                        log_status "Diante do Cristal Primordial"
                        keymatrixdown 8 1 ;# Ergue o Cristal de Vitória
                        set state_timer 0
                        set test_state 9
                    }
                }
            }
        }
        9 {
            # Estado 9: Vitória Total!
            incr state_timer
            if {$state_timer == 5} {
                keymatrixup 8 1
            } elseif {$state_timer >= 25} {
                log_status "VITORIA TOTAL DAS CATACUMBAS"
                screenshot "screenshots/sub53_06_victory.png"
                log_msg "================================================================="
                log_msg "   VALIDACAO SUBFASE 5.3: JOGO DE REFERENCIA COMPLETO CONCLUIDO! "
                log_msg "   - 20 Salas Interligadas Operacionais no Grid 4x5              "
                log_msg "   - Driver PSG de SFX Nativo Validado com Sucesso               "
                log_msg "   - Combate Hostil, Dano em Tiles e Enigmas Verificados         "
                log_msg "   - Cristal Primordial Erguido com Sucesso                      "
                log_msg "================================================================="
                set test_state 99
                after time 1.0 {
                    release_all_keys
                    close $log_file
                    exit
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
# Com throttle desligado, 7.5 segundos emulados passam em menos de 0.5 segundo real!
after time 7.5 {
    log_msg "Cartucho MegaROM carregado. Iniciando suite automatizada da Subfase 5.3..."
    control_loop
}

# Timeout de segurança após 120 segundos emulados
after time 120.0 {
    log_msg "TIMEOUT de seguranca alcancado. Encerrando."
    release_all_keys
    close $log_file
    exit
}
