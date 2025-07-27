package server

import (
	"context"
	"fmt"
	"gophkeeper/config"
	"gophkeeper/internal/logger"
	pbcs "gophkeeper/internal/protos/crypto"
	pbit "gophkeeper/internal/protos/items"
	pbus "gophkeeper/internal/protos/users"
	"gophkeeper/internal/server/controllers"
	cserv "gophkeeper/internal/server/services/crypto_service"
	iserv "gophkeeper/internal/server/services/item_service"
	userv "gophkeeper/internal/server/services/user_service"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

type Server interface {
	Create(us *userv.UserService) error
	Run() error
	Shutdown(ctx context.Context, idleConnsClosed chan struct{})
}

func CreateAndRun(us *userv.UserService, cs *cserv.CryptoService, is *iserv.ItemService) error {
	g, err := createGRPCServer(us, cs, is)
	if err != nil {
		return fmt.Errorf("create grpc server error: %w\n", err)
	}

	if err := g.Run(); err != nil {
		return fmt.Errorf("grpc server error: %w\n", err)
	}

	return nil
}

type GRPCServer struct {
	Server *grpc.Server
	Listen net.Listener

	US *userv.UserService
	CS *cserv.CryptoService
	IS *iserv.ItemService
}

func createGRPCServer(us *userv.UserService, cs *cserv.CryptoService, is *iserv.ItemService) (*GRPCServer, error) {
	uc := controllers.NewUserController(us)
	cc := controllers.NewCryptoController()
	ic := controllers.NewItemController(is)
	cnfg, err := config.GetServerConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get server config: %w", err)
	}
	listen, err := net.Listen("tcp", cnfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("create listener error: %w", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(controllers.AuthInterceptor),
	)
	pbus.RegisterUserControllerServer(s, uc)
	pbcs.RegisterCryptoControllerServer(s, cc)
	pbit.RegisterItemsControllerServer(s, ic)

	return &GRPCServer{
		Server: s,
		Listen: listen,

		US: us,
		CS: cs,
		IS: is,
	}, nil
}

func (s *GRPCServer) Run() error {
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

		return fmt.Errorf("failed to run grpc server: %w", err)
	}

	<-idleConnsClosed

	logger.Log.Info("Server shutted down gracefully")
	return nil
}

func (s *GRPCServer) Shutdown(ctx context.Context, idleConnsClosed chan struct{}) {
	logger.Log.Info("Shutdown server")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	//(*s.Storage).Close()

	s.Server.GracefulStop()

	close(idleConnsClosed)
}
