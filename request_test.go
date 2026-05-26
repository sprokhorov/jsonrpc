package jsonrpc

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRequest_Validate(t *testing.T) {
	id1 := json.RawMessage(`1`)
	idString := json.RawMessage(`"abc"`)
	idNull := json.RawMessage(`null`)

	tests := []struct {
		name    string
		req     Request
		wantErr bool
		errMsgs []string
	}{
		{
			name: "Valid request with numeric ID",
			req: Request{
				JSONRPC: "2.0",
				ID:      &id1,
				Method:  "test",
			},
			wantErr: false,
		},
		{
			name: "Valid request with string ID",
			req: Request{
				JSONRPC: "2.0",
				ID:      &idString,
				Method:  "test",
			},
			wantErr: false,
		},
		{
			name: "Valid request with null ID",
			req: Request{
				JSONRPC: "2.0",
				ID:      &idNull,
				Method:  "test",
			},
			wantErr: false,
		},
		{
			name: "Valid notification (missing ID)",
			req: Request{
				JSONRPC: "2.0",
				Method:  "test",
			},
			wantErr: false,
		},
		{
			name: "Invalid version",
			req: Request{
				JSONRPC: "1.0",
				ID:      &id1,
				Method:  "test",
			},
			wantErr: true,
			errMsgs: []string{ErrInvalidJSONRPCVersion.Error()},
		},
		{
			name: "Empty method",
			req: Request{
				JSONRPC: "2.0",
				ID:      &id1,
				Method:  "",
			},
			wantErr: true,
			errMsgs: []string{ErrEmptyMethod.Error()},
		},
		{
			name: "Invalid Params (not array or object)",
			req: Request{
				JSONRPC: "2.0",
				ID:      &id1,
				Method:  "test",
				Params:  json.RawMessage(`123`),
			},
			wantErr: true,
			errMsgs: []string{ErrInvalidParams.Error()},
		},
		{
			name: "Valid Params (array)",
			req: Request{
				JSONRPC: "2.0",
				ID:      &id1,
				Method:  "test",
				Params:  json.RawMessage(`[1, 2]`),
			},
			wantErr: false,
		},
		{
			name: "Valid Params (object)",
			req: Request{
				JSONRPC: "2.0",
				ID:      &id1,
				Method:  "test",
				Params:  json.RawMessage(`{"a": 1}`),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Request.Validate() error = %v, wantErr %v", err, tt.wantErr)
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

func TestRequest_IsNotification(t *testing.T) {
	id := json.RawMessage(`1`)
	req1 := Request{ID: &id}
	if req1.IsNotification() {
		t.Errorf("Expected IsNotification to be false for request with ID")
	}

	req2 := Request{ID: nil}
	if !req2.IsNotification() {
		t.Errorf("Expected IsNotification to be true for request without ID")
	}
}

func TestNewRequest(t *testing.T) {
	id := json.RawMessage(`"1"`)
	params := []byte(`{"foo":"bar"}`)
	req := NewRequest(&id, "test", params)

	if req.JSONRPC != JSONRPCVersion {
		t.Errorf("Expected JSONRPC %s, got %s", JSONRPCVersion, req.JSONRPC)
	}
	if string(*req.ID) != string(id) {
		t.Errorf("Expected ID %s, got %s", string(id), string(*req.ID))
	}
	if req.Method != "test" {
		t.Errorf("Expected method test, got %s", req.Method)
	}
	if string(req.Params) != string(params) {
		t.Errorf("Expected params %s, got %s", string(params), string(req.Params))
	}
}
