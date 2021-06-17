#card_file = File.open('./cards.txt', 'r')
#out_file = File.open('./cards.csv', 'w')

card_file = File.open('starter_cards.txt', 'r')
out_file = File.open('starter_cards.csv', 'w')

COLORS = {y: :yellow, g: :green, b: :blue, p: :pink}
SYM_KEY = COLORS.merge("*": :upgrades)
KEYS = COLORS.values.map {|color| "#{color}_in"} + COLORS.values.map {|color| "#{color}_out"} + ["upgrades"]

def parse_card(card)
  puts card

  inputs, outputs = card.split('>')
  costs = inputs.nil? ? {} : gems_to_hash(inputs.scan(/[0-9][y,g,b,p]/), 'in')
  results = gems_to_hash(outputs.scan(/[0-9][y,g,b,p,\*]/), 'out')

  return costs.merge(results)
end


def gems_to_hash(cost_arr, suffix=nil)
  return cost_arr.each_with_object({}) do |cost_str, hash|
    color = "#{ SYM_KEY[cost_str[1].to_sym]}"
    color += "_#{suffix}" if (suffix && COLORS.keys.include?(cost_str[1].to_sym))

    hash[color] = cost_str[0].to_i
  end
end




out_file.write("#{KEYS.join(',')}\n")
card_file.each_line do |card|
  card_hash = parse_card(card)
  puts card_hash

  card_row = KEYS.each_with_object([]) { |key, arr| arr << (card_hash[key] || 0).to_s}
  out_file.write(card_row.join(',') + "\n")
end

card_file.close
out_file.close

