package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
)

// Hardcoded credentials
const (
	dbPassword = "SuperSecret123!"  // VULNERABLE: Hardcoded password
	apiKey     = "123"
)

// SQL Injection vulnerability
func getUserByID(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")

	db, err := sql.Open("mysql", "root:"+dbPassword+"@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// VULNERABLE: String concatenation in SQL query
	query := "SELECT * FROM users WHERE id = " + userID
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	fmt.Fprintf(w, "User data: %v", rows)
}

// Command Injection vulnerability
func executeCommand(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")

	// VULNERABLE: User input in system command
	cmd := exec.Command("cat", "/var/log/"+filename)
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(output)
}

// Path Traversal vulnerability
func readFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")

	// VULNERABLE: No path validation
	content, err := os.ReadFile("/app/data/" + filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(content)
}

// XSS vulnerability
func displayMessage(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("msg")

	// VULNERABLE: Unescaped user input in HTML response
	html := "<html><body><h1>" + message + "</h1></body></html>"
	fmt.Fprint(w, html)
}

// Weak cryptography
func hashPassword(password string) string {
	// VULNERABLE: Using MD5 for password hashing
	// Should use bcrypt or argon2
	return fmt.Sprintf("%x", password) // Simplified for demo
}

// Insecure random number generation
func generateToken() string {
	// VULNERABLE: Using weak random source
	// Should use crypto/rand instead of math/rand
	return "token_12345" // Simplified for demo
}

// Sensitive data in logs
func logUserActivity(username, password string) {
	// VULNERABLE: Logging sensitive data
	fmt.Printf("User login attempt: username=%s, password=%s\n", username, password)
}

// Missing error handling
func connectDatabase() *sql.DB {
	db, _ := sql.Open("mysql", "root:password@/mydb")
	// VULNERABLE: Ignoring error return value
	return db
}

// Insecure TLS configuration
func setupHTTPS() {
	// VULNERABLE: Accepting any certificate
	// http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
	//     InsecureSkipVerify: true,
	// }
}

// Race condition
var counter int

func incrementCounter(w http.ResponseWriter, r *http.Request) {
	// VULNERABLE: Race condition - no mutex protection
	counter++
	fmt.Fprintf(w, "Counter: %d", counter)
}

// SSRF vulnerability
func fetchURL(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")

	// VULNERABLE: Fetching user-provided URL without validation
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	fmt.Fprintf(w, "Fetched: %v", resp.Status)
}

func main() {
	http.HandleFunc("/user", getUserByID)
	http.HandleFunc("/execute", executeCommand)
	http.HandleFunc("/read", readFile)
	http.HandleFunc("/display", displayMessage)
	http.HandleFunc("/counter", incrementCounter)
	http.HandleFunc("/fetch", fetchURL)

	fmt.Println("Server starting on port 8080...")
	http.ListenAndServe(":8080", nil)
}
