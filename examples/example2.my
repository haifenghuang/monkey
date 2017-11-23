# Test for function default parameter and variadic parameters
add = fn (x, y=5,z=7, args...) {
	w = x + y + z
	for i in args {
		w += i
	}
	return w
}

w = add(2,3,4,5,6,7)
println(w)

# Test for 'enum' type
LogOption = enum {
	Ldate         = 1 << 0,
	Ltime         = 1 << 1,
	Lmicroseconds = 1 << 2,
	Llongfile     = 1 << 3,
	Lshortfile    = 1 << 4,
	LUTC          = 1 << 5,
	LstdFlags     = 1 << 4 | 1 << 5
};
	

opt = LogOption.LstdFlags
println(opt)

println()
for s in LogOption.getNames() { //not ordered
	println(s)
}
println()

println()
for s in LogOption.getValues() { //not ordered
	println(s)
}
println()

println(LogOption.getName(LogOption.Lshortfile))
println()

// Test for 'qw'(like perl's 'qw')
for str in qw<abc, def, ghi, jkl, mno> { //allowed 'qw' pair is '{}', '<>', '()'
  println('str={str}')
}

newArr = qw(1,2,3.5)
fmt.printf("newArr=%v\n", newArr)

for ij in range(6,2) {
    println('ij = {ij}')
}

