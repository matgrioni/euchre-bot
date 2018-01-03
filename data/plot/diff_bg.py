#!/usr/bin/env python2

#
# Plots the distribution of differences for the given player implementations.
# Plots all the distributions right next to each other.
#

import argparse
import os

import matplotlib
import matplotlib.pyplot as plt
import numpy


WIDTH = 0.27
COLORS = sorted(matplotlib.colors.ColorConverter.colors.keys())


parser = argparse.ArgumentParser()
parser.add_argument('results', nargs='+',
                    help='The location of the results to display')
args = parser.parse_args()

if len(args.results) > len(COLORS):
    raise ArgumentError('Too many results. Can at most be {}'.format(len(COLORS)))


dists = []
keys = []
for r in args.results:
    keys.append(os.path.splitext(os.path.basename(r))[0])

    with open(r, 'r') as f:
        l = next(f).strip()
        items = l.split(',')

        d = map(lambda p: float(p), items[1:])
        dists.append(d)


x_pos = numpy.arange(len(dists[0]))

fig = plt.figure()
ax = fig.add_subplot(111)

ids = []
for (i, dist) in enumerate(dists):
    rect = ax.bar(x_pos + i * WIDTH, dist, WIDTH, color=COLORS[i],
                  align='center')
    ids.append(rect[0])

ax.set_ylabel('Percentage')
ax.set_xlabel('Difference between Player and Minimax')
ax.set_xticks(x_pos + WIDTH)
ax.set_xticklabels(range(7))
ax.legend(ids, keys)
plt.title('Distribution of Differences for Player Implementations')

plt.show()
