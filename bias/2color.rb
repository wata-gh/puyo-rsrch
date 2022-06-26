n = 0
all_result = []
while line = gets
  result = []
  puyos = line.chomp.chars
  colors = puyos.group_by(&:itself).keys
  last = []
  counter = 0
  puyos.each_with_index do |c, i|
    if last.include?(c)
      counter += 1
    else
      last.push(c)
      last = last.group_by(&:itself).keys.last(2)
      counter = 0
      (0..i).to_a.reverse.each do |n|
        if last.include?(puyos[n])
          counter += 1
        else
          break
        end
      end
    end
    result.push([i + 1, counter, last.dup])
  end
  all_result.push(result)
end

#all_result.each.with_index(1) do |result, n|
#  max = result.max_by {|v| v[1]}
#  puts "#{n} #{max[0]} #{max[1]} #{max[2].join}"
#end
r = all_result.group_by do |result|
  result.map {|v| v[1]}.max
end
r.each do |k, v|
  puts "#{k} #{v.length}"
end
