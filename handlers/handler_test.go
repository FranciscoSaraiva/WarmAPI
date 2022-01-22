package handlers

import (
	"net/http/httptest"
	"strings"
	"testing"

	"WarmAPI/config"
	"WarmAPI/models"

	"github.com/julienschmidt/httprouter"
)

func TestRetrieveCandidates(t *testing.T) {
	clearDB()
	c1 := models.ChallengerCandidate{ID: 1, FirstName: "Francisco", LastName: "Saraiva", EmailAddress: "francisco.saraiva@olx.com"}
	c2 := models.ChallengerCandidate{ID: 2, FirstName: "Bruno", LastName: "Almeida", EmailAddress: "bruno.almeida@olx.com"}
	c3 := models.ChallengerCandidate{ID: 3, FirstName: "Daniel", LastName: "Matos", EmailAddress: "daniel.matos@olx.com"}
	c4 := models.ChallengerCandidate{ID: 4, FirstName: "Artur", LastName: "Esteves", EmailAddress: "artur.esteves@olx.com"}
	c5 := models.ChallengerCandidate{ID: 5, FirstName: "Francisco", LastName: "Silva", EmailAddress: "francisco.silva@olx.com"}
	candidates := [5]models.ChallengerCandidate{c1, c2, c3, c4, c5}

	for _, v := range candidates {
		models.InsertCandidate(&v)
	}

	w := httptest.NewRecorder()
	RetrieveCandidates(w, nil, nil)
	response := w.Body.String()

	//fmt.Printf("%s", w.Body.String())

	if len(response) == 0 {
		t.Errorf("Expected the request body to have content but came empty")
	}
	if !strings.Contains(response, "\"message\":\"Retrieve Candidates\"") {
		t.Errorf("Expected the message to be \"Retrieve Candidates\" but it was different")
	}
	if !strings.Contains(response, "\"status\":200") {
		t.Errorf("Expected the status to be 200 but it was different")
	}
	if !strings.Contains(response, "Francisco") ||
		!strings.Contains(response, "Bruno") ||
		!strings.Contains(response, "Daniel") ||
		!strings.Contains(response, "Artur") {
		t.Errorf("Expected the names of the candidates to be in the response but they were not all in.")
	}
	if !strings.Contains(response, "Saraiva") ||
		!strings.Contains(response, "Almeida") ||
		!strings.Contains(response, "Matos") ||
		!strings.Contains(response, "Esteves") {
		t.Errorf("Expected the last names of the candidates to be in the response but they were not all in.")
	}
	if !strings.Contains(response, "francisco.saraiva@olx.com") ||
		!strings.Contains(response, "bruno.almeida@olx.com") ||
		!strings.Contains(response, "daniel.matos@olx.com") ||
		!strings.Contains(response, "artur.esteves@olx.com") {
		t.Errorf("Expected the emails of the candidates to be in the response but they were not all in.")
	}
}

func TestRetrieveCandidatesNotFound(t *testing.T) {
	clearDB()

	w := httptest.NewRecorder()
	RetrieveCandidates(w, nil, nil)
	response := w.Body.String()

	//fmt.Printf("%s", w.Body.String())

	if !strings.Contains(response, "\"message\":\"No candidates were found.\"") {
		t.Errorf("Expected message to be No candidates were found. but it was different")
	}
	if !strings.Contains(response, "\"status\":404") {
		t.Errorf("Expected status to be 404 but it was different")
	}
}

func TestRetrieveCandidateID(t *testing.T) {
	clearDB()
	c1 := models.ChallengerCandidate{ID: 1, FirstName: "Francisco", LastName: "Saraiva", EmailAddress: "francisco.saraiva@olx.com"}

	models.InsertCandidate(&c1)

	w := httptest.NewRecorder()
	p := httprouter.Params{httprouter.Param{Key: "id", Value: "1"}}

	RetrieveCandidateID(w, nil, p)
	response := w.Body.String()

	//fmt.Printf("%s", w.Body.String())

	if len(response) == 0 {
		t.Errorf("Expected the request body to have content but came empty")
	}
	if !strings.Contains(response, "\"message\":\"Candidate Retrieved\"") {
		t.Errorf("Expected the message to be \"Candidate Retrieved\" but it was different")
	}
	if !strings.Contains(response, "\"status\":200") {
		t.Errorf("Expected the status to be 200 but it was different")
	}
	if !strings.Contains(response, "\"first_name\":\"Francisco\"") {
		t.Errorf("Expected the names of the candidate to be in the response but it was not.")
	}
	if !strings.Contains(response, "\"last_name\":\"Saraiva\"") {
		t.Errorf("Expected the last name of the candidate to be in the response but it was not.")
	}
	if !strings.Contains(response, "\"email_address\":\"francisco.saraiva@olx.com\"") {
		t.Errorf("Expected the email of the candidate to be in the response but it was not.")
	}
}

