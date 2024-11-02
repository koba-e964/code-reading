def approx_gcd(d: list[int], approx_error: int) -> int:
    """
    Returns q where d[0] ~= qx and d[i]'s are close to multiples of x.
    The caller must find (d[0] + q // 2) // q if they want to find x.
    """
    l = len(d)
    M = Matrix(ZZ, l, l)
    M[0, 0] = approx_error
    for i in range(1, l):
        M[0, i] = d[i]
        M[i, i] = -d[0]
    L = M.LLL()
    for row in L:
        if row[0] != 0:
            quot = abs(row[0] // approx_error)
            return quot
