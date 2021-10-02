package main

import (
	"fmt"

	"github.com/logica0419/Q-n-A/router"
)

func main() {
	fmt.Print("Q'n'A - traP Anonymous Question Box Service")

	e := router.Setup()

	e.Logger.Panic(e.Start(":9000"))
}
