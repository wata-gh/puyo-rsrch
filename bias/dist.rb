no = 0
while puyos = gets&.chomp&.chars
  no += 1
  cnt = Hash.new(0)
  results = Hash.new {|h, k| h[k] = []}
  colors = puyos.group_by(&:itself).keys

  puyos.each.with_index(1) do |puyo, i|
    colors.each do |c|
      if puyo == c
        results[puyo] << [i, cnt[puyo]]
        cnt[puyo] = 0
        next
      end
      cnt[c] += 1
    end
  end

  cnt.reject{|c, v| v == 0}.each do |c, v|
    results[c] << [puyos.length, v]
  end

  color_max = results.map{|c, r| [c, r.max_by{|r| r[1]}]}.sort_by{|c, r| r[1]}.reverse
  puts "#{no} #{color_max.map{|c, r| [c, r.join(' ')].join(' ')}.join(' ')}"
  #puts "#{no} #{color_max.max_by{|v| v[1][1]}}"
end
