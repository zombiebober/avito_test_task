package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type App struct {
	Router *mux.Router
	DB *sql.DB
}


func (a *App)Initialize(host, port, user, password, dbname, sslmode string)  {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	var err error
	fmt.Println(connectionString)
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter().StrictSlash(true)
	a.initializeRoutes()
}



func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{})  {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) getAdvert(w http.ResponseWriter, r *http.Request){
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "Invalid advert ID")
		return
	}

	advert := Advert{ID : id}

	fields := Fields{numberPhotoLink: 1}
	params := r.URL.Query().Get("fields")

	if params != "" {
		params := strings.Split(params, ",")
		for _, field := range params {
			switch strings.ToLower(field) {
			case "description":
				fields.onDescription = true
			case "photo_links":
				fields.numberPhotoLink = 3
			default:
				respondWithError(w, http.StatusBadRequest, "Invalid fields(description, photo_links)")
			}
		}
	}

	if err := advert.getAdvertAdditionalFields(a.DB, fields); err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Advert not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	advert.ID = 0

	respondWithJSON(w, http.StatusOK, advert)
}

func (a *App) getAllAdverts(w http.ResponseWriter, r *http.Request)  {
	limit := 10
	var sort Sort
	params := r.URL.Query()
	page := 1
	if p := params.Get("page"); p != "" {
		var err error
		page, err = strconv.Atoi(p)
		if page < 1 || err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid page")
			return
		}
	}

	if column := params.Get("sort"); column != "" {
		switch strings.ToLower(column) {
		case "price":
			sort.Column = "Price"
		case "time_create", "time", "date", "creation_date":
			sort.Column = "time_create"
		default:
			respondWithError(w, http.StatusBadRequest,"Invalid sort")
			return
		}
	}

	if sortType := params.Get("sort_type"); sortType != ""{
		switch strings.ToLower(sortType) {
		case "ask":
			sort.Type = Ask
		case "desc":
			sort.Type = Desc
		default:
			respondWithError(w, http.StatusBadRequest,"Invalid sort_type")
			return
		}
	}

	if sort.Column == "" && sort.Type != None {
		respondWithError(w, http.StatusBadRequest,
			"You need to enter the parameter by which the sorting will be (price/creation_date)")
	}

	offset := limit * (page -1)

	adverts, err := getAllAdverts(a.DB, offset, limit, &sort)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(adverts) < 1 {
		respondWithError(w, http.StatusNotFound, "Page does not exist")
		return
	}

	respondWithJSON(w, http.StatusOK, adverts)
}

func (a * App) createAdvert(w http.ResponseWriter, r * http.Request)  {
	var advert Advert
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&advert); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request payload")
		return
	}
	defer r.Body.Close()
	if err := advert.Validate(); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := advert.createAdvert(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, Advert{ID: advert.ID})
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/advert/", a.createAdvert).Methods("POST")
	a.Router.HandleFunc("/advert/{id:[0-9]+}/", a.getAdvert).Methods("GET")
	a.Router.HandleFunc("/adverts/", a.getAllAdverts).Methods("GET")

}


func (a *App) Run(addr string)   {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}