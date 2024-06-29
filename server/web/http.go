package web

import (
	"fmt"
	"net/http"
	"os"
	//"path/filepath"
	//"strings"
)

// 处理根路径的请求
func handler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, r.URL.Path)

	// 获取请求的路径
	reqPath := r.URL.Path
	println(reqPath)
	/*
			// 检查路径前缀
			if !strings.HasPrefix(reqPath, "/lkjhgfdsa/mnbvcxz/qwertyuiop/") {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// 移除前缀
			reqPath = strings.TrimPrefix(reqPath, "/lkjhgfdsa/mnbvcxz/qwertyuiop/")

		// 将请求路径映射到站点目录
		siteDir := "../statics"
		filePath := filepath.Join(siteDir, reqPath)
		//println(filePath)
	*/
	http.ServeFile(w, r, "./statics/index.html") //filePath)
	filePath := reqPath
	// 检查文件是否存在并且是合法的文件
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) || info.IsDir() {
		http.Error(w, "File not found.\n"+filePath, http.StatusNotFound) //"File not found.",http.StatusNotFound)
		return
	}

	// 返回文件

}

// 启动HTTPS服务器
func StartServer() {
	// 注册处理函数到根路径
	http.HandleFunc("/", handler)

	// 指定证书和密钥文件
	/*certFile := "web/cert.pem"
	keyFile := "web/key.pem"*/

	// 启动HTTPS服务器
	addr := ":3333"
	fmt.Printf("HTTPS Server is listening on %s\n", addr)
	/*if err := http.ListenAndServeTLS(addr, certFile, keyFile, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}*/
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
