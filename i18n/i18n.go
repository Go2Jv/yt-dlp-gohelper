package i18n

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/Go2Jv/yt-dlp-gohelper/deps"
	"github.com/Go2Jv/yt-dlp-gohelper/ui"
)

type Lang string

const (
	Zh Lang = "zh"
	En Lang = "en"
	Ja Lang = "ja"
)

type InstallChoice string

const (
	InstallChoiceAuto   InstallChoice = "auto"
	InstallChoiceManual InstallChoice = "manual"
)

func DetectLanguage() Lang {
	s := strings.ToLower(strings.Join([]string{
		os.Getenv("LC_ALL"),
		os.Getenv("LC_MESSAGES"),
		os.Getenv("LANG"),
	}, " "))
	switch {
	case strings.Contains(s, "zh"):
		return Zh
	case strings.Contains(s, "ja"):
		return Ja
	default:
		return En
	}
}

func PromptLanguage(defaultLang Lang) Lang {
	d := defaultLang
	if runtime.GOOS == "windows" && os.Getenv("LANG") == "" && os.Getenv("LC_ALL") == "" {
		d = En
	}

	var def string
	switch d {
	case Zh:
		def = "1"
	case En:
		def = "2"
	case Ja:
		def = "3"
	default:
		def = "1"
	}

	fmt.Println("Language / 语言 / 言語")
	fmt.Println("1. 中文")
	fmt.Println("2. English")
	fmt.Println("3. 日本語")
	choice := ui.Prompt("Select (default "+def+"): ", def)
	switch strings.TrimSpace(choice) {
	case "1":
		return Zh
	case "2":
		return En
	case "3":
		return Ja
	default:
		return d
	}
}

type Messages struct {
	lang Lang
}

func NewMessages(lang Lang) *Messages {
	return &Messages{lang: lang}
}

func (m *Messages) Lang() Lang {
	return m.lang
}

func (m *Messages) Welcome(goos string) string {
	osName := DisplayOSName(goos)
	switch m.lang {
	case Zh:
		return "========== yt-dlp-gohelper ==========\n系统: " + osName + "\n"
	case Ja:
		return "========== yt-dlp-gohelper ==========\nOS: " + osName + "\n"
	default:
		return "========== yt-dlp-gohelper ==========\nOS: " + osName + "\n"
	}
}

func (m *Messages) Done() string {
	switch m.lang {
	case Zh:
		return "完成"
	case Ja:
		return "完了"
	default:
		return "Done"
	}
}

func (m *Messages) Error(err error) string {
	switch m.lang {
	case Zh:
		return "错误: " + err.Error()
	case Ja:
		return "エラー: " + err.Error()
	default:
		return "Error: " + err.Error()
	}
}

func (m *Messages) DepsMissing(state deps.State) string {
	var missing []string
	if state.MissingYtDlp() {
		missing = append(missing, "yt-dlp")
	}
	if state.MissingFfmpeg() {
		missing = append(missing, "ffmpeg")
	}

	switch m.lang {
	case Zh:
		return "检测到依赖缺失: " + strings.Join(missing, ", ")
	case Ja:
		return "依存関係が見つかりません: " + strings.Join(missing, ", ")
	default:
		return "Missing dependencies: " + strings.Join(missing, ", ")
	}
}

func (m *Messages) PromptInstallChoice() InstallChoice {
	switch m.lang {
	case Zh:
		fmt.Println("1. 自动安装缺失依赖")
		fmt.Println("2. 退出并自行安装 (Windows 记得配置 PATH)")
		choice := ui.Prompt("请选择 (默认1): ", "1")
		if strings.TrimSpace(choice) == "2" {
			return InstallChoiceManual
		}
		return InstallChoiceAuto
	case Ja:
		fmt.Println("1. 不足している依存関係を自動インストール")
		fmt.Println("2. 終了して自分でインストール (Windows は PATH を設定)")
		choice := ui.Prompt("選択 (デフォルト1): ", "1")
		if strings.TrimSpace(choice) == "2" {
			return InstallChoiceManual
		}
		return InstallChoiceAuto
	default:
		fmt.Println("1. Install missing dependencies")
		fmt.Println("2. Exit and install manually (Windows: remember PATH)")
		choice := ui.Prompt("Select (default 1): ", "1")
		if strings.TrimSpace(choice) == "2" {
			return InstallChoiceManual
		}
		return InstallChoiceAuto
	}
}

func (m *Messages) ManualInstallHint(goos string) string {
	switch goos {
	case "darwin":
		switch m.lang {
		case Zh:
			return "请先安装 Homebrew，然后执行: brew install yt-dlp ffmpeg"
		case Ja:
			return "Homebrew をインストールしてから実行: brew install yt-dlp ffmpeg"
		default:
			return "Install Homebrew first, then: brew install yt-dlp ffmpeg"
		}
	case "linux":
		switch m.lang {
		case Zh:
			return "请先安装 yt-dlp 与 ffmpeg (优先 apt): sudo apt-get install -y yt-dlp ffmpeg"
		case Ja:
			return "yt-dlp と ffmpeg をインストールしてください (apt 推奨): sudo apt-get install -y yt-dlp ffmpeg"
		default:
			return "Install yt-dlp and ffmpeg (prefer apt): sudo apt-get install -y yt-dlp ffmpeg"
		}
	case "windows":
		switch m.lang {
		case Zh:
			return "请先安装 yt-dlp 与 ffmpeg，并确保能在命令行直接运行 yt-dlp/ffmpeg (配置 PATH)"
		case Ja:
			return "yt-dlp と ffmpeg をインストールし、コマンドラインで yt-dlp/ffmpeg を実行できるよう PATH を設定してください"
		default:
			return "Install yt-dlp and ffmpeg, and ensure yt-dlp/ffmpeg are runnable in CMD/PowerShell (PATH)"
		}
	default:
		return ""
	}
}

func DisplayOSName(goos string) string {
	switch goos {
	case "darwin":
		return "macOS / 苹果系统"
	case "windows":
		return "Windows"
	case "linux":
		return "Linux"
	default:
		return goos
	}
}
