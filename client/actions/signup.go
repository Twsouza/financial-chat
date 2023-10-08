package actions

import (
	"bytes"
	"client/dto"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gobuffalo/buffalo"
)

func SignupIndexHandler(c buffalo.Context) error {
	return c.Render(http.StatusOK, r.HTML("signup/index.plush.html"))
}

func SignupHandler(c buffalo.Context) error {
	signupReq := &dto.SignupReq{}
	if err := c.Bind(signupReq); err != nil {
		return err
	}

	body, err := json.Marshal(signupReq)
	if err != nil {
		return err
	}

	client := http.Client{
		Timeout:   5 * time.Second,
		Transport: &http.Transport{},
	}

	endpoint := fmt.Sprintf("%s/signup", os.Getenv("API_URL"))
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		c.Flash().Add("danger", "Invalid credentials, please check the form")
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	return c.Redirect(http.StatusSeeOther, "homePath()")
}
