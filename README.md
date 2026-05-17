# yt-dlp-gohelper

![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go&logoColor=white)
![Release](https://img.shields.io/github/v/release/Go2Jv/yt-dlp-gohelper?display_name=tag)

A cross-platform interactive helper for `yt-dlp`.

## 下载 / Downloads

Latest Release (Releases page): https://github.com/Go2Jv/yt-dlp-gohelper/releases/latest

| OS | Arch | Download |
|---|---|---|
| Windows | x64 | https://github.com/Go2Jv/yt-dlp-gohelper/releases/latest |
| Linux | x64 | https://github.com/Go2Jv/yt-dlp-gohelper/releases/latest |
| macOS | Intel (x64) | https://github.com/Go2Jv/yt-dlp-gohelper/releases/latest |
| macOS | Apple Silicon (arm64) | https://github.com/Go2Jv/yt-dlp-gohelper/releases/latest |

## 支持 / Support

如果这个项目对你有帮助，可以请我喝杯咖啡。谢谢 :)

WeChat Pay：请在仓库中添加二维码图片 `donate/wechat.jpg` 后即可在 README 中展示。

## 中文

一个跨平台的 `yt-dlp` 交互式下载小工具：启动后会自动识别系统（Windows / Linux / macOS / 苹果系统）与语言（中文 / English / 日本語），检查 `yt-dlp` 与 `ffmpeg` 是否可用；依赖缺失时提供自动安装或退出自行安装两种模式。

### 功能

- 自动识别系统与语言（也会提示你选择语言）
- 启动 Logo + 主菜单（下载视频 / 下载失败处理）
- 检测 `yt-dlp` / `ffmpeg` 是否已安装且可运行
- 依赖缺失时可选择自动安装
  - Windows：优先使用 `winget`
  - Linux：优先使用 `apt-get`，否则仅支持用 `curl` 安装 `yt-dlp`（`ffmpeg` 需手动）
  - macOS / 苹果系统：使用 `brew`；如未安装 `brew` 会询问是否安装
- 下载流程（交互式）
  - 输入视频链接（必填）
  - 选择浏览器 Cookie（macOS 才会显示 Safari）
  - 选择分辨率预设
  - 选择保存目录与文件名模板
- 下载失败处理（v2 稳定版）
  - 下载失败后可选择重新下载
  - 不重新下载时会自动删除 temp 并清理历史记录
  - 不会写入历史记录文件（避免隐私与磁盘占用）

### 重要提示（中国大陆）

- 中国大陆用户下载 YouTube 等外网视频，请使用全局 VPN（全局模式）。仅“代理模式/分应用代理”经常会导致 `yt-dlp` / `ffmpeg` 访问失败或下载速度异常。

### Cookie 提示（Windows / macOS）

`--cookies-from-browser` 在部分系统/浏览器版本上可能不稳定，常见原因是浏览器锁定 cookie 数据库或权限不足。遇到失败时建议：

- 先完全关闭浏览器（确保后台进程也退出）后重试
- 或者改为手动导出 cookies.txt，然后使用 `yt-dlp --cookies cookies.txt ...`
- 相关参考（yt-dlp Issue）：Windows Chromium/Edge 读取 cookie 的权限问题（#7271）与 Firefox 的 Permission denied（#15760）
  - https://github.com/yt-dlp/yt-dlp/issues/7271
  - https://github.com/yt-dlp/yt-dlp/issues/15760

### 运行（源码）

```bash
go run .
```

或编译后运行：

```bash
go build -o yt-dlp-gohelper .
./yt-dlp-gohelper
```

---

## English

On startup it detects your OS (Windows / Linux / macOS) and language (中文 / English / 日本語), checks whether `yt-dlp` and `ffmpeg` are available, and lets you either auto-install missing dependencies or exit and install them manually.

### Features (v2 stable)

- Startup logo + main menu (download / failed tasks)
- Failed download handling: retry or delete temp automatically
- No history file is written (to avoid privacy and disk usage)

### Notes (Mainland China)

- If you are in Mainland China and you download YouTube/overseas sites, use a system-wide VPN (global mode). Proxy-only modes may cause `yt-dlp` / `ffmpeg` to fail or become unstable.

### Cookies From Browser (Windows / macOS)

`--cookies-from-browser` may be unstable on some OS/browser versions (database locked / permission issues). If it fails:

- Fully close the browser (including background processes) and retry
- Or export cookies to `cookies.txt` and use `yt-dlp --cookies cookies.txt ...`
- References (yt-dlp issues): #7271, #15760
  - https://github.com/yt-dlp/yt-dlp/issues/7271
  - https://github.com/yt-dlp/yt-dlp/issues/15760

### Run from source

```bash
go run .
```

---

## 日本語

起動時に OS（Windows / Linux / macOS）と言語（中文 / English / 日本語）を判定し、`yt-dlp` と `ffmpeg` の有無を確認します。不足している場合は自動インストールするか、終了して手動インストールするかを選べます。

### 機能（v2 安定版）

- 起動ロゴ + メインメニュー（ダウンロード / 失敗タスク処理）
- 失敗時に再試行、または temp を自動削除
- 履歴ファイルは保存しません（プライバシー/容量のため）

### 注意（中国本土）

- 中国本土から YouTube など海外サイトをダウンロードする場合、VPN を全体モード（グローバル）にしてください。プロキシのみのモードだと `yt-dlp` / `ffmpeg` が失敗することがあります。

### ブラウザ Cookie について（Windows / macOS）

`--cookies-from-browser` は一部の環境で不安定です（DB ロック/権限など）。失敗する場合：

- ブラウザを完全に終了して再試行
- もしくは cookies.txt を手動でエクスポートして `yt-dlp --cookies cookies.txt ...` を使用
- 参考（yt-dlp issue）：#7271, #15760
  - https://github.com/yt-dlp/yt-dlp/issues/7271
  - https://github.com/yt-dlp/yt-dlp/issues/15760
