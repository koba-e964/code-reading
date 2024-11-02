import time
from Crypto.Util.number import bytes_to_long
from Crypto.PublicKey import RSA
from gen_rsa import genkey, gensig
load("approx_gcd.sage")


def try_one(rsa: RSA.RsaKey, num_sigs: int) -> bool:
    e = 11
    sigs: list[int] = []
    for _ in range(num_sigs):
        sigs.append(bytes_to_long(gensig(rsa)))
    diff = [abs(sigs[i]**e - sigs[i - 1]**e) for i in range(1, num_sigs)]
    q = approx_gcd(diff, 2**256)
    found_n = (diff[0] + q // 2) // q
    if found_n == rsa.n:
        return True
    q = (found_n + rsa.n // 2) // rsa.n
    print(f'q = round(found_n / rsa.n); {q.bit_length() = }')
    return False


def main() -> None:
    start = time.time()
    rsa = genkey()
    num_sigs_min = 14
    num_sigs_max = 14
    for num_sigs in range(num_sigs_min, num_sigs_max + 1):
        print(f'# ({time.time() - start:.2f}s) {num_sigs = }')
        success = 0
        count = 1000
        for iter in range(count):
            if iter % 10 == 0:
                print(f'({time.time() - start:.2f}s) {iter = }')
            if try_one(rsa, num_sigs):
                success += 1
        print(f'{success = }')
        print(f'{count = }')


if __name__ == "__main__":
    main()
