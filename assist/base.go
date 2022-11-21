package assist

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os/exec"
	"strings"
	"syscall"
)

type base struct {
}

func (s *base) runShell(arg ...string) ([]byte, error) {
	args := append([]string{"-nologo", "-noprofile"}, arg...)
	buf := &bytes.Buffer{}
	cmd := exec.Command("powershell", args...)
	cmd.Stdout = buf
	cmd.Stderr = buf

	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		output, _ := s.toUtf8(buf.Bytes())
		return nil, fmt.Errorf("%s", strings.TrimSpace(string(output)))
	}

	return s.toUtf8(buf.Bytes())
}

func (s *base) toUtf8(v []byte) ([]byte, error) {
	return simplifiedchinese.GB18030.NewDecoder().Bytes(v)
}

func (s *base) getExitCode(err error) int {
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exit.Sys().(syscall.WaitStatus); ok {
				return status.ExitStatus()
			}
		}
	}

	return 0
}

func (s *base) getFields(line string, sep string) []string {
	fields := make([]string, 0)
	items := strings.Split(strings.ReplaceAll(line, "\t", " "), sep)
	c := len(items)
	for i := 0; i < c; i++ {
		item := strings.TrimSpace(items[i])
		if len(item) < 1 {
			continue
		}
		fields = append(fields, item)
	}

	return fields
}

func (s *base) uniqueId(id string, a ...interface{}) string {
	o := fmt.Sprint(a...)
	v := fmt.Sprintf("%s-%s", id, o)
	h := md5.New()

	_, err := io.WriteString(h, v)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
