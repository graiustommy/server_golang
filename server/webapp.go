package main
import(
"fmt"
"net/http"
"html/template"
"database/sql"
"log"
_ "github.com/lib/pq"
"github.com/gorilla/mux"
)

type Product struct{
    Id int
    Model string
    Company string
    Price int
}

var database *sql.DB

func DeleteHandler(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    id := vars["id"]
    _, err := database.Exec("delete from Products where id=$1", id)
    if err != nil{ log.Println(err) }
    http.Redirect(w, r, "/", 301)
}

func EditPage(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    id := vars["id"]
    row := database.QueryRow("select * from Products where id=$1", id)
    prod := Product{}
    err := row.Scan(&prod.Id, &prod.Model, &prod.Company, &prod.Price)
    if err != nil{
        log.Println(err)
        http.Error(w, http.StatusText(404), http.StatusNotFound)
    }else{
        tmpl, _ := template.ParseFiles("templates/edit.html")
        tmpl.Execute(w, prod)
    }
}

func EditHandler( w http.ResponseWriter, r *http.Request){
    err := r.ParseForm()
    if err != nil{
        log.Println(err)
    }
    id := r.FormValue("id")
    company := r.FormValue("company")
    model := r.FormValue("model")
    price := r.FormValue("price")

    _, err = database.Exec("update Products set model=$1, company=$2, price=$3 where id = $4", model,company,price,id)
    if err != nil{
        log.Println(err)
    }
    http.Redirect(w,r,"/",301)
}

func CreateHandler(w http.ResponseWriter, r *http.Request){
    if r.Method == "POST"{
        err := r.ParseForm()
        if err != nil{
            log.Println(err)
        }
        model := r.FormValue("model")
        company := r.FormValue("company")
        price := r.FormValue("price")
        _, err = database.Exec("insert into Products (model, company, price) values ($1, $2, $3)", model, company, price)
        if err != nil{
            log.Println(err)
        }
        http.Redirect(w, r, "/", 301)
    }else{
        http.ServeFile(w,r, "templates/create.html")
    }
}

func IndexHandler(w http.ResponseWriter, r *http.Request){
    rows, err := database.Query("select * from Products")
    if err != nil{ log.Println(err)}
    defer rows.Close()
    products := []Product{}
    for rows.Next(){
        p := Product{}
        err := rows.Scan(&p.Id, &p.Model, &p.Company, &p.Price)
        if err != nil{fmt.Println(err)
            continue}
        products = append(products, p)
    } 
    tmpl, _ := template.ParseFiles("templates/index.html")
    tmpl.Execute(w, products)
}

func main(){
    connStr := "user=postgres password=a12l40a15a44 dbname=productsdb sslmode=disable"    
    db, err := sql.Open("postgres", connStr)
    if err != nil{ log.Println(err)}
    defer db.Close()
    database = db

    router := mux.NewRouter()
    router.HandleFunc("/", IndexHandler)
    router.HandleFunc("/create", CreateHandler)
    router.HandleFunc("/edit/{id:[0-9]+}", EditPage).Methods("GET")
    router.HandleFunc("/edit/{id:[0-9]+}", EditHandler).Methods("POST")
    router.HandleFunc("/delete/{id:[0-9]+}", DeleteHandler)

    http.Handle("/", router)

    fmt.Println("Server is listening...")
    http.ListenAndServe(":3000", nil)
}

