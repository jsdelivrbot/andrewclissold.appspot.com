package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/code", codeHandler)
	http.HandleFunc("/theory", theoryHandler)
	http.HandleFunc("/music", musicHandler)
	http.HandleFunc("/snips", snipsHandler)

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("/home/aclissold/Code/public_html/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("/home/aclissold/Code/public_html/js"))))
	http.Handle("/ico/", http.StripPrefix("/ico/", http.FileServer(http.Dir("/home/aclissold/Code/public_html/ico"))))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "/home/aclissold/Code/public_html/index.html")
}

func codeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/home/aclissold/Code/public_html/code.html")
}

func theoryHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/home/aclissold/Code/public_html/theory.html")
}

func musicHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/home/aclissold/Code/public_html/music.html")
}

func snipsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/home/aclissold/Code/public_html/snips.html")
}
