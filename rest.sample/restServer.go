package main

import (
	"encoding/json"
	//	"fmt"
	"io"
	"net/http"
	//	"os"
	//	"runtime"
	"time"

	l4g "github.com/alecthomas/log4go"
	"github.com/garyburd/redigo/redis"
)

var (
	RedisClient *redis.Pool
	REDIS_HOST  string
	REDIS_DB    int

//	db          *sql.DB
)

type StructResult struct {
	Cd  int
	Msg string
}

func init() {

	l4g.AddFilter("stdout", l4g.DEBUG, l4g.NewConsoleLogWriter())
	l4g.AddFilter("file", l4g.DEBUG, l4g.NewFileLogWriter("rest.log", false))
	//	db, err := sql.Open("postgres", "user=postgres dbname=meidb sslmode=disable")
}

func main() {

	startServer()

	defer l4g.Close()
	//	defer db.Close()
}

func login(w http.ResponseWriter, req *http.Request) {

	userName := req.FormValue("userName")
	pwd := req.FormValue("pwd")
	if userName == "" || pwd == "" {
		io.WriteString(w, getResult(500, "Login is failed,please use correct username or password."))
		l4g.Log(l4g.INFO, "", "Login is failed,please use correct username or password.")
	} else {
		io.WriteString(w, getResult(200, userName))
		l4g.Log(l4g.INFO, "", "Login is successful."+getResult(200, userName))

	}

}

func getResult(cd int, msg string) string {
	rslt := StructResult{
		Cd:  cd,
		Msg: msg,
	}
	ret, _ := json.Marshal(rslt)
	return string(ret)
}

func startServer() {

	addr := "127.0.0.1:12345"
	http.HandleFunc("/login", login)
	l4g.Log(l4g.INFO, "", "Server start.")

	go ListenAndServe(addr, nil, time.Second*10)

	//	http.HandleFunc("/login", login)

	//	l4g.Log(l4g.INFO, "", "Server start.")

	//	err := http.ListenAndServe(":12345", nil)
	//	if err != nil {
	//		l4g.Log(l4g.ERROR, "startServer", "Failed to start service."+err.Error())
	//	}

}

func ListenAndServe(addr string, handler http.Handler, timeout time.Duration) error {
	server := &http.Server{
		Addr:        addr,
		Handler:     handler,
		ReadTimeout: timeout,
	}
	return server.ListenAndServe()
}
