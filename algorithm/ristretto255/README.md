# Ristretto255 について


$\mathcal{E}$ を群 Edwards25519 とする。
$$-x^2 + y^2 = 1-\frac{121665}{121666}x^2y^2$$
$P \in \mathcal{E}$ を $P_y = 3$ で $P_x$ が偶数である点とする。このとき、 $lP$ は位数が 8 である。よって $P$ は位数が $8l$ ($l$ は素数) であり、 $\mathcal{E}$ が位数 $8l$ の巡回群であることがわかる。
- $\mathcal{E}$ の生成元、および本当に位数が $8l$ であることについては [prove_order.log](./prove_order.log) を参照すること。

$2\mathcal{E}/\mathcal{E}[4]$ は位数が素数 $l$ である。 $\mathcal{E}$ の上で計算するのはやりにくく脆弱性の原因となるので、この群の上で計算ができるようにラップするのが Ristretto [[Rist]] の役目である。ここで、記法は以下の通り。
- $r\mathcal{E} := \{rx \mid x \in \mathcal{E}\}$
- $\mathcal{E}[r] := \{x \in \mathcal{E} \mid rx = O\}$

主に [[Rist]] の流れに沿ってまとめる。また数式処理のために Sage を使う。バージョンは 10.2 である。
```sage
sage: version()
'SageMath version 10.2, Release Date: 2023-12-03'
```

ここでは $x^2$ の係数が -1 である twisted Edwards curve しか扱わないことに注意。これに当てはまるのは例えば Edwards25519 [[BDLSY2011]] で、当てはまらないのは例えば Ed448-Goldilocks [[Ham2015Ed448]] など。

## データ
Curve25519 [[Ber2006]], Edwards25519 [[BDLSY2011]], ristretto255 [[Rist]] の構成で使われる曲線は以下の通り。
$$p = 2^{255}-19$$
$$A = 486662$$
$$\mathcal{J}: t^2 = s^4 + As^2 + 1$$
$$\mathcal{E}: -x^2 + y^2 = 1-\frac{A-2}{A+2}x^2y^2$$
$$\mathcal{M}: v^2 = u^3 + Au^2+ u$$
データの表現は以下の通り。
- 直列化表現: $\mathcal{J}/\mathcal{J}[2]$ (の s 座標、[参考](https://github.com/bwesterb/go-ristretto/blob/v1.2.3/edwards25519/curve.go#L230-L235))
- 内部表現: $2\mathcal{E}/\mathcal{E}[4]$

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
[[BDLSY2011]] に群演算が載っている。
$d := -(A-2)/(A+2)$ とおくと $d$ は ${\bmod}\ p$ で平方非剰余である必要がある。 $p = 2^{255}-19, A=486662$ では満たされている。
```
sage: p = 2**255-19
sage: A = 486662
sage: d = -GF(p)(A-2) / (A+2)
sage: is_square(d)
False
```


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
$$\theta(s, t) := \left(\frac{\sqrt{-(A+2)}}{2} \cdot \frac{2s}{t}, \frac{1-s^2}{1+s^2}\right)$$
で与えられる。これの双対同種写像は以下である。

$$\hat\theta: \mathcal{E} \to \mathcal{J}$$
$$\hat\theta(x, y) := \left(\frac{2}{\sqrt{-(A+2)}} \cdot \frac{xy}{1+x^2}, \frac{y^2-x^2}{1+x^2}\right)$$

$\hat\theta \circ \theta = 2_{\mathcal{J}}$ かつ $\theta \circ \hat\theta = 2_{\mathcal{E}}$ が成立する。ここで $2_G$ は加法群 $G$ における 2 倍写像である。

$\mathcal{E}$ と $\mathcal{M}$ は同型である。同型写像は次数 1 の同種写像である。

### decoding
32 バイトのデータを $2\mathcal{E}/\mathcal{E}[4]$ の元に変換する。

TODO

### encoding
$2\mathcal{E}/\mathcal{E}[4] \to 2\mathcal{E}/\mathcal{E}[2] \to 2\mathcal{E} \to \mathcal{J}$
TODO

# 疑問点
## それぞれの群の位数は?
実験 ([phi.log](./phi.log)) によるとおそらくすべて $8l$ と思われる。 (phi.sage, phi.log) ただし、同型ではない。 $\mathcal{E}$ は巡回群であるため位数 2 の点をちょうど 1 個持つが、 $\mathcal{J}$ は 3 個持つ。

- $\mathcal{J} \simeq \mathbb{Z}/4 \times \mathbb{Z}/2 \times \mathbb{Z}/l$
- $\mathcal{J}/\mathcal{J}[2] \simeq \mathbb{Z}/2 \times \mathbb{Z}/l$
- $\mathcal{E} \simeq \mathbb{Z}/8 \times \mathbb{Z}/l$
- $2\mathcal{E} \simeq \mathbb{Z}/4 \times \mathbb{Z}/l$
- $2\mathcal{E}/\mathcal{E}[4] \simeq \mathbb{Z}/l$


## φ の像は 2E の部分群なのか?
$\mathcal{J} \simeq \mathbb{Z}/4 \times \mathbb{Z}/2 \times \mathbb{Z}/l$ であるから $\mathcal{J}[4l] = \mathcal{J}$ が成立する。そのため、 $\phi(\mathcal{J}) = \phi(\mathcal{J}[4l]) \subseteq \mathcal{E}[4l] = 2\mathcal{E}$ が言える。

また、 $\phi(\mathcal{J}[2]) = \{(x, y) \in \mathcal{E} \mid (x, y) = (0, 1), (0, -1)\} = \mathcal{E}[2]$ が成立する。

## φ と φ-hat は合成すると 2 倍写像になる。直列化はどのように行うのか?
$\hat\phi$ を使うのではなく、 $\phi$ の逆像のうち適切なものを選ぶことで行われる。

# 参考文献

[[BDLSY2011]] Daniel J. Bernstein, Niels Duif, Tanja Lange, Peter Schwabe, Bo-Yin Yang. High-speed high-security signatures. Journal of Cryptographic Engineering 2 (2012), 77-89. Document ID: a1a62a2f76d23f65d622484ddd09caf8. URL: https://cr.yp.to/papers.html#ed25519.

[[Ber2006]] D. J. Bernstein. Curve25519: new Diffie-Hellman speed records. Proceedings of PKC 2006, to appear. Document ID: 4230efdfa673480fc079449d90f322c0. URL: https://cr.yp.to/papers.html#curve25519.

[BDLSY2011]: https://ed25519.cr.yp.to/ed25519-20110926.pdf

[Ber2006]: https://cr.yp.to/ecdh/curve25519-20060209.pdf

[Ham2015Decaf]: https://www.shiftleft.org/papers/decaf/decaf.pdf

[Ham2015Ed448]: https://eprint.iacr.org/2015/625

[Rist]: https://ristretto.group/ristretto.html
