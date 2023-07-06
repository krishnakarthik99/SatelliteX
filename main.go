package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/IBM/cloudant-go-sdk/cloudantv1"
	"github.com/IBM/go-sdk-core/core"

	// "github.com/thedevsaddam/renderer"
	"html/template"

	"golang.org/x/crypto/bcrypt"
)

func connection(filename string) *cloudantv1.CloudantV1 {
	var method = "connection()"
	log.Printf("Entered : %s method", method)
	// Fetch CLOUDANT_URL,CLOUDANT_APIKEY from properties file
	// CLOUDANT_URL, CLOUDANT_APIKEY, CLOUDANT_DATABASE := loadProperties(filename)
	CLOUDANT_URL := "https://80adaedf-127d-4ef0-925b-4fd1c5ccb529-bluemix.cloudantnosqldb.appdomain.cloud"
	CLOUDANT_APIKEY := "hlaiPdLs7PoDiyfeGND4q88mubU39QF78VFYUUJjtaxC"
	CLOUDANT_DATABASE := "userdata"
	// Printing the properties fetched
	fmt.Println("CLOUDANT_URL :", CLOUDANT_URL)
	fmt.Println("CLOUDANT_APIKEY :", CLOUDANT_APIKEY)
	fmt.Println("CLOUDANT_DATABASE :", CLOUDANT_DATABASE)
	authenticator := &core.IamAuthenticator{
		ApiKey: CLOUDANT_APIKEY,
	}
	service, err := cloudantv1.NewCloudantV1(
		&cloudantv1.CloudantV1Options{
			URL:           CLOUDANT_URL,
			Authenticator: authenticator,
		},
	)
	if err != nil {
		panic(err)
	}

	// var indexField cloudantv1.IndexField
	// indexField.SetProperty("FirstName", core.StringPtr("asc"))

	// postIndexOptions := service.NewPostIndexOptions(
	// 	"userdata",
	// 	&cloudantv1.IndexDefinition{
	// 		Fields: []cloudantv1.IndexField{
	// 			indexField,
	// 		},
	// 	},
	// )
	// postIndexOptions.SetDdoc("json-index")
	// postIndexOptions.SetName("getUserByRole")
	// postIndexOptions.SetType("json")

	// indexResult, response, err := service.PostIndex(postIndexOptions)
	// fmt.Print(response)
	// if err != nil {
	// 	panic(err)
	// }

	// b, _ := json.MarshalIndent(indexResult, "", "  ")
	// fmt.Println(string(b))

	log.Printf("Exited : %s method", method)
	return service
}

// func loadProperties(filename string) {
// 	panic("unimplemented")
// }

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func userSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Signup")
	if r.Method == "GET" {
		http.ServeFile(w, r, "form.html")
	} else {
		r.ParseForm()
		password, _ := HashPassword(r.FormValue("password"))
		email := r.FormValue("email")
		// fmt.Print("email")
		// fmt.Printf("%T", email)
		service := connection("properties")
		userDOC := cloudantv1.Document{
			ID: &email,
		}
		userDOC.SetProperty("FirstName", r.FormValue("firstname"))
		userDOC.SetProperty("LastName", r.FormValue("lastname"))
		userDOC.SetProperty("Password", password)
		userDOC.SetProperty("EmployeId", r.FormValue("employeid"))
		userDOC.SetProperty("Role", r.FormValue("role"))
		userDOC.SetProperty("BU", r.FormValue("bu"))
		userDOC.SetProperty("WorkLocation", r.FormValue("worklocation"))
		postDocumentOptions := service.NewPostDocumentOptions(
			"userdata",
		)
		postDocumentOptions.SetDocument(&userDOC)
		documentResult, response, err := service.PostDocument(postDocumentOptions)
		if err != nil {
			fmt.Print("error", err)
		}

		b, _ := json.MarshalIndent(documentResult, "", "  ")
		fmt.Println(string(b), response)
		http.ServeFile(w, r, "success.html")
	}

}

func userQuery(w http.ResponseWriter, r *http.Request) {
	fmt.Print("userQuery")
	if r.Method == "GET" {
		http.ServeFile(w, r, "query.html")
	} else {
		service := connection("properties")
		r.ParseForm()
		role := r.FormValue("role")

		postFindOptions := service.NewPostFindOptions(
			"userdata",
			map[string]interface{}{
				"Role": map[string]string{
					"$eq": role,
				},
			},
		)
		postFindOptions.SetFields(
			[]string{"FirstName", "LastName", "Email", "Role", "BU", "WorkLocation", "EmployeId"},
		)

		findResult, response, err := service.PostFind(postFindOptions)
		fmt.Print("Response starting", response, "response ended")
		fmt.Println(findResult)
		if err != nil {
			panic(err)
		}
		b, _ := json.MarshalIndent(findResult, "", "  ")
		fmt.Printf("%T", b)

		type Data struct {
			Bookmark string            `json:"bookmark"`
			Docs     []json.RawMessage `json:"docs"`
		}
		var data Data
		error := json.Unmarshal(b, &data) //convert byte to json
		fmt.Print(error)
		str := data.Docs
		j, err := json.Marshal(&str) //convert data.docs which is of data type json.RawMessage to string
		t, _ := template.ParseFiles("roles.html")
		t.Execute(w, string(j))
	}
}

func userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Print("login")
	if r.Method == "GET" {
		http.ServeFile(w, r, "login.html")
	} else {
		r.ParseForm()
		id := r.FormValue("email")
		pwd := r.FormValue("password")
		service := connection("properties")
		getDocumentOptions := service.NewGetDocumentOptions(
			"userdata",
			id,
		)
		document, response, err := service.GetDocument(getDocumentOptions)
		fmt.Print(response)
		if err != nil {
			fmt.Print(err)
			fmt.Print("Email Not Found")
			http.ServeFile(w, r, "error.html")
		} else {
			fmt.Println(document.GetProperty("Password"))
			if CheckPasswordHash(pwd, document.GetProperty("Password").(string)) {
				http.ServeFile(w, r, "success.html")
			} else {
				fmt.Print("Password Mismatch")
				http.ServeFile(w, r, "error.html")
			}
		}

	}
}

// func userDelete(w http.ResponseWriter, r *http.Request) {
// 	fmt.Print("delete")
// 	if r.Method == "GET" {
// 		http.ServeFile(w, r, "login.html")
// 	} else {
// 		r.ParseForm()
// 		id := r.FormValue("email")

// 	}
// }

func main() {
	http.HandleFunc("/signup", userSignup)
	http.HandleFunc("/login", userLogin)
	http.HandleFunc("/query", userQuery)
	fmt.Println(http.ListenAndServe(":8080", nil))
}
