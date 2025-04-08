package home

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Home() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, `
			<!DOCTYPE html>
			<html>
			<head>
				<title>Swagger UI</title>
				<link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist/swagger-ui.css" />
			</head>
			<body>
				<div id="swagger-ui"></div>
				<script src="https://unpkg.com/swagger-ui-dist/swagger-ui-bundle.js"></script>
				<script>
					window.onload = () => {
						window.ui = SwaggerUIBundle({
							url: "/docs/openapi.yaml",
							dom_id: "#swagger-ui"
						});
					};
				</script>
			</body>
			</html>
		`)
	}
}
