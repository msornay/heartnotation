package annotation

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func verifyCode409(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, http.StatusText(409), 409)
		return
	}
}

func verifyDBConnection() *sql.DB {
	db, err := sql.Open("postgres", "user=postgres password=postgres dbname=heartnotation sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func CreateAnnotation(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	db := verifyDBConnection()

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	var annotation Annotation
	creationDate := time.Now()
	annotation.CreationDate = creationDate.Format("2006-01-02")
	annotation.EditDate = creationDate.Format("2006-01-02")
	json.Unmarshal(bodyBytes, &annotation)

	if annotation.OrganizationID != 0 {
		annotation.ProcessID = 2
	} else {
		annotation.ProcessID = 1
	}

	if annotation.Comment == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	fmt.Printf("%d", annotation.IDAnnotationParent)
	if annotation.IDAnnotationParent != 0 {
		_, err := db.Exec("INSERT INTO ANNOTATION VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)", annotation.IDAnnotation, annotation.IDAnnotationParent, annotation.OrganizationID, annotation.ProcessID, annotation.IDSignal, annotation.Comment, annotation.CreationDate, annotation.EditDate, true)
		verifyCode409(w, err)
	} else {
		_, err := db.Exec("INSERT INTO ANNOTATION VALUES($1, NULL, $2, $3, $4, $5, $6, $7, $8)", annotation.IDAnnotation, annotation.OrganizationID, annotation.ProcessID, annotation.IDSignal, annotation.Comment, annotation.CreationDate, annotation.EditDate, true)
		verifyCode409(w, err)
	}

	fmt.Fprintf(w, "Annotation %d successfully created \n", annotation.IDAnnotation)
}
