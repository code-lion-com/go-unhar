package goUnhar

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

var testFile = "har_test.har"

func TestOpen(t *testing.T) {
	har := &Har{}
	err := har.Open(testFile)
	if err != nil {
		t.Fatal(err)
	}
	if har.Log.Creator.Name != "WebInspector" {
		t.Fatalf("Failed to parse expected data")
	}
	fmt.Printf(
		"HTTP Archive (HAR) version: %s\nCreator: %s (%s)\n",
		har.Log.Version, har.Log.Creator.Name, har.Log.Creator.Version,
	)
}

func TestOverview(t *testing.T) {
	har := &Har{}
	err := har.Open(testFile)
	if err != nil {
		t.Fatal(err)
	}
	if len(har.Log.Entries) == 0 {
		t.Fatal("No entries found")
	}
	for _, entry := range har.Log.Entries {
		status := entry.Response.Status
		if status == 0 {
			t.Fatal("Missing status")
		}
		method := entry.Request.Method
		if method == "" {
			t.Fatal("Missing method")
		}
		url := entry.Request.URL
		if url == "" {
			t.Fatal("Missing URL")
		}
		time := float64(entry.Time)
		time = math.Round(time)
		if time == 0 {
			t.Fatal("Missing time")
		}
		mimetype := entry.Response.Content.MimeType
		ctype := mimetype[strings.LastIndex(mimetype, "/")+1:]
		if ctype == "" {
			t.Fatal("Missing mimetype")
		}
		fmt.Printf("%d %s %dms %s %s\n", status, method, int(time), ctype, url)
	}
}

func TestFindSentCookies(t *testing.T) {
	har := &Har{}
	err := har.Open(testFile)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range har.Log.Entries {
		if len(entry.Request.Cookies) > 0 {
			fmt.Printf("\n🍪 %d found %s\n", len(entry.Request.Cookies), entry.Request.URL)
		} else {
			t.Fatal("Missing cookies")
		}
		cookies := entry.Request.Cookies
		for nCookie, cookie := range cookies {
			fmt.Printf("  🍪 #%d %v\n", nCookie+1, NVP{Name: cookie.Name, Value: cookie.Value})
		}
	}
}
