[linked_hash_map:0.5.4](https://github.com/contain-rs/linked-hash-map/tree/0531e100ef052fd49b2f465abf96cd88aea84692/src)

```
// This type exists only to support borrowing `KeyRef`s, which cannot be borrowed to `Q` directly
// due to conflicting implementations of `Borrow`. The layout of `&Qey<Q>` must be identical to
// `&Q` in order to support transmuting in the `Qey::from_ref` method.
```

ちょっと変更して試してみると:

```
error[E0119]: conflicting implementations of trait `std::borrow::Borrow<KeyRef<_>>` for type `KeyRef<_>`
  --> linked_hash_map-mod.rs:96:1
   |
96 | impl<K, Q: ?Sized> Borrow<Q> for KeyRef<K> where K: Borrow<Q> {
   | ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^
   |
   = note: conflicting implementation in crate `core`:
           - impl<T> Borrow<T> for T
             where T: ?Sized;
```
(`run.sh` を実行すれば試せる)


昔からこのデータ構造は難しいと思って身構えていたが、コードを読んでみたらそれほどでもなかった。反省。

## 概要
HashMap + エントリ同士をつなぐリンク。

free: 使用済みの node を覚えておいて後で再利用する。K や V は uninitialized な状態である。特に、K や V を drop してはならない。新しく node として転用したい場合は `ptr::write` を使う。

## 細かい点
`detach`, `attach`, `with_map` は　`unsafe fn` であるべき。内部関数であっても dangling pointer を渡されたら UB を起こすので。
