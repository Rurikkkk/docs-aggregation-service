package main

import "docs-aggregation-service/internal/app"

func main() {
	context := app.NewContext()
	context.HTTPServer().ListenAndServe()
}
