//test csv reader
let r = newCsvReader("./examples/test.csv")
if r == nil {
    printf("newCsv returns err, message:%s\n", r.message())
}

r.setOptions({"Comma":";", "Comment":"#"})

ra = r.readAll()
if (ra == nil) {
    printf("readAll returns err, message:%s\n", ra.message())
}

for line in ra {
	println(line)
	for record in line {
		println("	", record)
	}
}

r.close() //do not forget to close the reader

//test csv writer
let ofile = newFile("./examples/demo.csv", "a+")
//note:we do not test if the file is successfully opened or not, but in real code, you should do it

let w = newCsvWriter(ofile)
w.setOptions({"Comma":"	"})
w.write(["1", "2", "3"])
w.writeAll([["4", "5", "6"],["7", "8", "9"],["10", "11", "12"]])
w.flush()

ofile.close() //do not forget to clsoe the file

