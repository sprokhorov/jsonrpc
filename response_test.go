package jsonrpc

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestErrorCode_Message(t *testing.T) {
	tests := []struct {
		code ErrorCode
		want string
	}{
		{CodeParseError, "Parse error"},
		{CodeInvalidRequest, "Invalid Request"},
		{CodeMethodNotFound, "Method not found"},
		{CodeInvalidParams, "Invalid params"},
		{CodeInternalError, "Internal error"},
		{ErrorCode(123), "Unknown error"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.code.Message(); got != tt.want {
				t.Errorf("ErrorCode.Message() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewResponseError(t *testing.T) {
	t.Run("with []byte data", func(t *testing.T) {
		data := []byte(`{"info":"test"}`)
		err := NewResponseError(CodeInvalidParams, data)
		if string(err.Data) != string(data) {
			t.Errorf("Expected data %s, got %s", string(data), string(err.Data))
		}
	})

	t.Run("with struct data", func(t *testing.T) {
		data := struct{ Foo string }{Foo: "bar"}
		err := NewResponseError(CodeInvalidParams, data)
		expected := `{"Foo":"bar"}`
		if string(err.Data) != expected {
			t.Errorf("Expected data %s, got %s", expected, string(err.Data))
		}
	})

	t.Run("with nil data", func(t *testing.T) {
		err := NewResponseError(CodeInvalidParams, nil)
		if err.Data != nil {
			t.Errorf("Expected nil data, got %s", string(err.Data))
		}
	})
}

func TestNewSuccessResponse(t *testing.T) {
	id := json.RawMessage(`1`)
	result := map[string]int{"val": 42}
	resp := NewSuccessResponse(id, result)

	if resp.JSONRPC != JSONRPCVersion {
		t.Errorf("Expected JSONRPC %s, got %s", JSONRPCVersion, resp.JSONRPC)
	}
	if string(resp.ID) != "1" {
		t.Errorf("Expected ID 1, got %s", string(resp.ID))
	}
	expectedResult := `{"val":42}`
	if string(resp.Result) != expectedResult {
		t.Errorf("Expected result %s, got %s", expectedResult, string(resp.Result))
	}
	if resp.Error != nil {
		t.Errorf("Expected nil error, got %v", resp.Error)
	}
}

func TestNewErrorResponse(t *testing.T) {
	id := json.RawMessage(`1`)
	data := "some error details"
	resp := NewErrorResponse(id, CodeMethodNotFound, data)

	if resp.JSONRPC != JSONRPCVersion {
		t.Errorf("Expected JSONRPC %s, got %s", JSONRPCVersion, resp.JSONRPC)
	}
	if resp.Result != nil {
		t.Errorf("Expected nil result, got %s", string(resp.Result))
	}
	if resp.Error == nil {
		t.Fatal("Expected non-nil error")
	}
	if resp.Error.Code != CodeMethodNotFound {
		t.Errorf("Expected code %v, got %v", CodeMethodNotFound, resp.Error.Code)
	}
	expectedData := `"some error details"`
	if string(resp.Error.Data) != expectedData {
		t.Errorf("Expected error data %s, got %s", expectedData, string(resp.Error.Data))
	}
}

func TestResponse_Validate(t *testing.T) {
	id1 := json.RawMessage(`1`)
	idString := json.RawMessage(`"abc"`)
	idNull := json.RawMessage(`null`)

	tests := []struct {
		name    string
		res     Response
		wantErr bool
		errMsgs []string
	}{
		{
			name: "Valid response with numeric ID",
			res: Response{
				JSONRPC: "2.0",
				ID:      id1,
				Result:  json.RawMessage(`"ok"`),
			},
			wantErr: false,
		},
		{
			name: "Valid response with string ID",
			res: Response{
				JSONRPC: "2.0",
				ID:      idString,
				Result:  json.RawMessage(`"ok"`),
			},
			wantErr: false,
		},
		{
			name: "Valid response with null ID",
			res: Response{
				JSONRPC: "2.0",
				ID:      idNull,
				Result:  json.RawMessage(`"ok"`),
			},
			wantErr: false,
		},
		{
			name: "Valid response with error",
			res: Response{
				JSONRPC: "2.0",
				ID:      id1,
				Error:   &ResponseError{Code: CodeInternalError, Message: "Internal error"},
			},
			wantErr: false,
		},
		{
			name: "Invalid version",
			res: Response{
				JSONRPC: "1.0",
				ID:      id1,
				Result:  json.RawMessage(`"ok"`),
			},
			wantErr: true,
			errMsgs: []string{ErrInvalidJSONRPCVersion.Error()},
		},
		{
			name: "Empty ID",
			res: Response{
				JSONRPC: "2.0",
				Result:  json.RawMessage(`"ok"`),
			},
			wantErr: true,
			errMsgs: []string{ErrEmptyId.Error()},
		},
		{
			name: "Result and Error set",
			res: Response{
				JSONRPC: "2.0",
				ID:      id1,
				Result:  json.RawMessage(`"ok"`),
				Error:   &ResponseError{Code: CodeInternalError, Message: "Internal error"},
			},
			wantErr: true,
			errMsgs: []string{ErrResultAndError.Error()},
		},
		{
			name: "Neither Result nor Error set",
			res: Response{
				JSONRPC: "2.0",
				ID:      id1,
			},
			wantErr: true,
			errMsgs: []string{ErrResultOrError.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.res.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Response.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				errMsg := err.Error()
				for _, expected := range tt.errMsgs {
					if !strings.Contains(errMsg, expected) {
						t.Errorf("Expected error message to contain '%s', got '%s'", expected, errMsg)
					}
				}
			}
		})
	}
}
