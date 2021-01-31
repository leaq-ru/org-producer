package healthz

import (
	"net/http"
	"strings"
)

func (h Healthz) Serve() (err error) {
	mux := http.NewServeMux()
	mux.Handle("/healthz", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, e := w.Write(nil)
		if e != nil {
			h.Logger.Error().Err(e).Send()
		}
	}))

	return http.ListenAndServe(strings.Join([]string{
		"0.0.0.0",
		h.Port,
	}, ":"), mux)
}
