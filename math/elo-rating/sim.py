"""
With ChatGPT's advice
https://chatgpt.com/share/67b03256-93e8-8010-9d52-f480559a1246
"""

import numpy as np
import matplotlib.pyplot as plt


def elo_rating_change(ra: float, rb: float, k: float, result: float) -> float:
    """
    Calculate the new Elo rating for player A after a match.
    ra: Current rating of player A
    rb: Current rating of player B
    k: K-factor (default 32)
    result: 1 if player A wins, 0 if loses, 0.5 if draw
    """
    ea = 1 / (1 + 10 ** ((rb - ra) / 400))
    new_ra = ra + k * (result - ea)
    return new_ra


def simulate_elo(true_skill: float, initial_rating: float, rounds: int,
                 opponents_rating: float, k: float) -> np.ndarray:
    """
    Simulate the Elo rating changes over several rounds against opponents.
    true_skill: The actual skill level of the player
    initial_rating: Initial Elo rating of the player
    rounds: Number of rounds to simulate
    opponents_rating: Rating of the opponents
    K: K-factor
    """
    ratings = [initial_rating]
    for _ in range(1, rounds + 1):
        # Probability of winning based on true skill vs opponent's rating
        win_prob = 1 / (1 + 10 ** ((opponents_rating - true_skill) / 400))
        result = np.random.binomial(1, win_prob)  # 1 if win, 0 if lose
        new_rating = elo_rating_change(
            ratings[-1], opponents_rating, k, result)
        ratings.append(new_rating)
    return np.array(ratings)


def estimate_normal_params(data: np.ndarray) -> tuple[float, float]:
    """
    Estimate the mean and variance of a dataset assuming a normal distribution.
    data: Input data array
    Returns: (mean, variance)
    """
    mean = np.mean(data)
    variance = np.var(data, ddof=1)  # Unbiased variance
    return mean, variance


# Simulation parameters
TRUE_SKILL = 1000
INITIAL_RATING = TRUE_SKILL
ROUNDS_TO_CHECK = [1000, 10000, 100000, 1000000]
OPPONENTS_RATINGS = [800, 1000, 1200]
KS = [16, 32, 64]


def main() -> None:
    """
    main function
    """
    rounds = max(ROUNDS_TO_CHECK)
    # Run simulation
    for opponents_rating in OPPONENTS_RATINGS:
        for k in KS:
            ratings = simulate_elo(TRUE_SKILL, INITIAL_RATING,
                                rounds, opponents_rating, k)

            # Plot rating distribution at different rounds
            print(f'\nOpponent Rating: {opponents_rating}, K: {k}')
            for r in ROUNDS_TO_CHECK:
                mean, variance = estimate_normal_params(ratings[:r+1])
                print(f'Round {r}: Mean={mean:.2f}, Variance={variance:.2f}, '
                    f'Std Dev={np.sqrt(variance):.2f}')

            plt.figure(figsize=(10, 5))
            for r in ROUNDS_TO_CHECK:
                plt.hist(ratings[:r+1], bins=30, alpha=0.5, label=f'Round {r}')

            plt.xlabel('Elo Rating')
            plt.ylabel('Frequency')
            plt.title(
                f'Elo Rating Distribution at Different Rounds (vs {opponents_rating}) (K={k})')
            plt.legend()
            plt.grid(True)
            plt.savefig(f'elo_rating_distribution-vs-{opponents_rating}-k-{k}.png')


if __name__ == '__main__':
    main()
