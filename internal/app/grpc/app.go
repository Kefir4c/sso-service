package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	authgrpc "github.com/Kefir4c/sso-service/internal/grpc/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// App represents gRPC server with middleware.
type App struct {
	log    *slog.Logger
	server *grpc.Server
	port   int
}

// New creates gRPC server with recovery and logging interceptors.
// Registers auth service.
func New(log *slog.Logger, auth authgrpc.Auth, port int, timeout time.Duration) *App {

	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived,
			logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			log.Error("recovery of panic", slog.Any("panic", p))
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
	))

	authgrpc.Register(gRPCServer, auth, timeout)

	return &App{
		log:    log,
		server: gRPCServer,
		port:   port,
	}
}

// InterceptorLogger adapts slog for grpc-middleware.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, level logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(level), msg, fields...)
	})
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run "

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))

	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	a.log.Info("grpc server started", slog.String("port", l.Addr().String()))

	if err := a.server.Serve(l); err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.stop"

	a.log.With(
		slog.String("op", op),
	).Info("stopping grpc server", slog.Int("port", a.port))

	a.server.GracefulStop()
}
