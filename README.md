---

# 🧠 Go LLM Tool Calling Runtime (Plugin + Service Discovery)

> A production-grade **LLM Agent Execution Runtime** built with Go
> Supports **tool calling, hot-pluggable plugins, service discovery, and context management**

---

## 🚀 项目简介

本项目是一个 **LLM Agent 执行引擎**，用于：

* 驱动 LLM 与外部工具（Tool）协同工作
* 提供可控、可扩展、可观测的 Agent Runtime
* 支持插件化 Tool + 服务发现（Consul / etcd 等）
* 构建 AI 基础设施级别的执行层

> ⚠️ 本项目不是 ChatBot，而是 AI 系统的 **Infra 层**，类似操作系统内核，用于管理 LLM 执行、工具调度和上下文状态。

---

## 🎯 核心概念

### 1️⃣ LLM

Large Language Model，如：

* GPT-4
* Claude

> **能力**：生成文本、推理、总结
> **限制**：无法访问实时数据、外部系统、数据库或执行操作

---

### 2️⃣ Tool（插件）

> **Tool = LLM 无法完成的外部能力**

常见类型：

| 类型      | 示例            |
| ------- | ------------- |
| 🌐 数据   | search / API  |
| 🗄️ 数据库 | SQL 查询        |
| 🧮 计算   | calculator    |
| 📧 操作   | send email    |
| 🔍 检索   | vector search |

**特点：**

* 独立进程（隔离）
* 可热插拔
* 支持多语言
* 通过 gRPC 插件化调用

---

### 3️⃣ Agent Core

> **Agent Core = 大脑 / 控制循环 / 状态管理**

职责：

* 调用 LLM
* 判断是否调用 Tool
* 管理上下文（State / Messages）
* 决定何时结束

❌ 不负责：

* 执行 Tool
* 调用外部系统

---

### 4️⃣ Executor

> **Executor = 工人 / 执行器**

职责：

* 执行 Tool（通过 RPC）
* 控制 timeout / retry / rate limit

❌ 不负责：

* 决策
* 状态管理
* 调用 LLM

---

### 5️⃣ State / Memory

> **上下文系统 = LLM + Agent 的记忆**

* 保存用户消息 / Assistant 回复 / Tool 执行结果
* 支持 future long-term memory / RAG
* Tool 执行结果必须由 Agent Core 注入

---

### 6️⃣ Tool Registry + Service Discovery

> 动态发现、注册和管理插件

* 使用服务发现（如 Consul）
* 支持工具上线/下线
* 健康检查、隔离失败 Tool

---

## 🔄 系统架构图

```text
                   ┌───────────────────────┐
                   │    Agent Runtime      │
                   │   (Agent Core + LLM) │
                   └─────────┬────────────┘
                             │
                             ▼
                    ┌───────────────────┐
                    │   Tool Registry   │  ← 服务发现（Consul / etcd）
                    └─────────┬─────────┘
                              │
            ┌─────────────────┼─────────────────┐
            │                 │                 │
 ┌──────────▼─────────┐ ┌─────▼──────────┐ ┌────▼──────────┐
 │  Search Tool        │ │ DB Tool       │ │ Python Tool   │
 │ (gRPC Plugin)       │ │ (gRPC Plugin) │ │ (gRPC Plugin) │
 └─────────────────────┘ └───────────────┘ └───────────────┘
```

---

## 🔄 Tool 执行流程

1. 用户输入 → Agent Core
2. Agent Core 调用 LLM → LLM 决策是否调用 Tool
3. 如果需要 Tool → Executor 调用插件
4. Tool 返回结果 → Agent Core 注入 State
5. Agent Core 再次调用 LLM → 生成最终结果

```go
resp := LLM(state)
if resp.ToolCall != nil {
    result := executor.Execute(resp.ToolCall)
    state.AddObservation(resp.ToolCall.Name, result)
} else {
    return resp.Answer
}
```

---

## 📁 详细项目目录

