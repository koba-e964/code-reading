# ELO Rating
ChatGPT に作ってもらった。
https://chatgpt.com/share/67b03256-93e8-8010-9d52-f480559a1246

K = 32 のときレーティングの標準偏差は 53 程度。

```
$ python3 sim.py

Opponent Rating: 800, K: 16
Round 1000: Mean=983.97, Variance=1284.77, Std Dev=35.84
Round 10000: Mean=997.33, Variance=1557.39, Std Dev=39.46
Round 100000: Mean=1001.19, Variance=1468.55, Std Dev=38.32
Round 1000000: Mean=1002.04, Variance=1401.45, Std Dev=37.44

Opponent Rating: 800, K: 32
Round 1000: Mean=992.98, Variance=2108.27, Std Dev=45.92
Round 10000: Mean=1010.15, Variance=2728.36, Std Dev=52.23
Round 100000: Mean=1005.00, Variance=2869.49, Std Dev=53.57
Round 1000000: Mean=1005.14, Variance=2865.02, Std Dev=53.53

Opponent Rating: 800, K: 64
Round 1000: Mean=991.25, Variance=5151.41, Std Dev=71.77
Round 10000: Mean=1001.10, Variance=6131.46, Std Dev=78.30
Round 100000: Mean=1006.93, Variance=5779.94, Std Dev=76.03
Round 1000000: Mean=1008.44, Variance=5832.55, Std Dev=76.37

Opponent Rating: 1000, K: 16
Round 1000: Mean=993.14, Variance=764.01, Std Dev=27.64
Round 10000: Mean=996.74, Variance=1407.75, Std Dev=37.52
Round 100000: Mean=1000.40, Variance=1438.47, Std Dev=37.93
Round 1000000: Mean=999.76, Variance=1420.84, Std Dev=37.69

Opponent Rating: 1000, K: 32
Round 1000: Mean=998.70, Variance=3021.17, Std Dev=54.97
Round 10000: Mean=1002.54, Variance=2734.68, Std Dev=52.29
Round 100000: Mean=1000.72, Variance=2809.18, Std Dev=53.00
Round 1000000: Mean=1000.28, Variance=2906.13, Std Dev=53.91

Opponent Rating: 1000, K: 64
Round 1000: Mean=993.52, Variance=5979.86, Std Dev=77.33
Round 10000: Mean=997.05, Variance=5886.84, Std Dev=76.73
Round 100000: Mean=1000.35, Variance=5911.23, Std Dev=76.88
Round 1000000: Mean=999.66, Variance=6027.58, Std Dev=77.64

Opponent Rating: 1200, K: 16
Round 1000: Mean=1008.75, Variance=1196.17, Std Dev=34.59
Round 10000: Mean=998.22, Variance=1486.51, Std Dev=38.56
Round 100000: Mean=999.62, Variance=1459.92, Std Dev=38.21
Round 1000000: Mean=997.89, Variance=1416.97, Std Dev=37.64

Opponent Rating: 1200, K: 32
Round 1000: Mean=997.51, Variance=2288.64, Std Dev=47.84
Round 10000: Mean=996.23, Variance=3519.89, Std Dev=59.33
Round 100000: Mean=995.02, Variance=2919.62, Std Dev=54.03
Round 1000000: Mean=995.76, Variance=2873.78, Std Dev=53.61

Opponent Rating: 1200, K: 64
Round 1000: Mean=984.61, Variance=5637.88, Std Dev=75.09
Round 10000: Mean=996.18, Variance=5737.06, Std Dev=75.74
Round 100000: Mean=990.43, Variance=5891.47, Std Dev=76.76
Round 1000000: Mean=991.79, Variance=5849.84, Std Dev=76.48
```
