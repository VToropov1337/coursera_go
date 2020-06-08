package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

type Row struct {
	Name      string
	Id        int    `xml:"id"`
	Age       int    `xml:"age"`
	Gender    string `xml:"gender"`
	About     string `xml:"about"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:last_name`
}

type Users struct {
	Users []Row `xml:"row"`
}

type TestCase struct {
	SearchRequest SearchRequest
	StatusCode    int
}

var users Users

func SearchServer(w http.ResponseWriter, r *http.Request) {

	parseDataSet()
	query_limit := r.FormValue("limit")
	limit, err := strconv.Atoi(query_limit)
	if err != nil {
		panic(err)
	}
	query := r.FormValue("query")

	if query == "list" {
		users, _ := json.Marshal(users.Users[:limit])
		w.WriteHeader(http.StatusOK)
		w.Write(users)
	}

}

func parseDataSet() {
	xmlFile, err := os.Open("dataset.xml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully opened xml file")
	defer xmlFile.Close()

	bytes, _ := ioutil.ReadAll(xmlFile)

	xml.Unmarshal(bytes, &users)
}

func TestFindUsersHeaders(t *testing.T) {
	cases := []TestCase{
		TestCase{
			SearchRequest: SearchRequest{
				Limit: 20,
				Query: "list",
			},
			StatusCode: 200,
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	c := &SearchClient{URL: ts.URL}
	for _, v := range cases {
		c.FindUsers(v.SearchRequest)

	}
}
