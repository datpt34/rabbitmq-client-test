package main

import (
	"fmt"
	"net/http"
)

// Struct để đọc dữ liệu yêu cầu xác thực từ RabbitMQ
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Authenticated bool `json:"authenticated"`
}

func main() {
	http.HandleFunc("/auth/user", authUser)
	http.HandleFunc("/auth/vhost", authDefault)
	http.HandleFunc("/auth/resource", authDefault)
	http.HandleFunc("/auth/topic", authDefault)

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func authUser(w http.ResponseWriter, r *http.Request) {
	fmt.Printf(">>>>Request: %s, path: %s\n", r.Method, r.URL.Path)
	// Chỉ chấp nhận yêu cầu POST
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Đọc dữ liệu yêu cầu

	r.ParseForm()
	authRequest := AuthRequest{
		Username: r.Form.Get("username"),
		Password: r.Form.Get("password"),
	}
	for key, value := range r.Form {
		fmt.Printf("Key: %s, Value: %s\n", key, value)
	}
	//if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
	//	fmt.Println("error NewDecoder: ", err)
	//	http.Error(w, "Bad request", http.StatusBadRequest)
	//	return
	//}

	// Thực hiện xác thực, ví dụ: kiểm tra người dùng trong cơ sở dữ liệu
	if authRequest.Username == "rabbitmq_user" && authRequest.Password == "rabbitmq_password" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `allow`)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `deny`)
	}
}
func authDefault(w http.ResponseWriter, r *http.Request) {
	fmt.Printf(">>>>Request: %s, path: %s\n", r.Method, r.URL.Path)
	r.ParseForm()
	for key, value := range r.Form {
		fmt.Printf("Key: %s, Value: %s\n", key, value)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `allow`)

}
