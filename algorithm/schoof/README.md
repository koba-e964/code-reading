有限体上の楕円曲線の点の個数の計算を行う Schoof のアルゴリズムについて書く。

# 疑問点
## 1: division polynomial が既約でないとき
楕円曲線 $E: y^{2} = x^{3}+7$ を考える。これの 3 次の division polynomial は $\psi_{3}(x) = 3x^{4} + 84x = 3x(x^3+28)$ である。これは有理数の根 $x=0$ を持つので、 ${} \bmod \psi_{3}(x)$ で考えると 0 でない多項式が逆元を持たないことがありえる。

$q := 31, E' := E(\mathbb{F}_{q})$ とする。
$P := (x,y)$ とし、Frobenius 写像 $\phi: E' \to E'$ を適用して $\psi_{3}(x)$ で mod をとった結果を見ると
- $\phi(P) = (x^{31}, y^{31}) \equiv (25x, y) \pmod{\psi_{3}(x), q}$
  - $(-28)^{10} \equiv 25 \pmod{31}$
- $\phi^2(P) = (x^{961}, y^{961}) \equiv (5x, y) \pmod{\psi_{3}(x), q}$
  - $(-28)^{320} \equiv 5 \pmod{31}$

が成立する。 $qP \equiv P = (x, y) \pmod{\psi_{3}(x), q}$ であるため、 $\phi^2(P) + qP \equiv (5x,y) + (x,y)$ の計算には x 座標が異なるが差が $\mathbb{F}_{q}[x]/\psi_{3}(x)\mathbb{F}_{q}[x]$ のゼロ因子であるような点の加算が必要である。
- naïve な実装では x 座標の差が可逆ではないという理由で無限遠点 $O$ を返してしまう。無限遠点はもちろん $0\phi(P)$ に等しいため、 $t \equiv 0 \pmod{3}$ が成立すると誤認してしまう。
- 一方、より納得感のある解釈では傾きは $(y-y) / (x-5x) = 0$ であるとみなせるので、和 $(x_3, y_3)$ は $x_3 = 0^2 - 5x-x = 25x$ (in $\mathbb{F}_{q}[x]$), $y_3 = -y$ と計算できる。 $(x_3, y_3) \equiv -\phi(P) \pmod{\psi_{3}(x), q}$ が成立するため、 $\phi$ のトレース $t$ は $t \equiv -1 \pmod{3}$ を満たす。 (実際には $t=11$ であるため正しい。)

実際には、 $\psi_l(x)$ が $\mathbb{F}_{q}[x]$ で既約でないのは the rule rather than the exception であり、 $\mathbb{F}_{q}$ に根を持つこともかなりある。例は以下:
- $\psi_3(x) = 3x(x^3 + 28)$
- $\psi_5(x) \equiv 5(x^3 + 9)(x^3 + 11)(x^3 + 21)(x^3 + 26) \pmod{31}$
- $\psi_7(x) \equiv 7(x + 4)(x + 7)(x + 20)(x^3 + 25)(x^6 + 16x^5 + 21x^4 + 18x^3 + 14x^2 + 30x + 9)(x^6 + 18x^5 + 29x^4 + 18x^3 + 8x^2 + 6x + 9)(x^6 + 28x^5 + 12x^4 + 18x^3 + 9x^2 + 26x + 9) \pmod{31}$

$\psi_l(x)$ の任意の約数 $f(x)$ に対して、 $(x, y)$ の $E(\mathbb{F}_{q}[x]/f(x)\mathbb{F}_{q}[x])$ における位数は $l$ であり、計算に代用できる。
- $(x, y)$ を $l$ 倍すると $O$ になるのは明らか。 $l$ は素数なので、 $(x, y)$ の位数は 1 または $l$ である。1 とすると $(x, y)$ が無限遠点に等しいことになり矛盾。

これを使うと先ほどの $\psi_3(x)$ の例は以下のように解釈できる:
- ${\bmod}\ x$ のとき: $(x, y) \equiv (0, \sqrt{7})$ である。これは普通に位数が 3 であるし、$(5x, y) + (x, y) \equiv 2(0, \sqrt{7})$ であるためただの二倍算である。
- ${\bmod}\ (x^3+28)$ のとき: $x-5x$ は可逆元であるため普通に計算できる。
```
  phi^2(x) + q x = [25*x, 30]
```

要するに別々の事情が 2 個重なっていただけであった。

## 2: division polynomial が有理根を持つとき
上の理屈だと $x - a$ を法として計算すればよいので計算が簡単である。
しかしこのような a は有限個しか存在しない。 $(x, y)$ は x-a を法とするとき、楕円曲線上にあるかそれの [twist](https://safecurves.cr.yp.to/twist.html) の上にあるかのどちらかである。どちらにせよ有限位数なので、 $l$ が位数を割る必要がある以上そのような $l$ は有限個しかない。

例:
楕円曲線 $E: y^{2} = x^{3}+7$ を考え、 $q := 31, E' := E(\mathbb{F}_{q})$ とする。 $t = 31 + 1 - 21 = 11$ より twist の位数は $31 + 1 + 11 = 43$ である。 $\psi_l(x)$ が有理根を持つような $l$ は 3, 7, 43 に限られる。

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

## 3: modular polynomial
modular polynomial [[KIY]] について。

# 参考文献
[[KIY]]: 小暮淳, 伊豆哲也, and 横山和弘. "Atkin, Elkies らによる Schoof のアルゴリズム改良の実装について (数式処理における理論と応用の研究)." 数理解析研究所講究録 1038 (1998): 230-243.

[[Mus2005]]: Musiker, Gregg. "Schoof's Algorithm for Counting Points on $E(\mathbb{F}_{q})$." (2005).

[KIY]: https://repository.kulib.kyoto-u.ac.jp/dspace/bitstream/2433/61961/1/1038-33.pdf

[Mus2005]: https://www-users.cse.umn.edu/~musiker/schoof.pdf
