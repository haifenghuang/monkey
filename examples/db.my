let dbOp = fn() {
	os.remove("./foo.db") //delete `foo.db` file
	let db = dbOpen("sqlite3", "./foo.db")
	if (db == nil) {
		println("DB open failed, error:", db.message())
		return false
	}
	defer db.close()
	let sqlStmt = `create table foo (id integer not null primary key, name text);delete from foo;`
	let exec_ret = db.exec(sqlStmt)
	if (exec_ret == nil) {
		println("DB exec failed! error:", exec_ret.message())
		return false
	}

	let tx = db.begin()
	if (tx == nil) { 
		println("db.Begin failed!, error:", tx.message())
		return false
	}

	let stmt = tx.prepare(`insert into foo(id, name) values(?, ?)`)
	if (stmt == nil) {
		println("tx.Prepare failed!, error:", stmt.message())
		return false
	}

	defer stmt.close()
	let i = 0
	for (i = 0; i < 105; i++) {
		let name = "hello" + i
		if (i>100) {
			//insert `null` value. There are six predefined values:
			let rs = stmt.exec(i, sql.STRING_NULL)
		} else {
			let rs = stmt.exec(i, name)
		}
		
		if (rs == nil) {
			println("statement exec failed, error:", rs.message())
			return false
		}
	} //end for

	tx.commit()

	let id, name = 0, ""
	let rows = db.query("select id, name from foo")
	if (rows == nil) {
		println("db queue failed, error:", rows.message())
		return false
	}
	defer rows.close()
	while (rows.next()) {
		rows.scan(id, name)
		if (name.valid()) { //check if it's `null`
			println(id, "|", name)
		} else {
			println(id, "|", "null")
		}
	}
	return true
}

let ret = dbOp()
if (ret == nil) {
	os.exit(1)
}

os.exit()
