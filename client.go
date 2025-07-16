package zabbix_go

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	apiURL     string
	apiToken   string
	httpClient *http.Client
}

func NewClient(zabbixURL, apiToken string) (*Client, error) {
	if zabbixURL == "" || apiToken == "" {
		return nil, fmt.Errorf("URL da API Zabbix e Token não podem ser vazios")
	}

	c := &Client{
		apiURL:   zabbixURL,
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Valida a conexão e o token
	if _, err := c.GetAPIInfo(context.Background()); err != nil {
		return nil, fmt.Errorf("falha ao validar a sessão Zabbix: %w", err)
	}

	return c, nil
}

func (c *Client) do(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	requestBody := RPCRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		Auth:    c.apiToken,
		ID:      1,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("falha ao converter requisição para JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("falha ao criar requisição HTTP: %w", err)
	}
	req.Header.Set("Content-Type", "application/json-rpc")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("falha ao executar requisição HTTP: %w", err)
	}
	defer resp.Body.Close()

	var rpcResponse RPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&rpcResponse); err != nil {
		return nil, fmt.Errorf("falha ao decodificar resposta JSON-RPC: %w", err)
	}

	if rpcResponse.Error != nil {
		return nil, rpcResponse.Error
	}

	return rpcResponse.Result, nil
}

// --- Implementação dos Métodos ---

func (c *Client) GetAPIInfo(ctx context.Context) (string, error) {
	result, err := c.do(ctx, "apiinfo.version", map[string]string{})
	if err != nil {
		return "", err
	}
	var version string
	if err := json.Unmarshal(result, &version); err != nil {
		return "", err
	}
	return version, nil
}

func (c *Client) ListHostGroups(ctx context.Context) ([]HostGroup, error) {
	params := map[string]interface{}{"output": "extend", "sortfield": "name"}
	result, err := c.do(ctx, "hostgroup.get", params)
	if err != nil {
		return nil, err
	}
	var groups []HostGroup
	if err := json.Unmarshal(result, &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

func (c *Client) ListHostsByGroupID(ctx context.Context, groupIDs []string) ([]Host, error) {
	params := map[string]interface{}{"output": []string{"hostid", "host", "name"}, "groupids": groupIDs, "sortfield": "name"}
	result, err := c.do(ctx, "host.get", params)
	if err != nil {
		return nil, err
	}
	var hosts []Host
	if err := json.Unmarshal(result, &hosts); err != nil {
		return nil, err
	}
	return hosts, nil
}

func (c *Client) ListItemsByHostID(ctx context.Context, hostIDs []string) ([]Item, error) {
	params := map[string]interface{}{"output": "extend", "hostids": hostIDs, "sortfield": "name"}
	result, err := c.do(ctx, "item.get", params)
	if err != nil {
		return nil, err
	}
	var items []Item
	if err := json.Unmarshal(result, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (c *Client) ListRecentAlertsByHostID(ctx context.Context, hostIDs []string) ([]Alert, error) {
	params := map[string]interface{}{
		"output":            "extend",
		"hostids":           hostIDs,
		"sortfield":         []string{"lastchange"},
		"sortorder":         "DESC",
		"filter":            map[string]string{"value": "1"}, // Apenas triggers em estado de PROBLEMA
		"selectHosts":       []string{"host"},
		"skipDependent":     "true",
		"expandDescription": "true",
	}
	result, err := c.do(ctx, "trigger.get", params)
	if err != nil {
		return nil, err
	}
	var alerts []Alert
	if err := json.Unmarshal(result, &alerts); err != nil {
		return nil, err
	}
	return alerts, nil
}
