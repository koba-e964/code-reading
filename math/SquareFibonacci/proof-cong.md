### (cong-L)
- (cong-L) $2 | k$ のとき、$L_{n+2k} \equiv -L_n \pmod{L_k}$

証明

$L_k | L_{n+2k} + L_n$ が言えれば良い。
$k$ が偶数であり $L_k = \phi^k + \overline\phi^k = \phi^k + \phi^{-k} = \overline\phi^k + \overline\phi^{-k}$ が成立することに注意すると、

$$\begin{align*}
L_{n+2k} + L_n &= \phi^{n+2k} + \overline\phi^{n+2k} + \phi^n + \overline\phi^n \\\\
&= \phi^{n+2k} + \phi^n + \overline\phi^{n+2k} + \overline\phi^n \\\\&= (\phi^k+\phi^{-k})\phi^{n+k} + (\overline\phi^k+\overline\phi^{-k})\overline\phi^{n+k} \\\\
&= (\phi^k+\phi^{-k})(\phi^{n+k} + \overline\phi^{n+k}) \\\\
&= L_kL_{n+k}
\end{align*}$$

### (cong-F)
- (cong-F) $2 | k$ のとき、$F_{n+2k} \equiv -F_n \pmod{L_k}$

証明

上と同様である。ただし $\sqrt{5}$ が分母に来ているのでそこの処理が面倒。

$$\begin{align*}
F_{n+2k} + F_n &= \frac{\phi^{n+2k} - \overline\phi^{n+2k} + \phi^n - \overline\phi^n}{\sqrt{5}} \\\\
&= \frac{\phi^{n+2k} + \phi^n - \overline\phi^{n+2k} - \overline\phi^n}{\sqrt{5}} \\\\
&= \frac{(\phi^k+\phi^{-k})(\phi^{n+k} - \overline\phi^{n+k})}{\sqrt{5}} \\\\
&= L_k F_{n+k}
\end{align*}$$
