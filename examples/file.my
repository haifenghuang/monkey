filename = "./file.log"

fileread = fn() {
	file = newFile(filename, "r")
	if (file == nil) {
		println("opening ", filename, "for reading failed, error:", file.message())
		return false
	}
	defer file.close()

	while ((l = file.readLine()) != nil) {
		println(l)
	}

	//Position to beginning of file
	file.seek(os.SEEK_SET, 0)
	do {
		x = file.readRune()
		if (x != nil) {
			println(x)
		} else {
			break
		}
	}
	return true
}


filewrite = fn() {
	file = newFile(filename, "w+")
	if (file == nil) {
		println("opening ", filename, "for writing failed, error:", file.message())
		return false
	}
	defer file.close()

	file.writeLine("This is the first line")
	file.writeLine("This is the second line")
	file.writeString("这是第三行\n")
	//Using `fmt.fprintf` function to write to file
	fmt.fprintf(file, "This is the %s line\n", "fourth")

	return true
}


ret = filewrite()
if (ret == false) {
	os.exit(1)
}

time.sleep(2)
ret = fileread()
if (ret == false) {
	os.exit(1)
}

os.exit(0)