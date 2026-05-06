# 开发工作流程 (Development Workflow)

## Order App 标准开发流程

---

## 一、完整开发流程图

```
┌─────────────────────────────────────────────────────────────────┐
│                        Sprint 周期 (2周)                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Day 1                Day 2-9              Day 10              │
│    │                      │                    │                 │
│    ▼                      ▼                    ▼                 │
│ ┌──────────┐        ┌──────────┐         ┌──────────┐           │
│ │ Planning │        │  开发    │         │ Review   │           │
│ │  计划会  │        │  +      │         │  +       │           │
│ │          │   ►    │  测试   │    ►    │  Retro   │           │
│ │ - 选任务 │        │          │         │  回顾    │           │
│ │ - 估点  │        │ - 开发   │         │          │           │
│ │ - 分配  │        │ - Code   │         │ - 演示   │           │
│ └──────────┘        │   Review │         │ - 反思   │           │
│                     │ - 测试   │         │ - 改进   │           │
│                     └──────────┘         └──────────┘           │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    Daily Standup (每天)                  │    │
│  │              10:00 - 15分钟站立会议                      │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 二、任务生命周期

```
待办 (Todo) ──► 进行中 (In Progress) ──► 待 Review ──► 已完成 (Done)
    │                  │                    │              │
    │                  │                    │              │
    ▼                  ▼                    ▼              ▼
 规划中            开发中               Code Review       验收通过
```

### 状态定义

| 状态 | 说明 | 操作人 |
|------|------|--------|
| Todo | 已规划，待开始 | PM |
| In Progress | 正在开发 | Dev |
| In Review | 已完成，待 Review | Dev |
| Tested | 测试通过 | QA |
| Done | 已验收完成 | PM/QA |

---

## 三、Git 工作流程

### 3.1 分支策略

```
main (生产环境 - 受保护)
  │
  └── develop (开发分支 - 受保护)
        │
        ├── sprint/1-mvp-account (Sprint 1)
        │     ├── sprint/1-account-crud
        │     └── sprint/1-transaction-crud
        │
        ├── sprint/2-data-management (Sprint 2)
        │     ├── sprint/2-reports
        │     └── sprint/2-search
        │
        └── ... 更多 Sprint 分支
```

### 3.2 开发步骤

```bash
# 1. 从 develop 创建功能分支
git checkout develop
git pull origin develop
git checkout -b sprint/1-account-crud

# 2. 开发功能
git add .
git commit -m "feat: 添加账户 CRUD 功能"

# 3. 推送分支
git push origin sprint/1-account-crud

# 4. 创建 Pull Request
# 在 GitHub/GitLab 创建 PR，目标是 develop

# 5. Code Review 通过后合并
# PR 被 Reviewer 批准后，合并到 develop

# 6. 删除分支
git branch -d sprint/1-account-crud
git push origin --delete sprint/1-account-crud
```

### 3.3 PR 模板

```markdown
## 描述
[简要描述这个 PR 的内容]

## 关联任务
- [ ] Task ID: S1-T1

## 改动内容
- [ ] 新增功能
- [ ] 修改功能
- [ ] 修复 Bug

## 测试情况
- [ ] 单元测试通过
- [ ] 功能测试通过
- [ ] 手动测试通过

## 截图 (如有 UI 变更)
[添加截图]

## Checklist
- [ ] 代码遵循规范
- [ ] 已添加必要的注释
- [ ] 没有 TODO 或 FIXME
- [ ] 相关文档已更新
```

---

## 四、代码审查清单

### 4.1 代码审查要点

| 类别 | 检查项 |
|------|--------|
| **功能** | 实现是否符合需求？逻辑是否正确？ |
| **可读性** | 命名是否清晰？代码是否简洁？ |
| **可维护性** | 是否重复代码？是否过度设计？ |
| **性能** | 是否有性能问题？ |
| **安全** | 是否有安全漏洞？ |
| **测试** | 是否有必要的测试？测试是否充分？ |

### 4.2 审查通过标准

- ✅ 至少 1 人 Approved
- ✅ 所有讨论已解决
- ✅ CI/CD 流水线通过
- ✅ 无 P0/P1 Bug

---

## 五、测试策略

### 5.1 测试金字塔

```
        ┌─────────┐
        │   E2E   │   少量端到端测试
        │   Test  │
        ├─────────┤
        │ Integration│  中等集成测试
        │   Test    │
        ├───────────┤
        │   Unit    │   大量单元测试
        │   Test    │
        └───────────┘
