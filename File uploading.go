package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var uploadFormTemplate = []byte(`
<html>
	<body>
	<form action="/upload" method = "post"
	enctype="multipart/form-data">
		Image: <input type="file" name="my_file">
		<input type="submit" value="Upload">
	</form>
	</body>
</html>
`)

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write(uploadFormTemplate)
}

func uploadPage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(5 * 1024 * 1024)
	file, fileHeader, err := r.FormFile("my_file")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Fprintf(w, "Filename: %v\n", fileHeader.Filename)
	fmt.Fprintf(w, "Header: %#v\n", fileHeader.Header)

	hasher := md5.New()
	io.Copy(hasher, file)

	fmt.Fprintf(w, "md5 %x", hasher.Sum(nil))
}

type Params struct {
	Id   int
	User string
}

func uploadRawBody(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
		return
	}

	p := &Params{}
	json.Unmarshal(body, p)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	fmt.Fprintf(w, "content type = %#v\n",
		r.Header.Get("Content-Type"))
	fmt.Fprintf(w, "params: %#v\n", p)
}

func main() {
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/upload", uploadPage)
	http.HandleFunc("/raw_body", uploadRawBody)

	addr := ":8080"
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		panic(err)
	}
}
