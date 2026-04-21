package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// FuzzFn is called when the MCP client invokes the "fuzz" tool.
type FuzzFn func(args map[string]interface{}) (string, error)

type mcpRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type mcpResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
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

func dispatch(req mcpRequest, fuzzFn FuzzFn) (interface{}, *rpcError) {
	switch req.Method {
	case "initialize":
		return map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities":    map[string]interface{}{"tools": map[string]interface{}{}},
			"serverInfo":      map[string]interface{}{"name": "cfuzz", "version": "2.0.0"},
		}, nil
	case "tools/list":
		return map[string]interface{}{"tools": []interface{}{fuzzToolSchema()}}, nil
	case "tools/call":
		var p struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
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
		return map[string]interface{}{
			"content": []interface{}{
				map[string]interface{}{"type": "text", "text": output},
			},
		}, nil
	default:
		return nil, &rpcError{Code: -32601, Message: "method not found: " + req.Method}
	}
}

func fuzzToolSchema() map[string]interface{} {
	return map[string]interface{}{
		"name":        "fuzz",
		"description": "Fuzz a shell command by substituting wordlist entries for the FUZZ keyword and returning filtered results.",
		"inputSchema": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"command":      map[string]interface{}{"type": "string", "description": "Shell command containing FUZZ keyword"},
				"wordlist":     map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}, "description": "Words to substitute for FUZZ"},
				"threads":      map[string]interface{}{"type": "integer", "description": "Max concurrent workers (default 50)"},
				"timeout":      map[string]interface{}{"type": "integer", "description": "Per-command timeout in seconds (default 30)"},
				"success_only": map[string]interface{}{"type": "boolean", "description": "Only return results with exit code 0"},
				"stdout_word":  map[string]interface{}{"type": "string", "description": "Only return results where stdout contains this word"},
				"ai_filter":    map[string]interface{}{"type": "string", "description": "Natural language filter (requires ANTHROPIC_API_KEY)"},
			},
			"required": []string{"command", "wordlist"},
		},
	}
}
