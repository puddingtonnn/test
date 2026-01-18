package handler

func (h *Handler) RegisterRoutes() {
	h.router.HandleFunc("/", h.handleRoot)
}