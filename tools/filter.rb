while line = gets
  fu = line.chomp
  puyos = fu.tr('a', '')
  if puyos.length == 17
    g = puyos.chars.group_by(&:itself).values.map(&:length).sort.reverse
    if g[0] == 5 && g[1] == 4
      puts fu
    end
  end
end

