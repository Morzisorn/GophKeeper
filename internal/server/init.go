package server

import (
	"context"
	"gophkeeper/config"
	"gophkeeper/internal/logger"
	pbus "gophkeeper/internal/protos/users"
	pbcs "gophkeeper/internal/protos/crypto"
	pbit "gophkeeper/internal/protos/items"
	"gophkeeper/internal/server/controllers"
	cserv "gophkeeper/internal/server/services/crypto_service"
	userv "gophkeeper/internal/server/services/user_service"
	iserv "gophkeeper/internal/server/services/item_service"
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

func CreateAndRun(us *userv.UserService, cs *cserv.CryptoService, is *iserv.ItemService) {
	var g *GRPCServer

	g = createGRPCServer(us, cs, is)

	g.Run()
}

type GRPCServer struct {
	Server  *grpc.Server
	Listen  net.Listener

	US *userv.UserService
	CS *cserv.CryptoService
	IS *iserv.ItemService
}

func createGRPCServer(us *userv.UserService, cs *cserv.CryptoService, is *iserv.ItemService) *GRPCServer {
	uc := controllers.NewUserController(us)
	cc := controllers.NewCryptoController()
	ic := controllers.NewItemController(is)
	cnfg := config.GetServerConfig()
	listen, err := net.Listen("tcp", cnfg.Addr)
	if err != nil {
		logger.Log.Fatal("create listener error", zap.Error(err))
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(controllers.AuthInterceptor),
	)
	pbus.RegisterUserControllerServer(s, uc)
	pbcs.RegisterCryptoControllerServer(s, cc)
	pbit.RegisterItemsControllerServer(s, ic)

	return &GRPCServer{
		Server: s, 
		Listen:listen,

		US: us,
		CS: cs,
		IS: is,
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
