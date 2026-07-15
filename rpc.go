package main

import (
	"encoding/json"
	"fmt"
)

type rpcRequest struct {
	ID     string      `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type rpcError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type rpcResponse struct {
	ID     string           `json:"id"`
	Result *json.RawMessage `json:"result,omitempty"`
	Error  *rpcError        `json:"error,omitempty"`
}

type layoutExportResult struct {
	Layout struct {
		Root LayoutNode `json:"root"`
	} `json:"layout"`
}

type setSplitRatioParams struct {
	PaneID string  `json:"pane_id,omitempty"`
	TabID  string  `json:"tab_id,omitempty"`
	Path   []bool  `json:"path"`
	Ratio  float64 `json:"ratio"`
}

var seq int

func nextID() string {
	seq++
	return fmt.Sprintf("ec-%d", seq)
}

func callHerdr(tr Transport, method string, params interface{}) (*rpcResponse, error) {
	req := rpcRequest{
		ID:     nextID(),
		Method: method,
		Params: params,
	}
	raw, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}
	respRaw, err := tr.Call(raw)
	if err != nil {
		return nil, fmt.Errorf("transport call: %w", err)
	}
	var resp rpcResponse
	if err := json.Unmarshal(respRaw, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("herdr API error: %s: %s", resp.Error.Code, resp.Error.Message)
	}
	return &resp, nil
}

func exportLayout(tr Transport, paneID string) (*LayoutNode, error) {
	params := map[string]string{}
	if paneID != "" {
		params["pane_id"] = paneID
	}
	resp, err := callHerdr(tr, "layout.export", params)
	if err != nil {
		return nil, err
	}
	var result layoutExportResult
	if err := json.Unmarshal(*resp.Result, &result); err != nil {
		return nil, fmt.Errorf("parse layout.export result: %w", err)
	}
	return &result.Layout.Root, nil
}

func setSplitRatio(tr Transport, paneID string, path []bool, ratio float64) error {
	params := setSplitRatioParams{
		Path:  path,
		Ratio: ratio,
	}
	if paneID != "" {
		params.PaneID = paneID
	}
	_, err := callHerdr(tr, "layout.set_split_ratio", params)
	return err
}
