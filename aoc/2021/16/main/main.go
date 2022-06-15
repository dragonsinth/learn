package main

import (
	"fmt"
	"strings"
)

var (
	charMap = map[byte][]byte{
		'0': {0, 0, 0, 0},
		'1': {0, 0, 0, 1},
		'2': {0, 0, 1, 0},
		'3': {0, 0, 1, 1},
		'4': {0, 1, 0, 0},
		'5': {0, 1, 0, 1},
		'6': {0, 1, 1, 0},
		'7': {0, 1, 1, 1},
		'8': {1, 0, 0, 0},
		'9': {1, 0, 0, 1},
		'A': {1, 0, 1, 0},
		'B': {1, 0, 1, 1},
		'C': {1, 1, 0, 0},
		'D': {1, 1, 0, 1},
		'E': {1, 1, 1, 0},
		'F': {1, 1, 1, 1},
	}
)

func main() {
	samples := []string{
		`D2FE28`,
		`38006F45291200`,
		`EE00D40C823060`,
		`8A004A801A8002F478`,
		`620080001611562C8802118E34`,
		`C0015000016115A2E0802F182340`,
		`A0016C880162017C3686B18A3D4780`,

		`C200B40A82`,
		`04005AC33890`,
		`880086C3E88112`,
		`CE00C43D881120`,
		`D8005AC2A8F0`,
		`F600BC2D8F`,
		`9C005AC2F8F0`,
		`9C0141080250320F1802104A08`,
	}

	for _, input := range samples {
		fmt.Println(input)

		var in []byte
		for _, c := range []byte(strings.TrimSpace(input)) {
			v, ok := charMap[c]
			if !ok {
				panic(c)
			}
			in = append(in, v...)
		}

		s := &stream{
			in:  in,
			pos: 0,
		}

		p := s.parsePacket()
		s.ensureTrailingZeros()
		fmt.Printf("version=%d, value=%d\n", p.versionSum(), p.value())
	}
}

type stream struct {
	in  []byte
	pos int
}

func (s *stream) next() byte {
	ret := s.in[s.pos]
	s.pos++
	return ret
}

func (s *stream) parsePacket() *packet {
	var ret packet
	ret.version = s.parseNumber(3)
	ret.typeId = s.parseNumber(3)

	if ret.typeId == 4 {
		ret.literal = s.parseLiteral()
	} else {
		typeId := s.next()
		if typeId == 0 {
			length := s.parseNumber(15)
			dst := s.pos + length
			for s.pos < dst {
				ret.subs = append(ret.subs, s.parsePacket())
			}
		} else {
			nPackets := s.parseNumber(11)
			for i := 0; i < nPackets; i++ {
				ret.subs = append(ret.subs, s.parsePacket())
			}
		}
	}

	return &ret
}

func (s *stream) parseNumber(nBits int) int {
	ret := 0
	for i := 0; i < nBits; i++ {
		ret <<= 1
		ret += int(s.next())
	}
	return ret
}

func (s *stream) parseLiteral() int {
	ret := 0
	for {
		hasMore := s.next()
		ret <<= 4
		ret += s.parseNumber(4)
		if hasMore == 0 {
			return ret
		}
	}
}

func (s *stream) ensureTrailingZeros() {
	//rem := len(s.in) - s.pos
	//if rem > 3 {
	//	panic(fmt.Sprintf("%d bits remain", rem))
	//}
	for s.pos < len(s.in) {
		c := s.next()
		if c != 0 {
			panic("non-zero bits remain")
		}
	}
}

type packet struct {
	version int
	typeId  int
	literal int

	subs []*packet
}

func (p *packet) versionSum() int {
	ret := p.version
	for _, sub := range p.subs {
		ret += sub.versionSum()
	}
	return ret
}

func (p *packet) value() int {
	if p.typeId == 4 {
		return p.literal
	}

	var vals []int
	for _, sub := range p.subs {
		vals = append(vals, sub.value())
	}

	switch p.typeId {
	case 0: // +
		sum := 0
		for _, v := range vals {
			sum += v
		}
		return sum
	case 1: // *
		prod := 1
		for _, v := range vals {
			prod *= v
		}
		return prod
	case 2: // min
		min := vals[0]
		for _, v := range vals {
			if v < min {
				min = v
			}
		}
		return min
	case 3: // max
		max := vals[0]
		for _, v := range vals {
			if v > max {
				max = v
			}
		}
		return max
	case 5: // >
		return boolToInt(vals[0] > vals[1])
	case 6: // <
		return boolToInt(vals[0] < vals[1])
	case 7: // ==
		return boolToInt(vals[0] == vals[1])
	default:
		panic(p.typeId)
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func (p *packet) String() string {
	var buf strings.Builder
	p.ToString("", &buf)
	return buf.String()
}

func (p *packet) ToString(indent string, buf *strings.Builder) {
	buf.WriteString(indent)
	if p.typeId == 4 {
		buf.WriteString(fmt.Sprintf("v=%d, type=%d, literal=%d\n", p.version, p.typeId, p.literal))
	} else {
		buf.WriteString(fmt.Sprintf("v=%d, type=%d\n", p.version, p.typeId))
		indent += "  "
		for _, sub := range p.subs {
			sub.ToString(indent, buf)
		}
	}
}
