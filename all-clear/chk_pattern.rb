patterns = []
while line = gets
  puyos = line.chomp.tr('[]', '').split(' ')
  a = puyos[0]
  b = ''
  c = ''
  d = ''
  normalize = ''
  n = puyos.map { |p|
    case p
    when a
      'a'
    when b
      'b'
    when c
      'c'
    when d
      'd'
    else
      if b == ''
        b = p
        'b'
      elsif c == ''
        c = p
        'c'
      else
        d = p
        'd'
      end
    end
  }.join
  patterns.push(n)
end
puts patterns.sort.uniq
