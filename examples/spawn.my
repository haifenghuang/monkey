# Test "sync"'s Wait Object(code from go's source :src/sync/example_test.go)
let ExampleWaitGroup = fn() {
	let wg = newWaitGroup()
	let urls = ["http://www.baidu.org/","http://www.csdn.net/","http://www.sina.com.cn/"]
	for url in urls {
		// Increment the WaitGroup counter.
		wg.add(1)
		// Launch a goroutine to fetch the URL.
		spawn fn(url) {
			println('In spawn, url={url}')
			// Decrement the counter when the goroutine completes.
			defer fn(wg) {
				wg.done()
			}(wg)
			// Fetch the URL.
			http.get(url)
		}(url)
	}
	// Wait for all HTTP fetches to complete.
	wg.wait()
}

ExampleWaitGroup()