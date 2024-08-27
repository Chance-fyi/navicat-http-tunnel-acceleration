package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"mime/multipart"
)

type formData struct {
	Key    string
	Action string
	Query  []string
}

func parseFormData(data []byte, boundary string) *formData {
	buffer := bytes.NewBuffer(data)
	reader := multipart.NewReader(buffer, boundary)

	var (
		host         string
		port         string
		login        string
		pass         string
		db           string
		encodeBase64 string
		FormData     formData
	)

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		n := part.FormName()
		v, err := io.ReadAll(part)
		if err != nil {
			log.Fatal(err)
		}
		switch n {
		case "host":
			host = string(v)
		case "port":
			port = string(v)
		case "login":
			login = string(v)
		case "password":
			pass = string(v)
		case "db":
			db = string(v)
		case "encodeBase64":
			encodeBase64 = string(v)
		case "actn":
			FormData.Action = string(v)
		case "q[]":
			FormData.Query = append(FormData.Query, string(v))
		}
	}
	FormData.Key = fmt.Sprintf("%s:%s:%s:%s:%s", host, port, login, pass, db)

	if encodeBase64 != "" {
		for i, s := range FormData.Query {
			// base64 decode
			decodeString, _ := base64.StdEncoding.DecodeString(s)
			FormData.Query[i] = string(decodeString)
		}
	}
	return &FormData
}
