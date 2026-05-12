# 运行测试指南

本文档展示如何运行 TwelveDataClient 的测试套件并获取详细输出。

## 基础测试命令

### 1. 运行所有测试 (详细输出)

```bash
go test -v ./test/
```

输出示例:
```
=== RUN   TestWebSocketService_SubscriptionMessageGeneration
    benchmark_test.go:145: 期望 Action 为 "subscribe"，得到 "subscribe"
--- PASS: TestWebSocketService_SubscriptionMessageGeneration (0.00s)
=== RUN   TestPricing_ValidPrice
--- PASS: TestPricing_ValidPrice (0.00s)
```

### 2. 运行所有测试 (简洁输出)

```bash
go test ./test/
```

输出示例:
```
ok      twelve_data_client/test     0.123s
```

---

## 基准测试命令

### 3. 运行所有基准测试

```bash
go test -v -bench=. ./test/
```

输出示例:
```
BenchmarkNewSubscribe_Latency-8                2500000    480 ns/op    200 B/op    5 allocs/op
BenchmarkSubscriptionMarshal_Latency-8         1500000    620 ns/op    240 B/op    6 allocs/op
BenchmarkPriceEventUnmarshal_Latency-8         1200000    890 ns/op    88 B/op     2 allocs/op
```

### 4. 运行所有基准测试 (带内存报告)

```bash
go test -v -bench=. -benchmem ./test/
```

输出示例:
```
BenchmarkNewSubscribe_Latency-8                2500000    480 ns/op    200 B/op    5 allocs/op
BenchmarkSubscriptionMarshal_Latency-8         1500000    620 ns/op    240 B/op    6 allocs/op
```

### 5. 运行特定基准测试

```bash
# 只运行延迟基准
go test -v -bench=Latency ./test/

# 只运行性能基准
go test -v -bench=Performance ./test/

# 只运行价格基准
go test -v -bench=Pricing ./test/

# 只运行并发基准
go test -v -bench=Concurrent ./test/

# 只运行 WebSocket 服务基准
go test -v -bench=WebSocketService ./test/
```

### 6. 运行特定基准并设置时间

```bash
# 运行时间加倍以获得更准确的结果
go test -v -bench=. -benchtime=5s ./test/

# 运行次数指定
go test -v -bench=. -benchtime=10000x ./test/
```

---

## 综合全面的测试输出

### 7. 完整的综合测试报告

```bash
go test -v -bench=. -benchmem -run=. -timeout=30s ./test/ 2>&1 | tee test_report.txt
```

这会:
- 运行所有测试和基准
- 显示内存分配信息
- 生成 test_report.txt 文件保存输出
- 同时在终端显示

### 8. 测试覆盖率报告

```bash
go test -v -cover ./test/
```

输出示例:
```
coverage: 78.5% of statements
ok      twelve_data_client/test     0.234s
```

### 9. 详细的覆盖率文件

```bash
go test -coverprofile=coverage.out ./test/
go tool cover -html=coverage.out -o coverage.html
```

然后用浏览器打开 coverage.html 查看可视化覆盖率。

---

## 分类测试运行

### 10. 只运行单元测试 (不运行基准)

```bash
go test -v -run=^Test ./test/
```

### 11. 只运行基准 (不运行单元测试)

```bash
go test -v -bench=. -run=^$ ./test/
```

### 12. 只运行特定类别的测试

```bash
# 只运行 WebSocket 服务测试
go test -v -run=WebSocketService ./test/

# 只运行价格数据测试
go test -v -run=Pricing ./test/

# 只运行错误处理测试
go test -v -run=Error ./test/

# 只运行订阅响应测试
go test -v -run=SubscriptionResponse ./test/
```

---

## 高级输出选项

### 13. CPU 和内存分析

```bash
go test -v -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof ./test/
go tool pprof -text cpu.prof > cpu_analysis.txt
go tool pprof -text mem.prof > mem_analysis.txt
```

### 14. 比较两次基准测试结果

```bash
# 运行第一次基准
go test -bench=. -benchmem ./test/ > bench1.txt

# 修改代码...

# 运行第二次基准
go test -bench=. -benchmem ./test/ > bench2.txt

# 比较结果
benchstat bench1.txt bench2.txt
```

### 15. 追踪信息输出

```bash
go test -v -trace=trace.out ./test/
go tool trace trace.out
```

---

## 脚本化测试运行

### 16. 完整的测试脚本 (Linux/macOS)

