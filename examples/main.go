package main

import (
	"fmt"
	"strconv"
	"strings"

	hook "github.com/ms-dosx86/gohook"
	flag "github.com/spf13/pflag"
)

func main() {
	hook.SetLogLevel(hook.Silent)
	binds := []string{}
	kinds := []string{}
	flag.StringArrayVar(&binds, "bind", []string{}, "Usage: --bind <key>")
	flag.StringArrayVar(&kinds, "kind", []string{}, "Usage: --kind <type>")
	flag.Parse()

	if len(binds) != len(kinds) {
		panic("bind and kind must have the same length")
	}

	kindsInt := make([]uint8, len(kinds))
	for i, bindType := range kinds {
		n, err := strconv.Atoi(bindType)
		if err != nil {
			panic(err)
		}
		kindsInt[i] = uint8(n)
	}

	for i, bind := range binds {
		if isMouseBind(bind) {
			registerMouseBind(bind, kindsInt[i])
		} else {
			registerKeyBind(bind, kindsInt[i])
		}
	}

	fmt.Println("Starting...")
	s := hook.Start()
	defer hook.End()
	<-hook.Process(s)
}

func isMouseBind(bind string) bool {
	return bind == "mleft" || bind == "mright" || strings.HasPrefix(bind, "wheel")
}

func registerMouseBind(bind string, kind uint8) {
	hookKind := hook.Kind(kind)

	switch hookKind {
	case hook.MouseDown:
		if err := hook.Register(hook.MouseDown, []string{bind}, func(e hook.Event) {
			fmt.Printf("%d:%s\n", e.Kind, bind)
		}); err != nil {
			panic(err)
		}
	case hook.MouseHold:
		if err := hook.Register(hook.MouseHold, []string{bind}, func(e hook.Event) {
			fmt.Printf("%d:%s\n", e.Kind, bind)
		}); err != nil {
			panic(err)
		}
	default:
		panic(fmt.Sprintf("invalid bindType: %d", kind))
	}
}

func registerKeyBind(bind string, kind uint8) {
	parts := strings.Split(bind, ",")

	hookKind := hook.Kind(kind)
	switch hookKind {
	case hook.KeyDown:
		if err := hook.Register(hook.KeyDown, parts, func(e hook.Event) {
			fmt.Printf("%d:%s\n", e.Kind, bind)
		}); err != nil {
			panic(err)
		}
	case hook.KeyUp:
		if err := hook.Register(hook.KeyUp, parts, func(e hook.Event) {
			fmt.Printf("%d:%s\n", e.Kind, bind)
		}); err != nil {
			panic(err)
		}
	default:
		panic(fmt.Sprintf("invalid bindType: %d", kind))
	}
}
