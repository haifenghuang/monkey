// Here we use 'using' statement, so we do not need to call infile.close().
// When finished running 'using' statement, it will automatically call infile.close().
using (infile = newFile("./file.demo", "r")) {
	if (infile == nil) {
		println("opening 'file.demo' for reading failed, error:", infile.message())
		os.exit(1)
	}

	let line;
	let num = 0
	//Read file by using extraction operator(">>")
	while (infile>>line != nil) {
		num++
		printf("%d	%s\n", num, line)
	}
}
