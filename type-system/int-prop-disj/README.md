# 直観主義命題論理の disjunction property について

## 概要
[直観主義命題論理](https://en.wikipedia.org/wiki/Intuitionistic_logic#Syntax)において、以下の性質が成り立つ。
- $A, B$ を論理式とする。このとき $\vdash A \vee B$ であれば、 $\vdash A$ または $\vdash B$ のどちらかは成り立つ。これを [disjunction property](https://en.wikipedia.org/wiki/Disjunction_and_existence_properties) と呼ぶ。

これは古典論理では成り立たない。例えば命題変数 $p$ に対して $\vdash p \vee \neg p$ が成立するが、 $\vdash p$ も $\vdash \neg p$ も成立しない。

このページでは、disjunction property の証明を行う。

## 構文論

### 論理式
以下の規則で作れるものだけが論理式である。論理式全体の集合を $\mathrm{Fml}$ と表記する。

- $\bot$
- $p \in \mathrm{PVar}$ に対し、 $p$
- $A, B \in \mathrm{Fml}$ に対し、
  - $A \wedge B$
  - $A \vee B$
  - $A \to B$

また、以下の略記を導入する。

- $\neg A :\equiv A \to \bot$
- $\top :\equiv \bot \to \bot \equiv \neg \bot$

### 推論規則
略

## 意味論
### Heyting 代数
#### 定義
半順序集合 $(H, \le)$ が Heyting 代数であるとは、以下の演算を持つことをいう。
- $\bot_H$: $\forall x \in H\ldotp \bot_H \le x$
- $\wedge$: $\forall x \in H\ldotp x \le (a \wedge b) \Leftrightarrow (x \le a) \wedge (x \le b)$
- $\vee$: $\forall x \in H\ldotp (a \vee b) \le x \Leftrightarrow (a \le x) \wedge (b \le x)$
- $\to$: $\forall x \in H\ldotp x \le (a \to b) \Leftrightarrow (x \wedge a) \le b$

この条件から、必然的に $\forall x \in H\ldotp x \le \top_H$ を満たす $\top_H$ が存在することに注意。 $\top_H := \bot_H \to \bot_H$ とおけば $x \le \top_H$ は $(x \wedge \bot_H) \le \bot_H$ と同値であり、つまり恒真である。


このページの議論において、Heyting 代数として**有限なものだけを考えて問題ない**。このような性質は**有限モデル性** (finite model property) と呼ばれ、全く自明ではないがここでは証明しない。

#### 準同型
略

### モデル
$e: \mathrm{PVar} \to H$ を**環境** (environment) と呼び、ペア $(H, e)$ を**モデル**と呼ぶ。
モデル $(H, e)$ および論理式 $A \in \mathrm{Fml}$ に対して、モデル $(H, e)$ における $A$ の**表示** (denotation) $\mathrm{eval}_H(e, A) \in H$ を以下のように定義する:

- $\mathrm{eval}_H(e, \bot) := \bot_H$
- $p \in \mathrm{PVar}$ に対して $\mathrm{eval}_H(e, p) := e(p)$
- $A, B \in \mathrm{Fml}$ に対して
  - $\mathrm{eval}_H(e, A \wedge B) := \mathrm{eval}_H(e, A) \wedge \mathrm{eval}(e, B)$
  - $\mathrm{eval}_H(e, A \vee B) := \mathrm{eval}_H(e, A) \vee \mathrm{eval}_H(e, B)$
  - $\mathrm{eval}_H(e, A \to B) := \mathrm{eval}_H(e, A) \to \mathrm{eval}_H(e, B)$

このように定めると、以下の補題が成立する。
#### 補題 1
モデル $(H_1, e)$ および Heyting 代数の準同型 $f: H_1 \to H_2$ があるとする。このとき、 $f(\mathrm{eval} _ {H _ 1}(e, A)) = \mathrm{eval} _ {H _ 2}(f \circ e, A)$ が成り立つ。

証明: 略

論理式 $A$ とモデル $(H, e)$ に対して $\mathrm{eval}_H(e, A) = \top_H$ が成立するとき、 $A$ は $(H, e)$ の上で**妥当** (valid) であるといい、 $(H,e)\Vdash A$ と表記する。また、論理式 $A$ とモデル $(H, e)$ に対して $\mathrm{eval}_H(e, A) \neq \top_H$ が成立する時、 $(H,e)\nVdash A$ と表記し、モデル $(H, e)$ のことを $A$ の **countermodel** と呼ぶ。

### 健全性
任意の論理式 $A$ およびモデル $(H, e)$ に対して、 $\vdash A$ であれば $(H, e) \Vdash A$ が成立する。この性質のことを直観主義命題論理の**健全性** (soundness) と呼ぶ。

直観主義命題論理の健全性から、countermodel を持つ論理式 $A$ について $\nvdash A$ である。これを利用して、 $\nvdash A$ を証明するために countermodel を見つけるという方法が使える。

### 完全性
$A$ を任意の論理式とする。このとき、すべてのモデル $(H, e)$ に対して $(H, e) \Vdash A$ が成立するのであれば、 $\vdash A$ が成立する。このような性質は**完全性** (completeness) と呼ばれ、一般に証明するのが難しい。

対偶をとると、 $\nvdash A$ であれば $A$ は countermodel を持つといえる。

## disjunction property の証明
https://math.stackexchange.com/questions/2000978/proof-of-the-disjunction-property の方針で証明する。

(Kripke model で $W \cup V \cup \lbrace u\rbrace$ を構成するのは、Heyting 代数だと $H_1 \times H_2 + \lbrace \top\rbrace$ を作ることに相当する。)

### 補題 2
$H$ を Heyting 代数とし、 $t_{H}: H + \lbrace \top\rbrace \to H$ を $t_H(h) := h, t_H(\top) := \top_H$ で定める。このとき $t_{H}$ は Heyting 代数の準同型である。

証明: 略


### 定理 3 (disjunction property)
論理式 $A, B$ に対し、 $\nvdash A$ かつ $\nvdash B$ であれば、 $\nvdash A \vee B$ が成立する。

証明:

完全性から $A$ や $B$ は countermodel を持つ。
$A$ の countermodel を $(H_1, e_1)$、 $B$ の countermodel を $(H_2, e_2)$ として、 $A \vee B$ の countermodel を構成しよう。

$H := H_1 \times H_2 + \lbrace \top\rbrace$ 上の環境 $e: \mathrm{PVar} \to H_1 \times H_2 + \lbrace \top\rbrace$ を $e(p) := (e_1(p), e_2(p))$ で定める。このとき、
$$\mathrm{eval} _ {H_1}(e_1, A) \neq \top_1, \mathrm{eval} _ {H_2}(e_2, B) \neq \top_2$$
であり、補題 1, 2 から $t _ {H _ 1\times H _ 2}(\mathrm{eval} _ H(e, X)) = \mathrm{eval} _ {H _ 1 \times H _ 2}(t _ {H _ 1\times H _ 2} \circ e, X)$ であるため、
$$t _ {H _ 1\times H _ 2}(\mathrm{eval} _ H(e, A)) = (\mathrm{eval} _ {H _ 1}(e _ 1, A), \mathrm{eval} _ {H _ 2}(e _ 2, A)) \ne (\top _ 1, \top _ 2)$$
$$t _ {H _ 1\times H _ 2}(\mathrm{eval} _ H(e, B)) = (\mathrm{eval} _ {H _ 1}(e _ 1, B), \mathrm{eval} _ {H _ 2}(e _ 2, B)) \ne (\top _ 1, \top _ 2)$$
から
$$\mathrm{eval}_H(e, A) \le (\top_1, \top_2), \mathrm{eval}_H(e, B) \le (\top_1, \top_2)$$
が言え、 $$\mathrm{eval}_H(e, A \vee B) = \mathrm{eval}_H(e, A) \vee \mathrm{eval}_H(e, B) \le (\top_1, \top_2)$$ が言える。つまり $\mathrm{eval}_H(e, A \vee B) \ne \top$ であり $(H, e)$ は $A \vee B$ の countermodel である。健全性から $\nvdash A \vee B$ である。

### 適用例
$p \in \mathrm{PVar}$ とする。 $p$ の countermodel の例は $(\lbrace0 < 1\rbrace, p \mapsto 0)$ であり、 $\neg p$ の countermodel の例は $(\lbrace0 < 1\rbrace, p \mapsto 1)$ である。
ここから $p \vee \neg p$ の countermodel を構成しよう。

定理 3 の構成だと $(H, e) = (\lbrace (0,0), (0,1), (1,0), (1,1), \top\rbrace, p \mapsto (0, 1))$ となる。
- $\mathrm{eval}_H(e, p) = (0, 1)$
- $\mathrm{eval}_H(e, \neg p) = (0, 1) \to (0, 0) = (1, 0)$
- $\mathrm{eval}_H(e, p \vee \neg p) = (0, 1) \vee (1, 0) = (1, 1) < \top$

よって $(H, e) \nVdash p \vee \neg p$ が成立する。よって $\nvdash p \vee \neg p$ も成立する。

**注意**: 実際にはもっと小さい countermodel $(\lbrace 0 < 1 < 2\rbrace, p \mapsto 1)$ が存在する。

## 参考文献
- ipc_botの解説 決定版! [https://zenn.dev/qnighy/articles/103144250d0a9d](https://zenn.dev/qnighy/articles/103144250d0a9d)
