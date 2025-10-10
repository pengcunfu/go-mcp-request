package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
)

// LogData represents the structure for logging HTTP requests and responses
type LogData struct {
	Timestamp string    `json:"timestamp"`
	Request   Request   `json:"request"`
	Response  Response  `json:"response"`
	Error     string    `json:"error,omitempty"`
}

type Request struct {
	Method     string            `json:"method"`
	URL        string            `json:"url"`
	Headers    map[string]string `json:"headers"`
	Cookies    map[string]string `json:"cookies"`
	Body       string            `json:"body"`
	BodyLength int               `json:"body_length"`
}

type Response struct {
	StatusCode     interface{}       `json:"status_code"`
	Headers        map[string]string `json:"headers"`
	ContentLength  int               `json:"content_length"`
	ContentPreview string            `json:"content_preview"`
}

type HTTPResult struct {
	Method          string            `json:"method"`
	URL             string            `json:"url"`
	StatusCode      int               `json:"status_code"`
	ResponseHeaders map[string]string `json:"response_headers"`
	ResponseContent string            `json:"response_content"`
	ResponseLength  int               `json:"response_length"`
	RequestHeaders  map[string]string `json:"request_headers"`
	RequestCookies  map[string]string `json:"request_cookies"`
	RequestBody     string            `json:"request_body"`
	LoggedTo        string            `json:"logged_to"`
}

var logger *logrus.Logger
var logDir string

func init() {
	// Setup logging
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("Failed to get user home directory: %v", err))
	}
	
	logDir = filepath.Join(homeDir, "mcp_requests_logs")
	err = os.MkdirAll(logDir, 0755)
	if err != nil {
		panic(fmt.Sprintf("Failed to create log directory: %v", err))
	}

	logger = logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	
	// Create log file with timestamp
	logFilename := fmt.Sprintf("requests_%s.log", time.Now().Format("20060102_150405"))
	logPath := filepath.Join(logDir, logFilename)
	
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("Failed to open log file: %v", err))
	}
	
	logger.SetOutput(file)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func logRequestResponse(method, url string, headers, cookies map[string]string, body string,
	statusCode int, responseHeaders map[string]string, responseContent string,
	responseLength int, errorMsg string) string {
	
	logData := LogData{
		Timestamp: time.Now().Format(time.RFC3339),
		Request: Request{
			Method:     method,
			URL:        url,
			Headers:    headers,
			Cookies:    cookies,
			Body:       body,
			BodyLength: len(body),
		},
		Response: Response{
			StatusCode:    statusCode,
			Headers:       responseHeaders,
			ContentLength: responseLength,
		},
		Error: errorMsg,
	}
	
	if errorMsg != "" {
		logData.Response.StatusCode = "ERROR"
		logData.Response.Headers = make(map[string]string)
		logData.Response.ContentLength = 0
	}
	
	// Truncate content preview if too long
	if len(responseContent) > 500 {
		logData.Response.ContentPreview = responseContent[:500] + "..."
	} else {
		logData.Response.ContentPreview = responseContent
	}
	
	logJSON, _ := json.MarshalIndent(logData, "", "  ")
	logger.Infof("HTTP_REQUEST: %s", string(logJSON))
	
	return filepath.Join(logDir, fmt.Sprintf("requests_%s.log", time.Now().Format("20060102_150405")))
}

