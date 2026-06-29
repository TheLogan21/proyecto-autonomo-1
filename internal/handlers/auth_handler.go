package handlers

import (
	"encoding/json"
	"net/http"

	"ebook-system/internal/middleware"
	"ebook-system/internal/services"
)

type AuthHandler struct {
	userService services.UserService
}

func NewAuthHandler(us services.UserService) *AuthHandler {
	return &AuthHandler{userService: us}
}

func (h *AuthHandler) LoginView(w http.ResponseWriter, r *http.Request) {
	render(w, r, "login", nil)
}

func (h *AuthHandler) LoginPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user, err := h.userService.Login(r.Context(), r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		render(w, r, "login", map[string]string{"Error": err.Error()})
		return
	}

	token := middleware.CreateSession(user.ID)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})
	http.Redirect(w, r, "/catalog", http.StatusSeeOther)
}

func (h *AuthHandler) RegisterView(w http.ResponseWriter, r *http.Request) {
	render(w, r, "register", nil)
}

func (h *AuthHandler) RegisterPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	err := h.userService.Register(r.Context(), r.FormValue("username"), r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		render(w, r, "register", map[string]string{"Error": err.Error()})
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("session"); err == nil {
		middleware.DestroySession(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// RegisterJSON procesa el registro de un nuevo usuario desde la API.
// Valida los datos requeridos y guarda el usuario mediante la interfaz del servicio.
func (h *AuthHandler) RegisterJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método no permitido"})
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decodificar el JSON de entrada
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cuerpo JSON inválido o malformado"})
		return
	}

	// Validación técnica defensiva de campos obligatorios
	if req.Username == "" || req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Todos los campos son requeridos"})
		return
	}

	// Invocar lógica de negocio en el servicio
	err := h.userService.Register(r.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuario registrado exitosamente"})
}

// LoginJSON valida las credenciales y genera un token de sesión.
// Soporta el retorno del token en la cabecera/JSON y la cookie para vistas tradicionales.
func (h *AuthHandler) LoginJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método no permitido"})
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decodificar JSON de credenciales
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cuerpo JSON inválido o malformado"})
		return
	}

	// Validación técnica defensiva
	if req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Email y contraseña obligatorios"})
		return
	}

	// Autenticar usuario con el servicio
	user, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Credenciales inválidas"})
		return
	}

	// Crear sesión y cookie correspondientes
	token := middleware.CreateSession(user.ID)
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
		"user":  user,
	})
}


func (h *AuthHandler) LogoutJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if cookie, err := r.Cookie("session"); err == nil {
		middleware.DestroySession(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Sesión cerrada exitosamente"})
}

