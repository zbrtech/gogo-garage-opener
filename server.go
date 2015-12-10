package main
import (
	"github.com/emicklei/go-restful"
	"net/http"
	 log "github.com/Sirupsen/logrus"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"github.com/spf13/viper"
)

func main() {
	log.SetLevel(log.DebugLevel)
	wsContainer := restful.NewContainer()
	os.Remove("gogo-garage-opener.db")
	db, err := sql.Open("sqlite3", "gogo-garage-opener.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	_, err = db.Exec("CREATE TABLE user (email TEXT NOT NULL PRIMARY KEY, password TEXT, longitude REAL, latitude REAL, time DATETIME, duration INTEGER, distance INTEGER);")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE event (timestamp DATETIME  NOT NULL PRIMARY KEY, email TEXT);")
	if err != nil {
		log.Fatal(err)
	}

	// Need to execute in $GOPATH/src/benjefferies/gogo-garage-opener for some reason... look into
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	userDao := UserDao{*db};
	u := UserResource{userDao}
	e := EventResource{eventDao:EventDao{*db}, userDao:userDao}

	u.Register(wsContainer)
	e.Register(wsContainer)

	server := &http.Server{Addr: ":8080", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
