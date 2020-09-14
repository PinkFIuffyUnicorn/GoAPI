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

	if updateResponse.FieldsUpdated != 0 || updateResponse.ID != "5f5e5c9f1265d2ffb73f33af" || updateResponse.ID == "" {
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

func TestAddGroup(t *testing.T) {
	var addResponse addResponse
	var groupAlreadyAddedResponse groupAlreadyAddedResponse
	var jsonStr = []byte(`{"Name":"skupina1"}`)

	req, err := http.NewRequest("POST", "/groups", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	addGroup(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v weant %v", status, http.StatusOK)
	}

	temp := strings.Contains(rr.Body.String(), "InsertedID")

	if temp {
		reqBody, _ := ioutil.ReadAll(rr.Body)
		json.Unmarshal(reqBody, &addResponse)
		if addResponse.InsertedID == "" {
			t.Error("Handler returned empty response, expected InsertedID attribute")
		}
	} else {
		reqBody, _ := ioutil.ReadAll(rr.Body)
		json.Unmarshal(reqBody, &groupAlreadyAddedResponse)
		if groupAlreadyAddedResponse.Response != "Group Name already exists" {
			t.Errorf("Handler returned wrong responsem expected 'Group Name already exists' got %v", groupAlreadyAddedResponse.Response)
		}
	}
}

func TestGetGroup(t *testing.T) {
	req, err := http.NewRequest("GET", "/groups", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	getGroup(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v weant %v", status, http.StatusOK)
	}
}

func TestUpdateGroup(t *testing.T) {
	var updateResponse updateResponse
	var jsonStr = []byte(`{"Name":"skupina122322"}`)

	req, err := http.NewRequest("PUT", "/groups", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "5f5e066649137f8066b54529"})
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	updateGroup(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	reqBody, _ := ioutil.ReadAll(rr.Body)
	json.Unmarshal(reqBody, &updateResponse)

	if updateResponse.FieldsUpdated > 0 || updateResponse.ID != "5f5e066649137f8066b54529" || updateResponse.ID == "" {
		t.Errorf("Handler returned unexpected body:\nExpected 0 (FieldsUpdated) got %v\nExpected 5f5e066649137f8066b54529 (ID) got %v", updateResponse.FieldsUpdated, updateResponse.ID)
	}
}

func TestDeleteGroup(t *testing.T) {
	var deletedResponse deletedResponse
	req, err := http.NewRequest("DELETE", "/groups", nil)
	if err != nil {
		t.Fatal(err)
	}
	req = mux.SetURLVars(req, map[string]string{"id": "5f5e066649137f8066b54529"})
	rr := httptest.NewRecorder()
	deleteGroup(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	reqBody, _ := ioutil.ReadAll(rr.Body)
	json.Unmarshal(reqBody, &deletedResponse)

	if deletedResponse.deletedCount != 0 || deletedResponse.ID != "5f5e066649137f8066b54529" {
		t.Errorf("Handler returned unexpected body:\nExpected 0 (deletedCount) got %v\nExpected 5f5e066649137f8066b54529 (ID) got %v", deletedResponse.deletedCount, deletedResponse.ID)
	}
}
