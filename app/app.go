package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Go2Jv/yt-dlp-gohelper/deps"
	"github.com/Go2Jv/yt-dlp-gohelper/downloader"
	"github.com/Go2Jv/yt-dlp-gohelper/i18n"
	"github.com/Go2Jv/yt-dlp-gohelper/install"
	"github.com/Go2Jv/yt-dlp-gohelper/ui"
)

func Run() int {
	lang := i18n.PromptLanguage(i18n.Zh)
	msg := i18n.NewMessages(lang)

	ui.PrintLogo()
	fmt.Println(msg.Welcome(runtime.GOOS))

	for {
		switch promptMainMenu(msg) {
		case 1:
			if !ensureDeps(msg) {
				continue
			}
			if err := runDownloadFlow(msg, runtime.GOOS); err != nil {
				fmt.Println(msg.Error(err))
			}
		default:
			return 0
		}
	}
}

func ensureDeps(msg *i18n.Messages) bool {
	state := deps.Check()
	if !state.MissingAny() {
		return true
	}

	fmt.Println(msg.DepsMissing(state))
	choice := msg.PromptInstallChoice()
	switch choice {
	case i18n.InstallChoiceAuto:
		if err := install.Ensure(msg, state); err != nil {
			fmt.Println(msg.Error(err))
			fmt.Println(msg.ManualInstallHint(runtime.GOOS))
			return false
		}
	case i18n.InstallChoiceManual:
		fmt.Println(msg.ManualInstallHint(runtime.GOOS))
		return false
	default:
		return false
	}

	if err := deps.VerifyAll(); err != nil {
		fmt.Println(msg.Error(err))
		fmt.Println(msg.ManualInstallHint(runtime.GOOS))
		return false
	}

	return true
}

func Exit(code int) {
	os.Exit(code)
}

func promptMainMenu(msg *i18n.Messages) int {
	switch msg.Lang() {
	case i18n.Zh:
		fmt.Println("1. 下载视频")
		fmt.Println("0. 退出")
		choice := ui.Prompt("请选择 (默认1): ", "1")
		return parseChoice(choice, 1, 1)
	case i18n.Ja:
		fmt.Println("1. 動画をダウンロード")
		fmt.Println("0. 終了")
		choice := ui.Prompt("選択 (デフォルト1): ", "1")
		return parseChoice(choice, 1, 1)
	default:
		fmt.Println("1. Download video")
		fmt.Println("0. Exit")
		choice := ui.Prompt("Select (default 1): ", "1")
		return parseChoice(choice, 1, 1)
	}
}

func runDownloadFlow(msg *i18n.Messages, goos string) error {
	req := downloader.PromptRequest(msg, goos)

	for {
		_, err := downloader.RunRequest(msg, req)
		if err == nil {
			fmt.Println(msg.Done())
			return nil
		}

		if !promptRetryAfterFailure(msg) {
			temps, _ := findTempFiles(req.OutputDir)
			deleteFiles(temps)
			return err
		}
	}
}

func promptRetryAfterFailure(msg *i18n.Messages) bool {
	switch msg.Lang() {
	case i18n.Zh:
		fmt.Print("下载失败，是否重新下载？(y/n, 默认y): ")
		return ui.PromptYesNo("", true)
	case i18n.Ja:
		fmt.Print("失敗しました。再試行しますか？(y/n, デフォルトy): ")
		return ui.PromptYesNo("", true)
	default:
		fmt.Print("Download failed. Retry? (y/n, default y): ")
		return ui.PromptYesNo("", true)
	}
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
	if n < 0 || n > max {
		return fallback
	}
	return n
}

func findTempFiles(outputDir string) ([]string, error) {
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, 16)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !looksLikeTemp(name) {
			continue
		}
		out = append(out, filepath.Join(outputDir, name))
		if len(out) >= 200 {
			break
		}
	}
	return out, nil
}

func deleteFiles(files []string) {
	for _, f := range files {
		_ = os.Remove(f)
	}
}

func looksLikeTemp(name string) bool {
	if strings.Contains(name, ".part") {
		return true
	}
	switch {
	case strings.HasSuffix(name, ".ytdl"):
		return true
	case strings.HasSuffix(name, ".aria2"):
		return true
	case strings.HasSuffix(name, ".temp"):
		return true
	default:
		return false
	}
}
