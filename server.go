package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("uploading file...")
	r.ParseMultipartForm(10 << 20) //file size 10 mb

	//ensure ki sirf POST request ko handle kare
	if r.Method != http.MethodPost {
		http.Error(w, "sirf POST method allowed hai", http.StatusMethodNotAllowed)
		return
	}

	//ek file handler initiate karenge
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error retrieving the file : ", err)
		return
	}

	defer file.Close()

	//ensure that directory exists
	os.MkdirAll("./ctr/", os.ModePerm)

	//saving the file to ./ctr/ directory
	dst, err := os.Create("./ctr/" + handler.Filename)
	if err != nil {
		fmt.Println("Error creating file ", err)
		return
	}
	defer dst.Close()

	//now copying the file
	_, err = io.Copy(dst, file)
	if err != nil {
		fmt.Println("Error saving the file ", err)
		return
	}
	fmt.Fprintf(w, "%s uploaded of size %d ", handler.Filename, handler.Size)

}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	//pehle GET request confirm karenge
	if r.Method != http.MethodGet {
		http.Error(w, "Bhai sirf GET method allowed hai", http.StatusMethodNotAllowed)
		return
	}

	//ek query initiate karenge yahan pe
	file := r.URL.Query().Get("file")
	if file == "" {
		http.Error(w, "file parameter missing hai", http.StatusBadRequest)
		return
	}

	path := "./ctr/" + file
	f, err := os.Open(path)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Disposition", "attachement;filename="+file)
	w.Header().Set("Content-Type", "application/octet-stream")

	io.Copy(w, f)
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	// isme GET ya POST req kuch bhi chalega to check karne ki jarurat nhi hai

	//again ek query initiate karenge
	file := r.URL.Query().Get("file")

	if file == "" {
		http.Error(w, "file parameter wala error hai", http.StatusBadRequest)
		return
	}

	err := os.Remove("./ctr/" + file)
	if err != nil {
		http.Error(w, "Error while deleting file", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s deleted successfully", file)
}

func listAll(w http.ResponseWriter, r *http.Request) {
	file, err := os.ReadDir("./ctr/")
	if err != nil {
		http.Error(w, "Error while reading the folder", http.StatusInternalServerError)
		return
	}

	//files downlaod krne ke liye ek hyperlink embed kare denge sare files me taki file pe click krte hi download ho jaaye
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintln(w, "<h1> availabe files for campus tv & radio club")

	//ab ek loop banake sare files print kar denge
	for _, f := range file {
		name := f.Name()
		fmt.Fprintf(w, `<li><a href="/download?file=%s">%s</a></li>`, name, name)
	}
	fmt.Fprintln(w, "</ul>")
}

func main() {
	http.HandleFunc("/upload", uploadFile)
	http.HandleFunc("/download", downloadFile)
	http.HandleFunc("/delete", deleteFile)
	http.HandleFunc("/", listAll)
	fmt.Println("running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
