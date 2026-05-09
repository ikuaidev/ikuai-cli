---
name: ikuai-version
description: Show ikuai-cli build metadata such as version, commit, and build date.
---

# Version

人工检查时默认先显示版本号：

```bash
ikuai-cli version
```

发布构建的版本号来自 git tag，例如 `v1.0.8`。

Agent 或脚本需要解析当前 CLI 构建信息时，优先使用 JSON：

```bash
ikuai-cli version --format json
ikuai-cli version --format yaml
ikuai-cli version --columns name,version,commit,date
```
