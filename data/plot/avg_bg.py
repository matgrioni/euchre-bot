#!/usr/bin/env python2

#
# Plot a bar graph comparing the averages of the different data files provided.
# Provide the results for all the players that you wish to plot and a bar graph
# will be output. The name of the category will be the name of the file.
#

import argparse
import os

import matplotlib.pyplot as plt

parser = argparse.ArgumentParser()
parser.add_argument('results', nargs='+', help='List of results to graph')
args = parser.parse_args()


avgs = []
names = []
for r in args.results:
    names.append(os.path.splitext(os.path.basename(r))[0])
    with open(r, 'r') as f:
        l = next(f).strip()
        items = l.split(',')
        avgs.append(float(items[0]))

x_pos = range(len(avgs))

plt.bar(x_pos, avgs, align='center', alpha=0.5)
plt.xticks(x_pos, names)
plt.ylabel('Difference to minimax')
plt.title('Average Difference Between Player Implementations and Minimax Player')
plt.show()
