# 楕円曲線などへのハッシュ関数 (hash_to_curve) について
[BLS12-381](https://hackmd.io/@benjaminion/bls12-381) などの曲線によって実現できるペアリングベースの暗号方式では、hash_to_curve の操作が必要になる。

- 署名: $H(m) \in E_2$ をメッセージ m に対して定まる $G_2$ の点とする。このとき $\sigma = \mathrm{sk}H(m)$ とすれば $e(g_1,\sigma) = e(\mathrm{pk}, H(m))$ が成立する。

すぐに思いつくナイーブな方法は以下である。 (hash and increment と呼ばれることもある)

1. 何らかの暗号学的ハッシュ関数を使い、 $x := H(m) \bmod p$ とする。
2. $x$ から順番に $x+1, x+2,\ldots$ と試していき、 $y^2 = f(x)$ となる $y$ が存在する (言い換えると $f(x)$ が ${\bmod}\ p$ で平方剰余である) ような $x$ が見つかるまで待つ。

この方法はただの確率的多項式時間アルゴリズムであり、多項式時間アルゴリズムではない。このことによって、暗号技術的には以下の問題が発生する。
- ナイーブに実装すると、試した回数が実行時間から推測できてしまう。これにより $m$ などが推測できる可能性がある。
- それを避けるために常に定数回試す (例えば 100 回) という方式を使うと、上の問題は避けられるが今度は実行時間が 100 倍になってしまう。

これらを避けられるように、多項式時間のアルゴリズムが必要である。

[[RFC9360]] で紹介されているアルゴリズムのいくつかを紹介する。

## 処理の流れ
メッセージ -> ハッシュ値 -> F_p の元 -> 曲線上の点 という変換を行う。 F_p の元を得るまでは難しいところはないので、F_p の元 -> 曲線上の点 の部分だけ書く。

### Elligator 2 [[BHKL2013]] について
$f(x) = x^3 + Ax^2 + Bx$ として、モンゴメリー曲線 $y^2 = f(x)$ に対して適用できる。少し変形すれば位数 2 の点を持つすべての楕円曲線に適用できる。
- $f(x)$ が有理根を持てば良い。その有理根を $a$ とすると $(a,0)$ の位数は 2 である。

楕円曲線に対してあらかじめ以下を計算しておく。
- $u$: 小さい平方非剰余。たとえば Curve25519 の場合は $u=2$ である。

$r \in \mathbb{F} _ {p}$ に対して点を割り当てる手順は以下である:
1. $v := -A / (1 + ur^2)$ とする。このように定めると $\chi(v)\chi(-v-A) = -1$ が満たされる。
  - $\chi$ は平方剰余かどうかを返す関数である。 $v$ が平方剰余の時 $\chi(v) = 1$ で、そうでないとき $\chi(v) = -1$ である。ただし $\chi(0) = 0$ とする。
2. $\chi(f(v))\chi(f(-v-A)) = -1$ である。 $f(v)$ と $f(-v-A)$ のうちどちらか一方だけが平方剰余なので、平方剰余である方を使って点を返す。
  - $f(x) = x(x^2 + Ax + B)$ が成立する。 $x^2 + Ax + B$ の部分は $x = v, -v-A$ のどちらでも値が同じであるため、 $\chi(f(v))\chi(f(-v-A)) = \chi(v^2 + Av + B)^2 \chi(v)\chi(-v-A) = \chi(v)\chi(-v-A) = -1$ が成立する。

### SWU について
TODO

# 疑問点
## 1 isogenous <=> 位数が同じ
https://www.johndcook.com/blog/2019/04/21/what-is-an-isogeny/ に言及があった。
> 有限体上の楕円曲線 E1 と E2 について、それらの位数が等しいことと、同種写像 E1 -> E2 が存在することは同値。

# 参考文献
[[BHKL2013]] Bernstein, Daniel J., et al. "Elligator: elliptic-curve points indistinguishable from uniform random strings." Proceedings of the 2013 ACM SIGSAC conference on Computer & communications security. 2013.

[[RFC9360]] Faz-Hernandez, A., Scott, S., Sullivan, N., Wahby, R., and C. Wood, "Hashing to Elliptic Curves", RFC 9380, DOI 10.17487/RFC9380, August 2023, <https://www.rfc-editor.org/info/rfc9380>.

[BHKL2013]: https://dl.acm.org/doi/pdf/10.1145/2508859.2516734

[RFC9360]: https://www.rfc-editor.org/info/rfc9380
