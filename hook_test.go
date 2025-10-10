package hook

import (
	"testing"
	"time"
)

const (
	TIMEOUT = 2 * time.Second
)

func TestKeyDown(t *testing.T) {
	done := make(chan bool)
	err := Register(KeyDown, []string{"a"}, func(e Event) {
		done <- true
	})
	if err != nil {
		t.Fatal(err)
	}
	ch := Start()
	defer End()
	Process(ch)

	go func() {
		ch <- Event{
			Rawcode: Keycode["a"],
			Kind:    KeyDown,
		}
	}()

	select {
	case <-time.After(TIMEOUT):
		t.Fatal("Timeout waiting for keydown")
	case <-done:
		t.Log("KeyDown received")
	}
}

func TestKeyDownWithModifier(t *testing.T) {
	done := make(chan bool)
	err := Register(KeyDown, []string{"ctrl", "a"}, func(e Event) {
		done <- true
	})
	if err != nil {
		t.Fatal(err)
	}
	ch := Start()
	defer End()
	Process(ch)

	go func() {
		ch <- Event{
			Rawcode: Keycode["ctrl"],
			Kind:    KeyDown,
		}
		time.Sleep(100 * time.Millisecond)
		ch <- Event{
			Rawcode: Keycode["a"],
			Kind:    KeyDown,
		}
	}()

	select {
	case <-time.After(TIMEOUT):
		t.Fatal("Timeout waiting for keydown")
	case <-done:
		t.Log("KeyDown received")
	}
}

func TestKeyUp(t *testing.T) {
	done := make(chan bool)
	err := Register(KeyUp, []string{"delete"}, func(e Event) {
		done <- true
	})
	if err != nil {
		t.Fatal(err)
	}
	ch := Start()
	defer End()
	Process(ch)

	go func() {
		ch <- Event{
			Rawcode: Keycode["delete"],
			Kind:    KeyUp,
		}
	}()

	select {
	case <-time.After(TIMEOUT):
		t.Fatal("Timeout waiting for keyup")
	case <-done:
		t.Log("KeyUp received")
	}
}

func TestKeyUpWithModifier(t *testing.T) {
	done := make(chan bool)
	err := Register(KeyUp, []string{"ctrl", "a"}, func(e Event) {
		done <- true
	})
	if err != nil {
		t.Fatal(err)
	}
	ch := Start()
	defer End()
	Process(ch)

	go func() {
		ch <- Event{
			Rawcode: Keycode["ctrl"],
			Kind:    KeyUp,
		}
		ch <- Event{
			Rawcode: Keycode["a"],
			Kind:    KeyUp,
		}
	}()

	select {
	case <-time.After(TIMEOUT):
		t.Fatal("Timeout waiting for keyup")
	case <-done:
		t.Log("KeyUp received")
	}
}

func TestMouseDown(t *testing.T) {
	done := make(chan bool)
	err := Register(MouseDown, []string{"mleft"}, func(e Event) {
		done <- true
	})
	if err != nil {
		t.Fatal(err)
	}
	ch := Start()
	defer End()
	Process(ch)

	go func() {
		ch <- Event{
			Button: MouseMap["left"],
			Kind:   MouseDown,
		}
	}()

	select {
	case <-time.After(TIMEOUT):
		t.Fatal("Timeout waiting for mousedown")
	case <-done:
		t.Log("MouseDown received")
	}
}

func TestMouseUp(t *testing.T) {
	done := make(chan bool)
	err := Register(MouseUp, []string{"mright"}, func(e Event) {
		done <- true
	})
	if err != nil {
		t.Fatal(err)
	}
	ch := Start()
	defer End()
	Process(ch)

	go func() {
		ch <- Event{
			Button: MouseMap["right"],
			Kind:   MouseUp,
		}
	}()

	select {
	case <-time.After(TIMEOUT):
		t.Fatal("Timeout waiting for mouseup")
	case <-done:
		t.Log("MouseUp received")
	}
}

func TestCombinationsLimit(t *testing.T) {
	// Should fail if more than 4 keys are provided
	err := Register(KeyDown, []string{"ctrl", "a", "b", "c", "d"}, func(e Event) {})
	if err != nil {
		t.Log("That was expected")
	} else {
		t.Fatal("Expected error, got none")
	}

	// Should succeed if less than 4 keys are provided
	err = Register(KeyDown, []string{"ctrl", "a", "b", "c"}, func(e Event) {})
	if err != nil {
		t.Fatal("Expected no error, got", err)
	}
}

func TestMouseDownWithMouseUp(t *testing.T) {
	mouseDownOccurred := false
	mouseUpOccurred := make(chan bool)
	err := Register(MouseDown, []string{"mleft"}, func(e Event) {
		mouseDownOccurred = true
	})
	if err != nil {
		t.Fatal("Could not register mouse down callback: ", err)
	}
	err = Register(MouseUp, []string{"mleft"}, func(e Event) {
		mouseUpOccurred <- true
	})
	if err != nil {
		t.Fatal("Could not register mouse up callback: ", err)
	}
	ch := Start()

	defer End()

	Process(ch)

	go func() {
		ch <- Event{
			Button: MouseMap["left"],
			Kind:   MouseDown,
		}
		ch <- Event{
			Button: MouseMap["right"],
			Kind:   MouseUp,
		}
	}()

	select {
	case <-time.After(time.Second * 1):
		if !mouseDownOccurred {
			t.Fatal("Timeout waiting for mouse events")
		}
	case <-mouseUpOccurred:
		t.Fatal("MouseUp occurred for the wrong button")
	}
}

