# go-dry

**DRY (Don't Repeat Yourself) utility package for Go**

[![Go Reference](https://pkg.go.dev/badge/github.com/ungerik/go-dry.svg)](https://pkg.go.dev/github.com/ungerik/go-dry)
[![Go Report Card](https://goreportcard.com/badge/github.com/ungerik/go-dry)](https://goreportcard.com/report/github.com/ungerik/go-dry)

**Note:** This package replaces `github.com/ungerik/go-quick`

## Overview

`go-dry` is a collection of utility functions and types that eliminate common boilerplate code in Go applications. It provides helpers for:

- **String operations** - parsing, formatting, searching, transformation
- **File I/O** - reading/writing various formats (JSON, XML, CSV, config files)
- **HTTP utilities** - compression, JSON/XML marshaling, request handling
- **Error handling** - error collection, wrapping, panic helpers
- **Byte operations** - encoding, compression, hashing
- **Reflection helpers** - struct field manipulation, sorting
- **Concurrency** - thread-safe wrappers for common types
- **Encryption** - AES encryption/decryption with pooling
- **Random generation** - cryptographically secure random strings
- **Debug utilities** - stack traces, pretty printing, debug mutexes

## Installation

```bash
go get github.com/ungerik/go-dry
```

## Quick Examples

### String Operations

```go
// Parse with zero values on error
age := dry.StringToInt("25")
price := dry.StringToFloat("19.99")
enabled := dry.StringToBool("true")

// Find text between tokens
content, remainder, found := dry.StringFindBetween(html, "<title>", "</title>")

// Convert case styles
camel := dry.StringToUpperCamelCase("user_name") // "UserName"
```

### File Operations

```go
// Read/Write JSON
var config Config
dry.FileUnmarshallJSON("config.json", &config)
dry.FileSetJSONIndent("output.json", data, "  ")

// Read lines
lines, _ := dry.FileGetNonEmptyLines("data.txt")

// Copy files and directories
dry.FileCopy("source.txt", "dest.txt")
dry.FileCopyDir("source_dir", "dest_dir")
```

### HTTP Helpers

```go
// Respond with compressed JSON
dry.HTTPRespondMarshalJSON(data, w, r)

// Post JSON to endpoint
err := dry.HTTPPostJSON("https://api.example.com/data", payload)

// Unmarshal request body
var input RequestData
dry.HTTPUnmarshalRequestBodyJSON(r, &input)
```

### Error Handling

```go
// Collect multiple errors
errs := dry.NewErrorList()
errs.Collect(operation1())
errs.Collect(operation2())
if err := errs.Err(); err != nil {
    log.Fatal(err)
}

// Panic on error with stack trace
dry.PanicIfErr(file.Close())
```

### Compression Pools

```go
// Efficient gzip compression with pooling
var buf bytes.Buffer
writer := dry.Gzip.GetWriter(&buf)
writer.Write(data)
dry.Gzip.ReturnWriter(writer)
compressed := buf.Bytes()
```

### Thread-Safe Types

```go
counter := dry.NewSyncInt(0)
counter.Add(1)
current := counter.Get()

config := dry.NewSyncMap()
config.AddString("key", "value")
```

## Documentation

Full API documentation is available at:
- [pkg.go.dev/github.com/ungerik/go-dry](https://pkg.go.dev/github.com/ungerik/go-dry)
- [godoc.org/github.com/ungerik/go-dry](https://godoc.org/github.com/ungerik/go-dry)

## Key Features

### String Utilities
- Conversion: `StringToInt`, `StringToFloat`, `StringToBool`
- Formatting: `StringMarshalJSON`, `StringPrettifyJSON`, `StringCSV`
- Searching: `StringFind`, `StringFindBetween`, `StringInSlice`
- Transformation: `StringToUpperCamelCase`, `StringToLowerCamelCase`
- HTML/XML: `StringStripHTMLTags`, `StringReplaceHTMLTags`
- Hashing: `StringMD5Hex`, `StringSHA1Base64` (non-cryptographic use only)

### File Operations
- Universal reader supporting files and URLs
- JSON/XML/CSV marshaling and unmarshaling
- Line-by-line reading with `FileGetLines`, `FileGetNonEmptyLines`
- Config file parsing (key=value format)
- Compression: deflate and gzip
- Checksums: MD5, CRC64
- File utilities: `FileExists`, `FileIsDir`, `FileTouch`, `FileTimeModified`

### HTTP Utilities
- Automatic gzip/deflate compression with `HTTPCompressHandler`
- JSON/XML response helpers with compression
- Form POST/PUT with status code returns
- Request body unmarshaling

### Error Handling
- `ErrorList` for collecting multiple errors
- `PanicIfErr` with stack traces
- `FirstError`, `LastError` for error sequences
- `AsError` for interface{} to error conversion

### Byte Operations
- Base64/Hex encoding and decoding
- Compression: `BytesDeflate`, `BytesGzip`
- MD5 hashing (checksums only)
- Head/Tail operations like Unix commands
- Map and Filter functions

### Reflection Helpers
- Struct field manipulation from string maps
- Generic sorting with reflection
- Exported field enumeration

### Concurrency
- `SyncBool`, `SyncInt`, `SyncFloat`, `SyncString` - thread-safe primitives
- `SyncMap`, `SyncStringMap` - thread-safe maps
- `SyncPoolMap` - thread-safe pool management
- `DebugMutex`, `DebugRWMutex` - mutexes with logging

### Encryption
- AES encryption/decryption with cipher block pooling
- Support for AES-128, AES-192, AES-256

### I/O Utilities
- `CountingReader`, `CountingWriter`, `CountingReadWriter`
- `ReaderFunc`, `WriterFunc` - function types implementing interfaces
- `ReadLine`, `WriteFull` helpers

### Debug & Development
- `StackTrace`, `StackTraceLine` - runtime stack inspection
- `PrettyPrintAsJSON` - formatted JSON output
- `Nop` - dummy function to avoid unused import errors

## Security Notes

⚠️ **Cryptographic Hash Functions**: This package uses MD5 and SHA1 for **non-cryptographic purposes only** (checksums, cache keys). Do not use these functions for security-sensitive operations like password hashing or digital signatures.

⚠️ **Input Validation**: Functions like `StringToInt`, `StringToFloat`, etc. return zero values on parse errors. For production code where error detection is critical, use the standard library's parsing functions instead.

⚠️ **File Operations**: Functions accepting `filenameOrURL` parameters can read from arbitrary URLs and local files. Validate and sanitize inputs in security-sensitive contexts.

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestName
```

## License

See LICENSE file for details.
