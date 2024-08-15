package conf

import (
	_ "embed"
	"fmt"
	"github.com/zwh8800/RGTV/util"
	"html/template"
	"net/http"
	"os"
)

const (
	addr = ":8080"
)

var (
	server *http.Server
	mux    *http.ServeMux
)

//go:embed home.html
var homePageHtml []byte

func StartNetConfServer() {
	mux = http.NewServeMux()
	mux.HandleFunc("/", homePage)
	mux.HandleFunc("/save", saveConf)
	server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

}

func StopNetConfServer() {
	if server == nil {
		return
	}
	server.Close()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	err := template.Must(template.New("home").Parse(string(homePageHtml))).Execute(w, config)
	if err != nil {
		fmt.Println(err)
	}
}

func saveConf(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	config.LiveSourceUrl = r.Form.Get("live_source_url")
	config.EPGUrl = r.Form.Get("epg_url")
	data := util.ToPrettyJson(config)
	os.WriteFile("tv.json", []byte(data), 0644)
	http.Redirect(w, r, "/", http.StatusFound)
}
