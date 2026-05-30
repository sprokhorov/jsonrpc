package jsonrpc

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Error codes are defined in the JSON-RPC specification.
type ErrorCode int

// Message returns the message associated with the given code.
func (c ErrorCode) Message() string {
	switch c {
	case CodeParseError:
		return "Parse error"

	case CodeInvalidRequest:
		return "Invalid Request"

	case CodeMethodNotFound:
		return "Method not found"

	case CodeInvalidParams:
		return "Invalid params"

	case CodeInternalError:
		return "Internal error"

	default:
		return "Unknown error"
	}
}

// Error implements the error interface.
func (c ErrorCode) Error() string {
	return c.Message()
}

// GoError returns an error for the given JSON-RPC error code.
func GoError(c ErrorCode) error {
	return errors.New(c.Message())
}

const (
	// Standard JSON-RPC 2.0 errors

	CodeParseError     ErrorCode = -32700
	CodeInvalidRequest ErrorCode = -32600
	CodeMethodNotFound ErrorCode = -32601
	CodeInvalidParams  ErrorCode = -32602
	CodeInternalError  ErrorCode = -32603
)

// NewResponse creates a new Response object with the provided ID, result, and error.
// Most users should use NewSuccessResponse or NewErrorResponse instead.
func NewResponse(id json.RawMessage, result json.RawMessage, err *ResponseError) *Response {
	return &Response{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Result:  result,
		Error:   err,
	}
}

// NewSuccessResponse creates a new success Response object.
// The result is marshaled to JSON. If marshalling fails, an internal error response is returned.
func NewSuccessResponse(id json.RawMessage, result any) *Response {
	res, err := json.Marshal(result)
	if err != nil {
		return NewErrorResponse(id, CodeInternalError, err.Error())
	}
	return NewResponse(id, res, nil)
}

// NewErrorResponse creates a new error Response object.
// The data is optional and is marshaled to JSON if provided.
func NewErrorResponse(id json.RawMessage, code ErrorCode, data any) *Response {
	return NewResponse(id, nil, NewResponseError(code, data))
}

// Response is a JSON-RPC response object.
type Response struct {
	// JSONRPC is the version of the JSON-RPC protocol. Must be exactly "2.0".
	JSONRPC string `json:"jsonrpc"`
	// ID is the identifier established by the Client. It must be the same as
	// the ID in the Request it is responding to. If there was an error in
	// detecting the id in the Request object (e.g. Parse error/Invalid Request),
	// it MUST be Null.
	ID json.RawMessage `json:"id"`
	// Result is the result of the method invocation. Required on success.
	Result json.RawMessage `json:"result,omitempty"`
	// Error is the error object if the method invocation failed. Required on error.
	Error *ResponseError `json:"error,omitempty"`
}

// Validate validates the response object following the JSON-RPC 2.0 specification.
// It returns an error if the response is invalid. Both server and client must
// validate the response before sending it.
func (res *Response) Validate() error {
	var errs ValidationErrors
	if res.JSONRPC != JSONRPCVersion {
		errs.Append(ErrInvalidJSONRPCVersion)
	}
	if len(res.ID) == 0 {
		errs.Append(ErrEmptyId)
	}
	if res.Error != nil && res.Result != nil {
		errs.Append(ErrResultAndError)
	}
	if res.Error == nil && res.Result == nil {
		errs.Append(ErrResultOrError)
	}
	if errs.HasErrors() {
		return &errs
	}
	return nil
}

func (res *Response) HTTPCode() int {
	if res.Error != nil {
		switch res.Error.Code {
		case CodeParseError, CodeInvalidRequest:
			return http.StatusBadRequest
		case CodeMethodNotFound:
			return http.StatusNotFound
		case CodeInternalError:
			return http.StatusInternalServerError
		case CodeInvalidParams:
			return http.StatusUnprocessableEntity
		}
	}
	return http.StatusOK
}

// ResponseError is a JSON-RPC error object.
type ResponseError struct {
	// Code is a number that indicates the error type that occurred.
	Code ErrorCode `json:"code"`
	// Message is a string providing a short description of the error.
	Message string `json:"message"`
	// Data is a structured value that contains additional information about the error.
	Data json.RawMessage `json:"data,omitempty"`
}

// NewResponseError creates a new ResponseError object.
// The data is optional and is marshaled to JSON if provided.
func NewResponseError(code ErrorCode, data any) *ResponseError {
	var rawData json.RawMessage
	if data != nil {
		if d, ok := data.([]byte); ok {
			rawData = json.RawMessage(d)
		} else if d, ok := data.(json.RawMessage); ok {
			rawData = d
		} else {
			if b, err := json.Marshal(data); err == nil {
				rawData = json.RawMessage(b)
			}
		}
	}

	return &ResponseError{
		Code:    code,
		Message: code.Message(),
		Data:    rawData,
	}
}
