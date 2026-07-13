package mobile

// Key names for MobileDriver.KeyPress.
//
// Deliberately tiny (mirrors platform-python AXI-1145): the earlier speculative
// HID constants (HOME, volume, media keys) are gone. Grow this entry by entry,
// in lockstep with the named-key table on the device side, as real needs appear.
const (
	// KeyEnter submits forms / fires the on-screen keyboard's Go / Search action.
	KeyEnter = "enter"
)
