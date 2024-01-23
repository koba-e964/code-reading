[go@go1.21.6 の regexp パッケージ](https://cs.opensource.google/go/go/+/refs/tags/go1.21.6:src/regexp/) を読む。

## regexp/syntax
依存関係は以下。なおデータ型の変換を担当する関数・ファイルは辺で表している。
![./go-regexp-syntax-dep.png](./go-regexp-syntax-dep.png)

それぞれのファイルの簡単な説明は以下。

- regexp.go: 型 syntax.Regexp を公開する。
- prog.go: 型 syntax.Prog を公開する。
- parse.go: 正規表現の文字列を syntax.Regexp に変換する。また型 charGroup を定義する。
- simplify.go: syntax.Regexp の単純化を担当する。
- perl_groups.go: Perl 互換用のグループ (charGroup) を定義する。
- compile.go: syntax.Regexp を syntax.Prog に変換する。

### parse.go
正規表現の文字列を syntax.Regexp に変換する。また型 charGroup を定義する。

正規表現はほとんど LL(1) (直観: 1 文字読めば構文木のトップレベルの要素がわかる) 文法なので、以下の要領で構文解析ができる。

- 構文解析の要素にあたる関数を `func Elem(t string) (value 値の型, rest string, err error)` で定義する。
- このように、文法に従って再帰的に構文解析するパーサーを作る。
  ```go
  func Expr(t string) (e *expr, rest string, err error) {
    if t[0] == '+' {
        val, rest, err := UnaryPlus(t)
        if err != nil {
            return nil, "", err
        }
        t = rest
        return ..., t, nil
    }
    ...
  }
  ```

ファイル内の大まかな処理の内容は以下の通り:

- [L14-L124](https://cs.opensource.google/go/go/+/refs/tags/go1.21.6:src/regexp/syntax/parse.go;l=14-124): 型や定数の定義。
- [L126-L289](https://cs.opensource.google/go/go/+/refs/tags/go1.21.6:src/regexp/syntax/parse.go;l=126-289): parser 自体の定義と便利メソッド。
- [L291-L575](https://cs.opensource.google/go/go/+/refs/tags/go1.21.6:src/regexp/syntax/parse.go;l=291-575): parser 内部の stack を扱うメソッド。
- [L577-L883](https://cs.opensource.google/go/go/+/refs/tags/go1.21.6:src/regexp/syntax/parse.go;l=577-883): parser 内部で使う便利メソッド。
- [L885-L1823](https://cs.opensource.google/go/go/+/refs/tags/go1.21.6:src/regexp/syntax/parse.go;l=885-1823): syntax.Parse の本丸。`(値 値の型, rest string, err error)` を返す関数が多い。
- [L1825-L2115](https://cs.opensource.google/go/go/+/refs/tags/go1.21.6:src/regexp/syntax/parse.go;l=1825-2115): 文字類 (class) を扱う関数。
