# Schoof のアルゴリズム
有限体上の楕円曲線の点の個数の計算を行う Schoof のアルゴリズムについて書く。
楕円曲線 $E: y^{2} = x^{3}+ax^{2} + bx + c$ を考え、それの有限体の上の還元 $E' := E(\mathbb{F} _ {q})$ を考える。 $|E'|$ を求める手法を理解するのがこのページの目標である。

数式処理のために Sage を使う。バージョンは 9.8 である。
```sage
sage: version()
'SageMath version 9.8, Release Date: 2023-02-11'
```

たちどころに明らかなことは以下である。
- $E'$ は有限集合であり、大きさは $q^2 + 1$ 以下である。
  - 無限遠点を除くと $\mathbb{F} _ {q} \times \mathbb{F} _ {q}$ の部分集合なので当然。
- $E'$ の大きさは $2q + 1$ 以下である。
  - 無限遠点を除くとそれぞれの x に対して y は高々 2 通りしかない。

$\overline{E'}$ を $\overline{E'} := E(\overline{\mathbb{F} _ {q}})$ と定義する。
Frobenius 写像 $\phi: \overline{E'} \to \overline{E'}$ を $\phi(x,y) := (x^q, y^q)$ で定義すると、 $\phi^2 -t \phi + q = 0$ が成立する有理整数 $t$ が存在する。
- この式の左辺に形式的に $\phi \leftarrow 1$ を代入して得られる式 $1-t+q$ が求める位数である。
- $E' = E(\mathbb{F} _ {q})$ の上で考えると $\phi = 1$ であるため、 $1-t+q=0$ となり (上の結果を認めれば) 当たり前であることに注意。拡大体の上で考えるべきである。
- $|t| \le 2 \sqrt{q}$ が成立する。これから、$|E'|$ は $q+1\pm 2\sqrt{q}$ の範囲内にあることがわかる。(Hasse の定理)

楕円曲線 $E$ を固定し、その定義体を $K$ とする。division polynomial $\psi_l(x)$ とは、 $l$ 分点の x 座標で消える多項式である。つまり $lP = O$ なる $P$ 全てにわたる $x-P_x$ の積の何らかのスカラー倍である。 $P_x$ そのものは $K$ の元とは限らず一般には $\overline{K}$ の元であるが、division polynomial は $K$ 係数である。

環 $\mathbb{F} _ {q}[x, y]/(y^2 - x^3 - ax^2 - bx - c, \psi_l(x))$ 上で演算を行い、上の $t$ を ${\bmod}\ l$ で求めようというのが基本方針である。より詳しくは以下。
1. $\phi^2((x, y)) + q(x, y)$ を計算する。
2. $\phi^2((x, y)) + q(x, y) = t_l\phi((x, y))$ なる $t_l$ がどれなのか、 $-(l-1)/2 \le t_l \le (l-1)/2$ の範囲で全探索する。
3. $t \equiv t_l \pmod{l}$ という情報が得られた。これを十分多数の $l$ に対して行い中国剰余定理から $t$ を求める。

詳しくは [[KIY]] などを参照されたい。

## 疑問点
### 1: division polynomial が既約でないとき
楕円曲線 $E: y^{2} = x^{3}+7$ を考える。これの 3 次の division polynomial は $\psi_{3}(x) = 3x^{4} + 84x = 3x(x^3+28)$ である。これは有理数の根 $x=0$ を持つので、 ${} \bmod \psi_{3}(x)$ で考えると 0 でない多項式が逆元を持たないことがありえる。

$q := 31, E' := E(\mathbb{F} _ {q})$ とする。
$P := (x,y)$ とし、Frobenius 写像 $\phi: \overline{E'} \to \overline{E'}$ を適用して $\psi_{3}(x)$ で mod をとった結果を見ると
- $\phi(P) = (x^{31}, y^{31}) \equiv (25x, y) \pmod{\psi_{3}(x), q}$
  - $(-28)^{10} \equiv 25 \pmod{31}$
- $\phi^2(P) = (x^{961}, y^{961}) \equiv (5x, y) \pmod{\psi_{3}(x), q}$
  - $(-28)^{320} \equiv 5 \pmod{31}$

