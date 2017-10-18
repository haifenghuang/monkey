################################################################
# Same as httpserver2.my file, but use fmt.fprintf function
################################################################
let todos = []
let showToDo = fn(writer, request) {
	fmt.fprintf(writer, "<html><head><title>ToDo List</title></head><body><h1>ToDo List</h1><ul>\n")
	for todo in todos {
	    fmt.fprintf(writer, "<li> %s </li>\n", todo)
	}
	fmt.fprintf(writer, "</ul>")
	
	fmt.fprintf(writer,"<h2>Add ToDo</h2>")
	fmt.fprintf(writer,`<form method="post" action="/todos/new">`)
	fmt.fprintf(writer,`<input type="text" name="todo">`)
	fmt.fprintf(writer,`<input type="submit" name="Add">`)
	fmt.fprintf(writer,`</form>`)
	
	fmt.fprintf(writer,"</body>")
	fmt.fprintf(writer,"</html>")
}

let addToDo = fn(writer, request) {
	let receiveValue = request.formValue("todo")
	todos = receiveValue + todos
	http.redirect(writer, request, "/todos", 303)
}

http.handleFunc("/todos", showToDo)
http.handleFunc("/todos/new", addToDo)
let ret = http.listenAndServe("127.0.0.1:9090")
if (ret == false) {
	println("listenAndServe failed, error:", ret.message())
}

