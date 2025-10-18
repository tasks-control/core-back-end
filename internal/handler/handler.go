package handler

import (
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
	v1 "github.com/tasks-control/core-back-end/api/v1"
)

type Handler struct {
	Service any
}

func (h *Handler) GetAlive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PostAuthRefresh(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PostAuthRegister(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) GetBoards(w http.ResponseWriter, r *http.Request, params v1.GetBoardsParams) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PostBoards(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) DeleteBoardsIdBoard(w http.ResponseWriter, r *http.Request, idBoard openapi_types.UUID) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) GetBoardsIdBoard(w http.ResponseWriter, r *http.Request, idBoard openapi_types.UUID) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PutBoardsIdBoard(w http.ResponseWriter, r *http.Request, idBoard openapi_types.UUID) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) DeleteBoardsIdBoardMembersIdMember(w http.ResponseWriter, r *http.Request, idBoard openapi_types.UUID, idMember openapi_types.UUID) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PostCards(w http.ResponseWriter, r *http.Request, params v1.PostCardsParams) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) DeleteCardsIdCard(w http.ResponseWriter, r *http.Request, idCard openapi_types.UUID) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) GetCardsIdCard(w http.ResponseWriter, r *http.Request, idCard openapi_types.UUID) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PutCardsIdCard(w http.ResponseWriter, r *http.Request, idCard openapi_types.UUID) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PostLists(w http.ResponseWriter, r *http.Request, params v1.PostListsParams) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) DeleteListsIdList(w http.ResponseWriter, r *http.Request, idList openapi_types.UUID) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) GetListsIdList(w http.ResponseWriter, r *http.Request, idList openapi_types.UUID) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PutListsIdList(w http.ResponseWriter, r *http.Request, idList openapi_types.UUID) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) DeleteMembersBoardsIdBoardStar(w http.ResponseWriter, r *http.Request, idBoard openapi_types.UUID) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PostMembersBoardsIdBoardStar(w http.ResponseWriter, r *http.Request, idBoard openapi_types.UUID) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PostMembersBoardsNameBoardUniqueJoin(w http.ResponseWriter, r *http.Request, nameBoardUnique string) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) GetMembersMe(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PutMembersMe(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func NewHandler(s any) *Handler {
	return &Handler{Service: s}
}
