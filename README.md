### HyperLogLog

Imagine a few consecutive coin flips. Each flip is independent of the next one. There is no memory, no weights, they are virtually just as random as the last one. These are the kinds of processes that can be described using Brownian process, or a "random walk". 

In a single flip, the probability of observing Tails (0) is $1/2$ (50%).
The probability of observing 2 consecutive Tails (00) is $(1/2)^2 = 1/4$ (25%).
The probability of observing $k$ consecutive Tails is $(1/2)^k$.