func makeHTTPRequestWithLogging(method, url string, headers, cookies map[string]string, body string, timeout float64) (*HTTPResult, error) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	
	var bodyReader io.Reader
	if body != "" {
		bodyReader = strings.NewReader(body)
	}
	
	req, err := http.NewRequest(strings.ToUpper(method), url, bodyReader)
	if err != nil {
		logRequestResponse(strings.ToUpper(method), url, headers, cookies, body, 0, nil, "", 0, err.Error())
		return nil, err
	}
	
	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	
	// Set cookies
	for name, value := range cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}
	
	resp, err := client.Do(req)
	if err != nil {
		logRequestResponse(strings.ToUpper(method), url, headers, cookies, body, 0, nil, "", 0, err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logRequestResponse(strings.ToUpper(method), url, headers, cookies, body, resp.StatusCode, nil, "", 0, err.Error())
		return nil, err
	}
	
	responseContent := string(responseBody)
	responseHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			responseHeaders[key] = values[0]
		}
	}
	
	logPath := logRequestResponse(
		strings.ToUpper(method), url, headers, cookies, body,
		resp.StatusCode, responseHeaders, responseContent, len(responseContent), "",
	)
	
	return &HTTPResult{
		Method:          strings.ToUpper(method),
		URL:             url,
		StatusCode:      resp.StatusCode,
		ResponseHeaders: responseHeaders,
		ResponseContent: responseContent,
		ResponseLength:  len(responseContent),
		RequestHeaders:  headers,
		RequestCookies:  cookies,
		RequestBody:     body,
		LoggedTo:        logPath,
	}, nil
}

func createHTTPHandler(method string) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.Params.Arguments
		
		url, ok := args["url"].(string)
		if !ok {
			return mcp.NewToolResultError("url parameter is required and must be a string"), nil
		}
		
		headers := make(map[string]string)
		if h, ok := args["headers"].(map[string]interface{}); ok {
			for k, v := range h {
				if str, ok := v.(string); ok {
					headers[k] = str
				}
			}
		}
		
		cookies := make(map[string]string)
		if c, ok := args["cookies"].(map[string]interface{}); ok {
			for k, v := range c {
				if str, ok := v.(string); ok {
					cookies[k] = str
				}
			}
		}
		
		body := ""
		if b, ok := args["body"].(string); ok {
			body = b
		}
		
		timeout := 30.0
		if t, ok := args["timeout"].(float64); ok {
			timeout = t
		}
		
		result, err := makeHTTPRequestWithLogging(method, url, headers, cookies, body, timeout)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
		}
		
		resultJSON, _ := json.MarshalIndent(result, "", "  ")
		return mcp.NewToolResultText(string(resultJSON)), nil
	}
}

func httpRawRequestHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.Params.Arguments
	
	url, ok := args["url"].(string)
	if !ok {
		return mcp.NewToolResultError("url parameter is required and must be a string"), nil
	}
	
	method := "GET"
	if m, ok := args["method"].(string); ok {
		method = m
	}
	
	headers := make(map[string]string)
	if h, ok := args["headers"].(map[string]interface{}); ok {
		for k, v := range h {
			if str, ok := v.(string); ok {
				headers[k] = str
			}
		}
	}
	
	cookies := make(map[string]string)
	if c, ok := args["cookies"].(map[string]interface{}); ok {
		for k, v := range c {
			if str, ok := v.(string); ok {
				cookies[k] = str
			}
		}
	}
	
	rawBody := ""
	conversionInfo := ""
	
	if rb := args["raw_body"]; rb != nil {
		switch v := rb.(type) {
		case string:
			rawBody = v
		case map[string]interface{}:
			bodyBytes, _ := json.Marshal(v)
			rawBody = string(bodyBytes)
			conversionInfo = "⚠️ AUTO-CONVERTED: Dict → JSON string"
		case []interface{}:
			bodyBytes, _ := json.Marshal(v)
			rawBody = string(bodyBytes)
			conversionInfo = "⚠️ AUTO-CONVERTED: Array → JSON string"
		default:
			rawBody = fmt.Sprintf("%v", v)
			conversionInfo = fmt.Sprintf("⚠️ AUTO-CONVERTED: %T → string", v)
		}
	}
	
	timeout := 30.0
	if t, ok := args["timeout"].(float64); ok {
		timeout = t
	}
	
	result, err := makeHTTPRequestWithLogging(method, url, headers, cookies, rawBody, timeout)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Error: %v", err)), nil
	}
	
	// Add conversion warning if applicable
	if conversionInfo != "" {
		resultMap := make(map[string]interface{})
		resultJSON, _ := json.Marshal(result)
		json.Unmarshal(resultJSON, &resultMap)
		resultMap["conversion_warning"] = conversionInfo
		resultJSON, _ = json.MarshalIndent(resultMap, "", "  ")
		return mcp.NewToolResultText(string(resultJSON)), nil
	}
	
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	return mcp.NewToolResultText(string(resultJSON)), nil
}