func TestRetrieveCandidateIDErrorOnlyNumbers(t *testing.T) {
	clearDB()
	c1 := models.ChallengerCandidate{ID: 1, FirstName: "Francisco", LastName: "Saraiva", EmailAddress: "francisco.saraiva@olx.com"}

	models.InsertCandidate(&c1)

	w := httptest.NewRecorder()
	p := httprouter.Params{httprouter.Param{Key: "id", Value: "1a"}}

	RetrieveCandidateID(w, nil, p)
	response := w.Body.String()

	//fmt.Printf("%s", w.Body.String())

	if len(response) == 0 {
		t.Errorf("Expected the request body to have content but came empty")
	}
	if !strings.Contains(response, "\"message\":\"The id must be all numbers\"") {
		t.Errorf("Expected the message to be \"The id must be all numbers\" but it was different")
	}
	if !strings.Contains(response, "\"status\":400") {
		t.Errorf("Expected the status to be 400 but it was different")
	}
}

func TestRetrieveCandidateIDErrorNotFound(t *testing.T) {
	clearDB()
	w := httptest.NewRecorder()
	p := httprouter.Params{httprouter.Param{Key: "id", Value: "1"}}

	RetrieveCandidateID(w, nil, p)
	response := w.Body.String()

	//fmt.Printf("%s", w.Body.String())

	if len(response) == 0 {
		t.Errorf("Expected the request body to have content but came empty")
	}
	if !strings.Contains(response, "\"message\":\"Candidate not found\"") {
		t.Errorf("Expected the message to be \"Candidate not found\" but it was different")
	}
	if !strings.Contains(response, "\"status\":404") {
		t.Errorf("Expected the status to be 404 but it was different")
	}
}

func TestRetrieveCandidateFromDB(t *testing.T) {
	clearDB()
	c1 := models.ChallengerCandidate{ID: 1, FirstName: "Francisco", LastName: "Saraiva", EmailAddress: "francisco.saraiva@olx.com"}
	models.InsertCandidate(&c1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/candidate", strings.NewReader("{\"email\":\"francisco.saraiva@olx.com\"}"))
	RetrieveCandidate(w, r, nil)
	response := w.Body.String()

	//fmt.Printf("%s", w.Body.String())

	if len(response) == 0 {
		t.Errorf("Expected the request body to have content but came empty")
	}
	if !strings.Contains(response, "\"message\":\"Retrieved candidate\"") {
		t.Errorf("Expected the message to be \"Retrieved candidate \" but it was different")
	}
	if !strings.Contains(response, "\"status\":200") {
		t.Errorf("Expected the status to be 200 but it was different")
	}
	if !strings.Contains(response, "\"first_name\":\"Francisco\"") {
		t.Errorf("Expected the names of the candidate to be in the response but it was not.")
	}
	if !strings.Contains(response, "\"last_name\":\"Saraiva\"") {
		t.Errorf("Expected the last name of the candidate to be in the response but it was not.")
	}
	if !strings.Contains(response, "\"email_address\":\"francisco.saraiva@olx.com\"") {
		t.Errorf("Expected the email of the candidate to be in the response but it was not.")
	}
}

