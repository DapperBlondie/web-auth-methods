package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"
)

type AppConf struct {
}

type Status struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

var Conf *AppConf

func NewConfiguration() {
	Conf = &AppConf{}
}

func dResponseWriter(w http.ResponseWriter, data interface{}, HStat int) error {
	dataType := reflect.TypeOf(data)
	if dataType.Kind() == reflect.String {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/text")

		w.Write([]byte(data.(string)))
		return nil
	} else if reflect.PtrTo(reflect.TypeOf(reflect.Struct)) == dataType {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/json")

		outData, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
			return err
		}

		w.Write(outData)
		return nil
	} else if reflect.Struct == dataType.Kind() {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/json")

		outData, err := json.MarshalIndent(data, "", "\t")
		if err != nil {
			log.Println(err.Error())
			w.Write([]byte(err.Error()))
			return err
		}

		w.Write(outData)
		return nil
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
