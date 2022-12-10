package api

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"gzipResponse"
	"net/http"
	"something"
	"strings"
)

func ServeAPIUseGZip(response http.ResponseWriter, request *http.Request) {
	if !strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
		ServeAPI(response, request)
	} else {
		response.Header().Set("Content-Encoding", "gzip")
		gzipping := gzip.NewWriter(response)
		defer gzipping.Close()
		ServeAPI(gzipResponse.Writer{Writer: gzipping, ResponseWriter: response}, request)
	}
}

// ServeAPI - serves the api request by handling the response in a certain way
func ServeAPI(response http.ResponseWriter, request *http.Request) {
	// Redirect if not trailing slash
	// if the url contains / in the end, then we will ignore it for the API listener (return same result no matter if there is / or not
	urlPath := request.URL.Path
	if hasTrailingSlash(urlPath) == false {
		addTrailingSlashRedirect(urlPath, response, request)
	}

	// set response header as json
	response.Header().Set("Content-Type", "application/json")
	response.Header().Add("Cache-Control", "max-age=31536000, immutable")

	// response
	var jsonResponse any = JSONMessageResponse{false, "Not set"}

	// api commands (URLs)
	switch request.URL.Path {

	case "/api/test/list/":
		list, err := something.List(true)
		if err != nil {
			// error getting data, sending error
			jsonResponse = JSONMessageResponse{false, err.Error()}
		} else {
			// successfully sending data
			jsonResponse = JSONDataResponse{true, list}
		}
		break
	default:
		// if not part of api
		jsonResponse = JSONMessageResponse{false, "Part of the api does not exist"}
		break
	}

	// Marshal JSON
	jsonFormattedResponse := jsonMarshalled(jsonResponse)
	// Send Response
	sendResponse(response, jsonFormattedResponse)
}

func jsonMarshalled(jsonResponse any) string {
	jsonFormattedResponse, err := json.Marshal(jsonResponse)
	//jsonResponseString, err := json.MarshalIndent(jsonResponse, "", "\t")
	if err != nil {
		fmt.Println("Could not json.Marshal - error: " + err.Error())
	}
	return string(jsonFormattedResponse)
}

func sendResponse(response http.ResponseWriter, jsonFormattedResponse string) {
	printText, err := fmt.Fprintln(response, jsonFormattedResponse)
	if err != nil {
		// if there was and error when writing response
		fmt.Println("Could not fmt.Fprintln - error: " + err.Error())
		fmt.Println("formattedPrintNewLine: " + string(rune(printText)))
	}
}

func hasTrailingSlash(urlPath string) bool {
	if urlPath[len(urlPath)-1] == '/' {
		return true
	} else {
		return false
	}
}

func addTrailingSlashRedirect(urlPath string, response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, urlPath+"/", 301) // redirect url to without trailing slash, permanent redirect
}
