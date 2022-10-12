package alterationcontroller

import (
	"fmt"
	"gocopyjpayroll/models/fileprocessmodel"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"text/template"
	"time"
)

var fileprocessModel = fileprocessmodel.New()

type M map[string]interface{}

func AlterationReport(w http.ResponseWriter, r *http.Request) {

	var data = M{"name": "HRD"}

	var tmpl = template.Must(template.ParseFiles(
		"views/templates/_header.html",
		"views/templates/_navbar.html",
		"views/alteration/alterationreport.html",
		"views/templates/_footer.html",
	))

	var err = tmpl.ExecuteTemplate(w, "alterationreport", data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func AlterationDownload(w http.ResponseWriter, r *http.Request) {
	periode := r.FormValue("periode")
	re, err := regexp.Compile(`[^\w]`)
	if err != nil {
		log.Fatal(err)
	}
	periode = re.ReplaceAllString(periode, "")
	path := fmt.Sprintf("C:/go/gocopyjpayroll/files/alteration/Employee_Report_Alteration_%s.xlsx", periode)
	f, err := os.Open(path)
	if f != nil {
		defer f.Close()
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	time.Sleep(2 * time.Second)
	contentDisposition := fmt.Sprintf("attachment; filename=Employee_Report_Alteration_%s.xlsx", periode)
	w.Header().Set("Content-Disposition", contentDisposition)

	if _, err := io.Copy(w, f); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
