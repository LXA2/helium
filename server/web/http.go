package web

import (
	"fmt"
	"helium/session"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	//"path/filepath"
	"strings"
)

// 处理根路径的请求
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n---------------------------------------------------------------")
	// 打印请求方法
	fmt.Println("Method:", r.Method)

	// 打印请求路径
	fmt.Println("Path:", r.URL.Path)

	// 打印查询参数
	fmt.Println("Query Parameter 'name':", r.URL.Query().Get("name"))

	// 打印请求头
	fmt.Println("Headers:", r.Header) //.Get("User-Agent"))

	// 打印请求头
	fmt.Println("Cookie:", r.Header.Get("Cookie"))

	//
	//fmt.Printf("%T",r.RemoteAddr)//string
	fmt.Println("Address:", r.RemoteAddr)

	var ipAddress []string
	ipAddressAndPort := r.RemoteAddr
	if ipAddressAndPort[0] == 91 { //==91时为ipv6的情况，第一个字符为"[",ascii为91
		ipAddress = []string{"ipv6",strings.Split(ipAddressAndPort, "]")[0][1:],strings.Split(ipAddressAndPort, "]")[1][1:]} // 创建 IP 地址的切片
	} else { //ipv4
		//ipAddressAndPort := strings.Split(r.RemoteAddr, ":")[0] // 提取 IP 地址部分
		ipAddress = []string{"ipv4",strings.Split(ipAddressAndPort, ":")[0],strings.Split(ipAddressAndPort, ":")[1]} 
	}

	// 读取并打印请求主体
	body, err := io.ReadAll(r.Body)
	if err == nil {
		fmt.Println("Body:", string(body))
	}
	r.Body.Close()

	// 获取请求的路径
	reqPath := r.URL.Path
	println(reqPath)

	//w.Header().Set("Content-Type", "application/pdf")
	//w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Type", "text/html")

	cookie1,ifSetNewCookie1 := session.GetCookie_1(&r.Header, &ipAddress)
	/*now := time.Now()
	cookie := &http.Cookie{
		Name:   "time",
		Value:  now.Format("2006-01-02 15:04:05"), // 07:00"),
		MaxAge: 20,
		//Secure: true,//Open when https is available
	}*/
	if ifSetNewCookie1 {
		http.SetCookie(w, cookie1)
	}
	

	// http.ServeFile(w, r, "./statics/index.html") //filePath)
	filePath := "./statics" + reqPath
	fmt.Printf("%T", reqPath)
	// 检查文件是否存在并且是合法的文件
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) || info.IsDir() {
		//http.Error(w, "File not found.\n"+filePath, http.StatusNotFound) //"File not found.",http.StatusNotFound)
		return
	}

	// 返回文件
	http.ServeFile(w, r, filePath) //"./statics/aaa.html") // filePath)
}

func StartServer() {
	// 创建一个新的ServeMux
	mux := http.NewServeMux()

	// 注册处理函数
	mux.HandleFunc("/", handler)
	//mux.HandleFunc("/handler2", handler2)

	// 创建一个http.Server实例
	server := &http.Server{
		Addr:         ":3344",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		//MaxHeaderBytes: 1 << 20,
	}

	// 启动服务器
	log.Println("Starting server on :3344")
	//err := server.ListenAndServeTLS("web/cert.pem", "web/key.pem")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
