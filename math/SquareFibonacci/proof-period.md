### (per-L2) $3 | m \Leftrightarrow 2 | L_m$
$L_{n}$ を mod 2 で見ると $0, 1, 1$ の周期的な列である。これは $3|n$ のとき、およびそのときにのみ 0 である。

### (per-L3) $m \equiv 2 \pmod{4} \Leftrightarrow 3 | L_m$
$L_{n}$ を mod 3 で見ると $2, 1, 0, 1$ の周期的な列である。

### (per-L4) $2 | k, 3 \not | k$ のとき、 $L_{k} \equiv 3 \pmod{4}$

$L_{n}$ を mod 4 で見ると $2, 1, 3, 0, 3, 3$ の周期的な列である。

### (per-L8) $L_{n+12} \equiv L_n \pmod{8}$
$\phi^6 = 5 + 8\phi$ であることに注意。イデアル $8\mathbb{Z}[\phi]$ を法として考えると、 $\phi^6 \equiv 5 \pmod{8\mathbb{Z}[\phi]}$ なので $\phi^{12} \equiv 5^2 \equiv 1 \pmod{8\mathbb{Z}[\phi]}$ である。

これに注目すると、

$$\begin{align*}
L_{n+12} &= \phi^{n+12} + \overline\phi^{n+12} \\\\
&\equiv \phi^n + \overline\phi^n \\\\
&\equiv L_n \pmod{8\mathbb{Z}[\phi]}
\end{align*}$$

有理整数 a, b に対して $a \equiv b \pmod{8\mathbb{Z}[\phi]}$ (つまり $a - b \in 8\mathbb{Z}[\phi]$) のとき $a \equiv b \pmod 8$ であることは、 $8\mathbb{Z}[\phi] \cap \mathbb{Z} = 8\mathbb{Z}$ から明らか。
