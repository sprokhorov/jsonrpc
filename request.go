package jsonrpc

import (
	"bytes"
	"encoding/json"
)

// NewRequest creates a new Request object.
func NewRequest(id *json.RawMessage, method string, params []byte) *Request {
	return &Request{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Method:  method,
		Params:  json.RawMessage(params),
	}
}

// Request is a JSON-RPC request object.
type Request struct {
	// JSONRPC is the version of the JSON-RPC protocol. Must be exactly "2.0".
	JSONRPC string `json:"jsonrpc"`
	// ID is an identifier established by the Client. It MUST contain a String,
	// Number, or NULL value if included. If it is not included it is assumed
	// to be a notification.
	ID *json.RawMessage `json:"id,omitempty"`
	// Method is a string containing the name of the method to be invoked.
	Method string `json:"method"`
	// Params is a structured value that holds the parameter values to be used
	// during the invocation of the method.
	Params json.RawMessage `json:"params,omitempty"`
}

// IsNotification returns true if the request is a notification (has no ID).
func (req *Request) IsNotification() bool {
	return req.ID == nil
}

// Validate validates the request object following the JSON-RPC 2.0 specification.
// It returns an error if the request is invalid. Both server and client must
// validate the request before sending it.
func (req *Request) Validate() error {
	var errs ValidationErrors
	if req.JSONRPC != JSONRPCVersion {
		errs.Append(ErrInvalidJSONRPCVersion)
	}
	if req.Method == "" {
		errs.Append(ErrEmptyMethod)
	}

	if req.Params != nil {
		p := bytes.TrimSpace(req.Params)
		if len(p) > 0 && p[0] != '[' && p[0] != '{' {
			errs.Append(ErrInvalidParams)
		}
	}

	if errs.HasErrors() {
		return &errs
	}
	return nil
}
