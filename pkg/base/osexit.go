package base

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
)

func OsExitAndWaitPressIfWindows(code int) {
	if runtime.GOOS == "windows" {
		_, _ = fmt.Fprintf(os.Stderr, "Press Enter to exit...")
		r := bufio.NewReader(os.Stdin)
		_, _ = r.ReadByte()
	}
	os.Exit(code)
}
