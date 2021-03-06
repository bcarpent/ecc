#!/usr/bin/env python3

# This script makes use of another module: common.py, which can be
# found on GitHub:
#
#  https://github.com/andreacorbellini/ecc/blob/master/logs/common.py
#
# You must place that module on the same directory of this script
# prior to running it.

import random

from common import inverse_mod, tinycurve as curve


class PollardRhoSequence:

    def __init__(self, point1, point2):
        self.point1 = point1
        self.point2 = point2

        # Generate random pair (a1, b1)
        self.add_a1 = random.randrange(1, curve.n)
        self.add_b1 = random.randrange(1, curve.n)

#TBD
#        self.add_a1 = 3
#        self.add_b1 = 4

        # Compute a1P + b1Q
        print('')
        print('')
        print('')
        print('Generating new sequence....................')
        print('a1 = ', self.add_a1)
        print('b1 = ', self.add_b1)
        self.a1P = curve.mult(self.add_a1, point1)
        print('a1P = ', self.a1P)
        print('')
        print('')
        self.b1Q = curve.mult(self.add_b1, point2)
        print('b1Q = ', self.b1Q)
        print('')
        print('')
        self.add_x1 = curve.add(self.a1P, self.b1Q)

#        self.add_x1 = curve.add(
#            curve.mult(self.add_a1, point1),
#            curve.mult(self.add_b1, point2),
#        )

        print('X1 = ', self.add_x1)
        print('')
        print('')

        # Generate random pair (a2, b2)
        self.add_a2 = random.randrange(1, curve.n)
        self.add_b2 = random.randrange(1, curve.n)

#TBD
#        self.add_a2 = 3
#        self.add_b2 = 3

        # Compute a2P + b2Q
        print('a2 = ', self.add_a2)
        print('b2 = ', self.add_b2)
        self.a2P = curve.mult(self.add_a2, point1)
        self.b2Q = curve.mult(self.add_b2, point2)
        print('a2P = ', self.a2P)
        print('b2Q = ', self.b2Q)
        self.add_x2 = curve.add(self.a2P, self.b2Q)

#        self.add_x2 = curve.add(
#            curve.mult(self.add_a2, point1),
#            curve.mult(self.add_b2, point2),
#        )

        print('X2 = ', self.add_x2)
        print('')
        print('')

    def __iter__(self):
        # Partition the curve into 3 segments
        partition_size = curve.p // 3 + 1

        x = None
        a = 0
        b = 0

        # Start with 0P, then add 1P to get 1P, where i == 0.
        # Next iteration is 2P (doubling 1P to get 2P), where i == 1.
        # Next iteration is 3P (adding P to 2P), where i == 2.
        while True:
            if x is None:
                i = 0
                print('Partition size = ', partition_size)
            else:
                i = x[0] // partition_size
                print('Partition size = ', partition_size)

            if i == 0:
                # x is either the point at infinity (None), or is in the first
                # third of the plane (x[0] <= curve.p / 3).
                a += self.add_a1
                b += self.add_b1
                print('Iterating, i = 0, a =', a, ', b =', b)
                x = curve.add(x, self.add_x1)
            elif i == 1:
                # x is in the second third of the plane
                # (curve.p / 3 < x[0] <= curve.p * 2 / 3).
                a *= 2
                b *= 2
                print('Iterating, i = 1, a =', a, ', b =', b)
                x = curve.double(x)
            elif i == 2:
                # x is in the last third of the plane (x[0] > curve.p * 2 / 3).
                a += self.add_a2
                b += self.add_b2
                print('Iterating, i = 2, a =', a, ', b =', b)
                x = curve.add(x, self.add_x2)
            else:
                raise AssertionError(i)

            print('x =', x)
            a = a % curve.n
            b = b % curve.n
            print('a, b = ', a, b)
            print('')
            print('')
            yield x, a, b


def log(p, q, counter=None):
    assert curve.is_on_curve(p)
    assert curve.is_on_curve(q)

    # Pollard's Rho may fail sometimes: it may find a1 == a2 and b1 == b2,
    # leading to a division by zero error. Because PollardRhoSequence uses
    # random coefficients, we have more chances of finding the logarithm
    # if we try again, without affecting the asymptotic time complexity.
    # We try at most three times before giving up.
    for i in range(3):
        sequence = PollardRhoSequence(p, q)

        tortoise = iter(sequence)
        hare = iter(sequence)

        print('Iterate over the sequences...................')
        # The range is from 0 to curve.n - 1, but actually the algorithm will
        # stop much sooner (either finding the logarithm, or failing with a
        # division by zero).
        for j in range(curve.n):
            print('One tortoise step:')
            x1, a1, b1 = next(tortoise)

            # Hare skips over a step
            print('Two hare steps')
            x2, a2, b2 = next(hare)
            x2, a2, b2 = next(hare)

            print('a1 =', a1, ', b1 =', b1, ', a2 =', a2, ' b2 =', b2)

            # We have found a match (detected a cycle) if x1 == x2
            if x1 == x2:
                if b1 == b2:
                    # This would lead to a division by zero. Try with
                    # another random sequence.
                    break

                x = (a1 - a2) * inverse_mod(b2 - b1, curve.n)
                logarithm = x % curve.n
                steps = i * curve.n + j + 1
                return logarithm, steps

    raise AssertionError('logarithm not found')


def main():
    x = random.randrange(1, curve.n)
    p = curve.g
    q = curve.mult(x, p)

    print('Curve: {}'.format(curve))
    print('Curve order: {}'.format(curve.n))
    print('p = (0x{:x}, 0x{:x})'.format(*p))
    print('q = (0x{:x}, 0x{:x})'.format(*q))
    print(x, '* p = q')

    y, steps = log(p, q)
    print('log(p, q) =', y)
    print('Took', steps, 'steps')

    assert x == y


if __name__ == '__main__':
    main()
