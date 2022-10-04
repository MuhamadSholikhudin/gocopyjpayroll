package main

import (
	"fmt"
	"net/http"

	salarycontroller "github.com/jeypc/go-crud-modal/controllers/salarycontroller"
)

func main() {

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/salaryindex", salarycontroller.SalaryIndex)
	http.HandleFunc("/", salarycontroller.Index)
	http.HandleFunc("/salary/get_form", salarycontroller.GetForm)
	http.HandleFunc("/salary/store", salarycontroller.Store)
	http.HandleFunc("/salary/delete", salarycontroller.Delete)

	http.HandleFunc("/jpayroll/salary", salarycontroller.Salary)
	http.HandleFunc("/jpayroll/download", salarycontroller.DownloadSalary)

	fmt.Println("Listen on port 10.10.42.6:3000")
	http.ListenAndServe(":3000", nil)

}
