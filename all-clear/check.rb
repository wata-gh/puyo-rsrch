def check(i, puyos)
  g = puyos.group_by(&:itself)
  if g.values.all? { |v| v.length >= 4 }
    puts "#{i} #{puyos.length / 2} #{puyos.join}"
    return true
  end
  false
end

i = 1
while line = gets
  (3..255).select(&:odd?).each do |n|
    if check(i, line[..n].chars)
      break
    end
  end
  i += 1
end
