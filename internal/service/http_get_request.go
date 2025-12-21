package engine

import (
	"context"
	"net/http"

	"github.com/ManoloEsS/burrow/internal/domain"
)

func ReqFromStruct(ctx context.Context, reqStruct *domain.Request) (*http.Response, error) {
	request, err := http.NewRequestWithContext(context.Background(), reqStruct.Method, reqStruct.URL, nil)
	if err != nil {
		return nil, err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}