が成立する。 $qP \equiv P = (x, y) \pmod{\psi_{3}(x), q}$ であるため、 $\phi^2(P) + qP \equiv (5x,y) + (x,y)$ の計算には x 座標が異なるが差が $\mathbb{F} _ {q}[x]/\psi_{3}(x)\mathbb{F} _ {q}[x]$ のゼロ因子であるような点の加算が必要である。
- naïve な実装では x 座標の差が可逆ではないという理由で無限遠点 $O$ を返してしまう。無限遠点はもちろん $0\phi(P)$ に等しいため、 $t \equiv 0 \pmod{3}$ が成立すると誤認してしまう。
- 一方、より納得感のある解釈では傾きは $(y-y) / (x-5x) = 0$ であるとみなせるので、和 $(x_3, y_3)$ は $x_3 = 0^2 - 5x-x = 25x$ (in $\mathbb{F} _ {q}[x]$), $y_3 = -y$ と計算できる。 $(x_3, y_3) \equiv -\phi(P) \pmod{\psi_{3}(x), q}$ が成立するため、 $\phi$ のトレース $t$ は $t \equiv -1 \pmod{3}$ を満たす。 (実際には $t=11$ であるため正しい。)

実際には、 $\psi_l(x)$ が $\mathbb{F} _ {q}[x]$ で既約でないのは the rule rather than the exception であり、 $\mathbb{F} _ {q}$ に根を持つこともかなりある。例は以下:
- $\psi_3(x) = 3x(x^3 + 28)$
- $\psi_5(x) \equiv 5(x^3 + 9)(x^3 + 11)(x^3 + 21)(x^3 + 26) \pmod{31}$
- $\psi_7(x) \equiv 7(x + 4)(x + 7)(x + 20)(x^3 + 25)(x^6 + 16x^5 + 21x^4 + 18x^3 + 14x^2 + 30x + 9)(x^6 + 18x^5 + 29x^4 + 18x^3 + 8x^2 + 6x + 9)(x^6 + 28x^5 + 12x^4 + 18x^3 + 9x^2 + 26x + 9) \pmod{31}$

$\psi_l(x)$ の任意の約数 $f(x)$ に対して、 $(x, y)$ の $E(\mathbb{F} _ {q}[x]/f(x)\mathbb{F} _ {q}[x])$ における位数は $l$ であり、計算に代用できる。
- $(x, y)$ を $l$ 倍すると $O$ になるのは明らか。 $l$ は素数なので、 $(x, y)$ の位数は 1 または $l$ である。1 とすると $(x, y)$ が無限遠点に等しいことになり矛盾。

これを使うと先ほどの $\psi_3(x)$ の例は以下のように解釈できる:
- ${\bmod}\ x$ のとき: $(x, y) \equiv (0, \sqrt{7})$ である。これは普通に位数が 3 であるし、$(5x, y) + (x, y) \equiv 2(0, \sqrt{7})$ であるためただの二倍算である。
- ${\bmod}\ (x^3+28)$ のとき: $x-5x$ は可逆元であるため普通に計算できる。
```
  phi^2(x) + q x = [25*x, 30]
```

要するに別々の事情が 2 個重なっていただけであった。

