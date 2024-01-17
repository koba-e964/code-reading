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
	Letters []HuffmanSymbol
}

func (f *FixedHuffman) Show() []RangeDescription {
	var result []RangeDescription
	for _, letter := range f.Letters {
		result = append(result, letter.Show()...)
	}
	return result
}

type DynamicHuffman struct {
	Ints    []FixedInt
	Letters []HuffmanSymbol
}

func (d *DynamicHuffman) Show() []RangeDescription {
	var result []RangeDescription
	for _, i := range d.Ints {
		result = append(result, i.Show()...)
	}
	for _, letter := range d.Letters {
		result = append(result, letter.Show()...)
	}
	return result
}

type FixedInt struct {
	StartPos   uint64
	Length     uint64
	Type       string
	Value      int
	Calculated int
}

func (f *FixedInt) Show() []RangeDescription {
	return []RangeDescription{
		{
			StartPos:    f.StartPos,
			Length:      f.Length,
			Description: fmt.Sprintf("%s %d", f.Type, f.Value),
		},
	}
}

type HuffmanSymbol struct {
	StartPos    uint64
	Length      uint64
	Type        string
	HuffmanCode string
	Value       int
	Calculated  int
}

func (f *HuffmanSymbol) Show() []RangeDescription {
	description := fmt.Sprintf("%s \"%s\" -> %d", f.Type, f.HuffmanCode, f.Value)
	if f.Calculated != 0 {
		description += fmt.Sprintf(" (calculated: %d)", f.Calculated)
	}
	return []RangeDescription{{StartPos: f.StartPos, Length: f.Length, Description: description}}
}

