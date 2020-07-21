package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/pat"
	"github.com/urfave/negroni"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = oauth2.Config{
	RedirectURL:  "http://localhost:3000/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateStateOauthCookie(w)
	url := googleOauthConfig.AuthCodeURL(state)            //user를 어떤 경로로 보내야할지 알려준다.
	http.Redirect(w, r, url, http.StatusTemporaryRedirect) //user를 이쪽 경로로 바로 redirect, 구글 로그인 폼이 뜬다.
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour) //24시간 유효

	b := make([]byte, 16) //16바이트 수
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
}

func googleAuthCallback(w http.ResponseWriter, r *http.Request) {
	oauthstate, _ := r.Cookie("oauthstate")       //request에 쿠키 담겨있음, Google에서 state를 request로 보내줌
	if r.FormValue("state") != oauthstate.Value { //우리가 로그인한거 아닐 때(공격시도)
		log.Printf("invalid google oauth state cookie: %s state: %s\n", oauthstate.Value)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	data, err := getGoogleUserInfo(r.FormValue("code")) //이 code를 가지고 userinfo를 요청해서 get
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect) //에러 생기면 메인으로 redirect
		return
	}
	fmt.Fprint(w, string(data)) //에러 없으면 user한테 다시 userinfo를 보내줌
}

//userinfo를 request하는 경로, access_token을 우리가 넣어서 보내면 Google이 userinfo 보내줌
const oauthGoogleURLAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func getGoogleUserInfo(code string) ([]byte, error) {
	// context는 스레드 간의 데이터 주고받을 때 사용되는 safe 저장공간
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}

	resp, err := http.Get(oauthGoogleURLAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get UserInfo %s\n", err.Error())
	}

	return ioutil.ReadAll(resp.Body)
}

func main() {
	mux := pat.New()
	mux.HandleFunc("/auth/google/login", googleLoginHandler)
	mux.HandleFunc("/auth/google/callback", googleAuthCallback)

	n := negroni.Classic()
	n.UseHandler(mux)
	http.ListenAndServe(":3000", n)
}
