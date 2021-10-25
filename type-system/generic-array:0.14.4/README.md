[generic-array:0.14.4](https://github.com/fizyk20/generic-array/tree/0.14.4/src)

## 概要
`[T; N]` のようなもの。Rust では型引数として整数を取ることができなかった (TODO ... まで) ので、それに対するワークアラウンドとして実装された。

`typenum` で定義されたコンパイル時整数を表す型 (`U0`, `U1`, ... などの `Unsigned` を実装した型) に対して、`ArrayLength<T>` というトレイトを実装している。

`typenum` が提供する符号なし整数は二進数の形で提供されているので、それと同じ構造で配列が定義されている。具体的には、符号なし整数は以下のコンストラクタで定義されている:
```
UTerm : 0
UInt<U, B0> : 2 * U
UInt<U, B1> : 2 * U + 1
```

これの各コンストラクタに対し、`<U as ArrayLength<T>>::ArrayType` は以下のように定義される:
```
U = UTerm : [T; 0]
U = UInt<N, B0> : (N::ArrayType, N::ArrayType)
U = UInt<N, B1> : (N::ArrayType, N::ArrayType, T)
```

(なお、タプルで表記したのはわかりやすさのためで、実際には専用の構造体 `GenericArrayImplEven`, `GenericArrayImplOdd` が定義されている。)


### 疑問点
アラインメントは配列と同じなのか? ... 同じ
(TODO なぜなのか書く)
