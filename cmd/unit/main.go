package main

import "github.com/devdammit/shekel/cmd/unit/app"

func main() {
	wg := app.Run()

	wg.Wait()
}
