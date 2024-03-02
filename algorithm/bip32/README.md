# 乱数源 (シードフレーズなど) から鍵を導出する BIP 32

## 仕様
### 概要
https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki
色々な方法での鍵の導出 (決定的な計算) をできるようにする。
- 乱数の種 -> マスター秘密鍵
  - HMAC-512 を使って生成
- 秘密鍵, index -> 子の秘密鍵 (non-hardened)
  - mpriv, mpub を鍵ペアとした時、mpriv + (mpub, index から生成した決定的乱数) mod n
- 秘密鍵, index -> 子の秘密鍵 (hardened)
  - mpriv, mpub を鍵ペアとした時、mpriv + (mpriv, index から生成した決定的乱数) mod n
- 公開鍵, index -> 子の公開鍵 (non-hardened)
  - `秘密鍵 -> 子の秘密鍵 (non-hardened)` で得られるものと整合的なので、mpub + (mpub, index から生成した決定的乱数)B である。

細かい背景情報は以下。
- 曲線: [secp256k1](https://www.secg.org/sec2-v2.pdf)
  - $p = 2^{256} - 2^{32} - 2^9 - 2^8 - 2^7 - 2^6 - 2^4 - 1$
  - 方程式: $y^2 = x^3+7$
  - cofactor = 1
  - B: ベースポイント
  - n: 位数

鍵から鍵の導出は深さ 255 まで許される。

### セキュリティー
親の鍵ペア (mpriv, mpub) に対して non-hardened な child priv は cpriv=(mpub の関数)+k という形になっていて、 mpub を公開情報と仮定した時 cpriv が漏れたら mpriv も漏れる。

## 既存の実装
### Go
https://github.com/tyler-smith/go-bip32
- 秘密鍵と公開鍵が同じ型
  - Go の標準ライブラリーでは PrivateKey, PublicKey というように分けている
  - [https://soatok.blog/2021/01/20/please-stop-encrypting-with-rsa-directly/](https://soatok.blog/2021/01/20/please-stop-encrypting-with-rsa-directly/)
    > In my opinion, we should stop shipping cryptography interfaces that…
    > - Allow public/private keys to easily get confused (e.g. lack of type safety)

- 暗号技術の操作にふさわしくない big.Int を使っている
  - [https://pkg.go.dev/math/big#Int](https://pkg.go.dev/math/big#Int)
    > Note that methods may leak the Int's value through timing side-channels. Because of this and because of the scope and complexity of the implementation, Int is not well-suited to implement cryptographic operations. The standard library avoids exposing non-trivial Int methods to attacker-controlled inputs and the determination of whether a bug in math/big is considered a security vulnerability might depend on the impact on the standard library.
  - 参考資料: [simple-minds-think-alike.moritamorie.com/entry/crypto-ecdh (Internet Archive)](https://web.archive.org/web/20230601112654if_/https://simple-minds-think-alike.moritamorie.com/entry/crypto-ecdh)

https://github.com/libs4go/crypto/tree/main/elliptic
- 暗号技術の操作にふさわしくない big.Int を使っている

### Rust
https://docs.rs/crate/bip32/0.5.1
- かなりセキュリティーのことを真面目に考えていそう
  - ゼロ化、一定時間演算などを使っている

### Python
https://github.com/scgbckbone/btc-hd-wallet/blob/master/btc_hd_wallet/bip32.py
- int を使っているので多分タイミング攻撃に弱い
- Python で一定時間演算を行うのは難しいらしい
  - [https://securitypitfalls.wordpress.com/2018/08/03/constant-time-compare-in-python/](https://securitypitfalls.wordpress.com/2018/08/03/constant-time-compare-in-python/)
