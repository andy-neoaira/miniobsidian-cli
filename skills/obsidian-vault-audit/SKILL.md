---
name: obsidian-vault-audit
version: 1.0.0
description: "只读盘点 Obsidian Vault：列目录、搜索 TODO/主题/status，并生成整理建议。"
metadata:
  category: "obsidian"
  requires:
    bins: ["obs-cli"]
---

# obsidian-vault-audit

用于对 Obsidian Vault 做只读盘点、巡检和整理建议生成。

## 适用场景

当用户表达以下意图时使用：

- 检查我的 vault
- 找所有 TODO
- 看看哪些笔记提到了某个主题
- 帮我盘点项目笔记
- 查找 FIXME、draft、stale note 或未整理内容

## 需要收集的信息

- `directory`：要盘点的目录，可选；不提供则从 vault 根目录开始。
- `search_terms`：搜索词列表，可选。
- `audit_type`：`todo`、`topic`、`mentions`、`status` 或 `structure`。
- `max_results`：最大结果数，默认 25，最高建议 100。
- `vault`：目标 vault，可选。

## 操作流程

### 1. 查看目录结构

```bash
obs-cli list
obs-cli list "Projects"
```

指定 vault：

```bash
obs-cli list "Projects" --vault "<vault>"
```

### 2. 搜索关键词

```bash
obs-cli search-content "TODO" --format json --page 1 --page-size 100
```

常见 audit 类型默认搜索词：

- `todo`：`TODO`、`FIXME`、`待办`、`未完成`
- `status`：`status:`、`draft`、`active`、`done`
- `topic`：使用用户提供的主题词
- `structure`：优先 `list` 目录，再按需要搜索

### 3. 查看特定笔记及反向链接

```bash
obs-cli print "<note>" --mentions
```

## 规则

- 本 skill 默认只读，不修改 vault。
- 不调用 `create`、`move`、`delete`、`frontmatter --edit`。
- 如果用户要求“顺便整理 / 移动 / 标记”，应建议切换到 `obsidian-inbox-triage` 或 `obsidian-project-note`，并先确认。
- 不要一次打印大量完整笔记；优先输出统计、列表和摘要。
- `max_results` 不应超过 `search-content --page-size` 上限 100。

## 输出格式

输出应包含：

- 盘点范围
- 使用的搜索词或 audit 类型
- 命中列表和分类统计
- 高优先级整理建议
- 可执行的后续 obs-cli 命令，例如移动、追加 daily、更新 frontmatter 等，但不要直接执行这些写入命令
