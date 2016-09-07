package tureng

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Result struct {
	Exception    string `json:"ExceptionMessage"`
	IsSuccessful bool   `json:"IsSuccessful"`
	MobileResult struct {
		IsFound  int `json:"IsFound"`
		IsTRToEN int `json:"IsTRToEN"`
		Results  []struct {
			Category string `json:"CategoryEN"`
			Term     string `json:"Term"`
			TypeEN   string `json:"TypeEN"`
		} `json:"Results"`
	} `json:"MobileResult"`
}

type Req struct {
	Term string `json:"Term"`
	Code string `json:"Code"`
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func PrepareReq(word string) (*bytes.Buffer, error) {
	code := getMD5Hash(fmt.Sprintf("%s%s", word, WTF))
	req := Req{word, code}
	jsonStr, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(jsonStr), nil
}

func Translate(reqString *bytes.Buffer) (*Result, error) {
	resp, err := http.Post(URL, BODY_TYPE, reqString)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	rr := &Result{}
	err = json.Unmarshal(body, rr)
	if err != nil {
		return nil, err
	}

	if !rr.IsSuccessful {
		return nil, errors.New("Tureng request is not successfull!")
	}

	if rr.MobileResult.IsFound != 1 {
		return nil, errors.New("No results!")
	}

	if !rr.IsSuccessful || rr.MobileResult.IsFound != 1 {
		return nil, errors.New("Not succesfull")
	}
	return rr, nil
}

const (
	URL       = "http://ws.tureng.com/TurengSearchServiceV4.svc/Search"
	WTF       = "46E59BAC-E593-4F4F-A4DB-960857086F9C"
	BODY_TYPE = "application/json"
)
