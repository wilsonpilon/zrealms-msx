// ____________________________
// Z-Realm (zrealm-msx) - Binary Disk Loader para MSX-DOS 2
//─────────────────────────────────────────────────────────────────────────────
#include "loader.h"

static MasterHeader  s_MasterHeader;
static ResourceEntry s_ResourceDirectory[MAX_RESOURCES_DIR];

// Carrega o HEADER.BIN e faz o streaming de GAME.DAT para os segmentos do Mapper
bool LOADER_LoadGame(const c8* headerPath, const c8* dataPath)
{
	u8  hFile;
	u16 bytesRead;
	u16 seg;

	// 1. Abre e lê HEADER.BIN
	hFile = DOS_OpenHandle(headerPath, O_RDONLY);
	if (hFile == 0 || hFile == HANDLE_INVALID)
	{
		DOS_StringOutput("Erro: Nao foi possivel abrir HEADER.BIN!\r\n$");
		return FALSE;
	}

	bytesRead = DOS_ReadHandle(hFile, &s_MasterHeader, sizeof(MasterHeader));
	if (bytesRead != sizeof(MasterHeader))
	{
		DOS_CloseHandle(hFile);
		DOS_StringOutput("Erro: Falha na leitura do MasterHeader!\r\n$");
		return FALSE;
	}

	// Valida assinatura 'ZR01' e versão
	if (s_MasterHeader.Magic[0] != LOADER_MAGIC_0 ||
	    s_MasterHeader.Magic[1] != LOADER_MAGIC_1 ||
	    s_MasterHeader.Magic[2] != LOADER_MAGIC_2 ||
	    s_MasterHeader.Magic[3] != LOADER_MAGIC_3 ||
	    s_MasterHeader.Version != LOADER_VERSION)
	{
		DOS_CloseHandle(hFile);
		DOS_StringOutput("Erro: Formato ou versao de HEADER.BIN invalido!\r\n$");
		return FALSE;
	}

	// Lê o diretório de recursos
	if (s_MasterHeader.ResourceCount > MAX_RESOURCES_DIR)
	{
		DOS_CloseHandle(hFile);
		DOS_StringOutput("Erro: Excesso de recursos em HEADER.BIN!\r\n$");
		return FALSE;
	}

	if (s_MasterHeader.ResourceCount > 0)
	{
		u16 dirBytes = s_MasterHeader.ResourceCount * sizeof(ResourceEntry);
		bytesRead = DOS_ReadHandle(hFile, s_ResourceDirectory, dirBytes);
		if (bytesRead != dirBytes)
		{
			DOS_CloseHandle(hFile);
			DOS_StringOutput("Erro: Falha na leitura do diretorio de recursos!\r\n$");
			return FALSE;
		}
	}

	DOS_CloseHandle(hFile);

	// 2. Aloca segmentos de Memory Mapper
	if (!MAPPER_Init(s_MasterHeader.SegmentCount))
	{
		DOS_StringOutput("Erro: Mapper insuficiente para os segmentos do jogo!\r\n$");
		return FALSE;
	}

	// 3. Abre e carrega GAME.DAT diretamente nos segmentos mapeados na Página 2
	hFile = DOS_OpenHandle(dataPath, O_RDONLY);
	if (hFile == 0 || hFile == HANDLE_INVALID)
	{
		DOS_StringOutput("Erro: Nao foi possivel abrir GAME.DAT!\r\n$");
		return FALSE;
	}

	for (seg = 0; seg < s_MasterHeader.SegmentCount; seg++)
	{
		// Mapeia o segmento lógico 'seg' na Página 2 (0x8000 - 0xBFFF)
		MAPPER_SetPage2(seg);

		// Lê 16.384 bytes em dois blocos de 8192 bytes
		bytesRead = DOS_ReadHandle(hFile, (void*)0x8000, 8192);
		if (bytesRead != 8192)
		{
			DOS_CloseHandle(hFile);
			DOS_StringOutput("Erro: Leitura incompleta de GAME.DAT (bloco 1)!\r\n$");
			return FALSE;
		}

		bytesRead = DOS_ReadHandle(hFile, (void*)0xA000, 8192);
		if (bytesRead != 8192)
		{
			DOS_CloseHandle(hFile);
			DOS_StringOutput("Erro: Leitura incompleta de GAME.DAT (bloco 2)!\r\n$");
			return FALSE;
		}
	}

	DOS_CloseHandle(hFile);
	return TRUE;
}

// Retorna ponteiro constante para o MasterHeader carregado
const MasterHeader* LOADER_GetMasterHeader(void)
{
	return &s_MasterHeader;
}

// Busca uma entrada de recurso no diretório pelo tipo e ID lógico
const ResourceEntry* LOADER_FindResource(u8 type, u16 id)
{
	u16 i;
	for (i = 0; i < s_MasterHeader.ResourceCount; i++)
	{
		if (s_ResourceDirectory[i].Type == type && s_ResourceDirectory[i].ID == id)
		{
			return &s_ResourceDirectory[i];
		}
	}
	return NULL;
}

// Mapeia o segmento que contém o recurso na Página 2 e retorna o ponteiro absoluto
void* LOADER_MapResource(u8 type, u16 id)
{
	const ResourceEntry* entry = LOADER_FindResource(type, id);
	if (!entry)
	{
		return NULL;
	}

	MAPPER_SetPage2(entry->Segment);
	return (void*)(0x8000 + entry->Offset);
}

// Atalho para mapear e obter ponteiro para uma Sala (BinaryRoom)
BinaryRoom* LOADER_GetRoom(u16 roomID)
{
	return (BinaryRoom*)LOADER_MapResource(RES_TYPE_ROOM, roomID);
}

// Atalho para mapear e obter ponteiro para um Tileset (BinaryTileset)
BinaryTileset* LOADER_GetTileset(u16 tilesetID)
{
	return (BinaryTileset*)LOADER_MapResource(RES_TYPE_TILESET, tilesetID);
}
