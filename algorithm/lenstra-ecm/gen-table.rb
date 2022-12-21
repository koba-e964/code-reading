include Math

def ecm_complexity(d)
    lnp = log(10) * d
    exp(sqrt(2) * sqrt(lnp * log(lnp)))
end

def format(number)
    lognumber = log10(number).floor
    mantissa = number * 10 ** -lognumber
    sprintf("%.2f * 10^%d", mantissa, lognumber)
end

puts "| p | L_p(1/2, sqrt(2)) |"
puts "| - | - |"
for i in 1..14
    log10p = i * 5
    puts "| 10^#{log10p} | #{format(ecm_complexity(log10p))} |"
end
