package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type PushEvent struct {
	Ref string `json:"ref"`
}

func gitPull(repoDir string) error {
	token := os.Getenv("GIT_PAT")
	if token == "" {
		return fmt.Errorf("GIT_PAT environment variable not set")
	}

	repoURL := fmt.Sprintf("https://%s@github.com/rozdolsky33/monitor-repo-service.git", token)

	cmd := exec.Command("git", "-C", repoDir, "pull", repoURL)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error pulling from git: %v - %s", err, string(output))
	}
	log.Printf("Git pull output: %s", string(output))
	return nil
}

func runScript(scriptPath string) error {
	cmd := exec.Command("/bin/sh", scriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing script: %v - %s", err, string(output))
	}
	log.Printf("Script output: %s", string(output))
	return nil
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var event PushEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	if event.Ref == "refs/heads/main" {
		log.Println("Push event detected on main branch!")

		// Pull the latest changes from the repository
		repoDir := "/root/monitor-repo-service"
		scriptPath := "/root/monitor-repo-service/upload.sh"
		if err := gitPull(repoDir); err != nil {
			log.Printf("Error pulling from git: %v", err)
			http.Error(w, "Failed to pull from git", http.StatusInternalServerError)
			return
		}

		// Run the script and capture its output
		if err := runScript(scriptPath); err != nil {
			log.Printf("Error executing script: %v", err)
			http.Error(w, "Failed to execute script", http.StatusInternalServerError)
			return
		}

		log.Println("Script executed successfully!")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received"))
}

func main() {
	http.HandleFunc("/webhook", handleWebhook)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
