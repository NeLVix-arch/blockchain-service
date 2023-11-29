package middlewares

import (
	"log"
	"net/http"
)

// recoverMiddleware - middleware для восстановления от паники
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Отложенная функция для перехвата паники
		defer func() {
			if err := recover(); err != nil {
				// Логируем ошибку
				log.Printf("Паника: %v\n", err)
				// Отправляем ответ с кодом 500
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			}
		}()

		// Вызываем следующий обработчик
		next.ServeHTTP(w, r)
	})
}
