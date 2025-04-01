# Instrument Trace

这是一个用于追踪Go程序函数调用的工具，支持自动注入追踪代码和手动添加追踪代码两种方式。通过该工具，你可以清晰地了解程序的函数调用链和执行流程，对于调试和理解复杂程序特别有帮助。

## 项目架构

```
instrument_trace/
├── Makefile              # 项目构建脚本
├── cmd/
│   └── instrument/
│       └── main.go      # 命令行工具入口
├── example_test.go      # 示例测试
├── go.mod              # Go模块定义
├── instrumenter/       # 自动注入相关代码
│   ├── ast/
│   │   └── ast.go     # AST处理相关代码
│   └── instrumenter.go # 注入器实现
└── trace.go           # 核心追踪实现
```

## 快速开始

1. 克隆项目
```bash
git clone https://github.com/your-username/instrument_trace.git
cd instrument_trace
```

2. 构建项目
```bash
make build
```

3. 运行测试
```bash
make test
```

## 使用示例

### 1. 手动使用追踪

```go
package main

import "github.com/secret-deus/instrument_trace"

func foo() {
    defer trace.Trace()()
    bar()
}

func bar() {
    defer trace.Trace()()
    // 你的代码
}

func main() {
    foo()
}
```

输出示例：
```
g[00001]:    ->main.foo
g[00001]:        ->main.bar
g[00001]:        <-main.bar
g[00001]:    <-main.foo
```

### 2. 使用命令行工具注入追踪代码

1. 对单个文件注入追踪代码：
```bash
./bin/instrument -w main.go
```

2. 对特定包注入追踪代码：
```bash
./bin/instrument -w -pkg mypackage ./mypackage
```

3. 预览注入结果（不写入文件）：
```bash
./bin/instrument main.go
```

### 3. 并发场景示例

```go
package main

import (
    "sync"
    "github.com/secret-deus/instrument_trace"
)

func worker() {
    defer trace.Trace()()
    // 工作代码
}

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            worker()
        }()
    }
    wg.Wait()
}
```

输出示例：
```
g[00002]:    ->main.worker
g[00003]:    ->main.worker
g[00004]:    ->main.worker
g[00002]:    <-main.worker
g[00003]:    <-main.worker
g[00004]:    <-main.worker
```

## 命令行工具参数

```bash
instrument [-w] [-pkg package] [path...]
```

参数说明：
- `-w`：将修改写入源文件（默认输出到标准输出）
- `-pkg`：指定要处理的包名（可选）
- `path`：要处理的文件或目录路径

## 特性

- 支持函数调用追踪
- 自动注入追踪代码
- 支持处理单个文件或整个目录
- 支持选择性处理特定包
- 自动添加必要的导入语句
- 保留原有的注释和格式
- 支持并发场景
- 线程安全

## 注意事项

1. 自动注入时会跳过以下情况：
   - main 函数
   - 已经包含 `defer trace.Trace()()` 的函数
   - 不匹配指定包名的文件

2. 建议在开发环境中使用，不要在生产环境中启用追踪

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License 
