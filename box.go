package box

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
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
	AccessToken  string        `json:"access_token"`
	ExpiresIn    int           `json:"expires_in"`
	RestrictedTo []interface{} `json:"restricted_to"` //trying to force RestrictedTo to be a string and not an array
	RefreshToken string        `json:"refresh_token"`
	TokenType    string        `json:"token_type"`
}

// UserResponse - response from box api for Users
type UserResponse struct {
	TotalCount int64  `json:"total_count"`
	Entries    []User `json:"entries"`
	Limit      int64  `json:"limit"`
	Offset     int64  `json:"offset"`
}

//GroupResponse - reponse from box api for groups
type GroupResponse struct {
	TotalCount int64   `json:"total_count"`
	Entries    []Group `json:"entries"`
	Limit      int64   `json:"limit"`
	Offset     int64   `json:"offset"`
}

//User struct to identiy user information
type User struct {
	MethodType    string  `json:"type"`
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Login         string  `json:"login"`
	CreatedAt     string  `json:"created_at"`
	ModifiedAt    string  `json:"modified_at"`
	Language      string  `json:"language"`
	Timezone      string  `json:"timezone"`
	SpaceAmount   float64 `json:"space_amount"`
	SpaceUsed     float64 `json:"space_used"`
	MaxUploadSize float64 `json:"max_upload_size"`
	Status        string  `json:"status"`
	JobTitle      string  `json:"job_title"`
	Phone         string  `json:"phone"`
	Address       string  `json:"address"`
	AvatarURL     string  `json:"avatar_url"`
}

//Group struct to identify group information
type Group struct {
	MethodType             string `json:"type"`
	ID                     string `json:"id"`
	Name                   string `json:"name"`
	CreatedAt              string `json:"created_at"`
	ModifiedAt             string `json:"modified_at"`
	Provenance             string `json:"provenance"`
	ExternalSyncIdentifier string `json:"external_sync_identifier"`
	Description            string `json:"description"`
}

// execute get against box api
func (token *RefreshToken) get(method string) ([]byte, error) {
	url := boxAPIURL + method

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	if err != nil {
		debug(httputil.DumpRequestOut(req, true))
	}
	resp, err := client.Do(req)
	if err != nil {
		debug(httputil.DumpResponse(resp, true))
		log.Fatalln(err)
	}
	//fmt.Println(resp)
	fmt.Println("Refreshing Box API Token.")
	switch resp.StatusCode {
	case 200:

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return body, nil

	case 401:
		fmt.Println("Box Token is invalid. Please Renew and Try again.")
		log.Fatalln("Invalid Box Token.")
		return nil, err
	}
	return nil, nil
}

//RefreshTokens - Refresh Box API tokens to continue working with API
func (token *RefreshToken) RefreshTokens(clientID string, clientSecret string) {
	urlStr := boxAPIURL + "/oauth2/token"
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Add("refresh_token", token.RefreshToken)
	data.Add("client_id", clientID)
	data.Add("client_secret", clientSecret)

	fmt.Println("Updating Box Token")
	client := &http.Client{}
	req, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Println("Resonse Status Code ", resp.StatusCode)
	defer resp.Body.Close()
	switch resp.StatusCode {
	case 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		//fmt.Println(string(body))
		//var tmptoken RefreshToken

		err = json.Unmarshal(body, &token)
		if err != nil {
			fmt.Println(err)
			switch v := err.(type) {
			case *json.SyntaxError:
				fmt.Println(string(body[v.Offset-40 : v.Offset]))
			}
		}
		//fmt.Println("Access Token: ", token.AccessToken)
		//fmt.Println("Refresh Token: ", token.RefreshToken)
	case 400:
		fmt.Println("Refresh token has expired")
		log.Fatalln(resp.Body)
	default:
		fmt.Println("Something bad happened updating Refresh Token")
	}
}

//GetUsers -  Get list of Box Users
func (token *RefreshToken) GetUsers(limit, offset int) (UserResponse, int) {
	//var total int
	//var limit = 100
	//var offset int

	urlStr := fmt.Sprintf("%s/2.0/users?limit=%d&offset=%d", boxAPIURL, limit, offset)
	//fmt.Println("urlStr:", urlStr)
	//body, err := token.get("/2.0/users")
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	if err != nil {
		debug(httputil.DumpRequestOut(req, true))
	}
	resp, err := client.Do(req)
	//fmt.Println("Box Response: ", resp)
	if err != nil {
		debug(httputil.DumpResponse(resp, true))
		log.Fatalln(err)
	}

	var data UserResponse
	switch resp.StatusCode {
	case 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln("err")
		}

		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println(err)
			switch v := err.(type) {
			case *json.SyntaxError:
				fmt.Println(string(body[v.Offset-40 : v.Offset]))
			}
		}
		//return data
	case 401:
		fmt.Println("Box Token is invalid")
	}
	/*body, err := token.get(urlStr)
	if err != nil {
		log.Fatal(err)
	}*/

	return data, resp.StatusCode
}

//GetEnterpriseUsers - Return a struct of all box users
func (token *RefreshToken) GetEnterpriseUsers() []User {
	enterprise, _ := token.GetUsers(0, 0)
	var boxusers []User
	offset := 0
	limit := 1000
	for int64(offset) < enterprise.TotalCount {
		tempusers, _ := token.GetUsers(limit, offset)
		for _, item := range tempusers.Entries {
			boxusers = append(boxusers, item)
		}
		//boxusers = append(boxusers, tempusers.Entries)
		offset += limit
	}
	return boxusers
}

