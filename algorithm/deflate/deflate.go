package main

import (
	"fmt"
)

// RangeDescription is a struct for describing a range of bits.
//
// StartPos is the absolute position of the first bit.
type RangeDescription struct {
	StartPos    uint64
	Length      uint64
	Description string
}

// Show is an interface for showing the content of a deflate block or element.
type Show interface {
	Show() []RangeDescription
}

type Deflate struct {
	Pos        uint64
	Header     DeflateHeader
	ContentPos uint64
	Content    Show
	EndPos     uint64
}

func (d *Deflate) Show() []RangeDescription {
	var result []RangeDescription
	result = append(result, RangeDescription{d.Pos, 3, "deflate header"})
	result = append(result, d.Header.Show()...)
	result = append(result, RangeDescription{d.ContentPos, d.EndPos - d.ContentPos, "deflate content"})
	result = append(result, d.Content.Show()...)
	return result
}

type DeflateHeader struct {
	StartPos uint64
	BFINAL   byte
	BTYPE    byte
}

func (d *DeflateHeader) Show() []RangeDescription {
	return []RangeDescription{
		{StartPos: d.StartPos, Length: 1, Description: fmt.Sprintf("BFINAL %d", d.BFINAL)},
		{StartPos: d.StartPos + 1, Length: 2, Description: fmt.Sprintf("BTYPE %d", d.BTYPE)},
	}
}

type Uncompressed struct {
	StartPos uint64
	Length   uint16
	NLength  uint16
	Literal  []byte
}

func (u *Uncompressed) Show() []RangeDescription {
	var result []RangeDescription
	result = append(result, RangeDescription{StartPos: u.StartPos, Length: 16, Description: "Length"})
	result = append(result, RangeDescription{StartPos: u.StartPos + 16, Length: 16, Description: "NLength (inverted)"})
	result = append(result, RangeDescription{
		StartPos:    u.StartPos + 32,
		Length:      uint64(len(u.Literal)) * 8,
		Description: fmt.Sprintf("literal %dbytes", len(u.Literal)),
	})
	return result
}

type FixedHuffman struct {
	Letters []FixedAlphabet
}

func (f *FixedHuffman) Show() []RangeDescription {
	var result []RangeDescription
	for _, letter := range f.Letters {
		result = append(result, letter.Show()...)
	}
	return result
}

type FixedAlphabet struct {
	StartPos    uint64
	Length      uint64
	Type        string
	HuffmanCode string
	Value       int
	Calculated  int
}

func (f *FixedAlphabet) Show() []RangeDescription {
	description := fmt.Sprintf("%s \"%s\" -> %d", f.Type, f.HuffmanCode, f.Value)
	if f.Calculated != 0 {
		description += fmt.Sprintf(" (calculated: %d)", f.Calculated)
	}
	return []RangeDescription{{StartPos: f.StartPos, Length: f.Length, Description: description}}
}

// formatAsBEBits formats the given value as a big endian bit string.
func formatAsBEBits(value uint64, length int) string {
	result := ""
	for i := length - 1; i >= 0; i-- {
		if (value>>i)&1 == 1 {
			result += "1"
		} else {
			result += "0"
		}
	}
	return result
}

var (
	ErrInvalidDeflate = fmt.Errorf("invalid deflate")
	ErrInvalidLength  = fmt.Errorf("invalid length")
	ErrInvalidBtype   = fmt.Errorf("invalid BTYPE")
)

