while line = gets
  fu = line.chomp
  puyos = fu.tr('a', '')
  if puyos.length == 20
    g = puyos.chars.group_by(&:itself).values.map(&:length).sort.reverse
    puts fu
  end
end

