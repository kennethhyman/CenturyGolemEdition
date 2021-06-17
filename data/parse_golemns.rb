golem_file = File.open('./golems.txt', 'r')
out_file = File.open('./golems.csv', 'w')

COLORS = {y: :yellow, g: :green, b: :blue, p: :pink}
KEYS = COLORS.values + [:points]

def parse_golem(golem)
  puts golem
  costs = golem.scan(/[0-9][y,g,b,p]/)
  value = golem.match(/([0-9]+)$/)[1]

  return costs_to_hash(costs).merge!(points: value.to_i)
end


def costs_to_hash(cost_arr)
  return cost_arr.each_with_object({}) do |cost_str, hash|
    color = COLORS[cost_str[1].to_sym]
    hash[color] = cost_str[0].to_i
  end
end

golem_file.each_line do |golem|
  golem_hash = parse_golem(golem)
  puts golem_hash
  golem_row = KEYS.each_with_object([]) { |key, arr| arr << (golem_hash[key] || 0).to_s}
  puts golem_row.class
  out_file.write(golem_row.join(',') + "\n")
  
end

golem_file.close
out_file.close

