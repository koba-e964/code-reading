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

[[RFC1951]]: TODO

[RFC1951]: https://datatracker.ietf.org/doc/html/rfc1951
