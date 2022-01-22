package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"WarmAPI/config"
	"WarmAPI/models"

	"github.com/julienschmidt/httprouter"
	"sync"
)

// WarmResponse that represents the response in the JSON requests.
type WarmResponse struct {
	Message   string                       `json:"message"`
	Status    int                          `json:"status"`
	Candidate *models.ChallengerCandidate  `json:"candidate,omitempty"`
	Data      []models.ChallengerCandidate `json:"data,omitempty"`
}

/*
*	Public Functions
 */

// RetrieveCandidates function that handles the endpoint for retrieving all the candidates
// inside the database.
func RetrieveCandidates(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)

	candidates, _ := models.GetAllCandidates()
	fmt.Printf("%v", len(*candidates))
	if len(*candidates) == 0 {
		w.WriteHeader(http.StatusNotFound)

		jsonResp, err := json.Marshal(
			WarmResponse{
				Message: "No candidates were found.",
				Status:  http.StatusNotFound})

		config.CheckError(err)
		fmt.Fprintf(w, "%s", jsonResp)
		return
	}

	data := []models.ChallengerCandidate{}

	for _, v := range *candidates {
		data = append(data, v)
	}

	jsonResp, _ := json.Marshal(
		WarmResponse{
			Message: "Retrieve Candidates",
			Status:  http.StatusOK,
			Data:    data})

	fmt.Fprintf(w, "%s", jsonResp)
}

