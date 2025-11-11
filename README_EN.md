<div align="center">

# 🚀 MCP HTTP Request Tool

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.21-00ADD8?logo=go)](https://go.dev/)
[![GitHub release](https://img.shields.io/github/v/release/pengcunfu/go-mcp-request)](https://github.com/pengcunfu/go-mcp-request/releases)
[![GitHub stars](https://img.shields.io/github/stars/pengcunfu/go-mcp-request)](https://github.com/pengcunfu/go-mcp-request/stargazers)

English | [简体中文](./README.md)

A full-featured HTTP client MCP (Model Context Protocol) server designed for API testing, web automation, and security testing.

Complete HTTP tools with detailed logging capabilities, written in Go for superior performance and simple deployment.

</div>

---

## ✨ Features

- 🔧 **Complete HTTP Method Support** - GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
- 🔒 **Advanced Security Testing** - Raw request tool supporting penetration testing, SQL injection, XSS testing
- ⚙️ **Full Parameter Support** - All methods support Headers, Cookies, Body, and timeout settings
- 📝 **Automatic Logging** - All requests and responses automatically logged to `~/mcp_requests_logs/`
- 🎯 **Precision Guaranteed** - Raw mode preserves every character without automatic encoding
- 🔌 **MCP Compatible** - Perfect compatibility with Claude Desktop, Cursor, and other MCP clients
- ⚡ **High Performance** - Go implementation, single binary, no runtime dependencies
- 🌍 **Cross-Platform** - Supports Windows, Linux, macOS

## 📦 Installation

### Option 1: Download Pre-compiled Binary (Recommended)

Download the latest version for your operating system from the [Releases page](https://github.com/pengcunfu/go-mcp-request/releases):

- **Windows**: `mcp-request-windows-amd64.exe`
- **Linux**: `mcp-request-linux-amd64`
- **macOS**: `mcp-request-darwin-amd64`

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/pengcunfu/go-mcp-request.git
cd go-mcp-request

# Build
go build -o mcp-request main.go

# Or use go install
go install github.com/pengcunfu/go-mcp-request@latest
```

## 🚀 Usage

### Configuration for Claude Desktop

Edit the configuration file `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or `%APPDATA%\Claude\claude_desktop_config.json` (Windows):

```json
{
  "mcpServers": {
    "mcp-request": {
      "command": "/path/to/mcp-request",
      "args": []
    }
  }
}
```

### Configuration for Cursor

Edit the configuration file `~/.cursor/mcp_servers.json`:

```json
{
  "mcpServers": {
    "mcp-request": {
      "command": "/path/to/mcp-request",
      "type": "stdio"
    }
  }
}
```

> **Tip**: Replace `/path/to/mcp-request` with the actual path to the binary file. If added to PATH environment variable, you can use `mcp-request` directly.

## 🛠️ Available Tools

| Tool Name | Description | Use Cases |
|-----------|-------------|-----------|
| `http_get` | Full-featured GET request | Retrieve resources, API queries |
| `http_post` | Full-featured POST request | Create resources, form submission |
| `http_put` | Full-featured PUT request | Update resources |
| `http_delete` | Full-featured DELETE request | Delete resources |
| `http_patch` | Full-featured PATCH request | Partial resource updates |
| `http_head` | Full-featured HEAD request | Get response headers |
| `http_options` | Full-featured OPTIONS request | Check supported methods |
| `http_raw_request` | Raw HTTP request | Security testing, penetration testing |

## 📖 Usage Examples

### Basic GET Request

```python
# Simple GET request
http_get("https://api.example.com/users")

# GET request with headers
http_get(
    url="https://api.example.com/users",
    headers={"Authorization": "Bearer token123"}
)
```

### POST Request

```python
# POST with JSON data
http_post(
    url="https://api.example.com/login",
    body='{"username":"test","password":"test"}',
    headers={"Content-Type": "application/json"}
)

# POST with form data
http_post(
    url="https://api.example.com/form",
    body="name=John&age=30",
    headers={"Content-Type": "application/x-www-form-urlencoded"}
)
```

### Security Testing (Raw Request)

```python
# SQL injection testing
http_raw_request(
    url="https://test-site.com/search",
    method="POST",
    raw_body="q=test' OR 1=1--",
    headers={"Content-Type": "application/x-www-form-urlencoded"}
)

# XSS testing
http_raw_request(
    url="https://test-site.com/comment",
    method="POST",
    raw_body='comment=<script>alert("XSS")</script>',
    headers={"Content-Type": "application/x-www-form-urlencoded"}
)
```

## 🔒 Security Testing Features

The `http_raw_request` tool is designed specifically for security research and penetration testing:

### Core Advantages

- ✅ **Absolute Precision** - Every byte is preserved exactly without any modification
- ✅ **No Automatic Encoding** - Special characters (`'`, `"`, `\`, `%`, `&`, `=`) are sent as-is
- ✅ **Complete Headers** - No truncation of long cookies or tokens
- ✅ **Raw Payloads** - Perfect for all types of security testing

### Applicable Test Scenarios

- 🎯 SQL Injection Testing
- 🎯 XSS (Cross-Site Scripting) Testing
- 🎯 CSRF (Cross-Site Request Forgery) Testing
- 🎯 Command Injection Testing
- 🎯 HTTP Request Smuggling Testing
- 🎯 Custom Protocol Testing

> ⚠️ **Disclaimer**: This tool is intended for legal security testing and research purposes only. Ensure you have permission to test the target system. Unauthorized testing may violate laws.

## 📝 Logging

All HTTP requests and responses are automatically logged for debugging and auditing purposes.

### Log Configuration

- **Storage Location**: `~/mcp_requests_logs/`
- **File Format**: JSON format, structured storage
- **File Naming**: `requests_YYYYMMDD_HHMMSS.log`
- **Included Information**:
  - Timestamp
  - Complete request information (method, URL, headers, body)
  - Complete response information (status code, headers, body)
  - Error information (if any)

### Viewing Logs

```bash
# View latest logs in real-time
tail -f ~/mcp_requests_logs/requests_*.log

# View logs for a specific date
cat ~/mcp_requests_logs/requests_20250111_*.log

# Format view with jq
cat ~/mcp_requests_logs/requests_*.log | jq .
```

## 💻 System Requirements

### Runtime Requirements

- **Operating System**: Windows / Linux / macOS
- **Architecture**: amd64 / arm64
- **Runtime Dependencies**: None (single binary)

### Build Requirements (Only for building from source)

- **Go Version**: ≥ 1.21
- **Dependency Management**: Go Modules

## 📄 License

This project is licensed under the [Apache License 2.0](LICENSE).

```
Copyright 2025 pengcunfu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!

### How to Contribute

1. Fork this repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Code of Conduct

- This tool is intended for **legal security testing** and **API testing** purposes only
- Ensure you have permission to test the target system
- Comply with all applicable laws and regulations
- Respect others' privacy and data security

## ⭐ Star History

If this project helps you, please give it a Star ⭐️

## 📧 Contact

- **Author**: pengcunfu
- **GitHub**: [@pengcunfu](https://github.com/pengcunfu)
- **Repository**: [go-mcp-request](https://github.com/pengcunfu/go-mcp-request)

## 🙏 Acknowledgments

- [MCP Go SDK](https://github.com/mark3labs/mcp-go) - Go implementation of the MCP protocol
- [Logrus](https://github.com/sirupsen/logrus) - Structured logging library

---

<div align="center">

**[⬆ Back to Top](#-mcp-http-request-tool)**

</div>
