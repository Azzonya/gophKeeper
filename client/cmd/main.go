package main

import "gophKeeper/client/internal/app"

func main() {
	a := &app.App{}

	a.Init()
	a.Start()
	a.Listen()
	a.Exit()
}