// RetrieveCandidateID function that handles the endpoint for retrieving a certain candidate
// via his ID in the database.
func RetrieveCandidateID(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {

	candidateID := p.ByName("id")

	if !checkIDOnlyNumbers(candidateID) {
		w.WriteHeader(http.StatusBadRequest)

		jsonResp, err := json.Marshal(
			WarmResponse{
				Message: "The id must be all numbers",
				Status:  http.StatusBadRequest})

		config.CheckError(err)
		fmt.Fprintf(w, "%s", jsonResp)
		return
	}

	id, _ := strconv.Atoi(candidateID)

	candidate, err := models.GetCandidateByID(id)
	if config.CheckError(err) {
		w.WriteHeader(http.StatusNotFound)

		jsonResp, err := json.Marshal(
			WarmResponse{
				Message: "Candidate not found",
				Status:  http.StatusNotFound})

		config.CheckError(err)
		fmt.Fprintf(w, "%s", jsonResp)
		return
	}

	w.WriteHeader(http.StatusOK)

	jsonResp, err := json.Marshal(
		WarmResponse{
			Message:   "Candidate Retrieved",
			Status:    http.StatusOK,
			Candidate: candidate,
		})

	config.CheckError(err)
	fmt.Fprintf(w, "%s", jsonResp)
	return
}

// RetrieveCandidate function that handles the endpoint for retrieving a certain candidate
// via its email either in the database, or via requesting the Greenhouse Harvest API.
func RetrieveCandidate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	testStruct := make(map[string]interface{})

	//Empty body verification
	body, err := ioutil.ReadAll(r.Body)
	config.CheckError(err)
	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		jsonResp, err := json.Marshal(WarmResponse{Message: "The Json body is empty", Status: http.StatusBadRequest})
		config.CheckError(err)
		fmt.Fprintf(w, "%s", jsonResp)
		return
	}

	//Empty object verification
	err = json.Unmarshal(body, &testStruct)
	if len(testStruct) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		jsonResp, err := json.Marshal(WarmResponse{Message: "The Json object is empty", Status: http.StatusBadRequest})
		config.CheckError(err)
		fmt.Fprintf(w, "%s", jsonResp)
		return
	}

	//Empty email field
	if len(testStruct["email"].(string)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		jsonResp, err := json.Marshal(WarmResponse{Message: "The email of the candidate is empty", Status: http.StatusBadRequest})
		config.CheckError(err)
		fmt.Fprintf(w, "%s", jsonResp)
		return
	}

	/***
	* FINDING THE CANDIDATE WITHIN THE DATABASE
	****/
	fmt.Printf("\n\n ** CHECKING THE DATABASE ** \n\n")

	candidate, err := models.GetCandidateByEmail(testStruct["email"].(string))
	// If the candidate was found within the database
	if !config.CheckError(err) {
		w.WriteHeader(http.StatusOK)

		jsonResp, err := json.Marshal(
			WarmResponse{
				Message:   "Retrieved candidate",
				Status:    http.StatusOK,
				Candidate: candidate})

		config.CheckError(err)
		fmt.Fprintf(w, "%s", jsonResp)
		fmt.Printf("%s", jsonResp)
		return
	}

	fmt.Printf("\n\n ** CHECKING THE HARVEST API ** \n\n")
	/***
	* FINDING THE CANDIDATE WITHIN THE HARVEST API
	****/
	challengerCandidate := &models.ChallengerCandidate{}
	pageCount := 0
	for {
		if challengerCandidate.ID == 0 {
			fmt.Printf("Candidate still not found, moving to page [%d]", pageCount+1)
			pageCount++
			client := &http.Client{}
			req, _ := http.NewRequest("GET",
				config.HarvestURL+
					config.PerPageParam+"500&"+
					config.PageParam+strconv.Itoa(pageCount)+"&"+
					config.CreatedAfterParam+config.GetDateForRequest(5),
				nil)
			fmt.Printf("\nURL: %s", req.URL)
			req.SetBasicAuth("1a466133e09f2e2c700f867ead70b002", "")
			reply, _ := client.Do(req)
			fmt.Printf("\n Request Arrived!")

			candidates := []models.Candidate{}

			defer reply.Body.Close()
			body, err := ioutil.ReadAll(reply.Body)
			config.CheckError(err)

			if len(body) == 2 { // 2 because it's empty with "[]"
				w.WriteHeader(http.StatusNotFound)
				jsonResp, err := json.Marshal(
					WarmResponse{
						Message: "No candidate was found within the Greenhouse database.",
						Status:  http.StatusNotFound})
				config.CheckError(err)
				fmt.Fprintf(w, "%s", jsonResp)
				return
			}

			json.Unmarshal(body, &candidates)

			for _, candidate := range candidates {
				for _, email := range candidate.EmailAddresses {
					time.Sleep(5 * time.Millisecond)
					//fmt.Printf("\nComparing Emails: -> [%s] - [%s]", email, testStruct["email"])
					if email.Value == testStruct["email"] {
						fmt.Printf("\n\n-----> EMAIL [%s] FOUND <------\n\n", testStruct["email"])

						challengerCandidate = convertCandidateToChallenger(&candidate)
						models.InsertCandidate(challengerCandidate)

						time.Sleep(1 * time.Second)
						w.WriteHeader(http.StatusOK)
						jsonResp, err := json.Marshal(
							WarmResponse{
								Message:   "Retrieved candidate",
								Status:    http.StatusOK,
								Candidate: challengerCandidate})
						config.CheckError(err)
						fmt.Fprintf(w, "%s", jsonResp)
						fmt.Printf("%s", jsonResp)
						return
					}
				}
			}
			fmt.Printf("\n\n ****************************** \n\n")
			fmt.Printf(" ** > PAGE [%d] ENDED < **", pageCount)
			fmt.Printf("\n\n ****************************** \n\n")
		}
	}
}

