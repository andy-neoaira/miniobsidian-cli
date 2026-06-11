---
name: obsidian-knowledge-search
version: 1.0.0
description: "从 Obsidian Vault 检索相关笔记，打印内容并按需附带 linked mentions 作为上下文。"
metadata:
  category: "obsidian"
  requires:
    bins: ["obs-cli"]
---

# obsidian-knowledge-search

用于从本地 Obsidian Vault 中搜索知识，并把相关笔记内容取出来作为回答、研究或分析的上下文。

## 适用场景

当用户表达以下意图时使用：

- 在我的 Obsidian 里找
- 搜索 vault
- 找相关笔记
- 把相关笔记内容拿出来给我
- 从我的知识库里查一下
- 给当前问题找本地知识库上下文

## 需要收集的信息

- `query`：搜索关键词，必需。
- `top_n`：返回数量，默认 5，通常不超过 10。
- `vault`：目标 vault，可选。
- `include_mentions`：是否用 `print --mentions` 附带反向链接。
- `open_selected`：是否打开选中的笔记。

## 操作流程

### 1. 先用 JSON 搜索

```bash
obs-cli search-content "<query>" --format json --page 1 --page-size 10
```

指定 vault：

```bash
obs-cli search-content "<query>" --format json --page 1 --page-size 10 --vault "<vault>"
```

### 2. 根据结果选择相关笔记

优先选择：

- 文件名匹配 query 的笔记
- 内容片段明显相关的笔记
- 用户指定目录或项目下的笔记
- 重复出现频率较高的笔记

### 3. 打印笔记内容

```bash
obs-cli print "<note>"
```

附带 linked mentions：

```bash
obs-cli print "<note>" --mentions
```

### 4. 可选打开笔记

```bash
obs-cli open "<note>"
obs-cli open "<note>" --editor
```

## 规则

- 默认只读，不修改 vault。
- 搜索时优先使用 `--format json`，便于结构化处理。
- 不要一次打印大量完整笔记；先汇总命中列表，再按 `top_n` 读取。
- `--mentions` 可能输出较多内容，只在用户要求或上下文需要时使用。
- 只有用户明确要求打开时才调用 `open`。

## 输出格式

回答中包含：

- 命中的 note 列表
- 每个 note 的相关片段摘要
- 已打印 / 建议读取的 note
- 如果使用了 `--mentions`，说明反向链接摘要
- 如果没有命中，给出可尝试的其他关键词