//UpdateUserSpace - Update the UserSpace for a box account. Required userid int, and newsize int.
func (token RefreshToken) UpdateUserSpace(userid string, newsize int) int {
	boxURL := boxAPIURL + "/2.0/users/" + userid
	//log.Printf("Box Update User Space URL: %s and New Size: %s", boxURL, strconv.Itoa(newsize))

	query := fmt.Sprintf("{\"space_amount\": %d}", newsize)

	data := []byte(query)
	client := &http.Client{}
	r, err := http.NewRequest("PUT", boxURL, bytes.NewBuffer(data)) // <-- URL-encoded payload
	r.Header.Add("Authorization", "Bearer "+token.AccessToken)
	r.Header.Add("Content-Length", strconv.Itoa(len(data)))

	//useful in debuging the http put request and responses
	if err != nil {
		debug(httputil.DumpRequestOut(r, true))

	}
	resp, err := client.Do(r)
	if err != nil {
		defer resp.Body.Close()
		debug(httputil.DumpResponse(resp, true))
		fmt.Printf("Changed %s storage size to %d\n", userid, newsize)
		log.Printf("Changed %s storage size to %d\n", userid, newsize)
	} //else {
	//debug(httputil.DumpResponse(resp, true))
	//}

	//fmt.Println(resp.Status)
	return resp.StatusCode
}

//ChangeUserStatus - Change status on box account for user
func (token *RefreshToken) ChangeUserStatus(userid string, newstatus string) int {
	boxURL := boxAPIURL + "/2.0/users/" + userid

	query := fmt.Sprintf("{\"status\": \"%s\"}", newstatus)

	data := []byte(query)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", boxURL, bytes.NewBuffer(data))
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	req.Header.Add("Content-Lenght", strconv.Itoa(len(data)))

	if err != nil {
		debug(httputil.DumpRequestOut(req, true))
	}

	resp, err := client.Do(req)
	if err != nil {
		defer resp.Body.Close()
		debug(httputil.DumpResponse(resp, true))

		//log.Fatalln(err)
	}
	switch resp.StatusCode {
	case 200:
		fmt.Printf("Changed user %s to %s\n", userid, newstatus)
		log.Printf("Changed user %s to %s\n", userid, newstatus)
	default:
		fmt.Printf("ERROR: Unable to change user %s to %s\n", userid, newstatus)
	}
	return resp.StatusCode
}

//RollOffAccount - Roll Off Account to a personal account
func (token RefreshToken) RollOffAccount(userid string) int {
	boxURL := boxAPIURL + "/2.0/users/" + userid

	payload := strings.NewReader("{\"enterprise\": null}")

	req, err := http.NewRequest("PUT", boxURL, payload)

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	if err != nil {
		debug(httputil.DumpRequestOut(req, true))
	}
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		debug(httputil.DumpResponse(res, true))
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

	return res.StatusCode

}

//GetGroups - Retuns a list of current groups from box
func (token RefreshToken) GetGroups(limit, offset int) (GroupResponse, int) {
	url := fmt.Sprintf("%s/2.0/groups?limit=%d&offset=%d&fields=external_sync_identifier,provenance,description,name,created_at,modified_at", boxAPIURL, limit, offset)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	if err != nil {
		debug(httputil.DumpRequestOut(req, true))
	}
	resp, err := client.Do(req)
	if err != nil {
		debug(httputil.DumpResponse(resp, true))
	}

	var data GroupResponse
	switch resp.StatusCode {
	case 200:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Fatalln(err)
			switch v := err.(type) {
			case *json.SyntaxError:
				fmt.Println(string(body[v.Offset-40 : v.Offset]))
			}
		}
	case 401:
		fmt.Println("Box Token is invalid")
	}
	return data, resp.StatusCode

}

//DeleteGroup - delete box group
func (token RefreshToken) DeleteGroup(groupid string) int {
	url := fmt.Sprintf("%s/2.0/groups/%s", boxAPIURL, groupid)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)
	if err != nil {
		debug(httputil.DumpRequestOut(req, true))
	}
	resp, err := client.Do(req)
	if err != nil {
		debug(httputil.DumpResponse(resp, true))
	}
	return resp.StatusCode
}

//CreateGroup - create box group
func (token RefreshToken) CreateGroup(name string) {
	urlStr := fmt.Sprintf("%s/2.0/groups", boxAPIURL)
	client := &http.Client{}
	json := fmt.Sprintf("{\"name\": \"%s\"}", name)
	payload := strings.NewReader(json)
	req, err := http.NewRequest("POST", urlStr, payload)
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	if err != nil {
		debug(httputil.DumpRequestOut(req, true))
	}
	resp, err := client.Do(req)
	if err != nil {
		debug(httputil.DumpResponse(resp, true))
	}

	if resp.StatusCode == 201 {
		fmt.Println("Success: Created group ", name)
	}

	fmt.Println("response status: ", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response body: ", string(body))
}

//CreateGroupFull - create box group with all attributes
func (token RefreshToken) CreateGroupFull(name string, provenance string, externalsyncidentifier string) {
	urlStr := fmt.Sprintf("%s/2.0/groups", boxAPIURL)
	client := &http.Client{}
	json := fmt.Sprintf("{\"name\": \"%s\", \"provenance\": \"%s\", \"external_sync_identifier\": \"%s\"}", name, provenance, externalsyncidentifier)
	payload := strings.NewReader(json)
	req, err := http.NewRequest("POST", urlStr, payload)
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	if err != nil {
		debug(httputil.DumpRequestOut(req, true))
	}
	resp, err := client.Do(req)
	if err != nil {
		debug(httputil.DumpResponse(resp, true))
	}

	if resp.StatusCode == 201 {
		fmt.Println("Success: Created group: ", name)
	}

	fmt.Println("response status: ", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response body: ", string(body))
}
func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}
