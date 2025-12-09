//++++++++++++++++++++++++++++++++++++++++
// 《Go语言Web编程实战》源码
//++++++++++++++++++++++++++++++++++++++++
// 作者公众号：源码大数据
// 自媒体账号（抖音、视频号、快手、B站、知乎）：廖显东-ShirDon
// Blog:https://www.shirdon.com/
// 仓库地址：https://gitee.com/shirdonl/goWebActualCombatV2
// 仓库地址：https://github.com/shirdonl/goWebActualCombatV2
//++++++++++++++++++++++++++++++++++++++++

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

// const clientID = "<your client id>"
const clientID = "Ova5Jf"

// const clientSecret = "<your client secret>"
const clientSecret = "sdfgbdfgf96849ad"

func hello(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("hello.html")
		t.Execute(w, nil)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	}
}

func main() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/", hello)
	http.HandleFunc("/hello", hello)

	httpClient := http.Client{}
	http.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		code := r.FormValue("code")

		//完整URL请查看GitHub官网
		reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?"+
			"client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
		req, err := http.NewRequest(http.MethodPost, reqURL, nil)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}
		req.Header.Set("accept", "application/json")

		res, err := httpClient.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		defer res.Body.Close()

		var t AccessTokenResponse
		if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
			fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
			w.WriteHeader(http.StatusBadRequest)
		}

		w.Header().Set("Location", "/hello.html?access_token="+t.AccessToken)
		w.WriteHeader(http.StatusFound)
	})

	http.ListenAndServe(":8087", nil)
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}
