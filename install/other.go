//go:build !darwin && !windows && !linux
// +build !darwin,!windows,!linux

package install

import (
	"fmt"

	"github.com/Go2Jv/yt-dlp-gohelper/deps"
	"github.com/Go2Jv/yt-dlp-gohelper/i18n"
)

func Ensure(msg *i18n.Messages, state deps.State) error {
	return fmt.Errorf("%w: unsupported OS", ErrManualInstallRequired)
}
