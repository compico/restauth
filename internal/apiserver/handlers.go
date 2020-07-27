package apiserver

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/compico/restauth/internal/db"
	"github.com/compico/restauth/internal/middleware"
	"github.com/compico/restauth/internal/user"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

var users = user.InitUsers()

func AddNewUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	guid := ps.ByName("guid")
	uuid, err := user.VerifyGUID(guid)
	if err != nil {
		// Если UUID(GUID) не вальдный - возвращает ошибку
		// и отменяет авторизацию
		fmt.Fprintf(w, "{error:%v}", err)
		return
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = "unknown"
	}
	accesscfg := middleware.JWTConfig{
		Data: jwt.MapClaims{
			"guid":       uuid.String(),
			"ip":         ip,
			"createtime": time.Now().Unix(),
			"lifetime":   time.Now().Add(730 * time.Hour).Unix(),
		},
		Method: jwt.SigningMethodHS256,
	}
	refreshcfg := middleware.JWTConfig{
		Data: jwt.MapClaims{
			"guid":       uuid,
			"ip":         ip,
			"createtime": time.Now().Unix(),
		},
		Method: jwt.SigningMethodES256,
	}
	usr := user.User{
		Access:  make([]*jwt.Token, 0),
		Refresh: make([]*jwt.Token, 0),
		GUID:    uuid.String(),
	}
	usr.Access = append(usr.Access, accesscfg.NewToken())
	fmt.Println("[DEBUG] Create access token")
	usr.Access = append(usr.Refresh, refreshcfg.NewToken())
	fmt.Println("[DEBUG] Create refresh token")
	users.User = append(users.User, usr)

	for i := 0; i < len(users.User); i++ {
		fmt.Printf("User №%v:\n", i)
		for j := 0; j < len(users.User[i].Access); j++ {
			x, err := users.User[i].Access[j].SigningString()
			if err != nil {
				return
			}
			fmt.Printf("\tAscess: %v", x)
		}
		for j := 0; j < len(users.User[i].Refresh); j++ {
			y, err := users.User[i].Refresh[j].SigningString()
			if err != nil {
				return
			}
			fmt.Printf("\tRefresh: %v", y)
		}
	}
}

func Test(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	guid := ps.ByName("guid")
	client := db.NewClient()
	err := client.Connect()
	if err != nil {
		fmt.Fprintf(w, "{ \"error\" : \"%v\" }", err)
	}
	defer func() {
		err := client.Disconnect()
		if err != nil {
			fmt.Printf("[ERROR]%v\n", err)
		}
	}()
	coll, err := client.HaveCollection(guid)
	if err != nil {
		fmt.Fprintf(w, "{ \"error\" : \"%v\" }", err)
	}
	fmt.Fprintf(w, "%#v", coll)
}
