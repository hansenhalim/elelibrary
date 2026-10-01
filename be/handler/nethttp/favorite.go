package nethttp

import "net/http"

func (h *Handler) AddFavorite(w http.ResponseWriter, r *http.Request) {
	if err := h.addFavorite.Execute(r.Context(), r.PathValue("bookId")); err != nil {
		writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	if err := h.removeFavorite.Execute(r.Context(), r.PathValue("bookId")); err != nil {
		writeError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
