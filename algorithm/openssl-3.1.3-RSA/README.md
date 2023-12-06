[OpenSSL における RSA の鍵生成の実装 (openssl-3.1.3)](https://github.com/openssl/openssl/blob/openssl-3.1.3/crypto/rsa/rsa_gen.c)
## 概要
RSA の鍵生成 (素数 2 個の生成) アルゴリズムの紹介。
「ランダムに生成した整数が素数であれば OK」より複雑な生成プロセスを行い、[p+1 法や p-1 法](https://wacchoz.hatenablog.com/entry/2019/01/20/120000)を防いでいる。

### 鍵生成 (ossl_rsa_sp800_56b_generate_key)

ソースコード: https://github.com/openssl/openssl/blob/openssl-3.1.3/crypto/rsa/rsa_sp800_56b_gen.c#L356-L420

[[FIPS186-4]] B.3.6 を使用し、確率的素数を 2 個生成する。

### [[FIPS186-4]] B.3.6
ソースコード: https://github.com/openssl/openssl/blob/openssl-3.1.3/crypto/bn/bn_rsa_fips186_4.c#L184-L252

入力:
- nlen: $n$ に要求されるビットサイズ
- $e$: 公開鍵の指数

出力:
- $p, q$: 2 つの確率的素数。

内部で [[FIPS186-4]] C.9 を呼ぶ。

### [[FIPS186-4]] C.9
ソースコード: https://github.com/openssl/openssl/blob/openssl-3.1.3/crypto/bn/bn_rsa_fips186_4.c#L275-L406


入力:
- $nlen$: $n$ に要求されるサイズ
- $r_1, r_2$: 確率的素数

出力:
- $p$: 確率的素数。
  - $p \equiv 1 \pmod{r_1}, p \equiv -1 \pmod{r_2}$ を満たす。おそらく p+1 法や p-1 法などによる素因数分解を防ぐことを目的として、 $p \pm 1$ が大きい素因数を持つようにするため。
  - $\mathrm{gcd}(p-1, e) = 1$
  - $p < 2^{nlen/2}$
- $X$: $p$ の生成時に使った整数の乱数。
  - $2^{nlen/2-1/2} < X < 2^{nlen/2}$
  - $X \le p < X + 2r_1r_2$


[[FIPS186-4]]: National Institute of Standards and Technology (2013) Security Requirements for Cryptographic Modules. (Department of Commerce, Washington, D.C.), Federal Information Processing Standards Publications (FIPS PUBS) 186-4.

[FIPS186-4]: https://www.omgwiki.org/dido/doku.php?id=dido:public:ra:xapend:xapend.b_stds:tech:nist:dss
