# go-unhar
Zero dependency golang module and CLI to handle [HTTP Archive (HAR)](https://w3c.github.io/web-performance/specs/HAR/Overview.html) files.

**go-unhar can be used as:**
1. A lightweight go module to parse HAR files.
2. A standalone CLI tool to extract (think unzip) all resources (images, webpages, api answers) from a given HAR file.

## go-unhar as a module
### Install
```bash
go get github.com/code-lion-com/go-unhar
```
### Load from file
Load HAR file to [`goUnhar.Har`](https://github.com/code-lion-com/go-unhar/blob/main/types.go) struct.
```go
// Load HAR from file
har := &goUnhar.Har{}
har.Open("myFile.har")
```

### Examples
#### Display HAR tool name and version
```go
fmt.Printf(
    "HTTP Archive (HAR) version: %s\nCreator: %s (%s)\n",
    har.Log.Version, har.Log.Creator.Name, har.Log.Creator.Version,
)
//  HTTP Archive (HAR) version: 
//  Creator: WebInspector (537.36)
```

#### Overview all Entries

```go
for _, entry := range har.Log.Entries {
    status := entry.Response.Status
    method := entry.Request.Method
    url := entry.Request.URL
    time := float64(entry.Time)
    time = math.Round(time)
    mimetype := entry.Response.Content.MimeType
    ctype := mimetype[strings.LastIndex(mimetype, "/")+1:]
    fmt.Printf("%d %s %dms %s %s\n", status, method, int(time), ctype, url)
}
// [...]
// 200 GET 174ms javascript https://test.com/script.js
// 200 GET 86ms svg+xml https://test.com/logo.svg
// [...]
```

#### Display Cookies Sent to Server by Entry
```go
for _, entry := range har.Log.Entries {
    if len(entry.Request.Cookies) > 0 {
        fmt.Printf("\n🍪 %d found %s\n", len(entry.Request.Cookies), entry.Request.URL)
    }
    cookies := entry.Request.Cookies
    for nCookie, cookie := range cookies {
        fmt.Printf("  🍪 #%d %v\n", nCookie+1, goUnhar.NVP{Name: cookie.Name, Value: cookie.Value})
    }
}
// 🍪 3 found https://test.com/api/me
//   🍪 #1 {SESSIONID 123456789}
//   🍪 #2 {lang v=2&lang=en-us}
//   🍪 #3 {bcookie "v=1&00000000-0000-0000-0000-00000000"}
//
// 🍪 17 found https://ads.com/tracker
//   🍪 #1 {uid 123456789}
//   🍪 #2 {lang en-us}
//   [...]
```

## UnHAR as a standalone CLI

```bash
# Add the cli tool to your gopath
go get github.com/code-lion-com/go-unhar/cmd/unhar

cd ~/Desktop  # working dir
unhar www.example.com.har  # extract to working dir
ls ~/Desktop/www.example.com  # list working dir
# ~/Desktop/www.example.com
# - index.html
# - about.html
# -- api
#    |-- users/me
# -- static
#    |-- me.png
#    |-- style.css
```

## Credits and similar projects
- https://github.com/mrichman/hargo
  Hargo parses HAR files, can convert to curl format, and serve as a load test driver.
- https://github.com/arp242/har
  Read HAR ("HTTP Archive format") archives
