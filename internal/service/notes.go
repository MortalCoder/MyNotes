package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type noteReq struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type noteResp struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (s *Service) CreateNote(c echo.Context) error {
	uid, err := s.userIDFromToken(c)
	if err != nil {
		return s.unauthorized(c)
	}

	var req noteReq
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil || req.Title == "" || req.Body == "" {
		return c.JSON(http.StatusBadRequest, &Response{ErrorMessage: InvalidParams})
	}

	id, err := s.notesRepo.Create(c.Request().Context(), uid, req.Title, req.Body)
	if err != nil {
		s.logger.Errorf("notes.Create: %v", err)
		return c.JSON(http.StatusInternalServerError, &Response{ErrorMessage: InternalServerError})
	}

	return c.JSON(http.StatusCreated, map[string]any{"id": id})
}

func (s *Service) GetNoteByID(c echo.Context) error {
	uid, err := s.userIDFromToken(c)
	if err != nil {
		return s.unauthorized(c)
	}
	nid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &Response{ErrorMessage: InvalidParams})
	}

	n, err := s.notesRepo.Get(c.Request().Context(), uid, nid)
	if err != nil {
		s.logger.Errorf("notes.Get: %v", err)
		return c.JSON(http.StatusInternalServerError, &Response{ErrorMessage: InternalServerError})
	}
	if n == nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, noteResp{ID: n.ID, Title: n.Title, Body: n.Body})
}

func (s *Service) UpdateNote(c echo.Context) error {
	uid, err := s.userIDFromToken(c)
	if err != nil {
		return s.unauthorized(c)
	}
	nid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &Response{ErrorMessage: InvalidParams})
	}

	var req noteReq
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil || req.Title == "" || req.Body == "" {
		return c.JSON(http.StatusBadRequest, &Response{ErrorMessage: InvalidParams})
	}

	aff, err := s.notesRepo.Update(c.Request().Context(), uid, nid, req.Title, req.Body)
	if err != nil {
		s.logger.Errorf("notes.Update: %v", err)
		return c.JSON(http.StatusInternalServerError, &Response{ErrorMessage: InternalServerError})
	}
	if aff == 0 {
		return c.NoContent(http.StatusNotFound)
	}
	return c.NoContent(http.StatusNoContent)
}

func (s *Service) DeleteNote(c echo.Context) error {
	uid, err := s.userIDFromToken(c)
	if err != nil {
		return s.unauthorized(c)
	}
	nid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &Response{ErrorMessage: InvalidParams})
	}

	aff, err := s.notesRepo.Delete(c.Request().Context(), uid, nid)
	if err != nil {
		s.logger.Errorf("notes.Delete: %v", err)
		return c.JSON(http.StatusInternalServerError, &Response{ErrorMessage: InternalServerError})
	}
	if aff == 0 {
		return c.NoContent(http.StatusNotFound)
	}
	return c.NoContent(http.StatusNoContent)
}

func (s *Service) ListNotes(c echo.Context) error {
	uid, err := s.userIDFromToken(c)
	if err != nil {
		return s.unauthorized(c)
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	list, err := s.notesRepo.List(c.Request().Context(), uid, limit, offset)
	if err != nil {
		s.logger.Errorf("notes.List: %v", err)
		return c.JSON(http.StatusInternalServerError, &Response{ErrorMessage: InternalServerError})
	}

	out := make([]noteResp, 0, len(list))
	for _, n := range list {
		out = append(out, noteResp{ID: n.ID, Title: n.Title, Body: n.Body})
	}
	return c.JSON(http.StatusOK, out)
}
