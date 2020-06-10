package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

type Row struct {
	Name      string
	Id        int    `xml:"id"`
	Age       int    `xml:"age"`
	Gender    string `xml:"gender"`
	About     string `xml:"about"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
}

type ListRow struct {
	Rows []Row `xml:"row"`
}

var list ListRow
var users []User

func SearchServer(w http.ResponseWriter, r *http.Request) {

	query_limit := r.FormValue("limit")
	limit, err := strconv.Atoi(query_limit)
	if err != nil {
		panic(err)
	}
	query := r.FormValue("query")

	if query == "__timeout" {
		time.Sleep(time.Second * 2)
		w.WriteHeader(http.StatusRequestTimeout)
		return
	}

	if query == "__bad_token" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if query == "__fatal_error" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if query == "__bad_request" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	parseDataSet()

	if query == "users" {
		limit_users, _ := json.Marshal(users[:limit])
		w.WriteHeader(http.StatusOK)
		w.Write(limit_users)
	}

}

func parseDataSet() {
	xmlFile, err := os.Open("dataset.xml")
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()

	bytes, _ := ioutil.ReadAll(xmlFile)

	xml.Unmarshal(bytes, &list)

	for _, v := range list.Rows {
		users = append(users, User{
			Id:     v.Id,
			Name:   v.FirstName + " " + v.LastName,
			Age:    v.Age,
			About:  v.About,
			Gender: v.Gender,
		})

	}

}

func TestFindUsersTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	client := &SearchClient{
		URL: ts.URL,
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 20, Query: "__timeout"})

	if resp != nil {
		t.Errorf("Expected nil, got response: %v with err: %v", resp, err)
	}

	if err == nil {
		t.Errorf("Expected err, got nil : %v ", err)
	}
}

func TestFindUsersBadRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	client := &SearchClient{
		URL: ts.URL,
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 20, Query: "__bad_request"})

	if resp != nil {
		t.Errorf("Expected nil, got response: %v with err: %v", resp, err)
	}
}

func TestFindUsersFatalError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	client := &SearchClient{
		URL: ts.URL,
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 20, Query: "__fatal_error"})
	if resp != nil {
		t.Errorf("Expected nil, got response: %v with err: %v", resp, err)
	}
}

func TestFindUsersBadToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	client := &SearchClient{
		URL: ts.URL,
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 20, Query: "__bad_token"})
	if resp != nil {
		t.Errorf("Expected nil, got response: %v with err: %v", resp, err)
	}
}

func TestFindUsersBadOrderField(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	client := &SearchClient{
		URL: ts.URL,
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 20, OrderField: "ErrorBadOrderField"})
	if resp != nil {
		t.Errorf("Expected nil, got response: %v with err: %v", resp, err)
	}
}

func TestFindUsersBadUnknownError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	client := &SearchClient{
		URL: "/",
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 20})
	if resp != nil {
		t.Errorf("Expected nil, got response: %v with err: %v", resp, err)
	}
}

func TestFindUsersGoodCase(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	client := &SearchClient{
		URL: ts.URL,
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 2, Query: "users"})
	if err != nil {
		t.Errorf("Error: %v. ", err)
	}
	if resp == nil {
		t.Errorf("Expected resp, got: %v. ", resp)
	}

	if len(resp.Users) != 2 {
		t.Errorf("Expected 2 users, got: %v. ", len(resp.Users))
	}
}

func TestFindUsersNegativeLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	client := &SearchClient{
		URL: ts.URL,
	}

	resp, err := client.FindUsers(SearchRequest{Limit: -1, Query: "users"})
	if resp != nil {
		t.Errorf("Expected nil, got: %v. ", resp)
	}
	if err == nil {
		t.Errorf("Expected error, got: %v. ", err)
	}
}

func TestFindUsersOverLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	client := &SearchClient{
		URL: ts.URL,
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 26, Query: "users"})
	if err != nil {
		t.Errorf("Expected nil, got: %v. ", err)
	}

	if len(resp.Users) != 25 {
		t.Errorf("Expected 25, got: %v. ", len(resp.Users))
	}
}

func TestFindUsersLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	client := &SearchClient{
		URL: ts.URL,
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 26, Query: "users"})
	if err != nil {
		t.Errorf("Expected nil, got: %v. ", err)
	}

	if len(resp.Users) != 25 {
		t.Errorf("Expected 25, got: %v. ", len(resp.Users))
	}
}

func TestFindUsersNegativeOffset(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	client := &SearchClient{
		URL: ts.URL,
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 100, Offset: -1, Query: "users"})
	if err == nil {
		t.Errorf("Expected nil, got: %v. ", err)
	}
	if resp != nil {
		t.Errorf("Expected nil, got: %v. ", err)
	}
}

func TestFindUsersZeroLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	defer ts.Close()

	client := &SearchClient{
		URL: ts.URL,
	}

	resp, err := client.FindUsers(SearchRequest{Limit: 0, Query: "users"})
	if resp.NextPage != true {
		t.Errorf("Expected true, got: %v. ", resp)
	}
	if err != nil {
		t.Errorf("Expected nil, got: %v. ", resp)
	}
}

