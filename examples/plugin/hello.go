// SPDX-License-Identifier: Apache-2.0

package hello

import (
	"fmt"
	"net/http"

	"github.com/RestGW/api-routerd/cmd/share"
)

type Hello struct {
	Cmd  string `json:"cmd"`
	Text string `json:"text"`
}

func (r *Hello) SayHello(rw http.ResponseWriter) error {
	fmt.Println(r)

	g := Hello{
		Cmd:  r.Cmd,
		Text: r.Text,
	}

	return share.JsonResponse(g, rw)
}
