package main

import "fmt"

func main() {
	b := new(BitWriter)

	b.Bits("1")        // BFINAL=1
	b.Bits("10")       // BTYPE=01
	b.Bits("00110001") // Lit 0x01 = Direct 0x01
	b.Bits("00110010") // Lit 0x02 = Direct 0x02
	b.Bits("00110011") // Lit 0x03 = Direct 0x03
	b.Bits("00110100") // Lit 0x04 = Direct 0x04
	b.Bits("0000001")  // Lit 0x101 = <length=3>
	b.Bits("00011")    // Dist 0x03 = Distance 0x04
	b.Bits("0000011")  // Lit 0x103 = <length=5>
	b.Bits("00010")    // Dist 0x02 = Distance 0x03
	b.Bits("00110110") // Lit 0x06 = Direct 0x06
	b.Bits("0000011")  // Lit 0x103 = <length=5>
	b.Bits("00100")    // Dist 0x04 = Distance [0x05, 0x07), Extra = 1
	b.Bits("1")        // DExtra: Distance = 0x06
	b.Bits("00110011") // Lit 0x03 = Direct 0x03
	b.Bits("00110010") // Lit 0x02 = Direct 0x02
	b.Bits("00110011") // Lit 0x03 = Direct 0x03
	b.Bits("00110001") // Lit 0x01 = Direct 0x01
	b.Bits("00110010") // Lit 0x02 = Direct 0x02
	b.Bits("00110011") // Lit 0x03 = Direct 0x03
	b.Bits("0000000")  // Lit 0x100 = end

	bs := b.Emit()
	// [99 100 98 102 1 98 48 98 3 147 204 76 204 140 76 204]
	// Test case from https://github.com/Frommi/miniz_oxide/blob/0.7.0/miniz_oxide/src/deflate/core.rs#L2456
	fmt.Println(bs)
}
