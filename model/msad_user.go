package model

import "strings"

type MsAdUser struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Account string `json:"account"`
}

type MsAdUserCollection []*MsAdUser

func (x MsAdUserCollection) Len() int { return len(x) }
func (x MsAdUserCollection) Less(i, j int) bool {
	return strings.ToLower(x[i].Account) < strings.ToLower(x[j].Account)
}
func (x MsAdUserCollection) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
