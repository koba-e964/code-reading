import sys


d = 11
j_inv = -2**15


# Slightly modified version of https://stackoverflow.com/questions/356090/how-to-compute-the-nth-root-of-a-very-big-integer
def find_invpow(x: int, n: int) -> int:
    """Finds the integer component of the n'th root of x,
    an integer such that y ** n <= x < (y + 1) ** n.
    """
    high = 1
    while high ** n <= x:
        high *= 2
    low = high//2
    while low < high:
        mid = (low + high) // 2
        if low < mid and mid**n < x:
            low = mid
        elif high > mid and mid**n > x:
            high = mid
        else:
            return mid
    return mid + 1


# Find a prime of form dx^2 + dx + (d+1)//4, in the range [lo, hi]
def find_suitable_prime(lo: int, hi: int) -> tuple[int, int] | None:
    # Solve dx^2 + dx + (d+1)//4 >= lo
    s = find_invpow((4 * max(lo, 1) - 1) // d, 2)
    if s * s < (4 * max(lo, 1) - 1) // d:
        s += 1
    x = s // 2
    while True:
        val = x * d * (x + 1) + (d + 1) // 4
        if val > hi:
            break
        if val.is_prime():
            return (val, x)
        x += 1
    return None


def gen_anomalous_curve(p: int, j_inv: int) -> tuple[int, int]:
    a = GF(p)(-3*j_inv)/GF(p)(j_inv-1728)
    b = GF(p)(2*j_inv)/GF(p)(j_inv-1728)
    E = EllipticCurve(GF(p), [a, b])
    P = E.random_point()
    if P * p == E(0):
        return (a, b)
    r = 2
    while GF(p)(r).is_square():
        r += 1
    return (a * r**2, b * r**3)


if __name__ == "__main__":
    lo = 2**512
    hi = lo*2
    res = find_suitable_prime(lo, hi)
    if res is None:
        print("No suitable prime found.")
        sys.exit(1)
    p, x = res
    a, b = gen_anomalous_curve(p, j_inv)
    E = EllipticCurve(GF(p), [a, b])
    P = E.random_point()
    assert P * p == E(0)
    print(f'x = {x}')
    print(f'p = {d} * x * (x + 1) + {(d + 1) // 4}')
    print(f'a = {a}')
    print(f'b = {b}')
