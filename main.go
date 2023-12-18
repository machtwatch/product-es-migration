package main

import (
	"product-es-migration/cmd"
	"product-es-migration/config"
)

func main() {
	config.SetConfig(".", ".env")
	cmd.Execute()
}
