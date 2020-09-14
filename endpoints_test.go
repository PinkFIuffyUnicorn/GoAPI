package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetUser(t *testing.T) {
	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	getUser(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v weant %v", status, http.StatusOK)
	}
}

func TestAddUser(t *testing.T) {
	var addResponse addResponse
	var jsonStr = []byte(`{"Email":"janez.novak@gmail.com","Password":"janez123","Name":"skupina1"}`)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	addUser(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v weant %v", status, http.StatusOK)
	}

	reqBody, _ := ioutil.ReadAll(rr.Body)
	reqBody = []byte(strings.ToLower(string(reqBody)))
	json.Unmarshal(reqBody, &addResponse)

	if addResponse.InsertedID == "" {
		t.Error("Handler returned empty response")
	}

}

// TODO
func TestUpdateUser(t *testing.T) {
	var updateResponse updateResponse
	var jsonStr = []byte(`{"Email":"janez.novak@gmail.comdd","Name":"skupina1222"}`)

	req, err := http.NewRequest("PUT", "/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "5f5e5c9f1265d2ffb73f33af"})
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	updateUser(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	reqBody, _ := ioutil.ReadAll(rr.Body)
	json.Unmarshal(reqBody, &updateResponse)

	if updateResponse.FieldsUpdated != 0 && updateResponse.ID != "5f5e5c9f1265d2ffb73f33af" {
		t.Log(updateResponse.FieldsUpdated, updateResponse.ID)
		t.Error("Handler returned unexpected body")
	}
}

func TestDeleteUser(t *testing.T) {
	var deletedResponse deletedResponse
	req, err := http.NewRequest("DELETE", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "5f5e5c292d07120dc8f25894"})
	rr := httptest.NewRecorder()
	deleteUser(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	reqBody, _ := ioutil.ReadAll(rr.Body)
	json.Unmarshal(reqBody, &deletedResponse)

	if deletedResponse.deletedCount != 0 || deletedResponse.ID != "5f5e5c292d07120dc8f25894" {
		t.Errorf("Handler returned unexpected body:\nExpected 0 (deletedCount) got %v\nExpected 5f5e5c292d07120dc8f25894 (ID) got %v", deletedResponse.deletedCount, deletedResponse.ID)
	}
}
