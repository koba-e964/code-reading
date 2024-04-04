# モジュラー曲線について


## 注意点
数式処理のために Sage を使う。バージョンは 10.2 である。
```python
sage: version()
'SageMath version 10.2, Release Date: 2023-12-03'
```

また流儀は以下の通り。
- 用語・記法:
  - $j$: [$j$-不変量](https://en.wikipedia.org/wiki/J-invariant)を指す。

# 概要
**古典的モジュラー多項式** (classical modular polynomial) とは、 $\Phi_N(j(N\tau), j(\tau)) = 0$ が満たされる多項式 $\Phi_N(x,y)$ である [[BLS2012]]。 $\Phi_N(x,y)$ は対称多項式である。方程式 $\Phi_N(x, y) = 0$ で定義される曲線を**古典的モジュラー曲線** (classical modular curve) と呼ぶ。

# 実験
## 実験 1
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

ともかく、 $j(2\tau) = j(\sqrt{-3}) = 54000$ は $\Phi_2(x, j(\tau)) = 0$ の根であるため、 $\Phi_2(j(2\tau), j(\tau)) = 0$ が成立する。なお残りの根は $j(\tau/2)$ と $j((\tau+1)/2)$ である ([[KIY1998]] 2.2.1) が、これらはすべて $j(\sqrt{-3})$ に等しいので $\Phi_2(x, j(\tau)) = 0$ が重根を持つことと整合する。
- $j(\tau/2) = j((1+\sqrt{-3})/4) = j(-4/(1+\sqrt{-3})) = j(-1+\sqrt{-3}) = j(\sqrt{-3})$
- $j((\tau+1)/2) = j((-1+\sqrt{-3})/4) = j(-4/(-1+\sqrt{-3})) = j(1+\sqrt{-3}) = j(\sqrt{-3})$

## 実験 2
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

## 実験 3
$N$ が素数でない場合の検証のため、 $N=4$ で試す。

phi4.sage を実行した結果、 $\Phi(j(\tau), j(\tau+1/2)) = 0$ であることが分かる。

```
$ sage phi4.sage
Phi_4(j(q), j(-q)) = O(q^4)
```

# 疑問点
## 1: モジュラー多項式の計算法は?
誰でも思いつきそうな方法として、 $j(\tau)$ と $j(N\tau)$ の q-展開を比較して係数を合わせる方法 ([[BCRS1999]] など) がある。

$j(\tau)$ の q-展開が十分な長さ必要である。これは例えば https://mathoverflow.net/questions/71704/computing-the-q-series-of-the-j-invariant などで計算できる。ファイル jinv.go で実装した。

よりよい手法は [[BLS2012]] などを参照されたい。

## 2: [モジュラー性定理](https://ja.wikipedia.org/wiki/%E8%B0%B7%E5%B1%B1%E2%80%93%E5%BF%97%E6%9D%91%E4%BA%88%E6%83%B3)との関連は?

モジュラー性定理は、以下のような定理である。
> 有理数体 $\mathbb{Q}$ 上のすべての楕円曲線はモジュラーである。

楕円曲線が**モジュラーである**という単語の定義には以下のような流儀がある。
- 古典的モジュラー曲線 $\Phi_N(x, y) = 0$ からの全射が存在する。([Wikipedia](https://en.wikipedia.org/w/index.php?title=Modularity_theorem&oldid=1212586185))
- あるモジュラー形式 $f$ が存在し (**付随する**モジュラー形式とよばれる) q-展開の係数を $f = \sum a_n(f) q^n$ としたとき $a_p = p+1 - |E(\mathbb{F} _ p)|$ が成り立つ。[[DdB2011]]
- L 関数を使ったもの TODO

TODO
全射の構成方法とか、具体例とかに触れたい
tsujimotter 氏のブログ ([モジュラー曲線(2)：合同部分群とモジュラー方程式 - tsujimotterのノートブック](https://tsujimotter.hatenablog.com/entry/modular-curve-2)) にない例を取り上げる。

## 3: ある種の j-不変量が代数的整数であることについて
$E$ が虚数乗法を持つとき (つまり、 $\tau$ が虚二次体 $K=\mathbb{Q}(\sqrt{-D})$) の有理数でない元であるとき) には、 $j(\tau)$ は代数的整数である。
- 略証: $\sqrt{-D}\tau = a\tau+b, \sqrt{-D} = c\tau + d$ とする。行列 $\alpha$ を $\alpha = (a,b;c,d)$ で定めると、 $j(\alpha\tau) = j(\tau)$ であり $0 = \Phi_{D}(j(\alpha\tau), j(\tau)) = \Phi_{D}(j(\tau), j(\tau))$ である。 $\Phi_{D}(x, x)$ はモニックで有理整数係数の多項式であるため、その根たる $j(\tau)$ は代数的整数である。
  - 詳しくは [[Sil1994]] の Cor. II.6.3.1 を参照されたい。

## 4: 商空間を曲線と呼んでいいのか?
- 商空間
- コンパクトリーマン面

TODO: わかりやすい文献

# 参考文献
[[BCRS1999]] Blake, Ian F., et al. "On the computation of modular polynomials for elliptic curves." HP Laboratories Technical Report, to appear (1999).

[[BLS2012]] Bröker, Reinier, Kristin Lauter, and Andrew Sutherland. "Modular polynomials via isogeny volcanoes." Mathematics of Computation 81.278 (2012): 1201-1231.

[[CMP]] https://math.mit.edu/~drew/ClassicalModPolys.html

[[DdB2011]] TODO

[[KIY1998]] 小暮淳, 伊豆哲也, and 横山和弘. "Atkin, Elkies らによる Schoof のアルゴリズム改良の実装について (数式処理における理論と応用の研究)." 数理解析研究所講究録 1038 (1998): 230-243.

[[Sil1994]] TODO

[BLS2012]: https://arxiv.org/abs/1001.0402

[CMP]: https://math.mit.edu/~drew/ClassicalModPolys.html

[KIY1998]: https://repository.kulib.kyoto-u.ac.jp/dspace/bitstream/2433/61961/1/1038-33.pdf

[Sil1994]: https://link.springer.com/book/10.1007/978-1-4612-0851-8

[a]: https://math.mit.edu/classes/18.783/2017/LectureNotes25.pdf

[DdB2011]: https://www.universiteitleiden.nl/binaries/content/assets/science/mi/scripties/dobbendebruynbach.pdf

[BCRS1999]: https://www.math.uwaterloo.ca/~mrubinst/publications/phi.pdf
