# Sprint 0: 准备阶段 - 详细任务

**时间:** Week 1-2
**目标:** 完成开发环境搭建和架构设计

---

## 任务详情

### S0-T1: Flutter 项目初始化

**负责人:** Dev1 + Dev2
**预估:** 2 人天
**依赖:** 无

**任务内容:**
- [ ] 安装 Flutter SDK 3.x
- [ ] 安装必要插件 (Flutter DevTools, Dart)
- [ ] 创建 Flutter 项目 `flutter create order_app`
- [ ] 配置项目名称、包名
- [ ] 创建目录结构:
  ```
  lib/
    core/
      constants/
      utils/
      theme/
      widgets/
    data/
      database/
      models/
      repositories/
    domain/
      entities/
      repositories/
      usecases/
    presentation/
      pages/
      widgets/
      providers/
  ```
- [ ] 添加基础依赖到 pubspec.yaml:
  - sqflite
  - path_provider
  - provider
  - intl
  - fl_chart
  - csv

**验收标准:**
- `flutter doctor` 无错误
- `flutter run` 可以启动
- 目录结构符合架构要求

---

### S0-T2: SQLite 数据库配置

**负责人:** Dev1
**预估:** 1 人天
**依赖:** S0-T1

**任务内容:**
- [ ] 创建 `database_helper.dart`
- [ ] 配置数据库版本管理
- [ ] 创建表结构:
  - accounts 表
  - categories 表
  - transactions 表
  - products 表
  - warehouses 表
  - inventory_flows 表
  - budgets 表
  - cost_categories 表
- [ ] 添加数据库索引
- [ ] 编写数据库迁移脚本

**验收标准:**
- 数据库可以正常创建
- 表结构符合 Schema 设计

---

### S0-T3: 项目架构设计

**负责人:** Tech Lead
**预估:** 2 人天
**依赖:** S0-T1, S0-T2

**任务内容:**
- [ ] 设计整体架构图
- [ ] 定义模块边界
- [ ] 设计数据流
- [ ] 定义接口规范
- [ ] 编写架构文档 ARCHITECTURE.md

**验收标准:**
- 架构文档完整
- 团队成员理解架构

---

### S0-T4: 数据库 Schema 评审

**负责人:** Tech Lead + 全员
**预估:** 1 人天
**依赖:** S0-T3

**任务内容:**
- [ ] 评审 accounts 表结构
- [ ] 评审 categories 表结构
- [ ] 评审 transactions 表结构
- [ ] 评审 products/warehouses/inventory 表结构
- [ ] 评审 budgets 表结构
- [ ] 确认最终 Schema

**验收标准:**
- 所有 Schema 已确认
- 无遗漏字段

---

### S0-T5: V0.1 核心页面设计稿

**负责人:** Designer
**预估:** 3 人天
**依赖:** 无

**任务内容:**
- [ ] 设计首页 (账户总览)
- [ ] 设计添加交易页
- [ ] 设计账户列表页
- [ ] 设计分类管理页
- [ ] 创建 UI 设计规范:
  - 颜色系统
  - 字体规范
  - 间距规范
  - 组件规范

**验收标准:**
- Figma 设计稿完成
- 设计稿已评审通过

---

### S0-T6: PRD 详细需求文档

**负责人:** PM
**预估:** 3 人天
**依赖:** 无

**任务内容:**
- [ ] 编写产品概述
- [ ] 编写用户画像
- [ ] 编写功能需求 (V0.1 - V1.0)
- [ ] 编写非功能需求
- [ ] 编写用户故事
- [ ] 编写原型图说明

**验收标准:**
- PRD 文档完整
- 团队成员理解需求

---

### S0-T7: 测试用例编写

**负责人:** QA
**预估:** 2 人天
**依赖:** S0-T5, S0-T6

**任务内容:**
- [ ] 编写账户管理测试用例
- [ ] 编写交易记录测试用例
- [ ] 编写分类管理测试用例
- [ ] 编写测试用例评审

**验收标准:**
- 测试用例覆盖 V0.1 功能
- 用例评审通过

---

### S0-T8: CI/CD 流水线配置

**负责人:** Dev1
**预估:** 2 人天
**依赖:** S0-T1

**任务内容:**
- [ ] 配置 GitHub Actions
- [ ] 配置 Flutter 分析工作流
- [ ] 配置测试工作流
- [ ] 配置构建工作流
- [ ] 配置部署工作流

**验收标准:**
- CI 流水线正常运行
- 代码推送后自动测试

---

## Sprint 0 交付物清单

| 交付物 | 负责人 | 截止日期 |
|--------|--------|----------|
| Flutter 项目骨架 | Dev1/2 | Day 3 |
| 数据库配置 | Dev1 | Day 4 |
| 架构设计文档 | Tech Lead | Day 5 |
| UI 设计稿 | Designer | Day 7 |
| PRD 文档 | PM | Day 7 |
| 测试用例 | QA | Day 8 |
| CI/CD 配置 | Dev1 | Day 8 |
| Sprint 0 评审 | 全体 | Day 10 |

---

## Sprint 0 评审检查点

- [ ] Flutter 环境可用
- [ ] 项目结构符合架构
- [ ] 数据库 Schema 确认
- [ ] 设计稿评审通过
- [ ] PRD 评审通过
- [ ] 测试用例评审通过
- [ ] CI/CD 正常运行
