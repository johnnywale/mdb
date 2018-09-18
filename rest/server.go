package rest

import (
	"fmt"
	"github.com/johnnywale/mdb/dao"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type restServer struct {
	dbDao *dao.FdbDao
}

func (rs *restServer) Read(path []string) {

	rs.dbDao.Load(path)

}

func YourHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Gorilla!\n"))
}

//type Handler interface {
//	ServeHTTP(ResponseWriter, *Request)
//}

var myHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user")
	fmt.Fprintf(w, "This is an authenticated request")
	fmt.Fprintf(w, "Claim content:\n")
	for k, v := range user.(*jwt.Token).Claims.(jwt.MapClaims) {
		fmt.Fprintf(w, "%s :\t%#v\n", k, v)
	}
})

func (rs *restServer) Start() {
	r := mux.NewRouter()
	mw := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})

	//ServeHTTP(ResponseWriter, *Request)
	var get = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//		user := context.Get(r, "user")
		//		fmt.Fprintf(w, "This is an authenticated request")
		//		fmt.Fprintf(w, "Claim content:\n")
		//		for k, v := range user.(*jwt.Token).Claims.(jwt.MapClaims) {
		//			fmt.Fprintf(w, "%s :\t%#v\n", k, v)
		//		}
	})

	r.Handle("/", get)
	r.Handle("/api", negroni.New(
		negroni.HandlerFunc(mw.HandlerWithNext),
		negroni.Wrap(myHandler),
	))
	log.Fatal(http.ListenAndServe(":8000", r))
}

func NewRestServer(db *dao.FdbDao) *restServer {
	return &restServer{dbDao: db}
}
