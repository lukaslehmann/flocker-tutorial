package main

import (
	"encoding/json"
	// "fmt"
	"database/sql"
	"github.com/coopernurse/gorp"
	"github.com/go-martini/martini"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	// "os/exec"
)

var dbmap *gorp.DbMap

//createContainer is json payload used in post
type CreateContainerStruct struct {
	Guid           string
	Template_guid  string `json:"template_guid"`
	Description    string
	Runtime_params map[string]string
	Config_params  map[string]string
}

type Container struct {
	Guid           string    `db:"guid"`
	Description    string    `db:"description"`
	Templates_guid string    `db:"templates_guid"`
	Host_guid      string    `db:"hosts_guid"`
	Port           int       `db:"port"`
	Created_at     time.Time `db:"created_at"`
}

type ContainerConfigParams struct {
	Container_guid string `db:"containers_guid"`
	Name           string `db:"name"`
	Value          string `db:"value"`
}

//containerStartHandler starts the container
func containerStartHandler(res http.ResponseWriter, r *http.Request) {
	// cmd := exec.Command("flocker-deploy", "postgres-deployment.yml", "postgres-application.yml")
	// out, err := cmd.CombinedOutput()
	// if err != nil {
	// 	log.Print(err)
	// 	log.Printf(string(out))
	// 	res.WriteHeader(500)
	// 	fmt.Fprint(res, string(out))
	// } else {
	// 	fmt.Fprint(res, `{ "success": "true" }`)
	// }

	var createContainer CreateContainerStruct
	body, err := ioutil.ReadAll(r.Body)
	if err != null {
		res.WriteHeader(500)
		fmt.Fprint(res, string(err))
		checkErr(err, "read body failed")
	}

	err = json.Unmarshal(body, &createContainer)
	if err != null {
		res.WriteHeader(500)
		fmt.Fprint(res, string(err))
		checkErr(err, "read body failed")
	}
	checkErr(err, "unmarshal failed")

	log.Println(createContainer.Runtime_params["memory"])
	log.Println(createContainer.Config_params["password"])

	container := newContainer(createContainer.Guid, createContainer.Description, "1", "1")
	if err != null {
		res.WriteHeader(500)
		fmt.Fprint(res, string(err))
		checkErr(err, "read body failed")
	}

	//insert container
	err = dbmap.Insert(&container)
	checkErr(err, "Insert failed")

	//insert container params
	for name, value := range createContainer.Config_params {
		containerConfigParams := &ContainerConfigParams{createContainer.Guid, name, value}
		err = dbmap.Insert(containerConfigParams)
		checkErr(err, "Insert failed")
	}

}

func newContainer(guid, description, templates_guid, host_guid string) Container {
	//get port for container
	return Container{
		Guid:           guid,
		Description:    description,
		Templates_guid: templates_guid,
		Host_guid:      host_guid,
		Port:           3306,
		Created_at:     time.Now().UTC(),
	}
}

func main() {
	// initialize the DbMap
	dbmap = initDb()

	m := martini.Classic()
	m.Post("/containers", containerStartHandler)
	m.Run()
}

func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("mysql", "root:mysecretpassword@tcp(172.16.255.250:3306)/service_manager")
	checkErr(err, "DB Connection failed")

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	// add a table, setting the table name to 'posts' and
	dbmap.AddTableWithName(Container{}, "containers")
	dbmap.AddTableWithName(ContainerConfigParams{}, "config_params_containers")
	return dbmap
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Println(msg, err)
	}
}
