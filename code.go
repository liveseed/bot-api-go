package bot

import (
	"context"
	"encoding/json"
)

func ReadMultisigByCode(ctx context.Context, codeId string) (*MultisigRequest, error) {
	body, err := Request(ctx, "GET", "/codes/"+codeId, nil, "")
	if err != nil {
		return nil, ServerError(ctx, err)
	}
	var resp struct {
		Data  *MultisigRequest `json:"data"`
		Error Error            `json:"error"`
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, BadDataError(ctx)
	}
	if resp.Error.Code > 0 {
		return nil, resp.Error
	}
	return resp.Data, nil
}