package assist

import (
	"bytes"
	"fmt"
	"github.com/csby/gwin/model"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

const (
	dataTypeA = "A"
)

// Dns
// https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/dnscmd#dnscmd-enumrecords-command
type Dns struct {
	base

	ZoneName string
}

func (s *Dns) GetRecords() ([]*model.DnsRecord, error) {
	output, err := s.runCmd("/EnumRecords", s.ZoneName, ".", "/Type", dataTypeA, "/Child")
	if err != nil {
		return nil, err
	}

	return s.getRecords(output), nil
}

func (s *Dns) AddRecord(name, data string) error {
	_, err := s.runCmd("/RecordAdd", s.ZoneName, name, dataTypeA, data)
	return err
}

func (s *Dns) DeleteRecord(name, data string) error {
	_, err := s.runCmd("/RecordDelete", s.ZoneName, name, dataTypeA, data, "/f")
	return err
}

func (s *Dns) getRecords(text []byte) []*model.DnsRecord {
	results := make([]*model.DnsRecord, 0)
	if len(text) < 1 {
		return results
	}

	name := ""
	reader := &bytes.Buffer{}
	reader.Write(text)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if len(line) < 5 {
			continue
		}

		record := s.getRecord(name, line)
		if record != nil {
			results = append(results, record)
			name = record.Name
		}
	}

	return results
}

func (s *Dns) getRecord(name string, line string) *model.DnsRecord {
	if len(line) < 19 {
		return nil
	}

	// linux-dev	3600 A 192.168.123.201
	// win2016		3600 A 192.168.123.101
	// 				3600 A 172.16.22.182
	fields := make([]string, 0)
	items := strings.Split(strings.ReplaceAll(line, "\t", " "), " ")
	c := len(items)
	for i := 0; i < c; i++ {
		item := strings.TrimSpace(items[i])
		if len(item) < 1 {
			continue
		}
		fields = append(fields, item)
	}
	count := len(fields)
	if count < 3 {
		return nil
	}

	record := &model.DnsRecord{
		Name: name,
	}
	recordType := ""
	if count == 4 {
		time, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil
		}
		if time < 1 {
			return nil
		}
		record.Name = fields[0]
		recordType = fields[2]
		record.Data = fields[3]
	} else if count == 3 {
		time, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil
		}
		if time < 1 {
			return nil
		}
		recordType = fields[1]
		record.Data = fields[2]
	}
	if recordType != dataTypeA {
		return nil
	}

	return record
}

func (s *Dns) runCmd(args ...string) ([]byte, error) {
	buf := &bytes.Buffer{}
	cmd := exec.Command("dnscmd", args...)
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
