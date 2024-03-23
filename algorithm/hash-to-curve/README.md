# 楕円曲線などへのハッシュ関数 (hash_to_curve) を一定時間で計算する
[BLS12-381](https://hackmd.io/@benjaminion/bls12-381) などの曲線によって実現できるペアリングベースの暗号方式では、hash_to_curve の操作が必要になる。

- 署名: $H(m) \in E_2$ をメッセージ m に対して定まる $G_2$ の点とする。このとき $\sigma = \mathrm{sk}H(m)$ とすれば $e(g_1,\sigma) = e(\mathrm{pk}, H(m))$ が成立する。

すぐに思いつくナイーブな方法は以下である。 (hash and increment と呼ばれることもある)

1. 何らかの暗号学的ハッシュ関数を使い、 $x := H(m) \bmod p$ とする。
2. $x$ から順番に $x+1, x+2,\ldots$ と試していき、 $y^2 = f(x)$ となる $y$ が存在する (言い換えると $f(x)$ が ${\bmod}\ p$ で平方剰余である) ような $x$ が見つかるまで待つ。

この方法はただの確率的多項式時間アルゴリズムであり、多項式時間アルゴリズムではない。このことによって、暗号技術的には以下の問題が発生する。
- ナイーブに実装すると、試した回数が実行時間から推測できてしまう。これにより $m$ などが推測できる可能性がある。
- それを避けるために常に定数回試す (例えば 100 回) という方式を使うと、上の問題は避けられるが今度は実行時間が 100 倍になってしまう。

これらを避けられるように、多項式時間のアルゴリズムが必要である。そのようなアルゴリズムは constant-time hashing と呼ばれることもある。

[[RFC9360]] で紹介されているアルゴリズムのいくつかを紹介する。

## 注意点
- 用語・記法:
  - $x \in \mathbb{F} _ p$ の **canonical な表現**とは、 $x$ の同値類の中で最小の非負整数である。つまり $0, 1, 2, \ldots, p-1$ のいずれかである。
  - 素数 $p$ に対して $\mathbb{F} _ p$ の元が**正** (**positive**) であるとは、それが $0, 2, 4, \ldots, p-1$ のいずれかであることをいう。簡単にいうと canonical な表現を有理整数と解釈したものが偶数であること。一般の有限体に対する正の定義は疑問点 3 を参照すること。
  - $\mathbb{F} _ q ^ +$ で、 $\mathbb{F} _ q$ の正の元を集めた集合を表す。つまり、素数 $p$ に対して $\mathbb{F} _ p^+ := \lbrace 0, 2, 4, \ldots, p-1\rbrace$ である。

## 処理の流れ
メッセージ -> ハッシュ値 -> F_q の元 -> 曲線上の点 という変換を行う。 F_q の元を得るまでは難しいところはないので、F_q の元 -> 曲線上の点 の部分だけ書く。

### Elligator 2 [[BHKL2013]] について
$f(x) = x^3 + Ax^2 + Bx$ として、モンゴメリー曲線 $y^2 = f(x)$ に対して適用できる。少し変形すれば位数 2 の点を持つすべての楕円曲線に適用できる。
- $f(x)$ が有理根を持てば良い。その有理根を $a$ とすると $(a,0)$ の位数は 2 である。

楕円曲線の定義体 ($\mathbb{F} _ {q}$) に対してあらかじめ以下を計算しておく。
- $u$: 小さい平方非剰余。たとえば Curve25519 の場合は $q = p = 2^{255} - 19$ であるため $u=2$ とできる。

$r \in \mathbb{F} _ {q}$ に対して点を割り当てる手順は以下である:
1. $v := -A / (1 + ur^2)$ とする。このように定めると $\chi(v)\chi(-v-A) = -1$ が満たされる。
    - $\chi$ は平方剰余かどうかを返す関数である。 $v$ が平方剰余の時 $\chi(v) = 1$ で、そうでないとき $\chi(v) = -1$ である。ただし $\chi(0) = 0$ とする。
2. $\chi(f(v))\chi(f(-v-A)) = -1$ である。 $f(v)$ と $f(-v-A)$ のうちどちらか一方だけが平方剰余なので、平方剰余である方を使って点を返す。
    - $f(x) = x(x^2 + Ax + B)$ が成立する。 $x^2 + Ax + B$ の部分は $x = v, -v-A$ のどちらでも値が同じであるため、 $\chi(f(v))\chi(f(-v-A)) = \chi(v^2 + Av + B)^2 \chi(v)\chi(-v-A) = \chi(v)\chi(-v-A) = -1$ が成立する。
    - y 座標の符号は、 $x = v$ のとき正、 $x = -v-A$ のとき負になるようにとる。

### Simplified SWU [[WB2019]] について
$f(x) = x^3 + ax + b$ としたとき、 $ab \ne 0 \in \mathbb{F} _ {q}$ であれば適用できる。

楕円曲線の定義体 ($\mathbb{F} _ {q}$) に対してあらかじめ以下を計算しておく。
- $\xi$: 小さい平方非剰余。例えば secp256k1 や BLS12-381 の場合は以下である。
  - secp256k1: $q = p = 2^{256} - 2^{32} - 977 \equiv 3 \pmod 4$ なので $\xi := -1$ ととれる。
  - BLS12-381: $G_2$ の方でやるとする。 $q = p^2$ であり $p = (x^4 - x^2 + 1)(x - 1)^2/3 + x, 2^{16} | x$ であるため $p \equiv 3 \pmod{8}$ であり $\xi := 1 + \sqrt{-1}$ ととれる。
    - $\xi^{(p^2 - 1)/2} = (2\sqrt{-1})^{(p^2 - 1)/4} = -2^{(p^2 - 1)/4} = -(2^{p-1})^{(p+1)/4} = -1$

$t \in \mathbb{F} _ {q}$ に対して点を割り当てる手順は以下である:
1. $u := \xi t^2$ とおく。 $t \neq 0$ ならこの $u$ は平方非剰余である。
2. $f(ux) = u^3f(x)$ は 1 次方程式であり、 $x$ の項は 0 でない。これを $x$ について解くと $x = -(b/a)(1 + 1/(u^2+u))$ であり、これを $X_0(t)$ と呼ぶ。
3. $X_1(t) := uX_0(t)$ とする。 $\chi(f(X_0(t)))\chi(f(X_1(t))) = -1$ であるため、 $f(X_0(t))$ と $f(X_1(t))$ のどちらかは平方剰余なのでそちらを採用する。
    - $\chi(f(X_0(t)))\chi(f(X_1(t))) = \chi(f(X_0(t)))^2\chi(u^3) = \chi(u) = -1$
    - y 座標の符号は、 $x = X_0(t)$ のとき正、 $x = X_1(t)$ のとき負になるようにとる。

### Simplified SWU with ab = 0 [[WB2019]] について
上の Simplified SWU は、そのままだと ab = 0 である曲線、例えば $y^2 = x^3 + b$ (secp256k1 や BLS12-381) などに適用できない。
そこで、isogenous curve とそこからの同種写像を作って、そこから点を移すことで実現する。([[WB2019]] の Section 4)

[[WB2019]] の Subsection 4.3 に詳しいパラメーターが記載されている。たとえば BLS12-381 の $E_2: y^2 = x^3 + 4(1+\sqrt{-1})$ には $y^2 = x^3 + 240\sqrt{-1}x + 1012(1+\sqrt{-1})$ がある。

# 疑問点
## 1: isogenous <=> 位数が同じ
https://www.johndcook.com/blog/2019/04/21/what-is-an-isogeny/ に言及があった。
> 有限体上の楕円曲線 E1 と E2 について、それらの位数が等しいことと、同種写像 E1 -> E2 が存在することは同値。

## 2: isogenous な曲線をどうやって見つける?
おそらく https://web.archive.org/web/20220211033616/https://aghitza.github.io/publication/translation_velu/ などを読めば良さそう。([[WB2019]] の参考文献 [Vél71])

## 3: 素体でない有限体の上の符号はどう定義する?
例えば [[WB2019]] の Section 2 では、 $\gamma_0 + \gamma_1 \sqrt{d}$ に対して $\gamma_1$ の符号を先に見て、 $\gamma_1 = 0$ のときは $\gamma_0$ の符号を見るという方式でやっている。

## 4: これらの写像は単射なのか、全射なのか?
どちらでもない。例えば Elligator 2 の場合、同じ $ur^2$ に対してそれを実現する $r$ は 2 通りあることが普通。そう考えると像の大きさは $(p+1)/2$ である。

無限遠点に行くこともないので $(p+1)/2 \le |E| - 1$ である。Hasse の定理および位数 2 の点が存在することから $|E| \ge \max(p - 2 \sqrt{p}, 1) + 1$ であるため、 $(p+1)/2 \le \max(p - 2 \sqrt{p}, 1)$ が成立してほしい。
これは $p \ge 18$ であれば常に成立する。

# 参考文献
[[BHKL2013]] Bernstein, Daniel J., et al. "Elligator: elliptic-curve points indistinguishable from uniform random strings." Proceedings of the 2013 ACM SIGSAC conference on Computer & communications security. 2013.

[[RFC9360]] Faz-Hernandez, A., Scott, S., Sullivan, N., Wahby, R., and C. Wood, "Hashing to Elliptic Curves", RFC 9380, DOI 10.17487/RFC9380, August 2023, <https://www.rfc-editor.org/info/rfc9380>.

[[WB2019]] Wahby, Riad S., and Dan Boneh. "Fast and simple constant-time hashing to the BLS12-381 elliptic curve." Cryptology ePrint Archive (2019).

[BHKL2013]: https://dl.acm.org/doi/pdf/10.1145/2508859.2516734

[RFC9360]: https://www.rfc-editor.org/info/rfc9380

[WB2019]: https://eprint.iacr.org/2019/403
