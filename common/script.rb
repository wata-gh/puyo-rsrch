@r = Random.new(Time.now.to_i)
def delay
  i = @r.rand(10)
  if i < 7
    return 0.1
  elsif i < 9
    return 0.2
  end
  return 0.3
end
puts 'tell application "System Events"'
puts 'delay 3'
while line = gets
  line.chars.each do |c|
    if c != ' '
      puts "delay #{delay}"
    end
    if c == '"'
      c = '\"'
    end
    puts "keystroke \"#{c}\""
  end
end
puts 'end tell'
