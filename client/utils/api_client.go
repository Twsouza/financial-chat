package utils

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gobuffalo/buffalo"
)

func CallApi(c buffalo.Context, path string, method string, body []byte) (*http.Response, error) {
	client := http.Client{
		Timeout:   5 * time.Second,
		Transport: &http.Transport{},
	}

	endpoint := fmt.Sprintf("%s%s", os.Getenv("API_URL"), path)
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	setCookies(c, req)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func setCookies(c buffalo.Context, req *http.Request) {
	// ignoring err because the auth would have redirect to the login page if doesn't have the cookies
	tokenCookie, _ := c.Cookies().Get("token")
	userIDCookie, _ := c.Cookies().Get("userId")
	usernameCookie, _ := c.Cookies().Get("username")

	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tokenCookie,
	})

	req.AddCookie(&http.Cookie{
		Name:  "userId",
		Value: userIDCookie,
	})

	req.AddCookie(&http.Cookie{
		Name:  "username",
		Value: usernameCookie,
	})
}
