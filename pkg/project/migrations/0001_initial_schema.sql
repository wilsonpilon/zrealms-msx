-- Schema DDL inicial para projetos Z-Realm (zrealm-msx)
-- Suporte completo a MSX 2 (V9938 / SCREEN 4 / Sprites Modo 2 / Memory Mapper)

-- Controle de migrações
CREATE TABLE IF NOT EXISTS schema_migrations (
    version INTEGER PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Configurações e metadados globais do projeto
CREATE TABLE IF NOT EXISTS project_settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

-- Conjuntos de padrões de tiles (8x8 pixels)
CREATE TABLE IF NOT EXISTS tilesets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT
);

-- Definição individual de tiles (V9938 SCREEN 4 / Graphic 3)
CREATE TABLE IF NOT EXISTS tiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    tileset_id INTEGER NOT NULL REFERENCES tilesets(id) ON DELETE CASCADE,
    tile_index INTEGER NOT NULL,          -- Índice de 0 a 255 no tileset local
    pattern_bytes BLOB NOT NULL,          -- 8 bytes (1 bit por pixel, monocromático por linha)
    color_bytes BLOB NOT NULL,            -- 8 bytes (2 cores por linha: alto nibble foreground, baixo background)
    collision_type INTEGER DEFAULT 0,     -- 0: Passável, 1: Sólido, 2: Água, 3: Dano, etc.
    animation_next_tile_id INTEGER REFERENCES tiles(id) ON DELETE SET NULL,
    UNIQUE(tileset_id, tile_index)
);

-- Sprites (16x16 pixels - Modo 2 do V9938)
CREATE TABLE IF NOT EXISTS sprites (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    pattern_bytes BLOB NOT NULL,          -- 32 bytes (4 blocos de 8x8 para formar 16x16)
    color_bytes BLOB NOT NULL             -- 16 bytes (cor individual para cada scanline de 0 a 15)
);

-- Scripts e Roteiros de Eventos
CREATE TABLE IF NOT EXISTS scripts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    source_code TEXT NOT NULL,            -- Código de script em alto nível
    bytecode BLOB                         -- Bytecode compilado para execução na VM do MSX
);

-- Salas / Telas Estáticas (Grid de 32x18 tiles)
CREATE TABLE IF NOT EXISTS rooms (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    world_x INTEGER NOT NULL,
    world_y INTEGER NOT NULL,
    name TEXT NOT NULL,
    tileset_id INTEGER NOT NULL REFERENCES tilesets(id),
    tile_matrix BLOB NOT NULL,            -- 576 bytes (32 colunas * 18 linhas de índices de tiles)
    north_room_id INTEGER REFERENCES rooms(id) ON DELETE SET NULL,
    south_room_id INTEGER REFERENCES rooms(id) ON DELETE SET NULL,
    east_room_id  INTEGER REFERENCES rooms(id) ON DELETE SET NULL,
    west_room_id  INTEGER REFERENCES rooms(id) ON DELETE SET NULL,
    UNIQUE(world_x, world_y)
);

-- Atores, Entidades e Triggers posicionados nas salas
CREATE TABLE IF NOT EXISTS entities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    room_id INTEGER NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    pos_x INTEGER NOT NULL CHECK(pos_x >= 0 AND pos_x < 32),
    pos_y INTEGER NOT NULL CHECK(pos_y >= 0 AND pos_y < 18),
    sprite_id INTEGER REFERENCES sprites(id) ON DELETE SET NULL,
    behavior_type INTEGER NOT NULL DEFAULT 0, -- 0: NPC Estático, 1: NPC Errante, 2: Baú, 3: Gatilho, etc.
    event_script_id INTEGER REFERENCES scripts(id) ON DELETE SET NULL
);

-- Tabela de Textos e Diálogos (suporta compressão DTE futura)
CREATE TABLE IF NOT EXISTS string_table (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    context_tag TEXT NOT NULL,
    text_content TEXT NOT NULL
);

-- Classes de Heróis
CREATE TABLE IF NOT EXISTS hero_classes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    base_hp INTEGER NOT NULL,
    base_mp INTEGER NOT NULL,
    base_atk INTEGER NOT NULL,
    base_def INTEGER NOT NULL
);

-- Tabela de Itens
CREATE TABLE IF NOT EXISTS items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    item_type INTEGER NOT NULL,           -- 0: Arma, 1: Armadura, 2: Consumível, 3: Chave, 4: Missão
    modifier_stat INTEGER DEFAULT 0,
    modifier_value INTEGER DEFAULT 0,
    price INTEGER DEFAULT 0
);
