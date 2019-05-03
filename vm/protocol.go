// Package vm defines the VirtualMachine interface, as well as standard helper
// methods for configuring the standard WASM VM.
package vm

// Protocol defines a virtual machine protocol.
type Protocol struct {
	LanguagesSupported []string `json:"languages"` // Languages supported by the VirtualMachine (e.g. wasm).
}
