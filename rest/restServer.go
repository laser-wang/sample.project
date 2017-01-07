// 包说明见doc.go
// 作者：laser.wang

package rest

import (
	"encoding/json"
	//	"fmt"
	"io"
	"net/http"
	//	"os"
	//	"runtime"
	"time"

	"database/sql"
	//数据库连接池
	//	"github.com/lib/pq"
	_ "github.com/lib/pq"

	"laserTools"

	//日志
	l4g "github.com/alecthomas/log4go"
	//redis驱动
	"github.com/garyburd/redigo/redis"
	//平滑重启和安全关闭服务
	"github.com/tabalt/gracehttp"
)

var (
	gRedisClient *redis.Pool
	gRedisHost   string
	gRedisDB     int

	gRestDB *sql.DB
)

type structResult struct {
	Cd  int
	Msg string
}

//初始化，包括日志、数据库等
func init() {

	var err error

	//初始化日志
	l4g.AddFilter("stdout", l4g.DEBUG, l4g.NewConsoleLogWriter())
	l4g.AddFilter("file", l4g.DEBUG, l4g.NewFileLogWriter("rest.log", false))

	// 初始化数据库,已自带连接池
	gRestDB, err = sql.Open("postgres", "postgres://avcp_work:avcp_work@192.168.0.10/avcp_work?sslmode=disable")
	// 连接的最大空闲时间(可选)
	gRestDB.SetConnMaxLifetime(300 * time.Second)
	// 默认不限制
	gRestDB.SetMaxOpenConns(2000)
	// 默认为0
	gRestDB.SetMaxIdleConns(100)
	checkErr(err, check_flag_exit)

	// 初始化redis
	gRedisHost = "192.168.0.10:6379"
	gRedisDB = 1
	// 建立连接池
	gRedisClient = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     10,
		MaxActive:   100,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", gRedisHost)
			if err != nil {
				return nil, err
			}
			// 选择db
			c.Do("SELECT", gRedisDB)
			return c, nil
		},
	}

}

// StartServer 启动服务入口
func StartServer() {

	// 连接数
	gRestDB.SetMaxOpenConns(100) // 最大连接数
	gRestDB.SetMaxIdleConns(20)  // 最大空闲数

	startHttp()
	defer releaseResource()
}

// login 用户登陆，登陆后会生成token，返回给客户端，下次访问非login接口时header需要包含正确的token.
// 如果token不正确，则不会返回正确结果
func login(w http.ResponseWriter, req *http.Request) {

	//	loginToken := req.Header.Get("loginToken")
	//	l4g.Log(l4g.INFO, "", "token:"+loginToken)

	userCode := req.FormValue("userCode")
	pwd := req.FormValue("pwd")

	if userCode == "" || pwd == "" {
		io.WriteString(w, getResult(500, "Login is failed,please use correct username or password."))
		l4g.Log(l4g.INFO, "", "Login is failed,please use correct username or password.")
		return
	}

	var err error
	//查询数据
	rows, err := gRestDB.Query("SELECT user_name FROM xx_user where user_code = $1 and pwd = $2", userCode, pwd)

	//	rows, err := gRestDB.Query("SELECT user_name FROM xx_user")

	checkErr(err, check_flag_log)

	if rows.Next() == true {
		var username string
		err = rows.Scan(&username)
		checkErr(err, check_flag_log)

		loginToken := laserTools.GetToken()
		redisConn := gRedisClient.Get()

		redisConn.Do("SET", userCode, loginToken)
		redisConn.Do("EXPIRE", userCode, 30)

		defer redisConn.Close()
		io.WriteString(w, getResult(200, loginToken))
		l4g.Log(l4g.INFO, "", "Login is successful."+getResult(200, username))
	} else {
		io.WriteString(w, getResult(500, "Login is failed,please use correct username or password."))
		l4g.Log(l4g.INFO, "", "Login is failed,please use correct username or password.")
	}

	defer rows.Close()

}

//把接口返回值封装成json串
func getResult(cd int, msg string) string {
	rslt := structResult{
		Cd:  cd,
		Msg: msg,
	}
	ret, _ := json.Marshal(rslt)
	return string(ret)
}

// startHttp 启动http服务
func startHttp() {

	//	addr := "127.0.0.1:12345"

	// test
	//http://192.168.0.10:12345/login?userName=ccc&pwd=rrr
	//kill -SIGUSR2 $pid 	平滑重启
	//kill $pid				关闭服务

	http.HandleFunc("/login", login)
	l4g.Log(l4g.INFO, "", "Server start.")

	err := gracehttp.ListenAndServe(":12345", nil)
	if err != nil {

	}

	l4g.Log(l4g.INFO, "", "Server stop")

	//	http.HandleFunc("/login", login)

	//	l4g.Log(l4g.INFO, "", "Server start.")

	//	err := http.ListenAndServe(":12345", nil)
	//	if err != nil {
	//		l4g.Log(l4g.ERROR, "startServer", "Failed to start service."+err.Error())
	//	}

}

//释放日志文件、数据库连接等资源
func releaseResource() {
	l4g.Close()
	gRestDB.Close()
	gRedisClient.Close()
}

//func ListenAndServe(addr string, handler http.Handler, timeout time.Duration) error {
//	server := &gracehttp.Server{
//		Addr:        addr,
//		Handler:     handler,
//		ReadTimeout: timeout,
//	}
//	return server.ListenAndServe()
//}

//错误处理函数，程序中错误分2种，一种需要终止程序，一种仅仅是记录错误日志
//flag:  1:exit  2:log only
func checkErr(err error, flag int) {

	if err != nil {
		switch flag {
		case check_flag_exit:
			l4g.Log(l4g.ERROR, "", err.Error())
			panic(err)
		case check_flag_log:
			l4g.Log(l4g.ERROR, "", err.Error())
		default:
			l4g.Log(l4g.ERROR, "", err.Error())
		}
	}

}

// setToken 登陆后设置token
//func setToken(String userCode) {

//}