// parseLengthDistancePair parses a length-distance pair.
func parseLengthDistancePair(b *BitReader, startPos, firstHuffmanCode uint64, firstVal int) ([]FixedAlphabet, error) {
	var letters []FixedAlphabet
	// 3.2.5. Compressed blocks (length and distance codes)
	// length
	if firstVal < 257 || firstVal > 285 {
		return nil, ErrInvalidLength
	}
	var length int
	var extraLengthBits int
	if firstVal < 265 {
		// [3, 11)
		length = firstVal - 254
		extraLengthBits = 0
	} else if firstVal < 269 {
		// [11, 19)
		length = 11 + (firstVal-265)*2
		extraLengthBits = 1
	} else if firstVal < 273 {
		// [19, 35)
		length = 19 + (firstVal-269)*4
		extraLengthBits = 2
	} else if firstVal < 277 {
		// [35, 67)
		length = 35 + (firstVal-273)*8
		extraLengthBits = 3
	} else if firstVal < 281 {
		// [67, 131)
		length = 67 + (firstVal-277)*16
		extraLengthBits = 4
	} else if firstVal < 285 {
		// [131, 258)
		length = 131 + (firstVal-281)*32
		extraLengthBits = 5
	} else {
		length = 258
		extraLengthBits = 0
	}
	letters = append(letters, FixedAlphabet{
		StartPos:    startPos,
		Length:      b.Position() - startPos,
		Type:        "Length",
		HuffmanCode: formatAsBEBits(firstHuffmanCode, int(b.Position()-startPos)),
		Value:       firstVal,
		Calculated:  length,
	})
	startPos = b.Position()
	extra := b.IntBE(extraLengthBits)
	length += int(extra)
	letters = append(letters, FixedAlphabet{
		StartPos:    startPos,
		Length:      b.Position() - startPos,
		Type:        "LExtra",
		HuffmanCode: formatAsBEBits(extra, extraLengthBits),
		Value:       int(extra),
		Calculated:  length,
	})
	// distance
	startPos = b.Position()
	distanceCode := int(b.IntBE(5))
	if distanceCode >= 30 {
		return nil, ErrInvalidLength
	}
	var distance int
	if distanceCode < 4 {
		extraLengthBits = 0
		distance = distanceCode + 1
	} else if distanceCode < 6 {
		extraLengthBits = 1
		distance = (distanceCode-4)*2 + 5
	} else if distanceCode < 8 {
		extraLengthBits = 2
		distance = (distanceCode-6)*4 + 9
	} else if distanceCode < 10 {
		extraLengthBits = 3
		distance = (distanceCode-8)*8 + 17
	} else if distanceCode < 12 {
		extraLengthBits = 4
		distance = (distanceCode-10)*16 + 33
	} else if distanceCode < 14 {
		extraLengthBits = 5
		distance = (distanceCode-12)*32 + 65
	} else if distanceCode < 16 {
		extraLengthBits = 6
		distance = (distanceCode-14)*64 + 129
	} else if distanceCode < 18 {
		extraLengthBits = 7
		distance = (distanceCode-16)*128 + 257
	} else if distanceCode < 20 {
		extraLengthBits = 8
		distance = (distanceCode-18)*256 + 513
	} else if distanceCode < 22 {
		extraLengthBits = 9
		distance = (distanceCode-20)*512 + 1025
	} else if distanceCode < 24 {
		extraLengthBits = 10
		distance = (distanceCode-22)*1024 + 2049
	} else if distanceCode < 26 {
		extraLengthBits = 11
		distance = (distanceCode-24)*2048 + 4097
	} else if distanceCode < 28 {
		extraLengthBits = 12
		distance = (distanceCode-26)*4096 + 8193
	} else {
		extraLengthBits = 13
		distance = (distanceCode-28)*8192 + 16385
	}
	letters = append(letters, FixedAlphabet{
		StartPos:    startPos,
		Length:      b.Position() - startPos,
		Type:        "Distance",
		HuffmanCode: formatAsBEBits(uint64(distanceCode), 5),
		Value:       distanceCode,
		Calculated:  distance,
	})
	startPos = b.Position()
	extra = b.IntBE(extraLengthBits)
	letters = append(letters, FixedAlphabet{
		StartPos:    startPos,
		Length:      b.Position() - startPos,
		Type:        "DExtra",
		HuffmanCode: formatAsBEBits(extra, extraLengthBits),
		Value:       int(extra),
		Calculated:  distance,
	})
	filteredLetters := []FixedAlphabet{}
	for _, letter := range letters {
		if letter.Length > 0 {
			filteredLetters = append(filteredLetters, letter)
		}
	}
	return filteredLetters, nil
}

func ParseDeflate(stream []byte) ([]Deflate, error) {
	var deflates []Deflate
	b := NewBitReader(stream)
	final := false
	for !final {
		var deflate Deflate
		deflate.Pos = b.Position()
		deflate.Header.StartPos = b.Position()
		deflate.Header.BFINAL = byte(b.Int(1))
		if deflate.Header.BFINAL == 1 {
			final = true
		}
		deflate.Header.BTYPE = byte(b.Int(2))
		if deflate.Header.BTYPE == 0 {
			// no compression, so skip to byte boundary
			b.SkipToByteBoundary()
		}
		deflate.ContentPos = b.Position()
		if deflate.Header.BTYPE == 0 {
			// no compression
			startPos := b.Position()
			length := b.Int(16)
			nLength := b.Int(16)
			if length != nLength^0xffff {
				return nil, ErrInvalidLength
			}
			deflate.Content = &Uncompressed{StartPos: startPos, Literal: b.Bytes(int(length))}
		} else if deflate.Header.BTYPE == 1 {
			// fixed huffman codes
			var letters []FixedAlphabet
			finished := false
			for !finished {
				startPos := b.Position()
				head := int(b.IntBE(7))
				var val int
				var huffmanCode uint64
				if head < 0x18 {
					huffmanCode = uint64(head)
					val = head + 0x100
				} else if head < 0x60 {
					aux := int(b.IntBE(1))
					huffmanCode = uint64(head*2 + aux)
					val = (head-0x18)*2 + aux
				} else if head < 0x64 {
					aux := int(b.IntBE(1))
					huffmanCode = uint64(head*2 + aux)
					val = (head-0x60)*2 + aux + 280
				} else {
					aux := int(b.IntBE(2))
					huffmanCode = uint64(head*4 + aux)
					val = (head-0x64)*4 + aux + 144
				}
				// val is in [0, 288)
				var ty string
				if val == 256 {
					ty = "End"
					finished = true
				} else if val < 256 {
					// literal
					ty = "Literal"
				} else if val < 286 {
					// length-distance pair
					parsedLetters, err := parseLengthDistancePair(b, startPos, huffmanCode, val)
					if err != nil {
						return nil, err
					}
					letters = append(letters, parsedLetters...)
					continue
				}
				letters = append(letters, FixedAlphabet{
					StartPos:    startPos,
					Length:      b.Position() - startPos,
					Type:        ty,
					HuffmanCode: formatAsBEBits(huffmanCode, int(b.Position()-startPos)),
					Value:       val,
				})
			}
			deflate.Content = &FixedHuffman{Letters: letters}
		} else if deflate.Header.BTYPE == 2 {
			panic("TODO: not implemented")
		} else {
			return nil, ErrInvalidBtype
		}
		deflate.EndPos = b.Position()
		deflates = append(deflates, deflate)
	}
	if b.Remaining() >= 8 {
		return nil, ErrInvalidDeflate
	}
	return deflates, nil
}
