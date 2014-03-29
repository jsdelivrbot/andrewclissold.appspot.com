package andrewclissold

import "github.com/russross/blackfriday"
import "html/template"
import "io/ioutil"
import "net/http"
import "os"
import "strings"

func init() {
	http.HandleFunc("/", rootHandler)

	http.HandleFunc("/code", pageHandler)
	http.HandleFunc("/theory", pageHandler)
	http.HandleFunc("/music", pageHandler)
	http.HandleFunc("/snips", postHandler)

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

func pageHandler(w http.ResponseWriter, r *http.Request) {
	title := strings.ToUpper(string(r.URL.Path[1])) + r.URL.Path[2:]

	templates.ExecuteTemplate(w, "header.html", &info{title, ie})
	templates.ExecuteTemplate(w, r.URL.Path[1:]+".html", nil)
	templates.ExecuteTemplate(w, "footer.html", title)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path[1:]
	title := strings.ToUpper(string(r.URL.Path[1])) + r.URL.Path[2:]

	posts := make([][]byte, 0)

    // Find all posts within the directory
	dir := "posts/" + path + "/"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set up the markdown renderer
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_GITHUB_BLOCKCODE
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Open and read the file
		file, err := os.Open(dir + file.Name())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fi, err := file.Stat()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := make([]byte, fi.Size())
		file.Read(data)

		// Render the file's markdown
		post := blackfriday.Markdown(data, renderer, blackfriday.EXTENSION_FENCED_CODE)

		posts = append(posts, post)
	}

	templates.ExecuteTemplate(w, "header.html", &info{title, ie})
	templates.ExecuteTemplate(w, path+".tmpl", nil)
	for i, post := range posts {
		w.Write(post)
		if i < len(posts)-1 {
			w.Write([]byte("<hr>"))
		}
	}

	templates.ExecuteTemplate(w, "footer.html", title)
}

type info struct {
	Title string
	IE    template.HTML
}

var ie template.HTML = `
    <!--[if lt IE 9]>
      <script src="js/html5shiv.min.js"></script>
      <script src="js/respond.min.js"></script>
    <![endif]-->`

var templates = template.Must(template.ParseFiles(
	"header.html",

	"code.html", "theory.html", "music.html", "tmpl/snips.tmpl",

	"footer.html"))
