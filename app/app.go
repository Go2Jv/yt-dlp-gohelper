package app

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/Go2Jv/yt-dlp-gohelper/deps"
	"github.com/Go2Jv/yt-dlp-gohelper/downloader"
	"github.com/Go2Jv/yt-dlp-gohelper/history"
	"github.com/Go2Jv/yt-dlp-gohelper/i18n"
	"github.com/Go2Jv/yt-dlp-gohelper/install"
	"github.com/Go2Jv/yt-dlp-gohelper/ui"
)

func Run() int {
	lang := i18n.PromptLanguage(i18n.Zh)
	msg := i18n.NewMessages(lang)

	ui.PrintLogo()
	fmt.Println(msg.Welcome(runtime.GOOS))

	store := history.NewStore()
	if auto := loadAndPruneTasks(store); len(auto) > 0 {
		if shouldHandleNow(msg, auto) {
			_ = handleFailures(msg, store, runtime.GOOS)
		}
	}

	for {
		switch promptMainMenu(msg) {
		case 1:
			if !ensureDeps(msg) {
				continue
			}
			if err := runDownloadFlow(msg, store, runtime.GOOS); err != nil {
				fmt.Println(msg.Error(err))
			}
		case 2:
			_ = handleFailures(msg, store, runtime.GOOS)
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
		fmt.Println("2. 下载失败处理 (重试 / 删除 temp)")
		fmt.Println("0. 退出")
		choice := ui.Prompt("请选择 (默认1): ", "1")
		return parseChoice(choice, 2, 1)
	case i18n.Ja:
		fmt.Println("1. 動画をダウンロード")
		fmt.Println("2. 失敗タスク処理 (再試行 / temp削除)")
		fmt.Println("0. 終了")
		choice := ui.Prompt("選択 (デフォルト1): ", "1")
		return parseChoice(choice, 2, 1)
	default:
		fmt.Println("1. Download video")
		fmt.Println("2. Failed tasks (retry / delete temp)")
		fmt.Println("0. Exit")
		choice := ui.Prompt("Select (default 1): ", "1")
		return parseChoice(choice, 2, 1)
	}
}

func shouldHandleNow(msg *i18n.Messages, tasks []history.Task) bool {
	switch msg.Lang() {
	case i18n.Zh:
		fmt.Printf("检测到未完成任务: %d 个，是否现在处理？(y/n, 默认y): ", len(tasks))
		return ui.PromptYesNo("", true)
	case i18n.Ja:
		fmt.Printf("未完了のタスクが見つかりました: %d 件、今処理しますか？(y/n, デフォルトy): ", len(tasks))
		return ui.PromptYesNo("", true)
	default:
		fmt.Printf("Found unfinished tasks: %d. Handle now? (y/n, default y): ", len(tasks))
		return ui.PromptYesNo("", true)
	}
}

func runDownloadFlow(msg *i18n.Messages, store *history.Store, goos string) error {
	req := downloader.PromptRequest(msg, goos)
	task := history.Task{
		ID:         history.NewID(),
		URL:        req.URL,
		Args:       req.Args,
		OutputDir:  req.OutputDir,
		OutputTmpl: req.OutputTmpl,
		Status:     history.StatusRunning,
	}
	_ = store.Upsert(task)

	for {
		out, err := downloader.RunRequest(msg, req)
		if err == nil {
			_ = store.Delete(task.ID)
			fmt.Println(msg.Done())
			return nil
		}

		task.Status = history.StatusFailed
		task.Error = err.Error()
		task.FragmentsTotal = parseFragmentsTotal(out)
		if temps, e := history.FindTempFiles(task.OutputDir, task.CreatedAt); e == nil {
			task.TempFiles = temps
		}
		_ = store.Upsert(task)

		if !promptRetryAfterFailure(msg) {
			if len(task.TempFiles) == 0 {
				if temps, e := history.FindTempFiles(task.OutputDir, task.CreatedAt); e == nil {
					task.TempFiles = temps
				}
			}
			history.DeleteFiles(task.TempFiles)
			_ = store.Delete(task.ID)
			return err
		}
		task.Status = history.StatusRunning
		task.Error = ""
		_ = store.Upsert(task)
	}
}

func handleFailures(msg *i18n.Messages, store *history.Store, goos string) error {
	tasks := loadAndPruneTasks(store)
	if len(tasks) == 0 {
		switch msg.Lang() {
		case i18n.Zh:
			fmt.Println("没有需要处理的失败任务。")
		case i18n.Ja:
			fmt.Println("処理する失敗タスクはありません。")
		default:
			fmt.Println("No failed tasks to handle.")
		}
		return nil
	}

	switch msg.Lang() {
	case i18n.Zh:
		fmt.Println("\n--- 未完成任务 ---")
	case i18n.Ja:
		fmt.Println("\n--- 未完了タスク ---")
	default:
		fmt.Println("\n--- Unfinished tasks ---")
	}
	for i, t := range tasks {
		fmt.Printf("%d. %s\n", i+1, shortTaskLine(t))
	}

	def := "1"
	var prompt string
	switch msg.Lang() {
	case i18n.Zh:
		prompt = "请选择任务编号 (默认1，输入0返回): "
	case i18n.Ja:
		prompt = "タスク番号 (デフォルト1、0で戻る): "
	default:
		prompt = "Select task (default 1, 0 to back): "
	}
	choice := strings.TrimSpace(ui.Prompt(prompt, def))
	if choice == "0" {
		return nil
	}
	idx := parseChoice(choice, len(tasks), 1) - 1
	if idx < 0 || idx >= len(tasks) {
		return nil
	}
	t := tasks[idx]

	switch msg.Lang() {
	case i18n.Zh:
		fmt.Println("1. 重新下载")
		fmt.Println("2. 删除 temp 并清理记录")
		fmt.Println("0. 返回")
		choice = ui.Prompt("请选择 (默认1): ", "1")
	case i18n.Ja:
		fmt.Println("1. 再試行")
		fmt.Println("2. temp削除 + 記録削除")
		fmt.Println("0. 戻る")
		choice = ui.Prompt("選択 (デフォルト1): ", "1")
	default:
		fmt.Println("1. Retry")
		fmt.Println("2. Delete temp & remove record")
		fmt.Println("0. Back")
		choice = ui.Prompt("Select (default 1): ", "1")
	}

	action := strings.TrimSpace(choice)
	switch action {
	case "2":
		if len(t.TempFiles) == 0 {
			if temps, e := history.FindTempFiles(t.OutputDir, t.CreatedAt); e == nil {
				t.TempFiles = temps
			}
		}
		history.DeleteFiles(t.TempFiles)
		return store.Delete(t.ID)
	case "1", "":
		if !ensureDeps(msg) {
			return nil
		}
		return retryTask(msg, store, goos, t)
	default:
		return nil
	}
}

func retryTask(msg *i18n.Messages, store *history.Store, goos string, t history.Task) error {
	_ = goos
	req := downloader.Request{
		URL:        t.URL,
		Args:       t.Args,
		OutputDir:  t.OutputDir,
		OutputTmpl: t.OutputTmpl,
	}
	t.Status = history.StatusRunning
	t.Error = ""
	_ = store.Upsert(t)

	for {
		out, err := downloader.RunRequest(msg, req)
		if err == nil {
			_ = store.Delete(t.ID)
			fmt.Println(msg.Done())
			return nil
		}

		t.Status = history.StatusFailed
		t.Error = err.Error()
		t.FragmentsTotal = parseFragmentsTotal(out)
		if temps, e := history.FindTempFiles(t.OutputDir, t.CreatedAt); e == nil {
			t.TempFiles = temps
		}
		_ = store.Upsert(t)

		if !promptRetryAfterFailure(msg) {
			if len(t.TempFiles) == 0 {
				if temps, e := history.FindTempFiles(t.OutputDir, t.CreatedAt); e == nil {
					t.TempFiles = temps
				}
			}
			history.DeleteFiles(t.TempFiles)
			_ = store.Delete(t.ID)
			return err
		}
		t.Status = history.StatusRunning
		t.Error = ""
		_ = store.Upsert(t)
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

func shortTaskLine(t history.Task) string {
	s := t.URL
	if len(s) > 80 {
		s = s[:77] + "..."
	}
	if t.OutputDir != "" {
		s += " | " + t.OutputDir
	}
	return s
}

func loadAndPruneTasks(store *history.Store) []history.Task {
	tasks, err := store.Load()
	if err != nil {
		return nil
	}

	out := make([]history.Task, 0, len(tasks))
	for _, t := range tasks {
		if t.Status != history.StatusFailed && t.Status != history.StatusRunning {
			continue
		}
		if len(t.TempFiles) == 0 && t.OutputDir != "" {
			if temps, e := history.FindTempFiles(t.OutputDir, t.CreatedAt); e == nil {
				t.TempFiles = temps
				_ = store.Upsert(t)
			}
		}
		if len(t.TempFiles) > 0 && history.AnyExists(t.TempFiles) {
			out = append(out, t)
			continue
		}
		_ = store.Delete(t.ID)
	}
	return out
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

var reFragA = regexp.MustCompile(`Downloading fragment\s+\d+\s+of\s+(\d+)`)
var reFragB = regexp.MustCompile(`Downloading fragment\s+\d+\/(\d+)`)

func parseFragmentsTotal(output string) int {
	m := reFragA.FindAllStringSubmatch(output, -1)
	max := 0
	for _, mm := range m {
		if len(mm) != 2 {
			continue
		}
		v := atoi(mm[1])
		if v > max {
			max = v
		}
	}
	m = reFragB.FindAllStringSubmatch(output, -1)
	for _, mm := range m {
		if len(mm) != 2 {
			continue
		}
		v := atoi(mm[1])
		if v > max {
			max = v
		}
	}
	return max
}

func atoi(s string) int {
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0
		}
		n = n*10 + int(r-'0')
	}
	return n
}
