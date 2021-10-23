source: [typenum:1.14.0](https://github.com/paholg/typenum/tree/v1.14.0/src)

## 概要
trait system を駆使して、コンパイル時整数を実現している。

整数は型として表現される。

型は基本的に singleton である。
基本的に、以下の形式でコンパイル時に型を得られるようになっている。

```
// A XXX B = C
impl XXX<B> for A {
  Output = C
  fn xxx() -> Self::Output {
    /* 唯一の値が返る */
  }
}

// operator_aliases.rs
// 便利 type。XXX に対しての YYY はアドホックに命名されている。(BitAnd -> And, Add -> Sum, Gcd -> Gcf など。)
pub type YYY<A, B> = <A as XXX<B>>::Output;
```

型システムを使う都合上、第一引数での場合分けが必須。


## 個別モジュール
### [bit.rs](https://github.com/paholg/typenum/blob/v1.14.0/src/bit.rs)
TODO

### [uint.rs](https://github.com/paholg/typenum/blob/v1.14.0/src/uint.rs)
符号なし整数 (unsigned integer) を扱う。

#### データ表現形式
二進数。
Uterm: Unsigned ... 0
UInt<U: Unsigned, B: Bit> ... 2 * U + B (B = B0 or B1)

符号なし整数を表す型は `Unsigned` を実装する。

UInt の定義段階で `U: Unsigned, B: Bit` という制限を課すべきであるが、`but enforcing these bounds causes linear instead of logarithmic scaling` ということらしい。同様に、leading zero がないという制約もあるが、それは型では表現されていない。
TODO: linear scaling の理由を調べる。

#### 演算
|演算| 性質|
|:---:|:---:|
|Add, Mul, Or, Shl, Shr, GetBit | 普通|
|Sub, And, Xor, SetBit | leading zero の関係で trim が必要|
|Cmp, Pow, Min, Max | 下の桁から再帰する場合、内部状態を持つ必要あり|
|BitDiff | ?|
|Div, Rem, PartialDiv, Gcd, Sqrt, Logarithm2 | 高度|

### [marker_traits.rs](https://github.com/paholg/typenum/blob/v1.14.0/src/marker_traits.rs)
crate で定義された型しか実装できないような仕掛け (`Sealed`) を利用している。

### [operator_aliases.rs](https://github.com/paholg/typenum/blob/v1.14.0/src/operator_aliases.rs)
概要で述べた alias が定義されている。面倒だからね。

### [private.rs](https://github.com/paholg/typenum/blob/v1.14.0/src/private.rs)
> `**Ignore me!**`
>
> `Just look away.`
>
> `Loooooooooooooooooooooooooooooooooook awaaaaaaaaaaaayyyyyyyyyyyyyyyyyyyyyyyyyyyyy...`

ok you got it

PrivateXXX 系の trait の定義が行われている。実際の実装は `uint.rs` でやられている。

`InternalMarker`, `Internal`: TODO

`Trim`, `TrimTrailingZeros`, `Invert`: TODO

### [int.rs](https://github.com/paholg/typenum/blob/v1.14.0/src/int.rs)
TODO

### [type_operators.rs](https://github.com/paholg/typenum/blob/v1.14.0/src/type_operators.rs)
TODO

### [array.rs](https://github.com/paholg/typenum/blob/v1.14.0/src/array.rs)
コンパイル時整数からなるコンパイル時配列 (を表現する型) を提供する。

TODO


### [lib.rs](https://github.com/paholg/typenum/blob/v1.14.0/src/lib.rs)
re-exporting, `Ordering` の定義、`sealed` モジュールの定義。