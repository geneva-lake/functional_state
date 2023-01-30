package general

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

func MakeHTTPRequest[R any, A any](ctx context.Context, method, url string, request *R) (*A, error) {
	var buf *bytes.Buffer
	if request != nil {
		buf = new(bytes.Buffer)
		json.NewEncoder(buf).Encode(request)
	}
	r, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	answer := new(A)
	err = json.NewDecoder(res.Body).Decode(answer)
	if err != nil {
		return nil, err
	}
	return answer, nil
}
