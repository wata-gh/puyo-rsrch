while line = gets
  puyos = line.chomp.chars
  odd = []
  even = []
  puyos.each.with_index(1) do |puyo, i|
    if i % 2 == 0
      even << puyo
    else
      odd << puyo
    end
  end
  odds = odd.each_slice(16).to_a
  evens = even.each_slice(16).to_a

  odds.each_with_index do |odd, i|
    puts evens[i].join('.')
    puts odd.join('.')
    puts
  end
end
