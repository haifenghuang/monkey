let todos = []
let showToDo = fn(writer, request) {
	writer.write("<html>")
	writer.write("<head><title>ToDo List</title></head>")
	
	writer.write("<body>")
	writer.write("<h1>ToDo List</h1>")
	
	writer.write("<ul>\n")
	for todo in todos {
	    writer.write("<li>" + todo + "</li>\n")
	}
	writer.write("</ul>")
	
	writer.write("<h2>Add ToDo</h2>")
	writer.write(`<form method="post" action="/todos/new">`)
	writer.write(`<input type="text" name="todo">`)
	writer.write(`<input type="submit" name="Add">`)
	writer.write(`</form>`)
	
	writer.write("</body>")
	writer.write("</html>")
}

let addToDo = fn(writer, request) {
	let receiveValue = request.formValue("todo")
	todos = receiveValue + todos
	http.redirect(writer, request, "/todos", 303)
}

# In browser, navigate to http://localhost:9090/todos page
http.handleFunc("/todos", showToDo)
http.handleFunc("/todos/new", addToDo)
let ret = http.listenAndServe("127.0.0.1:9090")
if (ret == false) {
	println("listenAndServe failed, error:", ret.message())
}

