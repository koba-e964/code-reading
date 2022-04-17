# Rust 1.60.0 に対する所感
[Rust 1.60.0 のリリースノート](https://github.com/rust-lang/rust/blob/master/RELEASES.md#version-1600-2022-04-07)のうち気になったものに対する所感を書く。

## 計算量
BTreeMap, BTreeSet に計算量の保証が行われた。
https://doc.rust-lang.org/1.60.0/std/collections/struct.BTreeMap.html
> Iterators obtained from functions such as BTreeMap::iter, BTreeMap::values, or BTreeMap::keys produce their items in order by key, and take worst-case logarithmic and amortized constant time per item returned.

PR: https://github.com/rust-lang/rust/pull/92706
競技プログラミングでこの関数を使うとき、以前は保証がなくて実装を見なければ安心して使えなかったが、これで安心して使えるようになった。


## 機能
### [Port cargo from toml-rs to toml_edit](https://github.com/rust-lang/cargo/pull/10086)

これがあると `cargo add` のように、Cargo.lock の構造を保存しながらエントリを追加することなどができる。そのうち公式に入る可能性が高そう。

## 使い勝手

### [Make rustc use `RUST_BACKTRACE=full` by default](https://github.com/rust-lang/rust/pull/93566)

個人的にこのような使いやすいデフォルト値は結構大事だと思っているので、真似していきたい。
