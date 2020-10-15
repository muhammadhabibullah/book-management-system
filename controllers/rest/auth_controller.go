package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"book-management-system/controllers/rest/requests"
)

// AuthController will handle authentication & authorization domain requests
type AuthController struct{}

// NewAuthController returns new AuthController
func NewAuthController(
	route *mux.Router,
) *AuthController {
	ctrl := &AuthController{}

	v1Route := route.PathPrefix("/v1").Subrouter()

	v1AuthRoute := v1Route.PathPrefix("/auth").Subrouter()

	v1AuthRoute.HandleFunc("", ctrl.Login).Methods(http.MethodPost)

	return ctrl
}

// Login handle login request
// @Summary Login
// @Description Login
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body requests.Login true "Request Body"
// @Success 200 {object} models.Book "OK"
// @Failure 400 {object} responses.ErrorResponse "Bad Request"
// @Failure 403 {object} responses.ErrorResponse "Forbidden"
// @Failure 500 {object} responses.ErrorResponse "Internal Server Error"
// @Router /v1/book [post]
func (ctrl *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest requests.Login
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginRequest); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
}
