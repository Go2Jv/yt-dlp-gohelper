package app

import (
	"fmt"
	"os"
	"runtime"

	"github.com/Go2Jv/yt-dlp-gohelper/deps"
	"github.com/Go2Jv/yt-dlp-gohelper/downloader"
	"github.com/Go2Jv/yt-dlp-gohelper/i18n"
	"github.com/Go2Jv/yt-dlp-gohelper/install"
)

func Run() int {
	lang := i18n.DetectLanguage()
	lang = i18n.PromptLanguage(lang)
	msg := i18n.NewMessages(lang)

	fmt.Println(msg.Welcome(runtime.GOOS))

	state := deps.Check()
	if !state.MissingAny() {
		return runDownloader(msg)
	}

	fmt.Println(msg.DepsMissing(state))
	choice := msg.PromptInstallChoice()
	switch choice {
	case i18n.InstallChoiceAuto:
		if err := install.Ensure(msg, state); err != nil {
			fmt.Println(msg.Error(err))
			fmt.Println(msg.ManualInstallHint(runtime.GOOS))
			return 1
		}
	case i18n.InstallChoiceManual:
		fmt.Println(msg.ManualInstallHint(runtime.GOOS))
		return 1
	default:
		return 1
	}

	if err := deps.VerifyAll(); err != nil {
		fmt.Println(msg.Error(err))
		fmt.Println(msg.ManualInstallHint(runtime.GOOS))
		return 1
	}

	return runDownloader(msg)
}

func runDownloader(msg *i18n.Messages) int {
	if err := downloader.RunInteractive(msg, runtime.GOOS); err != nil {
		fmt.Println(msg.Error(err))
		return 1
	}
	fmt.Println(msg.Done())
	return 0
}

func Exit(code int) {
	os.Exit(code)
}
