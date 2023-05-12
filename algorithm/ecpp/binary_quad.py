# pylint: disable=C0103
"""
This script handles positive definite quadratic forms.
"""

import math
import num_theory

def discriminant(form: tuple[int, int, int]) -> int:
    (a, b, c) = form
    return b * b - 4 * a * c


def is_positive_form(form: tuple[int, int, int]) -> bool:
    """
    Is a form positive definite?
    """
    (a, b, c) = form
    return a >= 0 and c >= 0 and b * b < 4 * a * c


def is_reduced(abc: tuple[int, int, int]) -> bool:
    """
    Is abc a reduced binary quadratic form?
    https://mathworld.wolfram.com/ReducedBinaryQuadraticForm.html
    """
    (a, b, c) = abc
    absb = abs(b)
    return absb <= a <= c and (b >= 0 if a in (c, absb) else True)


def classes(d: int) -> list[tuple[int, int, int]]:
    """
    Find all classes with discriminant d.
    (i.e. triples (a, b, c) such that b^2 - 4ac = d and ax^2 + bxy + cy^2 is reduced)
    d must be negative.
    """
    assert d < 0
    assert d % 4 in (0, 1)
    lim = -d // 3
    absb = 0
    result = []
    while absb * absb <= lim:
        if (absb - d) % 2 != 0:
            absb += 1
            continue
        ac = (absb * absb - d) // 4
        for a in range(max(1, absb), ac + 1):
            if ac < a * a:
                break
            if ac % a != 0:
                continue
            c = ac // a
            if math.gcd(a, absb, c) != 1:
                continue
            result.append((a, absb, c))
            if a not in (c, absb) and absb != 0:
                result.append((a, -absb, c))
        absb += 1
    return result


def reduce(form: tuple[int, int, int]) -> tuple[int, int, int]:
    """
    Reduces a form.
    """
    # https://metaphor.ethz.ch/x/2021/hs/401-3110-71L/ex/Eleventh.pdf
    # Theorem 3.3
    assert is_positive_form(form)
    (a, b, c) = form
    while not is_reduced((a, b, c)):
        delta = (b + 2 * c - 1) // (2 * c)
        newb = -b + 2 * delta * c
        if newb > c:
            newb -= 2 * c
            delta -= 1
        newc = a - delta * b + delta * delta * c
        (a, b, c) = (c, newb, newc)
    return (a, b, c)


def mul_class(form1: tuple[int, int, int], form2: tuple[int, int, int]) -> tuple[int, int, int]:
    """
    Multiplies two forms and returns its reduced form.
    https://en.wikipedia.org/wiki/Binary_quadratic_form#Composition
    """
    assert discriminant(form1) == discriminant(form2)
    disc = discriminant(form1)
    (a1, b1, _) = form1
    (a2, b2, _) = form2
    bmu = (b1 + b2) // 2
    e = math.gcd(a1, a2, bmu)
    a = (a1 // e) * (a2 // e)
    (tmp, g) = num_theory.inverse_with_gcd(bmu // e, 2 * a)
    constraints = [(b1, 2 * a1 // e), (b2, 2 * a2 // e), ((disc + b1 * b2) // (2 * e) // g * tmp, 2 * a // g)]
    (b, _) = num_theory.solve_equation(constraints)
    c = (b * b - disc) // (4 * a)
    return reduce((a, b, c))


def genera(d: int) -> int:
    clz = classes(d)
    doubled = set()
    for c in clz:
        doubled.add(mul_class(c, c))
    quotient = len(clz) // len(doubled)
    return bin(quotient - 1).count("1")
