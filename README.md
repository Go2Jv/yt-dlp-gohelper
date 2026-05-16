# yt-dlp-gohelper

## 中文

一个跨平台的 `yt-dlp` 交互式下载小工具：启动后会自动识别系统（Windows / Linux / macOS）与语言（中文 / English / 日本語），检查 `yt-dlp` 与 `ffmpeg` 是否可用；依赖缺失时提供自动安装或退出自行安装两种模式。

### 功能

- 自动识别系统与语言（也会提示你选择语言）
- 检测 `yt-dlp` / `ffmpeg` 是否已安装且可运行
- 依赖缺失时可选择自动安装
  - Windows：优先使用 `winget`
  - Linux：优先使用 `apt-get`，否则仅支持用 `curl` 安装 `yt-dlp`（`ffmpeg` 需手动）
  - macOS：使用 `brew`；如未安装 `brew` 会询问是否安装
- 下载流程（交互式）
  - 输入视频链接（必填）
  - 选择浏览器 Cookie（macOS 才会显示 Safari）
  - 选择分辨率预设
  - 选择保存目录与文件名模板

### 依赖说明

本程序只是帮你调用 `yt-dlp`，实际下载能力取决于你本机的 `yt-dlp` 与 `ffmpeg`。

### 运行

```bash
go run .
```

或编译后运行：

```bash
go build -o yt-dlp-gohelper .
./yt-dlp-gohelper
```

### 常见问题

- Windows 下安装完仍提示找不到 `yt-dlp/ffmpeg`
  - 请新开一个终端再试
  - 确认 `yt-dlp` 与 `ffmpeg` 已加入 PATH（能在 CMD/PowerShell 直接运行）

---

## English

A cross-platform interactive helper for `yt-dlp`. On startup it detects your OS (Windows / Linux / macOS) and language (中文 / English / 日本語), checks whether `yt-dlp` and `ffmpeg` are available, and lets you either auto-install missing dependencies or exit and install them manually.

### Features

- Detect OS and language (and asks you to confirm the language)
- Check `yt-dlp` / `ffmpeg` availability (must be runnable)
- Auto-install missing dependencies
  - Windows: prefers `winget`
  - Linux: prefers `apt-get`, otherwise supports installing `yt-dlp` via `curl` only (`ffmpeg` must be installed manually)
  - macOS: uses `brew`; if `brew` is missing it will ask whether to install it
- Interactive download flow
  - Video URL (required)
  - Browser cookies (Safari is shown on macOS only)
  - Quality preset
  - Output directory and file name template

### Run

```bash
go run .
```

Or build and run:

```bash
go build -o yt-dlp-gohelper .
./yt-dlp-gohelper
```

### Troubleshooting

- On Windows, if it still says `yt-dlp/ffmpeg` not found after installation:
  - Open a new terminal and retry
  - Ensure `yt-dlp` and `ffmpeg` are in PATH (runnable in CMD/PowerShell)
