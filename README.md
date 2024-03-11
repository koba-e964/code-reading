# code-reading

[koba](https://github.com/koba-e964) がコードリーディング、勉強などをして分かった結果をまとめる。

# 目次
- アルゴリズム
  - 暗号理論関係
    - [Schoof のアルゴリズム](algorithm/schoof)
    - [ペアリング曲線の構成について](algorithm/optimal-ate-pairing/)
    - [OpenSSH における ed25519 の実装 (@V_8_8_P1)](algorithm/OpenSSH:V_8_8_P1-ed25519/)
    - [Smart-attack](algorithm/smart-attack/)
    - [乱数源 (シードフレーズなど) から鍵を導出する BIP 32](algorithm/bip32/)
    - [素数位数の群を高速に実現する Ristretto (`ristretto255`) について](algorithm/ristretto255/)
  - [Reed-Solomon 符号の実装](algorithm/reed-solomon/)
  - [素数性の証明アルゴリズム ECPP](algorithm/ecpp/)
  - [gzip, deflate フォーマットの中身](algorithm/deflate/)
- データ構造
  - [Linked Hash Map (linked_hash_map:0.5.4)](data-structure/linked-hash-map:0.5.4/)
- 型理論
  - [Java Generics are Turing Complete (2016)](type-system/java-generics-are-turing-complete/)
  - [generic-array:0.14.4 (Rust)](type-system/generic-array:0.14.4/)
- 数学
  - [フィボナッチ数の平方数は 0, 1, 144 のみ](math/SquareFibonacci/)

# このリポジトリーのポリシー
- 外部から頂いた PR はマージしない。
- 誤字脱字などを見つけていただいた際に PR などを頂くのは歓迎するが、その変更を自分でコミットしてその PR はマージせずに閉じる。
- 何かを書いてほしいという要望は受け付けない。
- 誤字脱字の指摘は PR か各種 SNS のアカウントまで。
