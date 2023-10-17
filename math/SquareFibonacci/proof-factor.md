# 証明
### (F-even)
- $F_{2m} = F_mL_m$

式変形をすれば明らか。
$$\begin{align*}
F_{2m} &= (\phi^{2m}-\overline\phi^{2m})/\sqrt{5}\\\\
&= (\phi^{m}+\overline\phi^{m})(\phi^{m}-\overline\phi^{m})/\sqrt{5} \\\\
&= L_mF_m
\end{align*}$$

### (gcd-3)
- $\mathrm{gcd}(F_{3m}, L_{3m}) = 2$
### (gcd-not-3)
- $3 \not| n$ のとき、 $\mathrm{gcd}(F_{n}, L_{n}) = 1$

似たような命題なので同時に証明する。

環 $\mathbb{Z}[\phi]$ 上で議論する。有理整数 n に対して
$$\begin{align*}
\mathrm{gcd}(F_{n}, L_{n}) &= \mathrm{gcd}(F_{n}, L_{n} + \sqrt{5}F_n) \\\\
&= \mathrm{gcd}(F_{n}, 2\phi^n) \\\\
&= \mathrm{gcd}(F_{n}, 2)
\end{align*}$$

が成立する。そのため、 $\mathrm{gcd}(F_{n}, L_{n})$ は 1 または 2 である。

$F_{n}$ を mod 2 で見ると $0, 1, 1$ の周期的な列である。これは $3|n$ のとき、およびそのときにのみ 0 である。
