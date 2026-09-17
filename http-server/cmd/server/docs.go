package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *httpServer) docs(r *mux.Router) {
	/* Swagger json */
	r.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs/"))))
	r.HandleFunc("/doc.api/swagger/json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger/swagger_v1.json")
	})

	r.HandleFunc("/doc.api/swagger/gui", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		html := `<!DOCTYPE html>
				<html lang="en">
				<head>
					<meta charset="utf-8" />
					<meta name="viewport" content="width=device-width, initial-scale=1" />
					<meta name="description" content="SwaggerUI" />
					<title>NeoSync API</title>
					<link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css" />
				</head>
				<body>
					<div id="swagger-ui"></div>
					<script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js" crossorigin></script>
					<script>
						window.onload = () => {
							window.ui = SwaggerUIBundle({
								url: '/api/v1/doc.api/swagger/json',
								dom_id: '#swagger-ui',
								requestInterceptor: (req) => {
									req.credentials = 'include';
									return req;
								}
							});
						};
					</script>
				</body>
				</html>`
		_, _ = w.Write([]byte(html))
	})
}
