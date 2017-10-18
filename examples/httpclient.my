fn() {
	let response = http.get("http://localhost:8080/hello")
	defer response.closeBody()
	let body = response.readAll()
	println(body)
}()