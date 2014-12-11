package box

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var boxAPIURL = "https://api.box.com"

// Tokens for Box API.  accessToken, refreshToken, clientID, and clientSecret
type Tokens struct {
	AccessToken  string
	RefreshToken string
	ClientID     string
	ClientSecret string
}

// RefreshToken - struct for containing the refresh token response from box
type RefreshToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RestrictedTo string `json:"restricted_to"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
}

// UserResponse - response from box api for Users
type UserResponse struct {
	TotalCount int64  `json:"total_count"`
	Entries    []User `json:"entries"`
	Limit      int64  `json:"limit"`
	Offset     int64  `json:"offset"`
}

//User struct to identiy user information
type User struct {
	MethodType    string `json:"type"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Login         string `json:"login"`
	CreatedAt     string `json:"created_at"`
	ModifiedAt    string `json:"modified_at"`
	Language      string `json:"language"`
	Timezone      string `json:"timezone"`
	SpaceAmount   int64  `json:"space_amount"`
	SpaceUsed     int64  `json:"space_used"`
	MaxUploadSize int64  `json:"max_upload_size"`
	Status        string `json:"status"`
	JobTitle      string `json:"job_title"`
	Phone         string `json:"phone"`
	Address       string `json:"address"`
	AvatarURL     string `json:"avatar_url"`
}

// execute get against box api
func (token Tokens) get(method string) ([]byte, error) {
	url := boxAPIURL + method

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Println(resp)
	//fmt.Println("Reponse Status Code ", resp.StatusCode)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil

}

//RefreshTokens - Refresh Box API tokens to continue working with API
func (token Tokens) RefreshTokens() {
	urlStr := boxAPIURL + "/api/oauth2/token"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Add("refresh_token", token.RefreshToken)
	data.Add("client_id", token.ClientID)
	data.Add("client_secret", token.ClientSecret)

	client := &http.Client{}
	req, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Resonse Status Code ", resp.StatusCode)
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(string(body))
	default:
		fmt.Println("Something bad happened updating Refresh Token")
	}
}

//GetUsers -  Get list of Box Users
func (token Tokens) GetUsers(limit, offset int) UserResponse {

	//var total int
	//var limit = 100
	//var offset int

	urlStr := fmt.Sprintf("/2.0/users?limit=%d&offset=%d", limit, offset)
	//fmt.Println("urlStr:", urlStr)
	//body, err := token.get("/2.0/users")
	body, err := token.get(urlStr)

	var data UserResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		switch v := err.(type) {
		case *json.SyntaxError:
			fmt.Println(string(body[v.Offset-40 : v.Offset]))
		}
	}

	return data
	/*if data.TotalCount > data.Offset {
	    fmt.Println("Still more records!")
	  }
	  //total = data.TotalCount
	  //fmt.Println(data.TotalCount)
	  //fmt.Println(data.Limit)
	  //fmt.Println(data.Offset)
	  //fmt.Println(data)
	  fmt.Println("Total Count: ", data.TotalCount )
	  fmt.Println("Limit: ", data.Limit)
	  fmt.Println("Offset: ", data.Offset)
	  for _, user :=range data.Entries{
	  fmt.Println(user.MethodType, user.ID, user.Name, user.Login, user.CreatedAt, user.ModifiedAt, user.Language, user.Timezone, user.SpaceAmount, user.SpaceUsed, user.MaxUploadSize, user.Status, user.JobTitle, user.Phone, user.Address, user.AvatarURL )
	  }*/
}

//UpdateUserSpace - Update the UserSpace for a box account. Required userid int, and newsize int.
func (token Tokens) UpdateUserSpace(userid, newsize int) {
	boxURL := boxAPIURL + "/v1/users/" + string(userid)
	data := url.Values{}
	data.Set("space_amount", string(newsize))

	client := &http.Client{}
	r, _ := http.NewRequest("PUT", boxURL, bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	r.Header.Add("Authorization", "Bearer "+token.AccessToken)
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)
	fmt.Println(resp.Status)
}
