package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	custommiddleware "github.com/wonjinsin/simple-chatbot/internal/handler/http/middleware"
	"github.com/wonjinsin/simple-chatbot/internal/usecase"
)

// NewRouter creates and configures a new chi router
func NewRouter(
	inquirySvc usecase.InquiryService,
) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(custommiddleware.TrID())
	r.Use(custommiddleware.CORS())
	r.Use(middleware.RealIP)
	r.Use(custommiddleware.HTTPLogger())
	r.Use(middleware.Recoverer)

	// Controllers
	healthCtrl := NewHealthController()
	inquiryCtrl := NewInquiryController(inquirySvc)

	// Routes
	r.Get("/healthz", healthCtrl.Check)

	// Inquiry routes
	r.Route("/inquiry", func(r chi.Router) {
		r.Post("/ask", inquiryCtrl.Ask)
		r.Post("/embed/origins", inquiryCtrl.EmbedInquiryOrigins)
	})

	return r
}
