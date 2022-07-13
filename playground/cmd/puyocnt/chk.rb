n = 0
while line = gets
  n += 1
  puyos = line.chomp.chars
  cnt = puyos.group_by(&:itself).map {|v| v[1].length}
  if cnt.max == 68 && cnt.min == 60
    puts n
  end
end
