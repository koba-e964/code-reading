# 概要
gzip や deflate のフォーマットの雑な紹介を行う。

# フォーマットの説明
## gzip フォーマット
gzip フォーマットは、ヘッダー、圧縮データ、フッターの3つの部分から成り立っており、member と呼ばれている。([[RFC1952]] の 2.3. Member format)
- ヘッダー: ファイルのメタデータ。10 バイト。
  - |名前|オフセット (バイト)|長さ|説明|
    |--|--|--|--|
    |ID|0|2|固定値 (`0x1f, 0x8b`)|
    |CM|2|1|圧縮アルゴリズム。多くの場合で 8 = deflate|
    |FLG|3|1|フラグ|
    |MTIME|4|4|Unix epoch からの秒数|
    |XFL|8|1|追加のフラグ|
    |OS|9|1|OS。3 = Unix|

    なお、FLG.FEXTRA FLG.FNAME FLG.COMMENT はここでは省略する。
- 圧縮データ: 実際のデータ (deflate によるもの)
  - deflate はビット指向の圧縮アルゴリズムだが、gzip に格納される時点でバイト単位になるようパディングされる。
- フッター: チェックサムと圧縮前のデータ長。8 バイト。
  - |名前|オフセット (バイト)|長さ|説明|
    |--|--|--|--|
    |CRC32|0|4|圧縮前のデータの CRC32 チェックサム。|
    |ISIZE|4|4|圧縮前のデータの長さ mod $2^{32}$|

なお、仕様ではこれら (members) の複数の連結が許可されているが、対応しているソフトウェアがあまりないようである。
> A gzip file consists of a series of "members" (compressed data sets).

## deflate フォーマット
deflate フォーマットは、ヘッダーと圧縮データの2つの部分から成り立っている。**ビット指向**のフォーマットである。([[RFC1951]] の 3. Detailed specification)
- ヘッダー: 圧縮情報
- 圧縮データ: ブロックの繰り返し。([[RFC1951]] の 3.2. Compressed block format)

### Huffman coding
[[RFC1951]] の 3.2.2. Use of Huffman coding in the "deflate" format の説明は理解に時間がかかったため、理解を助けるための追加の説明を行う。

Huffman code の canonical な割り当ては、以下のように (精度付き) 二進小数の範囲の割り当てだと思うとわかりやすい。

|ビット数|頻度|範囲|要素一覧|
|--|--|--|--|
|1|0|[0b0.0, 0b0.0)|(空集合)|
|2|1|[0b0.00, 0b0.01)|0b0.00|
|3|5|[0b0.010, 0b0.111)|0b0.010, 0b0.011, 0b0.100, 0b0.101, 0b0.110|
|4|2|[0b0.1110, 0b1.0000)|0b0.1110, 0b0.1111|

ビット数が短いものに小さい実数が割り当てられる。また、同じビット数の文字の中では若いものに小さい実数が割り当てられる。このように範囲の割り当てを定める方法はほとんど明らかだろう。

実際に next_code, tree[I].code を計算するステップをこの方法で解釈しよう。`2)` は以下のように解釈できる。(next_code は範囲の最初の符号語)

|N|bl_count|next_code|範囲|
|--|--|--|--|
|1|0|0 = 0b0|[0b0.0, 0b0.0)|
|2|1|0 = 0b00|[0b0.00, 0b0.01)|
|3|5|2 = 0b010|[0b0.010, 0b0.111)|
|4|2|14 = 0b1110|[0b0.1110, 0b1.0000)|

`3)` が終わった後は以下のようになる。(next_code は範囲の終わりを示すようになる)

|N|bl_count|next_code|範囲|
|--|--|--|--|
|1|0|0 = 0b0|[0b0.0, 0b0.0)|
|2|1|1 = 0b01|[0b0.00, 0b0.01)|
|3|5|7 = 0b111|[0b0.010, 0b0.111)|
|4|2|16 = 0b10000|[0b0.1110, 0b1.0000)|

また、符号語は以下のようになる。

|文字|長さ|符号語|二進小数としての表現|
|--|--|--|--|
|A|3|010|`0b0.010`|
|B|3|011|`0b0.011`|
|C|3|100|`0b0.100`|
|D|3|101|`0b0.101`|
|E|3|110|`0b0.110`|
|F|2|00|`0b0.00`|
|G|4|1110|`0b0.1110`|
|H|4|1111|`0b0.1111`|

実数として重複がないこと、先頭から二進小数として読んでいって実数の範囲によって読むべきビット数が決定できることは明らか。

これらの情報はすべて符号語長の長さの配列 [3, 3, 3, 3, 3, 2, 4, 4] から復元可能である。

