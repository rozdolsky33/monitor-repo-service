package cmd

//package main

//
//import (
//	"encoding/json"
//	"fmt"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"os/exec"
//	"time"
//)
//
//const (
//	repo          = "rozdolsky33/monitor-repo-service"
//	branch        = "main"
//	apiURL        = "https://api.github.com/repos/" + repo + "/commits/" + branch
//	checkInterval = 60 * time.Second
//)
//
//type Commit struct {
//	SHA string `json:"sha"`
//}
//
//func gitPull(repoDir string) error {
//	cmd := exec.Command("git", "-C", repoDir, "pull")
//	output, err := cmd.CombinedOutput()
//	if err != nil {
//		return fmt.Errorf("error pulling from git: %v - %s", err, string(output))
//	}
//	log.Printf("Git pull output: %s", string(output))
//	return nil
//}
//
//func getLatestCommit() (string, error) {
//	resp, err := http.Get(apiURL)
//	if err != nil {
//		return "", fmt.Errorf("error fetching commit data: %v", err)
//	}
//	defer resp.Body.Close()
//
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return "", fmt.Errorf("error reading response body: %v", err)
//	}
//
//	var commit Commit
//	if err := json.Unmarshal(body, &commit); err != nil {
//		return "", fmt.Errorf("error unmarshalling response: %v", err)
//	}
//
//	return commit.SHA, nil
//}
//
//func main() {
//	repoDir := "~/"
//	scriptPath := "./upload.sh"
//
//	var lastCommit string
//
//	// Get the initial latest commit
//	latestCommit, err := getLatestCommit()
//	if err != nil {
//		log.Fatalf("Failed to get latest commit: %v", err)
//	}
//	lastCommit = latestCommit
//
//	for {
//		time.Sleep(checkInterval)
//
//		latestCommit, err = getLatestCommit()
//		if err != nil {
//			log.Printf("Error getting latest commit: %v", err)
//			continue
//		}
//
//		if latestCommit != lastCommit {
//			log.Printf("New commit detected: %s", latestCommit)
//
//			if err := gitPull(repoDir); err != nil {
//				log.Printf("Error pulling from git: %v", err)
//				continue
//			}
//
//			err := exec.Command("/bin/sh", scriptPath).Run()
//			if err != nil {
//				log.Printf("Error executing script: %v", err)
//			}
//
//			lastCommit = latestCommit
//		}
//	}
//}
