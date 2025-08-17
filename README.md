# gohook

This is a fork of [gohook](https://github.com/robotn/gohook) that is updated to use Go 1.25.0 and has been rewritten to completely to fix the issues with mouse events.

## Requirements (Linux):

[Robotgo-requirements-event](https://github.com/go-vgo/robotgo#requirements)

## Install:

With Go module support (Go 1.11+), just import:

```go
import "github.com/ms-dosx86/gohook"
```

## Examples:

```Go
package main

import (
	"fmt"

	hook "github.com/ms-dosx86/gohook"
)

func main() {
	add()

	low()
}

func add() {
	fmt.Println("--- Please press ctrl + shift + q to stop hook ---")
	hook.Register(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("ctrl-shift-q")
		hook.End()
	})

	fmt.Println("--- Please press w---")
	hook.Register(hook.KeyDown, []string{"w"}, func(e hook.Event) {
		fmt.Println("keyDown: ", "w")
	})

	hook.Register(hook.KeyUp, []string{"w"}, func(e hook.Event) {
		fmt.Println("keyUp: ", "w")
	})

	s := hook.Start()
	<-hook.Process(s)
}

func low() {
	evChan := hook.Start()
	defer hook.End()

	for ev := range evChan {
		fmt.Println("hook: ", ev)
	}
}

```

Based on [libuiohook](https://github.com/kwhat/libuiohook).