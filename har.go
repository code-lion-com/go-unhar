package goUnhar

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Parser for .har files.
func (har *Har) Parse(data []byte) (*Har, error) {
	err := json.Unmarshal(data, &har)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	return har, nil
}

// Open and parse .har file.
func (har *Har) Open(filePath string) error {
	fp, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	data, err := io.ReadAll(fp)
	fp.Close()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	har, err = har.Parse(data)
	if err != nil {
		return err
	}
	return err
}

// Extract all the files.
func (h *Har) Write(filePath string, verbose bool) error {
	root := filepath.Base(filePath)
	if ext := filepath.Ext(root); ext != "" {
		root = root[:len(root)-len(ext)]
	}

	for _, e := range h.Log.Entries {
		path := root + "/" + strings.TrimPrefix(strings.TrimPrefix(e.Request.URL, "http://"), "https://")
		if verbose {
			fmt.Println("  ", path)
		}

		for _, h := range e.Response.Headers {
			if h.Name == "Content-Disposition" {
				for _, v := range strings.Split(h.Value, ";") {
					v := strings.TrimSpace(v)
					if strings.HasPrefix(v, "filename=") {
						v = filepath.Clean(strings.Trim(v, `"`))
						v = strings.ReplaceAll(v, "/", "")
						path = filepath.Dir(path) + "/" + v
					}
				}
				break
			}
		}
		if strings.HasSuffix(path, "/") {
			path += "index.html"
		}

		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return fmt.Errorf("mkdir: %w", err)
		}

		var resp []byte
		if e.Response.Content.Encoding == "base64" {
			resp, err = base64.StdEncoding.DecodeString(e.Response.Content.Text)
			if err != nil {
				return fmt.Errorf("base64: %w", err)
			}
		} else {
			resp = []byte(e.Response.Content.Text)
		}

		err = os.WriteFile(path, resp, 0644)
		if err != nil {
			return fmt.Errorf("write: %w", err)
		}
	}

	fmt.Printf("Extracted %d files\n", len(h.Log.Entries))
	return nil
}
