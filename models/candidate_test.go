package models

import (
	"testing"

	"WarmAPI/config"
)

func TestGetAllCandidates(t *testing.T) {
	clearDB()
	c1 := ChallengerCandidate{ID: 1, FirstName: "Francisco", LastName: "Saraiva", EmailAddress: "francisco.saraiva@olx.com"}
	c2 := ChallengerCandidate{ID: 2, FirstName: "Bruno", LastName: "Almeida", EmailAddress: "bruno.almeida@olx.com"}
	c3 := ChallengerCandidate{ID: 3, FirstName: "Daniel", LastName: "Matos", EmailAddress: "daniel.matos@olx.com"}
	c4 := ChallengerCandidate{ID: 4, FirstName: "Artur", LastName: "Esteves", EmailAddress: "artur.esteves@olx.com"}
	c5 := ChallengerCandidate{ID: 5, FirstName: "Francisco", LastName: "Silva", EmailAddress: "francisco.silva@olx.com"}
	candidates := [5]ChallengerCandidate{c1, c2, c3, c4, c5}

	for _, v := range candidates {
		InsertCandidate(&v)
	}

	list, err := GetAllCandidates()

	if config.CheckError(err) {
		t.Errorf("Expected no error but got error [%s]", err)
	}
	if len(*list) != 5 {
		t.Errorf("Expected the size of the candidates to be 5 but got [%d]", len(*list))
	}

	count := 1
	for _, v := range *list {
		if v.ID != count {
			t.Errorf("Expected the ID of the candidate to be %d but it was %d", count, v.ID)
		}
		if count == 1 && v.FirstName != "Francisco" {
			t.Errorf("Expected the name of the first candidate to be Francisco but it was %s", v.FirstName)
		}
		if count == 2 && v.FirstName != "Bruno" {
			t.Errorf("Expected the name of the first candidate to be Bruno but it was %s", v.FirstName)
		}
		if count == 3 && v.FirstName != "Daniel" {
			t.Errorf("Expected the name of the first candidate to be Daniel but it was %s", v.FirstName)
		}
		if count == 4 && v.FirstName != "Artur" {
			t.Errorf("Expected the name of the first candidate to be Artur but it was %s", v.FirstName)
		}
		count++
	}

}

func TestGetCandidateByID(t *testing.T) {
	clearDB()
	c1 := ChallengerCandidate{ID: 1, FirstName: "Francisco", LastName: "Saraiva", EmailAddress: "francisco.saraiva@olx.com"}
	InsertCandidate(&c1)
	idWanted := 1

	candidate, err := GetCandidateByID(idWanted)

	if config.CheckError(err) {
		t.Errorf("Expected no error but got [%s]", err)
	}

	if candidate.ID != idWanted {
		t.Errorf("Expected the id to be %d but it was %d", idWanted, candidate.ID)
	}
}

func TestGetCandidateByIDNotFound(t *testing.T) {
	clearDB()
	c1 := ChallengerCandidate{ID: 1, FirstName: "Francisco", LastName: "Saraiva", EmailAddress: "francisco.saraiva@olx.com"}
	InsertCandidate(&c1)
	idWanted := 2

	candidate, err := GetCandidateByID(idWanted)

	if !config.CheckError(err) {
		t.Errorf("Expected an error but it came empty")
	}
	if candidate != nil {
		t.Errorf("Expected the candidate to be nil but it was created.")
	}
}

func TestGetCandidateByEmail(t *testing.T) {
	clearDB()
	c1 := ChallengerCandidate{ID: 1, FirstName: "Francisco", LastName: "Saraiva", EmailAddress: "francisco.saraiva@olx.com"}
	InsertCandidate(&c1)
	emailWanted := "francisco.saraiva@olx.com"

	candidate, err := GetCandidateByEmail("francisco.saraiva@olx.com")

	if config.CheckError(err) {
		t.Errorf("Expected no rror but got [%s]", err)
	}

	if candidate.EmailAddress != emailWanted {
		t.Errorf("Expected the email to be %s but it was %s", emailWanted, candidate.EmailAddress)
	}
}

func TestGetCandidateByEmailNotFound(t *testing.T) {
	clearDB()
	c1 := ChallengerCandidate{ID: 1, FirstName: "Francisco", LastName: "Saraiva", EmailAddress: "francisco.saraiva@olx.com"}
	InsertCandidate(&c1)
	emailWanted := "aaaa"

	candidate, err := GetCandidateByEmail(emailWanted)

	if !config.CheckError(err) {
		t.Errorf("Expected an error but it came empty")
	}
	if candidate != nil {
		t.Errorf("Expected the candidate to be nil but it was created.")
	}
}

func TestInsertCandidate(t *testing.T) {
	clearDB()
	c1 := ChallengerCandidate{ID: 1, FirstName: "Francisco", LastName: "Saraiva", EmailAddress: "francisco.saraiva@olx.com"}
	InsertCandidate(&c1)
	idWanted := 1
	firstNameWanted := "Francisco"
	lastNameWanted := "Saraiva"
	emailWanted := "francisco.saraiva@olx.com"

	candidate, err := GetCandidateByID(1)

	if config.CheckError(err) {
		t.Errorf("Expected no rror but got [%s]", err)
	}

	if candidate.ID != idWanted {
		t.Errorf("Expected the id to be %d but it was %d", idWanted, candidate.ID)
	}
	if candidate.FirstName != firstNameWanted {
		t.Errorf("Expected the first name to be %s but it was %s", firstNameWanted, candidate.FirstName)
	}
	if candidate.LastName != lastNameWanted {
		t.Errorf("Expected the last name to be %s but it was %s", lastNameWanted, candidate.LastName)
	}
	if candidate.EmailAddress != emailWanted {
		t.Errorf("Expected the email to be %s but it was %s", emailWanted, candidate.EmailAddress)
	}
}

func clearDB() {
	collection := config.GetCollection("candidates")
	collection.DropCollection()
}
