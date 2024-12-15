package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"os"
)

const (
	maxDbBatch = 100 // ensure we test cursoring
)

func NewDropBoxClient(ctx context.Context) (*DropBoxClient, error) {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("DROPBOX_CLIENT_ID"),
		ClientSecret: os.Getenv("DROPBOX_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://api.dropboxapi.com/oauth2/token",
		},
	}

	tok := &oauth2.Token{
		RefreshToken: os.Getenv("DROPBOX_REFRESH_TOKEN"),
	}

	// make sure it works
	ts := conf.TokenSource(ctx, tok)
	if _, err := ts.Token(); err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	oauthCl := oauth2.NewClient(nil, ts)
	retryCl := retryablehttp.NewClient()
	retryCl.HTTPClient = oauthCl
	retryCl.Logger = nil
	return &DropBoxClient{cl: retryCl}, nil
}

type DropBoxClient struct {
	cl *retryablehttp.Client
}

type File struct {
	Tag            string `json:".tag"`
	Name           string `json:"name"`
	PathLower      string `json:"path_lower"`
	PathDisplay    string `json:"path_display"`
	ID             string `json:"id"`
	ClientModified string `json:"client_modified"`
	ServerModified string `json:"server_modified"`
	Rev            string `json:"rev"`
	Size           int    `json:"size"`
	IsDownloadable bool   `json:"is_downloadable"`
	ContentHash    string `json:"content_hash"`
}

func (db *DropBoxClient) IterateFiles(ctx context.Context, limit int, cb func(context.Context, File) error) error {
	type ListRequest struct {
		Path                        string `json:"path"`
		Recursive                   bool   `json:"recursive"`
		Limit                       int    `json:"limit"`
		IncludeNonDownloadableFiles bool   `json:"include_non_downloadable_files"`
	}

	type ListResponse struct {
		Entries []File `json:"entries"`
		Cursor  string `json:"cursor"`
		HasMore bool   `json:"has_more"`
	}

	type ContinueRequest struct {
		Cursor string `json:"cursor"`
	}

	req := &ListRequest{
		Path:                        "",
		Recursive:                   true,
		Limit:                       min(limit, maxDbBatch),
		IncludeNonDownloadableFiles: false,
	}
	var rsp ListResponse
	err := db.post(ctx, "/files/list_folder", req, &rsp)
	if err != nil {
		return fmt.Errorf("dropbox list files: %w", err)
	}

	for {
		for _, e := range rsp.Entries {
			if err := cb(ctx, e); err != nil {
				return fmt.Errorf("callback error: %w", err)
			}
			limit--
			if limit == 0 {
				return nil
			}
		}

		if !rsp.HasMore {
			return nil
		}

		req := &ContinueRequest{Cursor: rsp.Cursor}
		rsp = ListResponse{}
		err := db.post(ctx, "/files/list_folder/continue", req, &rsp)
		if err != nil {
			return fmt.Errorf("dropbox iterate files: %w", err)
		}
	}
}

func (db *DropBoxClient) post(ctx context.Context, path string, req any, rsp any) error {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		panic(err) // should never happen
	}
	httpReq, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, "https://api.dropboxapi.com/2"+path, reqBytes)
	if err != nil {
		panic(err) // should never happen
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpRsp, err := db.cl.Do(httpReq)
	if err != nil {
		return fmt.Errorf("dropbox http request error: %w", err)
	}
	body, err := io.ReadAll(httpRsp.Body)
	if err != nil {
		return fmt.Errorf("dropbox http body read error: %w", err)
	}
	if httpRsp.StatusCode != http.StatusOK {
		return fmt.Errorf("dropbox http bad status code: %d\n%s", httpRsp.StatusCode, string(body))
	}

	if rsp != nil {
		if err := json.Unmarshal(body, rsp); err != nil {
			return fmt.Errorf("dropbox could not unmarshal: %w\n%s", err, string(body))
		}
	}

	return nil
}

type DownloadRequest struct {
	Path string `json:"path"`
}

func (db *DropBoxClient) Download(ctx context.Context, path string) ([]byte, error) {
	req := &DownloadRequest{
		Path: path,
	}
	buf, err := db.content(ctx, "/files/download", req, nil)
	if err != nil {
		return nil, fmt.Errorf("dropbox list files: %w", err)
	}
	return buf, nil
}

func (db *DropBoxClient) content(ctx context.Context, path string, req any, rsp any) ([]byte, error) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		panic(err) // should never happen
	}
	httpReq, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPost, "https://content.dropboxapi.com/2"+path, nil)
	if err != nil {
		panic(err) // should never happen
	}
	httpReq.Header.Set("Dropbox-API-Arg", string(reqBytes))

	httpRsp, err := db.cl.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("dropbox http request error: %w", err)
	}
	body, err := io.ReadAll(httpRsp.Body)
	if err != nil {
		return nil, fmt.Errorf("dropbox http body read error: %w", err)
	}

	if httpRsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dropbox http bad status code: %d\n%s", httpRsp.StatusCode, string(body))
	}

	if rsp != nil {
		rspJson := httpRsp.Header.Get("Dropbox-API-Result")
		if err := json.Unmarshal([]byte(rspJson), rsp); err != nil {
			return nil, fmt.Errorf("dropbox could not unmarshal: %w\n%s", err, rspJson)
		}
	}

	return body, nil
}
