package handlers

import (
	"encoding/json"
	"errors"
	"github.com/DapperBlondie/web-auth-methods/src/repo"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"reflect"
)

type AppConf struct {
	ScsManager *scs.SessionManager
	DRepo      *repo.DBRepo
}

type Status struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

var Conf *AppConf

func NewConfiguration(manager *scs.SessionManager, repo *repo.DBRepo) {
	Conf = &AppConf{
		ScsManager: manager,
		DRepo:      repo,
	}
}

func dResponseWriter(w http.ResponseWriter, data interface{}, HStat int) error {
	dataType := reflect.TypeOf(data)
	if dataType.Kind() == reflect.String {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/text")

		_, err := w.Write([]byte(data.(string)))
		return err
	} else if reflect.PtrTo(reflect.TypeOf(reflect.Struct)) == dataType {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/json")

		outData, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
			return err
		}

		_, err = w.Write(outData)
		return err
	} else if reflect.Struct == dataType.Kind() {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/json")

		outData, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
			return err
		}

		_, err = w.Write(outData)
		return err
	}

	return errors.New("we could not be able to support data type that you passed")
}

func (conf *AppConf) CheckStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, errors.New(http.MethodGet+" use this method").Error(), http.StatusMethodNotAllowed)
		log.Println(errors.New(http.MethodGet + " use this method").Error())
		return
	}

	stat := &Status{
		Ok:      true,
		Message: "Just check the status",
	}

	err := dResponseWriter(w, stat, http.StatusOK)
	if err != nil {
		log.Println(err.Error())
		return
	}

	return
}
