package zabbix_go

import (
	"encoding/json"
	"fmt"
)

type RPCRequest struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Auth    string      `json:"auth,omitempty"`
	ID      int         `json:"id"`
}

type RPCResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
	ID      int             `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (e *RPCError) Error() string {
	return fmt.Sprintf("API Error: %s (Code: %d) - %s", e.Message, e.Code, e.Data)
}

type HostGroup struct {
	ID   string `json:"groupid"`
	Name string `json:"name"`
}

type Host struct {
	ID   string `json:"hostid"`
	Host string `json:"host"`
	Name string `json:"name"`
}

type Item struct {
	ID        string `json:"itemid"`
	Name      string `json:"name"`
	Key       string `json:"key_"`
	LastValue string `json:"lastvalue"`
	LastClock string `json:"lastclock"`
}

type Alert struct {
	ID          string `json:"triggerid"` // Em Zabbix, alertas vêm de Triggers
	Description string `json:"description"`
	Priority    string `json:"priority"`
	LastChange  string `json:"lastchange"`
	Value       string `json:"value"` // 1 = PROBLEMA, 0 = OK
	State       string `json:"state"` // 1 = ATUALIZADO RECENTEMENTE, 0 = DESCONHECIDO
}
