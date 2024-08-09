# Streamline NTRU Prime
[[BBCC+2020]] の 2.3 に沿って説明を行う。

# sntrup761
パラメーターは以下の通り。
- $p = 761, q = 4591$

sntrup761 の公開鍵は [1158 バイトである](https://github.com/openssh/openssh-portable/blob/V_9_8_P1/sntrup761.c#L254)。
- 公開鍵は $[-(q-1)/2, (q-1)/2]$ の範囲内の整数 $p$ 個である。

1158 はほとんど $(p \log_2 q)/8$ に等しい。
```irb
$ irb
irb(main):001> include Math
=> Object
irb(main):005> 761 * log2(4591) / 8
=> 1157.156882177264
```

また暗号文は 1039 バイトである。これは素の暗号文 (1007 バイト) と Confirm (32 バイト) の連結である。

- 素の暗号文は $[-(q-1)/6, (q-1)/6]$ の範囲内の整数 $p$ 個 (と同等のもの) である。
- Confirm はハッシュ値である。 `Hash(0x02 || Hash(0x03 || r) || Hash(0x04 || pk))` で計算される。ここで r は平文、pk は公開鍵、Hash(b) は [SHA-512(b)[0:32] である](https://github.com/openssh/openssh-portable/blob/V_9_8_P1/sntrup761.c#L685-L696)。

1007 はほとんど $(p \log_2 ((q-1)/3 + 1))/8$ に等しい。
```irb
irb(main):007> 761 * log2(4591/3+1) / 8
=> 1006.4470962334174
```

## エンコーディング方法
特殊。Range encoding のような手法を使い、情報理論的にほぼ最適なデータ量になるように整数列をバイト列へとエンコードする。

# 参考文献

[BBCC+2020]: https://ntruprime.cr.yp.to/nist/ntruprime-20201007.pdf

[[BBCC+2020]] Daniel J. Bernstein, Billy Bob Brumley, Ming-Shing Chen, Chitchanok Chuengsatiansup, Tanja Lange, Adrian Marotzke, Bo-Yuan Peng, Nicola Tuveri, Christine van Vredendaal, Bo-Yin Yang, NTRU Prime: round 3 (2020) URL: https://ntruprime.cr.yp.to/nist/ntruprime-20201007.pdf
