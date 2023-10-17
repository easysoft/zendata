package server

import (
	"fmt"

	serverConfig "github.com/easysoft/zendata/internal/server/config"
	logUtils "github.com/easysoft/zendata/pkg/utils/log"
	"github.com/easysoft/zendata/pkg/utils/vari"
	"github.com/facebookgo/inject"
)

// Server for admin server
type Server struct {
	Config *serverConfig.Config `inject:""`
}

func InitServer(config *serverConfig.Config) (server *Server, err error) {
	var g inject.Graph

	server = &Server{}

	// inject objects
	if err := g.Provide(
		&inject.Object{Value: config},
		&inject.Object{Value: vari.DB},
		&inject.Object{Value: server},
	); err != nil {
		logUtils.PrintErrMsg(fmt.Sprintf("provide usecase objects to the Graph: %v", err))
	}
	err = g.Populate()
	if err != nil {
		logUtils.PrintErrMsg(fmt.Sprintf("populate the incomplete Objects: %v", err))
	}

	return
}