### 2: division polynomial が有理根を持つとき
ある $v \in \mathbb{F} _ {q}$ に対して $\psi_l(v) = 0$ (in $\mathbb{F} _ {q}$) が成立する場合について。
上の理屈だと $x - v$ を法として計算すればよいので計算が簡単である。
しかしこのような $v$ は有限個しか存在しない。 $(x, y)$ は $x-v$ を法とするとき、楕円曲線上にあるかそれの [twist](https://safecurves.cr.yp.to/twist.html) の上にあるかのどちらかである。元の楕円曲線も twist も有限位数なので、 $l$ が位数を割る必要がある以上そのような $l$ は有限個しかない。

例:
楕円曲線 $E: y^{2} = x^{3}+7$ を考え、 $q := 31, E' := E(\mathbb{F} _ {q})$ とする。 $t = 31 + 1 - 21 = 11$ より twist の位数は $31 + 1 + 11 = 43$ である。 $\psi_l(x)$ が有理根を持つような $l$ は 3, 7, 43 に限られる。

同じような事情が任意の有限次拡大についていえる。 $l=5$ のとき $\psi_5$ は 3 次の既約多項式の積に分解されるのであった。以下の計算によると、 $E(\mathbb{F} _ {31^3})$ の twist は位数が 5 の倍数であるため、当然位数 5 の点を含む。

```console
sage: q = 31^3
sage: E = EllipticCurve(GF(q), [0, 0, 0, 0, 7])
sage: len(E.rational_points())
29484
sage: factor(len(E.rational_points()))
2^2 * 3^4 * 7 * 13
sage: 31^3 + 1 - 29484
308
sage: 31^3 + 1 + 308
30100
sage: factor(31^3 + 1 + 308)
2^2 * 5^2 * 7 * 43
```

# modular polynomial を用いた高速化
modular polynomial [[KIY]] について。

実数体上で議論を行うため、[Weierstraß の 楕円関数](https://ja.wikipedia.org/w/index.php?title=%E3%83%B4%E3%82%A1%E3%82%A4%E3%82%A8%E3%83%AB%E3%82%B7%E3%83%A5%E3%83%88%E3%83%A9%E3%82%B9%E3%81%AE%E6%A5%95%E5%86%86%E5%87%BD%E6%95%B0&oldid=95955129)から逃れることはできない。

素数 $l$ に対する modular polynomial $\Phi_l(x, y)$ は、 $x$ および $y$ に関して $l+1$ 次の多項式である。具体的な値は [[CMP]] で計算されている。
modular polynomial には以下の著しい性質がある。
- 以下の等式が成立する:
  $$\Phi_l(x, j(\tau)) = (x-j(l\tau)) \prod_{0 \le i < l} (x-j\left(\frac{\tau+i}{l}\right))$$
-  $\Phi_l(x, y)$ の係数は有理整数なので、 $\Phi_l(x, j(\tau))$ の係数は $\mathbb{Z}[j(\tau)]$ の元である。

また、 $j \neq 0, 1728$ のとき、以下の 3 つは同値であるようである。 ([[Sch1995], Prop. 6.2] など)
- $\phi$ の特性多項式が因数分解できること ($t^2-4q$ が ${\bmod}\ l$ の平方剰余であること)
- modular polynomial が有理根を持つこと
- division polynomial が $(l-1)/2$ 次の因子を持つこと

これらいずれか一つ (つまり全て) を満たすときには division polynomial の因子を見つけることができ、上の Schoof のアルゴリズムの ${\bmod}\ l$ の部分を高速化できる。基本的な流れは以下。

1. $E$ と isogenous である曲線 $\tilde{E}$ の j-不変量 $\tilde \jmath$ を求める。
    - $E$ の $\tau$ を $\tau$ と置くと、 $\tilde \jmath = j(l\tau) \vee \exists i\ldotp \tilde \jmath = j((\tau+i)/l)$ が成立する。これらは modular polynomial の有理根として求められる。
2. $\tilde \jmath$ から曲線 $\tilde{E}$ を復元する。
3. $F_l$ を求める。 $F_l(x) = \prod_{1 \le i \le (l-1)/2} (x-\wp(i\tau_1/l))$ は有理整数の係数を持ち、次数は $(l-1)/2$ であり、division polynomial の因子である。 $F_l(\wp(z))$ は関数の間のある等式を満たすので、そこから $F_l$ の係数を求める。

また $j = 0, 1728$ のときには Cornacchia's algorithm が適用できるようである。([[Sch1995], Section 4])

詳しくは [[Sch1995]] を参照されたい。

## 疑問点
### modular polynomial に j(\tau) を代入する方法
$y^2 = x^3 + bx + c$ に対して、 $j = 1728 \cdot 4b^3/(4b^3+27c^2)$ である。楕円曲線が $\mathbb{F} _ {q}$ 上で定義されていればこれは普通に $\mathbb{F} _ {q}$ の元なので、特に対処に困ることはない。
なお、 $\tau$ や nome $q = \exp(2\pi \sqrt{-1} \tau)$ を陽に計算することはないことに注意。


# 参考文献
[[CMP]] https://math.mit.edu/~drew/ClassicalModPolys.html

[[KIY]] 小暮淳, 伊豆哲也, and 横山和弘. "Atkin, Elkies らによる Schoof のアルゴリズム改良の実装について (数式処理における理論と応用の研究)." 数理解析研究所講究録 1038 (1998): 230-243.

[[Mus2005]] Musiker, Gregg. "Schoof's Algorithm for Counting Points on $E(\mathbb{F} _ {q})$." (2005).

[[Sch1995]] Schoof, René. "Counting points on elliptic curves over finite fields." Journal de théorie des nombres de Bordeaux 7.1 (1995): 219-254.

[CMP]: https://math.mit.edu/~drew/ClassicalModPolys.html

[KIY]: https://repository.kulib.kyoto-u.ac.jp/dspace/bitstream/2433/61961/1/1038-33.pdf

[Mus2005]: https://www-users.cse.umn.edu/~musiker/schoof.pdf

[Sch1995]: http://www.numdam.org/item/JTNB_1995__7_1_219_0/