```bash
#!/bin/bash

echo "========================================="
echo "TwelveData WebSocket 客户端 - 完整测试"
echo "========================================="
echo ""

echo "运行时间: $(date)"
echo ""

echo "1. 运行所有单元测试..."
go test -v -run=^Test ./test/ 2>&1 | tee test_unit.log
echo ""

echo "2. 运行所有基准测试..."
go test -v -bench=. -benchmem ./test/ 2>&1 | tee test_bench.log
echo ""

echo "3. 生成覆盖率报告..."
go test -coverprofile=coverage.out ./test/
go tool cover -html=coverage.out -o coverage.html
echo "覆盖率报告已生成: coverage.html"
echo ""

echo "4. 生成测试统计..."
echo ""
echo "总体结果统计:"
grep -E "^(ok|FAIL)" test_unit.log test_bench.log
echo ""

echo "========================================="
echo "测试完成!"
echo "========================================="
```

保存为 `run_tests.sh` 并运行:
```bash
chmod +x run_tests.sh
./run_tests.sh
```

### 17. 完整的测试脚本 (Windows PowerShell)

```powershell
# run_tests.ps1

Write-Host "=========================================" -ForegroundColor Green
Write-Host "TwelveData WebSocket 客户端 - 完整测试" -ForegroundColor Green
Write-Host "=========================================" -ForegroundColor Green
Write-Host ""

$timestamp = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
Write-Host "运行时间: $timestamp"
Write-Host ""

Write-Host "1. 运行所有单元测试..." -ForegroundColor Cyan
go test -v -run=^Test ./test/ 2>&1 | Tee-Object -FilePath test_unit.log
Write-Host ""

Write-Host "2. 运行所有基准测试..." -ForegroundColor Cyan
go test -v -bench=. -benchmem ./test/ 2>&1 | Tee-Object -FilePath test_bench.log
Write-Host ""

Write-Host "3. 生成覆盖率报告..." -ForegroundColor Cyan
go test -coverprofile=coverage.out ./test/
go tool cover -html=coverage.out -o coverage.html
Write-Host "覆盖率报告已生成: coverage.html" -ForegroundColor Yellow
Write-Host ""

Write-Host "=========================================" -ForegroundColor Green
Write-Host "测试完成!" -ForegroundColor Green
Write-Host "=========================================" -ForegroundColor Green
```

运行:
```powershell
.\run_tests.ps1
```

---

## 输出说明

### 测试输出格式

```
=== RUN   TestName
    file.go:123: 测试日志输出
--- PASS: TestName (0.00s)
```

| 字段 | 含义 |
|-----|------|
| RUN | 测试开始 |
| PASS | 测试通过 |
| FAIL | 测试失败 |
| SKIP | 测试跳过 |
| 0.00s | 测试耗时 |

### 基准输出格式

```
BenchmarkName-8    2500000    480 ns/op    200 B/op    5 allocs/op
└────┬──────┘      └───┬───┘  └──┬──┘      └─┬──┘      └──┬───┘
  测试名称        迭代次数   单次耗时    内存分配量    分配次数
```

---

## 常见用例

### 快速检查 (5 秒)

```bash
go test -v ./test/
```

### 详细基准 (1 分钟)

```bash
go test -v -bench=. -benchmem -benchtime=10s ./test/
```

### 完整报告 (2 分钟)

```bash
go test -v -bench=. -benchmem -cover -timeout=120s ./test/ 2>&1 | tee full_report.txt
```

### 持续集成 (CI/CD)

```bash
go test -v -race -cover -timeout=60s ./test/
```

---

## 解读测试结果

### 成功的测试输出

```
ok      twelve_data_client/test     0.123s  coverage: 85.2%
```

表示:
- 所有测试通过
- 耗时 0.123 秒
- 代码覆盖率 85.2%

### 失败的测试输出

```
FAIL    twelve_data_client/test     0.456s
```

表示:
- 存在失败的测试
- 耗时 0.456 秒
- 需要查看详细的失败信息

### 基准性能指标

| 指标 | 含义 | 预期值 |
|-----|------|--------|
| ns/op | 纳秒/操作 | 越小越好 |
| B/op | 字节/操作 | 越小越好 |
| allocs/op | 分配次数/操作 | 越少越好 |

---

## 最推荐的命令

### 日常开发

```bash
go test -v ./test/
```

### 提交前检查

```bash
go test -v -bench=. -benchmem -cover ./test/
```

### 性能优化

```bash
go test -v -bench=. -benchmem -benchtime=10s -cpuprofile=cpu.prof ./test/
go tool pprof -http=:8080 cpu.prof
```

### CI/CD 流程

```bash
go test -v -race -cover -timeout=60s ./test/ && \
go test -v -bench=. -benchmem ./test/
```
