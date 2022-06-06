def all_clearable?(i, puyos)
  puyos.group_by(&:itself).values.all? { |v| v.length >= 4 }
end

i = 1
while line = gets
  (3..255).select(&:odd?).each do |n|
    puyos = line[..n].chars
    if all_clearable?(i, puyos)
      puts "#{i} #{puyos.length / 2} #{puyos.join}"
      break
    end
  end
  i += 1
end