func RetrieveCandidateRoutine(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	testStruct := make(map[string]interface{})

	//Empty body verification
	body, err := ioutil.ReadAll(r.Body)
	config.CheckError(err)
	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		jsonResp, err := json.Marshal(WarmResponse{Message: "The Json body is empty", Status: http.StatusBadRequest})
		config.CheckError(err)
		fmt.Fprintf(w, "%s", jsonResp)
		return
	}

	//Empty object verification
	err = json.Unmarshal(body, &testStruct)
	if len(testStruct) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		jsonResp, err := json.Marshal(WarmResponse{Message: "The Json object is empty", Status: http.StatusBadRequest})
		config.CheckError(err)
		fmt.Fprintf(w, "%s", jsonResp)
		return
	}

	//Empty email field
	if len(testStruct["email"].(string)) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		jsonResp, err := json.Marshal(WarmResponse{Message: "The email of the candidate is empty", Status: http.StatusBadRequest})
		config.CheckError(err)
		fmt.Fprintf(w, "%s", jsonResp)
		return
	}

	email := testStruct["email"].(string)

	/********************************************
	*********************************************
	* FINDING THE CANDIDATE WITHIN THE DATABASE *
	*********************************************
	*********************************************/

	fmt.Printf("\n\n ** CHECKING THE DATABASE ** \n\n")

	candidateDB, err := models.GetCandidateByEmail(email)
	// If the candidate was found within the database
	if !config.CheckError(err) {
		w.WriteHeader(http.StatusOK)

		jsonResp, err := json.Marshal(
			WarmResponse{
				Message:   "Retrieved candidate",
				Status:    http.StatusOK,
				Candidate: candidateDB})

		config.CheckError(err)
		fmt.Fprintf(w, "%s", jsonResp)
		fmt.Printf("%s", jsonResp)
		return
	}

	/***********************************************
	************************************************
	* FINDING THE CANDIDATE WITHIN THE HARVEST API *
	************************************************
	************************************************/

	fmt.Printf("\n\n ** CHECKING THE HARVEST API ** \n\n")
	fmt.Printf("Looking for candidate with email [ %s ] in the Harvest API", email)
	var wg sync.WaitGroup
	pageNumbers := 5
	wg.Add(pageNumbers)
	chn := make(chan models.Candidate)

	count := 0
	for {
		if count == pageNumbers {
			break
		}
		count++
		go FetchHarvest(&chn, count, email, &wg)
	}

	candidate := <-chn
	wg.Wait()
	fmt.Printf("\n HERE END \n\n\n\n\n")

	if candidate.ID == 0 {
		fmt.Printf("\nNo candidate was found with the email [ %s ]", email)
		w.WriteHeader(http.StatusNotFound)
		jsonResp, err := json.Marshal(
			WarmResponse{
				Message: "No candidate was found within the Greenhouse database.",
				Status:  http.StatusNotFound})
		config.CheckError(err)
		fmt.Fprintf(w, "%s", jsonResp)
		return
	}

	json.Unmarshal(body, &candidate)

	challengerCandidate := convertCandidateToChallenger(&candidate)
	models.InsertCandidate(challengerCandidate)

	w.WriteHeader(http.StatusOK)
	jsonResp, err := json.Marshal(
		WarmResponse{
			Message:   "Retrieved candidate",
			Status:    http.StatusOK,
			Candidate: challengerCandidate})
	config.CheckError(err)
	fmt.Fprintf(w, "%s", jsonResp)
	fmt.Printf("%s", jsonResp)
	return
}

func FetchHarvest(ch *chan models.Candidate, count int, email string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	url := config.HarvestURL +
		config.PerPageParam + "500&" +
		config.PageParam + strconv.Itoa(count) + "&" +
		config.CreatedAfterParam + config.GetDateForRequest(5)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(config.TokenURL, "")
	reply, _ := client.Do(req)

	cands := []models.Candidate{}

	defer reply.Body.Close()
	body, _ := ioutil.ReadAll(reply.Body)

	if len(body) != 2 { //not empty like -> []
		json.Unmarshal(body, &cands)
		for _, v := range cands {
			for _, e := range v.EmailAddresses {
				if e.Value == email {
					fmt.Printf("\nEmail >> %s << has been found.", email)
					*ch <- v
					close(*ch)
				}
			}
		}
	}
}

/*
*	Private Functions
 */

// checkIDOnlyNumbers function that checks if the id of the candidate only has numbers in its string.
func checkIDOnlyNumbers(id string) bool {
	for _, v := range id {
		if v >= 'a' && v <= 'z' || v >= 'A' && v <= 'Z' {
			return false
		}
	}
	return true
}

// convertCandidateToChallenger function that converts a candidate in the struct type Candidate
// to the struct type ChallengerCandidate.
func convertCandidateToChallenger(candidate *models.Candidate) *models.ChallengerCandidate {
	challengerCandidate := models.ChallengerCandidate{}
	var application struct {
		Job    string `json:"job,omitempty"`
		Status string `json:"status,omitempty"`
	}

	challengerCandidate.ID = candidate.ID
	challengerCandidate.FirstName = candidate.FirstName
	challengerCandidate.LastName = candidate.LastName
	challengerCandidate.EmailAddress = candidate.EmailAddresses[0].Value //first email of the candidate

	for _, apps := range candidate.Applications {
		application.Status = apps.Status
		for _, jobs := range apps.Jobs {
			application.Job = jobs.Name
		}
		challengerCandidate.Applications = append(challengerCandidate.Applications, application)
	}

	return &challengerCandidate
}
