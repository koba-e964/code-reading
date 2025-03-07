# Anomalous curve の生成

[[LMVV2005]] の手法を採用した。

## 使い方
```bash
sage gen.sage
```
$p \in [2^{512}, 2^{513}]$ であるような anomalous curve (であって Smart attack が可能なもの) を返す。範囲を変えたい場合はコード内の `lo` と `hi` を修正すること。

[LMVV2005]: https://infoscience.epfl.ch/entities/publication/cec333e1-711d-4224-8df0-d69d154e892f