# 例
フォーマットの例を示す。

infgen というツールを使用した:
<https://github.com/madler/infgen> (version: [3.2](https://github.com/madler/infgen/commit/2d2300507d24b398dfc7482f3429cc0061726c8b))

コンパイル方法は以下:
```
cc -I/opt/homebrew/opt/zlib/include -L/opt/homebrew/opt/zlib/lib infgen.c -o infgen -lz
```

依存関係は以下:
```
$ otool -L ./infgen
./infgen:
	/opt/homebrew/opt/zlib/lib/libz.1.dylib (compatibility version 1.0.0, current version 1.3.0)
	/usr/lib/libSystem.B.dylib (compatibility version 1.0.0, current version 1336.61.1)
```

結果は以下:
```console
$ echo 123123123123123123123 | gzip -9 | ./infgen -idds
! infgen 3.2 output
!
time 1704260380		! [UTC Wed Jan  3 05:39:40 2024]
xfl 2
gzip
!
last			! 1
fixed			! 01
literal '1		! 10000110
literal '2		! 01000110
literal '3		! 11000110
literal '1		! 10000110
match 17 3		! 01000 0 0011000
literal 10		! 01011100
end			! 0000000
! stats literals 8.0 bits each (40/5)
! stats matches 77.3% (1 x 17.0)
! stats inout 7:7 (6) 22 0
			! 0
! stats total inout 7:7 (6) 22
! stats total block average 22.0 uncompressed
! stats total block average 6.0 symbols
! stats total literals 8.0 bits each
! stats total matches 77.3% (1 x 17.0)
!
crc
length
```

```console
$ echo 123123123123123123123 | gzip -9 | hexdump -C
00000000  1f 8b 08 00 fc 59 96 65  02 03 33 34 32 36 c4 40  |.....Y.e..3426.@|
00000010  5c 00 5e 96 a9 24 16 00  00 00                    |\.^..$....|
0000001a
```
また、Go で作ったツールで中身の解析を行うと以下のようになっている。

```
$ echo 123123123123123123123 | gzip -9 | go run .
0000: 1f 8b       ID (0x1f 0x8b)
0002: 08          CM
0003: 00          FLG
0004: 06 c4 a5 65 MTIME 2024-01-15T23:47:18Z
0008: 02          XFL
0009: 03          OS
000a: 33 34 32 36 Data
      c4 40 5c 00 
0012: 5e 96 a9 24 CRC32 0x24a9965e
0016: 16 00 00 00 isize 22
[0000,0003)  3   0000: 33    .....011          deflate header
[0000,0001)  1   0000: 33    .......1          BFINAL 1
[0001,0003)  2   0000: 33    .....01.          BTYPE 1
[0003,003f)  60  0000: 33 34 00110... 00110100 deflate content
                       32 36 00110010 00110110 
                       c4 40 11000100 01000000 
                       5c 00 01011100 .0000000 
[0003,000b)  8   0000: 33 34 00110... .....100 Literal "01100001" -> 49
[000b,0013)  8   0001: 34 32 00110... .....010 Literal "01100010" -> 50
[0013,001b)  8   0002: 32 36 00110... .....110 Literal "01100011" -> 51
[001b,0023)  8   0003: 36 c4 00110... .....100 Literal "01100001" -> 49
[0023,002a)  7   0004: c4 40 11000... ......00 Length "0001100" -> 268 (calculated: 17)
[002a,002b)  1   0005: 40    .....0..          LExtra "0" -> 0 (calculated: 17)
[002b,0030)  5   0005: 40    01000...          Distance "00010" -> 2 (calculated: 3)
[0030,0038)  8   0006: 5c    01011100          Literal "00111010" -> 10
[0038,003f)  7   0007: 00    .0000000          End "0000000" -> 256
```

# 無圧縮つ gzip を作るツールについて
以下のコマンドでインストールできる。
```bash
go install github.com/koba-e964/code-reading/algorithm/deflate/cmd/gnozip
```

# 参考資料

[[RFC1951]]: Deutsch, P., "DEFLATE Compressed Data Format Specification version 1.3", RFC 1951, DOI 10.17487/RFC1951, May 1996, <https://www.rfc-editor.org/info/rfc1951>.

[[RFC1952]]: Deutsch, P., "GZIP file format specification version 4.3", RFC 1952, DOI 10.17487/RFC1952, May 1996, <https://www.rfc-editor.org/info/rfc1952>.

[RFC1951]: https://datatracker.ietf.org/doc/html/rfc1951
[RFC1952]: https://datatracker.ietf.org/doc/html/rfc1952
