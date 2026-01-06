# Operator Hub

[English](README.md)

Operator Hub 是一个开源的算子与工具管理平台，旨在连接大语言模型（LLMs）与实际业务能力。通过支持 [Model Context Protocol (MCP)](https://modelcontextprotocol.io/)，它提供了一套标准化的机制来注册、管理和执行各种算子（Operators）及工具（Tools），帮助开发者快速构建强大的 AI Agent 应用。

## 核心组件

本项目包含两个主要组件：

### 1. Operator Integration (`operator-integration`)
核心集成服务平台，负责算子和工具的全生命周期管理。
- **算子管理**：支持算子的注册、版本控制、发布与下架。
- **工具箱（Toolbox）**：支持将多个工具组合成工具箱，便于统一管理和调用。
- **MCP 支持**：作为 MCP Server，向 LLM 提供标准化的工具调用接口。
- **多协议适配**：支持 HTTP、SSE 等多种通信协议。
- **权限控制**：内置基于策略的访问控制机制。

### 2. Operator App (`operator-app`)
应用端运行时与示例实现。
- 提供了一个轻量级的算子执行环境。
- 展示了如何集成和使用 Operator Hub 的核心能力。
- 包含 MCP 客户端与服务端的交互示例。

## 特性

- **标准化接口**：基于 MCP 协议，实现模型与工具的解耦。
- **灵活扩展**：支持多种编程语言编写的算子（如 Go, Python 等）。
- **可观测性**：集成了 OpenTelemetry，提供全链路追踪能力。
- **高性能**：基于 Go 语言开发，具备高并发处理能力。

## 快速开始

### 前置要求
- Go 1.23+
- MySQL / MariaDB / Dameng DB
- Redis

### 编译与运行

#### 运行 Operator Integration
```bash
cd operator-integration
# 安装依赖
go mod tidy
# 编译
go build -o operator-integration server/main.go
# 运行 (需配置相应的配置文件)
./operator-integration
```

#### 运行 Operator App
```bash
cd operator-app
# 安装依赖
go mod tidy
# 编译
go build -o operator-app server/main.go
# 运行
./operator-app
```

## 贡献

欢迎提交 Pull Request 或 Issue！

## 许可证

[Apache-2.0](LICENSE)
