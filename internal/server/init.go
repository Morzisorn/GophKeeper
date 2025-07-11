package server

import (
	"context"
	"gophkeeper/config"
	"gophkeeper/internal/logger"
	pbus "gophkeeper/internal/protos/users"
	pbcs "gophkeeper/internal/protos/crypto"
	"gophkeeper/internal/server/controllers"
	cserv "gophkeeper/internal/server/services/crypto_service"
	userv "gophkeeper/internal/server/services/user_service"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server interface {
	Create(us *userv.UserService)
	Run()
	Shutdown(ctx context.Context, idleConnsClosed chan struct{})
}

func CreateAndRun(us *userv.UserService, cs *cserv.CryptoService) {
	var g *GRPCServer

	g = createGRPCServer(us, cs)

	g.Run()
}

type GRPCServer struct {
	Server  *grpc.Server
	Listen  net.Listener

	US *userv.UserService
	CS *cserv.CryptoService
}

func createGRPCServer(us *userv.UserService, cs *cserv.CryptoService) *GRPCServer {
	uc := controllers.NewUserController(us)
	cc := controllers.NewCryptoController()
	cnfg := config.GetServerConfig()
	listen, err := net.Listen("tcp", cnfg.Addr)
	if err != nil {
		logger.Log.Fatal("create listener error", zap.Error(err))
	}

	s := grpc.NewServer()
	pbus.RegisterUserControllerServer(s, uc)
	pbcs.RegisterCryptoControllerServer(s, cc)

	return &GRPCServer{
		Server: s, 
		Listen:listen,

		US: us,
		CS: cs,
	}
}

func (s *GRPCServer) Run() {
	logger.Log.Info("Run grpc server")

	idleConnsClosed := make(chan struct{})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.Shutdown(ctx, idleConnsClosed)
	}()

	if err := s.Server.Serve(s.Listen); err != nil {
		logger.Log.Fatal("Failed to run grpc server", zap.Error(err))
	}

	<-idleConnsClosed

	logger.Log.Info("Server shutted down gracefully")
}

func (s *GRPCServer) Shutdown(ctx context.Context, idleConnsClosed chan struct{}) {
	logger.Log.Info("Shutdown server")

	//(*s.Storage).Close()

	s.Server.GracefulStop()

	close(idleConnsClosed)
}