func TestPreventKeyDownSpamming(t *testing.T) {
	count := 0
	err := Register(KeyDown, []string{"a"}, func(e Event) {
		count++
	})
	if err != nil {
		t.Fatal(err)
	}
	ch := Start()
	defer End()
	Process(ch)

	go func() {
		// counts
		ch <- Event{
			Rawcode: Keycode["a"],
			Kind:    KeyDown,
		}
		// doesn't count
		ch <- Event{
			Rawcode: Keycode["a"],
			Kind:    KeyDown,
		}
		// doesn't count
		ch <- Event{
			Rawcode: Keycode["a"],
			Kind:    KeyDown,
		}
		// triggering key up for a different key
		// should clear the buffer
		ch <- Event{
			Rawcode: Keycode["b"],
			Kind:    KeyUp,
		}
		// counts
		ch <- Event{
			Rawcode: Keycode["a"],
			Kind:    KeyDown,
		}
		// triggering a mouse event should not clear the buffer
		ch <- Event{
			Button: MouseMap["left"],
			Kind:   MouseDown,
		}
		// doesn't count
		ch <- Event{
			Rawcode: Keycode["a"],
			Kind:    KeyDown,
		}
		// triggering key up for the same key
		// should clear the buffer
		ch <- Event{
			Rawcode: Keycode["a"],
			Kind:    KeyUp,
		}
		// counts
		ch <- Event{
			Rawcode: Keycode["a"],
			Kind:    KeyDown,
		}

		// Total count should be 3
	}()

	select {
	case <-time.After(TIMEOUT):
		t.Fatal("Timeout waiting for keydown")
	case <-time.After(100 * time.Millisecond):
		if count != 3 {
			t.Fatal("Expected 3 keydowns, got", count)
		}
		t.Log("KeyDown received")
	}
}

func TestKeyUpDoesNotFireForDifferentKey(t *testing.T) {
	done := make(chan bool)
	err := Register(KeyUp, []string{"a"}, func(e Event) {
		done <- true
	})
	if err != nil {
		t.Fatal(err)
	}
	ch := Start()
	defer End()
	Process(ch)

	go func() {
		// press and hold 'a'
		ch <- Event{Rawcode: Keycode["a"], Kind: KeyDown}
		// release a different key 'g' (should not trigger 'a' keyup)
		ch <- Event{Rawcode: Keycode["g"], Kind: KeyUp}
		// now release 'a' to ensure the correct event still fires
		time.Sleep(50 * time.Millisecond)
		ch <- Event{Rawcode: Keycode["a"], Kind: KeyUp}
	}()

	select {
	case <-time.After(TIMEOUT):
		t.Fatal("Timeout waiting for correct keyup sequence")
	case <-done:
		t.Log("KeyUp for 'a' received only after releasing 'a'")
	}
}

func TestComboDoesNotRetriggerWithUnrelatedKeys(t *testing.T) {
	count := 0
	err := Register(KeyDown, []string{"ctrl", "a"}, func(e Event) {
		count++
	})
	if err != nil {
		t.Fatal(err)
	}
	ch := Start()
	defer End()
	Process(ch)

	go func() {
		// hold the combo ctrl+a
		ch <- Event{Rawcode: Keycode["ctrl"], Kind: KeyDown}
		time.Sleep(10 * time.Millisecond)
		ch <- Event{Rawcode: Keycode["a"], Kind: KeyDown}
		// press unrelated keys while holding the combo
		ch <- Event{Rawcode: Keycode["g"], Kind: KeyDown}
		ch <- Event{Rawcode: Keycode["g"], Kind: KeyUp}
		// release combo
		ch <- Event{Rawcode: Keycode["a"], Kind: KeyUp}
		ch <- Event{Rawcode: Keycode["ctrl"], Kind: KeyUp}
	}()

	select {
	case <-time.After(TIMEOUT):
		t.Fatal("Timeout waiting for combo test")
	case <-time.After(100 * time.Millisecond):
		if count != 1 {
			t.Fatalf("Expected combo to trigger once, got %d", count)
		}
	}
}

func TestAKeyHandlersNotTriggeredByG(t *testing.T) {
	countDown := 0
	countUp := 0

	err := Register(KeyDown, []string{"a"}, func(e Event) {
		countDown++
	})
	if err != nil {
		t.Fatal(err)
	}

	err = Register(KeyUp, []string{"a"}, func(e Event) {
		countUp++
	})
	if err != nil {
		t.Fatal(err)
	}

	ch := Start()
	defer End()
	Process(ch)

	go func() {
		ch <- Event{Rawcode: Keycode["g"], Kind: KeyDown}
		ch <- Event{Rawcode: Keycode["g"], Kind: KeyUp}
	}()

	select {
	case <-time.After(TIMEOUT):
		t.Fatal("Timeout waiting for non-trigger condition")
	case <-time.After(100 * time.Millisecond):
		if countDown != 0 || countUp != 0 {
			t.Fatalf("Expected no 'a' callbacks, got KeyDown=%d KeyUp=%d", countDown, countUp)
		}
	}
}
