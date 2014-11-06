package middleware

import (
  "net/http"
  "appengine"
  "appengine/datastore"
  "strconv"
  "fmt"
  "io/ioutil"
  "strings"
  "encoding/json"
  "github.com/gorilla/securecookie"
  "vidao/martini-tools"
  )

/*type Loggedinusers struct {
	UID    int64
	SID    int64
	Extime  int64
}*/


var key1 = []byte("5916569511133184")
var key2 = []byte("4776259720577024")
var CookieHandler = securecookie.New(key1, key2)

type loggedinusers tools.Loggedinusers
type users tools.Users 

func SessionAuth(w http.ResponseWriter, r *http.Request){
  var userid string
	c := appengine.NewContext(r)
  	cookie, err := r.Cookie("session")
  	var sessionid string
	if err == nil {
		cookieValue := make(map[string]string)
		if err = CookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			sessionid = cookieValue["sessionid"]
			sid, _ := strconv.ParseInt(sessionid, 10, 64)
			qClient_user := datastore.NewQuery("login").
								Filter("SID =", sid)
			var currentuser []loggedinusers
			_, err := qClient_user.GetAll(c, &currentuser)
			if err != nil{
				fmt.Fprint(w,err)
				return 
			}
			userid = strconv.FormatInt(currentuser[0].UID, 10)
		}
	}

  if(userid=="0" || userid==""){
    http.Error(w, err.Error(), http.StatusUnauthorized)
  }else{
  	(*r).Header.Set("SessionID",sessionid)
    (*r).Header.Set("UserID",userid)
  }
  
  return
}

func BasicAuth(w http.ResponseWriter, r *http.Request){
  c := appengine.NewContext(r)
  username:=r.Header.Get("Username")
  password:=tools.Hash256(r.Header.Get("Password"))

  qClient_user := datastore.NewQuery("User").
    Filter("Username =", username).
    Filter("Password =", password)

  var currentuser []users
  qClient_user.GetAll(c, &currentuser)
  userid := strconv.FormatInt(currentuser[0].UID, 10)

  if(userid=="0" || userid==""){
    http.Error(w,"Unauthorized", http.StatusUnauthorized)
  }else{
    (*r).Header.Add("UserID",userid)
  }

  return
}

func Translator(w http.ResponseWriter, r *http.Request) {
  if strings.Contains(r.Header.Get("Content-Type"), "json") == true{
  		jsonBinInfo, err := ioutil.ReadAll(r.Body)
  		r.Body.Close()
  		if err != nil {
  			fmt.Fprintln(w, err)
  		}

      (*r).Header.Add("InComingData",fmt.Sprintf("%s",jsonBinInfo))

  }else if strings.Contains(r.Header.Get("Content-Type"), "form-urlencoded") == true{
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
    m:=make(map[string]string)
    for i, v:= range r.Form{
      m[i]=v[0]
    }
    b,_:=json.Marshal(m)
    (*r).Header.Add("InComingData",fmt.Sprintf("%s",b))
  }else{
    http.Error(w,"Invalid Input Request", http.StatusBadRequest)
  }
}
