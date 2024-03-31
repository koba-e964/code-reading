# モジュラー曲線について

数式処理のために Sage を使う。バージョンは 10.2 である。
```sage
sage: version()
'SageMath version 10.2, Release Date: 2023-12-03'
```

## 注意点
- 用語・記法:
  - $j$: [$j$-不変量](https://en.wikipedia.org/wiki/J-invariant)を指す。

# 概要
**古典的モジュラー多項式** (classical modular polynomial) とは、 $\Phi_N(j(N\tau), j(\tau)) = 0$ が満たされる多項式 $\Phi_N(x,y)$ である。 $\Phi_N(x,y)$ は対称多項式である。方程式 $\Phi_N(x, y) = 0$ で定義される曲線を**古典的モジュラー曲線** (classical modular curve) と呼ぶ。

# 実験 1
$\Phi_2(x, y) = (x^3 + y^3)  -162000(x^2+y^2) + 1488xy(x+y) - x^2y^2 + 8748000000(x+y) + 40773375xy - 157464000000000$ である [[CMP]]。

$\Phi_2(j(2\tau), j(\tau)) = 0$ が成立することを見る。そのため、 $\tau = (1+\sqrt{3})/2$ を代入する。

$j(\tau) = 0$ であるため、 $\Phi_2(x, j(\tau)) = \Phi_2(x, 0) = x^3 - 162000x^2 + 8748000000x - 157464000000000 = (x-54000)^3$ が成立する。

ここで、 $j(\sqrt{-3}) = 54000$ である。計算は以下の通り。
```python
sage: E = EllipticCurve_from_j(54000)
sage: E.j_invariant()
54000
sage: emb = QQ.embeddings(ComplexField())[0]
sage: L = E.period_lattice(emb)
sage: L.tau()
1.73205080756888*I
sage: L.tau()^2
-3.00000000000000
```

直接的には以下のように計算できる。([以下のようにしか計算できない。](https://ask.sagemath.org/question/61800/compute-the-equation-for-elliptic-curves-from-lattices/)) $q$ についてのローラン展開を使う都合上、 $|q| = |\exp(2\pi\sqrt{-1} \tau)|$ が小さいほど、つまり $\mathop{\mathrm{Im}} \tau$ が大きいほど精度が良い。
```python
sage: def j(tau):
....:     return j_invariant_qexp(prec=100).laurent_polynomial()(exp(2*pi*i * tau)).n(prec=100)
....: 
sage: j(sqrt(3)*i)
54000.000000000000000019205480
sage: j((1+sqrt(3)*i)/2) # actual value = 0
-5.3221448610759974101135450830e261 - 9.0246840405273862296685225022e230*I
```

ともかく、 $j(2\tau) = j(\sqrt{-3}) = 54000$ は $\Phi_2(x, j(\tau)) = 0$ の根であるため、 $\Phi_2(j(2\tau), j(\tau)) = 0$ が成立する。なお残りの根は $j(\tau/2)$ と $j((\tau+1)/2)$ である ([[KIY]] 2.2.1) が、これらはすべて $j(\sqrt{-3})$ に等しいので $\Phi_2(x, j(\tau)) = 0$ が重根を持つことと整合する。
- $j(\tau/2) = j((1+\sqrt{-3})/4) = j(-4/(1+\sqrt{-3})) = j(-1+\sqrt{-3}) = j(\sqrt{-3})$
- $j((\tau+1)/2) = j((-1+\sqrt{-3})/4) = j(-4/(-1+\sqrt{-3})) = j(1+\sqrt{-3}) = j(\sqrt{-3})$

# 実験 2
$\Phi_2(x, y) = (x^3 + y^3)  -162000(x^2+y^2) + 1488xy(x+y) - x^2y^2 + 8748000000(x+y) + 40773375xy - 157464000000000$ であった [[CMP]]。

$j(\tau)$ と $j(2\tau)$ の q-展開を考えると、 $j(2\tau)$ の方は $j(\tau)$ に登場する $q$ に $q^2$ を代入したものである。 ($q(2\tau) = \exp(2\pi\sqrt{-1}\cdot 2\tau) = q^2$ から)

このことを利用し、 $\Phi_2(j(2\tau), j(\tau)) = 0$ を以下のように検証できる。

```python
sage: j1 = j_invariant_qexp(prec = 10)
sage: j2 = j1.V(2)
sage: j1
q^-1 + 744 + 196884*q + 21493760*q^2 + 864299970*q^3 + 20245856256*q^4 + 333202640600*q^5 + 4252023300096*q^6 + 44656994071935*q^7 + 401490886656000*q^8 + 3176440229784420*q^9 + O(q^10)
sage: j2
q^-2 + 744 + 196884*q^2 + 21493760*q^4 + 864299970*q^6 + 20245856256*q^8 + 333202640600*q^10 + 4252023300096*q^12 + 44656994071935*q^14 + 401490886656000*q^16 + 3176440229784420*q^18 + O(q^20)
sage: (j1^3+j2^3) - 162000*(j1^2+j2^2) + 1488*j1*j2*(j1+j2) - j1^2*j2^2 + 8748000000*(j1+j2) + 40773375*j1*j2 - 157464000000000
O(q^5)
```

# 疑問点
## 1: モジュラー多項式の計算法は?
TODO

## 2: [モジュラー性定理](https://ja.wikipedia.org/wiki/%E8%B0%B7%E5%B1%B1%E2%80%93%E5%BF%97%E6%9D%91%E4%BA%88%E6%83%B3)との関連は?
TODO

# 参考文献
[[BLS]] Bröker, Reinier, Kristin Lauter, and Andrew Sutherland. "Modular polynomials via isogeny volcanoes." Mathematics of Computation 81.278 (2012): 1201-1231.

[[CMP]] https://math.mit.edu/~drew/ClassicalModPolys.html

[[KIY]] 小暮淳, 伊豆哲也, and 横山和弘. "Atkin, Elkies らによる Schoof のアルゴリズム改良の実装について (数式処理における理論と応用の研究)." 数理解析研究所講究録 1038 (1998): 230-243.

[BLS]: https://arxiv.org/abs/1001.0402

[CMP]: https://math.mit.edu/~drew/ClassicalModPolys.html

[KIY]: https://repository.kulib.kyoto-u.ac.jp/dspace/bitstream/2433/61961/1/1038-33.pdf
