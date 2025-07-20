package routes

import (
	"net/http"

	"github.com/greatdaveo/SendlyPay/controllers"
)

func WhatsAppRoutes() {
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to SendlyPay!!!"))
	})

	http.HandleFunc("/webhook/whatsapp", controllers.HandleIncomingMessages)

}
