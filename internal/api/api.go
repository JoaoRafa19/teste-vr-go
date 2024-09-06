package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/JoaoRafa19/teste-vr-go/internal/store/pgstore"
	"github.com/JoaoRafa19/teste-vr-go/internal/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5"
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

	r.Post("/echo", h.echo)
	r.Get("/dashboard", h.handleGetDashboardInfo)
	r.Route("/aluno", func(r chi.Router) {
		r.Get("/", h.handleGetAllStudents)
		r.Get("/{codigo}", h.handleGetStudent)
		r.Post("/", h.handleCreateAluno)
		r.Patch("/{codigo}", h.handleUpdateStudent)
		r.Post("/{codigo}/matricula", h.handleMatricula)
		r.Delete("/{codigo}", h.handleDeleteAluno)
		r.Get("/search", h.handleSearchAlunos)
	})
	r.Route("/curso", func(r chi.Router) {
		r.Post("/", h.handleCreateCurso)
		r.Get("/", h.handleGetAllCursos)
		r.Get("/{codigo}", h.handleGetCurso)
		r.Patch("/{codigo}", h.handleUpdateCurso)
		r.Delete("/{codigo}", h.handleDeleteCurso)
	})

	h.r = r
	return h
}

func (h *apiHandler) handleGetDashboardInfo(w http.ResponseWriter, r *http.Request) {

	dashboardInfo, err := h.q.GetDashBoardInfo(r.Context())
	if err != nil {
		fmt.Printf("Failed to get dashboard info: %v\n", err)
		returnError(w, 500)
		return
	}

	convertToMatriculas := func(data interface{}) ([]MatriculaPorCurso, error) {
		var result []MatriculaPorCurso
		if items, ok := data.([]interface{}); ok {
			for _, item := range items {
				if m, ok := item.(map[string]interface{}); ok {
					curso := m["curso"].(string)
					totalMatriculas := int64(m["total_matriculas"].(float64))
					result = append(result, MatriculaPorCurso{
						Curso:           curso,
						TotalMatriculas: totalMatriculas,
					})
				}
			}
		} else {
			return nil, fmt.Errorf("invalid type for matriculas_por_curso")
		}
		return result, nil
	}

	convertToAlunos := func(data interface{}) ([]Aluno, error) {
		var result []Aluno
		if items, ok := data.([]interface{}); ok {
			for _, item := range items {
				if m, ok := item.(map[string]interface{}); ok {
					nome := m["nome"].(string)
					codigo := int64(m["codigo"].(float64))
					result = append(result, Aluno{
						Nome:   nome,
						Codigo: codigo,
					})
				}
			}
		} else {
			return nil, fmt.Errorf("invalid type for alunos")
		}
		return result, nil
	}

	matriculasPorCurso, err := convertToMatriculas(dashboardInfo.MatriculasPorCurso)
	if err != nil {
		fmt.Printf("MatriculasPorCurso is empty or invalid: %v\n", err)
	}

	alunosComMatricula, err := convertToAlunos(dashboardInfo.AlunosComMatricula)
	if err != nil {
		fmt.Printf("AlunosComMatricula is empty or invalid: %v\n", err)
	}

	alunosSemMatricula, err := convertToAlunos(dashboardInfo.AlunosSemMatricula)
	if err != nil {
		fmt.Printf("AlunosSemMatricula is empty or invalid: %v\n", err)
	}

	res := ResponseDashBoardInfoRow{
		TotalAlunos:        dashboardInfo.TotalAlunos.(int64),
		TotalCursos:        dashboardInfo.TotalCursos.(int64),
		TotalMatriculas:    dashboardInfo.TotalMatriculas.(int64),
		MatriculasPorCurso: matriculasPorCurso,
		AlunosComMatricula: alunosComMatricula,
		AlunosSemMatricula: alunosSemMatricula,
	}

	fmt.Println(res)
	returnData(w, res)
}

