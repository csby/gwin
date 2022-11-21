package assist

import (
	"crypto/tls"
	"fmt"
	"github.com/csby/gwin/model"
	"github.com/go-ldap/ldap/v3"
	"sort"
	"strings"
)

type MsAd struct {
	Host     string
	Port     int
	Base     string
	Account  string
	Password string
}

func (s *MsAd) GetAllUsers() ([]*model.MsAdUser, error) {
	l, err := s.open(true)
	if err != nil {
		return nil, err
	}
	defer l.Close()

	filter := fmt.Sprintf("(&(&(objectCategory=%s)(objectClass=%s)))", "Person", "user")
	attributes := []string{"dn", "objectSid", "sAMAccountName", "displayName"}
	request := ldap.NewSearchRequest(
		s.Base,
		ldap.ScopeWholeSubtree, ldap.DerefAlways, 0, 0, false,
		filter,
		attributes,
		nil)

	rs, err := l.Search(request)
	if err != nil {
		return nil, err
	}

	items := make(model.MsAdUserCollection, 0)
	c := len(rs.Entries)
	for i := 0; i < c; i++ {
		entry := rs.Entries[i]
		if entry == nil {
			continue
		}

		item := &model.MsAdUser{}
		item.Account = entry.GetAttributeValue("sAMAccountName")
		item.Id = s.decodeSID(entry.GetRawAttributeValue("objectSid"))
		item.Name = entry.GetAttributeValue("displayName")
		if len(item.Name) < 1 {
			item.Name = item.Account
		}

		items = append(items, item)
	}

	sort.Sort(items)
	return items, nil
}

func (s *MsAd) open(bind bool) (*ldap.Conn, error) {
	var (
		conn *ldap.Conn
		err  error
	)
	server := fmt.Sprintf("%s:%d", s.Host, s.Port)
	//conn, err = ldap.DialURL(server)
	if s.Port == 636 {
		conn, err = ldap.DialTLS("tcp", server, &tls.Config{InsecureSkipVerify: true})
	} else {
		conn, err = ldap.Dial("tcp", server)
	}
	if err != nil {
		return nil, err
	}

	if bind {
		err = conn.Bind(s.Account, s.Password)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}

	return conn, nil
}

func (s *MsAd) decodeSID(sid []byte) string {
	if len(sid) < 28 {
		return ""
	}
	strSid := strings.Builder{}
	strSid.WriteString("S-")

	revision := int(sid[0])
	strSid.WriteString(fmt.Sprint(revision))

	countSubAuths := int(sid[1] & 0xFF)
	authority := int(0)
	for i := 2; i <= 7; i++ {
		shift := uint(8 * (5 - (i - 2)))
		authority |= int(sid[i]) << shift
	}
	strSid.WriteString("-")
	strSid.WriteString(fmt.Sprintf("%x", authority))

	offset := 8
	size := 4
	for j := 0; j < countSubAuths; j++ {
		subAuthority := 0
		for k := 0; k < size; k++ {
			subAuthority |= (int(sid[offset+k]) & 0xFF) << uint(8*k)
		}
		strSid.WriteString("-")
		strSid.WriteString(fmt.Sprint(subAuthority))
		offset += size
	}

	return strSid.String()
}
