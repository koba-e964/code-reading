コード: https://github.com/rust-lang/regex/tree/1.10.3


https://docs.rs/regex-automata/0.4.5/src/regex_automata/nfa/thompson/nfa.rs.html#1463-1469

`>` は start_unanchored (文字列の任意の位置から開始可能) の位置を、`^` は start_anchored (文字列の先頭からのみ開始可能) の位置を指す。
```
PikeVM { config: Config { match_kind: None, pre: None }, nfa: thompson::NFA(
>000000: binary-union(7, 1)
 000001: \x00-\xFF => 0
 000002: a => 8
 000003: b => 4
 000004: c => 5
 000005: binary-union(6, 8)
 000006: d => 5
^000007: binary-union(2, 3)
 000008: MATCH(0)

transition equivalence classes: ByteClasses(0 => [\x00-`], 1 => [a], 2 => [b], 3 => [c], 4 => [d], 5 => [e-\xFF], 6 => [EOI])
)
 }
```
