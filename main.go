package main

import "ak/services"

func main() {
	// _ = cmd.Execute()
	items := services.FetchStep1()
	services.FetchStep2(items)
}
