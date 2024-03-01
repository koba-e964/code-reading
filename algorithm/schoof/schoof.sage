class ModCurve:
    """
    defining_poly: polynomial of degree 3
    mod_poly: polynomial in x modulo which all polynomial computations will be done
    example:
    p = 31
    R.<x> = PolynomialRing(GF(p))
    curve = ModCurve(R, x^3+x, 3 * x^4 + 84 * x)
    curve.add([x, 1], [x, 1])
    """
    def __init__(self, PolyRing, defining_poly, mod_poly):
        q = PolyRing.base().cardinality()
        if q % 2 == 0:
            raise "q should be odd"
        self.PolyRing = PolyRing
        self.defining_poly = defining_poly
        self.mod_poly = mod_poly

    def add(self, p1, p2):
        S = self.PolyRing.quotient(self.mod_poly)
        if p1 == None:
            return p2
        if p2 == None:
            return p1
        x1, y1 = p1
        x2, y2 = p2
        grad = None
        try:
            if x1 == x2:
                if y1 != y2:
                    return None
                grad = S(3 * x1 * x1) / S(y1 * 2) / S(self.defining_poly)
            else:
                grad = S(y2 - y1) / S(x2 - x1)
        except ZeroDivisionError:
            # If the zero division occurs, it means the result is O.
            return None
        grad = S.lift(grad)
        x3 = S(grad * grad * self.defining_poly - x1 - x2)
        y3 = S(-grad * (x3 - x1) - y1)
        return [S.lift(x3), S.lift(y3)]

    def mul(self, point, scalar):
        cur = None
        neg = scalar < 0
        if neg:
            scalar = -scalar
        for _ in range(0, scalar):
            cur = self.add(cur, point)
        if neg and cur is not None:
            cur = [cur[0], -cur[1]]
        return cur

    def frobenius(self, point):
        if point is None:
            return None
        px, py = point
        q = self.PolyRing.base().cardinality()
        S = self.PolyRing.quotient(self.mod_poly)
        new_px = S(px)^q
        new_py = S(py)^q * S(self.defining_poly)^((q-1)//2)
        return [S.lift(new_px), S.lift(new_py)]

    def on_curve(self, point):
        if point is None:
            return True
        px, py = point
        value = py^2 * self.defining_poly - self.defining_poly(px)
        return value % self.mod_poly == 0

p = 31
R.<x> = PolynomialRing(GF(p))
# division polynomial for y^2 = x^3+7 in Q
l_div_pairs = [
    (3, 3 * x^4 + 84 * x),
    (5, 5*x^12 + 2660*x^9 - 11760*x^6 - 548800*x^3 - 614656),
    (7, 7*x^24 + 27608*x^21 - 2101904*x^18 - 284585728*x^15 - 2228742656*x^12 - 26142548992*x^9 - 330576748544*x^6 - 661153497088*x^3 + 377801998336),
    (11, 11*x^60 + 650188*x^57 - 904424752*x^54 - 1362018710976*x^51 + 40649607043584*x^48 - 27469877472457728*x^45 - 4134364704189468672*x^42 - 64227041794669264896*x^39 - 8432457864221884416*x^36 + 85206511357804262981632*x^33 + 3138358694073278126882816*x^30 + 40385368497397499290451968*x^27 + 256894564925292514665037824*x^24 + 1234632848278370583136174080*x^21 + 10812377157691273558785785856*x^18 + 88664931644194734932559396864*x^15 + 324049937418286078324371357696*x^12 + 410562796998815176985413681152*x^9 + 67701756866706860751422750720*x^6 - 758259676907116840415934808064*x^3 - 87732524600823436081182539776),
    (3, x),
    (3, x^3 + 28),
    (7, x+4),
    (7, x+7),
    (7, x+20),
    (7, x^3+25),
]

for l, div in l_div_pairs:
    print(f"l = {l}")
    print(f"  div = {factor(div)}")
    curve = ModCurve(R, x^3+7, div)
    cur = None
    point = [x % div, 1]
    for i in range(1, l + 1):
        cur = curve.add(cur, point)
        print(i, cur, curve.on_curve(cur))
    assert curve.mul(point, l) is None

    fr = curve.frobenius(point) # (x^{q}, y^{q})
    fr2 = curve.frobenius(fr) # (x^{q^2}, y^{q^2})
    qp = curve.mul(point, p) # q(x, y)
    lhs = curve.add(fr2, qp)
    print(f"  phi(x) = {fr} {curve.on_curve(fr)}")
    print(f"  phi^2(x) = {fr2} {curve.on_curve(fr2)}")
    print(f"  q x = {qp}")
    print(f"  phi^2(x) + q x = {lhs}")
    for t in range(-(l-1)//2, (l-1)//2+1):
        rhs = curve.mul(fr, t)
        if l <= 3:
            print(f"    {lhs} {t} {rhs}")
        if lhs == rhs:
            print(f"      Found: {t} mod {l}")
