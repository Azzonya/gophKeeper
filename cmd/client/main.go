package main

import "gophKeeper/internal/client/app"

func main() {
	a := &app.App{}

	a.Init()
	a.PreStartHook()
	a.Start()
	a.Listen()
	a.Exit()
}
