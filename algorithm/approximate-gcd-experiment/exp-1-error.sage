import time
from Crypto.Util.number import bytes_to_long
from Crypto.PublicKey import RSA
from gen_rsa import genkey, gensig
load("approx_gcd.sage")


def try_and_find_error(rsa: RSA.RsaKey, num_sigs: int) -> int:
    """
    Returns how many bits bigger n is estimated than the actual n.
    """
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
    return q.bit_length()


def main() -> None:
    start = time.time()
    rsa = genkey()
    num_sigs_min = 3
    num_sigs_max = 13
    for num_sigs in range(num_sigs_min, num_sigs_max + 1):
        print(f'# ({time.time() - start:.2f}s) {num_sigs = }')
        count = 30
        error_sum = 0.0
        error_sqsum = 0.0
        for _ in range(count):
            res = try_and_find_error(rsa, num_sigs)
            error_sum += res
            error_sqsum += res**2
        avg = error_sum / count
        sample_std = ((error_sqsum / count - avg**2) * count / (count - 1.0)) ** 0.5
        avg_uncertainty = sample_std / count ** 0.5
        print(f'{avg = }')
        print(f'{avg_uncertainty = }')


if __name__ == "__main__":
    main()
