# This code is for experimental purpose.
# DO NOT USE THIS IN PRODUCTION.
# pylint: disable=C0103

import sys
import math
from Crypto.Util import number

def xy_to_zw(mo: int, point: tuple[int, int]) -> tuple[int, int]:
    (x, y) = point
    w = number.inverse(-y, mo)
    z = x * w % mo
    return (z, w)

class ECZW:
    def __init__(self, mo: int, a1: int, a2: int, a3: int, a4: int, a6: int):
        """y^2 + a1 * x * y + a3 * y = x^3 + a2 * x^2 + a4 * x + a6
        w = z^3 + a1 * z * w + a2 * z^2 * w + a3 * w^2 + a4 * z * w^2 + a6 * w^3
        """
        self.mo = mo
        self.a1 = a1
        self.a2 = a2
        self.a3 = a3
        self.a4 = a4
        self.a6 = a6

    @staticmethod
    def simplified(mo: int, a: int, b: int):
        """Simplified form: y^2 = x^3 + a * x + b
        w = z^3 + a * z * w^2 + b * w^3
        """
        return ECZW(mo, 0, 0, 0, a, b)

    def is_on(self, point: tuple[int, int]) -> bool:
        return self.g(point) == 0

    def g(self, point: tuple[int, int]) -> bool:
        mo = self.mo
        a1 = self.a1
        a2 = self.a2
        a3 = self.a3
        a4 = self.a4
        a6 = self.a6
        (z, w) = point
        rhs = z * z * z + a1 * z * w + a2 * z * z * w + a3 * w * w + a4 * z * w * w + a6 * w * w * w
        return (rhs - w) % mo

    def g_z(self, point: tuple[int, int]) -> int:
        """∂g/∂z(point)
        """
        (z, w) = point
        mo = self.mo
        a1 = self.a1
        a2 = self.a2
        a4 = self.a4
        return (3 * z * z + a1 * w + 2 * a2 * z * w + a4 * w * w) % mo

    def g_w(self, point: tuple[int, int]) -> int:
        """∂g/∂w(point)
        """
        (z, w) = point
        mo = self.mo
        a1 = self.a1
        a2 = self.a2
        a3 = self.a3
        a4 = self.a4
        a6 = self.a6
        return (a1 * z + a2 * z * z + 2 * a3 * w + 2 * a4 * z * w + 3 * a6 * w * w - 1) % mo

    def inv(self, p: tuple[int, int]) -> tuple[int, int]:
        """Computes -p
        """
        mo = self.mo
        a1 = self.a1
        a3 = self.a3
        (z, w) = p
        invden = number.inverse(a1 * z + a3 * w - 1, mo)
        return (z * invden % mo, w * invden % mo)

    def add(self, p1: tuple[int, int], p2: tuple[int, int]) -> tuple[int, int]:
        """Computes p1 + p2
        """
        mo = self.mo
        a1 = self.a1
        a2 = self.a2
        a3 = self.a3
        a4 = self.a4
        a6 = self.a6
        (z1, w1) = p1
        (z2, w2) = p2
        lam = None
        invlam = None
        if z1 == z2 and w1 == w2:
            nom = self.g_z(p1)
            den = -self.g_w(p1) % mo
            if math.gcd(den, mo) != 1:
                invlam = den * number.inverse(nom, mo) % mo
            else:
                lam = nom * number.inverse(den, mo) % mo
        elif math.gcd(abs(z2 - z1), mo) != 1:
            invlam = (z2 - z1) * number.inverse(w2 - w1, mo) % mo
        else:
            lam = (w2 - w1) * number.inverse(z2 - z1, mo) % mo
        if lam is not None:
            nu = (w1 - z1 * lam) % mo
            zsum = -(a1 * lam + a2 * nu + a3 * lam * lam + 2 * a4 * lam * nu + 3 * a6 * lam * lam * nu) \
                * number.inverse(1 + lam * (a2 + lam * (a4 + a6 * lam)), mo)
            z3 = -(z1 + z2 - zsum) % mo
            w3 = (lam * z3 + nu) % mo
        elif invlam is not None:
            mu = (z1 - invlam * w1) % mo
            wsum = -number.inverse(a6 + invlam * (a4 + invlam * (a2 + invlam)), mo) \
                * (a3 + mu * (a4 + 2 * a2 * invlam) + a1 * invlam + 3 * invlam * invlam * mu)
            w3 = -(w1 + w2 - wsum)
            w3 %= mo
            z3 = (invlam * w3 + mu) % mo
        else:
            z = z1
            z3 = z
            # TODO: a6 != 0 must hold
            wsum = (a3 + a4 * z) * number.inverse(-a6, mo) % mo
            w3 = (wsum - w1 - w2) % mo
        return self.inv((z3, w3))

    def mul(self, x: int, p: tuple[int, int]) -> tuple[int, int]:
        """Computes x * p
        """
        result = (0, 0)
        cur = p
        while x > 0:
            if x % 2 == 1:
                result = self.add(result, cur)
            cur = self.add(cur, cur)
            x //= 2
        return result

    def lift(self, less_mo: int, p: tuple[int, int]) -> tuple[int, int]:
        """Hensel lifting to mod less_mo^2
        """
        assert less_mo * less_mo == self.mo
        mo = self.mo
        g_z = self.g_z(p)
        g_w = self.g_w(p)
        (z, w) = p
        if g_z % less_mo != 0:
            newz = (z - self.g(p) * number.inverse(g_z, mo)) % mo
            assert self.is_on((newz, w))
            return (newz, w)
        neww = (w - self.g(p) * number.inverse(g_w, mo)) % mo
        assert self.is_on((z, neww))
        return (z, neww)


