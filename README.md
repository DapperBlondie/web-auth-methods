# web-auth-methods
 A Web application for implementing JWT, HMAC and OAUTH2 functionalities

## HAMC Signing Methods
 I use a database model for saving and getting any request data and hmac token 
 associate with it known as DataModel in repo pkg.<br>
 use a helper function for creating a unique key based on their own emails.
 ```go
// keyGeneratorByEmail a helper function for creating unique keys based on users emails
func keyGeneratorByEmail(mail string) string {
	key := uuid.FromBytesOrNil([]byte(mail))

	return key.String()
}
 ```
<br>
I use two function for saving and getting HMAC token from session that stores in browser.

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
