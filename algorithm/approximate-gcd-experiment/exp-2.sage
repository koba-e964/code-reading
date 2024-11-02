"""
Randomly generate numbers with the same bit length and see how long approx_gcd's results are.
"""
import os
import time
from Crypto.Util.number import bytes_to_long
load("approx_gcd.sage")


def try_and_find_magnitude(num_nums: int) -> int:
    """
    Returns how big approx_gcd() is.
    """
    nums: list[int] = []
    for _ in range(num_nums):
        nums.append(bytes_to_long(os.urandom(1250)))
    q = approx_gcd(nums, 2**256)
    found_n = (nums[0] + q // 2) // q
    return found_n.bit_length()


def main() -> None:
    start = time.time()
    num_nums_min = 2
    num_nums_max = 16
    for num_nums in range(num_nums_min, num_nums_max + 1):
        print(f'# ({time.time() - start:.2f}s) {num_nums = }')
        count = 50
        error_sum = 0.0
        error_sqsum = 0.0
        for _ in range(count):
            res = try_and_find_magnitude(num_nums)
            error_sum += res
            error_sqsum += res**2
        avg = error_sum / count
        sample_std = ((error_sqsum / count - avg**2) * count / (count - 1.0)) ** 0.5
        avg_uncertainty = sample_std / count ** 0.5
        print(f'{avg = }')
        print(f'{avg_uncertainty = }')


if __name__ == "__main__":
    main()
