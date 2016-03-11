package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

type GoDynErr struct {
	msg   string
	inner error
}

func (e *GoDynErr) Error() string {
	if e.inner == nil {
		return e.msg
	}

	return fmt.Sprintf("%v: %v", e.msg, e.inner.Error())
}

func main() {
	dns_service := os.Getenv("GODYN_SERVICE")
	dns_host := os.Getenv("GODYN_HOST")
	dns_user := os.Getenv("GODYN_USER")
	dns_pw := os.Getenv("GODYN_PW")
	if dns_service == "" || dns_host == "" || dns_user == "" || dns_pw == "" {
		log.Printf("Please ensure GODYN_SERVICE, GODYN_HOST, GODYN_USER, and GODYN_PW are set")
		return
	}

	log.Printf("Using %v service", dns_service)

	ip, err := GetExternalIP()
	if err != nil {
		log.Printf("Unable to get external IP: %v", err)
		return
	}

	log.Printf("External IP is %v", ip)

	reqstr := fmt.Sprintf("https://%v/nic/update?hostname=%v&amp;myip=%v", dns_service, dns_host, ip)

	log.Printf("Submitting update request to %v", reqstr)
	err = SubmitAuthenticatedRequest(reqstr, dns_user, dns_pw)
	if err != nil {
		log.Printf("Failed to update: %v", err)
		return
	}

	log.Printf("IP updated successfully")
}

var IPRegex = regexp.MustCompile(`(\d{1,3}\.){3}\d{1,3}`)

/*
 * GetExternalIP returns the current public IP of this host
 */
func GetExternalIP() (string, error) {
	resp, err := http.Get("http://checkip.dyndns.com/")
	if err != nil {
		return "", &GoDynErr{msg: "Unable to get IP response", inner: err}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", &GoDynErr{msg: "Unable to read IP response", inner: err}
	}

	log.Printf("%v", string(body))

	ip := IPRegex.FindString(string(body))
	if ip == "" {
		return "", &GoDynErr{msg: "IP not found in checker response"}
	}

	return ip, nil
}

func SubmitAuthenticatedRequest(url, user, pw string) error {
	client := &http.Client{
		CheckRedirect: nil,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &GoDynErr{msg: "Unable to create update request", inner: err}
	}

	req.SetBasicAuth(user, pw)

	resp, err := client.Do(req)
	if err != nil {
		return &GoDynErr{msg: "Unable to execute update request", inner: err}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &GoDynErr{msg: "Unable to read update response", inner: err}
	}

	log.Printf("Response from service: %v", string(body))

	return nil
}
