有限体上の楕円曲線の点の個数の計算を行う Schoof のアルゴリズムについて書く。

# 疑問点
楕円曲線 $E: y^{2} = x^{3}+7$ を考える。これの 3 次の division polynomial は $\psi_{3}(x) = 3x^{4} + 84x$ である。これは有理数の根 $x=0$ を持つので、 ${} \bmod \psi_{3}(x)$ で考えると 0 でない多項式が逆元を持たないことがありえる。

$q := 31, E' := E(\mathbb{F}_{q})$ とする。
$P := (x,y)$ とし、Frobenius 写像 $\phi: E' \to E'$ を適用して $\psi_{3}(x)$ で mod をとった結果を見ると
- $\phi(P) = (x^{31}, y^{31}) \equiv (25x, y) \pmod {\psi_{3}(x)}$
- $\phi^2(P) = (x^{961}, y^{961}) \equiv (5x, y) \pmod {\psi_{3}(x)}$

が成立する。 $qP \equiv P = (x, y) \pmod {\psi_{3}(x)}$ であるため、 $\phi^2(P) + qP \equiv (5x,y) + (x,y)$ の計算には x 座標が異なるが差が $\mathbb{F}_{q}[x]/\psi_{3}(x)\mathbb{F}_{q}[x]$ のゼロ因子であるような点の加算が必要である。
- naïve な実装では x 座標の差が可逆ではないという理由で無限遠点を返してしまう。無限遠点はもちろん $0\phi(P)$ に等しいため、 $t \equiv 0 \pmod{3}$ が成立すると誤認してしまう。
- 一方、より納得感のある解釈では傾きは $(y-y) / (x-5x) = 0$ であるとみなせるので、和 $(x_3, y_3)$ は $x_3 = 0^2 - 5x-x = 25x$ (in $\mathbb{F}_{q}[x]$), $y_3 = -y$ と計算できる。 $(x_3, y_3) \equiv -\phi(P) \pmod {\psi_{3}(x)}$ が成立するため、 $\phi$ のトレース $t$ は $t \equiv -1 \pmod{3}$ を満たす。 (実際には $t=11$ であるため正しい。)


# 参考文献
[[Mus2005]]: Musiker, Gregg. "Schoof's Algorithm for Counting Points on $E(\mathbb{F}_{q})$." (2005).

[Mus2005]: https://www-users.cse.umn.edu/~musiker/schoof.pdf
