# figma-cli

[English](README.md)

`figma-cli` 是一个面向 AI agent 和人类用户的 Figma REST API 命令行工具。

它的定位是稳定、可复现地读取 Figma 信息和导出素材：

- 直接粘贴 Figma file / design / proto URL。
- 自动解析 file key 和 node ID。
- 不打开 Figma app 也能读取文件和节点摘要。
- 导出节点图片 URL，或下载素材到本地。
- 列出 comments，给 review / 实现提供上下文。
- 列出本地 components、styles、variables。
- 对重复读取做有 TTL 的本地缓存。
- token 存入 OS keyring；CI 或一次性调用可使用 `FIGMA_TOKEN`。

这个 CLI 是刻意只读的：只做检视与导出，永不写回 Figma。读取基于 token、结果可复现——脚本、CI、无头 agent 环境里行为一致。

官方 API 边界、已有工具经验和 v1 非目标见 [docs/CAPABILITY_REVIEW.md](docs/CAPABILITY_REVIEW.md)。

## 安装

Homebrew（macOS）：

```bash
brew install --cask wohsj110/tap/figma-cli
```

Go 安装：

```bash
go install github.com/wohsj110/figma_cli/cmd/figma-cli@latest
```

或本地编译：

```bash
make build
./bin/figma-cli help
```

## 配置

先创建 Figma personal access token，然后写入 OS keyring：

```bash
printf %s "$FIGMA_TOKEN" | figma-cli init --token-stdin
```

CI 或一次性调用可以直接使用环境变量：

```bash
export FIGMA_TOKEN="figd_..."
```

`FIGMA_TOKEN` 优先级高于 keyring。CLI 不会打印 token。

本地 mock server 或代理集成测试可以覆盖 API 地址：

```bash
export FIGMA_API_BASE_URL="http://127.0.0.1:8080/v1"
```

## 命令

查看当前 Figma 用户：

```bash
figma-cli me
```

查看文件摘要：

```bash
figma-cli file get "https://www.figma.com/design/FILE_KEY/name?node-id=1-2"
figma-cli file get FILE_KEY --depth 2
```

查看节点摘要：

```bash
figma-cli node inspect "https://www.figma.com/design/FILE_KEY/name?node-id=1-2"
figma-cli node inspect FILE_KEY --node 1:2
```

导出节点图片 URL：

```bash
figma-cli image export FILE_KEY --node 1:2 --format png --scale 2
```

下载导出的图片：

```bash
figma-cli image export FILE_KEY --node 1:2 --format png --scale 2 --out ./figma-assets
```

列出 comments：

```bash
figma-cli comments list FILE_KEY
```

列出可复用设计元数据：

```bash
figma-cli components list FILE_KEY
figma-cli styles list FILE_KEY
figma-cli variables list FILE_KEY
```

agent 需要原始 API 响应时使用 `--raw`：

```bash
figma-cli node inspect FILE_KEY --node 1:2 --raw
```

本地缓存控制：

```bash
figma-cli file get FILE_KEY --cache-ttl 1h
figma-cli file get FILE_KEY --no-cache
```

## Agent Skills

用开放的 `skills` CLI 安装内置的 `figma-cli` skill（自动识别 Claude Code / Codex / Cursor）：

```bash
npx skills add wohsj110/figma_cli
```

### For LLM——整段复制给任意 agent 即可安装

```text
Install and verify the figma-cli skill:

1. Run: npx skills add wohsj110/figma_cli --yes
   (target agents with --agent codex|claude-code; add --global for a user-wide install)
2. Ensure the CLI binary exists: command -v figma-cli
   If missing: brew install --cask wohsj110/tap/figma-cli
   (no Homebrew: go install github.com/wohsj110/figma_cli/cmd/figma-cli@latest)
3. Verify credentials: figma-cli me
   If not configured, follow "Credential Setup" in the installed SKILL.md and ask the
   user for a Figma personal access token — never guess node ids or print token values.
```

skills.sh 上的 skill 页面（进入搜索/排行榜靠安装量累积）：https://www.skills.sh/wohsj110/figma_cli/figma-cli

## 输出契约

- 默认输出是简洁文本，适合人类阅读，也适合 agent 解析。
- `--raw` 输出格式化后的 Figma REST API JSON。
- `--verbose` 输出请求诊断，但不会输出 token。
- GET JSON 响应默认缓存 15 分钟。
- `--no-cache` 可关闭单次调用缓存。
- `--cache-ttl` 支持 Go duration 字符串，比如 `30s`、`15m` 或 `1h`。
- `FIGMA_API_BASE_URL` 可用于本地测试时覆盖 API host。
- 大文件内容需要通过 `--depth` 或 `--raw` 显式请求。

## 发布准备

GoReleaser 构建 GitHub release archives，并把 `figma-cli` cask 发布到 `wohsj110/homebrew-tap`。

```bash
go test ./...
make build
make mock-verify
make release-check
```

见 [docs/RELEASING.md](docs/RELEASING.md)。

## Roadmap

- design token 导出为 CSS / JSON。
- ~~仓库名和 tap 名确认后启用 Homebrew cask 发布。~~ 已交付：`brew install --cask wohsj110/tap/figma-cli`。
- Codex / Claude Code skills。

## 开发

```bash
make test
make build
make mock-verify
```
