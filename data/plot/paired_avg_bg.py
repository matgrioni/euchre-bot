#!/usr/bin/env python2

#
# Put the single and paired average results next to each other in a bar graph.
# Provide the pairs as a flat list, with every two consecutive files
# corresponding the unpaired / paired results.
#


import argparse
import os

import matplotlib.pyplot as plt
import numpy


WIDTH = 0.27
TYPES = ('unpaired', 'paired')
COLORS = ('r', 'b')

parser = argparse.ArgumentParser()
parser.add_argument('results', nargs='+', help='List of pairs of results.')
args = parser.parse_args()


unpaired = []
paired = []
keys = []
for i, r in enumerate(args.results):
    if i % 2 == 0:
        keys.append(os.path.splitext(os.path.basename(r))[0])

    with open(r, 'r') as f:
        line = next(f).strip()
        items = line.split(',')
        avg = float(items[0])

        if i % 2 == 0:
            unpaired.append(avg)
        else:
            paired.append(avg)

x_pos = numpy.arange(len(unpaired))


fig = plt.figure()
ax = fig.add_subplot(111)

ids = []
for i, dataset in enumerate((unpaired, paired)):
    rect = ax.bar(x_pos + i * WIDTH, dataset, WIDTH, color=COLORS[i],
                  align='center')
    ids.append(rect[0])


ax.set_xlabel('Player Implementation')
ax.set_ylabel('Difference to Minimax')
ax.set_xticks(x_pos + WIDTH / 2)
ax.set_xticklabels(keys)
ax.legend(ids, TYPES)
plt.title('Average Difference Between Player Implementations and Minimax Player for Paired and Unpaired Tests')

plt.show()
