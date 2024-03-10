# 素数位数の群を高速に実現する Ristretto255 について


$\mathcal{E}$ を群 Edwards25519 とする。
$$-x^2 + y^2 = 1-\frac{121665}{121666}x^2y^2$$
$P \in \mathcal{E}$ を $P_y = 3$ で $P_x$ が偶数である点とする。このとき、 $lP$ は位数が 8 である。よって $P$ は位数が $8l$ ($l$ は素数) であり、 $\mathcal{E}$ が位数 $8l$ の巡回群であることがわかる。
- $\mathcal{E}$ の生成元、および本当に位数が $8l$ であることについては [prove_order.log](/koba-e964/code-reading/tree/master/algorithm/ristretto255/prove_order.log) を参照すること。

$2\mathcal{E}/\mathcal{E}[4]$ は位数が素数 $l$ である。 $\mathcal{E}$ の上で計算するのはやりにくく脆弱性の原因となるので、この群の上で計算ができるようにラップするのが Ristretto [[Rist]] の役目である。ここで、記法は以下の通り。
- $r\mathcal{E} := \{rx \mid x \in \mathcal{E}\}$
- $\mathcal{E}[r] := \{x \in \mathcal{E} \mid rx = O\}$

主に [[Rist]] の流れに沿ってまとめる。また数式処理のために Sage を使う。バージョンは 10.2 である。
```sage
sage: version()
'SageMath version 10.2, Release Date: 2023-12-03'
```
## 注意点
- 曲線の限定について: ここでは $x^2$ の係数が -1 である twisted Edwards curve しか扱わないことに注意。これに当てはまるのは例えば Edwards25519 [[BDLSY2011]] で、当てはまらないのは例えば Ed448-Goldilocks [[Ham2015Ed448]] など。とくに -1 mod p が平方剰余である必要があるため、 $p \equiv 1 \pmod 4$ が要求されることに注意。
- 用語・記法:
  - $x \in \mathbb{F} _ p$ の **canonical な表現**とは、 $x$ の同値類の中で最小の非負整数である。つまり $0, 1, 2, \ldots, p-1$ のいずれかである。
  - $\mathbb{F} _ p$ の元が**正** (**positive**) であるとは、それが $0, 2, 4, \ldots, p-1$ のいずれかであることをいう。簡単にいうと canonical な表現を有理整数と解釈したものが偶数であること。
  - $\mathbb{F} _ p ^ +$ で、 $\mathbb{F} _ p$ の正の元を集めた集合を表す。つまり、 $\mathbb{F} _ p^+ := \{0, 2, 4, \ldots, p-1\}$ である。

