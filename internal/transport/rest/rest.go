package rest

import (
	"fmt"
	"html/template"
	"log"
	"module/internal/domain"

	"net/http"
	"strings"

	"github.com/patrickmn/go-cache"
)

func ServeTemplate(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Result string
		}{}
		if err := tmpl.Execute(w, data); err != nil {
			log.Printf("Error when displaying the template: %v", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	}
}
func SearchOrderData(w http.ResponseWriter, r *http.Request, c *cache.Cache, tmpl *template.Template) {
	if r.Method != "GET" {
		http.Error(w, "Only GET requests are allowed.", http.StatusMethodNotAllowed)
		return
	}

	orderUID := r.FormValue("orderUID")

	data := struct {
		Result string
	}{}

	if order, found := c.Get(fmt.Sprintf("order_%s", orderUID)); found {
		orderData := order.(domain.OrderCache)

		data.Result += fmt.Sprintf("Order cached: %+v\n", orderData)
		if delivery, found := c.Get(fmt.Sprintf("delivery_%d", orderData.Delivery)); found {
			str := fmt.Sprintf("Delivery cached: %+v\n", delivery)
			data.Result += str
		} else {
			fmt.Println("Delivery not found in cache")
		}

		if payment, found := c.Get(fmt.Sprintf("payment_%d", orderData.Payment)); found {
			str := fmt.Sprintf("Payment cached: %+v\n", payment)
			data.Result += str
		} else {
			fmt.Println("Payment not found in cache")
		}

		orderItemKey := fmt.Sprintf("order_item_%s_", orderUID)
		for key, value := range c.Items() {
			if strings.HasPrefix(key, orderItemKey) {
				orderItem := value.Object.(domain.OrderItemCache)
				if item, found := c.Get(fmt.Sprintf("item_%d", orderItem.ItemID)); found {
					str := fmt.Sprintf("Item cached: %+v\n", item)
					data.Result += str
				} else {
					fmt.Println("Item not found in cache")
				}
			}
		}
	} else {
		fmt.Println("Order not found in cache")
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error when displaying the template: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}
