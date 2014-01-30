package andrewclissold

import "net/http"

func init() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/code", codeHandler)
	http.HandleFunc("/theory", theoryHandler)
	http.HandleFunc("/music", musicHandler)
	http.HandleFunc("/snips", snipsHandler)

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/ico/", http.StripPrefix("/ico/", http.FileServer(http.Dir("ico"))))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func codeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "code.html")
}

func theoryHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "theory.html")
}

func musicHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "music.html")
}

func snipsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "snips.html")
}
