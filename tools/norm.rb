while line = gets
  fu = line.chomp
  table = {}
  colors = %w/b c d e/
  ns = ''
  puyos = fu.chars
  puyos[10*6] = 'b'
  puyos[11*6] = 'b'
  puyos[11*6+1] = 'b'
  puyos.each do |c|
    if c == 'a'
      ns += 'a'
    else
      v = table[c]
      if !v
        v = colors.shift
        table[c] = v
      end
      ns += v
    end
  end
  puts ns
end
