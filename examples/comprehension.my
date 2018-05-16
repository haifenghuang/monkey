//hash comprehension (from hash)
z1 = { v:k for k,v in {"key1":10, "key2":20, "key3":30}} //reverse key-value pair
println(z1)

//hash comprehension (from array)
z2 = {x:x**2 for x in [1,2,3]}
println(z2)

//hash comprehension (from .. range)
z3 = {x:x**2 for x in 5..7}
println(z3)

//hash comprehension (from string)
z4 = {x:x.upper() for x in "hi"}
println(z4)

//hash comprehension (from tuple)
z5 = {x+1:x+2 for x in (1,2,3)}
println(z5)
	
	
println("==================================")
//array
x = [[word.upper(), word.lower(), word.title()] for word in ["hello", "world", "good", "morning"]]
println(x)

//string
y = [ c.upper() for c in "huanghaifeng" where $_ % 2 != 0] //$_ is the index
println(y)

//range
w = [i + 1 for i in 1..10]
println(w)

//tuple
v = [x+1 for x in (12,34,56)]
println(v)

//hash
z = [v * 10 for k,v in {"key1":10, "key2":20, "key3":30}]
println(z)