def test():
    # y^2 = x^3 + 3 (mod 7), order = 13
    ec = ECZW.simplified(7, 0, 3)
    p = xy_to_zw(7, (1, 2))
    assert ec.mul(13, p) == (0, 0)
    c = (0, 0)
    for i in range(13):
        c = ec.add(c, p)
        assert ec.is_on(c)
    assert c == (0, 0)
    # An example found in https://www.hpl.hp.com/techreports/97/HPL-97-128.pdf
    mo = 43
    ec = ECZW(mo, 0, -4, 0, -128, -432)
    order = 0
    for i in range(43):
        for j in range(43):
            if ec.is_on((i, j)):
                order += 1
    assert order == 43
    p = xy_to_zw(mo, (0, 16))
    assert ec.is_on(p)
    assert ec.mul(43, p) == (0, 0)
    c = (0, 0)
    for i in range(43):
        c = ec.add(c, p)
        assert ec.is_on(c)

    # Testing of Smart's attack
    mo = 43 * 43
    ec = ECZW(mo, 0, -4, 0, -128, -432) # EC mod p^2
    ec_red = ECZW(43, 0, -4, 0, -128, -432) # reduced mod p
    p = xy_to_zw(43, (0, 16))
    lp = ec.lift(43, p)
    assert ec.is_on(lp)
    for r in range(10, 43):
        q = ec_red.mul(r, p)
        lq = ec.lift(43, q)
        plp = ec.mul(43, lp)
        plq = ec.mul(43, lq)
        v1 = plp[0] // 43
        v2 = plq[0] // 43
        disclog = v2 * number.inverse(v1, 43) % 43
        assert disclog == r


def main():
    if len(sys.argv) >= 2 and sys.argv[1] == "test":
        test()
        return

    mo = 43 * 43
    ec = ECZW(mo, 0, -4, 0, -128, -432) # EC mod p^2
    ec_red = ECZW(43, 0, -4, 0, -128, -432) # reduced mod p
    p = xy_to_zw(43, (0, 16))
    lp = ec.lift(43, p)
    assert ec.is_on(lp)
    q = xy_to_zw(43, (12, 1))
    lq = ec.lift(43, q)
    assert ec.is_on(lq)
    plp = ec.mul(43, lp)
    plq = ec.mul(43, lq)
    v1 = plp[0] // 43
    v2 = plq[0] // 43
    disclog = v2 * number.inverse(v1, 43) % 43
    print(f"discrete log = {disclog}")
    assert ec_red.mul(disclog, p) == q


if __name__ == "__main__":
    main()
