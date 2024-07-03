package session

//   cookie_1  is the cookie that is included in the response headers (unsafe)

import (
	"fmt"
	"time"
	//"fmt"
	"net/http"
)

func GetCookie_1(headers *http.Header, ipvAndIpAndPort *[]string,) (*http.Cookie,bool){
	cookieInHeaders:= headers.Get("Cookie")
	fmt.Println("cookie in headers:",cookieInHeaders)
	now := time.Now()
	cookie := &http.Cookie{
		Name:   "time",
		Value:  now.Format("2006-01-02 15:04:05"), // 07:00"),
		MaxAge: 20,
		//Secure: true,//Open when https is available
	}
	return cookie,true
}