## データ
Curve25519 [[Ber2006]], Edwards25519 [[BDLSY2011]], ristretto255 [[Rist]] の構成で使われる曲線は以下の通り。
$$p = 2^{255}-19$$
$$A = 486662$$
$$\mathcal{J}: t^2 = s^4 + As^2 + 1$$
$$\mathcal{E}: -x^2 + y^2 = 1-\frac{A-2}{A+2}x^2y^2$$
$$\mathcal{M}: v^2 = u^3 + Au^2+ u$$
データの表現は以下の通り。
- 直列化表現: $\mathcal{J}/\mathcal{J}[2] \simeq \{(s,t) \in \mathcal{J} \mid s \in \mathbb{F} _ p^+, s/t \in (1/\sqrt{-(A+2)})\mathbb{F} _ p^+\} \cap \phi^{-1}(\{(x,y) \mid xy \in \mathbb{F} _ p^+, y \neq 0\})$ (の s 座標、[参考](https://github.com/bwesterb/go-ristretto/blob/v1.2.3/edwards25519/curve.go#L230-L235))
- 内部表現: $2\mathcal{E}/\mathcal{E}[4] \simeq 2\mathcal{E} \cap \{(x,y) \mid x \in \mathbb{F} _ p^+, xy \in \mathbb{F} _ p^+, y \neq 0\}$

### 曲線
#### Jacobi quartic curve
[Wikipedia の記事](https://en.wikipedia.org/wiki/Jacobian_curve)を参考にしている。

$t^2 = s^4 + As^2 + 1$ という方程式で表せる曲線を **Jacobi quartic curve** と呼ぶ。この曲線は 2 個の無限遠点 $(\infty, \pm \infty)$ を含む。

加法公式は以下の通り。
$$s_3 = \frac{s_1t_2 + s_2t_1}{1-(s_1s_2)^2}$$
$$t_3 = \frac{(1 + (s_1s_2)^2)(t_1t_2 + As_1s_2) + 2s_1s_2(s_1^2+s_2^2)}{(1-(s_1s_2)^2)^2}$$

逆元は以下の通り。
$$(s_i, t_i) = (-s, t)$$

また単位元は $(0, 1)$ である。

$\mathcal{J}[2]$ の要素および Jacobi quartic における無限遠点の扱いについて、以下のような projective coordinate で扱うと分かりやすい。
$$Y^2 = X^4 + AX^2Z^2 + Z^4, (s, t) = (X/Z, Y/Z^2),\\
(X, Y, Z) = (aX, a^2Y, aZ), (X, Y, Z) \ne (0, 0, 0)$$

加法公式は以下の通り。
$$X_3 = X_1Y_2Z_2 + X_2Y_1Z_1$$
$$Y_3 = ((Z_1Z_2)^2 + (X_1X_2)^2)(Y_1Y_2 + AX_1X_2Z_1Z_2) \\+ 2X_1X_2Z_1Z_2(X_1^2Z_2^2+X_2^2Z_1^2)$$
$$Z_3 = (Z_1Z_2)^2 - (X_1X_2)^2$$

逆元は以下の通り。
$$(X_i, Y_i, Z_i) = (-X, Y, Z)$$

また単位元は $(0, 1, 1)$ である。

$\mathcal{J}[2] = \{(0, 1, 1), (0, -1, 1), (1, \pm 1, 0)\}$ が成立する。
なお、 $(-X, Y, -Z) = (X, Y, Z)$ に注意すること。特に $(0, 1, -1) = (0, 1, 1)$ にハマりやすい。
- $(0, 1, 1)$: 位数 1
- $(0, -1, 1), (1, \pm 1, 0)$: 位数 2

なお、projective coordinate で計算してわかることとして、t 座標は正の無限大と負の無限大を区別する意味があるが s 座標はそうではない。

#### Twisted Edwards curve
[[BDLSY2011], Section 2] に群演算が載っている。
$d := -(A-2)/(A+2)$ とおくと $d$ は ${\bmod}\ p$ で平方非剰余である必要がある。 $p = 2^{255}-19, A=486662$ では満たされている。
```
sage: p = 2**255-19
sage: A = 486662
sage: d = -GF(p)(A-2) / (A+2)
sage: is_square(d)
False
```

$\mathcal{E}[4] = \{(0, 1), (0, -1), (\pm \sqrt{-1}, 0)\}$ が成立する。ここで、
- $\sqrt{-1}$ は $-1$ の平方根であって $\mathbb{F} _ p^+$ の元であるもの、具体的には `0x2b8324804fc1df0b2b4d00993dfbd7a72f431806ad2fe478c4ee1b274a0ea0b0` を指す。(https://github.com/bwesterb/go-ristretto/blob/v1.2.3/edwards25519/field_radix51.go#L16-L20, [encoding.log](/koba-e964/code-reading/tree/master/algorithm/ristretto255/encoding.log))
- $(0, 1)$ は位数 1、 $(0, -1)$ は位数 2、 $(\pm \sqrt{-1}, 0)$ は位数 4 である。
- $(x, y) + (0, 1) = (x, y)$
- $(x, y) + (0, -1) = (-x, -y)$
  - $x \ne 0$ の場合は $x \in \mathbb{F} _ p^+$ かどうかがフリップする。 $x = 0$ の場合は $y = 1$ かどうかがフリップする。
- $(x, y) + (\pm \sqrt{-1}, 0) = (\pm \sqrt{-1}y, \pm \sqrt{-1}x)$ (複号同順)
  - $xy \ne 0$ の場合は $xy \in \mathbb{F} _ p^+$ かどうかがフリップする。 $xy = 0$ の場合は $x = 0$ かどうかや $y = 0$ かどうかがフリップする。


$2\mathcal{E}/\mathcal{E}[4]$ から $2\mathcal{E}/\mathcal{E}[2]$ への変換 (代表元の選択) を行う際は、 $xy \in \mathbb{F} _ p^+ \wedge x \neq 0$ である方を選ぶ。(そうでなければ $(\sqrt{-1}, 0)$ を足す。)

$2\mathcal{E}/\mathcal{E}[2]$ から $2\mathcal{E}$ への変換 (代表元の選択) を行う際は、 $x \in \mathbb{F} _ p^+ \wedge y \ne -1$ である方を選ぶ。(そうでなければ $(0, -1)$ を足す。)

#### Montgomery curve
[[Ber2006]] の参考文献などを参照すること。

## データ間の変換
$\mathcal{J}, \mathcal{E}, \mathcal{M}$ の間には**同種写像** (**isogeny**) が存在する。同種写像はほとんど群準同型のようなものである。

以下では $-(A+2)$ が ${\bmod}\ p$ で平方剰余であることを仮定する。 $p = 2^{255}-19, A=486662$ では満たされている。
```
sage: pow(-486664, 2**254-10, 2**255-19)
1
```

$\theta: \mathcal{J} \to \mathcal{E}$ は次数 2 の同種写像であり、
$$\theta(s, t) := \left(\sqrt{-(A+2)} \cdot \frac{s}{t}, \frac{1-s^2}{1+s^2}\right)$$
で与えられる。これの双対同種写像は以下である。

$$\hat\theta: \mathcal{E} \to \mathcal{J}$$
$$\hat\theta(x, y) := \left(\frac{2}{\sqrt{-(A+2)}} \cdot \frac{xy}{1+x^2}, \frac{y^2-x^2}{1+x^2}\right)$$

$\hat\theta \circ \theta = 2_{\mathcal{J}}$ かつ $\theta \circ \hat\theta = 2_{\mathcal{E}}$ が成立する。ここで $2_G$ は加法群 $G$ における 2 倍写像である。

$\mathcal{E}$ と $\mathcal{M}$ は同型である。同型写像は次数 1 の同種写像である。

### decoding
32 バイトのデータを $2\mathcal{E}/\mathcal{E}[4]$ の元に変換する。

コードは https://github.com/bwesterb/go-ristretto/blob/v1.2.3/edwards25519/curve.go#L224-L268 である。
1. 32 バイトのデータを受け取り、リトルエンディアンの整数として解釈して $s$ にする。 $0 \le s < p, 2|s$ が満たされているかどうか確認し、満たされていなかったらエラー。満たされていたら今後 $s$ は $\mathbb{F} _ p^+$ の元として扱う。
2. $-d(1-s^2)^2 - (1+s^2)^2$ が ${\bmod}\ p$ で平方剰余かどうか調べる。平方剰余でなかったらエラー。
  - $d = -(A-2)/(A+2)$
3. $y := (1-s^2)/(1+s^2)$ とする。また $x := 2s/\sqrt{-d(1-s^2)^2 - (1+s^2)^2}$ とする。
4. $xy \in \mathbb{F} _ p^+$ かどうか調べる。そうでなかったらエラー。
5. $(x, y)$ を出力する。

### encoding
$2\mathcal{E}/\mathcal{E}[4]$ の元 $(x, y)$ を 32 バイトのデータに変換する。

大まかに言って $2\mathcal{E}/\mathcal{E}[4] \to 2\mathcal{E}/\mathcal{E}[2] \to 2\mathcal{E} \to \mathcal{J}/\langle(0, -1)\rangle$ の変換を行う。
1. $xy \not \in \mathbb{F} _ p^+ \vee y = 0$ が成り立っていたら[^rotate]、 $(\sqrt{-1}, 0)$ を足す。つまり $(x,y)$ を $(\sqrt{-1}y, \sqrt{-1}x)$ にする。
  - 実装: https://github.com/bwesterb/go-ristretto/blob/v1.2.3/edwards25519/curve.go#L302-L309
2. $x \not \in \mathbb{F} _ p^+ \vee y = -1$ が成り立っていたら、 $(0, -1)$ を足す。つまり $(x,y)$ を $(-x, -y)$ にする。
  - 実装: https://github.com/bwesterb/go-ristretto/blob/v1.2.3/edwards25519/curve.go#L311-L314
3. $s = \pm\sqrt{(1-y)/(1+y)}$ として、$s \in \mathbb{F} _ p^+$ なる方を取る。
  - 実装: https://github.com/bwesterb/go-ristretto/blob/v1.2.3/edwards25519/curve.go#L316-L321
4. $s$ の canonical な表現をリトルエンディアンで表現し、32 バイトの列を得る。
  - 実装: https://github.com/bwesterb/go-ristretto/blob/v1.2.3/edwards25519/curve.go#L323

[^rotate]: 実装では $y = 0$ かどうかのテストをしていないが問題ない。詳しくは疑問点 4 を参照すること。
# 疑問点
## 1: それぞれの群の位数は?
実験 ([phi.log](/koba-e964/code-reading/tree/master/algorithm/ristretto255/phi.log)) によるとおそらくすべて $8l$ と思われる。 (phi.sage, phi.log) ただし、同型ではない。 $\mathcal{E}$ は巡回群であるため位数 2 の点をちょうど 1 個持つが、 $\mathcal{J}$ は 3 個持つ。

- $\mathcal{J} \simeq \mathbb{Z}/4 \times \mathbb{Z}/2 \times \mathbb{Z}/l$
- $\mathcal{J}/\mathcal{J}[2] \simeq \mathbb{Z}/2 \times \mathbb{Z}/l$
- $\mathcal{E} \simeq \mathbb{Z}/8 \times \mathbb{Z}/l$
- $2\mathcal{E} \simeq \mathbb{Z}/4 \times \mathbb{Z}/l$
- $2\mathcal{E}/\mathcal{E}[4] \simeq \mathbb{Z}/l$

## 2: φ の像は 2E の部分群なのか?
$\mathcal{J} \simeq \mathbb{Z}/4 \times \mathbb{Z}/2 \times \mathbb{Z}/l$ であるから $\mathcal{J}[4l] = \mathcal{J}$ が成立する。そのため、 $\phi(\mathcal{J}) = \phi(\mathcal{J}[4l]) \subseteq \mathcal{E}[4l] = 2\mathcal{E}$ が言える。

また、以下の事実および[同型定理](https://ja.wikipedia.org/wiki/%E5%90%8C%E5%9E%8B%E5%AE%9A%E7%90%86#%E5%AE%9A%E7%90%861)から $\mathcal{J}/\mathcal{J}[2]$ と $2\mathcal{E}/\mathcal{E}[2]$ が同型であること、同型が $\phi$ から誘導されることがわかる。
- $\phi^{-1}(\mathcal{E}[2]) = \mathcal{J}[2]$ が成立する。そのため $p \circ \phi: \mathcal{J} \to 2\mathcal{E}/\mathcal{E}[2]$ の核は $\mathcal{J}[2]$ である。
- $\phi: \mathcal{J} \to 2\mathcal{E}$ は全射である。
  - [phi_surj.log](/koba-e964/code-reading/tree/master/algorithm/ristretto255/phi_surj.log) から、 $s = 6$ を満たす点 $P \in \mathcal{J}$ に対して $\phi(P) \in 2\mathcal{E}$ は位数が $4l$ である。 $2\mathcal{E} = \langle \phi(P) \rangle$ が成立するので、任意の $Q \in 2\mathcal{E}$ に対して対して、 $Q = k\phi(P)$ なる $k$ が存在しその $k$ を使って $Q = \phi(kP)$ と書けることから、 $\phi$ は全射。

## 3: φ と φ-hat は合成すると 2 倍写像になる。直列化はどのように行うのか?
$\hat\phi$ を使うのではなく、 $\phi$ の逆像のうち適切なものを選ぶことで行われる。

## 4: encoding の際に xy が正かどうかだけ検査しても問題ないのか? xy = 0 のときは x または y の符号も見るべきなのでは?
結論としては RFC の実装に従う限り問題ない。
$xy = 0$ が満たされるとき invsqrt は 0 になるため、その他の分母は軒並み 0 になる。特に den_inv も 0 であるため s も 0 である。これは $(0, 1) \in \mathcal{J}$ および $(0, 1) \in 2\mathcal{E}$ を表現するため問題ない。
- 実装例: https://github.com/bwesterb/go-ristretto/blob/v1.2.3/edwards25519/curve.go#L302-L309
- RFC の該当箇所: [[RFC9496], 4.3.2. Encode]

# 参考文献

[[BDLSY2011]] Daniel J. Bernstein, Niels Duif, Tanja Lange, Peter Schwabe, Bo-Yin Yang. High-speed high-security signatures. Journal of Cryptographic Engineering 2 (2012), 77-89. Document ID: a1a62a2f76d23f65d622484ddd09caf8. URL: https://cr.yp.to/papers.html#ed25519.

[[Ber2006]] D. J. Bernstein. Curve25519: new Diffie-Hellman speed records. Proceedings of PKC 2006, to appear. Document ID: 4230efdfa673480fc079449d90f322c0. URL: https://cr.yp.to/papers.html#curve25519.

[[Ham2015Decaf]] Hamburg, Mike. "Decaf: Eliminating cofactors through point compression." Advances in Cryptology--CRYPTO 2015: 35th Annual Cryptology Conference, Santa Barbara, CA, USA, August 16-20, 2015, Proceedings, Part I 35. Springer Berlin Heidelberg, 2015.

[[Ham2015Ed448]] Hamburg, Mike. "Ed448-Goldilocks, a new elliptic curve." Cryptology ePrint Archive (2015).

[[RFC9496]] de Valence, H., Grigg, J., Hamburg, M., Lovecruft, I., Tankersley, G., and F. Valsorda, "The ristretto255 and decaf448 Groups", RFC 9496, DOI 10.17487/RFC9496, December 2023, <https://www.rfc-editor.org/info/rfc9496>.

[[Rist]] https://ristretto.group/ristretto.html

[BDLSY2011]: https://ed25519.cr.yp.to/ed25519-20110926.pdf

[Ber2006]: https://cr.yp.to/ecdh/curve25519-20060209.pdf

[Ham2015Decaf]: https://www.shiftleft.org/papers/decaf/decaf.pdf

[Ham2015Ed448]: https://eprint.iacr.org/2015/625

[RFC9496]: https://datatracker.ietf.org/doc/rfc9496/

[Rist]: https://ristretto.group/ristretto.html
