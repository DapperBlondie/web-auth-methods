package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DapperBlondie/web-auth-methods/src/repo"
	"github.com/alexedwards/scs/v2"
	"github.com/dgrijalva/jwt-go"
	"hash"
	"log"
	"net/http"
	"reflect"
	"time"
)

type AppConf struct {
	ScsManager   *scs.SessionManager
	DRepo        *repo.DBRepo
	HashFunction func() hash.Hash
}

type Status struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

var Conf *AppConf

// NewConfiguration use for creating a new configuration for our handlers
func NewConfiguration(manager *scs.SessionManager, repo *repo.DBRepo) {
	Conf = &AppConf{
		ScsManager:   manager,
		DRepo:        repo,
		HashFunction: sha256.New,
	}
}

// dResponseWriter use for writing response to the user
func dResponseWriter(w http.ResponseWriter, data interface{}, HStat int) error {
	dataType := reflect.TypeOf(data)
	if dataType.Kind() == reflect.String {
		w.WriteHeader(HStat)
		w.Header().Set("Content-Type", "application/text")

		_, err := w.Write([]byte(data.(string)))
		return err
	} else if reflect.PtrTo(dataType).Kind() == dataType.Kind() {
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

// CreateJWTToken a helper function for creating JWT token
func CreateJWTToken(user *repo.DataModel) string {
	claims := &repo.UserClaims{
		StandardClaims: &jwt.StandardClaims{
			Audience:  "localhost:8080",
			ExpiresAt: time.Now().Add(time.Hour * 24).UnixNano(),
			Id:        string(rune(user.ID)),
			IssuedAt:  time.Now().UnixNano(),
			Issuer:    "localhost:8080",
			NotBefore: time.Now().UnixNano(),
			Subject:   "Jwt token for login",
		},
		Email: user.Mail,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedT, err := token.SignedString(user.Key)
	if err != nil {
		log.Println(err.Error())
		return ""
	}

	return signedT
}

// keyGeneratorByEmail a helper function for creating unique keys based on users emails
func keyGeneratorByEmail(mail string) string {
	en := base64.StdEncoding

	return en.EncodeToString([]byte(mail))
}

// CheckStatusHandler just for checking the API handlers
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

// SaveHmacToken use for saving HMAC token based on sha hash functions for user
func (conf *AppConf) SaveHmacToken(w http.ResponseWriter, r *http.Request) {
	var user *repo.DataModel = &repo.DataModel{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error in parsing the body; "+err.Error(), http.StatusInternalServerError)
		return
	}
	user.Key = keyGeneratorByEmail(user.Mail)
	signToken, err := conf.SignWithHmac(user.Mail, user.Key)
	if err != nil {
		http.Error(w, err.Error()+"; in signing with hmac", http.StatusInternalServerError)
		return
	}

	err = conf.DRepo.SaveUserWithHAMCMethod(user)
	if err != nil {
		log.Println(err.Error() + "; error in saving user in db")
		http.Error(w, err.Error()+"; error in saving user in db", http.StatusInternalServerError)
		return
	}

	conf.ScsManager.Put(r.Context(), "hmac-token", signToken)
	conf.ScsManager.Put(r.Context(), "user-mail", user.Mail)
	return
}

// SignWithHmac use for creating HMAC tokens
func (conf *AppConf) SignWithHmac(userMail string, key string) (string, error) {
	h := hmac.New(conf.HashFunction, []byte(key))

	_, err := h.Write([]byte(userMail))
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// GetAndCheckHmacToken use for getting and checking the HMAC token that we store it in cookies
func (conf *AppConf) GetAndCheckHmacToken(w http.ResponseWriter, r *http.Request) {
	userEmail, ok := conf.ScsManager.Get(r.Context(), "user-mail").(string)
	if !ok {
		http.Error(w, "Something went wrong; ", http.StatusInternalServerError)
		return
	}

	hmacToken, ok := conf.ScsManager.Get(r.Context(), "hmac-token").(string)
	if !ok {
		http.Error(w, "Something went wrong; ", http.StatusInternalServerError)
		return
	}

	userKey, err := conf.DRepo.GetUserByItsEmailHMACMethod(userEmail)
	if err != nil {
		log.Println(err.Error())
		return
	}

	user := &repo.DataModel{
		ID:        0,
		Mail:      userEmail,
		HmacToken: hmacToken,
		Key:       userKey,
	}

	err = dResponseWriter(w, user, http.StatusOK)
	if err != nil {
		log.Println(err.Error())
		return
	}

	return
}

// SaveJWTToken use for store JWT token in scs.Session
func (conf *AppConf) SaveJWTToken(w http.ResponseWriter, r *http.Request) {
	var user *repo.DataModel = &repo.DataModel{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Error in parsing the body; "+err.Error(), http.StatusInternalServerError)
		return
	}
	user.Key = keyGeneratorByEmail(user.Mail)

	t := CreateJWTToken(user)
	conf.ScsManager.Put(r.Context(), "jwt-token", t)
	conf.ScsManager.Put(r.Context(), "user-key", user.Key)

	err = dResponseWriter(w, user, http.StatusOK)
	if err != nil {
		log.Println(err.Error())
		return
	}

	return
}

// ParseJWTToken use for parsing JWT token with claims
func (conf *AppConf) ParseJWTToken(w http.ResponseWriter, r *http.Request) {
	jt, ok := conf.ScsManager.Get(r.Context(), "jwt-token").(string)
	key, ok := conf.ScsManager.Get(r.Context(), "jwt-token").(string)
	if !ok {
		http.Error(w, "Something Went Wrong", http.StatusInternalServerError)
		return
	}

	jwtToken, err := jwt.ParseWithClaims(jt, &repo.UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Valid {
			return []byte(key), nil
		}
		return nil, errors.New("token is not valid")
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
	}

	userClaims := jwtToken.Claims.(*repo.UserClaims)

	err = dResponseWriter(w, userClaims.Email, http.StatusOK)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}
