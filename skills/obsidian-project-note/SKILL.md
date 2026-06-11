---
name: obsidian-project-note
version: 1.0.0
description: "创建或更新项目级 Obsidian 笔记，并维护 project/status/type 等 frontmatter。"
metadata:
  category: "obsidian"
  requires:
    bins: ["obs-cli"]
---

# obsidian-project-note

用于把项目分析、架构说明、session log、TODO 或决策记录写入 Obsidian 项目笔记。

## 适用场景

当用户表达以下意图时使用：

- 给这个项目建一份 Obsidian 笔记
- 把代码分析结果写入 Vault
- 更新项目知识库
- 把这次项目 review 存档
- 为某个项目追加 session log / decision / TODO

## 需要收集的信息

- `project_name`：项目名称，必需。
- `note_type`：`overview`、`architecture`、`session-log`、`todo` 或 `decision`。
- `content`：要写入的 Markdown 内容。
- `mode`：`create-only`、`append` 或 `overwrite`。
- `target_folder`：默认 `Projects/<project_name>`。
- `open_after`：是否写入后打开编辑器。
- `vault`：目标 vault，可选。

## 默认路径约定

- overview：`Projects/<project_name>/Overview.md`
- architecture：`Projects/<project_name>/Architecture.md`
- session-log：`Projects/<project_name>/Session Log.md`
- todo：`Projects/<project_name>/TODO.md`
- decision：`Projects/<project_name>/Decisions.md`

这些目录是 skill 层约定，不是 obs-cli 内置特殊目录。

## 操作流程

### 创建或更新项目笔记

```bash
printf '%s' "<markdown>" | obs-cli create "Projects/<project_name>/Overview.md" --content-file -
```

追加：

```bash
printf '%s' "<markdown>" | obs-cli create "Projects/<project_name>/Session Log.md" --content-file - --append
```

覆盖：

```bash
printf '%s' "<markdown>" | obs-cli create "Projects/<project_name>/Overview.md" --content-file - --overwrite
```

### 设置 frontmatter

```bash
obs-cli frontmatter "<project-note>" --edit --key type --value project
obs-cli frontmatter "<project-note>" --edit --key project --value "<project_name>"
obs-cli frontmatter "<project-note>" --edit --key status --value active
```

也可以根据 `note_type` 设置：

```bash
obs-cli frontmatter "<project-note>" --edit --key note_type --value "<note_type>"
```

### 可选打开编辑器

```bash
obs-cli open "<project-note>" --editor
```

## 规则

- 默认不覆盖已有笔记；只有用户明确要求替换、重写、覆盖时才使用 `--overwrite`。
- session log 默认使用 `--append`。
- 写入成功后再更新 frontmatter。
- 如果用户没有给内容，应先要求用户提供内容或说明要生成哪类项目笔记。
- 多行 Markdown 必须优先用 `--content-file -` 从 stdin 传入，避免 shell 引号破坏内容。
- 不调用 `delete`。

## 输出格式

完成后汇报：

- project note path
- 使用的更新模式：create-only / append / overwrite
- 写入或更新的 frontmatter 字段
- 后续建议维护的关联笔记，例如 Architecture、Decisions、TODO
