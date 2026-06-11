---
name: obsidian-daily-log
version: 1.0.0
description: "把工作日志、会议记录、session 总结或临时流水账追加到 Obsidian Daily Note。"
metadata:
  category: "obsidian"
  requires:
    bins: ["obs-cli"]
---

# obsidian-daily-log

用于通过 `obs-cli daily` 创建 / 打开今天的 Daily Note，并向其中追加内容。

## 适用场景

当用户表达以下意图时使用：

- 写到今天的 daily note
- 记录今天完成了什么
- 把这次 session 总结追加到日记
- 写一条工作日志
- 追加站会记录、会议摘要、TODO 或临时流水账

## 需要收集的信息

- `content`：要追加的内容。
- `section`：可选，例如 `Work Log`、`Meeting`、`Summary`、`TODO`。
- `format`：`bullet`、`section` 或 `raw`。
- `open_after`：追加后是否打开编辑器。
- `vault`：目标 vault，可选。

## 内容格式

### bullet

适合一条或多条简短记录：

```markdown
- 09:30 standup
- 完成 obs-cli skill 规划
```

### section

适合 session 总结、会议记录等结构化内容：

```markdown
## Session Summary

...
```

### raw

用户给什么就写什么，不额外改写。

## 操作流程

### 追加内容到 Daily Note

```bash
printf '%s' "<content>" | obs-cli daily --content-file -
```

指定 vault：

```bash
printf '%s' "<content>" | obs-cli daily --content-file - --vault "<vault>"
```

### 追加后打开编辑器

```bash
printf '%s' "<content>" | obs-cli daily --content-file - --editor
```

### 只打开今天的 Daily Note

```bash
obs-cli daily --editor
```

## 规则

- Daily Note 的目录、日期格式、模板由 obs-cli / Obsidian 配置决定；不要硬编码 `Dailies/`。
- 用户明确要求“记录 / 写到 daily / 追加到日记”时，可以执行写入。
- 只有用户明确要求打开编辑器时才加 `--editor`。
- 多行内容或用户原文必须优先用 `--content-file -` 从 stdin 传入，不要直接拼接进 shell 字符串。
- 不使用 `create` 直接写 daily 文件；优先使用 `obs-cli daily`。

## 输出格式

完成后汇报：

- 已追加的内容摘要
- 使用的格式：bullet / section / raw
- 是否已打开编辑器
- 如能从命令输出或配置推断，则说明 Daily Note 位置
