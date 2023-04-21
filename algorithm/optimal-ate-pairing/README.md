## 概要
BLS signature や BN signature に必要な pairing 写像の実例である、optimal ate pairing について書く。

## 準備
pairing 写像: 写像 $e: G_1 \times G_2 \to G_T$ であり以下を満たすもの:
- 双線型性: $e(P_1 - P_2, Q) = e(P_1, Q)e(P_2, Q)^{-1}$ かつ $e(P, Q_1 - Q_2) = e(P, Q_1)e(P, Q_2)^{-1}$
- 非退化性: $e(P, Q) \neq 1$ なる $P, Q$ が存在すること

(登場する群はすべてアーベル群であるが、典型的には $G_1, G_2$ は加法的に書かれ $G_T$ は乗法的に書かれることに注意。)

因子: 零点や極の情報を簡潔に表すためのもの。
形式的な定義は [Wikipedia](https://en.wikipedia.org/wiki/Divisor_(algebraic_geometry)) を見ること。以下は非形式的な記述のみを行う。

以下は楕円曲線の上のことだけを考える。$\mathrm{div}(f)$ の形で表せる因子のことを**主因子** (principal divisor) と呼ぶ。(代数的整数論におけるイデアルと主イデアルの関係と同じ。) 主因子の基本的な性質として以下のようなものがある。
- 係数の和が 0: 例えば $2[P] - [O]$ のような divisor は $\mathrm{div}(f)$ としては得られない。
- 係数の重みつき和が $O$: $[P] + [Q] - [P + Q]$ のような divisor は重みつき和が $P + Q - (P + Q) = O$ なので OK。$[P] - [O]$ は重みつき和が $P - O = P$ なので、$P \neq O$ のときは $\mathrm{div}(f)$ としては得られない。

無限遠点での ord をどう求めるべきか?
1. 係数の和が 0 であることを利用する: $P = (a, b)$ に対して $\mathrm{div}(x - a) = [P] + [-P] - 2[O]$ である。$[O]$ の係数 -2 は他の項の係数の和が 2 であることから導出できる。
2. 正攻法。$\mathrm{ord}_{O}(x) = -2$ と $\mathrm{ord}_{O}(y) = -3$ を利用する。これ自体は Riemann-Roch の定理によって構成できるようである。([ACDFLNV2006] の 4.4.2.a を参照)

フロベニウス写像 $\phi_q: \mathbb{F}_{q^k} \to \mathbb{F}_{q^k}$ を $\phi(x) := x^q$ で定義するとこれは体の自己同型である。
また楕円曲線上の点 $(x, y)$ に対して $\phi_q(x, y) := (x^q, y^q)$ と定義すると $\phi_q$ は群の準同型 $E(\mathbb{F}_{q^k}) \to E(\mathbb{F}_{q^k})$ である。これは核が有限で全射であることが証明できるようであり、このような写像を同種写像 (isogeny) と呼ぶ。

## Tate pairing

具体的な構成は https://hackmd.io/@jpw/bn254 や https://hackmd.io/@benjaminion/bls12-381 が参考になる。ここでは [[Ver2008]] の Section 2.1 に沿うことにする。

前提条件として以下がある。
- $k \ge 2$
- $E$ は楕円曲線
- $r$ は $r | \#E(\mathbb{F}_{q})$ なる素数
- $r | q^k - 1, r^2 \not | q^k - 1$
- $G_1 = E(\mathbb{F}_{q^k})[r]$, $G_2 = E(\mathbb{F}_{q^k})/rE(\mathbb{F}_{q^k})$, $G_T = \mathbb{F}_{q^k}^*$

まずは基本的な関数から。$l_{P, Q}$ と $v_P$ を以下のように定義する。
- $l_{P, Q}$: $P, Q$ を通る直線において 0 をとるような関数
- $v_P$: $P$ と $-P$ を通る直線において 0 をとるような関数 ($x - P_x$)

これらの関数の因子は以下である:
- $\mathrm{div}(l_{P, Q}) = [P] + [Q] + [-(P + Q)] - 3[O]$
- $\mathrm{div}(v_P) = [P] + [-P] - 2[O]$

続いて Miller function $f_{s,P}$ を以下で定義する。
- 要件: $\mathrm{div}(f_{s,P}) = s[P] - [sP] - (s-1)[O]$
- 定義: $f_{1,P} := 1, f_{i+j,P} := f_{i,P}f_{j,P} l_{iP,jP}/v_{(i+j)P}$

こうすると $t(P, Q) := f_{r, P}(Q)^{(q^k-1)/r}$ がペアリング写像である。

Barreto-Naehrig curve では以下のようにパラメータが決められている。
- $x \in \mathbb{Z}_{\ge 0}$
- $k = 12$
- $q = p = 36x^4 + 36x^3 + 24x^2 + 6x + 1$
- $r = p - 6x^2 = 36x^4 + 36x^3 + 18x^2 + 6x + 1$

その一種、BN254 では以下のようにパラメータが決められている。
- $x = $ `4965661367192848881 = 0x44e992b44a6909f1 ~= 2^{62.1}`
## Optimal ate pairing
[[Ver2008]] の Sections 2.2, 4 に沿うことにする。

$G_1 = E[r] \cap \mathrm{ker}(\phi_q - 1) = E(\mathbb{F}_{q})[r]$, $G_2 = E[r] \cap \mathrm{ker}(\phi_q - q)$ とする。
先ほどの Tate pairing の引数を $G_2 \times G_1$ とした上で $m$ 乗することにする。
$$t(Q, P)^m = f_{r, Q}(P)^{m(q^k-1)/r}$$
これは $f_{mr, Q}(P)^{(q^k-1)/r}$ に等しい。$rQ = O$ であったため、一般に $f_{ab, P} = f_{a,P}^b f_{b,aP}$ という等式が成り立つことを使えばわかる。これもペアリング写像である。これにより、うまい $m$ を見つけて計算負荷が小さくなるようにする、ということを考えることができる。

$G_2 = E[r] \cap \mathrm{ker}(\phi_q - q)$ だったので $P \in G_2$ に対して $\phi_q(P) = qP$ が成立する。だから $q$ 倍の代わりにフロベニウス写像 $\phi_q$ を計算することができる。

ここで唐突だが、Barreto-Naehrig curve では $p^3 - p^2 + p + 6x+2 = (36x^4 + 36x^3 + 24x^2 + 6x + 1)^3 - (36x^4 + 36x^3 + 24x^2 + 6x + 1)^2 + (36x^4 + 36x^3 + 24x^2 + 6x + 1) + 6x+2$ は $r = 36x^4 + 36x^3 + 18x^2 + 6x + 1$ の倍数である。これを利用して公式を作ることができる。

以上をまとめると、$m = (p^3 - p^2 + p + 6x+2)/r$ として $f_{mr, Q}(P)^{(q^k-1)/r}$ を計算することになる。

$f_{q^i, Q} = f_{1,Q}^{q^i} = 1$ とかが成り立つことに注意。
$f_{-q^i, Q} = f_{-1,Q}^{q^i} = 1/v_Q^{q^i}$ であり、$v_Q$ は $(q^k-1)/r$ 乗すると消えるので、これを利用して以下のように計算できる。
$$f_{mr, Q}(P)^{(p^k-1)/r} \\= \left(f_{6x+2,Q}(P)l_{-p^2Q, p^3Q}(P)l_{-p^2Q + p^3Q, pQ}(P)l_{-p^2Q + p^3Q + pQ, (6x+2)Q}(P)\right)^{(p^k-1)/r}$$

繰り返しになるがここで $p^iQ = \phi_p^i(Q)$ であることに注意。
また $-p^2Q + p^3Q + pQ + (6x+2)Q = O$ であることにも注意すると $l_{-p^2Q + p^3Q + pQ, (6x+2)Q}(P) = v_{(6x+2)Q}(P)$ であるため、やはり $(q^k-1)/r$ 乗すると消える。このようにすると $f$ 1 個、$l$ 2 個の積になる。https://hackmd.io/@jpw/bn254 では別の組み合わせ方を紹介しているが、いずれにせよ $f$ 1 個、$l$ 2 個の積である。

## 疑問点
### 1
[[Ver2008]] の 2.1:
ペアリングの型が $t: E(\mathbb{F}_{q^k})[r] \times E(\mathbb{F}_{q^k})/rE(\mathbb{F}_{q^k}) \to \mathbb{F}_{q^k}^*$ なのはなぜ?
- $P$ は $rP = O$ を満たす必要あり: こうすると $\mathrm{div}(f_{r, P}) = r[P] - r[O]$ が成立する 
- $Q$ は $rE(\mathbb{F}_{q^k})$ の差を無視できる: $f_{r, P}$ は群の準同型なので引数の $rE(\mathbb{F}_{q^k})$ の差は結果の $(\mathbb{F}_{q^k}^*)^r$ の差として現れ、この差は $(q^k-1)/r$ 乗すると消える。

[[Naeh2009]] の Prop. 1.33 から、$\mathrm{Pic}^0(E) \simeq E$ がわかるので、$E$ 上の演算の代わりに $\mathrm{Pic}^0(E)$ 上の演算を定義することができる。おそらくこれを使っている。
[[FR1994]] の構成が元になっているが、[[FR1994]] ではコホモロジー理論を使っていて難解。

### 2
[[Ver2008]] の 2.2:
$G_1 = E[r] \cap \mathrm{ker}(\phi_q - 1), G_2 = E[r] \cap \mathrm{ker}(\phi_q - q)$ というのはなぜ?

### 3
[[Ver2008]] の 4: $v_Q$ が final exponentiation で消えるのはなぜ?

### 4
Twisting とは何?

### 5
[[Ver2008]] の 2.2: $E$ は $E[\overline {\mathbb{F}_q}]$ でよいのか? これなら $E[r] \simeq (\mathbb{Z}/r\mathbb{Z})^2$ であるという定理があったはず。

## References
[[Ver2008]]: Frederik Vercauteren. Optimal pairings. Cryptology ePrint Archive, Report 2008/096, 2008.

[[Naeh2009]]: Michael Naehrig. Constructive and Computational Aspects of Cryptographic Pairings.

[ACDFLNV2006]:  R. Avanzi, H. Cohen, C. Doche, G. Frey, T. Lange, K. Nguyen, and F. Vercauteren. Handbook of elliptic and hyperelliptic curve cryptography. Discrete Mathematics and its Applications (Boca Raton). Chapman & Hall/CRC, Boca Raton, FL, 2006.

[[FR1994]]: G. Frey and H-G. Rück. A remark concerning m-divisibility and the discrete logarithm in the divisor class group of curves. Math. Comp., 62(206):865–874, 1994.

[Ver2008]: https://eprint.iacr.org/2008/096.pdf
[Naeh2009]: https://cryptosith.org/michael/data/thesis/2009-05-13-diss.pdf
[FR1994]: https://www.ams.org/journals/mcom/1994-62-206/S0025-5718-1994-1218343-6/S0025-5718-1994-1218343-6.pdf