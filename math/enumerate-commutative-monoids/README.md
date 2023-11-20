位数・冪等元の個数ごとに可換モノイドの列挙をして commutative-monoids.json に書き込む。位数 <= 5 の可換モノイドを列挙する。

以下の数列は https://oeis.org/A058142 と合致する。

```
$ cat commutative-monoids.json | jq '.orderwise[].idemwise[].monoids|length'
1
1
1
1
3
1
2
9
6
2
1
26
30
16
5
```
