package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// FuzzFn is called when the MCP client invokes the "fuzz" tool.
type FuzzFn func(args map[string]any) (string, error)

type mcpRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type mcpResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      any `json:"id"`
	Result  any `json:"result,omitempty"`
	Error   *rpcError   `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Serve starts the MCP server over stdio. Blocks until stdin is closed.
func Serve(fuzzFn FuzzFn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var req mcpRequest
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			continue
		}
		// Notifications have no id — no response required
		if req.ID == nil {
			continue
		}
		result, rpcErr := dispatch(req, fuzzFn)
		resp := mcpResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  result,
			Error:   rpcErr,
		}
		data, _ := json.Marshal(resp)
		fmt.Println(string(data))
	}
}

func dispatch(req mcpRequest, fuzzFn FuzzFn) (any, *rpcError) {
	switch req.Method {
	case "initialize":
		return map[string]any{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]any{"tools": map[string]any{}},
			"serverInfo":      map[string]any{"name": "cfuzz", "version": "2.0.0"},
		}, nil
	case "tools/list":
		return map[string]any{"tools": []any{fuzzToolSchema()}}, nil
	case "tools/call":
		var p struct {
			Name      string                 `json:"name"`
			Arguments map[string]any `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &p); err != nil {
			return nil, &rpcError{Code: -32600, Message: err.Error()}
		}
		if p.Name != "fuzz" {
			return nil, &rpcError{Code: -32601, Message: "unknown tool: " + p.Name}
		}
		output, err := fuzzFn(p.Arguments)
		if err != nil {
			return nil, &rpcError{Code: -32000, Message: err.Error()}
		}
		return map[string]any{
			"content": []any{
				map[string]any{"type": "text", "text": output},
			},
		}, nil
	default:
		return nil, &rpcError{Code: -32601, Message: "method not found: " + req.Method}
	}
}

func fuzzToolSchema() map[string]any {
	return map[string]any{
		"name":        "fuzz",
		"description": "Fuzz a shell command by substituting wordlist entries for the FUZZ keyword and returning filtered results.",
		"inputSchema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"command":      map[string]any{"type": "string", "description": "Shell command containing FUZZ keyword"},
				"wordlist":     map[string]any{"type": "array", "items": map[string]any{"type": "string"}, "description": "Words to substitute for FUZZ"},
				"threads":      map[string]any{"type": "integer", "description": "Max concurrent workers (default 50)"},
				"timeout":      map[string]any{"type": "integer", "description": "Per-command timeout in seconds (default 30)"},
				"success_only": map[string]any{"type": "boolean", "description": "Only return results with exit code 0"},
				"stdout_word":  map[string]any{"type": "string", "description": "Only return results where stdout contains this word"},
				"ai_filter":    map[string]any{"type": "string", "description": "Natural language filter (requires ANTHROPIC_API_KEY)"},
			},
			"required": []string{"command", "wordlist"},
		},
	}
}
