//go:build linux
// +build linux

package install

import (
	"fmt"

	"github.com/Go2Jv/yt-dlp-gohelper/deps"
	"github.com/Go2Jv/yt-dlp-gohelper/execx"
	"github.com/Go2Jv/yt-dlp-gohelper/i18n"
)

func Ensure(msg *i18n.Messages, state deps.State) error {
	if !state.MissingAny() {
		return nil
	}

	if _, err := execx.LookPath("apt-get"); err == nil {
		if err := execx.RunInteractive("sudo", "apt-get", "update"); err != nil {
			return err
		}

		args := []string{"apt-get", "install", "-y"}
		if state.MissingYtDlp() {
			args = append(args, "yt-dlp")
		}
		if state.MissingFfmpeg() {
			args = append(args, "ffmpeg")
		}
		if err := execx.RunInteractive("sudo", args...); err != nil {
			return err
		}

		return nil
	}

	if state.MissingFfmpeg() {
		return fmt.Errorf("%w: apt-get not found (ffmpeg needs manual install)", ErrManualInstallRequired)
	}

	if _, err := execx.LookPath("curl"); err != nil {
		return fmt.Errorf("%w: curl not found", ErrManualInstallRequired)
	}

	outPath := "/usr/local/bin/yt-dlp"
	if err := execx.RunInteractive("sudo", "curl", "-L", "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp", "-o", outPath); err != nil {
		return err
	}
	if err := execx.RunInteractive("sudo", "chmod", "a+rx", outPath); err != nil {
		return err
	}

	post := deps.Check()
	if post.MissingYtDlp() {
		return fmt.Errorf("%w: installed but still not found in PATH (yt-dlp)", ErrManualInstallRequired)
	}
	if post.MissingFfmpeg() {
		return fmt.Errorf("%w: missing ffmpeg", ErrManualInstallRequired)
	}
	return nil
}
