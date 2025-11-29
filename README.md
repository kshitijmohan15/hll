### HyperLogLog

Imagine a few consecutive coin flips. Each flip is independent of the next one. There is no memory, no weights, they are virtually just as random as the last one. These are the kinds of processes that can be described using Brownian process, or a "random walk". 

In a single flip, the probability of observing Tails (0) is $1/2$ (50%).
The probability of observing 2 consecutive Tails (00) is $(1/2)^2 = 1/4$ (25%).
The probability of observing $k$ consecutive Tails is $(1/2)^k$.

If the experimenter reports a maximum run of 2 Tails, you might infer they flipped the coin only a handful of times. A run of 2 is common; it happens every 4 sequences on average. If the experimenter reports a maximum run of 50 Tails, you would immediately deduce that they must have been flipping the coin for a very long time. The probability of flipping 50 consecutive tails is $1 / 2^{50}$, roughly one in a quadrillion. It is statistically impossible to observe such a rare event in a small sample size.

This is the core mechanism of HyperLogLog. We replace the coin with a hash function. A good hash function maps inputs (users, IP addresses) to a sequence of bits that are indistinguishable from random coin flips. By observing the maximum number of leading zeros in the binary representation of these hash values, we can estimate the cardinality of the input stream.

