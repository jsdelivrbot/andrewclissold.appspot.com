package andrewclissold

import (
	"bufio"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

func init() {
	http.HandleFunc("/", rootHandler)

	http.HandleFunc("/apps", tabHandler)
	http.HandleFunc("/music", tabHandler)
	http.HandleFunc("/snips", tabHandler)

	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("files"))))

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("fonts"))))
	http.Handle("/ico/", http.StripPrefix("/ico/", http.FileServer(http.Dir("ico"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("img"))))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	title := "Home"

	templates.ExecuteTemplate(w, "header.tmpl", &info{title, ie})
	templates.ExecuteTemplate(w, "index.html", nil)
	templates.ExecuteTemplate(w, "footer.tmpl", title)
}

func tabHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:]
	title := strings.ToUpper(string(r.URL.Path[1])) + r.URL.Path[2:]

	var posts Posts

	// Find all posts within the directory
	dir := "tabs/" + path + "/"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set up the markdown renderer
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

	for _, fi := range files {
		// Skip any files not ending in ".md"
		re := regexp.MustCompile(`\.md$`)
		if !re.Match([]byte(fi.Name())) {
			continue
		}

		// Open the file
		file, err := os.Open(dir + fi.Name())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		scanner := bufio.NewScanner(file)

		// Read the first line of the file as the post's creation date
		scanner.Scan()
		ref := "2 Jan 2006"
		date, err := time.Parse(ref, scanner.Text())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Don't display any posts intended for the future
		if date.After(time.Now()) {
			continue
		}

		// Read the rest of the file as the post itself
		var data []byte
		for scanner.Scan() {
			data = append(data, scanner.Bytes()...)
			data = append(data, '\n') // add back the stripped newlines
		}
		if err := scanner.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Render the file's markdown
		content := blackfriday.Markdown(data, renderer, blackfriday.EXTENSION_FENCED_CODE)

		// Add the parsed post to the collection of posts
		posts = append(posts, post{Content: content, Date: date})
	}

	// Sort the posts by creation date
	sort.Sort(posts)

	templates.ExecuteTemplate(w, "header.tmpl", &info{title, ie})
	w.Write([]byte(`<div class="` + strings.ToLower(title) + `">`))
	templates.ExecuteTemplate(w, path+".tmpl", nil)
	for i, post := range posts {
		if i == len(posts)-1 {
			w.Write([]byte(`<div class="post last">`))
		} else {
			w.Write([]byte(`<div class="post">`))
		}
		w.Write(post.Content)
		w.Write([]byte(`</div>`))
	}
	w.Write([]byte(`</div>`))
	templates.ExecuteTemplate(w, "footer.tmpl", title)
}

type post struct {
	Content []byte
	Date    time.Time
}

type Posts []post

// Satisfy sort.Interface to sort posts by creation date
func (p Posts) Len() int {
	return len(p)
}
func (p Posts) Less(i, j int) bool {
	// Return Date.After instead of Before to place them newest-first
	return p[i].Date.After(p[j].Date)
}
func (p Posts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
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
	"tmpl/header.tmpl",
	"index.html",
	"tmpl/apps.tmpl", "tmpl/music.tmpl", "tmpl/snips.tmpl",
	"tmpl/footer.tmpl"))
