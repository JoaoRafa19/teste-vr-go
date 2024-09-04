package api

// Requests
type RequestCreateCurso struct {
	Theme       string `json:"ementa"`
	Description string `json:"descricao"`
}
type RequestUpdateCurso struct {
	Theme       string `json:"ementa"`
	Description string `json:"descricao"`
}

type RequestCreateAluno struct {
	Name string `json:"nome"`
}

type RequestUpdateAluno struct {
	Name string `json:"nome"`
}
type RequestMatricula struct {
	CourseCodes []int `json:"cursos"`
}

// Responses
type CodigoMatricula struct {
	Codigo int `json:"codigo_curso"`
	Matricula int `json:"matricula"`
}
type ResponseMatricula struct {
	Enrolments []CodigoMatricula `json:"matriculas"`
}

type ResponseMatriculaError struct {
	Message string `json:"message"`
}

type ResponseCurso struct {
	Code        int32  `json:"codigo"`
	Description string `json:"descricao"`
	Theme       string `json:"ementa"`
	Enrolments  int32  `json:"matriculas"`
}

type ResponseStudent struct {
	Code int    `json:"codigo"`
	Name string `json:"nome"`
}

type ResponseCreateAluno struct {
	Code int32 `json:"codigo"`
}

type ResponseUpdateCurso struct {
	Code        int32  `json:"codigo"`
	Description string `json:"descricao"`
	Theme       string `json:"ementa"`
}
type ResponseCreateCurso struct {
	Code        int32  `json:"codigo"`
	Description string `json:"descricao"`
	Theme       string `json:"ementa"`
}
