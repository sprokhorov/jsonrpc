# jsonrpc

A lightweight, zero-dependency Go library providing primitives for working with JSON-RPC 2.0.

This package implements the core data structures and validation logic defined in the [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification), focusing on ease of use and strict compliance.

## Features

- **Strictly Spec-Compliant**: Follows JSON-RPC 2.0 rules for requests, responses, and notifications.
- **Flexible ID Handling**: Supports string, number, and null IDs out of the box using `json.RawMessage`.
- **Easy Notifications**: Native support for notifications (requests without an ID).
- **Automated Validation**: Built-in methods to validate protocol version, method presence, and parameter structure.
- **Developer-Friendly API**: High-level helpers like `NewSuccessResponse` and `NewErrorResponse` handle JSON marshalling for you.
- **Idiomatic Go**: Implements standard initialisms, `Unwrap` for multi-errors, and the `error` interface for RPC codes.

## Installation

```bash
go get github.com/sprokhorov/jsonrpc
```

## Quick Start

### 1. Parsing and Validating a Request

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/sprokhorov/jsonrpc"
)

func main() {
	// Example raw JSON request (numeric ID)
	raw := []byte(`{"jsonrpc": "2.0", "method": "sum", "params": [1, 2, 3], "id": 1}`)

	var req jsonrpc.Request
	if err := json.Unmarshal(raw, &req); err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	if err := req.Validate(); err != nil {
		fmt.Printf("Validation error: %v\n", err)
		return
	}

	fmt.Printf("Method: %s, ID: %s\n", req.Method, string(*req.ID))
}
```

### 2. Creating a Success Response

The `NewSuccessResponse` helper automatically marshals your result to JSON.

```go
type MyResult struct {
    Total int `json:"total"`
}

// id is usually taken from the request: req.ID
resp := jsonrpc.NewSuccessResponse(*req.ID, MyResult{Total: 6})

// Marshal to send back over the wire
data, _ := json.Marshal(resp)
fmt.Println(string(data)) 
// Output: {"jsonrpc":"2.0","id":1,"result":{"total":6}}
```

### 3. Creating an Error Response

```go
// Standard error without extra data
errResp := jsonrpc.NewErrorResponse(*req.ID, jsonrpc.CodeMethodNotFound, nil)

// Custom error with details (auto-marshaled)
details := map[string]string{"reason": "negative numbers not allowed"}
customErr := jsonrpc.NewErrorResponse(*req.ID, jsonrpc.CodeInvalidParams, details)
```

### 4. Handling Notifications

A notification is simply a request where the `id` field is omitted.

```go
req := jsonrpc.Request{
    JSONRPC: "2.0",
    Method:  "log_message",
    Params:  json.RawMessage(`["System started"]`),
}

if req.IsNotification() {
    fmt.Println("No response needed for this request")
}
```

## Standard Error Codes

The library provides the standard JSON-RPC 2.0 error codes, which also implement the `error` interface:

- `CodeParseError` (-32700)
- `CodeInvalidRequest` (-32600)
- `CodeMethodNotFound` (-32601)
- `CodeInvalidParams` (-32602)
- `CodeInternalError` (-32603)

## License

MIT
