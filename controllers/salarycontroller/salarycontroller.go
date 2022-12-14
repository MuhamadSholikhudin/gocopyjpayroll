package mahasiswacontroller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"gocopyjpayroll/entities"
	"gocopyjpayroll/models/fileprocessmodel"
)

var fileprocessModel = fileprocessmodel.New()

// Salary
func SalaryReport(w http.ResponseWriter, r *http.Request) {
	type M map[string]interface{}

	var data = M{"name": "HRD"}
	var tmpl = template.Must(template.ParseFiles(
		"views/templates/_header.html",
		"views/templates/_navbar.html",
		"views/salary/salaryreport.html",
		"views/templates/_footer.html",
	))

	var err = tmpl.ExecuteTemplate(w, "salaryreport", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SalaryDownload(w http.ResponseWriter, r *http.Request) {
	periode := r.FormValue("periode")
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	periode = re.ReplaceAllString(periode, "")
	path := fmt.Sprintf("C:/go/gocopyjpayroll/files/salary/Payroll_Salary_Report_M_%s.xlsx", periode)
	f, err := os.Open(path)
	if f != nil {
		defer f.Close()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	time.Sleep(2 * time.Second)
	contentDisposition := fmt.Sprintf("attachment; filename=Payroll_Salary_Report_M_%s.xlsx", periode)
	w.Header().Set("Content-Disposition", contentDisposition)

	if _, err := io.Copy(w, f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func SalaryForm(w http.ResponseWriter, r *http.Request) {
	type K map[string]interface{}

	data := map[string]interface{}{
		"data": template.HTML(GetDataSalary()),
	}

	var tmpl = template.Must(template.ParseFiles(
		"views/templates/_header.html",
		"views/templates/_navbar.html",
		"views/salary/salaryform.html",
		"views/templates/_footer.html",
	))

	var err = tmpl.ExecuteTemplate(w, "salaryform", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetDataSalary() string {

	buffer := &bytes.Buffer{}

	temp, _ := template.New("data.html").Funcs(template.FuncMap{
		"increment": func(a, b int) int {
			return a + b
		},
	}).ParseFiles("views/salary/data.html")

	var fileprocess []entities.Fileprocess
	err := fileprocessModel.FindAll(&fileprocess)
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"fileprocess": fileprocess,
	}

	temp.ExecuteTemplate(buffer, "data.html", data)

	return buffer.String()
}

func SalaryUpload(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	periode := r.FormValue("periode")
	category := r.FormValue("category")

	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	periode = re.ReplaceAllString(periode, "")

	fmt.Println(periode)
	fmt.Println(category)

	uploadedFile, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := handler.Filename
	if category != "" {
		filename = fmt.Sprintf("Payroll_%s_Report_M_%s%s", category, periode, filepath.Ext(handler.Filename))
	}

	var fileprocess entities.Fileprocess

	fileprocess.Periode = periode
	fileprocess.Category = category
	fileprocess.File = filename

	err = fileprocessModel.Create(&fileprocess)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	fileLocation := filepath.Join(dir, "files/salary", filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/jpayroll/salaryreport", 301)
}

func newFunction(fileprocess entities.Fileprocess) error {
	err := fileprocessModel.Create(&fileprocess)
	return err
}

func SalaryEdit(w http.ResponseWriter, r *http.Request) {

	queryString := r.URL.Query()
	id, err := strconv.ParseInt(queryString.Get("id"), 10, 64)

	var data map[string]interface{}
	var fileprocess entities.Fileprocess

	if err != nil {
		data = map[string]interface{}{
			"title":       "Tambah Data File Processs",
			"fileprocess": fileprocess,
		}
	} else {

		err := fileprocessModel.Find(id, &fileprocess)
		if err != nil {
			panic(err)
		}

		data = map[string]interface{}{
			"title":       "Edit Data fileprocess",
			"fileprocess": fileprocess,
		}
	}

	temp, _ := template.ParseFiles("views/salary/editform.html")
	temp.Execute(w, data)
}

func ResponseError(w http.ResponseWriter, code int, message string) {
	ResponseJson(w, code, map[string]string{"error": message})
}

func ResponseJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
