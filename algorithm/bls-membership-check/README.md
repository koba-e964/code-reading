# BLS curve における高速な所属判定

## 前提知識
BLS12-381 などの BLS 曲線は以下のように定義されている。[[Edg2023]]

- $u$: パラメーター。2 の高いベキを約数に持つべきである。
- $r = u^4 - u^2 + 1$ は素数。
- $p = (u-1)^2r/3 + u$ は素数。
- $E(\mathbb{F}_p): y^2 = x^3 + 4 \text{ in } \mathbb{F}_p$
  - 位数は $p-u = (u-1)^2r/3$ である。 $h_1 := (u-1)^2/3$ とすれば $h_1r$ と表せる。
- $E'(\mathbb{F} _ {p^2}): y^2 = x^3 + 4(1+i) \text{ in } \mathbb{F} _ {p^2}$
  - $i$ は $i^2 = -1$ を満たす。
  - 位数は $h_2r$ である。ただし $h_2 :=$ (複雑な $u$ の式) であり $\mathrm{gcd}(h_1,h_2) = 1$ を満たす。
- $\mathbb{G} _ 1 := E(\mathbb{F} _ p)[r]$ は位数 $r$ の巡回群。
- $\mathbb{G} _ 2 := E'(\mathbb{F} _ {p^2})[r]$ は位数 $r$ の巡回群。
- $\mathbb{G} _ T := (\mathbb{F} _ {p^{12}})^{\times}[r]$ は位数 $r$ の巡回群。

$E(\mathbb{F} _ p)$ や $E'(\mathbb{F} _ {p^2})$ の点が与えられたとき、それが $\mathbb{G}_1$ や $\mathbb{G}_2$ に属するとは限らないため、所属判定が必要である。そのような判定は $rP = O$ かどうかの判定で簡単にできるが、ここではそれよりも高速な手法を扱う。

## $\psi$ の定義
TODO

## $\mathbb{G}_1$
$E(\mathbb{F}_p)$ において $\psi(x, y) = (\beta x, y)$ が成立する。ただし $\beta$ は曲線ごとに定まる 1 の原始 3 乗根である。

$P \in \mathbb{G}_1$ に対して
$\psi(P) = -u^2P$ か判定する。


## $\mathbb{G}_2$
$Q \in \mathbb{G}_2$ に対して
$\psi(Q) = uP$ か判定する。

## $\mathbb{G}_T$
$w \in \mathbb{G}_T$ に対して
$w^p = w^u$ か判定する。

# 参考文献

[[Edg2023]] Edgington, Ben. "BLS12-381 for the Rest of Us." HackMD, <https://hackmd.io/@benjaminion/bls12-381>. Accessed 18 May 2025.

[[Sco2011]] Scott, Michael. "A note on group membership tests for $\mathbb{G} _ 1$, $\mathbb{G} _ 2$ and $\mathbb{G} _ T$ on BLS pairing-friendly curves." Cryptology ePrint Archive (2021).

[Sco2011]: https://eprint.iacr.org/2021/1130

[Edg2023]: https://hackmd.io/@benjaminion/bls12-381
