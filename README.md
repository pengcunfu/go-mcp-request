<div align="center">

# 🚀 MCP HTTP 请求工具

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-%3E%3D1.21-00ADD8?logo=go)](https://go.dev/)
[![GitHub release](https://img.shields.io/github/v/release/pengcunfu/go-mcp-request)](https://github.com/pengcunfu/go-mcp-request/releases)
[![GitHub stars](https://img.shields.io/github/stars/pengcunfu/go-mcp-request)](https://github.com/pengcunfu/go-mcp-request/stargazers)

[English](./README_EN.md) | 简体中文

为API测试、Web自动化和安全测试设计的全功能HTTP客户端MCP（模型上下文协议）服务器。

具备完整的HTTP工具和详细的日志记录功能，使用Go语言编写，性能卓越，部署简单。

</div>

---

## ✨ 特性

- 🔧 **完整的HTTP方法支持** - GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
- 🔒 **高级安全测试** - 原始请求工具，支持渗透测试、SQL注入、XSS测试
- ⚙️ **全参数支持** - 所有方法支持Headers、Cookies、Body、超时设置
- 📝 **自动日志记录** - 所有请求和响应自动记录到 `~/mcp_requests_logs/`
- 🎯 **精确保证** - 原始模式完全保留每个字符，无自动编码
- 🔌 **MCP兼容** - 完美兼容Claude Desktop、Cursor和其他MCP客户端
- ⚡ **高性能** - Go语言实现，单一二进制文件，无运行时依赖
- 🌍 **跨平台** - 支持Windows、Linux、macOS

## 📦 安装

### 方式一：下载预编译二进制文件（推荐）

从 [Releases 页面](https://github.com/pengcunfu/go-mcp-request/releases) 下载适合你操作系统的最新版本：

- **Windows**: `mcp-request-windows-amd64.exe`
- **Linux**: `mcp-request-linux-amd64`
- **macOS**: `mcp-request-darwin-amd64`

### 方式二：从源码构建

```bash
# 克隆仓库
git clone https://github.com/pengcunfu/go-mcp-request.git
cd go-mcp-request

# 构建
go build -o mcp-request main.go

# 或使用 go install
go install github.com/pengcunfu/go-mcp-request@latest
```

## 🚀 使用方法

### 在 Claude Desktop 中配置

编辑配置文件 `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) 或 `%APPDATA%\Claude\claude_desktop_config.json` (Windows):

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

### 在 Cursor 中配置

编辑配置文件 `~/.cursor/mcp_servers.json`:

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

> **提示**: 将 `/path/to/mcp-request` 替换为实际的二进制文件路径。如果已添加到PATH环境变量，可直接使用 `mcp-request`。

## 🛠️ 可用工具

| 工具名称 | 描述 | 适用场景 |
|---------|------|----------|
| `http_get` | 全功能GET请求 | 获取资源、API查询 |
| `http_post` | 全功能POST请求 | 创建资源、表单提交 |
| `http_put` | 全功能PUT请求 | 更新资源 |
| `http_delete` | 全功能DELETE请求 | 删除资源 |
| `http_patch` | 全功能PATCH请求 | 部分更新资源 |
| `http_head` | 全功能HEAD请求 | 获取响应头 |
| `http_options` | 全功能OPTIONS请求 | 检查支持的方法 |
| `http_raw_request` | 原始HTTP请求 | 安全测试、渗透测试 |

## 📖 使用示例

### 基础GET请求

```python
# 简单GET请求
http_get("https://api.example.com/users")

# 带请求头的GET请求
http_get(
    url="https://api.example.com/users",
    headers={"Authorization": "Bearer token123"}
)
```

### POST请求

```python
# JSON数据POST
http_post(
    url="https://api.example.com/login",
    body='{"username":"test","password":"test"}',
    headers={"Content-Type": "application/json"}
)

# 表单数据POST
http_post(
    url="https://api.example.com/form",
    body="name=John&age=30",
    headers={"Content-Type": "application/x-www-form-urlencoded"}
)
```

### 安全测试（原始请求）

```python
# SQL注入测试
http_raw_request(
    url="https://test-site.com/search",
    method="POST",
    raw_body="q=test' OR 1=1--",
    headers={"Content-Type": "application/x-www-form-urlencoded"}
)

# XSS测试
http_raw_request(
    url="https://test-site.com/comment",
    method="POST",
    raw_body='comment=<script>alert("XSS")</script>',
    headers={"Content-Type": "application/x-www-form-urlencoded"}
)
```

## 🔒 安全测试特性

`http_raw_request` 工具专为安全研究和渗透测试设计：

### 核心优势

- ✅ **绝对精确** - 每个字节完全保留，不做任何修改
- ✅ **无自动编码** - 特殊字符 (`'`, `"`, `\`, `%`, `&`, `=`) 原样发送
- ✅ **完整请求头** - 不截断长cookies或tokens
- ✅ **原始载荷** - 完美适用于各类安全测试

### 适用测试场景

- 🎯 SQL注入测试
- 🎯 XSS跨站脚本测试
- 🎯 CSRF跨站请求伪造测试
- 🎯 命令注入测试
- 🎯 HTTP请求走私测试
- 🎯 自定义协议测试

> ⚠️ **免责声明**: 此工具仅供合法的安全测试和研究使用。请确保你有权限测试目标系统。未经授权的测试可能违反法律。

## 📝 日志记录

所有HTTP请求和响应自动记录，便于调试和审计。

### 日志配置

- **存储位置**: `~/mcp_requests_logs/`
- **文件格式**: JSON格式，结构化存储
- **文件命名**: `requests_YYYYMMDD_HHMMSS.log`
- **包含信息**:
  - 时间戳
  - 完整的请求信息（方法、URL、请求头、请求体）
  - 完整的响应信息（状态码、响应头、响应体）
  - 错误信息（如有）

### 查看日志

```bash
# 实时查看最新日志
tail -f ~/mcp_requests_logs/requests_*.log

# 查看特定日期的日志
cat ~/mcp_requests_logs/requests_20250111_*.log

# 使用jq格式化查看
cat ~/mcp_requests_logs/requests_*.log | jq .
```

## 💻 系统要求

### 运行要求

- **操作系统**: Windows / Linux / macOS
- **架构**: amd64 / arm64
- **运行时依赖**: 无（单一二进制文件）

### 构建要求（仅从源码构建时）

- **Go版本**: ≥ 1.21
- **依赖管理**: Go Modules

## 📄 许可证

本项目采用 [Apache License 2.0](LICENSE) 开源协议。

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

## 🤝 贡献

欢迎贡献代码、报告问题或提出新功能建议！

### 如何贡献

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

### 行为准则

- 此工具仅用于**合法的安全测试**和**API测试**目的
- 请确保你有权限测试目标系统
- 遵守所有适用的法律法规
- 尊重他人的隐私和数据安全

## ⭐ Star History

如果这个项目对你有帮助，请给个 Star ⭐️

## 📧 联系方式

- **作者**: pengcunfu
- **GitHub**: [@pengcunfu](https://github.com/pengcunfu)
- **仓库**: [go-mcp-request](https://github.com/pengcunfu/go-mcp-request)

## 🙏 致谢

- [MCP Go SDK](https://github.com/mark3labs/mcp-go) - MCP协议的Go实现
- [Logrus](https://github.com/sirupsen/logrus) - 结构化日志库

---

<div align="center">

**[⬆ 回到顶部](#-mcp-http-请求工具)**

</div>
