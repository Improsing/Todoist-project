package handlers

import (
	"net/http"
	"time"

	"github.com/Improsing/go-final-project/utils"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	now, err := time.Parse("20060102", nowParam)
	if err != nil {
		http.Error(w, "Недопустимый параметр now", http.StatusBadRequest)
		return
	}
	nextDate, err := utils.NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} 
	w.Write([]byte(nextDate))
}