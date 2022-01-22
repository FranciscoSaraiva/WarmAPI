package models

import (
	"WarmAPI/config"

	"gopkg.in/mgo.v2/bson"
)

// Candidate that represents the candidate as the Greenhouse Harvest API returns it.
type Candidate struct {
	ID           int    `json:"id" bson:"_id,omitempty"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PhoneNumbers []struct {
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"phone_numbers"`
	EmailAddresses []struct {
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"email_addresses"`
	Applications []struct {
		Jobs []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"jobs"`
		Status string `json:"status"`
	} `json:"applications"`
}

// ChallengerCandidate that represents a candidate from the challenger.
type ChallengerCandidate struct {
	ID           int    `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	EmailAddress string `json:"email_address,omitempty"`
	Applications []struct {
		Job    string `json:"job,omitempty"`
		Status string `json:"status,omitempty"`
	} `json:"applications,omitempty"`
}

/*
*	Public Functions
 */

// GetAllCandidates function that returns a list of all the candidates in the database.
func GetAllCandidates() (*[]ChallengerCandidate, error) {
	collection := config.GetCollection("candidates")

	var candidates []ChallengerCandidate

	err := collection.Find(nil).All(&candidates)
	config.CheckError(err)

	return &candidates, nil
}

// GetCandidateByID function that returns a candidate in the database via the id.
func GetCandidateByID(id int) (*ChallengerCandidate, error) {
	collection := config.GetCollection("candidates")

	candidate := ChallengerCandidate{}
	query := bson.M{"_id": id}

	err := collection.Find(query).One(&candidate)
	if config.CheckError(err) {
		return nil, err
	}

	return &candidate, nil
}

// GetCandidateByEmail function that returns a candidate in the database via email.
func GetCandidateByEmail(email string) (*ChallengerCandidate, error) {
	collection := config.GetCollection("candidates")

	candidate := ChallengerCandidate{}
	query := bson.M{"emailaddress": email}

	err := collection.Find(query).One(&candidate)
	if config.CheckError(err) {
		return nil, err
	}

	return &candidate, nil
}

// InsertCandidate function that inserts a new Candidate in the database.
func InsertCandidate(candidate *ChallengerCandidate) error {
	collection := config.GetCollection("candidates")

	err := collection.Insert(&candidate)
	config.CheckError(err)
	return err
}
