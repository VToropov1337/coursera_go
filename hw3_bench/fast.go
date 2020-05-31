package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type User struct {
	Browsers []string `json:"browsers"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
}

type Users []User

func InArray(a []string, e string) bool {
	for _, x := range a {
		if x == e {
			return true
		}
	}
	return false
}

// вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	seenBrowsers := []string{}
	uniqueBrowsers := 0
	foundUsers := ""

	lines := strings.Split(string(fileContents), "\n")
	var users Users
	for _, line := range lines {
		var user User
		err := json.Unmarshal([]byte(line), &user)
		if err != nil {
			panic(err)
		}
		users = append(users, user)
	}

	for i, v := range users {
		isAndroid := false
		isMSIE := false
		for _, browser := range v.Browsers {
			if ok, err := regexp.MatchString("Android", browser); ok && err == nil {
				isAndroid = true
				if InArray(seenBrowsers, browser) {
					continue
				} else {
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}

			}
			if ok, err := regexp.MatchString("MSIE", browser); ok && err == nil {
				isMSIE = true
				if InArray(seenBrowsers, browser) {
					continue
				} else {
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}

			}

		}
		if isMSIE && isAndroid {
			v.Email = strings.Replace(v.Email, "@", " [at] ", 1)
			foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, v.Name, v.Email)

		}
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))

}
