# 素数性の証明アルゴリズム ECPP
## 楕円曲線の位数について

[Sil2016] Thm. V.2.3.1

$E = E(\mathbb{F} _ q)$ を楕円曲線とする。
$\phi_p : (x, y) \mapsto (x^q, y^q)$ をフロベニウス自己準同型とする。また $a = q + 1 - |E(\mathbb{F} _ q)|$ とする。
- $\phi_p^2 - a\phi_p + q = 0$ in $\mathrm{End}(E)$

重要事実:

- $\det \phi_l = \deg \phi$

## アルゴリズムの概要

与えられた正の有理整数 N (おそらく素数) の素数性を証明したい。

以下の定理が重要。

### [[AO1993]] Thm. 5.2 + Cor. 5.1: 
N を 6 と互いに素な正の有理整数、E を $\mathbb{Z}/N\mathbb{Z}$ 上の楕円曲線、 $P \in E$ をその上の点、m と s を $s | m$ を満たす正の有理整数とする。q | s なる素因数 q それぞれに対し $(m / q)P = (x_q : y_q : z_q)$ となるように $x_q, y_q, z_q$ を定める。さらに $mP = O_E$ かつすべての q に対し $\gcd(z_q, N) = 1$ を仮定する。
- このとき p が N の素因数であれば $|E(\mathbb{Z}/p\mathbb{Z})| \equiv 0 \pmod{s}$ である。
  - とくに $s \gt (N^{1/4}+1)^2$ であればそのような $p < N$ が存在せず N は素数である。


