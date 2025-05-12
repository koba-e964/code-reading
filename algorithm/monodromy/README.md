# Montgomery curve における division polynomial と doubling formula の関係

## 結論
Montgomery 曲線 $y^2 = x^3 + Ax^2 + x$ において、doubling formula は division polyomials を尊重する。


## 注意点

`division_polynomial` で division polynomial を計算する際、m が偶数の場合 `two_torsion_multiplicity` に注意すること。
```python
sage: E.division_polynomial(2, two_torsion_multiplicity=0) # division polynomial / (2y)
1
sage: E.division_polynomial(2, two_torsion_multiplicity=1) # 本来の division polynomial
2*y
sage: E.division_polynomial(2, two_torsion_multiplicity=2) # division polynomial * 2y
4*x^3 + 4*a*x^2 + 4*x
```
これは m が偶数の場合に division polynomial が x の多項式にならず y * (x の多項式) になってしまうことの対策と思われる。今回欲しいのは $\psi_n^2$ でありこれは常に x の多項式であるため問題はないが、 `division_polynomial` の呼び出し方には気をつける必要がある。主な対策法は 2 種類ある。
1. 常に `two_torsion_multiplicity=1` で呼び出して計算し、最終的に $y^2 - f(x)$ で割った余りをとる。
2. `two_torsion_multiplicity=2` と `two_torsion_multiplicity=0` を半々で呼び出し、因子の乗除の影響を取り去りながら x の多項式だけで処理する。

今回は 2. を採用する。
