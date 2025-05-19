u = -0xd201000000010000
r = u ** 4 - u ** 2 + 1
p = r * (u - 1) ** 2 // 3 + u
K1 = GF(p)
K2 = GF(p ** 2)
i = K2(-1).sqrt()
E1 = EllipticCurve(K1, [0, 0, 0, 0, 4])
E2 = EllipticCurve(K2, [0, 0, 0, 0, 4 * (1 + i)])

assert p.is_prime()
assert r.is_prime()

t1 = u + 1
y = (u - 1) * (2 * u**2 - 1) / 3
E1.set_order(p - u)
h1 = (p - u) // r
assert p - u == h1 * r

t2 = t1**2 - 2*p
t2p = (t2 - 3*t1*y)/2
h2 = (p**2 + 1 - t2p) // r
E2.set_order(h2 * r)
