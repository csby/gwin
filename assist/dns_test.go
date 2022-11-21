package assist

import (
	"testing"
)

func TestDns_getRecords(t *testing.T) {
	dns := &Dns{
		ZoneName: "csby.fun",
	}
	output, err := dns.runCmd("/EnumRecords", dns.ZoneName, ".", "/Type", dataTypeA, "/Child")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("output:")
	t.Log(string(output))

	results := dns.getRecords(output)
	c := len(results)
	t.Log("count: ", c)
	for i := 0; i < c; i++ {
		item := results[i]
		t.Logf("%3d %14s %s", i+1, item.Name, item.Data)
	}
}
