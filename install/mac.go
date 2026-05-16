//go:build darwin
// +build darwin

package install

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Go2Jv/yt-dlp-gohelper/deps"
	"github.com/Go2Jv/yt-dlp-gohelper/execx"
	"github.com/Go2Jv/yt-dlp-gohelper/i18n"
	"github.com/Go2Jv/yt-dlp-gohelper/ui"
)

func Ensure(msg *i18n.Messages, state deps.State) error {
	brew := findBrew()
	if brew == "" {
		switch msg.Lang() {
		case i18n.Zh:
			fmt.Println("未检测到 Homebrew (macOS / 苹果系统)")
		case i18n.Ja:
			fmt.Println("Homebrew が見つかりません (macOS)")
		default:
			fmt.Println("Homebrew not found (macOS)")
		}

		choice := ui.Prompt("Install Homebrew now? (y/N): ", "N")
		if choice != "y" && choice != "Y" {
			return fmt.Errorf("%w: homebrew missing", ErrManualInstallRequired)
		}
		if err := execx.RunInteractive("/bin/bash", "-c", `$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)`); err != nil {
			return err
		}
		brew = findBrew()
		if brew == "" {
			return fmt.Errorf("%w: brew installed but not found", ErrManualInstallRequired)
		}
	}

	if state.MissingAny() {
		args := []string{"install"}
		if state.MissingYtDlp() {
			args = append(args, "yt-dlp")
		}
		if state.MissingFfmpeg() {
			args = append(args, "ffmpeg")
		}
		if err := execx.RunInteractive(brew, args...); err != nil {
			return err
		}
	}

	return nil
}

func findBrew() string {
	if p, err := execx.LookPath("brew"); err == nil {
		return p
	}
	candidates := []string{
		"/opt/homebrew/bin/brew",
		"/usr/local/bin/brew",
	}
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			return filepath.Clean(c)
		}
	}
	return ""
}