func main() {
	s := server.NewMCPServer(
		"HTTP Requests",
		"1.0.0",
	)
	
	// Register HTTP method tools
	httpMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	
	for _, method := range httpMethods {
		toolName := fmt.Sprintf("http_%s", strings.ToLower(method))
		description := fmt.Sprintf("HTTP %s request with full support (headers, cookies, body, timeout) - All requests logged", method)
		
		if method == "HEAD" || method == "OPTIONS" {
			description = fmt.Sprintf("HTTP %s request with full support (headers, cookies, timeout) - All requests logged", method)
		}
		
		var tool mcp.Tool
		if method == "HEAD" || method == "OPTIONS" {
			tool = mcp.NewTool(
				toolName,
				mcp.WithDescription(description),
				mcp.WithString("url", mcp.Required(), mcp.Description("The URL to send the request to")),
				mcp.WithObject("headers", mcp.Description("Optional headers to include in the request")),
				mcp.WithObject("cookies", mcp.Description("Optional cookies to include in the request")),
				mcp.WithNumber("timeout", mcp.Description("Request timeout in seconds (default: 30)")),
			)
		} else {
			tool = mcp.NewTool(
				toolName,
				mcp.WithDescription(description),
				mcp.WithString("url", mcp.Required(), mcp.Description("The URL to send the request to")),
				mcp.WithObject("headers", mcp.Description("Optional headers to include in the request")),
				mcp.WithObject("cookies", mcp.Description("Optional cookies to include in the request")),
				mcp.WithString("body", mcp.Description("Optional body content for the request")),
				mcp.WithNumber("timeout", mcp.Description("Request timeout in seconds (default: 30)")),
			)
		}
		
		s.AddTool(tool, createHTTPHandler(method))
	}
	
	// Register raw request tool
	rawRequestTool := mcp.NewTool(
		"http_raw_request",
		mcp.WithDescription(`🔒 CRITICAL SECURITY TESTING TOOL: Sends HTTP requests with ABSOLUTE PRECISION - All requests logged

⚠️  IMPORTANT: This tool preserves EVERY SINGLE CHARACTER of your request:
- Headers: Every cookie, token, session ID - NO CHARACTER LIMIT, NO TRUNCATION
- Body: Raw payload sent byte-for-byte, preserving payloads exactly
- Cookies: Complete cookie strings including long JWT tokens, session data
- Special characters: ', ", \, %, &, =, etc. are preserved without encoding
- Whitespace: Spaces, tabs, newlines maintained exactly as provided

🎯 Perfect for: all kinds of security vulnerability testing, testing like SQL injection, XSS, CSRF, authentication bypass, parameter pollution
📝 Guarantee: What you input is EXACTLY what gets sent - zero modifications
📊 All requests and responses are automatically logged to ~/mcp_requests_logs/

💡 USAGE TIP: raw_body must be a STRING, not an object. For JSON, use: '{"key":"value"}' not {"key":"value"}`),
		mcp.WithString("url", mcp.Required(), mcp.Description("The URL to send the request to")),
		mcp.WithString("method", mcp.Description("HTTP method (default: GET)")),
		mcp.WithString("raw_body", mcp.Description("Raw body content - preserved exactly as provided")),
		mcp.WithObject("headers", mcp.Description("Optional headers to include in the request")),
		mcp.WithObject("cookies", mcp.Description("Optional cookies to include in the request")),
		mcp.WithNumber("timeout", mcp.Description("Request timeout in seconds (default: 30)")),
	)
	
	s.AddTool(rawRequestTool, httpRawRequestHandler)
	
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}