package fileserve

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type RequestParams struct {
	ContactID uint64
	BankName  string
}

func Parse(_ context.Context, _ *http.Request, params httprouter.Params) (*RequestParams, error) {
	contactStr := params.ByName("contact")
	bank := params.ByName("bank")

	contactID, err := strconv.ParseUint(contactStr, 10, 64)
	if err != nil {
		return nil, err
	}

	return &RequestParams{
		ContactID: contactID,
		BankName:  bank,
	}, nil
}
