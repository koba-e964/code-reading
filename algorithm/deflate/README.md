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

## 例
フォーマットの例を示す。

infgen というツールを使用した:
<https://github.com/madler/infgen/tree/master> (version: [3.2](https://github.com/madler/infgen/commit/2d2300507d24b398dfc7482f3429cc0061726c8b))

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

## 参考資料

[[RFC1951]]: TODO

[[RFC1952]]: TODO

[RFC1951]: https://datatracker.ietf.org/doc/html/rfc1951
[RFC1952]: https://datatracker.ietf.org/doc/html/rfc1952
