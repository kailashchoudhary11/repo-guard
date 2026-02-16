package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "github.com/kailashchoudhary11/repo-guard/db/generated"
	"github.com/kailashchoudhary11/repo-guard/initializers"
	"github.com/kailashchoudhary11/repo-guard/models"
	"github.com/kailashchoudhary11/repo-guard/services"
	"github.com/kailashchoudhary11/repo-guard/templates"
)

func Setup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fmt.Println("This is the setup route verify and then serve the webpage")
	code := r.URL.Query().Get("code")
	if code == "" {
		fmt.Println("Code query param is missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	installationID := r.URL.Query().Get("installation_id")
	if installationID == "" {
		fmt.Println("installation_id query param is missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := initializers.Oauth2Config.Exchange(ctx, code)
	if err != nil {
		fmt.Println("Invalid code", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	githubClient := initializers.GetClientWithToken(token.AccessToken)
	ok, err := services.IsAppInstallationOwner(ctx, githubClient)
	if err != nil {
		fmt.Println("Error in listing the user installations", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !ok {
		fmt.Println("Malicious User Not allowed to update this particular installation")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	template := templates.SetupPage(installationID)
	if err := template.Render(ctx, w); err != nil {
		fmt.Println("Error in rendering the template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func SetupSave(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		fmt.Println("Error parsing form:", err)
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	installationID := r.FormValue("installation_id")
	if installationID == "" {
		fmt.Println("installation_id is missing from form")
		http.Error(w, "Missing installation ID", http.StatusBadRequest)
		return
	}

	config := models.InstallationConfig{
		ShouldClose: r.FormValue("auto_close") == "on",
		Language:    r.FormValue("language"),
		Sensitivity: r.FormValue("sensitivity"),
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		fmt.Println("Error marshaling config:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	conn, err := initializers.GetDBClient(ctx)
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer conn.Close(ctx)

	queries := db.New(conn)
	_, err = queries.UpsertInstallationConfig(ctx, db.UpsertInstallationConfigParams{
		InstallationID: installationID,
		ConfigData:     configJSON,
		UpdatedBy:      installationID,
	})
	if err != nil {
		fmt.Println("Error saving config to database:", err)
		http.Error(w, "Failed to save configuration", http.StatusInternalServerError)
		return
	}

	fmt.Println("Configuration saved for installation:", installationID)
	http.Redirect(w, r, "/?saved=true", http.StatusSeeOther)
}
