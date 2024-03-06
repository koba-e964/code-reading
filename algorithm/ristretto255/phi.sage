primes = [5, 13, 17, 29, 37, 41, 53]
for p in primes:
    K = GF(p)
    print(f'p = {p}')
    for A in range(p):
        if K(A^2) == 4 or K(A^2) == 1:
            continue
        d = -K(A-2)/(A+2)
        if not is_square(-K(A+2)) or is_square(d):
            continue
        print(f'  p = {p}, A = {A}, d = {d}, is_square(d) = {is_square(d)}')

        montgomery_order = 1 # one infinity
        for u in range(p):
            v2 = K(u^3 + A * u^2 + u)
            if v2 == 0:
                montgomery_order += 1
            elif is_square(v2):
                montgomery_order += 2
        print(f'    montgomery_order = {montgomery_order}')

        edwards_order = 0
        for x in range(p):
            for y in range(p):
                xx = K(x)
                yy = K(y)
                if -xx^2 + yy^2 == 1 + d * xx^2 * yy^2:
                    edwards_order += 1
        print(f'    edwards_order = {edwards_order}')

        jacobi_quartic_order = 2 # two infinities
        for s in range(p):
            for t in range(p):
                ss = K(s)
                tt = K(t)
                if tt^2 == ss^4 + A * ss^2 + 1:
                    jacobi_quartic_order += 1
        print(f'    jacobi_quartic_order = {jacobi_quartic_order}')