func TestRetrieveCandidateFromAPI(t *testing.T) {
	clearDB()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/candidate", strings.NewReader("{\"email\":\"riblesfranco@gmail.com\"}"))
	RetrieveCandidate(w, r, nil)
	response := w.Body.String()

	//fmt.Printf("%s", w.Body.String())

	if len(response) == 0 {
		t.Errorf("Expected the request body to have content but came empty")
	}
	if !strings.Contains(response, "\"message\":\"Retrieved candidate\"") {
		t.Errorf("Expected the message to be \"Retrieved candidate \" but it was different")
	}
	if !strings.Contains(response, "\"status\":200") {
		t.Errorf("Expected the status to be 200 but it was different")
	}
	if !strings.Contains(response, "\"first_name\":\"Franco\"") {
		t.Errorf("Expected the names of the candidate to be in the response but it was not.")
	}
	if !strings.Contains(response, "\"last_name\":\"Ribles\"") {
		t.Errorf("Expected the last name of the candidate to be in the response but it was not.")
	}
	if !strings.Contains(response, "\"email_address\":\"riblesfranco@gmail.com\"") {
		t.Errorf("Expected the email of the candidate to be in the response but it was not.")
	}
}

func TestRetrieveCandidateFromAPINotFound(t *testing.T) {
	clearDB()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/candidate", strings.NewReader("{\"email\":\"email\"}"))
	RetrieveCandidate(w, r, nil)
	response := w.Body.String()

	//fmt.Printf("%s", w.Body.String())

	if len(response) == 0 {
		t.Errorf("Expected the request body to have content but came empty")
	}
	if !strings.Contains(response, "\"message\":\"No candidate was found within the Greenhouse database.\"") {
		t.Errorf("Expected the message to be \"No candidate was found within the Greenhouse database\" but it was different")
	}
	if !strings.Contains(response, "\"status\":404") {
		t.Errorf("Expected the status to be 404 but it was different")
	}
}

func TestRetrieveCandidateErrJsonBodyEmpty(t *testing.T) {
	clearDB()
	c1 := models.ChallengerCandidate{ID: 1, FirstName: "Francisco", LastName: "Saraiva", EmailAddress: "francisco.saraiva@olx.com"}
	models.InsertCandidate(&c1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/candidate", strings.NewReader(""))
	RetrieveCandidate(w, r, nil)
	response := w.Body.String()

	//fmt.Printf("%s", w.Body.String())

	if len(response) == 0 {
		t.Errorf("Expected the request body to have content but came empty")
	}
	if !strings.Contains(response, "\"message\":\"The Json body is empty\"") {
		t.Errorf("Expected the message to be \"The Json body is empty\" but it was different")
	}
	if !strings.Contains(response, "\"status\":400") {
		t.Errorf("Expected the status to be 400 but it was different")
	}
}

func TestRetrieveCandidateErrJsonObjectEmpty(t *testing.T) {
	clearDB()
	c1 := models.ChallengerCandidate{ID: 1, FirstName: "Francisco", LastName: "Saraiva", EmailAddress: "francisco.saraiva@olx.com"}
	models.InsertCandidate(&c1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/candidate", strings.NewReader("{}"))
	RetrieveCandidate(w, r, nil)
	response := w.Body.String()

	//fmt.Printf("%s", w.Body.String())

	if len(response) == 0 {
		t.Errorf("Expected the request body to have content but came empty")
	}
	if !strings.Contains(response, "\"message\":\"The Json object is empty\"") {
		t.Errorf("Expected the message to be \"The Json body is empty\" but it was different")
	}
	if !strings.Contains(response, "\"status\":400") {
		t.Errorf("Expected the status to be 400 but it was different")
	}
}

func TestRetrieveCandidateErrEmailIsEmpty(t *testing.T) {
	clearDB()
	c1 := models.ChallengerCandidate{ID: 1, FirstName: "Francisco", LastName: "Saraiva", EmailAddress: "francisco.saraiva@olx.com"}
	models.InsertCandidate(&c1)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/candidate", strings.NewReader("{\"email\":\"\"}"))
	RetrieveCandidate(w, r, nil)
	response := w.Body.String()

	//fmt.Printf("%s", w.Body.String())

	if len(response) == 0 {
		t.Errorf("Expected the request body to have content but came empty")
	}
	if !strings.Contains(response, "\"message\":\"The email of the candidate is empty\"") {
		t.Errorf("Expected the message to be \"The email of the candidate is empty\" but it was different")
	}
	if !strings.Contains(response, "\"status\":400") {
		t.Errorf("Expected the status to be 400 but it was different")
	}
}

func clearDB() {
	collection := config.GetCollection("candidates")
	collection.DropCollection()
}
