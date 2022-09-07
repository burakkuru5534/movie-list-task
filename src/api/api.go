package api

import (
	"encoding/json"
	"example.com/m/v2/src/helper"
	"example.com/m/v2/src/model"
	"github.com/Shyp/go-dberror"
	_ "github.com/letsencrypt/boulder/db"
	"net/http"
)

func MovieCreate(w http.ResponseWriter, r *http.Request) {

	var movie model.Movie

	err := helper.BodyToJsonReq(r, &movie)
	if err != nil {
		http.Error(w, "{\"error\": \"Bad request\"}", http.StatusBadRequest)
		return
	}

	err = movie.Create()
	if err != nil {
		dberr := dberror.GetError(err)
		switch e := dberr.(type) {
		case *dberror.Error:
			if e.Code == "23505" {
				http.Error(w, "{\"error\": \"Movie with that email already exists\"}", http.StatusForbidden)
				return
			}
		}

		http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
		return
	}

	respBody := struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		Time  int64  `json:"time"`
	}{
		ID:    movie.ID,
		Title: movie.Title,
		Time:  movie.Time,
	}

	json.NewEncoder(w).Encode(respBody)

}

func MovieUpdate(w http.ResponseWriter, r *http.Request) {

	var movie model.Movie

	//id := helper.StrToInt64(chi.URLParam(r, "id"))
	id := helper.StrToInt64(r.URL.Query().Get("id"))

	isExists, err := helper.CheckIfMovieExists(id)
	if err != nil {
		http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
		return
	}

	if !isExists {
		http.Error(w, "{\"error\": \"Movie with that id does not exist\"}", http.StatusNotFound)
		return
	}

	err = helper.BodyToJsonReq(r, &movie)
	if err != nil {
		http.Error(w, "{\"error\": \"Bad request\"}", http.StatusBadRequest)
		return
	}

	err = movie.Update(id)
	if err != nil {
		dberr := dberror.GetError(err)
		switch e := dberr.(type) {
		case *dberror.Error:
			if e.Code == "23505" {
				http.Error(w, "{\"error\": \"Movie with that email already exists\"}", http.StatusForbidden)
				return
			}
		}

		http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
		return
	}

	respBody := struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
		Time  int64  `json:"time"`
	}{
		ID:    id,
		Title: movie.Title,
		Time:  movie.Time,
	}
	json.NewEncoder(w).Encode(respBody)

}

func MovieDelete(w http.ResponseWriter, r *http.Request) {

	var movie model.Movie

	//id := helper.StrToInt64(chi.URLParam(r, "id"))
	id := helper.StrToInt64(r.URL.Query().Get("id"))

	isExists, err := helper.CheckIfMovieExists(id)
	if err != nil {
		http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
		return
	}

	if !isExists {
		http.Error(w, "{\"error\": \"Movie with that id does not exist\"}", http.StatusNotFound)
		return
	}

	err = movie.Delete(id)
	if err != nil {
		http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("ok")

}

func MovieGet(w http.ResponseWriter, r *http.Request) {

	var movie model.Movie

	id := helper.StrToInt64(r.URL.Query().Get("id"))
	//id := helper.StrToInt64(chi.URLParam(r, "id"))

	isExists, err := helper.CheckIfMovieExists(id)
	if err != nil {
		http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
		return
	}

	if !isExists {
		http.Error(w, "{\"error\": \"Movie with that id does not exist\"}", http.StatusNotFound)
		return
	}

	err = movie.Get(id)
	if err != nil {
		http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
		return
	}

	respBody := struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
		Time  int64  `json:"time"`
	}{
		ID:    id,
		Title: movie.Title,
		Time:  movie.Time,
	}
	json.NewEncoder(w).Encode(respBody)

}

func MovieList(w http.ResponseWriter, r *http.Request) {

	var movie model.Movie

	MovieList, err := movie.GetAll()
	if err != nil {
		http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(MovieList)

}

func MovieListByTime(w http.ResponseWriter, r *http.Request) {

	flightDuration := helper.StrToInt64(r.FormValue("duration"))

	var selectedMovies []model.Movie
	index1, index2 := 0, 0
	var movie model.Movie

	MovieList, err := movie.GetAllByTime()
	if err != nil {
		http.Error(w, "{\"error\": \"server error\"}", http.StatusInternalServerError)
		return
	}

	if MovieList != nil {

		for i, _ := range MovieList {
			for j := 0; j < len(MovieList); j++ {
				if i == j {
					continue
				}
				if MovieList[i].Time+MovieList[j].Time <= flightDuration-30 {
					index1 = i
					index2 = j
					break
				}
			}
			if (index1 > 0 || index2 > 0) && (index1 != index2) {
				break
			}

		}
	}

	if index1 > 0 || index2 > 0 {
		selectedMovies = append(selectedMovies, MovieList[index1])
		selectedMovies = append(selectedMovies, MovieList[index2])

		json.NewEncoder(w).Encode(selectedMovies)

	} else {
		json.NewEncoder(w).Encode("No movies founded")

	}

}
