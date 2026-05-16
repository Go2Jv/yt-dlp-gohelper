//go:build windows
// +build windows

package install

import (
	"fmt"
	"strings"

	"github.com/Go2Jv/yt-dlp-gohelper/deps"
	"github.com/Go2Jv/yt-dlp-gohelper/execx"
	"github.com/Go2Jv/yt-dlp-gohelper/i18n"
)

func Ensure(msg *i18n.Messages, state deps.State) error {
	if _, err := execx.LookPath("winget"); err != nil {
		return fmt.Errorf("%w: winget not found", ErrManualInstallRequired)
	}

	if state.MissingYtDlp() {
		if err := execx.RunInteractive("winget", "install", "--id", "yt-dlp.yt-dlp", "-e", "--source", "winget", "--accept-source-agreements", "--accept-package-agreements"); err != nil {
			return err
		}
	}
	if state.MissingFfmpeg() {
		if err := execx.RunInteractive("winget", "install", "--id", "Gyan.FFmpeg", "-e", "--source", "winget", "--accept-source-agreements", "--accept-package-agreements"); err != nil {
			return err
		}
	}

	post := deps.Check()
	if post.MissingAny() {
		var m []string
		if post.MissingYtDlp() {
			m = append(m, "yt-dlp")
		}
		if post.MissingFfmpeg() {
			m = append(m, "ffmpeg")
		}
		return fmt.Errorf("%w: installed but still not found in PATH (%s)", ErrManualInstallRequired, strings.Join(m, ", "))
	}
	return nil
}
