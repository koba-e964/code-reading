# pylint: disable=C0103
"""
This script handles number theoretic functions.
"""

from typing import Optional

def ext_gcd(x: int, y: int) -> tuple[int, int, int]:
    """
    Returns (g, a, b) s.t. g = x * a + y * b > 0
    """
    if y == 0:
        return (x, 1, 0) if x > 0 else (-x, -1, 0)
    q = x // y
    r = x % y
    (g, a, b) = ext_gcd(y, r)
    return (g, b, a - b * q)


def inverse_with_gcd(x: int, y: int) -> tuple[int, int]:
    """
    Returns a that satisfies ax = g (mod y) along with the minimum g with which this is possible.
    """
    assert y > 0
    (g, a, _) = ext_gcd(x % y, y)
    return (a % (y // g), g)

def solve_equation_2(ab: tuple[int, int], cd: tuple[int, int]) -> Optional[tuple[int, int]]:
    """
    >>> num_theory.solve_equation_2((1, 6), (3, 8))
    (19, 24)
    """
    (a, b) = ab
    (c, d) = cd
    (bm1, g) = inverse_with_gcd(b, d)
    (dm1, _) = inverse_with_gcd(d, b)
    if a % g != c % g:
        return None
    value = a % g + (a // g * dm1 * d + c // g * bm1 * b)
    lcm = b // g * d
    return (value % lcm, lcm)


def solve_equation(exprs: list[tuple[int, int]]) -> Optional[tuple[int, int]]:
    cur = (0, 1)
    for cons in exprs:
        cur = solve_equation_2(cur, cons)
        if cur is None:
            return None
    return cur
