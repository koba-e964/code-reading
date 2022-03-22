https://arxiv.org/abs/1605.05274 の実装。

```bash
javac Turing.java 2>&1| head -n 10
```
を実行すると以下のようなエラーが出てコンパイルが異常終了する。

```


The system is out of resources.
Consult the following stack trace for details.
java.lang.StackOverflowError
        at jdk.compiler/com.sun.tools.javac.code.Types$HashCodeVisitor.visitClassType(Types.java:4188)
        at jdk.compiler/com.sun.tools.javac.code.Types$HashCodeVisitor.visitClassType(Types.java:4181)
        at jdk.compiler/com.sun.tools.javac.code.Type$ClassType.accept(Type.java:1011)
        at jdk.compiler/com.sun.tools.javac.code.Types$UnaryVisitor.visit(Types.java:4980)
        at jdk.compiler/com.sun.tools.javac.code.Types$HashCodeVisitor.visitWildcardType(Types.java:4213)
```
