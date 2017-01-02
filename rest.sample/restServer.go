package main

import (
	"encoding/json"
	//	"fmt"
	"io"
	"net/http"
	//	"os"
	//	"runtime"
	//	"time"

	"database/sql"
	//数据库连接池
	"github.com/lib/pq"
	_ "github.com/lib/pq"

	//日志
	l4g "github.com/alecthomas/log4go"
	//redis驱动
	"github.com/garyburd/redigo/redis"
	//平滑重启和安全关闭服务
	"github.com/tabalt/gracehttp"
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

	defer releaseResource()

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

	//	addr := "127.0.0.1:12345"

	// test
	//http://192.168.0.10:12345/login?userName=ccc&pwd=rrr
	//kill -SIGUSR2 $pid 	平滑重启
	//kill $pid				关闭服务

	http.HandleFunc("/login", login)
	l4g.Log(l4g.INFO, "", "Server start.")

	err := gracehttp.ListenAndServe(":12345", nil)
	if err != nil {
		l4g.Log(l4g.INFO, "", "Server stop")
	}

	//	http.HandleFunc("/login", login)

	//	l4g.Log(l4g.INFO, "", "Server start.")

	//	err := http.ListenAndServe(":12345", nil)
	//	if err != nil {
	//		l4g.Log(l4g.ERROR, "startServer", "Failed to start service."+err.Error())
	//	}

}

func releaseResource() {
	l4g.Close()
	db.Close()
}

//func ListenAndServe(addr string, handler http.Handler, timeout time.Duration) error {
//	server := &gracehttp.Server{
//		Addr:        addr,
//		Handler:     handler,
//		ReadTimeout: timeout,
//	}
//	return server.ListenAndServe()
//}