### アルゴリズム ECPP(N)
再帰アルゴリズムである。(元論文の書き方はわかりにくいので、再帰アルゴリズムとしてわかりやすく記述する。)
1. $N \le N _ {\mathrm{small}}$ であれば素数性が証明できたとして終了。
2. $N = \pi \overline \pi$ なる整数 $\pi$ が存在するような虚二次体 $\mathbb{Q}(\sqrt{-D})$ を見つける。
3. $\mathbb{Q}(\sqrt{-D})$ の単数は $w = w(-D)$ 個あるが、それらを $\epsilon _ 1, \ldots, \epsilon _ w$ とする。 $m_i = N(\epsilon _ i \pi - 1)$ としうまく素因数分解できる $m_i$ を見つける。見つからなければ 2. に戻る。ここでの条件は $m_i = FN'$ であって、 $F$ は有理整数、N' はおそらく素数であり $N' \gt (N^{1/4} + 1)^2$ をみたす有理整数であるようなものである。
4. $H_D(X)$ を Hilbert class polynomial とする。 $H_D(X) = 0 \pmod{N}$ の根を 1 つ選び j と置く。このような根は N が素数であれば必ず $\mathbb{Z}/N\mathbb{Z}$ に存在するので、存在しなければ、あるいは計算途中でエラーが発生すれば、N は合成数である。
5. j-invariant が j であり、位数が $m_i$ であるような楕円曲線 $E(\mathbb{Z}/N\mathbb{Z})$ を構成し、その上の点 P を選ぶ。
6. $s = N', m = m_i$ として上の定理の前提条件を確かめる。つまり、 $(m/s)P = FP \neq O_E$ かつ $mP = m_iP = O_E$ を確かめる。失敗したら 2. に戻る。
7. ECPP(N') を実行して N' が素数であることを証明する。

## 楕円曲線の構成
この節では ECPP で重要な楕円曲線の構成を説明する。後で実際にここで構成した楕円曲線の上で計算する必要があるため、存在が証明できるだけでは足りず実際に構成する必要がある。


数式処理のために Sage を使う。バージョンは 9.8 である。
```sage
sage: version()
'SageMath version 9.8, Release Date: 2023-02-11'
```

### 定理 ([[AO1993]] 4.2)
H_D(X) mod p の根を任意にとって j と置く。j-invariant が j であり位数が p + 1 - tr(π) となるような $\mathbb{F} _ p$ 上の楕円曲線が存在する。

### 例1
p = 17, Δ = -8 とすると、 $4 \times 17 = 6^2 + 8 \times 2^2$ が成立する。よって p は $\mathbb{Q}(\sqrt{-2})$ の整数 $\pi = 3 + 2\sqrt{-2}$ のノルムである。

$J(\sqrt{-2}) = (5/3)^3 \equiv 4 \pmod{17}$ である。
```sage
sage: hilbert_class_polynomial(-8).change_ring(GF(17)).factor()
x + 7
sage: GF(17)(-7/1728)
4
```
そのため $k = J/(J-1) \equiv 7 \pmod{17}$ が成立する。c = 3 とすると c は mod 17 で平方非剰余であるため、 $y^2 = x^3 - 3kc^{2r} + 2kc^{3r}$ は以下のようになる:
1. $y^2 = x^3 + 13x + 14 \quad (r = 0)$
2. $y^2 = x^3 + 15x + 4 \quad (r = 1)$

前者の位数が 12 で後者の位数が 24 である。これは $p+1 \pm \mathrm{tr}(\pi) = 18 \pm 6 = 24,12$ と整合的である。

```sage
sage: E1 = EllipticCurve(GF(17), [0, 0, 0, 13, 14])
sage: len(E1.rational_points())
12
sage: E2 = EllipticCurve(GF(17), [0, 0, 0, 15, 4])
sage: len(E2.rational_points())
24
```

### 例2

$j(\sqrt{-5}) = (50 + 26\sqrt{5})^3$ らしい。
```sage
sage: (elliptic_j(sqrt(-5.0)))^(1/3)-sqrt(5.0)*26
50.0000000000000
```

$4 \times 29 = 6^2 + 20 \times 2^2$ であるため、29 は $\mathbb{Q}(\sqrt{-5})$ の整数 $\pi = 3 + 2 \sqrt{-5}$ のノルムである。

$\sqrt{5} \equiv 11, 16 \pmod{29}$ であり、それぞれに対して $j$ を計算すると $j \equiv (50 \pm 26 \times 11)^3 \equiv 12, 23$ である。これらに対する $k = j / (j - 1728)$ は $k \equiv 15, 28$ である。

$j \equiv 12, 23$ は以下によっても確かめられる:
```sage
sage: hilbert_class_polynomial(-20).change_ring(GF(29)).factor()
(x + 6) * (x + 17)
```

k の値が 15, 28 のどちらであるにしても $y^2 = x^3 - 3kx + 2k$ の位数は 36 であった。

```sage
sage: E1 = EllipticCurve(GF(29), [0, 0, 0, -45, 30])
sage: len(E1.rational_points())
36
sage: E2 = EllipticCurve(GF(29), [0, 0, 0, -84, 56])
sage: len(E2.rational_points())
36
```

## 確率の分析
h を類数、g を genera とすると g <= h かつ g | h が成り立つ。

## 疑問点
### 1
有理素数 p が K = Q(sqrt(-d)) で完全分解するとき、 K の類体を L としたら p は L で完全分解するか?

-> https://mathoverflow.net/questions/207922/quickest-and-or-most-elementary-proof-of-principal-iff-splits-completely 類体論から、K において p が主イデアル 2 個に分解するのであれば L においてそれぞれの因子が完全分解するはず。

### 2
そもそも isogeny の degree って何だったっけ?
-> 
まず [Sil2016] の I.1 での定義を復習する。
- アファイン代数的集合 V がアファイン多様体であるとは、I(V) が $\overline{K}[X]$ で素イデアルであることをいう。(K[X] で素イデアルであることでは足りない。)
- V/K が多様体 (V が K 上で定義された多様体) であるとする。V/K の**アファイン座標環**を $K[V] := K[X]/I(V/K)$ と定義する。この K[V] は整域であるが、K[V] の商体を K(V) とし、V/K の**関数体**と呼ぶ。
  - $\overline{K}[V]$ の元は関数 $V \to \overline{K}$ を惹き起こすことに注意。V の上では I(V) の差異は消えるため。


[Sil2016] の III.4 には以下のように書かれている:
- isogeny とは射 $\phi: E_1 \to E_2$ であって $\phi(O) = O$ であるものである。
- $\phi$ を isogeny とする。 $\phi \neq 0$ であれば $\phi$ は全射であるため(TODO: なんで?)自然に関数体の間の単射 $\phi^\ast: \overline{K}(E_2) \to \overline{K}(E_1)$ が得られる。 $\overline{K}(E_1)/\phi^\ast\overline{K}(E_2)$ は体の有限次拡大だが、それの次数を $\deg \phi := [\overline{K}(E_1) : \phi^\ast\overline{K}(E_2)]$ とする。また deg 0 = 0 とする。(TODO: II.2 の定義を使う)
  - 例えば $f: \mathbb{C}(\mathbb{P}^1) \to \mathbb{C}(\mathbb{P}^1)$ を $f(X, Y) := (X^2 - 2Y^2, Y^2)$ で定めると、 $\deg f = 2$ である。
- 任意の $Q \in E_2$ に対し $\deg_s \phi = \phi^{-1}(Q)$ (Thm. 4.10 (a))


### 3
H_D(X) = 0 mod N の根を a とおくと、a が j であるような楕円曲線の位数が N+1-tr(π) である (ことがある?) のはなぜ?

### 4
Class polynomial H_D(X) をどうやって計算するのか?

### 5

フロベニウス写像 $\phi_p$ について、勘違いしていたこと。

- $\phi_p$ は $\overline{\mathbb{F} _ p}$ では有限位数ではない。
- $\phi_p$ は $\overline{\mathbb{F} _ p}$ で可逆だが、だからといってノルムが $\pm 1$ とは限らない。

### 6
j-invariant に関する事実
- ω が虚二次体の数であるとき j(ω) は代数的整数である
- ω が虚二次体の数であるとき Q(j(ω),ω)/Q(ω) はヒルベルト拡大である

## References

[[AO1993]]: Atkin, A. Oliver L., and François Morain. "Elliptic curves and primality proving." Mathematics of computation 61.203 (1993): 29-68.

[Sil2016]: Silverman, Joseph H. The Arithmetic Of Elliptic Curves. 2nd ed., Springer-Verlag, GTM 106, 2016.

[[Tsu2017]]: ガウスの種の理論 (Genus Theory), https://tsujimotter.hatenablog.com/entry/genus-theory

[AO1993]: https://www.ams.org/journals/mcom/1993-61-203/S0025-5718-1993-1199989-X/S0025-5718-1993-1199989-X.pdf

[Tsu2017]: https://tsujimotter.hatenablog.com/entry/genus-theory
