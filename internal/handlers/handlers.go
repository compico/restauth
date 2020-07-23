package handlers

import (
	"fmt"
	"net/http"

	"github.com/compico/restauth/internal/user"
	"github.com/julienschmidt/httprouter"
)

func addNewUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, users *user.Users) {
	guid := ps.ByName("guid")
	uuid, err := user.VerifyGUID(guid)
	if err != nil {
		fmt.Fprintf(w, "[ERROR] %v", err)
		return
	}
	usr := user.User{
		GUID: uuid.String(),
	}
	users.User = append(users.User, usr)
}
