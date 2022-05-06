package api

type State struct {
	Code int            `json:"code,omitempty"`
	Data any            `json:"data"`
	Errs []string       `json:"errs,omitempty"`
	Meta map[string]any `json:"meta,omitempty"`
}
