class Array
  def sum
    reduce(:+)
  end

  def mean
    sum.to_f / size
  end

  def var
    m = mean
    reduce(0) { |a,b| a + (b - m) ** 2 } / size
  end

  def sd
    Math.sqrt(var)
  end
end


n = 0
while line = gets
  n += 1
  cnt = {}
  puyos = line.chomp.chars
  colors = puyos.group_by(&:itself).keys
  colors.each do |c|
    cnt[c] = 0
  end
  distance = []
  puyos.each_slice(2) do |pair|
    cnt[pair[0]] += 1
    cnt[pair[1]] += 1
    max = cnt.values.max

    # 平均
    #sum = cnt.to_a.inject(0.0) {|sum, v| sum += max - v[1]}
    #distance.push(sum.to_f / colors.length)
    # 一番離れた距離
    distance.push(cnt.to_a.map {|k, v| (max - v).to_f}.max)
    # 離れた距離の合計
    # sum = cnt.to_a.inject(0.0) {|sum, v| sum += max - v[1]}
    # distance.push(sum)
    # 多い順に離れた距離の合計
    #sum = 0
    #cnts = cnt.values.sort.reverse
    #cnts.each_with_index do |v, i|
    #  next if i == 0
    #  sum += cnts[i - 1] - v
    #end
    #distance.push(sum)
    # 離れた 2 色の距離の合計
    #max2 = cnt.to_a.map {|k, v| (max - v).to_f}.max(2)
    #distance.push([max2.sum, max2[0] - max2[1]])
  end
  #distance.each do |v|
  #  p v
  #end
  # max2
  #puts "#{n} #{distance.map{|v| v[0]}.max}"
  # max
  puts "#{n} #{distance.max}"
  #平均
  #puts "#{n} #{distance.sum / distance.length}"
  # 中央値
  #puts "#{n} #{distance.length % 2 == 0 ? distance[(distance.length / 2) - 1] : distance[distance.length / 2]}"
  # 分散
  #puts distance.var
end
