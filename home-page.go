package main

import (
	"net/http"
)

func handleHomePage(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	const page = `<html>
<head></head>
<body>
	<p> Hello from (ideally) Docker! I'm a Go server. </p>
	<p> Hello from Docker! I'm a Go server. </p>
	<p> Hi Docker, I pushed a new version </p>
</body>
</html>
`
	w.Write([]byte(page))
}
