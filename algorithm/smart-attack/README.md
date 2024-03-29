# Smart-attack
このページでは、ある種の楕円曲線に対する攻撃手法を扱う。この攻撃手法は Smart [[Sma1999]] により提案され、俗な呼び方として以下がある。
- Smart-ASS attack (https://safecurves.cr.yp.to/transfer.html)
- Additive transfer (https://safecurves.cr.yp.to/transfer.html)

この攻撃が適用できる楕円曲線は、実用上まったく役に立たない。

## 準備
Weierstrass 型の楕円曲線とは以下のような式で定義される曲線である:
$$E\colon y^2 + a_1xy + a_3y = x^3 + a_2x^2 + a_4x + a_6$$
有限体 $\mathbb{F} _ q$ 上で定義されたそれを $E(\mathbb{F} _ q)$ と表記する。

以下の Hasse の定理が成立する:
$$|E(\mathbb{F} _ q)| = q - t + 1$$

ここで $t$ はフロべニウス写像 $\phi$ の trace である。つまり $\phi^2 - t\phi + p = 0$ を満たす数 $t$ である。

この記事では $q$ は素数の場合、 $q = p$ の場合だけを扱う。

楕円曲線は非特異であるもののみを扱う。これは以下の条件と同値であるはず。
- 定義式を $f(x, y) = 0$ in $\mathbb{F} _ p$ とするとき、楕円曲線の上で $(f_x(x, y), f_y(x, y)) \neq (0, 0)$ in $\mathbb{F} _ p$ が成立する。

以降点の座標としては $\mathbb{F} _ p$ や $\mathbb{Z} _ p$ の元を考える。

## 無限遠点の周りの展開
[Sil2016] の Section IV.1 に沿って議論する。

基本的に話の流れはある環 $R$ に対して $pR/p^2R$ の元に対して計算をすることなので、マクローリン級数を扱う時は 2 次の近似をすることになる。この節では $R = \mathbb{F} _ p$ か $R = \mathbb{Z} _ p$ とする。(どちらでも問題ない議論をする)

Weierstrass 型の楕円曲線の方程式は以下だった:
$$y^2 + a_1xy + a_3y = x^3 + a_2x^2 + a_4x + a_6$$

これに対し $(z, w) = (-x/y, -1/y)$ ( $(x, y) = (z/w, -1/w)$ ) という変換を行うと、以下のようになる。
$$w^{-2}(1 - a_1z - a_3w) = w^{-3}(z^3 + a_2z^2w + a_4zw^2 + a_6w^3)$$
$$w - a_1zw - a_3w^2 = z^3 + a_2z^2w + a_4zw^2 + a_6w^3$$
$$w = z^3 + a_1zw + a_2z^2w + a_3w^2 + a_4zw^2 + a_6w^3$$

これにより $w \in z^3 + a_1z^4 + (z^5)$ がわかる。
この式で $(z, w) = (0, 0)$ は完全に妥当であり、これが無限遠点に相当する。

楕円曲線が非特異であるという条件は、zw 座標でも微分係数が非ゼロという条件と同値であるはず。(未検証)
具体的には $g(z, w) := z^3 + a_1zw + a_2z^2w + a_3w^2 + a_4zw^2 + a_6w^3 - w$ としたとき、 $g(z, w) = 0$ in $\mathbb{F} _ p$ ならば $(g_z(z, w), g_w(z, w)) \neq (0, 0)$ in $\mathbb{F} _ p$ であるという条件である。

楕円曲線上のアーベル群の演算を復習しよう。
- 単位元はもちろん無限遠点なので、 $(z, w) = (0, 0)$ である。
- 逆元について、直線 $x = x_1$ と曲線の共有点が重要である。一つの共有点が $(x_1, y_1)$ の場合、二次方程式の根と係数の関係からもう一つの共有点は $(x_1, -y_1 - a_1x_1 - a_3)$ であり、これが $(x_1, y_1)$ の逆元なのであった。
これを zw-座標で計算するとどうなるか? xy-座標で計算すると $(x_1, y_1) = (z_1/w_1, -1/w_1)$ の逆元は $(x_1, -y_1 - a_1x_1 - a_3) = (z_1/w_1, (1 - a_1z_1 - a_3w_1)/w_1)$ なので、これを zw-座標に直すと $(-z_1/(1 - a_1z_1 - a_3w_1), -w_1/(1 - a_1z_1 - a_3w_1))$ である。z 座標を $i(z_1)$ と表記すると $p | z_1$ であれば $i(z_1) \in -z_1 + p^2R$ であることに注意。
- 群の演算について、xy-平面での直線は zw-平面でも直線であることに注意すること。 $ax + by + c = 0$ は $az/w - b/w + c = 0$ に変形でき、これは $az + cw - b = 0$ に変形できるため。
$(z_1, w_1)$ と $(z_2, w_2)$ を通る直線を $w = \lambda z + \nu$ としたとき、これと $w = z^3 + a_1zw + a_2z^2w + a_3w^2 + a_4zw^2 + a_6w^3$ の共有点を求めたい。w に代入すると z に関する 3 次方程式が得られるので、3 次と 2 次の係数に対し根と係数の関係を使う。
もうひとつの共有点を $(z_3, w_3)$ と置くと
$z_3 = -z_1 - z_2 - (a_1\lambda + a_2\nu + a_3 \lambda^2 + 2a_4\lambda\nu + 3a_6\lambda^2\nu)/(1 + a_2\lambda + a_4\lambda^2 + a_6\lambda^3) = -z_1 - z_2 + O(z_1^2, z_1z_2, z_2^2)$
が成立する。 $(z_1, w_1)$ と $(z_2, w_2)$ の和の z 座標は $i(z_3)$ であるため、 $F(z_1, z_2) = i(z_3)$ である。
$p | \lambda^{-1}$ の場合 (直線が $w$ 軸に平行な場合も含む) は $z = \lambda^{-1}w + \mu$ という方程式を使って $w$ に関する三次方程式を解く必要があるが、それ以外は同じである。 
  - $p | z_1, z_2$ のとき、 $z_3 \in -z_1 - z_2 + p^2 R$ が成立する。 $(z_1, w_1)$ と $(z_2, w_2)$ の和の z 座標は $i(z_3)$ であるため、2 次以上の項を無視すると $F(z_1, z_2) = i(z_3) \in z_1 + z_2 + p^2 R$ が成立する。(TODO: $\lambda$ が p の負のベキを含む場合でも同じか調べる)
  - 同じ点の加算は接線を見る必要がある。 $(g_z(z, w), g_w(z, w)) \neq 0$ in $\mathbb{F} _ p$ だったので、 $g_z(z, w) \neq 0$ in $\mathbb{F} _ p$ であれば zw-平面の傾き $-g_w/g_z$ の直線を、 $g_w(z, w) \neq 0$ in $\mathbb{F} _ p$ であれば wz-平面の傾き $-g_z/g_w$ の直線を考えれば良い。

## lifting
$\mathbb{F} _ p$ 上で定義された zw-平面の楕円曲線 $E(\mathbb{F} _ p): w = z^3 + a_1zw + a_2z^2w + a_3w^2 + a_4zw^2 + a_6w^3$ に対し、同じ式で $\mathbb{Z} _ p$ 上で定義された曲線 $E(\mathbb{Z} _ p)$ を考えることができる。 $(Z, W) \in E(\mathbb{F} _ p)$ のとき、 $(Z', W') \in E(\mathbb{Z} _ p)$ であって $Z' \equiv Z, W' \equiv W \pmod p$ であるようなものが存在する。これを lifting と呼ぶことにする。
- 計算方法: [Hensel lifting](https://en.wikipedia.org/wiki/Hensel%27s_lemma) を用いる。 $\mathbb{Z} _ p$ 上で計算する。非特異性の仮定から $g_z(Z, W) \not \equiv 0 \pmod p$ と $g_w(Z, W) \not \equiv 0 \pmod p$ の少なくともどちらか一方は成り立つのだった。 $g_z(Z, W) \not \equiv 0 \pmod p$ の場合、 $Z' = Z - g(Z, W) / g_z(Z, W), W' = W$ とすれば $g(Z', W') \equiv 0 \pmod{p^2}$ が成立する。 $g_w(Z, W) \not \equiv 0 \pmod p$ の場合は $Z' = Z, W' = W - g(Z, W) / g_w(Z, W)$ とすれば $g(Z', W') \equiv 0 \pmod{p^2}$ が成立する。これの繰り返しで任意の k に対して mod $p^k$ で正しい値を計算することができ、これの極限を取った結果が $E(\mathbb{Z} _ p)$ の点となるが、今回は mod $p^2$ で十分である。

## 攻撃手法
[[Sma1999]] より。 $q = p$ として $|E(\mathbb{F} _ p)| = p$ 、つまり $t = 1$ であるような楕円曲線を考える。

1. $P$ と $Q$ は mod p で表現されているが、それを $\mathbb{Z} _ p$ の元に変換する。(実装上は mod $p^2$ でよい。) その後それらを zw-座標に変換する。(y 座標が mod p で非ゼロであったことに注意すると、これらの z 座標はかならず $\mathbb{Z} _ p$ の元である。)
2. $pP$ と $pQ$ を計算する。どちらの z 座標も $p\mathbb{Z} _ p$ の元であることに注意。これは仮に $\mathbb{F} _ p$ で計算していたら $pP$ と $pQ$ どちらも (0, 0) になっていたはずであるため。
3. $pP = (pv_1, w_1), pQ = (pv_2, w_2)$ として、 $(v_2/v_1) \bmod p$ を計算する。これが求めるべき離散対数である。

## 感想
- [Sil2016]:
  - Sections IV.3-6 について、valuation ring の性質や議論に慣れておらず読むのに苦労した。今回の場合では $R = \mathbb{Z} _ p, v = \mathrm{ord} _ p$ なのだから $v(p) = 1$ というのは決まっており Section IV.6 のほとんどの命題はより単純な形で書くことができる。

- 実装上の注意点: zw-座標で計算する場合でも曲線と直線の共有点を求める必要がある。以下が注意点。
  - 同じ点 2 個の加算の場合、接線を求めることになる。接線は w 軸に平行なこともあるので注意。
  - w 軸に平行な直線の場合、w に関する三次方程式を解くことになることに注意。
  - 同じ z に対する w は最大で 3 個存在することに注意。
  - xy-平面で y = 0 である点は zw-座標で無限遠点になってしまうことに注意。 $a_1 = a_3 = 0$ という条件が仮定できればこのような点における接線は y 軸に平行であるため 2 倍したら無限遠点になる。そのため奇素数位数という仮定の下ではこのような点は発生しない。
  - $\mathbb{Z} _ p$ 係数で計算する場合、計算途中の $\lambda$ などは $\mathbb{Q} _ p$ の元になる可能性がある。 $\lambda$ が $\mathbb{Q} _ p \setminus \mathbb{Z} _ p$ の元になってしまう場合の対処は「w 軸に平行な直線の場合」の対処と似ていて、w の方程式とみなすことである。

## 疑問点
### 1
位数が $p$ である必要があるのはなぜ? 位数が $r \neq p$ のときも $E(\mathbb{F} _ p)$ から持ち上げた点 $P$ に対して z が $p | z$ を満たしている点が $rP$ によって得られるのであれば問題ないのでは?

- 厳密に議論を行う。 $Q = mP$ in $E(\mathbb{F} _ p)$ から点を持ち上げると $Q - mP = R$ in $E(\mathbb{Z} _ p)$ ( $R \in E_1(\mathbb{Z} _ p)$ ) が言える。これを r 倍すると $rQ - mrP = rR$ in $E(\mathbb{Z} _ p)$ となるが、 $E_2(\mathbb{Z} _ p)/E_1(\mathbb{Z} _ p) \simeq \mathbb{Z}/p\mathbb{Z}$ が成立するため、 $r \neq p$ の場合は一般に $rR \not \in E_2(\mathbb{Z} _ p)$ であって R を求める手立てもないので、この方法は使えない。
### 2
zw-平面で考える時、位数が $r \neq p$ のときも $E(\mathbb{F} _ p)$ から持ち上げた点 $P$ に対して $rP = O = (0,0)$ となるのはなぜ? 群の演算 $F$ が $F(X, Y) \simeq X + Y$ を満たすなら位数 $p$ であるべきでは?

- $F(X, Y) \simeq X + Y$ が近似として意味を持つ (誤差が $p$ のより高いべきになる) のは $p | X, p | Y$ の時のみ。 $E(\mathbb{F} _ p)$ から持ち上げた点はそれが成り立たないはずなので、 $p$ 倍しても z 座標が $p$ の倍数になるとは限らない。

## References

[Sil2016]: Silverman, Joseph H. The Arithmetic Of Elliptic Curves. 2nd ed., Springer-Verlag, GTM 106, 2016.

[[Sma1999]]: Smart, Nigel P. "The discrete logarithm problem on elliptic curves of trace one." Journal of cryptology 12 (1999): 193-196.

[Sma1999]: https://iacr.org/cryptodb/data/paper.php?pubkey=14320
