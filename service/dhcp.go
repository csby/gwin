package main

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

type Dhcp struct {
}

func (s *Dhcp) GetFilters() ([]*Filter, error) {
	output, err := s.runCmd("show", "filter")
	if err != nil {
		return nil, err
	}

	return s.getFilters(output), err
}

func (s *Dhcp) AddFilter(v *Filter) error {
	if v == nil {
		return fmt.Errorf("filter is nil")
	}

	allow := "deny"
	if v.Allow {
		allow = "allow"
	}
	address := v.Address
	comment := v.Comment

	_, err := s.runCmd("add", "filter", allow, address, comment)
	return err
}

func (s *Dhcp) DeleteFilter(address string) error {
	if len(address) < 1 {
		return fmt.Errorf("address is empty")
	}

	_, err := s.runCmd("delete", "filter", address)
	return err
}

func (s *Dhcp) getFilters(text []byte) []*Filter {
	results := make([]*Filter, 0)
	if len(text) < 1 {
		return results
	}

	allow := false
	reader := &bytes.Buffer{}
	reader.Write(text)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if len(line) < 7 {
			continue
		}
		if strings.Index(line, "允许列表中所有") != -1 {
			allow = true
		} else if strings.Index(line, "拒绝列表中所有") != -1 {
			allow = false
		}

		filter := s.getFilter(allow, line)
		if filter != nil {
			results = append(results, filter)
		}
	}

	return results
}

func (s *Dhcp) getFilter(allow bool, line string) *Filter {
	if len(line) < 19 {
		return nil
	}

	fields := make([]string, 0)
	items := strings.Split(line, "\t")
	c := len(items)
	for i := 0; i < c; i++ {
		item := strings.TrimSpace(items[i])
		if len(item) < 1 {
			continue
		}
		fields = append(fields, item)
	}
	if len(fields) < 2 {
		return nil
	}

	index, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil
	}
	if index < 1 {
		return nil
	}

	address := fields[1]
	if len(address) != 17 {
		return nil
	}

	comment := ""
	if len(fields) > 2 {
		comment = fields[2]
	}

	return &Filter{
		Allow:   allow,
		Address: address,
		Comment: comment,
	}
}

func (s *Dhcp) runCmd(arg ...string) ([]byte, error) {
	args := append([]string{"dhcp", "server", "V4"}, arg...)
	buf := &bytes.Buffer{}
	cmd := exec.Command("netsh", args...)
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

func (s *Dhcp) toUtf8(v []byte) ([]byte, error) {
	return simplifiedchinese.GB18030.NewDecoder().Bytes(v)
}

func (s *Dhcp) getExitCode(err error) int {
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
