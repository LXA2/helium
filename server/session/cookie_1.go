package session

//   cookie_1  is the cookie that is included in the response headers (unsafe)

import (
	"fmt"
	"time"
	"net/http"
	"bufio"
	"crypto/rand"
	"encoding/hex"
	//"io/ioutil"
	"os"
	"strings"
)

func GetCookie_1(headers *http.Header, ipvAndIpAndPort *[]string) (*http.Cookie,bool){
	/*
	cookieInHeaders:= headers.Get("Cookie")
	fmt.Println("cookie in headers:",cookieInHeaders)
	now := time.Now()
	cookie := &http.Cookie{
		Name:   "time",
		Value:  now.Format("2006-01-02 15:04:05"), // 07:00"),
		MaxAge: 20,
		//Secure: true,//Open when https is available
	}
		*/
	cookie,b :=getCookie(headers,ipvAndIpAndPort)
	return cookie,b
}

// 生成指定长度的随机字符串
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// 读取存储 Cookie 的文件
func readCookiesFromFile(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}
		return nil, err
	}
	defer file.Close()

	cookies := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) == 2 {
			cookies[parts[0]] = parts[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cookies, nil
}

// 写入 Cookie 到文件
func writeCookieToFile(filePath, cookie, expiration string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s,%s\n", cookie, expiration))
	return err
}

// 记录客户端请求信息
func logClientRequest(logDir, cookie string, clientInfo []string, headers *http.Header) error {
	logFile := logDir + "/" + cookie
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s,%s,%s\nHeaders: %v\n",
		clientInfo[0], clientInfo[1], clientInfo[2], *headers))
	return err
}

// GetCookie 函数
func getCookie(headers *http.Header, ipvAndIpAndPort *[]string) (*http.Cookie, bool) {
	const cookieFilePath = "./session/cookies1.txt"
	const logDir = "./session/logs_1"

	// 读取现有的 Cookies
	cookies, err := readCookiesFromFile(cookieFilePath)
	if err != nil {
		fmt.Println("Error reading cookie file:", err)
		return nil, false
	}

	// 检查传入的 Headers 是否包含 Cookie
	cookieHeader := headers.Get("Cookie")
	if cookieHeader != "" {
		cookieParts := strings.Split(cookieHeader, "=")
		if len(cookieParts) == 2 {
			existingCookie := cookieParts[1]
			if _, exists := cookies[existingCookie]; exists {
				// 如果 Cookie 存在，记录请求信息并返回
				err := logClientRequest(logDir, existingCookie, *ipvAndIpAndPort, headers)
				if err != nil {
					fmt.Println("Error logging client request:", err)
				}
				return nil, false
			}
		}
	}

	// 生成新的 Cookie
	var newCookie string
	for {
		newCookie, err = generateRandomString(200)
		if err != nil {
			fmt.Println("Error generating random string:", err)
			return nil, false
		}
		if _, exists := cookies[newCookie]; !exists {
			break
		}
	}

	expiration := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	err = writeCookieToFile(cookieFilePath, newCookie, expiration)
	if err != nil {
		fmt.Println("Error writing to cookie file:", err)
		return nil, false
	}

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   newCookie,
		Expires: time.Now().Add(24 * time.Hour),
	}

	// 记录请求信息
	err = logClientRequest(logDir, newCookie, *ipvAndIpAndPort, headers)
	if err != nil {
		fmt.Println("Error logging client request:", err)
	}

	return cookie, true
}