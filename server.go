package main

import (
	_ "github.com/lib/pq"
	"encoding/json"
	"net/http"
	"strconv"
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"os"
	"github.com/auth0-community/auth0"
	"fmt"
	"gopkg.in/square/go-jose.v2"
	"time"
	"log"
	"io"

	"io/ioutil"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

type Pet struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Owner   string `json:"owner"`
}

type Pet1 struct {
	Id      int    `json:"id"`
	PetName    string `json:"petName"`
	PetTypeBreed   string `json:"typeBreed"`
	PetAge   string `json:"age"`
	PetDescription   string `json:"description"`
	OwnerName   string `json:"ownerName"`
	OwnerAddress   string `json:"ownerAddress"`
	OwnerPhone   string `json:"ownerPhone"`
	OwnerEmail   string `json:"ownerEmail"`
	OwnerPreferredContact   string `json:"preferredContact"`
	VetName   string `json:"vetName"`
	VetAddress   string `json:"vetAddress"`
	VetPhone   string `json:"vetPhone"`
	VetEmail   string `json:"vetEmail"`
	EmergencyProcedure   string `json:"emergencyProcedure"`
	PetFeedingTimes   string `json:"petFeedingTimes"`
	PetFoodAndQuantity   string `json:"petFoodAndQuantity"`
	PetAllowedTreats   string `json:"petAllowedTreats"`
	PetSleepLocation   string `json:"petSleepLocation"`
	PetBedtime   string `json:"petBedtime"`
	PetSleepingHabits   string `json:"petSleepingHabits"`
	PetOkWithKidsAndOtherAnimals   string `json:"petOkWithKidsAndOtherAnimals"`
	PetOkOffLeash   string `json:"petOkOffLeash"`
	PetSicknessHistory   string `json:"petSicknessHistory"`
	PetDesexed   string `json:"petDesexed"`
	PetVaccinated   string `json:"petVaccinated"`
	PetToiletTrained   string `json:"petToiletTrained"`
	PetChews   string `json:"petChews"`
	OwnerHappyToReimburse   string `json:"ownerHappyToReimburse"`
	PetAnxieties   string `json:"petAnxieties"`
	PetDisallowedActivities   string `json:"petDisallowedActivities"`
	PetOtherDetails   string `json:"petOtherDetails"`
}

type Stay struct {
	Id      int    `json:"id"`
	PetId    int `json:"petId"`
	FromDate    time.Time `json:"fromDate"`
	ToDate    time.Time `json:"toDate"`
	Rate    float64 `json:"rate"`
	OtherFee1    float64 `json:"otherFee1"`
	OtherFeeType1    string `json:"otherFeeType1"`
	OtherFee2    float64 `json:"otherFee2"`
	OtherFeeType2    string `json:"otherFeeType2"`
	OtherFee3    float64 `json:"otherFee3"`
	OtherFeeType3    string `json:"otherFeeType3"`
	Discount    float64 `json:"discount"`
	DiscountReason    string `json:"discountReason"`
	TotalCharged    float64 `json:"totalCharged"`
	TotalPaid    float64 `json:"totalPaid"`
}

var skipAuth  = true

// connect to the Db
func init() {
	var err error
	Db, err = sql.Open("postgres", "user=fido dbname=petstay password=woof sslmode=disable")
	if err != nil {
		panic(err)
	}
}

func InitLog(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {

	InitLog(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./views/")))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.Handle("/pets", authMiddleware(GetPetsHandler)).Methods("GET")
	r.Handle("/pet", authMiddleware(PostPetHandler)).Methods("POST")
	r.Handle("/pet/{id}", authMiddleware(GetPetHandler)).Methods("GET")
	r.Handle("/pet/{id}", authMiddleware(PutPetHandler)).Methods("PUT")
	r.Handle("/pet/{id}", authMiddleware(DeletePetHandler)).Methods("DELETE")

	r.Handle("/pets1", authMiddleware(GetPetsHandler1)).Methods("GET")
	r.Handle("/pet1", authMiddleware(PostPetHandler1)).Methods("POST")

	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))

}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if skipAuth {
			next.ServeHTTP(w, r)
			return
		}

		// JC: these are from Auth0 management console
		secret := []byte("pwsy9HvbACAKQYlw1Rp1EAL6ej2OfCZ3")
		audience := "https://jamez.com"
		issuer := "https://speak2jezza.au.auth0.com/"
		secretProvider := auth0.NewKeyProvider(secret)

		var audiences []string
		configuration := auth0.NewConfiguration(secretProvider, append(audiences, audience), issuer, jose.HS256)
		validator := auth0.NewValidator(configuration)

		token, err := validator.ValidateRequest(r)

		if err != nil {
			Error.Println(err)
			//fmt.Println(err)
			fmt.Println("Token is not valid:", token)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized\n"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

var GetPetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		Error.Println("Unable to determine pet id")
		return
	}

	pet, err := retrievePet(id)
	if err != nil {
		Error.Println("Error retrieving pet data")
		return
	}

	output, err := json.MarshalIndent(&pet, "", "\t\t")
	if err != nil {
		Error.Println("Unable to marshall JSON data")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
})

var PostPetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var pet Pet
	json.Unmarshal(body, &pet)
	err := pet.create()
	if err != nil {
		Error.Println("Unable to create pet record")
		return
	}
	w.WriteHeader(200)
	return
})

var PutPetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return
	}
	pet, err := retrievePet(id)
	if err != nil {
		Error.Println("Unable to update pet record")
		return
	}
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	json.Unmarshal(body, &pet)
	err = pet.updatePet()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
})

var DeletePetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return
	}
	pet, err := retrievePet(id)
	if err != nil {
		return
	}
	err = pet.deletePet()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
})

var GetPetsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vals := r.URL.Query()
	offsets :=  vals["offset"]

	offset := 0
	if offsets != nil {
		offset1, err := strconv.Atoi(offsets[0])
		if err != nil {
			return
		}
		offset = offset1
	}
	limits :=  vals["limit"]

	limit := 1
	if limits != nil {
		limit1, err := strconv.Atoi(limits[0])
		if err != nil {
			return
		}
		limit = limit1
	}

	pets, err := retrievePets(offset, limit)
	if err != nil {
		return
	}
	output, err := json.MarshalIndent(&pets, "", "\t\t")
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return

})
var GetPetsHandler1 = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vals := r.URL.Query()
	offsets :=  vals["offset"]

	offset := 0
	if offsets != nil {
		offset1, err := strconv.Atoi(offsets[0])
		if err != nil {
			return
		}
		offset = offset1
	}
	limits :=  vals["limit"]

	limit := 1
	if limits != nil {
		limit1, err := strconv.Atoi(limits[0])
		if err != nil {
			return
		}
		limit = limit1
	}

	pets, err := retrievePets1(offset, limit)
	if err != nil {
		return
	}
	output, err := json.MarshalIndent(&pets, "", "\t\t")
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return

})

var PostPetHandler1 = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var pet Pet1
	json.Unmarshal(body, &pet)
	err := pet.create()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
})