func (h *apiHandler) handleDeleteAluno(w http.ResponseWriter, r *http.Request) {
	codigo := chi.URLParam(r, "codigo")
	var codigoInt int32
	if _, err := fmt.Sscanf(codigo, "%d", &codigoInt); err != nil {
		returnError(w, http.StatusBadRequest)
		return
	}

	err := h.q.DeleteAluno(r.Context(), codigoInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			returnError(w, 404)
			return
		}
		returnError(w, 500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *apiHandler) handleDeleteCurso(w http.ResponseWriter, r *http.Request) {
	codigo := chi.URLParam(r, "codigo")
	var codigoInt int32
	if _, err := fmt.Sscanf(codigo, "%d", &codigoInt); err != nil {
		returnError(w, http.StatusBadRequest)
		return
	}

	err := h.q.DeleteCurso(r.Context(), codigoInt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			returnError(w, 404)
			return
		}
		returnError(w, 500)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *apiHandler) handleMatricula(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "codigo")
	if len(code) <= 0 {
		returnError(w, http.StatusBadRequest)
		return
	}

	var codigo int32
	if _, err := fmt.Sscanf(code, "%d", &codigo); err != nil {
		returnError(w, 500)
		return
	}

	var body RequestMatricula
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		returnError(w, http.StatusBadRequest)
		return
	}

	aluno, err := h.q.GetAluno(r.Context(), codigo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			returnError(w, 404)
			return
		}
		returnError(w, 500)
		return
	}

	matriculas, err := h.q.MatriculasPorAluno(r.Context(), aluno.Codigo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			returnError(w, 404)
			return
		}
		returnError(w, 500)
		return
	}
	if matriculas == 3 {
		//Um aluno não pode estar matriculado em mais de 3 cursos;
		w.WriteHeader(400)
		returnData(w, ResponseMatriculaError{
			Message: "Limite de matriculas atingido",
		})
		return
	}

	cursosMatriculado, err := h.q.CursosMatriculados(r.Context(), aluno.Codigo)
	if err != nil {
		slog.Error("Erro ao buscar os cursos do aluno", "error", err)
		return
	}

	var res ResponseMatricula

	for _, code := range body.CourseCodes {
		// Cursos não podem ter mais de 10 alunos matriculados (turma cheia);
		if utils.Contains(cursosMatriculado, int32(code)) {
			continue
		}

		curso, err := h.q.GetCurso(r.Context(), int32(code))
		if err != nil {
			slog.Error("Erro ao buscar curso", "error", err)
			continue
		}
		if curso.Matriculas < 10 {
			matricula, err := h.q.MatricularAluno(r.Context(), pgstore.MatricularAlunoParams{
				CodigoAluno: aluno.Codigo,
				CodigoCurso: curso.Codigo,
			})
			if err != nil {
				slog.Error("Erro ao matricular", "error", err)
				continue
			}
			res.Enrolments = append(res.Enrolments, CodigoMatricula{
				Codigo:    int(curso.Codigo),
				Matricula: int(matricula),
			})
		}
	}

	returnData(w, res)
}

func (h *apiHandler) handleUpdateCurso(w http.ResponseWriter, r *http.Request) {
	codigoCurso := chi.URLParam(r, "codigo")

	var codigo int32
	if _, err := fmt.Sscanf(codigoCurso, "%d", &codigo); err != nil {
		returnError(w, http.StatusBadRequest)
		return
	}
	_, err := h.q.GetCurso(r.Context(), codigo)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			returnError(w, 404)
			return
		}
		returnError(w, 500)
		return
	}

	var b RequestUpdateCurso
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "Invalid Request", 400)
		return
	}

	if len(b.Theme) == 0 || len(b.Description) == 0 {
		http.Error(w, "Invalid Request", 400)
		return
	}

	novoCurso, err := h.q.UpdateCurso(r.Context(), pgstore.UpdateCursoParams{
		Codigo:    codigo,
		Descricao: b.Description,
		Ementa:    b.Theme,
	})

	if err != nil {
		returnError(w, 500)
		return
	}
	var data = ResponseUpdateCurso{
		Code:        novoCurso.Codigo,
		Theme:       novoCurso.Ementa,
		Description: novoCurso.Descricao,
	}

	returnData(w, data)
}

func (h *apiHandler) handleUpdateStudent(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "codigo")
	if len(code) <= 0 {
		returnError(w, http.StatusBadRequest)
		return
	}

	var codigo int32
	if _, err := fmt.Sscanf(code, "%d", &codigo); err != nil {
		returnError(w, 500)
		return
	}

	var body RequestUpdateAluno
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		returnError(w, http.StatusBadRequest)
		return
	}

	_, err := h.q.GetAluno(r.Context(), codigo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			returnError(w, 404)
			return
		}
		returnError(w, 500)
		return
	}
	aluno, err := h.q.UpdateNomeAluno(r.Context(), pgstore.UpdateNomeAlunoParams{
		Codigo: codigo,
		Nome:   body.Name,
	})
	if err != nil {
		returnError(w, http.StatusInternalServerError)
		return
	}

	returnData(w, ResponseStudent{
		Code: int(aluno.Codigo),
		Name: aluno.Nome,
	})
}

