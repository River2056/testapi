package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
)

const (
	LOGIN_URL = "/api/login/to_login"
	GET       = "get"
	POST      = "post"
)

var environments map[string]string = map[string]string{
	"localhost": "http://localhost:8080",
	"thor":      "https://secure-thor.crm-alpha.com",
	"leo":       "https://secure-leo-ex.crm-alpha.com",
	"venus":     "https://secure-venus-ex.crm-alpha.com",
	"virgo":     "https://secure-virgo-ex.crm-alpha.com",
}

func logError(err error, msg string) {
	if err != nil {
		log.Fatalf("%v: %v\n", msg, err)
	}
}

func main() {
	env := flag.String("e", "localhost", "environment name")
	username := flag.String("u", "deposit@test.com", "login user")
	password := flag.String("p", "123Qwe", "login user password")
	urlEndpoint := flag.String("r", "/api/campaign/eligible-campaigns", "test url endpoint")
	dataPayload := flag.String("d", "", "post request body (payload)")
	// requestMethod := flag.String("m", "get", "request method")
	flag.Parse()

	payload := make(map[string]interface{})
	json.Unmarshal([]byte(*dataPayload), &payload)

	loginUrl := fmt.Sprintf("%v%v", environments[*env], LOGIN_URL)
	testEndPoint := fmt.Sprintf("%v%v", environments[*env], *urlEndpoint)

	fmt.Printf("env: %v\n", environments[*env])
	fmt.Printf("username: %v\n", *username)
	fmt.Printf("password: %v\n", *password)
	fmt.Printf("url endpoint: %v\n", *urlEndpoint)
	fmt.Printf("payload: %v\n", payload)
	var requestMethod string
	if len(payload) == 0 {
		requestMethod = GET
	} else {
		requestMethod = POST
	}
	// fmt.Printf("requestMethod: %v\n", *requestMethod)
	fmt.Printf("login url: %v\n", loginUrl)
	fmt.Printf("url to test: %v\n", testEndPoint)
	fmt.Println()

	result, err := exec.Command("node", "/root/go/src/testapi/rsa_components/index.js", *username).Output()
	logError(err, "error running node command")

	passwordHash := md5.Sum([]byte(*password))

	// login data
	loginData := map[string][]string{
		"userName_login": {strings.TrimSpace(string(result))},
		"password_login": {hex.EncodeToString(passwordHash[:])},
		"utc":            {"39600000"},
	}
	logError(err, "error while marshaling login json data")

	res, err := http.PostForm(loginUrl, url.Values(loginData))
	logError(err, "error while logging in")

	if res.StatusCode != 200 {
		log.Fatalf("login failed: %v\n", res)
	}

	var jsonRes map[string]interface{}
	json.NewDecoder(res.Body).Decode(&jsonRes)
	cookie := fmt.Sprintf("%v=%v", res.Cookies()[0].Name, res.Cookies()[0].Value)
	token := jsonRes["data"].(map[string]interface{})

	var req *http.Request
	switch requestMethod {
	case GET:
		req, err = http.NewRequest(requestMethod, testEndPoint, nil)
		logError(err, "error while making a get request")

	case POST:
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(payload)
		req, err = http.NewRequest(requestMethod, testEndPoint, b)
		logError(err, "error while making a post request")
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	req.Header.Add("cookie", cookie)
	req.Header.Add("token", token["accessToken"].(string))

	response, err := http.DefaultClient.Do(req)
	logError(err, "error while receiving a response from test end point: %v\n")
	defer response.Body.Close()

	var jsonResponse map[string]interface{}
	json.NewDecoder(response.Body).Decode(&jsonResponse)
	output, err := json.MarshalIndent(jsonResponse, "", "  ")
	logError(err, "error while marshaling json output")

	fmt.Println()
	fmt.Println(string(output))
}
