def shift(poly, root):
    K = poly.base_ring()
    R.<x> = PolynomialRing(K, 'x')
    cur = K(1)
    ans = K(0)
    for i in range(poly.degree() + 1):
        ans += cur * poly[i]
        cur *= x + root
    return ans


def to_montgomery(E):
    a1, a2, a3, a4, a6 = E.a_invariants()
    assert a1 == 0
    assert a3 == 0
    K = E.base()
    p = K.cardinality()
    R.<x> = PolynomialRing(K, 'x')
    poly = x**3 + K(a2)*x**2 + K(a4)*x + K(a6)
    assert E.order() == p + 1 # E is supersingular
    assert p % 8 == 3
    assert E.two_torsion_rank() >= 1

    root = poly.roots()[0][0]
    poly = shift(poly, root)
    coef1 = poly[1]
    scaler = coef1.sqrt() ** ((p + 1) // 2)
    assert scaler.is_square()
    return EllipticCurve(K, [0, poly[2] / scaler, 0, 1, 0])


E1 = EllipticCurve(GF(83), [0, 11, 0, 1, 0])
E = E1
for i in range(9):
    (gen,) = E.gens()
    order3 = gen * (E.order() // 3)
    E = E.isogeny(order3).codomain()
    E = to_montgomery(E)
    print(i, E)
