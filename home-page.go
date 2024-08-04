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
</body>
</html>
`
	w.Write([]byte(page))
}
