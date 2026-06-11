---
name: obsidian-capture
version: 1.0.0
description: "快速把想法、摘要、任务或 AI 输出捕获到 Obsidian Inbox 或指定笔记。"
metadata:
  category: "obsidian"
  requires:
    bins: ["obs-cli"]
---

# obsidian-capture

用于把内容快速写入 Obsidian。默认面向 Inbox 工作流，但也允许用户指定任意目标路径。

## 适用场景

当用户表达以下意图时使用：

- 记到 Obsidian
- 保存到我的笔记
- 放到 Inbox
- 把这段内容写成一条笔记
- 把 AI 输出、网页摘要、review 结论或临时任务保存起来

## 需要收集的信息

- `content`：要保存的内容，必需。
- `title`：笔记标题，可选。
- `target_path`：目标笔记路径，可选。
- `mode`：`create` 或 `append`。
- `open_after`：保存后是否打开编辑器。
- `vault`：目标 vault，可选；如提供则使用 `--vault`。

## 默认路径约定

如果用户没有指定路径：

- 单独成文：使用 `Inbox/<title>.md`。
- 追加流水账：使用 `Inbox.md`。

`Inbox` 是 skill 层约定，不是 obs-cli 内置特殊目录。

## 操作流程

### 创建单独捕获笔记

```bash
printf '%s' "<content>" | obs-cli create "Inbox/<title>.md" --content-file -
```

如果用户指定 vault：

```bash
printf '%s' "<content>" | obs-cli create "Inbox/<title>.md" --content-file - --vault "<vault>"
```

### 追加到统一 Inbox 文件

```bash
printf '%s' "<content>" | obs-cli create "Inbox.md" --content-file - --append
```

### 保存后打开编辑器

```bash
obs-cli open "<note-path>" --editor
```

或者创建时直接打开：

```bash
printf '%s' "<content>" | obs-cli create "<note-path>" --content-file - --open --editor
```

## 内容格式建议

- 临时想法：保持原文即可。
- 任务：优先转成 Markdown bullet。
- 摘要：建议包含来源、时间、摘要、后续动作。
- 多行内容或用户原文：优先使用 `--content-file -` 从 stdin 传入，避免 shell 引号破坏内容。

## 安全规则

- 不要默认使用 `--overwrite`。如果目标已存在，优先使用 `--append` 或询问用户。
- 只有用户明确要求打开 / 编辑 / review 时才使用 `--editor`。
- 如果目标路径是根据标题推断的，汇报最终路径。
- 不要把未转义的用户长文本直接拼接到 `--content "<content>"` 中。
- 不调用 `delete`。

## 输出格式

完成后汇报：

- 创建或追加的 note path
- 使用的是 create 还是 append
- 是否已打开编辑器
- 后续整理建议，例如可用 `obsidian-inbox-triage` 归档