// formatAsBits formats the given value as a little endian bit string.
func formatAsBits(value uint64, length int) string {
	result := ""
	for i := 0; i < length; i++ {
		if (value>>i)&1 == 1 {
			result += "1"
		} else {
			result += "0"
		}
	}
	return result
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

// Ref: https://github.com/madler/infgen/blob/2d2300507d24b398dfc7482f3429cc0061726c8b/infgen.c#L864-L931
type Huffman struct {
	// The number of codes for each code length.
	NumCodes []int
	// The table of (length, code) -> symbol.
	Symbols [][]int
}

func ConstructHuffman(lengths []int) *Huffman {
	n := len(lengths)
	maxLength := 0
	for _, length := range lengths {
		maxLength = max(maxLength, length)
	}
	numCodes := make([]int, maxLength+1)
	symbols := make([][]int, maxLength+1)
	for i := 0; i < n; i++ {
		numCodes[lengths[i]]++
	}
	code := 0
	for i := 1; i <= maxLength; i++ {
		code = (code + numCodes[i-1]) << 1
		symbols[i] = make([]int, 0, numCodes[i])
	}
	for i := 0; i < n; i++ {
		length := lengths[i]
		if length != 0 {
			symbols[length] = append(symbols[length], i)
		}
	}
	return &Huffman{
		NumCodes: numCodes,
		Symbols:  symbols,
	}
}

// Returns (length, symbol, rawBits, error)
func (h *Huffman) Decode(b *BitReader) (int, int, int, error) {
	code := 0
	length := 1
	first := 0
	for length < len(h.NumCodes) {
		code = code*2 + int(b.Int(1))
		count := h.NumCodes[length]
		if code < first+count {
			return length, h.Symbols[length][code-first], code, nil
		}
		first += count
		first <<= 1
		length++
	}
	return -1, -1, -1, fmt.Errorf("invalid huffman code")
}

// parseLengthDistancePair parses a length-distance pair.
func parseLengthDistancePair(b *BitReader, startPos, firstHuffmanCode uint64, firstHuffmanSymbol int, distanceHuffman *Huffman) ([]HuffmanSymbol, error) {
	var letters []HuffmanSymbol
	// 3.2.5. Compressed blocks (length and distance codes)
	// length
	if firstHuffmanSymbol < 257 || firstHuffmanSymbol > 285 {
		return nil, ErrInvalidLength
	}
	var length int
	var extraLengthBits int
	if firstHuffmanSymbol < 265 {
		// [3, 11)
		length = firstHuffmanSymbol - 254
		extraLengthBits = 0
	} else if firstHuffmanSymbol < 269 {
		// [11, 19)
		length = 11 + (firstHuffmanSymbol-265)*2
		extraLengthBits = 1
	} else if firstHuffmanSymbol < 273 {
		// [19, 35)
		length = 19 + (firstHuffmanSymbol-269)*4
		extraLengthBits = 2
	} else if firstHuffmanSymbol < 277 {
		// [35, 67)
		length = 35 + (firstHuffmanSymbol-273)*8
		extraLengthBits = 3
	} else if firstHuffmanSymbol < 281 {
		// [67, 131)
		length = 67 + (firstHuffmanSymbol-277)*16
		extraLengthBits = 4
	} else if firstHuffmanSymbol < 285 {
		// [131, 258)
		length = 131 + (firstHuffmanSymbol-281)*32
		extraLengthBits = 5
	} else {
		length = 258
		extraLengthBits = 0
	}
	letters = append(letters, HuffmanSymbol{
		StartPos:    startPos,
		Length:      b.Position() - startPos,
		Type:        "Length",
		HuffmanCode: formatAsBEBits(firstHuffmanCode, int(b.Position()-startPos)),
		Value:       firstHuffmanSymbol,
		Calculated:  length,
	})
	startPos = b.Position()
	extra := int(b.Int(extraLengthBits))
	length += extra
	letters = append(letters, HuffmanSymbol{
		StartPos:    startPos,
		Length:      b.Position() - startPos,
		Type:        "LExtra",
		HuffmanCode: formatAsBits(uint64(extra), extraLengthBits),
		Value:       extra,
		Calculated:  length,
	})
	// distance
	startPos = b.Position()
	distanceLength, distanceCode, distanceRaw, err := distanceHuffman.Decode(b)
	if err != nil {
		return nil, err
	}
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
	letters = append(letters, HuffmanSymbol{
		StartPos:    startPos,
		Length:      b.Position() - startPos,
		Type:        "Distance",
		HuffmanCode: formatAsBEBits(uint64(distanceRaw), distanceLength),
		Value:       distanceCode,
		Calculated:  distance,
	})
	startPos = b.Position()
	extra = int(b.Int(extraLengthBits))
	distance += extra
	letters = append(letters, HuffmanSymbol{
		StartPos:    startPos,
		Length:      b.Position() - startPos,
		Type:        "DExtra",
		HuffmanCode: formatAsBits(uint64(extra), extraLengthBits),
		Value:       extra,
		Calculated:  distance,
	})
	filteredLetters := []HuffmanSymbol{}
	for _, letter := range letters {
		if letter.Length > 0 {
			filteredLetters = append(filteredLetters, letter)
		}
	}
	return filteredLetters, nil
}

func parseBlock(b *BitReader, lengthHuffman *Huffman, distanceHuffman *Huffman) ([]HuffmanSymbol, error) {
	var letters []HuffmanSymbol
	finished := false
	for !finished {
		startPos := b.Position()
		huffmanLength, huffmanSymbol, huffmanCode, err := lengthHuffman.Decode(b)
		if err != nil {
			return nil, err
		}
		// huffmanSymbol is in [0, 288)
		var ty string
		if huffmanSymbol == 256 {
			ty = "End"
			finished = true
		} else if huffmanSymbol < 256 {
			// literal
			ty = "Literal"
		} else if huffmanSymbol < 286 {
			// length-distance pair
			parsedLetters, err := parseLengthDistancePair(b, startPos, uint64(huffmanCode), huffmanSymbol, distanceHuffman)
			if err != nil {
				return nil, err
			}
			letters = append(letters, parsedLetters...)
			continue
		}
		letters = append(letters, HuffmanSymbol{
			StartPos:    startPos,
			Length:      b.Position() - startPos,
			Type:        ty,
			HuffmanCode: formatAsBEBits(uint64(huffmanCode), huffmanLength),
			Value:       huffmanSymbol,
		})
	}
	return letters, nil
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
			lengthLengths := make([]int, 288)
			for i := 0; i < 144; i++ {
				lengthLengths[i] = 8
			}
			for i := 144; i < 256; i++ {
				lengthLengths[i] = 9
			}
			for i := 256; i < 280; i++ {
				lengthLengths[i] = 7
			}
			for i := 280; i < 288; i++ {
				lengthLengths[i] = 8
			}
			lengthHuffman := ConstructHuffman(lengthLengths)
			distanceLengths := make([]int, 30)
			for i := 0; i < 30; i++ {
				distanceLengths[i] = 5
			}
			distanceHuffman := ConstructHuffman(distanceLengths)
			letters, err := parseBlock(b, lengthHuffman, distanceHuffman)
			if err != nil {
				return nil, err
			}
			deflate.Content = &FixedHuffman{Letters: letters}
		} else if deflate.Header.BTYPE == 2 {
			// dynamic huffman codes
			ints := []FixedInt{}
			letters := []HuffmanSymbol{}
			startPos := b.Position()
			numLiteralLengthLengths := int(b.Int(5)) + 257
			ints = append(ints, FixedInt{
				StartPos:   startPos,
				Length:     5,
				Type:       "NumLitLen",
				Value:      numLiteralLengthLengths - 257,
				Calculated: numLiteralLengthLengths,
			})
			startPos = b.Position()
			numDistanceLengths := int(b.Int(5)) + 1
			ints = append(ints, FixedInt{
				StartPos:   startPos,
				Length:     5,
				Type:       "NumDist",
				Value:      numDistanceLengths - 1,
				Calculated: numDistanceLengths,
			})
			startPos = b.Position()
			numCodeLengths := int(b.Int(4)) + 4
			ints = append(ints, FixedInt{
				StartPos:   startPos,
				Length:     4,
				Type:       "NumCode",
				Value:      numCodeLengths - 4,
				Calculated: numCodeLengths,
			})
			codeLengths := make([]int, 19)
			codeOrder := []int{16, 17, 18, 0, 8, 7, 9, 6, 10, 5, 11, 4, 12, 3, 13, 2, 14, 1, 15}
			for i := 0; i < numCodeLengths; i++ {
				startPos = b.Position()
				codeLengths[codeOrder[i]] = int(b.Int(3))
				ints = append(ints, FixedInt{
					StartPos:   startPos,
					Length:     3,
					Type:       fmt.Sprintf("CodeLen[%d]", codeOrder[i]),
					Value:      codeLengths[codeOrder[i]],
					Calculated: codeLengths[codeOrder[i]],
				})
			}
			codeHuffman := ConstructHuffman(codeLengths)
			lengths := make([]int, 0, numLiteralLengthLengths+numDistanceLengths)
			for len(lengths) < numLiteralLengthLengths+numDistanceLengths {
				startPos = b.Position()
				codeLength, codeSymbol, codeRawBits, err := codeHuffman.Decode(b)
				if err != nil {
					return nil, err
				}
				letters = append(letters, HuffmanSymbol{
					StartPos:    startPos,
					Length:      uint64(codeLength),
					Type:        fmt.Sprintf("Code[%d]", len(lengths)),
					HuffmanCode: formatAsBEBits(uint64(codeRawBits), codeLength),
					Value:       codeSymbol,
				})
				startPos = b.Position()
				if codeSymbol < 16 {
					lengths = append(lengths, codeSymbol)
				} else if codeSymbol == 16 {
					if len(lengths) == 0 {
						return nil, ErrInvalidDeflate
					}
					lastLength := lengths[len(lengths)-1]
					repeat := 3 + int(b.Int(2))
					for i := 0; i < repeat; i++ {
						lengths = append(lengths, lastLength)
					}
					letters = append(letters, HuffmanSymbol{
						StartPos:    startPos,
						Length:      2,
						Type:        "Repeat Last",
						HuffmanCode: formatAsBits(uint64(repeat-3), 2),
						Value:       repeat - 3,
						Calculated:  repeat,
					})
				} else if codeSymbol == 17 {
					repeat := 3 + int(b.Int(3))
					for i := 0; i < repeat; i++ {
						lengths = append(lengths, 0)
					}
					letters = append(letters, HuffmanSymbol{
						StartPos:    startPos,
						Length:      3,
						Type:        "Repeat 0",
						HuffmanCode: formatAsBits(uint64(repeat-3), 3),
						Value:       repeat - 3,
						Calculated:  repeat,
					})
				} else if codeSymbol == 18 {
					repeat := 11 + int(b.Int(7))
					for i := 0; i < repeat; i++ {
						lengths = append(lengths, 0)
					}
					letters = append(letters, HuffmanSymbol{
						StartPos:    startPos,
						Length:      7,
						Type:        "Repeat 0",
						HuffmanCode: formatAsBits(uint64(repeat-11), 7),
						Value:       repeat - 11,
						Calculated:  repeat,
					})
				} else {
					return nil, ErrInvalidDeflate
				}
			}
			if len(lengths) > numLiteralLengthLengths+numDistanceLengths {
				return nil, ErrInvalidDeflate
			}
			lengthHuffman := ConstructHuffman(lengths[:numLiteralLengthLengths])
			distanceHuffman := ConstructHuffman(lengths[numLiteralLengthLengths:])
			blockLetters, err := parseBlock(b, lengthHuffman, distanceHuffman)
			if err != nil {
				return nil, err
			}
			letters = append(letters, blockLetters...)
			deflate.Content = &DynamicHuffman{Ints: ints, Letters: letters}
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
