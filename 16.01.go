package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/icza/bitio"
)

type PacketHeaders struct {
	version int
	typeId  int
}

type LiteralPacket struct {
	PacketHeaders
	value int
}

type OperatorPacket struct {
	PacketHeaders
	subpackets []Packet
}

type Packet interface {
	Version() int
	TypeId() int
	Print()
}

func (this PacketHeaders) Print() {
	fmt.Printf("-- Packet ------------\n")
	fmt.Printf("Version:  %d\n", this.version)
	fmt.Printf("Type ID:  %d\n", this.typeId)
}

func (this PacketHeaders) Version() int {
	return this.version
}

func (this PacketHeaders) TypeId() int {
	return this.typeId
}

func (this LiteralPacket) Print() {
	this.PacketHeaders.Print()
	fmt.Printf("Literal:  %d\n", this.value)
}

func (this OperatorPacket) Print() {
	this.PacketHeaders.Print()
	fmt.Printf("Packets:  %d\n", len(this.subpackets))

	for _, packet := range this.subpackets {
		packet.Print()
	}
}

func readFile(filename string) []string {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	lines := []string{}
	s := bufio.NewScanner(f)

	for s.Scan() {
		err = s.Err()
		if err != nil {
			log.Fatal(err)
		}

		lines = append(lines, s.Text())
	}

	return lines
}

func parseOperatorPacket(reader *bitio.CountReader, headers PacketHeaders) (*OperatorPacket, error) {
	packets := []Packet{}

	lenTypeId, err := reader.ReadBits(1)
	if err != nil {
		log.Fatal("Failed to read operator packet length id")
	}

	if lenTypeId == 0 {
		bitLength, err := reader.ReadBits(15)
		if err != nil {
			log.Fatal("Failed to read operator packet")
		}

		start := reader.BitsCount
		for start+int64(bitLength) > reader.BitsCount {
			packet, err := parsePacket(reader)
			if err != nil {
				log.Fatal("Failed to read subpacket")
			}

			packets = append(packets, packet)
		}

	} else if lenTypeId == 1 {
		numPackets, err := reader.ReadBits(11)
		if err != nil {
			log.Fatal("Failed to read operator packet")
		}

		for i := uint64(0); i < numPackets; i++ {
			packet, err := parsePacket(reader)
			if err != nil {
				log.Fatal("Failed to read subpacket")
			}

			packets = append(packets, packet)
		}

	} else {
		log.Fatalf("Unexpected operator packet type id %d\n", headers.typeId)
	}

	packet := OperatorPacket{
		headers,
		packets,
	}

	return &packet, nil
}

func parseLiteralPacket(reader *bitio.CountReader, headers PacketHeaders) (*LiteralPacket, error) {
	lastSegment := false
	value := 0

	for !lastSegment {
		last, err := reader.ReadBits(1)
		if err != nil {
			log.Print(err)
			log.Fatalf("Failed to parse literal packet 1")
		}
		lastSegment = (last == 0)

		segment, err := reader.ReadBits(4)
		if err != nil {
			log.Print(err)
			log.Fatalf("Failed to parse literal packet 2")
		}

		value = value << 4
		value += int(segment)
	}

	literal := LiteralPacket{headers, value}
	return &literal, nil
}

func parsePacket(reader *bitio.CountReader) (Packet, error) {
	version, err := reader.ReadBits(3)
	if err != nil {
		return nil, err
	}

	typeId, err := reader.ReadBits(3)
	if err != nil {
		return nil, err
	}

	headers := PacketHeaders{version: int(version), typeId: int(typeId)}
	var packet Packet

	if typeId == 4 {
		packet, err = parseLiteralPacket(reader, headers)
	} else {
		packet, err = parseOperatorPacket(reader, headers)
	}

	if err != nil {
		log.Println("Failed to parse packet")
	}

	return packet, err
}

func versionSum(packet Packet) int {
	sum := packet.Version()

	switch packet.(type) {
	case *OperatorPacket:
		op := packet.(*OperatorPacket)
		for _, p := range op.subpackets {
			sum += versionSum(p)
		}
	}

	return sum
}

func evaluate(packet Packet) int {
	value := 0

	switch packet.TypeId() {
	case 0: // sum
		for _, p := range packet.(*OperatorPacket).subpackets {
			value += evaluate(p)
		}

	case 1: // product
		subpackets := packet.(*OperatorPacket).subpackets
		value = evaluate(subpackets[0])
		for i := 1; i < len(subpackets); i++ {
			value *= evaluate(subpackets[i])
		}

	case 2: // minimum
		subpackets := packet.(*OperatorPacket).subpackets
		value = evaluate(subpackets[0])
		for i := 1; i < len(subpackets); i++ {
			newVal := evaluate(subpackets[i])
			if newVal < value {
				value = newVal
			}
		}

	case 3: // maximum
		subpackets := packet.(*OperatorPacket).subpackets
		value = evaluate(subpackets[0])
		for i := 1; i < len(subpackets); i++ {
			newVal := evaluate(subpackets[i])
			if newVal > value {
				value = newVal
			}
		}

	case 4: // literal
		value = packet.(*LiteralPacket).value

	case 5: // greater than
		subpackets := packet.(*OperatorPacket).subpackets
		if evaluate(subpackets[0]) > evaluate(subpackets[1]) {
			value = 1
		}

	case 6: // less than
		subpackets := packet.(*OperatorPacket).subpackets
		if evaluate(subpackets[0]) < evaluate(subpackets[1]) {
			value = 1
		}

	case 7: // greater than
		subpackets := packet.(*OperatorPacket).subpackets
		if evaluate(subpackets[0]) == evaluate(subpackets[1]) {
			value = 1
		}

	}

	return value
}

func main() {
	lines := readFile("16.txt")
	data, err := hex.DecodeString(lines[0])
	if err != nil {
		log.Fatal(err)
	}

	reader := bitio.NewCountReader(bytes.NewBuffer(data))
	packet, err := parsePacket(reader)

	if err != nil {
		log.Fatal(err)
	}

	//packet.Print()

	sum := versionSum(packet)
	fmt.Printf("Packet version sum: %d\n", sum)
	value := evaluate(packet)
	fmt.Printf("Packet value: %d\n", value)
}
