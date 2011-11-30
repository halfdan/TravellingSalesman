pm=[];
pc=[];

# PC
i = 0.1
while i <= 0.9
    i = i.round(3)
    pc << i
    i += 0.05
end

# PM
i = 0.005
while i <= 0.2
    i = i.round(4)
    pm << i
    i += 0.005
end

#puts pc
#puts pm

genecnt = 100
max_gen = 200

results = {}

pc.each do |c|
    results[c] = {}
    threads = []
    pm.each do |m|
        threads << Thread.new do
            results[c][m] = `./TravellingSalesman -pc=#{c} -pm=#{m} -ngenes #{genecnt} --maxgener #{max_gen} -cores=2 -runs=10 -file=36cities-border.txt`
        end
    end
    # Wait for all processes to complete
    threads.each do |t|
        t.join
    end
end


results.each do |pc, value|
    value.to_a.sort.each do |pm, data|
        puts "#{data}"
    end
    puts
end

