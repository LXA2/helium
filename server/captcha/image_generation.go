package main

import (
    "encoding/base64"
    "fmt"
    "github.com/afocus/captcha"
    "image/color"
	"image/jpeg"
    "log"
    "net/http"
    "strings"
)

func generateCaptcha(text string) (string, error) {
    cap := captcha.New()

    // 设置字体
    if err := cap.SetFont("./comic.ttf"); err != nil { // 请确保字体文件路径正确
        return "", err
    }

    // 设置图片尺寸
    cap.SetSize(240, 80)

    // 设置干扰强度
    cap.SetDisturbance(captcha.MEDIUM)

    // 设置前景色
    cap.SetFrontColor(color.RGBA{255, 255, 255, 255}, color.RGBA{0, 0, 0, 255})

    // 设置背景色
    cap.SetBkgColor(color.RGBA{100, 100, 100, 255})

    // 生成验证码图片
    // img, err := cap.CreateCustom(text)
    // if err != nil {
    //     return "", err
    // }
	img:= cap.CreateCustom(text)

    // 将图片编码为 base64
    var buffer strings.Builder
    encoder := base64.NewEncoder(base64.StdEncoding, &buffer)
    if err := jpeg.Encode(encoder, img, nil); err != nil {
        return "", err
    }

    return buffer.String(), nil
}

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
