package liner

import (
	"bytes"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

type coord struct {
	x, y int16
}
type smallRect struct {
	left, top, right, bottom int16
}

type consoleScreenBufferInfo struct {
	dwSize              coord
	dwCursorPosition    coord
	wAttributes         int16
	srWindow            smallRect
	dwMaximumWindowSize coord
}

func (s *State) cursorPos(x int) {
	var sbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(s.hOut), uintptr(unsafe.Pointer(&sbi)))
	procSetConsoleCursorPosition.Call(uintptr(s.hOut),
		uintptr(int(x)&0xFFFF|int(sbi.dwCursorPosition.y)<<16))
}

func (s *State) eraseLine() {
	var sbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(s.hOut), uintptr(unsafe.Pointer(&sbi)))
	var numWritten uint32
	procFillConsoleOutputCharacter.Call(uintptr(s.hOut), uintptr(' '),
		uintptr(sbi.dwSize.x-sbi.dwCursorPosition.x),
		uintptr(int(sbi.dwCursorPosition.x)&0xFFFF|int(sbi.dwCursorPosition.y)<<16),
		uintptr(unsafe.Pointer(&numWritten)))
}

func (s *State) eraseScreen() {
	var sbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(s.hOut), uintptr(unsafe.Pointer(&sbi)))
	var numWritten uint32
	procFillConsoleOutputCharacter.Call(uintptr(s.hOut), uintptr(' '),
		uintptr(sbi.dwSize.x)*uintptr(sbi.dwSize.y),
		0,
		uintptr(unsafe.Pointer(&numWritten)))
	procSetConsoleCursorPosition.Call(uintptr(s.hOut), 0)
}

func (s *State) moveUp(lines int) {
	var sbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(s.hOut), uintptr(unsafe.Pointer(&sbi)))
	procSetConsoleCursorPosition.Call(uintptr(s.hOut),
		uintptr(int(sbi.dwCursorPosition.x)&0xFFFF|(int(sbi.dwCursorPosition.y)-lines)<<16))
}

func (s *State) moveDown(lines int) {
	var sbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(s.hOut), uintptr(unsafe.Pointer(&sbi)))
	procSetConsoleCursorPosition.Call(uintptr(s.hOut),
		uintptr(int(sbi.dwCursorPosition.x)&0xFFFF|(int(sbi.dwCursorPosition.y)+lines)<<16))
}

func (s *State) emitNewLine() {
	// windows doesn't need to omit a new line
}

func (s *State) getColumns() {
	var sbi consoleScreenBufferInfo
	procGetConsoleScreenBufferInfo.Call(uintptr(s.hOut), uintptr(unsafe.Pointer(&sbi)))
	s.columns = int(sbi.dwSize.x)

	osVer := GetOSVer()
	if ver, err := strconv.Atoi(osVer); err != nil {
		if ver >= 10 {
			if s.columns > 1 {
				// Windows 10 needs a spare column for the cursor
				s.columns--
			}
		}
	}
}

/* Code stolen from 
   https://github.com/matishsiao/goInfo/blob/master/goInfo.go
  with minor modifications
*/
func GetOSVer() string {
	cmd := exec.Command("cmd","ver")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer 
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "0"
	}

	osStr := strings.Replace(out.String(),"\n","",-1)
	osStr = strings.Replace(osStr,"\r\n","",-1)

	r, _ := regexp.Compile("\\[(.*?) (.*?)\\]")
	matches := r.FindStringSubmatch(osStr)
	if len(matches) < 3 {
		return "0"
	}

	//only return the first part of the version
	ver := strings.Split(matches[2], ".")
	return ver[0]
} 
