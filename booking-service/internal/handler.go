package internal

type Handler struct {
	srv ServiceInterface
}

func NewHandler(srv ServiceInterface) *Handler {
	return &Handler{
		srv: srv,
	}
}