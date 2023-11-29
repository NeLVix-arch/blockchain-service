package middlewares

import (
	"net/http"

	"github.com/rs/cors"
)

// corsMiddleware - middleware для настройки корс
func CorsMiddleware(next http.Handler) http.Handler {
	// Создаем объект cors с желаемыми опциями
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost"}, // Разрешаем доступ только с локалхоста
		AllowedMethods: []string{"POST"},             // Разрешаем только пост запросы
		AllowedHeaders: []string{"Content-Type"},     // Разрешаем только указанные хедеры
	})

	// Возвращаем обернутый обработчик
	return c.Handler(next)
}
