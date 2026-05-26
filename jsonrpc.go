// Package jsonrpc provides primitives for working with JSON-RPC 2.0.
//
// It defines the core Request and Response structures, along with validation
// logic to ensure compliance with the JSON-RPC 2.0 specification.
//
// Key features:
//   - Strict adherence to JSON-RPC 2.0 specification.
//   - Support for flexible ID types (string, number, and null).
//   - Native support for Notifications (requests without an ID).
//   - Built-in validation for requests, responses, and parameters.
//   - Idiomatic Go API with automatic JSON marshalling for results and errors.
//
// More information about JSON-RPC 2.0 can be found at:
// https://www.jsonrpc.org/specification
package jsonrpc

const (
	// JSONRPCVersion is the supported version of the JSON-RPC protocol.
	JSONRPCVersion string = "2.0"
)
