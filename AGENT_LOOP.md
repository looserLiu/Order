# 常驻循环 Agent 协议 (Resident Loop Agent Protocol)

> 本文件由 Roo（常驻 agent 大脑）驱动。后台守护脚本 `scripts/loop-agent.sh`
> 负责心跳、静态检查与自动提交（单实例，pidfile 锁，已修复 fork 炸弹）；
> 真正的"重构/优化"决策与代码改动由 Roo 在会话中执行，并写入 `LOOP_LOG.md`。
> 若未来环境出现可用的编码引擎（如 `roo`/`claude` CLI 或 `ollama`），可通过 `ENGINE_CMD` 接管自主编码。

## 循环定义 (The Loop)

```
while now < DEADLINE:
    1. SCAN      扫描项目，识别下一个改进点（读代码 / 静态检查）
    2. REFACTOR  先完成待办中的重构项（架构/可维护性）
    3. OPTIMIZE  再持续做优化（性能/质量/去重/测试）
    4. VERIFY    前端: tsc --noEmit + eslint；后端: 精读 + (若 go 可用) go build/vet
    5. COMMIT    仅当验证通过才提交，并追加 LOOP_LOG.md
    6. LOG       记录本轮做了什么、下一步计划
```

- **截止时间 DEADLINE**: 2026-07-10 09:00 (Asia/Shanghai, UTC+8)
- **无需人工审批**: 所有改动由 agent 自主决定并提交。
- **安全原则**: 后端无编译器时只做保守/增量改动并精读确认；前端改动必须经 tsc/eslint 验证。

## 环境约束 (已探明)

- `go` 编译器：当前环境缺失（标准路径与 Spotlight 均未找到）→ 后端改动无法编译验证，须保守。
- `node` v18.15.0 可用；前端 `node_modules` 已装 → 前端可 `tsc --noEmit` / `eslint` 验证。
- `ollama` 不可用 → 无本地 LLM 自主引擎；Roo 会话即大脑。
- 守护进程：单实例常驻，心跳 + 前端检查 + 自动提交已验证改动。

## 待办清单 (Backlog)

状态: [ ] 待办 [-] 进行中 [x] 完成

### 阶段一：重构 (Refactor)

- [x] INFRA 初始化 git 仓库与 .gitignore，建立循环 harness（AGENT_LOOP/LOOP_LOG/loop-agent.sh）
- [x] FIX 修复守护脚本 fork 炸弹，确保单实例常驻
- [ ] R1 后端：将 3013 行单体 `handlers.go` 按资源拆分为同包多文件（auth/accounts/transactions/...）
- [ ] R2 后端：建立 `internal/repositories` 仓储层（接口 + GORM 实现），核心实体
- [ ] R3 后端：建立 `internal/services` 服务层，承载业务逻辑（从 handler 抽出）
- [ ] R4 后端：统一错误处理与响应封装（复用 `pkg/response`），消除重复样板
- [ ] R5 后端：请求 DTO 与校验集中化（复用 `pkg/validator`）
- [ ] R6 前端：API 客户端集中化与一致的错误/加载态处理（`services/api.ts`）
- [ ] R7 前端：抽取跨页面共享组件（表格/表单/弹窗/空态），消除重复

### 阶段二：优化 (Optimize) — 无限循环

- [ ] O1 后端：DB 查询优化（预加载、索引、消除 N+1、分页）
- [ ] O2 后端：事务边界与并发安全（批量操作、余额更新）
- [ ] O3 前端：路由级懒加载与代码分割（减少首屏体积）
- [ ] O4 前端：列表/计算记忆化（React.memo / useMemo / 虚拟列表）
- [ ] O5 前端：i18n 完整性与主题一致性
- [ ] O6 测试：后端 service 单测；前端关键组件/工具单测
- [ ] O7 文档：更新 ARCHITECTURE.md 反映真实实现
- [ ] O8 持续：扫描新发现的坏味道并修复（死代码、重复、魔法值、低效算法）

## 如何恢复 (Resume)

1. 查看 `LOOP_LOG.md` 最新一条，了解进度与下一步。
2. 运行 `bash scripts/loop-agent.sh` 重启守护（心跳 + 检查 + 自动提交）。
3. Roo 会话中继续按本文件 Backlog 推进。
