package main

import (
	"html/template"
	"log"
	"net"
	"net/http"
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
		IP:       GetOutboundIP().String(),
		Port:     r.RemoteAddr,
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

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", requestHandler)

	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}

/*
	import (
		"log"
		"net"
		"net/http"
		"strings"
	)

	func sayHello(w http.ResponseWriter, r *http.Request) {
		message := r.URL.Path
		message = strings.TrimPrefix(message, "/")
		message = "Hello " + message
		w.Write([]byte(message))
	}
	func GetOutboundIP() net.IP {
		conn, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		localAddr := conn.LocalAddr().(*net.UDPAddr)

		return localAddr.IP
	}
	func main() {
		http.HandleFunc("/", sayHello)
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic(err)
		}
		//	https://besthqwallpapers.com/Uploads/4-11-2017/27088/max-verstappen-dutch-racing-driver-red-bull-racing-formula-1-number-33.jpg
	}
*/
