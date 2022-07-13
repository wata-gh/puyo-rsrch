def load
  all_result = []
  no = 0
  while line = STDIN.gets
    cnt = {}
    result = {}
    no += 1
    puyos = line.chomp.chars
    colors = puyos.group_by(&:itself).keys
    colors.each do |c|
      cnt[c] = 0
      result[c] = []
    end

    puyos.each.with_index(1) do |c, i|
      colors.each do |l|
        if c == l
          result[c] << [i, cnt[c]]
          cnt[c] = 0
          next
        end
        cnt[l] = cnt[l] + 1
      end
    end

    cnt.each do |k, v|
      if v != 0
        result[k].push([puyos.length, v])
      end
    end

    result['no'] = no
    all_result.push(result)
  end

  all_result
end

puts ARGV[0]
File.open(ARGV[0], 'wb') do |f|
  f.write(Marshal.dump(load()))
end
