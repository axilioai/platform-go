package mobile

// Key names for MobileDriver.KeyPress.
//
// Deliberately tiny (mirrors platform-python AXI-1145): the earlier speculative
// constants (HOME, volume, media keys) are gone. Grow this entry by entry,
// in lockstep with the named-key table on the device side, as real needs appear.
// The contract enumerates the deployed names (KeyboardKeyPressParams.key in
// contracts/dcp-asyncapi.json); TestNamedKeysMatchContract fails when this
// list and that enum disagree, so a key the phone learns cannot stay a raw
// string here.
const (
	// KeyEnter submits forms / fires the on-screen keyboard's Go / Search action.
	KeyEnter = "enter"
	// KeyCapsLock toggles the on-screen keyboard's caps lock.
	KeyCapsLock = "capslock"
)
