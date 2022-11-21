package assist

import (
	"bytes"
	"fmt"
	"github.com/csby/gwin/model"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type Dhcp struct {
	base
}

func (s *Dhcp) GetFilters() ([]*model.DhcpFilter, error) {
	output, err := s.runCmd("show", "filter")
	if err != nil {
		return nil, err
	}

	return s.getFilters(output), err
}

func (s *Dhcp) AddFilter(v *model.DhcpFilter) error {
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

func (s *Dhcp) GetLeases() ([]*model.DhcpLease, error) {
	output, err := s.runShell("Get-DhcpServerV4Scope", "|", "Select", "ScopeId")
	if err != nil {
		return nil, err
	}
	results := make([]*model.DhcpLease, 0)
	scopes := s.getScopeIds(output)
	sc := len(scopes)
	if sc < 1 {
		return results, nil
	}

	comments := make(map[string]string)
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		fs, fe := s.GetFilters()
		if fe != nil {
			return
		}
		fc := len(fs)
		for fi := 0; fi < fc; fi++ {
			f := fs[fi]
			if f == nil {
				continue
			}
			comments[f.Address] = f.Comment
		}
	}()

	leases := make([][]*model.DhcpLease, 0)
	for si := 0; si < sc; si++ {
		scope := scopes[si]
		if len(scope) < 1 {
			continue
		}

		ls := make([]*model.DhcpLease, 0)
		li := len(leases)
		leases = append(leases, ls)
		waitGroup.Add(1)
		go func(scopeId string, index int) {
			defer waitGroup.Done()
			o, e := s.runShell("Get-DhcpServerV4Lease", "-ScopeId", scope, "|", "Select", "ClientId,IPAddress")
			if e != nil {
				return
			}
			leases[index] = s.getLeases(o)
		}(scope, li)
	}

	waitGroup.Wait()
	for gi := 0; gi < len(leases); gi++ {
		gs := leases[gi]
		for ri := 0; ri < len(gs); ri++ {
			r := gs[ri]
			if r == nil {
				continue
			}
			comment, ok := comments[r.Address]
			if ok {
				r.Comment = comment
			}
			results = append(results, r)
		}
	}

	return results, err
}

func (s *Dhcp) getFilters(text []byte) []*model.DhcpFilter {
	results := make([]*model.DhcpFilter, 0)
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

func (s *Dhcp) getFilter(allow bool, line string) *model.DhcpFilter {
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

	return &model.DhcpFilter{
		Allow:   allow,
		Address: strings.ToUpper(address),
		Comment: comment,
	}
}

func (s Dhcp) getScopeIds(text []byte) []string {
	results := make([]string, 0)
	if len(text) < 1 {
		return results
	}
	/*
		ScopeId
		-------
		172.16.11.0
		172.16.12.0
	*/

	reader := &bytes.Buffer{}
	reader.Write(text)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if len(line) < 8 {
			continue
		}
		if strings.Index(line, "ScopeId") == 0 ||
			strings.Index(line, "---") == 0 {
			continue
		}

		results = append(results, strings.TrimSpace(line))
	}

	return results
}

func (s Dhcp) getLeases(text []byte) []*model.DhcpLease {
	results := make([]*model.DhcpLease, 0)
	if len(text) < 1 {
		return results
	}
	/*
		ClientId          IPAddress
		--------          ---------
		90-94-97-8b-f5-f8 172.16.11.19
		9c-b6-d0-e8-38-47 172.16.11.20
	*/

	reader := &bytes.Buffer{}
	reader.Write(text)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if len(line) < 19 {
			continue
		}
		if strings.Index(line, "ClientId") == 0 ||
			strings.Index(line, "---") == 0 {
			continue
		}

		lease := s.getLease(line)
		if lease != nil {
			results = append(results, lease)
		}
	}

	return results
}

func (s *Dhcp) getLease(line string) *model.DhcpLease {
	if len(line) < 19 {
		return nil
	}
	/*
		90-94-97-8b-f5-f8 172.16.11.19
	*/

	fields := s.getFields(line, " ")
	if len(fields) < 2 {
		return nil
	}

	address := fields[0]
	if len(address) != 17 {
		return nil
	}

	return &model.DhcpLease{
		IpV4:    fields[1],
		Address: strings.ToUpper(address),
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
