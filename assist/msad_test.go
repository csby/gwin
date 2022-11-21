package assist

import (
	"testing"
)

func TestMsAd_GetAllUsers(t *testing.T) {
	ad := &MsAd{
		Host:     "localhost",
		Port:     636,
		Base:     "dc=csby,dc=fun",
		Account:  "CN=Administrator,CN=Users,DC=example,DC=com",
		Password: "***",
	}

	items, err := ad.GetAllUsers()
	if err != nil {
		t.Error(err)
		return
	}
	c := len(items)
	t.Log("count: ", c)
	for i := 0; i < c; i++ {
		item := items[i]
		if item == nil {
			continue
		}
		t.Logf("%3d  %20s  %s", i+1, item.Name, item.Id)
	}
}
