package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"

	"hackathon.com/pyz/dbm"
)

// Handler wrap an http.Handler with common helpers
type Handler struct {
	http.Handler
	DB dbm.DatabaseOperator
}

// ReadQueryParams is used to parse URL query params from request
func ReadQueryParams(r *http.Request, out interface{}) error {
	err := decoder.Decode(out, r.URL.Query())
	return err
}

// ReadJSON reads bytes from `r` and fill int `out` struct fields accordingly.
// used to parse http request body & initialize relevant param structs
func ReadJSON(r *http.Request, out interface{}) error {
	var contentType = r.Header.Get("Content-Type")
	var bytes []byte
	var err error

	switch {
	case strings.HasPrefix(contentType, "application/x-www-form-urlencoded"):
		/* Handle form params */
		if err = r.ParseForm(); err != nil {
			return err
		}

		if err = decoder.Decode(out, r.PostForm); err != nil {
			return err
		}

	case strings.HasPrefix(contentType, "multipart/form-data"):
		/* Handle form params */
		if err = r.ParseMultipartForm(32 << 20); err != nil {
			return err
		}

		if err = decoder.Decode(out, r.PostForm); err != nil {
			return err
		}

	default:
		if bytes, err = ioutil.ReadAll(r.Body); err != nil {
			return err
		}

		if len(bytes) <= 0 {
			return ErrEmptyRequestBody
		}

		if err = json.Unmarshal(bytes, &out); err != nil {
			return err
		}
	}

	if err = validate.Struct(out); err != nil {
		return err
	}

	return nil
}

// RollbackOnPanic makes sure database Tx is rollbacked on panic
func RollbackOnPanic(tx *sqlx.Tx) {
	if r := recover(); r != nil {
		if err := tx.Rollback(); err != nil {
			logger.Error("ErrTxRollback: ", err)
		}

		panic(r) // IMPORTANT keep panic to write JSON Error
	}
}

// WriteError write json response according to err
func WriteError(w http.ResponseWriter, err Error, locale string) {
	// log error for internal server errors (unexpected), otherwise warn
	if err.status == http.StatusInternalServerError {
		logger.Errorf("%s - %v", err.code, err)
	} else {
		logger.Warnf("%s - %v", err.code, err)
	}

	w.WriteHeader(err.status)
	WriteJSON(w, ErrorResponse{
		Success: false,
		Error:   err.code,
		Message: err.code,
	})
}

// WriteJSON generate JSON response in { "sucess": ..., } format
func WriteJSON(w http.ResponseWriter, response interface{}) {
	var err error
	var bytes []byte
	if bytes, err = json.Marshal(response); err != nil {
		logger.Error("ErrWriteJSON: ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

// SaveFile save file from multipart form format
func SaveFile(r *http.Request, dir string, uid int, prefix string, format string) (string, error) {
	var err error
	var file multipart.File
	if file, _, err = r.FormFile("file"); err != nil {
		return "", err
	}
	defer file.Close()

	// Create a temporary file within our temp-images directory that follows a particular naming pattern
	var tempFile *os.File
	if tempFile, err = ioutil.TempFile("assets", "upload-*.png"); err != nil {
		return "", err
	}
	// defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	// write this byte array to our temporary file
	tempFile.Write(fileBytes)

	// convert temp file to readable src file
	var srcFile *os.File
	if srcFile, err = os.Open(tempFile.Name()); err != nil {
		return "", err
	}
	defer srcFile.Close()

	// determine final output file
	var destFile *os.File
	var desFilename = "assets" + dir + fmt.Sprintf("%d", uid) + prefix + format
	if destFile, err = os.Create(desFilename); err != nil {
		return "", err
	}
	defer destFile.Close()

	// copy src file to output file
	if _, err = io.Copy(destFile, srcFile); err != nil {
		return "", err
	}

	// refresh final output file
	if err = destFile.Sync(); err != nil {
		return "", err
	}

	return desFilename, nil
}

// GetUID returns UID from request URL route parameter
func GetUID(r *http.Request) (int, error) {
	var uidParam string
	var err error
	var uid int

	if uidParam = chi.URLParam(r, "uid"); uidParam == "" {
		return 0, fmt.Errorf("the uid parameter cannot be empty")
	}

	if uid, err = strconv.Atoi(uidParam); err != nil {
		return 0, fmt.Errorf(uidParam + "is not a valid uid")
	}

	return uid, nil
}
