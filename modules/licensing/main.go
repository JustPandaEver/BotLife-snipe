package licensing

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func errorHandle(err error) {
	fmt.Println("[ERROR]An error has occurred. Please contact your seller.")
	os.Exit(0)
}

func CheckFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func CheckPresaleBotLicense(api string, ssl bool, silent bool) {
	if !CheckFileExist("license.dat") {
		if !silent {
			fmt.Println("license.dat not found.")
		}
		NewLicense()
		fmt.Println("Please run the program again.")
		os.Exit(0)
	}

	li, err := ioutil.ReadFile("license.dat") //string(li)
	if err != nil {
		if !silent {
			errorHandle(err)
		}
	}

	if ssl {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		data := url.Values{}
		data.Set("license", string(li))
		u, _ := url.ParseRequestURI(api + "/presalebot")
		urlStr := fmt.Sprintf("%v", u)
		r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(r)
		if err != nil {
			if !silent {
				fmt.Println("Unable to connect to license server.")
			}
			os.Exit(0)
		}
		defer resp.Body.Close()
		resp_body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			if string(resp_body) != "Good" {
				if string(resp_body) == "Expired" {
					if !silent {
						fmt.Println("License is Expired.")
					}
					os.Exit(0)
				} else {
					if !silent {
						fmt.Println("Connot verify license, Please contact your seller.")
					}
					os.Exit(0)
				}
			}
		}
	} else {
		client := &http.Client{}
		data := url.Values{}
		data.Set("license", string(li))
		u, _ := url.ParseRequestURI(api + "/presalebot")
		urlStr := fmt.Sprintf("%v", u)
		r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(r)
		if err != nil {
			if !silent {
				fmt.Println("Unable to connect to license server.")
			}
			os.Exit(0)
		}
		defer resp.Body.Close()
		resp_body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			if string(resp_body) != "Good" {
				if string(resp_body) == "Expired" {
					if !silent {
						fmt.Println("License is Expired.")
					}
					os.Exit(0)
				} else {
					if !silent {
						fmt.Println("Connot verify license, Please contact your seller.")
					}
					os.Exit(0)
				}
			}
		}
	}
}

func CheckTelegramBotLicense(api string, ssl bool, silent bool) {
	if !CheckFileExist("license.dat") {
		if !silent {
			fmt.Println("license.dat not found.")
		}
		NewLicense()
		fmt.Println("Please run the program again.")
		os.Exit(0)
	}

	li, err := ioutil.ReadFile("license.dat") //string(li)
	if err != nil {
		if !silent {
			errorHandle(err)
		}
	}

	if ssl {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		data := url.Values{}
		data.Set("license", string(li))
		u, _ := url.ParseRequestURI(api + "/telegrambot")
		urlStr := fmt.Sprintf("%v", u)
		r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(r)
		if err != nil {
			if !silent {
				fmt.Println("Unable to connect to license server.")
			}
			os.Exit(0)
		}
		defer resp.Body.Close()
		resp_body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			if string(resp_body) != "Good" {
				if string(resp_body) == "Expired" {
					if !silent {
						fmt.Println("License is Expired.")
					}
					os.Exit(0)
				} else {
					if !silent {
						fmt.Println("Connot verify license, Please contact your seller.")
					}
					os.Exit(0)
				}
			}
		}
	} else {
		client := &http.Client{}
		data := url.Values{}
		data.Set("license", string(li))
		u, _ := url.ParseRequestURI(api + "/telegrambot")
		urlStr := fmt.Sprintf("%v", u)
		r, _ := http.NewRequest("POST", urlStr, bytes.NewBufferString(data.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		resp, err := client.Do(r)
		if err != nil {
			if !silent {
				fmt.Println("Unable to connect to license server.")
			}
			os.Exit(0)
		}
		defer resp.Body.Close()
		resp_body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode == 200 {
			if string(resp_body) != "Good" {
				if string(resp_body) == "Expired" {
					if !silent {
						fmt.Println("License is Expired.")
					}
					os.Exit(0)
				} else {
					if !silent {
						fmt.Println("Connot verify license, Please contact your seller.")
					}
					os.Exit(0)
				}
			}
		}
	}
}

func NewLicense() {
	fmt.Println("Please enter your license: ")
	var license string

	fmt.Scanln(&license)
	license = strings.TrimRight(license, "\r\n")

	err := os.WriteFile("license.dat", []byte(license), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
