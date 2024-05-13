package main

import "net/http"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"order": "136987979637580771203159149270", "status": "PROCESSED", "accrual": 500}`))
		if err != nil {
			panic(err)
		}

		w.WriteHeader(http.StatusOK)
	})

	err := http.ListenAndServe("localhost:8081", nil)
	if err != nil {
		panic(err)
	}
}
