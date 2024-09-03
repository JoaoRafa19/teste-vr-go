package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/JoaoRafa19/teste-vr-go/internal/store/pgstore"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type apiHandler struct {
	q *pgstore.Queries
	r *chi.Mux
}

// Implements http handler
func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.r.ServeHTTP(w, r)
}

func NewHandler(q *pgstore.Queries) apiHandler {
	h := apiHandler{
		q: q,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	//#DEFINE routes
	r.Post("/echo", h.echo)
	// #

	r.Route("/aluno", func(r chi.Router) {
		r.Get("/", h.handleGetAllStudents)
	})

	h.r = r
	return h

}

func returnError(w http.ResponseWriter, status int) {

	type _Message struct {
		Error string `json:"error"`
	}
	var errorMessage _Message
	w.WriteHeader(status)

	errorMessage = _Message{
		Error: http.StatusText(status),
	}

	data, _ := json.Marshal(errorMessage)
	w.Write(data)
}

// func returnData(result []byte, w http.ResponseWriter) {
// 	w.Header().Set("Content-Type", "application/json")
// 	if _, err := w.Write(result); err != nil {
// 		slog.Error("failed to return response room", "error", err)
// 	}
// }

func (h *apiHandler) handleGetAllStudents(w http.ResponseWriter, r *http.Request) {
	students, err := h.q.GetAllStudents(r.Context())

	if err != nil {
		w.WriteHeader(500)
	}

	w.WriteHeader(200)
	type resType struct {
		Codigo int    `json:"codigo"`
		Nome   string `json:"nome"`
	}
	var data []resType
	for _, student := range students {
		data = append(data, resType{
			Codigo: int(student.Codigo),
			Nome:   student.Nome.String,
		})
	}

	//var marshaledData []byte
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(500)
	}

	//returnData(students)

}

func (h *apiHandler) echo(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Message string `json:"message"`
	}

	var b body

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		slog.Error("Unmarshal", "error", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(b.Message))
}
