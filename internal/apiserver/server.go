package apiserver

import (
	"net/http"

	"github.com/compico/restauth/internal/db"
)

type ApiServer struct {
	Server *http.Server
	DB     *db.DB
}

func InitializingServer(conf *http.Server) *ApiServer {
	server := new(ApiServer)
	server.DB = initializingMongoDBDriver()
	server.Server = conf
	return server
}
