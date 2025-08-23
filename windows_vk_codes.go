package hook

import "fmt"

// Windows Virtual Key Codes (VK codes) - what your system actually sends
// Source: https://learn.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes
var WindowsVKCodes = map[string]uint16{
	// Mouse buttons
	"mleft":   0x01, // VK_LBUTTON
	"mright":  0x02, // VK_RBUTTON
	"mcenter": 0x04, // VK_MBUTTON

	// Letter keys
	"a": 0x41, "b": 0x42, "c": 0x43, "d": 0x44, "e": 0x45,
	"f": 0x46, "g": 0x47, "h": 0x48, "i": 0x49, "j": 0x4A,
	"k": 0x4B, "l": 0x4C, "m": 0x4D, "n": 0x4E, "o": 0x4F,
	"p": 0x50, "q": 0x51, "r": 0x52, "s": 0x53, "t": 0x54,
	"u": 0x55, "v": 0x56, "w": 0x57, "x": 0x58, "y": 0x59,
	"z": 0x5A,

	// Number keys
	"0": 0x30, "1": 0x31, "2": 0x32, "3": 0x33, "4": 0x34,
	"5": 0x35, "6": 0x36, "7": 0x37, "8": 0x38, "9": 0x39,

	// Function keys
	"f1": 0x70, "f2": 0x71, "f3": 0x72, "f4": 0x73,
	"f5": 0x74, "f6": 0x75, "f7": 0x76, "f8": 0x77,
	"f9": 0x78, "f10": 0x79, "f11": 0x7A, "f12": 0x7B,

	// Navigation keys
	"insert":    0x2D, // VK_INSERT
	"delete":    0x2E, // VK_DELETE
	"home":      0x24, // VK_HOME
	"end":       0x23, // VK_END
	"page_up":   0x21, // VK_PRIOR
	"page_down": 0x22, // VK_NEXT
	"up":        0x26, // VK_UP
	"down":      0x28, // VK_DOWN
	"left":      0x25, // VK_LEFT
	"right":     0x27, // VK_RIGHT

	// Modifier keys
	"shift":     0x10, // VK_SHIFT
	"control":   0x11, // VK_CONTROL
	"ctrl":      0x11, // VK_CONTROL
	"alt":       0x12, // VK_MENU
	"tab":       0x09, // VK_TAB
	"enter":     0x0D, // VK_RETURN
	"escape":    0x1B, // VK_ESCAP
	"esc":       0x1B, // VK_ESCAP
	"backspace": 0x08, // VK_BACK
	"space":     0x20, // VK_SPACE

	// Special keys
	"caps_lock":    0x14, // VK_CAPITAL
	"num_lock":     0x90, // VK_NUMLOCK
	"scroll_lock":  0x91, // VK_SCROLL
	"print_screen": 0x2C, // VK_SNAPSHOT
	"pause":        0x13, // VK_PAUSE

	// Punctuation
	"minus":         0xBD, // VK_OEM_MINUS
	"equals":        0xBB, // VK_OEM_PLUS
	"comma":         0xBC, // VK_OEM_COMMA
	"period":        0xBE, // VK_OEM_PERIOD
	"slash":         0xBF, // VK_OEM_2
	"backslash":     0xDC, // VK_OEM_5
	"semicolon":     0xBA, // VK_OEM_1
	"quote":         0xDE, // VK_OEM_7
	"grave":         0xC0, // VK_OEM_3
	"open_bracket":  0xDB, // VK_OEM_4
	"close_bracket": 0xDD, // VK_OEM_6

	// Keypad keys
	"kp_0": 0x60, "kp_1": 0x61, "kp_2": 0x62, "kp_3": 0x63,
	"kp_4": 0x64, "kp_5": 0x65, "kp_6": 0x66, "kp_7": 0x67,
	"kp_8": 0x68, "kp_9": 0x69,
	"kp_multiply": 0x6A, "kp_add": 0x6B, "kp_separator": 0x6C,
	"kp_subtract": 0x6D, "kp_decimal": 0x6E, "kp_divide": 0x6F,
	"kp_enter": 0xE0, // Extended key

	// Extended keys (with 0xE0 prefix)
	"left_control":  0xA2, // VK_LCONTROL
	"right_control": 0xA3, // VK_RCONTROL
	"left_shift":    0xA0, // VK_LSHIFT
	"right_shift":   0xA1, // VK_RSHIFT
	"left_alt":      0xA4, // VK_LMENU
	"right_alt":     0xA5, // VK_RMENU
	"left_gui":      0x5B, // VK_LWIN
	"right_gui":     0x5C, // VK_RWIN
}

// Reverse mapping from VK code to key name
var WindowsVKCodeToName = make(map[uint16]string)

func init() {
	// Build reverse mapping
	for name, code := range WindowsVKCodes {
		WindowsVKCodeToName[code] = name
	}
}

// GetKeyName returns the key name for a given VK code
func GetWindowsVKKeyName(code uint16) string {
	if name, ok := WindowsVKCodeToName[code]; ok {
		return name
	}
	return fmt.Sprintf("unknown_%d", code)
}

// GetKeyCode returns the VK code for a given key name
func GetWindowsVKKeyCode(name string) (uint16, bool) {
	code, ok := WindowsVKCodes[name]
	return code, ok
}
