#!/usr/bin/env python3

# This script makes use of another module: common.py, which can be
# found on GitHub:
#
#  https://github.com/andreacorbellini/ecc/blob/master/logs/common.py
#
# You must place that module on the same directory of this script
# prior to running it.

import math
import random

from common import tinycurve as curve


def log(p, q):
    assert curve.is_on_curve(p)
    assert curve.is_on_curve(q)

    # Calculate m = sqrt(n)
    sqrt_n = int(math.sqrt(curve.n)) + 1

    print('Computing baby steps...................')
    # Compute the baby steps (jP) and store them in the 'precomputed' hash table.
    # The hash table is of size m with key-value pairs (jP, j). The table is indexed
    # by jP (key) with the corresponding entry's value set to j.
    jp = None
    precomputed = {None: 0}
    for j in range(1, sqrt_n):
        jp = curve.add(jp, p)
        precomputed[jp] = j

    # Now compute the giant steps (Q - mP) and check the hash table for any
    # matching point. Here we precompute s = -mP 
    print('Computing giant steps...................')
    jp = q
#    print('Multiply m = %d'%(sqrt_n))
#    print('with negP = ', curve.neg(p))
    s = curve.mult(sqrt_n, curve.neg(p))

    for i in range(sqrt_n):
        try:
            j = precomputed[jp]
        except KeyError:
            pass
        else:
            steps = sqrt_n + i
            logarithm = j + sqrt_n * i
            return logarithm, steps

        jp = curve.add(jp, s)

    raise AssertionError('logarithm not found')


def main():
    # Use a random number x to compute Q. Later we will use this 
    # number x to check our result y from y = log(P, Q).
    x = random.randrange(1, curve.n)
    print('Curve: {}'.format(curve))
    print('Curve order: {}'.format(curve.n))
    print('Random x = %d'%(x))

    p = curve.g
    q = curve.mult(x, p)
    print('p = (0x{:x}, 0x{:x})'.format(*p))
    print('q = (0x{:x}, 0x{:x})'.format(*q))
    print(x, '* p = q')

    y, steps = log(p, q)
    print('log(p, q) =', y)
    print('Took', steps, 'steps')

    # If x (the random number we chose) is equal to y
    # (the result of the baby-step, giant-step alg), then
    # our algorithm worked
    assert x == y


if __name__ == '__main__':
    main()
