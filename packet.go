package rawpkt

import (
	"encoding/binary"
)

type PacketHeader struct {
	Length       uint16
	IsEncrypted  bool
	IsCompressed bool
	Type         uint16
}

const HEADER_SIZE = 6

type Packet struct {
	data []byte
}

func ToPacket(data []byte) *Packet {
	if len(data) < HEADER_SIZE {
		return nil
	}

	return &Packet{data}
}

//  NewPacket generate new packet
func NewPacket(type_ uint16, isEncrypted bool, isCompressed bool) *Packet {
	bsSize := make([]byte, 2)
	binary.LittleEndian.PutUint16(bsSize, 6)
	if isEncrypted {
		bsSize = append(bsSize, byte(1))
	} else {
		bsSize = append(bsSize, byte(0))
	}
	if isCompressed {
		bsSize = append(bsSize, byte(1))
	} else {
		bsSize = append(bsSize, byte(0))
	}
	bsType := make([]byte, 2)
	binary.LittleEndian.PutUint16(bsType, type_)
	return &Packet{
		data: append(bsSize, bsType...),
	}
}

func (p *Packet) SetType(type_ uint16) {
	bsSize := make([]byte, 2)
	binary.LittleEndian.PutUint16(bsSize, 6)
	p.data[0] = bsSize[0]
	p.data[1] = bsSize[1]
}

func (p *Packet) Size() uint16 {
	return binary.LittleEndian.Uint16(p.data[:2])
}

func (p *Packet) setSize(size uint16) {
	binary.LittleEndian.PutUint16(p.data, size)
}

func (p *Packet) addSize(size uint16) {
	p.setSize(p.Size() + size)
}

func (p *Packet) removeSize(size uint16) {
	p.setSize(p.Size() - size)
	p.data = append(p.data[:HEADER_SIZE], p.data[HEADER_SIZE+size:]...)
}

func (p *Packet) Type() uint16 {
	return binary.LittleEndian.Uint16(p.data[4:6])
}

func (p *Packet) Buffer() []byte {
	return p.data
}

func (p *Packet) Stringify() string {
	return string(p.data[HEADER_SIZE:])
}

func (p *Packet) WriteRaw(data []byte) {
	p.data = append(p.data, data...)
}
