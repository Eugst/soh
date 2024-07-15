package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
    // Check if the request method is POST
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Get the uploaded file
    file, _, err := r.FormFile("upload")
    if err != nil {
        http.Error(w, "Get Uploaded file: "+err.Error(), http.StatusBadRequest)
        return
    }
    defer file.Close()

    // Create a new file on the server
    dst, err := os.Create("/tmp/test.csv")
    if err != nil {
        http.Error(w, "Create new file: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer dst.Close()

    // Upload the file to S3
    // upload2s3(file, handler, w) # add handler to line 22


    // Copy the uploaded file to the destination file
    _, err = io.Copy(dst, file)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // fmt.Fprintf(w, "File uploaded successfully")
    http.Redirect(w, r, "/show", http.StatusFound)
}

var templates = template.Must(template.ParseFiles("/app/index.html", "/app/show.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
    if err := templates.ExecuteTemplate(w, tmpl+".html", data); err != nil {
        log.Fatalln("Unable to execute template.")
    }
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{"Title": "index"}
    renderTemplate(w, "index", data)
}

func ShowHandler(w http.ResponseWriter, r *http.Request) {
    file, err := os.Open("/tmp/test.csv")
    defer file.Close()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    reader := csv.NewReader(file)
    reader.Comma = ','
    reader.FieldsPerRecord = 3
    reader.TrimLeadingSpace = true
    csvLines, err := reader.ReadAll()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    fmt.Println(csvLines)

    data := map[string]interface{}{"Title": "File info", "CSV_LENGTH": len(csvLines)}
    renderTemplate(w, "show", data)
}

func main() {
    http.HandleFunc("/", IndexHandler)
    http.HandleFunc("/upload", uploadFile)
    http.HandleFunc("/show", ShowHandler)
    http.ListenAndServe(":8080", nil)
}

// func upload2s3(file multipart.File, handler *multipart.FileHeader, w http.ResponseWriter) {
    // uploader := s3manager.NewUploader(awsConnectRegion(region))
    // _, err = uploader.Upload(&s3manager.UploadInput{
    // 	Bucket: aws.String(bucket),
    // 	Key:    aws.String(handler.Filename),
    // 	Body:   file,
    // })
    // if err != nil {
    // 	if awsErr, ok := err.(awserr.Error); ok {
    // 		log.Println(awsErr)
    // 	} else {
    // 		log.Println(err)
    // 	}

    // 	w.WriteHeader(http.StatusInternalServerError)
    // 	return
    // }
// }

// Function to connect to AWS for S3 upload. Requires AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables
// func awsConnectRegion(region string) *session.Session {
//     if region == "" {
//         region = awsRegion
//     }

//     if val, ok := sessions[region]; ok {
//         return val
//     } else {
//         sess, err := session.NewSession(
//             &aws.Config{
//                 Region: aws.String(region),
//             },
//         )
//         if err != nil {
//             log.Println(err)
//         }
//         sessions[region] = sess
//         return sess
//     }
// }
