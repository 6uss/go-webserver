package main

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
)

type NewsAggPage struct {
	IP       string
	Port     string
	Method   string
	Protocol string
	Host     string
	Agent    string
	RespTime float32
}

var indextmpl = template.Must(template.ParseFiles("index.html"))

func requestHandler(w http.ResponseWriter, r *http.Request) {

	p := NewsAggPage{
		IP:       getPublicIP(),
		Port:     getRemotePort(r.RemoteAddr),
		Method:   r.Method,
		Protocol: r.Proto,
		Host:     r.Host,
		Agent:    r.UserAgent(),
		RespTime: 2,
	}
	//t, _ := template.ParseFiles("index.html")
	if err := indextmpl.ExecuteTemplate(w, "index.html", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getPublicIP() string {
	url := "http://169.254.169.254/latest/meta-data/public-ipv4"

	req, _ := http.NewRequest("GET", url, nil)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return string(body)
}
func getRemotePort(remAddr string) string {
	re := regexp.MustCompile("(?:[0-9]+)$")
	match := re.FindStringSubmatch(remAddr)
	return match[0]

}
func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", requestHandler)

	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}
