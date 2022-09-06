package helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var App *app

type app struct {
	DB *DbHandle
}

func InitApp(db *DbHandle) error {
	App = &app{
		DB: db,
	}

	return nil
}

func BodyToJsonReq(r *http.Request, data interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return errors.New(fmt.Sprintf("Body unmarshall error %s", string(body)))
	}

	defer r.Body.Close()

	return nil
}

func StrToInt64(aval string) int64 {
	aval = strings.Trim(strings.TrimSpace(aval), "\n")
	i, err := strconv.ParseInt(aval, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func CheckIfMovieExists(id int64) (bool, error) {

	var isExists bool

	err := App.DB.Get(&isExists, "SELECT EXISTS (SELECT 1 FROM movie WHERE id = $1)", id)
	if err != nil {
		return false, err
	}
	return isExists, nil

}
