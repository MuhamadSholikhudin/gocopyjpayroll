package mahasiswacontroller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jeypc/go-crud-modal/entities"
	"github.com/jeypc/go-crud-modal/models/mahasiswamodel"
)

var mahasiswaModel = mahasiswamodel.New()

func Index(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"data": template.HTML(GetData()),
	}

	temp, _ := template.ParseFiles("views/salary/index.html")
	temp.Execute(w, data)
}

func SalaryIndex(w http.ResponseWriter, r *http.Request) {

	type M map[string]interface{}

	var tmpl, err = template.ParseGlob("views/salary/template/*")
	if err != nil {
		panic(err.Error())
		return
	}
	var data = M{"name": "Batman"}
	err = tmpl.ExecuteTemplate(w, "salaryreport", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Salary(w http.ResponseWriter, r *http.Request) {

	data := map[string]interface{}{
		"data": template.HTML(GetData()),
	}

	temp, _ := template.ParseFiles("views/salary/jpayroll.html")
	temp.Execute(w, data)
}

func DownloadSalary(w http.ResponseWriter, r *http.Request) {

	periode := r.FormValue("periode")
	path := fmt.Sprintf("C:/go/gocopyjpayroll/files/%s.xlsx", periode)
	f, err := os.Open(path)
	if f != nil {
		defer f.Close()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	time.Sleep(2 * time.Second)
	contentDisposition := fmt.Sprintf("attachment; filename=%s.xlsx", periode)
	w.Header().Set("Content-Disposition", contentDisposition)

	if _, err := io.Copy(w, f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetData() string {

	buffer := &bytes.Buffer{}

	temp, _ := template.New("data.html").Funcs(template.FuncMap{
		"increment": func(a, b int) int {
			return a + b
		},
	}).ParseFiles("views/salary/data.html")

	var mahasiswa []entities.Mahasiswa
	err := mahasiswaModel.FindAll(&mahasiswa)
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"mahasiswa": mahasiswa,
	}

	temp.ExecuteTemplate(buffer, "data.html", data)

	return buffer.String()
}

func GetForm(w http.ResponseWriter, r *http.Request) {

	queryString := r.URL.Query()
	id, err := strconv.ParseInt(queryString.Get("id"), 10, 64)

	var data map[string]interface{}
	var mahasiswa entities.Mahasiswa

	if err != nil {
		data = map[string]interface{}{
			"title":     "Tambah Data Mahasiswa",
			"mahasiswa": mahasiswa,
		}
	} else {

		err := mahasiswaModel.Find(id, &mahasiswa)
		if err != nil {
			panic(err)
		}

		data = map[string]interface{}{
			"title":     "Edit Data Mahasiswa",
			"mahasiswa": mahasiswa,
		}
	}

	temp, _ := template.ParseFiles("views/salary/form.html")
	temp.Execute(w, data)
}

func Store(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		r.ParseForm()
		var mahasiswa entities.Mahasiswa

		mahasiswa.NamaLengkap = r.Form.Get("nama_lengkap")
		mahasiswa.JenisKelamin = r.Form.Get("jenis_kelamin")
		mahasiswa.TanggalLahir = r.Form.Get("tanggal_lahir")
		mahasiswa.TempatLahir = r.Form.Get("tempat_lahir")
		mahasiswa.Alamat = r.Form.Get("alamat")

		id, err := strconv.ParseInt(r.Form.Get("id"), 10, 64)

		var data map[string]interface{}

		if err != nil {
			// insert data
			err := mahasiswaModel.Create(&mahasiswa)
			if err != nil {
				ResponseError(w, http.StatusInternalServerError, err.Error())
				return
			}
			data = map[string]interface{}{
				"message": "Data berhasil disimpan",
				"data":    template.HTML(GetData()),
			}
		} else {
			// mengupdate data
			mahasiswa.Id = id
			err := mahasiswaModel.Update(mahasiswa)
			if err != nil {
				ResponseError(w, http.StatusInternalServerError, err.Error())
				return
			}
			data = map[string]interface{}{
				"message": "Data berhasil diubah",
				"data":    template.HTML(GetData()),
			}
		}

		ResponseJson(w, http.StatusOK, data)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	id, err := strconv.ParseInt(r.Form.Get("id"), 10, 64)
	if err != nil {
		panic(err)
	}
	err = mahasiswaModel.Delete(id)
	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{
		"message": "Data berhasil dihapus",
		"data":    template.HTML(GetData()),
	}
	ResponseJson(w, http.StatusOK, data)
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
