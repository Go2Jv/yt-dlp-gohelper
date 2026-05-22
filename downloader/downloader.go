package downloader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Go2Jv/yt-dlp-gohelper/execx"
	"github.com/Go2Jv/yt-dlp-gohelper/i18n"
	"github.com/Go2Jv/yt-dlp-gohelper/ui"
)

type Request struct {
	URL          string
	Args         []string
	OutputDir    string
	OutputTmpl   string
}

func PromptRequest(msg *i18n.Messages, goos string) Request {
	lang := msg.Lang()

	switch lang {
	case i18n.Zh:
		fmt.Println("直接回车 = 使用默认值")
	case i18n.Ja:
		fmt.Println("Enter = デフォルトを使用")
	default:
		fmt.Println("Press Enter to use defaults")
	}

	var url string
	switch lang {
	case i18n.Zh:
		url = ui.PromptRequired("请输入 视频链接 (必填): ")
	case i18n.Ja:
		url = ui.PromptRequired("動画URL (必須): ")
	default:
		url = ui.PromptRequired("Video URL (required): ")
	}

	cookieBrowser := promptBrowser(lang, goos)
	format := promptQuality(lang)

	wd, _ := os.Getwd()
	var savePath string
	switch lang {
	case i18n.Zh:
		savePath = ui.Prompt("请输入保存位置 (默认当前文件夹): ", wd)
	case i18n.Ja:
		savePath = ui.Prompt("保存先 (デフォルト: 現在のフォルダ): ", wd)
	default:
		savePath = ui.Prompt("Save directory (default current): ", wd)
	}
	savePath = strings.TrimSpace(savePath)
	if savePath == "" {
		savePath = wd
	}

	var fileName string
	switch lang {
	case i18n.Zh:
		fileName = ui.Prompt("请输入保存文件名 (默认原视频名): ", "%(title)s")
	case i18n.Ja:
		fileName = ui.Prompt("ファイル名 (デフォルト: 元のタイトル): ", "%(title)s")
	default:
		fileName = ui.Prompt("File name (default original title): ", "%(title)s")
	}
	fileName = strings.TrimSpace(fileName)
	if fileName == "" {
		fileName = "%(title)s"
	}

	outTemplate := filepath.Join(savePath, fileName+".%(ext)s")

	args := make([]string, 0, 16)
	if cookieBrowser != "" {
		args = append(args, "--cookies-from-browser", cookieBrowser)
	}
	if goos == "windows" {
		args = append(args, "--windows-filenames")
	}
	args = append(args, "-f", format)
	args = append(args, "-o", outTemplate)
	args = append(args, url)

	return Request{
		URL:        url,
		Args:       args,
		OutputDir:  savePath,
		OutputTmpl: outTemplate,
	}
}

func RunRequest(msg *i18n.Messages, req Request) (string, error) {
	lang := msg.Lang()
	switch lang {
	case i18n.Zh:
		fmt.Println("开始下载...")
	case i18n.Ja:
		fmt.Println("ダウンロード開始...")
	default:
		fmt.Println("Downloading...")
	}

	return execx.RunInteractiveCapture(4*1024*1024, "yt-dlp", req.Args...)
}

func RunInteractive(msg *i18n.Messages, goos string) error {
	req := PromptRequest(msg, goos)
	_, err := RunRequest(msg, req)
	return err
}

func promptBrowser(lang i18n.Lang, goos string) string {
	type opt struct {
		label string
		val   string
	}

	opts := []opt{
		{label: "Chrome", val: "chrome"},
		{label: "Firefox", val: "firefox"},
		{label: "Edge", val: "edge"},
	}
	if goos == "windows" {
		opts = append(opts, opt{label: "360 (Chromium)", val: "chrome"})
	}
	if goos == "darwin" {
		opts = append(opts, opt{label: "Safari", val: "safari"})
	}
	opts = append(opts, opt{label: "None", val: ""})

	switch lang {
	case i18n.Zh:
		fmt.Println("\n--- 浏览器 Cookie 选择 ---")
	case i18n.Ja:
		fmt.Println("\n--- ブラウザ Cookie ---")
	default:
		fmt.Println("\n--- Browser Cookies ---")
	}

	for i, o := range opts {
		fmt.Printf("%d. %s\n", i+1, o.label)
	}
	def := fmt.Sprintf("%d", len(opts))
	choice := ui.Prompt("Select (default "+def+"): ", def)
	idx := parseChoice(choice, len(opts), len(opts)-1)
	return opts[idx].val
}

func promptQuality(lang i18n.Lang) string {
	type opt struct {
		label  string
		format string
	}
	opts := []opt{
		{label: "360P", format: "bestvideo[height<=360]+bestaudio/best"},
		{label: "720P", format: "bestvideo[height<=720]+bestaudio/best"},
		{label: "1080P", format: "bestvideo[height<=1080]+bestaudio/best"},
		{label: "1440P", format: "bestvideo[height<=1440]+bestaudio/best"},
		{label: "2160P", format: "bestvideo[height<=2160]+bestaudio/best"},
		{label: "Best", format: "best"},
	}

	switch lang {
	case i18n.Zh:
		fmt.Println("\n--- 分辨率选择 ---")
	case i18n.Ja:
		fmt.Println("\n--- 画質 ---")
	default:
		fmt.Println("\n--- Quality ---")
	}

	for i, o := range opts {
		fmt.Printf("%d. %s\n", i+1, o.label)
	}
	def := fmt.Sprintf("%d", len(opts))
	choice := ui.Prompt("Select (default "+def+"): ", def)
	idx := parseChoice(choice, len(opts), len(opts)-1)
	return opts[idx].format
}

func parseChoice(input string, max int, fallback int) int {
	input = strings.TrimSpace(input)
	if input == "" {
		return fallback
	}
	n := 0
	for _, r := range input {
		if r < '0' || r > '9' {
			return fallback
		}
		n = n*10 + int(r-'0')
	}
	if n <= 0 || n > max {
		return fallback
	}
	return n - 1
}
