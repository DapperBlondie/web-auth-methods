# web-auth-methods
 A Web application for implementing JWT, HMAC and OAUTH2 functionalities

***

## ResponseWriter
I implement a response writer that you can give it an ``` interface{} ``` which type of it can be ``` string ``` , ``` struct{} ``` ,
 ``` *struct{} ``` then with ``` package reflect ``` in golang function write response in appropriate format and HTTP header.
 
 ```go
 
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
 
 ```

## HMAC Signing Methods
 I use a database model for saving and getting any request data and hmac token 
 associate with it known as DataModel in repo pkg.<br>
 - I use a helper function for creating a unique key based on their own emails.
 - I use two function for saving and getting HMAC token from session that stores in browser.
 - You can find Everything About saving/getting or signing in the following peace of codes. 

```go
// keyGeneratorByEmail a helper function for creating unique keys based on users emails
func keyGeneratorByEmail(mail string) string {
	key := uuid.FromBytesOrNil([]byte(mail))

	return key.String()
}
 ```
<br>

```go
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

	conf.ScsManager.Put(r.Context(), "hmac-token", signToken)
	conf.ScsManager.Put(r.Context(), "user-mail", user.Mail)
	return
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

	user := &repo.DataModel{
		ID:        0,
		Mail:      userEmail,
		HmacToken: hmacToken,
	}

	err := dResponseWriter(w, user, http.StatusOK)
	if err != nil {
		log.Println(err.Error())
		return
	}

	return
}
```