```

### 5.2 测试覆盖率目标

| 层级 | 覆盖率目标 |
|------|-----------|
| Unit Test | 80%+ |
| Integration Test | 60%+ |
| E2E Test | 20%+ |

### 5.3 测试用例示例

```dart
// 账户余额计算测试
test('支出应减少账户余额', () {
  final account = Account(balance: 1000.0);
  final transaction = Transaction(
    accountId: account.id,
    amount: 100.0,
    type: 'expense',
  );

  account.applyTransaction(transaction);

  expect(account.balance, 900.0);
});
```

---

## 六、发布流程

### 6.1 发布检查清单

```markdown
## 发布前检查

### 代码
- [ ] 所有功能已合并到 develop
- [ ] develop 分支已通过所有测试
- [ ] 无已知 P0/P1 Bug
- [ ] 版本号已更新

### 测试
- [ ] 单元测试 100% 通过
- [ ] 集成测试 100% 通过
- [ ] E2E 测试完成
- [ ] 性能测试通过

### 文档
- [ ] CHANGELOG 已更新
- [ ] 用户文档已更新
- [ ] 部署文档已更新

### 沟通
- [ ] 利益相关者已通知
- [ ] 支持团队已培训
- [ ] 回滚计划已准备
```

### 6.2 发布步骤

```bash
# 1. 创建发布分支
git checkout develop
git pull origin develop
git checkout -b release/v0.1

# 2. 更新版本号
flutter pub run set_version 0.1.0

# 3. 合并到 main
git checkout main
git merge release/v0.1
git tag -a v0.1 -m "Release v0.1 MVP"
git push origin main
git push origin v0.1

# 4. 删除发布分支
git branch -d release/v0.1
```

---

## 七、问题处理

### 7.1 Bug 报告模板

```markdown
## Bug 标题
[简洁描述 Bug]

## 环境
- App 版本：
- 设备：
- 操作系统：

## 复现步骤
1. 打开 App
2. 进入 XXX 页面
3. 点击 XXX
4. 出现 XXX 问题

## 预期行为
[描述预期应该发生什么]

## 实际行为
[描述实际发生了什么]

## 截图/日志
[附加相关截图或日志]

## 严重程度
- [ ] P0 - App 崩溃
- [ ] P1 - 核心功能不可用
- [ ] P2 - 功能异常
- [ ] P3 - 轻微问题
```

### 7.2 Hotfix 流程

```bash
# 1. 从 main 创建 hotfix 分支
git checkout main
git pull origin main
git checkout -b hotfix/fix-login-crash

# 2. 修复并测试
# ... 修复代码 ...

# 3. 合并回 main
git checkout main
git merge hotfix/fix-login-crash
git tag -a v0.1.1 -m "Hotfix v0.1.1"
git push origin main
git push origin v0.1.1

# 4. 合并回 develop
git checkout develop
git merge hotfix/fix-login-crash

# 5. 删除 hotfix 分支
git branch -d hotfix/fix-login-crash
```

---

## 八、持续集成/部署

### 8.1 CI/CD 流水线

```yaml
# .github/workflows/flutter.yml
name: Flutter CI

on:
  push:
    branches: [develop, main]
  pull_request:
    branches: [develop]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: subosito/flutter-action@v1
      - run: flutter pub get
      - run: flutter analyze
      - run: flutter test

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: subosito/flutter-action@v1
      - run: flutter pub get
      - run: flutter build apk --release
```

### 8.2 部署环境

| 环境 | 分支 | 用途 |
|------|------|------|
| Development | develop | 开发测试 |
| Staging | release/* | 预发布测试 |
| Production | main | 正式环境 |

---

## 九、度量指标

### 9.1 流程度量

| 指标 | 目标 | 测量方式 |
|------|------|----------|
| Sprint 燃尽率 | 按计划完成 | 每日更新燃尽图 |
| 代码审查时长 | < 24 小时 | PR 创建到合并时间 |
| Bug 逃逸率 | < 5% | 生产 Bug / 测试发现 Bug |
| 部署频率 | 每 Sprint 1 次 | 发布记录 |

### 9.2 质量度量

| 指标 | 目标 | 测量方式 |
|------|------|----------|
| 测试覆盖率 | > 80% | CI 工具报告 |
| Code Review 通过率 | > 90% | PR 统计 |
| 技术债务 | < 5% | 代码扫描 |