```text
.
├── cmd/
│   └── server/                     # 启动入口
│       ├── main.go                  # Runtime 启动
│       └── config.go                # 配置加载
│
├── internal/
│   ├── agent/                       # Agent Core 模块
│   │   ├── loop.go                  # 主执行循环
│   │   ├── decision.go              # LLM 输出解析 & Tool 调用判断
│   │   ├── planner.go               # 可选：多步任务规划
│   │   └── context.go               # State / Memory 管理
│   │
│   ├── executor/                    # Tool 执行层
│   │   ├── executor.go              # 执行器入口
│   │   ├── runner.go                # 并发 + timeout + retry 控制
│   │   └── concurrency.go           # goroutine 管理
│   │
│   ├── registry/                    # Tool 注册与服务发现
│   │   ├── registry.go              # 内存注册表
│   │   ├── consul.go                # Consul 实现
│   │   └── healthcheck.go           # Tool 健康检查
│   │
│   ├── toolclient/                  # gRPC Tool 客户端
│   │   ├── client.go                # 调用封装
│   │   ├── plugin.go                # 动态注册/卸载插件
│   │   └── request.go               # ToolRequest / ToolResponse 封装
│   │
│   ├── llm/                         # LLM Client
│   │   ├── client.go                # 通用接口
│   │   ├── gpt.go                   # GPT 系列实现
│   │   └── prompt.go                # Prompt 构建工具
│   │
│   ├── state/                       # 上下文 / Memory
│   │   ├── state.go                 # 主结构
│   │   ├── message.go               # Message 类型
│   │   └── memory.go                # 长期记忆 / RAG
│   │
│   └── utils/                       # 工具函数
│       ├── logger.go                # 日志
│       ├── config.go                # 配置加载
│       └── error.go                 # 自定义错误类型
│
├── proto/                           # gRPC / Plugin 接口
│   ├── tool.proto                    # ToolService + Request / Response
│   └── tool.pb.go                    # Go 生成文件
│
├── pkg/
│   └── api/                          # HTTP/gRPC 对外接口
│       ├── server.go                 # 服务入口
│       ├── handler.go                # 请求处理
│       └── middleware.go             # 拦截 / auth / rate limit
│
├── configs/                          # 配置文件
│   ├── config.yaml                   # Runtime 配置
│   ├── tools.yaml                    # Tool 插件配置
│   └── logging.yaml                  # 日志配置
│
├── scripts/                          # 启动 & 测试脚本
│   ├── start.sh                       # 启动 Runtime
│   ├── start_tool.sh                  # 启动单个 Tool
│   └── test_rpc.sh                    # RPC 调试
│
└── README.md
```

---

## 🧭 Roadmap / TODO

### 🟢 基础（MVP）

* [ ] Agent Core + Loop
* [ ] Tool gRPC 调用
* [ ] ToolRegistry（静态配置）
* [ ] State / Memory 管理

### 🟡 插件化（核心）

* [ ] 服务发现（Consul / etcd）
* [ ] Tool 动态注册 / 卸载
* [ ] Tool 健康检查
* [ ] Tool 上下线动态更新

### 🟠 进阶

* [ ] JSON Schema Tool（function calling）
* [ ] 多 Tool 并发执行
* [ ] Streaming（LLM + Tool）
* [ ] Context Window 控制

### 🔴 Infra 级

* [ ] Observability（trace / metrics）
* [ ] Replay / 执行重放
* [ ] Tool 权限控制
* [ ] Rate limit / quota
* [ ] 分布式 Agent Runtime

---

## 🧠 Notes（关键理解 / 自己看的）

1. **Agent Core vs Executor**

   * Core = 决策 + 状态管理
   * Executor = 调用 Tool

2. **Tool 为什么不直接改 State**

   * Tool → 返回结果
   * Agent Core → 注入 State

3. **谁决定调用 Tool**

   * LLM 决定

4. **什么时候不调用 Tool**

   * LLM 自己可以完成任务（总结 / 推理 / 上下文已包含）

5. **Tool 本质**

   * 能力扩展（外部世界接口）

---