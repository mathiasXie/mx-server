package handler

import (
	"fmt"
	"net/http"
)

func (h *ChatHandler) Authenticate(header http.Header) error {

	fmt.Printf("%+v\n", header)
	// 进行认证
	return nil
}
