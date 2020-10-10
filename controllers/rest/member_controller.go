package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"book-management-system/entities/models"
	"book-management-system/usecases"
	"book-management-system/usecases/services"
)

// MemberController will handle member domain requests
type MemberController struct {
	memberService services.MemberService
}

// NewMemberController returns new MemberController
func NewMemberController(route *mux.Router, useCase *usecases.UseCase) *MemberController {
	ctrl := &MemberController{
		memberService: useCase.Service.MemberService,
	}

	v1Route := route.PathPrefix("/v1").Subrouter()

	v1MemberRoute := v1Route.PathPrefix("/member").Subrouter()
	v1MemberRoute.HandleFunc("", ctrl.CreateMember).Methods(http.MethodPost)
	v1MemberRoute.HandleFunc("", ctrl.GetMembers).Methods(http.MethodGet)
	v1MemberRoute.HandleFunc("", ctrl.UpdateMember).Methods(http.MethodPut)

	return ctrl
}

// CreateMember handle create member request
// @Summary Create a new member
// @Description Create a new member
// @Tags Member
// @Accept json
// @Produce json
// @Param request body models.Member true "Request Body"
// @Success 201 {object} models.Member "Created"
// @Failure 400 {object} responses.ErrorResponse "Bad Request"
// @Failure 500 {object} responses.ErrorResponse "Internal Server Error"
// @Router /v1/member [post]
func (ctrl *MemberController) CreateMember(w http.ResponseWriter, r *http.Request) {
	var member models.Member
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&member); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := ctrl.memberService.CreateMember(r.Context(), &member); err != nil {
		respondWithError(w, http.StatusInternalServerError,
			fmt.Sprintf("Failed create member: %s", err.Error()))
		return
	}

	respondWithJSON(w, http.StatusCreated, member)
}

// GetMembers handle get all members request
// @Summary Get all members
// @Description Get all members
// @Tags Member
// @Accept json
// @Produce json
// @Success 200 {object} models.Members "OK"
// @Failure 500 {object} responses.ErrorResponse "Internal Server Error"
// @Router /v1/member [get]
func (ctrl *MemberController) GetMembers(w http.ResponseWriter, r *http.Request) {
	members, err := ctrl.memberService.GetMembers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			fmt.Sprintf("Failed get members: %s", err.Error()))
		return
	}

	respondWithJSON(w, http.StatusOK, members)
}

// UpdateMember handle update member request
// @Summary Update a member
// @Description Update a member
// @Tags Member
// @Accept json
// @Produce json
// @Param request body models.Member true "Request Body"
// @Success 200 {object} models.Member "Updated"
// @Failure 400 {object} responses.ErrorResponse "Bad Request"
// @Failure 500 {object} responses.ErrorResponse "Internal Server Error"
// @Router /v1/member [put]
func (ctrl *MemberController) UpdateMember(w http.ResponseWriter, r *http.Request) {
	var member models.Member
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&member); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := ctrl.memberService.UpdateMember(r.Context(), &member); err != nil {
		respondWithError(w, http.StatusInternalServerError,
			fmt.Sprintf("Failed update member: %s", err.Error()))
		return
	}

	respondWithJSON(w, http.StatusOK, member)
}
