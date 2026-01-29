source: [eff v8.0.0](https://github.com/atnos-org/eff/tree/v8.0.0/core/shared/src/main/scala/org/atnos/eff)

# 概要
2.1 の `It i a` は Reader monad とは全然違う。 `It i a = a + (i -> It i a)` である。ADT であり、計算の履歴のようなものが残っている。

さらに、普通の Reader monad のような解釈 (`runReader`) とは別の解釈 (`feedAll`) もできる。

> Unlike the MTL Reader, `It i a` may be treated differently: each
> request gets a new value, as if read from an input stream

`It i` は自由モナドとして構築できる。 `data Reader i x = Get (i -> x)` として、 `It i = Free (Reader i)` とすればよい。


これのいいところは合成が容易いところ。 `Free F` と `Free G` を合成したい場合、 `Free (F+G)` でよい。

# 参考文献

- <https://github.com/atnos-org/eff/tree/v8.0.0/core/shared/src/main/scala/org/atnos/eff>

- <https://okmij.org/ftp/Haskell/extensible/more.pdf>
