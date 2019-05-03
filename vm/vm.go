// Package vm defines the VirtualMachine interface, as well as standard helper
// methods for configuring the standard WASM VM.
package vm

// VirtualMachine defines a Polaris Virtual Machine.
type VirtualMachine interface {
	GetVirtualMachineProtocol() *Protocol // Get virtual machine protocol.
}
