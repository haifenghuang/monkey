let handler = fn(writer, request) {
	writer.write("Hello HTTP client, I'm server")
}

http.listenAndServe("127.0.0.1:8080", handler)
//Same as above
//http.listenAndServe("127.0.0.1:8080", fn(writer, request) {
//	writer.write("Hello HTTP client, I'm server")
//})