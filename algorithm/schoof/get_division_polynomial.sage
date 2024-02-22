E = EllipticCurve(QQ, [0, 0, 0, 0, 7])
primes = [
    3,
    5,
    7,
    11,
]

for p in primes:
    print(p, E.division_polynomial(p))
