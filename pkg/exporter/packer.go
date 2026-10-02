package exporter

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Packer gerencia o empacotamento de recursos binários em segmentos contíguos de 16 KB (0x4000).
type Packer struct {
	segments       [][]byte
	currentSegment int
	currentOffset  int
	resources      []ResourceEntry
}

// NewPacker instancia um novo empacotador de segmentos para o MSX.
func NewPacker() *Packer {
	p := &Packer{
		segments:       make([][]byte, 0),
		currentSegment: -1,
		currentOffset:  0,
		resources:      make([]ResourceEntry, 0),
	}
	p.addSegment()
	return p
}

func (p *Packer) addSegment() {
	newSeg := make([]byte, SegmentSize)
	p.segments = append(p.segments, newSeg)
	p.currentSegment = len(p.segments) - 1
	p.currentOffset = 0
}

// WriteResource adiciona um recurso binário a um segmento de 16 KB.
// Se o recurso não couber no segmento atual, um novo segmento de 16 KB é iniciado automaticamente.
func (p *Packer) WriteResource(resType uint8, resID uint16, data []byte) (ResourceEntry, error) {
	dataLen := len(data)
	if dataLen > SegmentSize {
		return ResourceEntry{}, fmt.Errorf("recurso (tipo=%d, id=%d) excede o tamanho máximo de 16 KB: %d bytes", resType, resID, dataLen)
	}

	// Se não couber no segmento atual, abre novo segmento de 16 KB
	if p.currentOffset+dataLen > SegmentSize {
		p.addSegment()
	}

	entry := ResourceEntry{
		Type:    resType,
		Segment: uint8(p.currentSegment),
		ID:      resID,
		Offset:  uint16(p.currentOffset),
		Size:    uint16(dataLen),
	}

	copy(p.segments[p.currentSegment][p.currentOffset:], data)
	p.currentOffset += dataLen
	p.resources = append(p.resources, entry)

	return entry, nil
}

// SegmentCount retorna o número total de segmentos de 16 KB alocados.
func (p *Packer) SegmentCount() int {
	return len(p.segments)
}

// ResourceCount retorna a quantidade de recursos empacotados.
func (p *Packer) ResourceCount() int {
	return len(p.resources)
}

// GetSegments retorna a fatia de todos os segmentos de 16 KB gerados.
func (p *Packer) GetSegments() [][]byte {
	return p.segments
}

// GetFlatData retorna todos os segmentos concatenados em um único buffer alinhado a 16 KB (GAME.DAT).
func (p *Packer) GetFlatData() []byte {
	totalBytes := len(p.segments) * SegmentSize
	flat := make([]byte, totalBytes)
	for i, seg := range p.segments {
		copy(flat[i*SegmentSize:], seg)
	}
	return flat
}

// BuildMasterHeader monta a tabela mestre binária (HEADER.BIN) em formato Little-Endian (Z80 nativo).
// Estrutura:
// - 32 bytes: MasterHeader
// - N * 8 bytes: ResourceEntry (diretório de alocação de recursos)
func (p *Packer) BuildMasterHeader(initialRoomID uint16, initialHeroX, initialHeroY, initialTileset uint8) ([]byte, error) {
	buf := new(bytes.Buffer)

	// 1. Grava MasterHeader (32 bytes)
	var magicBytes [4]byte
	copy(magicBytes[:], HeaderMagic)

	if err := buf.WriteByte(magicBytes[0]); err != nil {
		return nil, err
	}
	_ = buf.WriteByte(magicBytes[1])
	_ = buf.WriteByte(magicBytes[2])
	_ = buf.WriteByte(magicBytes[3])

	if err := binary.Write(buf, binary.LittleEndian, uint16(FormatVersion)); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, uint16(len(p.segments))); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, uint16(len(p.resources))); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, initialRoomID); err != nil {
		return nil, err
	}
	_ = buf.WriteByte(initialHeroX)
	_ = buf.WriteByte(initialHeroY)
	_ = buf.WriteByte(initialTileset)

	// Preenchimento de alinhamento para 32 bytes (atualmente 15 bytes gravados, faltam 17)
	padding := make([]byte, 17)
	if _, err := buf.Write(padding); err != nil {
		return nil, err
	}

	// 2. Grava Diretório de Recursos (8 bytes por entrada)
	for _, res := range p.resources {
		_ = buf.WriteByte(res.Type)
		_ = buf.WriteByte(res.Segment)
		if err := binary.Write(buf, binary.LittleEndian, res.ID); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.LittleEndian, res.Offset); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.LittleEndian, res.Size); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
