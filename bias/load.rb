def load(dump_file)
  Marshal.load(File.binread(dump_file))
end

def sort(result)
  result.map { |r|
    r.reject{|c, v| c == 'no'}.map { |c, v|
      m = v.max_by {|i| i[1]}
      [r['no'], c, m]
    }
  }.map {|r|
    r.max_by{|v| v[2][1]}
  }.sort_by {|r|
    r[2][1]
  }
end
