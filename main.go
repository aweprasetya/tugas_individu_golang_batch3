package main 

import (
	"fmt"
	//"html"
	"net/http"
	"html/template"
	"path"
	"strings"
	"strconv"
	"time"
	"reflect"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"	
)

//type M map[string]interface{}

type EntitiesClass struct {
	Id int
	Name string 
	Assignee string 
	Status int
	Deadline string
}

func dbConnect() (db *sql.DB) {
	database, _ := sql.Open("sqlite3", "data.db")
	
	return database
}

func dbCreate() {
	database := dbConnect()
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS todo (id INTEGER PRIMARY KEY, name TEXT, assignee TEXT, status INTEGER, deadline INTEGER)")
	
	statement.Exec()

	defer database.Close()
}

func dbSelect() []EntitiesClass {
	database := dbConnect()
	rows, _ := database.Query("SELECT id, name, assignee, status, deadline FROM todo ORDER BY id") 

	dbdata := EntitiesClass{}
    dbdatas := []EntitiesClass{}

	for rows.Next() {
		var id int 
		var name string
		var assignee string
		var status int 
		var deadline int

		rows.Scan(&id, &name, &assignee, &status, &deadline)
		

		dbdata.Id = id
		dbdata.Name = name 
		dbdata.Assignee = assignee 
		dbdata.Status = status 

		dataBaru := time.Unix(int64(deadline), 0)	

		dbdata.Deadline = dataBaru.Format("January 2, 2006")

		dbdatas = append(dbdatas, dbdata)
	
	}

	defer database.Close()

	return dbdatas
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	filepath := path.Join("views","index.html")
	tmpl, _ := template.ParseFiles(filepath)

	table := dbSelect()

	tmpl.Execute(w, table)

}

func handleFormEntry(w http.ResponseWriter, r *http.Request) {
	filepath := path.Join("views","add_todo.html")
	tmpl, _ := template.ParseFiles(filepath)

	tmpl.Execute(w, nil)
}

func handleProcessEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		database := dbConnect()

		name_input := r.FormValue("todo_task") 
		assignee_input := r.FormValue("todo_assignee")
		date_input := r.FormValue("todo_date")
		
		date_array := strings.Split(date_input, "/")

		date_array_date, err := strconv.Atoi(date_array[0])
		if err != nil {
			fmt.Println("Can not convert date")
		}
		date_array_month, err := strconv.Atoi(date_array[1])
		if err != nil {
			fmt.Println("Can not convert month")
		}
		date_array_year, err := strconv.Atoi(date_array[2])
		if err != nil {
			fmt.Println("Can not convert year")
		}
		thisTime := time.Date(date_array_year, time.Month(date_array_month), date_array_date, 0, 0, 0, 0, time.Local)
		
		statement, _ := database.Prepare("INSERT INTO todo (name, assignee, status, deadline) VALUES (?, ?, ?, ?)")
		statement.Exec(name_input, assignee_input, 0, thisTime.Unix())

		defer database.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "", http.StatusBadRequest)
	}
	
}

func handleSetComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { 
		keys, ok := r.URL.Query()["id"]
		fmt.Println("Id = ", keys[0], ok)

		database := dbConnect()

		id, _ := strconv.Atoi(keys[0])
		statement, _ := database.Prepare("UPDATE todo SET status = ? WHERE id = ?")
		statement.Exec(1, id)

		defer database.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else {
		http.Error(w, "", http.StatusBadRequest)
	}
}

func handleFormEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { 
		keys, _ := r.URL.Query()["id"]
		
		database := dbConnect()

		rows, _ := database.Query("SELECT name, assignee, deadline FROM todo WHERE id = " + keys[0]) 

		var name string
		var assignee string
		var deadline int

		for rows.Next() {
			rows.Scan(&name, &assignee, &deadline)

		}

		dataDate := time.Unix(int64(deadline), 0)	
		
		year, month, day := dataDate.Date()
		
		getDay := strconv.Itoa(day)
		if day < 10 {
			getDay = string('0') + getDay // Ingat, char adalah rune di Go
		}
			
		getMonth := strconv.Itoa(int(month))
		if int(month) < 10 {
			getMonth = string('0') + getMonth // Ingat, char adalah rune di Go
		}
		getYear := strconv.Itoa(year)
		strDate := getDay + "/" + getMonth + "/" + getYear
		
		fmt.Println(reflect.TypeOf(dataDate))
		
		filepath := path.Join("views","edit_todo.html")
		tmpl, _ := template.ParseFiles(filepath)

		var data = map[string]interface{} {
			"id":keys[0],
			"name":name,
			"assignee":assignee,
			"deadline":strDate,
		}
		tmpl.Execute(w, data)

		defer database.Close()

	} else {
		http.Error(w, "", http.StatusBadRequest)
	}
	
}

func handleProcessEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		database := dbConnect()

		id_input := r.FormValue("edit_id")
		name_input := r.FormValue("edit_task") 
		assignee_input := r.FormValue("edit_assignee")
		date_input := r.FormValue("edit_date")
		
		date_array := strings.Split(date_input, "/")

		date_array_date, err := strconv.Atoi(date_array[0])
		if err != nil {
			fmt.Println("Can not convert date")
		}
		date_array_month, err := strconv.Atoi(date_array[1])
		if err != nil {
			fmt.Println("Can not convert month")
		}
		date_array_year, err := strconv.Atoi(date_array[2])
		if err != nil {
			fmt.Println("Can not convert year")
		}
		thisTime := time.Date(date_array_year, time.Month(date_array_month), date_array_date, 0, 0, 0, 0, time.Local)
		
		statement, _ := database.Prepare("UPDATE todo SET name=?, assignee=?, deadline=? WHERE id=? ")
		statement.Exec(name_input, assignee_input, thisTime.Unix(), id_input)

		defer database.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "", http.StatusBadRequest)
	}
	
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" { 
		keys, _ := r.URL.Query()["id"]
		
		database := dbConnect()

		statement, _ := database.Prepare("DELETE FROM todo WHERE id=? ")
		statement.Exec(keys[0])
		
		defer database.Close()

		http.Redirect(w, r, "/", http.StatusSeeOther)

	} else {
		http.Error(w, "", http.StatusBadRequest)
	}
	
}

func main() {
	
	dbCreate()
	
	
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/add", handleFormEntry)
	http.HandleFunc("/process_add", handleProcessEntry)
	http.HandleFunc("/set_complete/", handleSetComplete)
	http.HandleFunc("/edit/", handleFormEdit)
	http.HandleFunc("/process_edit", handleProcessEdit) 
	http.HandleFunc("/delete/", handleDelete)
	
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("assets"))))

	fmt.Println("Server starting at port 80")
	http.ListenAndServe(":80", nil)
}