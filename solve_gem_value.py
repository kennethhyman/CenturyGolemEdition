import csv
import numpy
from numpy.linalg import inv
import matplotlib.pyplot as plt
from matplotlib.pyplot import hist

golems = numpy.loadtxt('./data/golems.csv', skiprows=1, delimiter=',')
golems = golems[golems[:,-1].argsort()]
costs = golems[:,:-1]
points = golems[:,-1]

#print(golems)
#print(costs)
#print(str(costs.shape))
#print(str(points.shape))
#print(points)


#print(costs)
#print((points > 15).sum())
#
#cost_map = {}
#for value in points:
#    cost_map[value] = cost_map.get(value, 0) + 1
#
#print(cost_map.values())
#print(cost_map.keys())
#plt.bar(cost_map.keys(), cost_map.values(), facecolor='blue', alpha=0.5)
#plt.ylabel("Number of golems")
#plt.xlabel("Golem Value")
#plt.show()

#print(costs.shape)
#print(points.shape)
#print(numpy.linalg.solve(costs[0:4,:], points[0:4]))

a_transpose = numpy.matrix.transpose(costs)

a_trans_b = numpy.matmul(a_transpose, points)
a_trans_a = numpy.matmul(a_transpose, costs)
inv_a_trans_a = numpy.linalg.inv(a_trans_a)

print("a_trans_b: " + str(a_trans_b.shape))
print("a_trans_a: " + str(a_trans_a.shape))

approx_gem_value = numpy.matmul(a_trans_b, inv_a_trans_a)
colors = ['yellow', 'green', 'blue', 'pink']
print(dict(zip(colors,approx_gem_value)))

estimated_points = numpy.matmul(costs, approx_gem_value)
print(estimated_points.shape)
print(points.shape)
print(points - estimated_points)
print(costs[4])


yellows = costs[:,0]
greens = costs[:,1]
blues = costs[:,2]
pinks = costs[:,3]
total_gem_cost = yellows+greens+blues+pinks

width = 0.95
fig, ax = plt.subplots()

yellow_chart = ax.bar(range(len(points)), yellows, width, label='yellow gems', color="y")
green_chart = ax.bar(range(len(points)), greens, width, bottom=yellows, label='green gems', color='g')
blue_chart = ax.bar(range(len(points)), blues, width, label="blue gems", bottom=(yellows+greens), color='b')
pink_chart = ax.bar(range(len(points)), pinks, width, label="pink gems", bottom=(yellows+greens+blues), color='pink')

ax.set_ylabel('gem cost (count)')
ax.set_title('Golem cost by gem color')

#plt.axes().get_xaxis().get_ticks().set_visible(False)

for idx, point in enumerate(yellow_chart):
    plt.text(point.get_x(), total_gem_cost[idx], "%d" % points[idx], color="black", fontsize=8, fontweight="bold")

plt.show()