func (h *apiHandler) handleGetStudent(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "codigo")
	if len(code) <= 0 {
		returnError(w, http.StatusBadRequest)
		return
	}

	var codigo int32
	if _, err := fmt.Sscanf(code, "%d", &codigo); err != nil {
		returnError(w, 500)
		return
	}

	student, err := h.q.GetAluno(r.Context(), codigo)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			returnError(w, 404)
			return
		}
		returnError(w, 500)
		return
	}

	returnData(w, ResponseStudent{
		Code: int(student.Codigo),
		Name: student.Nome,
	})
}

func (h *apiHandler) handleGetCurso(w http.ResponseWriter, r *http.Request) {
	codigoCurso := chi.URLParam(r, "codigo")

	var codigo int32
	if _, err := fmt.Sscanf(codigoCurso, "%d", &codigo); err != nil {
		returnError(w, http.StatusBadRequest)
		return
	}
	curso, err := h.q.GetCurso(r.Context(), codigo)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			returnError(w, 404)
			return
		}
		returnError(w, 500)
		return
	}

	var res = ResponseCurso{
		Code:        curso.Codigo,
		Description: curso.Descricao,
		Theme:       curso.Ementa,
		Enrolments:  curso.Matriculas,
	}

	returnData(w, res)
}

func (h *apiHandler) handleGetAllCursos(w http.ResponseWriter, r *http.Request) {
	cursos, err := h.q.GetCursos(r.Context())

	if err != nil {
		returnError(w, 500)
		return
	}

	var response []ResponseCurso

	for _, c := range cursos {
		response = append(response, ResponseCurso{
			Code:        c.Codigo,
			Description: c.Descricao,
			Theme:       c.Ementa,
			Enrolments:  c.Matriculas,
		})
	}

	returnData(w, response)
}

func (h *apiHandler) handleGetAllStudents(w http.ResponseWriter, r *http.Request) {
	students, err := h.q.GetAllAlunos(r.Context())

	if err != nil {
		returnError(w, 500)
		return
	}

	w.WriteHeader(200)

	var data []ResponseStudent
	for _, student := range students {
		data = append(data, ResponseStudent{
			Code: int(student.Codigo),
			Name: student.Nome,
		})
	}

	returnData(w, data)
}

func (h *apiHandler) handleCreateAluno(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	var b RequestCreateAluno
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "Invalid body", 400)
		return
	}
	if len(b.Name) == 0 {
		http.Error(w, "Invalid body", 400)
		return
	}

	codigo, err := h.q.CreateAluno(r.Context(), b.Name)

	if err != nil {
		http.Error(w, "Invalid body", 400)
		return
	}

	returnData(w, ResponseCreateAluno{
		Code: codigo,
	})
}

func (h *apiHandler) handleCreateCurso(w http.ResponseWriter, r *http.Request) {

	var b RequestCreateCurso
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "Invalid Request", 400)
		return
	}
	if len(b.Theme) == 0 || len(b.Description) == 0 {
		http.Error(w, "Invalid Request", 400)
		return
	}
	var curso pgstore.Curso
	curso, err := h.q.CreateCurso(r.Context(), pgstore.CreateCursoParams{
		Descricao: b.Description,
		Ementa:    b.Theme,
	})
	if err != nil {
		http.Error(w, "Internal Error", 500)
		return
	}

	var response = ResponseCreateCurso{
		Code:        curso.Codigo,
		Description: curso.Descricao,
		Theme:       curso.Ementa,
	}
	returnData(w, response)
}

func (h *apiHandler) echo(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Message string `json:"message"`
	}
	defer r.Body.Close()

	var b requestBody

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		slog.Error("Unmarshal", "error", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(b.Message))
}

func returnData(w http.ResponseWriter, res any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		returnError(w, 500)
		return
	}
}

func returnError(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
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

func (h *apiHandler) handleSearchAlunos(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	termosBusca := queryParams.Get("q")

	if termosBusca == "" {
		http.Error(w, "Missing search terms", http.StatusBadRequest)
		return
	}

	alunos, err := h.q.SearchAlunos(r.Context(), termosBusca)
	if err != nil {
		slog.Error("Erro ao buscar alunos", "error", err)
		returnError(w, http.StatusInternalServerError)
		return
	}

	var data []ResponseStudent
	for _, aluno := range alunos {
		data = append(data, ResponseStudent{
			Name: aluno,
		})
	}

	returnData(w, data)
}
