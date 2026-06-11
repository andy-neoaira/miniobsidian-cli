---
name: obsidian-inbox-triage
version: 1.0.0
description: "整理 Obsidian Inbox：列出临时笔记、移动到正式目录，并更新状态 frontmatter。"
metadata:
  category: "obsidian"
  requires:
    bins: ["obs-cli"]
---

# obsidian-inbox-triage

用于整理 Inbox 中的临时笔记，将其移动到 Projects、Areas、Notes 或其他正式目录，并按需更新 frontmatter。

## 适用场景

当用户表达以下意图时使用：

- 帮我整理 Inbox
- 把这个临时笔记移动到项目目录
- 归档这条笔记
- 把 Inbox 里的内容分类
- 将捕获内容转为正式知识卡片

## 需要收集的信息

- `source_note`：源笔记路径，例如 `Inbox/tmp.md`。
- `target_path`：目标路径，例如 `Projects/MyProject/final.md`。
- `classification`：分类，可选，例如 `project`、`area`、`note`、`archive`。
- `status`：默认 `organized`。
- `open_after`：移动后是否打开编辑器。
- `vault`：目标 vault，可选。

## 操作流程

### 1. 如果源笔记不明确，先列出 Inbox

```bash
obs-cli list "Inbox"
```

### 2. 移动 / 重命名笔记

```bash
obs-cli move "<source-note>" "<target-note>"
```

`obs-cli move` 会移动文件并更新 vault 中指向该笔记的链接。

### 3. 更新 frontmatter 状态

```bash
obs-cli frontmatter "<target-note>" --edit --key status --value organized
```

如果有分类：

```bash
obs-cli frontmatter "<target-note>" --edit --key classification --value "<classification>"
```

### 4. 可选打开编辑器

```bash
obs-cli open "<target-note>" --editor
```

## 规则

- `Inbox` 是 skill 层约定，不是 obs-cli 内置特殊目录。
- 如果 source 或 target 是推断出来的，移动前需要向用户确认。
- 不要静默发明目标路径；除非用户明确给出或清楚表达分类规则。
- `move` 是可逆成本较高的操作，执行后要清楚汇报源路径和目标路径。
- 不调用 `delete`。

## 输出格式

完成后汇报：

- 移动前路径
- 移动后路径
- 更新的 frontmatter 字段
- 是否已打开编辑器
- 是否建议补充 tags、project、status 等 metadata
