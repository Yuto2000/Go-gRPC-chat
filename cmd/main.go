package main

import (
	"chat/client"
	"os"
)

func main() {
	os.Exit(client.NewChat().Run())
}
