# HTTP Router

Just another simple http router with middleware and error handling support.

## Examples

````
func main() {
	router := httprouter.NewRouter()

	helloHandler := func(responseWriter http.ResponseWriter, request *http.Request) error {
		_, _ = responseWriter.Write([]byte("Hello World!"))

		return nil
	}

	router.Get("/hello", httprouter.HandlerFunc(helloHandler))

	_ = http.ListenAndServe(":9015", router)
}
````