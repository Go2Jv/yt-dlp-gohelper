package deps

import (
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/Go2Jv/yt-dlp-gohelper/execx"
)

type State struct {
	YtDlp  bool
	Ffmpeg bool
}

func (s State) MissingAny() bool {
	return !s.YtDlp || !s.Ffmpeg
}

func (s State) MissingYtDlp() bool {
	return !s.YtDlp
}

func (s State) MissingFfmpeg() bool {
	return !s.Ffmpeg
}

func Check() State {
	_, ytErr := execx.LookPath("yt-dlp")
	_, ffErr := execx.LookPath("ffmpeg")
	return State{
		YtDlp:  ytErr == nil,
		Ffmpeg: ffErr == nil,
	}
}

func VerifyAll() error {
	if err := verifyOne("yt-dlp"); err != nil {
		return err
	}
	if err := verifyOne("ffmpeg"); err != nil {
		return err
	}
	return nil
}

func verifyOne(bin string) error {
	_, err := execx.LookPath(bin)
	if err != nil {
		return fmt.Errorf("%s not found in PATH", bin)
	}
	out, runErr := execx.RunCombined(15*time.Second, bin, "-version")
	if runErr != nil && bin == "yt-dlp" {
		out, runErr = execx.RunCombined(15*time.Second, bin, "--version")
	}
	if runErr != nil {
		if errors.Is(runErr, exec.ErrNotFound) {
			return fmt.Errorf("%s not found in PATH", bin)
		}
		return fmt.Errorf("%s exists but cannot run: %v\n%s", bin, runErr, out)
	}
	return nil
}
