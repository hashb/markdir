package main

import (
	"errors"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/russross/blackfriday/v2"
)

var bind = flag.String("bind", "127.0.0.1:19000", "port to run the server on")

func main() {
	flag.Parse()

	httpdir := http.Dir(".")
	handler := renderer{httpdir, http.FileServer(httpdir)}

	log.Println("Serving on http://" + *bind)
	log.Fatal(http.ListenAndServe(*bind, handler))
}

var outputTemplate = template.Must(template.New("base").Parse(`
<html>
  <head>
    <title>{{ .Path }}</title>
    <style>:root{--c-light-text: #333;--c-light-background: #fdfdfd;--c-light-secondary-background: #f8f8f8;--c-light-focus: #f0f8ff;--c-light-link-text: #06c;--c-dark-text: #ddd;--c-dark-background: #212121;--c-dark-secondary-background: #262626;--c-dark-focus: #f0f8ff;--c-dark-link-text: #26e}:root{--c-text: var(--c-light-text);--c-link-text: var(--c-light-link-text);--c-background: var(--c-light-background);--c-secondary-background: var(--c-light-secondary-background);--c-focus: var(--c-light-focus)}@media (prefers-color-scheme: dark){:root{--c-text: var(--c-dark-text);--c-link-text: var(--c-dark-link-text);--c-background: var(--c-dark-background);--c-secondary-background: var(--c-dark-secondary-background);--c-focus: var(--c-dark-focus)}}*,*:before,*:after{box-sizing:border-box}html,body{padding:0;margin:0;font-family:"Avenir", "Avenir Next", -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;color:var(--c-text)}body{display:grid;height:100vh;grid-template-rows:auto 1fr auto;background-color:var(--c-background)}main{max-width:37.5em}p,pre,code{line-height:1.5}a[href],a[href]:visited{color:var(--c-link-text)}a[href]:not(:hover){text-decoration:none}img{max-width:100%;height:auto}header,main,footer{padding:1em}main{padding-bottom:2em}main :first-child,main>article :first-child{margin-top:0}pre{font-size:14px;direction:ltr;text-align:left;white-space:pre;word-spacing:normal;word-break:normal;-moz-tab-size:2;tab-size:2;-webkit-hyphens:none;-moz-hyphens:none;-ms-hyphens:none;hyphens:none;padding:1em;margin:0.5em 0;overflow-x:auto;white-space:pre-wrap;white-space:-moz-pre-wrap;white-space:-pre-wrap;white-space:-o-pre-wrap;word-wrap:break-word}header>em{display:block;font-size:2em;margin:0.67em 0;font-weight:bold;font-style:normal}header nav ul{padding:0;list-style:none}header nav ul :first-child{margin-left:0}header nav li{display:inline-block;margin:0 0.25em}header nav li a{padding:0.25em 0.5em;border-radius:0.25em}header nav li a[href]:not(:hover){text-decoration:none}header nav li a[data-current="current item"]{background-color:var(--c-focus)}article{margin-bottom:1em;padding-bottom:1em}main>section>article>*{margin-top:0;margin-bottom:0.5em}div.post-footer{font-size:0.8rem}div.post-footer .meta{margin-top:0.25em}a[rel="tag"],a[rel="tag"]:visited{display:inline-block;vertical-align:text-top;text-transform:uppercase;letter-spacing:0.1em;font-size:0.8em;padding:0 0.5em;line-height:2em;height:2em;border:1px solid var(--c-background);color:var(--c-link-text);border-radius:0.25em;text-decoration:none;margin:0 0.5em 0.5em 0}a[rel="tag"]:hover{border:1px solid var(--c-link-text);background-color:var(--c-link-text);color:var(--c-light-background)}a[rel="tag"]:last-child{margin-right:0}form{display:grid;padding:2em 0}form label{display:none}input,textarea,button{width:100%;padding:1em;margin-bottom:1em;font-size:1rem;font-family:"Avenir", "Avenir Next", sans-serif}input,textarea{border:1px solid black}button{border:1px solid var(--c-link-text);background-color:var(--c-link-text);color:var(--c-background);cursor:pointer}@media screen and (min-width: 768px){:root{font-size:1.1rem}}div.archive-item h2{margin-top:5px;margin-bottom:5px}div.archive-item time{font:0.85em Monaco, Monospace}div.archive-item p{margin:0.3em}div.archive-month{margin-bottom:1.2em}figcaption{display:block;text-align:center;font-style:italic}.postnavigation{padding-top:10px;text-align:center;font-size:1.2em}.postnavigation .left{float:left}.postnavigation .right{float:right}ul.toc{padding:15px 15px 15px 25px}ul.toc ul{padding:0 0 3px 25px}.toc{display:inline-block;background-color:var(--c-background)}.highlight{margin-bottom:15px;border-radius:10px;background:var(--c-secondary-background)}.highlight .c{color:#998;font-style:italic}.highlight .err{color:#a61717;background-color:#e3d2d2}.highlight .k{font-weight:bold}.highlight .o{font-weight:bold}.highlight .cm{color:#998;font-style:italic}.highlight .cp{color:#999;font-weight:bold}.highlight .c1{color:#998;font-style:italic}.highlight .cs{color:#999;font-weight:bold;font-style:italic}.highlight .gd{color:#000;background-color:#fdd}.highlight .gd .x{color:#000;background-color:#faa}.highlight .ge{font-style:italic}.highlight .gr{color:#a00}.highlight .gh{color:#999}.highlight .gi{color:#000;background-color:#dfd}.highlight .gi .x{color:#000;background-color:#afa}.highlight .go{color:#888}.highlight .gp{color:#555}.highlight .gs{font-weight:bold}.highlight .gu{color:#aaa}.highlight .gt{color:#a00}.highlight .kc{font-weight:bold}.highlight .kd{font-weight:bold}.highlight .kp{font-weight:bold}.highlight .kr{font-weight:bold}.highlight .kt{color:#458;font-weight:bold}.highlight .m{color:#099}.highlight .s{color:#d14}.highlight .na{color:#008080}.highlight .nb{color:#0086b3}.highlight .nc{color:#458;font-weight:bold}.highlight .no{color:#008080}.highlight .ni{color:#800080}.highlight .ne{color:#900;font-weight:bold}.highlight .nf{color:#900;font-weight:bold}.highlight .nn{color:#555}.highlight .nt{color:#af00af}.highlight .nv{color:#008080}.highlight .ow{font-weight:bold}.highlight .w{color:#bbb}.highlight .mf{color:#099}.highlight .mh{color:#099}.highlight .mi{color:#099}.highlight .mo{color:#099}.highlight .sb{color:#d14}.highlight .sc{color:#d14}.highlight .sd{color:#d14}.highlight .s2{color:#d14}.highlight .se{color:#d14}.highlight .sh{color:#d14}.highlight .si{color:#d14}.highlight .sx{color:#d14}.highlight .sr{color:#009926}.highlight .s1{color:#d14}.highlight .ss{color:#990073}.highlight .bp{color:#999}.highlight .vc{color:#008080}.highlight .vg{color:#008080}.highlight .vi{color:#008080}.highlight .il{color:#099}</style>
  </head>
  <body>
    <main>
      {{ .Body }}
    </main>
  </body>
</html>
`))

type renderer struct {
	d http.Dir
	h http.Handler
}

func (r renderer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if !strings.HasSuffix(req.URL.Path, ".md") && !strings.HasSuffix(req.URL.Path, "/guide") {
		r.h.ServeHTTP(rw, req)
		return
	}

	// net/http is already running a path.Clean on the req.URL.Path,
	// so this is not a directory traversal, at least by my testing
	var pathErr *os.PathError
	input, err := ioutil.ReadFile("." + req.URL.Path)
	if errors.As(err, &pathErr) {
		http.Error(rw, http.StatusText(http.StatusNotFound)+": "+req.URL.Path, http.StatusNotFound)
		log.Printf("file not found: %s", err)
		return
	}

	if err != nil {
		http.Error(rw, "Internal Server Error: "+err.Error(), 500)
		log.Printf("Couldn't read path %s: %v (%T)", req.URL.Path, err, err)
		return
	}

	output := blackfriday.Run(input)

	rw.Header().Set("Content-Type", "text/html; charset=utf-8")

	outputTemplate.Execute(rw, struct {
		Path string
		Body template.HTML
	}{
		Path: req.URL.Path,
		Body: template.HTML(string(output)),
	})

}
