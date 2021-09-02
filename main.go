package main

// Very slim main.go file.
// We need 2 main.go files. One for lambda and one for dev.

func main() {
	// This function is defined 2 times.
	// Once in startProd and once in startDev.
	// This way we don't have to change the source code between
	// lambda and dev deploys.
	startup()
}
