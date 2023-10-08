package actions

import (
	"bytes"
	"client/dto"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gobuffalo/buffalo"
)

func LoginIndexHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("login/index.plush.html"))
}

func LoginHandler(c buffalo.Context) error {
	loginReq := &dto.LoginReq{}
	if err := c.Bind(loginReq); err != nil {
		log.Printf("Error binding the request: %v\n", err)
		c.Flash().Add("danger", "Error processing your request")
		return c.Redirect(http.StatusSeeOther, "/")
	}

	body, err := json.Marshal(loginReq)
	if err != nil {
		log.Printf("Error on json marshal: %v\n", err)
		c.Flash().Add("danger", "Error processing your request")
		return c.Redirect(http.StatusSeeOther, "/")
	}

	client := http.Client{
		Timeout:   5 * time.Second,
		Transport: &http.Transport{},
	}

	endpoint := fmt.Sprintf("%s/login", os.Getenv("API_URL"))
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error http new request: %v\n", err)
		c.Flash().Add("danger", "Error processing your request")
		return c.Redirect(http.StatusSeeOther, "/")
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Printf("Error client do request: %v\n", err)
		c.Flash().Add("danger", "Error processing your request")
		return c.Redirect(http.StatusSeeOther, "/")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		c.Flash().Add("danger", "Invalid email or password")
		return c.Redirect(http.StatusSeeOther, "/")
	}

	loginRes := dto.LoginRes{}
	if err := json.NewDecoder(res.Body).Decode(&loginRes); err != nil {
		c.Flash().Add("danger", "Invalid email or password")
		return c.Redirect(http.StatusUnauthorized, "/")
	}

	cookies := res.Cookies()
	for _, cookie := range cookies {
		c.Cookies().Set(cookie.Name, cookie.Value, 12*time.Hour)
	}

	return c.Redirect(http.StatusSeeOther, "homePath()")
}
