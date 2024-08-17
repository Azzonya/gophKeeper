package main

import "gophKeeper/server/internal/app"

func main() {
	a := &app.App{}

	a.Init()
	a.Start()
	a.Listen()
	a.Stop()
	a.Exit()
}
