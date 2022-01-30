[OpenSSH における ed25519 の実装 (V_8_8_P1)](https://github.com/openssh/openssh-portable/blob/V_8_8_P1/ed25519.c)
## 概要
https://ja.wikipedia.org/wiki/%E3%82%A8%E3%83%89%E3%83%AF%E3%83%BC%E3%82%BA%E6%9B%B2%E7%B7%9A%E3%83%87%E3%82%B8%E3%82%BF%E3%83%AB%E7%BD%B2%E5%90%8D%E3%82%A2%E3%83%AB%E3%82%B4%E3%83%AA%E3%82%BA%E3%83%A0 で紹介されているアルゴリズムの実装。

> Ed25519では、秘密のデータに依存した分岐命令と配列参照が用いられておらず、多くのサイドチャネル攻撃に耐性がある。

という特徴を持つ。


### [sc25519.c](https://github.com/openssh/openssh-portable/blob/V_8_8_P1/sc25519.c)
mod l = 2^252 +  27742317777372353535851937790883648493 (= 2^252 + 0x14def9dea2f79cd65812631a5cf5d3ed) での値。
ほとんどの関数で値に依存した分岐と配列参照を完全に避けており、そうでない関数には名前に `_vartime` がついている。

[m](https://github.com/openssh/openssh-portable/blob/V_8_8_P1/sc25519.c#L15-L16) に mod の値が書かれている。

#### 疑問点
mu とは?

### [fe25519.c](https://github.com/openssh/openssh-portable/blob/V_8_8_P1/fe25519.c)
mod 2^255 - 19 での値。
ほとんどの関数で値に依存した分岐と配列参照を完全に避けており、そうでない関数には名前に `_vartime` がついている。

注意点:
- 常に mod 2^255 - 19 で reduce しているわけではない。2^255 未満であればよいところではそうしている。
  - pack, unpack 時は流石に reduce する

### [ge25519.c](https://github.com/openssh/openssh-portable/blob/V_8_8_P1/ge25519.c)
mod 2^255 - 19 での twisted Edwards curve の上の点の計算。
ほとんどの関数で値に依存した分岐と配列参照を完全に避けており、そうでない関数には名前に `_vartime` がついている。
使われているテクニック:
- スカラーの掛け算: 愚直に実装すると偶数なら 2 倍だけ、奇数なら 2 倍 + 足し算という分岐が発生する。これを避けるため、基点 B に対して 1B, 8B, ..., 8^{84}B = 2^{248}B およびこれらそれぞれの点 x に対して x, 2x, 3x, 4x を事前に計算しておき、スカラーを各桁が負でもよいような八進数に直し、桁ごとに決定的に値を求めてまとめて足す。(実装: [L306-L321](https://github.com/openssh/openssh-portable/blob/V_8_8_P1/ge25519.c#L306-L321))
  - より詳しく: 255 bit 符号なし整数 s があったとき sB を計算したいとして、s の八進数表記を (s_0 s_1 ... s_84) と置く。(-4 <= s_i < 4)
  x = 8^{i}B に対して、0, x, 2x, 3x, 4x が書かれたテーブルを用意しておく。
    1. 一時変数 t を用意する。
    1. j = 0, ..., 4 に対して、abs(s_i) = j であれば t を jx で置き換える。これは分岐なしでできる。
    1. 一時変数 v を用意し -t を代入する。
    1. s_i < 0 であれば t を v で置き換える。これは分岐なしでできる。
   
    ここで得られた t をすべて足せば良い。


#### 疑問点
- (x, y, z, t) の z や t
- sc25519_window3 に渡せる s の範囲は? 255 bit しか見ていないのにどうやって 8l > 2^255 通りの値に対応しているのか? => sc25519 そのものは mod l なので 2^252 通り程度の値しかなく、問題ない。

### [ed25519.c](https://github.com/openssh/openssh-portable/blob/V_8_8_P1/ed25519.c)

本丸。

- [`crypto_sign_ed25519_keypair`](https://github.com/openssh/openssh-portable/blob/V_8_8_P1/ed25519.c#L26-L49): 鍵ペアを作る関数
- [`crypto_sign_ed25519`](https://github.com/openssh/openssh-portable/blob/V_8_8_P1/ed25519.c#L51-L101): 秘密鍵を使って署名を作る関数
- [`crypto_sign_ed25519_open`](https://github.com/openssh/openssh-portable/blob/V_8_8_P1/ed25519.c#L103-L144): 公開鍵を使って署名を検証する関数。`_vartime` 関数呼び放題。

#### 疑問点
- [L67-69](https://github.com/openssh/openssh-portable/blob/V_8_8_P1/ed25519.c#L67-L69) で s (256-bit) の下位 3 bit を 000 にし上位 2 bit を 01 にしている。下位 3 bit の方は、楕円曲線の cofactor が 8 = 2^3 なので基点により生成される部分群の上の点にするために必要と思われるが(TODO 裏をとる、基点の位数が l でなく 8l の可能性もあり)、上位 2 bit の方は意図がわからない。
  - 仕様で決められているのはわかった: [Bern2011#ed25519]
   [Bern2006#curve25519] によれば、s を 8 の倍数にすると small-subgroup attacks に対する耐性ができる。また上位 2 bit を 01 にするのもこれの時点で既に行われている。
   また、最上位ビットの位置を固定にすると doubling の回数が固定になり、タイミング攻撃に対する耐性ができる。ただしこれは doubling でスカラー倍を実装していた時期の名残であり、今の実装では必要ないはず。

- 署名の検証時に 2^c(...) = 0 という形の等式をチェックするという仕様になっている。しかし基点の位数は 2^cl ではなく l だったので、より厳しい (...) = 0 という等式が成立するはずで、仕様 [Bern2011#ed25519] はこちらをチェックすることも許容している。しかし 2^c(...) = 0 の方が弱い上により多くの計算を要するのに、誰がチェックしたいと思うだろうか? 実際 OpenSSH は (...) = 0 をチェックしている。

### References
[Bern2006#curve25519]: Daniel J. Bernstein. "Curve25519: new Diffie-Hellman speed records." Pages 207–228 in Public key cryptography—PKC 2006, 9th international conference on theory and practice in public-key cryptography, New York, NY, USA, April 24–26, 2006, proceedings, edited by Moti Yung, Yevgeniy Dodis, Aggelos Kiayias, Tal Malkin, Lecture Notes in Computer Science 3958, Springer, 2006, ISBN 3-540-33851-9.

[Bern2011#ed25519]: Daniel J. Bernstein, Niels Duif, Tanja Lange, Peter Schwabe, Bo-Yin Yang. "High-speed high-security signatures." Pages 124–142 in Cryptographic hardware and embedded systems—CHES 2011, 13th international workshop, Nara, Japan, September 28–October 1, 2011, proceedings, edited by Bart Preneel, Tsuyoshi Takagi, Lecture Notes in Computer Science 6917, Springer, 2011, ISBN 978-3-642-23950-2. Journal version: Journal of Cryptographic Engineering 2 (2012), 77–89.

[Bern2006#curve25519]: https://cr.yp.to/ecdh/curve25519-20060209.pdf
[Bern2011#ed25519]: https://ed25519.cr.yp.to/ed25519-20110926.pdf
