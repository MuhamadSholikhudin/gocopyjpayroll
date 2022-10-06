package main

import (
	"fmt"
	"net/http"

	alterationcontroller "gocopyjpayroll/controllers/alterationcontroller"
	attendancecontroller "gocopyjpayroll/controllers/attendancecontroller"
	customcontroller "gocopyjpayroll/controllers/customcontroller"
	salarycontroller "gocopyjpayroll/controllers/salarycontroller"
)

func main() {

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", salarycontroller.SalaryReport)

	//Salary
	http.HandleFunc("/jpayroll/salaryreport", salarycontroller.SalaryReport)
	http.HandleFunc("/jpayroll/salarydownload", salarycontroller.SalaryDownload)
	http.HandleFunc("/jpayroll/salaryform", salarycontroller.SalaryForm)
	http.HandleFunc("/jpayroll/salaryupload", salarycontroller.SalaryUpload)
	http.HandleFunc("/jpayroll/salaryedit", salarycontroller.SalaryEdit)
	// http.HandleFunc("/jpayroll/salarupdate", salarycontroller.SalaryUpdate)

	//Alteration
	http.HandleFunc("/jpayroll/alterationreport", alterationcontroller.AlterationReport)
	http.HandleFunc("/jpayroll/alterationdownload", alterationcontroller.AlterationDownload)
	// http.HandleFunc("/jpayroll/alterationform", alterationcontroller.AlterationForm)
	// http.HandleFunc("/jpayroll/alterationupload", alterationcontroller.AlterationUpload)
	// http.HandleFunc("/jpayroll/alterationedit", alterationcontroller.AlterationEdit)
	// http.HandleFunc("/jpayroll/alterationupdate", alterationcontroller.AlterationUpdate)

	//Attandance
	http.HandleFunc("/jpayroll/attendancereport", attendancecontroller.AttendanceReport)
	http.HandleFunc("/jpayroll/attendancedownload", attendancecontroller.AttendanceDownload)
	// http.HandleFunc("/jpayroll/attandanceform", attandancecontroller.AttandanceForm)
	// http.HandleFunc("/jpayroll/attandanceupload", attandancecontroller.AttandanceUpload)
	// http.HandleFunc("/jpayroll/attandanceedit", attandancecontroller.AttandanceEdit)
	// http.HandleFunc("/jpayroll/attandanceupdate", attandancecontroller.AttandanceUpdate)

	//Custom
	http.HandleFunc("/jpayroll/customreport", customcontroller.CustomReport)
	http.HandleFunc("/jpayroll/customdownload", customcontroller.CustomDownload)
	// http.HandleFunc("/jpayroll/customform", customcontroller.CustomForm)
	// http.HandleFunc("/jpayroll/customupload", customcontroller.CustomUpload)
	// http.HandleFunc("/jpayroll/customedit", customcontroller.CustomEdit)
	// http.HandleFunc("/jpayroll/customupdate", customcontroller.CustomUpdate)

	fmt.Println("Listen on port 10.10.42.6:3000")
	http.ListenAndServe(":3000", nil)

}
