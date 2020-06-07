package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	jsoniterator "github.com/json-iterator/go"
)

type User struct {
	Browsers []string `json:"browsers"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
}

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
	defer file.Close()

	// seenBrowsers := []string{}
	seenBrowsers := make([]string, 0, 1)
	uniqueBrowsers := 0
	foundUsers := ""

	var json = jsoniterator.ConfigCompatibleWithStandardLibrary
	var user User
	scanner := bufio.NewScanner(file)
	var i int
	for scanner.Scan() {

		err := json.Unmarshal([]byte(scanner.Text()), &user)
		if err != nil {
			panic(err)
		}
		isAndroid := false
		isMSIE := false
		for _, browser := range user.Browsers {
			if ok := strings.Contains(browser, "Android"); ok {
				isAndroid = true
				if InArray(seenBrowsers, browser) {
					continue
				} else {
					seenBrowsers = append(seenBrowsers, browser)
					uniqueBrowsers++
				}

			}
			if ok := strings.Contains(browser, "MSIE"); ok {
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
			user.Email = strings.Replace(user.Email, "@", " [at] ", 1)
			foundUsers += fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, user.Email)

		}
		i++

	}

	fmt.Fprintln(out, "found users:\n"+foundUsers)
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))

}
