package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func setupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		group, ctx := errgroup.WithContext(r.Context())

		var res1 int
		var res2 int

		group.Go(func() error {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.guildwars2.com/v2/build", nil)
			if err != nil {
				return err
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}

			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("bad status: %d", res.StatusCode)
			}

			var body buildApiResponse
			if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
				return err
			}

			res1 = body.Id
			return nil
		})

		group.Go(func() error {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://assetcdn.101.arenanetworks.com/latest/101", nil)
			if err != nil {
				return err
			}

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}

			defer res.Body.Close()

			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("bad status: %d", res.StatusCode)
			}

			b, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}

			versionNumStr, _, _ := strings.Cut(string(b), " ")
			versionNum, err := strconv.Atoi(versionNumStr)
			if err != nil {
				return err
			}

			res2 = versionNum
			return nil
		})

		if err := group.Wait(); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(`{"error":"Bad Gateway"}`))
			return
		}

		_ = json.NewEncoder(w).Encode(buildApiResponse{Id: max(res1, res2)})
	})

	return mux
}
