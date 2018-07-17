package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/taskcluster/statsum/payload"
)

func validateProjectName(project string) bool {
	for i := 0; i < len(project); i++ {
		b := project[i]
		if !('0' <= b && b <= '9') &&
			!('a' <= b && b <= 'z') &&
			!('A' <= b && b <= 'Z') &&
			b != '-' && b != '_' {
			return false
		}
	}
	return true
}

type contentType int

// List of content types
const (
	NoFormat contentType = iota
	JSONFormat
	MsgPackFormat
)

func detectContentType(contentType string) contentType {
	if contentType == "application/msgpack" || contentType == "application/x-msgpack" {
		return MsgPackFormat
	}
	if strings.HasPrefix(contentType, "application/json") &&
		(len(contentType) == 16 || contentType[16] == ';') {
		return JSONFormat
	}
	return NoFormat
}

func detectResponseType(r *http.Request, fallback contentType) contentType {
	// Default to fallback if we have nothing else
	responseType := fallback

	// Parse "Accept" header, prefer MsgPackFormat if accepted
	accept := r.Header.Get("Accept")
	for accept != "" {
		i := strings.IndexByte(accept, ',')
		if i == -1 {
			i = len(accept)
		}
		part := strings.TrimSpace(accept[:i])
		if strings.HasPrefix(part, "application/json") &&
			(len(part) == 16 || part[16] == ';') {
			responseType = JSONFormat
		}
		if strings.HasPrefix(part, "application/msgpack") &&
			(len(part) == 19 || part[19] == ';') {
			return MsgPackFormat
		}
		if strings.HasPrefix(part, "application/x-msgpack") &&
			(len(part) == 21 || part[21] == ';') {
			return MsgPackFormat
		}
		accept = accept[i:]
	}

	// If we don't have any idea, then we parse the request Content-Type
	if responseType == NoFormat {
		responseType = detectContentType(r.Header.Get("Content-Type"))
		// If all fails we fallback to JSON
		if responseType == NoFormat {
			responseType = JSONFormat
		}
	}

	// Return detected response type
	return responseType
}

func reply(w http.ResponseWriter, statusCode int, r payload.Response, responseType contentType) {
	// Write Content-Type header (fallback to JSON if NoFormat)
	b := []byte{}
	err := error(nil)
	switch responseType {
	case NoFormat, JSONFormat:
		w.Header().Set("Content-Type", "application/json")
		b, err = json.Marshal(r)
	case MsgPackFormat:
		w.Header().Set("Content-Type", "application/msgpack")
		b, err = r.MarshalMsg(b)
	}
	if err != nil {
		panic(err)
	}

	// Add Content-Length, send header and body
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Header().Set("Connection", "close")
	w.WriteHeader(statusCode)
	w.Write(b)
}
