package handlers

import (
	"fmt"
	"net/http"

	// "github.com/kailashchoudhary11/repo-guard/initializers"
	// "github.com/kailashchoudhary11/repo-guard/services"
	"github.com/kailashchoudhary11/repo-guard/templates"
)

func Setup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// fmt.Println("This is the setup route verify and then serve the webpage")
	// code := r.URL.Query().Get("code")
	// if code == "" {
	// 	fmt.Println("Code query param is missing")
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	return
	// }
	// token, err := initializers.Oauth2Config.Exchange(ctx, code)
	// if err != nil {
	// 	fmt.Println("Invalid code", err)
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	// githubClient := initializers.GetClientWithToken(token.AccessToken)
	// ok, err := services.IsAppInstallationOwner(ctx, githubClient)
	// if err != nil {
	// 	fmt.Println("Error in listing the user installations", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// if !ok {
	// 	fmt.Println("Malicious User Not allowed to update this particular installation")
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	template := templates.SetupPage()
	if err := template.Render(ctx, w); err != nil {
		fmt.Println("Error in rendering the template")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
