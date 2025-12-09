# 计划任务系统

本目录包含后端的计划任务实现。

## 架构

### Scheduler (调度器)

`Scheduler` 是一个通用的任务调度器，用于管理多个计划任务：

- **注册任务**: 通过 `Register()` 方法注册任务
- **启动调度**: 通过 `Start()` 方法启动所有任务
- **停止调度**: 通过 `Stop()` 方法优雅地停止所有任务
- **并发执行**: 每个任务在独立的 goroutine 中运行
- **错误处理**: 任务执行错误会被记录，不会影响其他任务

### Task 接口

所有计划任务都需要实现 `Task` 接口：

```go
type Task interface {
    Name() string              // 任务名称
    Run(ctx context.Context) error  // 执行任务
    Interval() time.Duration   // 执行间隔
}
```

## 现有任务

### FileCleanerTask (文件清理任务)

定期清理过期的文件：

- **名称**: `file-cleaner`
- **间隔**: 1 小时
- **功能**: 
  - 查找超过 `FILE_EXPIRE_TIME` 小时的文件
  - 删除物理文件
  - 删除数据库记录
  - 记录清理统计信息

## 添加新任务

1. 创建新的任务文件，实现 `Task` 接口：

```go
package tasks

import (
    "context"
    "time"
)

type MyTask struct {
    interval time.Duration
}

func NewMyTask() *MyTask {
    return &MyTask{
        interval: 30 * time.Minute,
    }
}

func (m *MyTask) Name() string {
    return "my-task"
}

func (m *MyTask) Interval() time.Duration {
    return m.interval
}

func (m *MyTask) Run(ctx context.Context) error {
    // 实现任务逻辑
    return nil
}
```

2. 在 `main.go` 中注册任务：

```go
myTask := tasks.NewMyTask()
scheduler.Register(myTask)
```

## 配置

任务可以通过环境变量配置：

- `FILE_EXPIRE_TIME`: 文件过期时间（小时），默认 6 小时

## 日志

所有任务执行都会记录日志：

- 任务开始执行
- 任务执行成功/失败
- 执行耗时
- 错误信息

## 优雅关闭

调度器支持优雅关闭：

- 调用 `Stop()` 时，会取消所有任务的 context
- 任务可以检查 `ctx.Done()` 来响应取消信号
- 所有任务完成后才会返回

