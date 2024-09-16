package main

import "gophKeeper/internal/server/app"

func main() {
	a := &app.App{}

	a.Init()
	a.Start()
	a.Listen()
	a.Stop()
	a.Exit()
}
