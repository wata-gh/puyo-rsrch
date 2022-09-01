clusters = {}
while line = gets
  c, p = line.chomp.split(' a')
  p = 'a' + p
  if clusters[p]
    clusters[p].push(c)
  else
    clusters[p] = [c]
  end
end
clusters.each do |k, v|
  if v.length > 1
    puts "#{k} #{v}"
  end
end
