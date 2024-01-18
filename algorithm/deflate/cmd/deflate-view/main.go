package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"

	gzip "github.com/koba-e964/code-reading/algorithm/deflate"
)

type TableEmitter struct {
	widths []int
}

// NewTableEmitter creates a new TableEmitter.
//
// widths is a slice of column widths.
// For example, []int{2, 3} means that the first column has width 2, the second column has width 3 and
// the other columns can have arbitrarily wide columns.
func NewTableEmitter(widths []int) *TableEmitter {
	return &TableEmitter{widths: widths}
}

func (t *TableEmitter) EmitRow(values ...[]string) {
	maxRow := 0
	for _, value := range values {
		maxRow = max(maxRow, len(value))
	}
	outputLines := make([]string, maxRow)
	for i, value := range values {
		for j := 0; j < maxRow; j++ {
			var line string
			if j < len(value) {
				line = value[j]
			}
			for i < len(t.widths) && len(line) < t.widths[i] {
				line += " "
			}
			outputLines[j] += line
		}
	}
	for _, line := range outputLines {
		fmt.Println(line)
	}
}

func emitRowForGzip(startpos uint64, bytes []byte, explanation string, a ...any) {
	widths := []int{6, 12}
	te := NewTableEmitter(widths)
	BytesWidth := 4
	values := [3][]string{}
	values[0] = []string{fmt.Sprintf("%04x: ", startpos)}
	var current string
	for index, b := range bytes {
		current += fmt.Sprintf("%02x ", b)
		if index%BytesWidth == BytesWidth-1 && index != len(bytes)-1 {
			values[1] = append(values[1], current)
			current = ""
		}
	}
	if current != "" {
		values[1] = append(values[1], current)
	}
	formatted := fmt.Sprintf(explanation, a...)
	values[2] = []string{formatted}
	te.EmitRow(values[:]...)
}

func main() {
	file := os.Stdin
	if len(os.Args) > 1 {
		currentFile, err := os.Open(os.Args[1])
		if err != nil {
			panic(err)
		}
		file = currentFile
	}
	defer file.Close()
	stream, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	entry, err := gzip.ParseGzip(stream)
	if err != nil {
		panic(err)
	}
	emitRowForGzip(0, entry.Header.ID[:], "ID (0x1f 0x8b)")
	emitRowForGzip(2, []byte{entry.Header.CompressionMethod}, "CM")
	emitRowForGzip(3, []byte{entry.Header.Flags}, "FLG")
	emitRowForGzip(4, stream[4:8], "MTIME %s", time.Unix(int64(entry.Header.Mtime), 0).UTC().Format(time.RFC3339))
	emitRowForGzip(8, []byte{entry.Header.XFL}, "XFL")
	emitRowForGzip(9, []byte{entry.Header.OS}, "OS")
	emitRowForGzip(entry.DataPos, entry.Data, "Data")
	emitRowForGzip(entry.FooterPos, binary.LittleEndian.AppendUint32(nil, entry.Footer.CRC32), "CRC32 0x%08x", entry.Footer.CRC32)
	emitRowForGzip(entry.FooterPos+4, binary.LittleEndian.AppendUint32(nil, entry.Footer.Isize), "isize %d", entry.Footer.Isize)

	// deflate
	deflates, err := gzip.ParseDeflate(entry.Data)
	if err != nil {
		panic(err)
	}
	te := NewTableEmitter([]int{
		13, 4,
		6, 6, 18,
	})
	for _, deflate := range deflates {
		for _, value := range deflate.Show() {
			rangeString := fmt.Sprintf("[%04x,%04x)", value.StartPos, value.StartPos+value.Length)
			lenString := fmt.Sprintf("%d", value.Length)
			bytePosStr := fmt.Sprintf("%04x:", value.StartPos/8)
			bytesStr := []string{}
			bitsStr := []string{}
			bytesCurrent := ""
			bitsCurrent := ""
			displayed := 0
			for i := value.StartPos / 8; i < (value.StartPos+value.Length+7)/8; i++ {
				bytesCurrent += fmt.Sprintf("%02x ", entry.Data[i])
				for j := 7; j >= 0; j-- {
					pos := 8*i + uint64(j)
					if pos >= value.StartPos && pos < value.StartPos+value.Length {
						bitsCurrent += fmt.Sprintf("%d", (entry.Data[i]>>j)&1)
					} else {
						bitsCurrent += "."
					}
				}
				bitsCurrent += " "
				displayed++
				if displayed%2 == 0 {
					bytesStr = append(bytesStr, bytesCurrent)
					bitsStr = append(bitsStr, bitsCurrent)
					bytesCurrent = ""
					bitsCurrent = ""
				}
			}
			if bytesCurrent != "" {
				bytesStr = append(bytesStr, bytesCurrent)
				bitsStr = append(bitsStr, bitsCurrent)
			}
			te.EmitRow([]string{rangeString}, []string{lenString},
				[]string{bytePosStr}, bytesStr, bitsStr,
				[]string{value.Description},
			)
		}
	}
}
