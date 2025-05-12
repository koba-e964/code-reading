R.<x,a> = PolynomialRing(QQ, 'x,a')

def psi2(E, n: int):
    return E.division_polynomial(n, x=x, two_torsion_multiplicity=0) \
        * E.division_polynomial(n, x=x, two_torsion_multiplicity=2)

def phi(E, n: int):
    first = x * psi2(E, n)
    second = E.division_polynomial(n + 1, x=x, two_torsion_multiplicity=0) \
        * E.division_polynomial(n - 1, x=x, two_torsion_multiplicity=2)
    return first - second


if __name__ == "__main__":
    E = EllipticCurve([0, a, 0, 1, 0])
    for n in range(2, 10):
        print(f'{n = }')
        if n == 2:
            print('  psi_{2*n}^2 =', psi2(E, 2*n))
            print(f'  phi_{2*n} =', phi(E, 2*n))
        phi_n = phi(E, n)
        psi_n_2 = psi2(E, n)
        if n == 2:
            print(f'  psi_{n}^2 =', psi_n_2)
            print(f'  phi_{n} =', phi_n)
        dbl_z = 4 * phi_n * (phi_n * phi_n + a * phi_n * psi_n_2 + psi_n_2 ** 2) * psi_n_2
        if n == 2:
            print(f'  DBL_z(phi_{n}, psi_{n}^2) =', dbl_z)
        print(f'  {psi2(E, 2*n) == dbl_z = }')
        dbl_x = (phi_n ** 2 - psi_n_2 ** 2) ** 2
        if n == 2:
            print(f'  DBL_x(phi_{n}, psi_{n}^2) =', dbl_x)
        print(f'  {phi(E, 2*n) == dbl_x = }')
