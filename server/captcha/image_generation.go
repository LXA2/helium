package main

import (
	"encoding/base64"
	"fmt"
	"image/color"
	"image/jpeg"
    "os"
	"github.com/afocus/captcha"
    "unicode/utf8"
	//"log"
	//"net/http"
	"strings"
)

func GenerateCaptcha(text string) (string, error) {
	cap := captcha.New()

	// 设置字体
	if err := cap.SetFont("./font/云峰静龙行书.ttf"); err != nil { // 请确保字体文件路径正确
		return "", err
	}
    println("length:",utf8.RuneCountInString(text))
	// 设置图片尺寸
	//cap.SetSize(240, 80)
    cap.SetSize(40*utf8.RuneCountInString(text), 60)

	// 设置干扰强度
	//cap.SetDisturbance(captcha.MEDIUM)
    cap.SetDisturbance(50)

	// 设置前景色
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255}, color.RGBA{0, 0, 0, 255})

	// 设置背景色
	cap.SetBkgColor(color.RGBA{100, 100, 100, 255})

	// 生成验证码图片
	// img, err := cap.CreateCustom(text)
	// if err != nil {
	//     return "", err
	// }
	img := cap.CreateCustom(text)
    
	// 将图片编码为 base64
	var buffer strings.Builder
	encoder := base64.NewEncoder(base64.StdEncoding, &buffer)
	if err := jpeg.Encode(encoder, img, nil); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

/*
	func captchaHandler(w http.ResponseWriter, r *http.Request) {
	    text := r.URL.Query().Get("text")
	    if text == "" {
	        http.Error(w, "text parameter is required", http.StatusBadRequest)
	        return
	    }

	    base64Image, err := generateCaptcha(text)
	    if err != nil {
	        http.Error(w, fmt.Sprintf("failed to generate captcha: %v", err), http.StatusInternalServerError)
	        return
	    }

	    w.Header().Set("Content-Type", "application/json")
	    w.Write([]byte(fmt.Sprintf(`{"captcha":"%s"}`, base64Image)))
	}

	func main() {
	    http.HandleFunc("/captcha", captchaHandler)
	    log.Fatal(http.ListenAndServe(":8080", nil))
	}
*/
func main() {
	var c string
	n, err := fmt.Scanln(&c)
	fmt.Println(c)
	if err != nil {
		fmt.Println("Error:", err, n)
	} else {
		fmt.Println("Input:", c, n)
	}
	captcha, err := GenerateCaptcha(c)
	if err != nil {
		fmt.Println("Failed to generate captcha:", err)
	} else {
		htmlImgTag := fmt.Sprintf("<p></p><p></p><p></p><p></p><p></p><p></p><img src=\"data:image/jpeg;base64,%s\" />", captcha)
		fmt.Println("HTML img tag:")
		fmt.Println(htmlImgTag)

        // 将 HTML img 标签追加到文件末尾
        filePath := "C:/Users/Administrator/Desktop/1.html"
        file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
        if err != nil {
            fmt.Println("Error opening file:", err)
            return
        }
        defer file.Close()

        if _, err := file.WriteString(htmlImgTag + "\n"); err != nil {
            fmt.Println("Error writing to file:", err)
        } else {
            fmt.Println("Successfully appended to file")
        }
	}
}
