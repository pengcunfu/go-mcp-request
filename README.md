# MCP HTTP 请求工具 (Go版本)

为API测试、Web自动化和安全测试设计的全功能HTTP客户端MCP（模型上下文协议）服务器，具备完整的HTTP工具和详细的日志记录功能。使用Go语言编写，性能更好，部署更简单。

## 特性

- **完整的HTTP方法支持**: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
- **高级安全测试**: 原始请求工具，用于渗透测试、SQL注入、XSS测试
- **全参数支持**: 所有方法支持Headers、Cookies、Body、超时设置
- **自动日志记录**: 所有请求和响应自动记录到 `~/mcp_requests_logs/`
- **精确保证**: 原始模式完全保留每个字符
- **MCP兼容**: 兼容Claude Code、Cursor和其他MCP客户端

## 安装

### 选项1：下载二进制文件
从 [releases 页面](https://github.com/pengcunfu/go-mcp-request/releases) 下载最新的二进制文件。

### 选项2：从源码构建
```bash
git clone https://github.com/pengcunfu/go-mcp-request.git
cd go-mcp-request
go build -o mcp-request main.go
```

## 使用方法

### 在 Cursor/Claude Code 中使用

添加到你的MCP配置文件 (`~/.cursor/mcp_servers.json` 或类似文件):

```json
{
  "mcpServers": {
    "mcp-request": {
      "command": "./mcp-request",
      "type": "stdio"
    }
  }
}
```

或者如果二进制文件在你的PATH中:

```json
{
  "mcpServers": {
    "mcp-request": {
      "command": "mcp-request",
      "type": "stdio"
    }
  }
}
```

### 可用工具

1. **http_get** - 全功能GET请求
2. **http_post** - 全功能POST请求
3. **http_put** - 全功能PUT请求
4. **http_delete** - 全功能DELETE请求
5. **http_patch** - 全功能PATCH请求
6. **http_head** - 全功能HEAD请求
7. **http_options** - 全功能OPTIONS请求
8. **http_raw_request** - 用于安全测试的原始HTTP请求

### 使用示例

```bash
# 基础GET请求
http_get("https://api.example.com/users")

# 带数据和请求头的POST请求
http_post(
    url="https://api.example.com/login",
    body='{"username":"test","password":"test"}',
    headers={"Content-Type": "application/json"}
)

# 使用原始请求进行安全测试
http_raw_request(
    url="https://vulnerable-site.com/search",
    method="POST", 
    raw_body="q=test' OR 1=1--",
    headers={"Content-Type": "application/x-www-form-urlencoded"}
)
```

## 安全测试特性

`http_raw_request` 工具专为安全测试设计：

- **绝对精确**: 每个字符完全保留
- **无编码**: 特殊字符 (', ", \\, %, &, =) 原样发送
- **完整请求头**: 不截断长cookies或tokens
- **原始载荷**: 完美适用于SQL注入、XSS、CSRF测试

## 日志记录

所有HTTP请求和响应自动记录到：

- **位置**: `~/mcp_requests_logs/`
- **格式**: JSON格式，包含时间戳和完整的请求/响应详情
- **文件名**: `requests_YYYYMMDD_HHMMSS.log`

查看日志：
```bash
tail -f ~/mcp_requests_logs/requests_*.log
```

## 系统要求

- Go ≥ 1.21 (从源码构建时需要)
- 无运行时依赖，单一二进制文件

## 许可证

MIT License

## 贡献

欢迎贡献！此工具仅用于防御性安全测试和合法的API测试目的。