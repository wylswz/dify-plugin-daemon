package service

import (
	"context"
	"errors"

	"github.com/langgenius/dify-plugin-daemon/internal/core/plugin_manager"
)

// SubmitToolInterruptOutOfSessionRequest is JSON: interrupt token and optional result map.
// Authenticate to the daemon with X-Api-Key (DIFY_PLUGIN_SERVER_KEY / same as SERVER_KEY).
type SubmitToolInterruptOutOfSessionRequest struct {
	Token  string         `json:"token" validate:"required"`
	Result map[string]any `json:"result"`
}

// SubmitToolInterruptResultOutOfSession forwards to Dify inner API (X-Inner-Api-Key) with body {token, result} only;
// Dify resolves tenant from the interrupt token. No plugin manifest lookup.
func SubmitToolInterruptResultOutOfSession(
	ctx context.Context,
	req *SubmitToolInterruptOutOfSessionRequest,
) (accepted bool, err error) {
	manager := plugin_manager.Manager()
	if manager == nil {
		return false, errors.New("plugin manager is not available")
	}
	inv := manager.BackwardsInvocation()
	inv.SetContext(ctx)
	if req.Result == nil {
		req.Result = map[string]any{}
	}
	data, err := inv.SubmitToolInterruptResultByTokenOnly(req.Token, req.Result)
	if err != nil {
		return false, err
	}
	if data == nil {
		return false, errors.New("dify response data is nil")
	}
	return data.Accepted, nil
}
