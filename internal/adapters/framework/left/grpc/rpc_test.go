package rpc

import (
	"context"
	"hex/internal/adapters/app/api"
	"hex/internal/adapters/core/arithmetic"
	"hex/internal/adapters/framework/left/grpc/proto/pb"
	"hex/internal/adapters/framework/right/db"
	"hex/internal/ports"
	"log"
	"net"
	"os"
    "testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const buffSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
    var err error
    lis = bufconn.Listen(buffSize)
    grpcServer := grpc.NewServer()

    // ports

    var dbaseAdapter ports.DbPort
    var core ports.ArithmeticPort
    var appAdapter ports.APIPort
    var gRPCAdapter ports.GRPCPort

    dbaseDriver := os.Getenv("DB_DRIVER")
    dsourceName := os.Getenv("DS_NAME")

    dbaseAdapter, err = db.NewAdapter(dbaseDriver, dsourceName)

    if err != nil {
        log.Fatalf("failed to init the db connection: %v", err)
    }

    core = arithmetic.NewAdapter()

    appAdapter = api.NewAdapter(dbaseAdapter, core)

    gRPCAdapter = NewAdapter(appAdapter)

    pb.RegisterArithmeticServiceServer(grpcServer, gRPCAdapter)

    go func() {
        if err := grpcServer.Serve(lis); err != nil {
            log.Fatalf("test server start error: %v", err)
        }
    }()
}

func buffDialer(context.Context, string) (net.Conn, error) {
    return lis.Dial()
}

func getGRPCConnection(ctx context.Context, t *testing.T) *grpc.ClientConn {
    conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(buffDialer), grpc.WithInsecure)

    if err != nil {
        log.Fatalf("failed to dial bufnet: %v", err)
    }

    return conn
}

func TestGetAddition(t *testing.T) {
    ctx := context.Background()
    conn := getGRPCConnection(ctx, t)
    defer conn.Close()

    client := pb.NewArithmeticServiceClient(conn)

    params := &pb.OperationParameters{
        A:1,
        B:1,
    }

    answer, err := client.GetAddition(ctx, params)

    if err != nil {
        log.Fatalf("expected: %v, got: %v", nil, err)
    }

    require.Equal(t, answer.Value, int32(2))
}


func TestGetSubtraction(t *testing.T) {
    ctx := context.Background()
    conn := getGRPCConnection(ctx, t)
    defer conn.Close()

    client := pb.NewArithmeticServiceClient(conn)

    params := &pb.OperationParameters{
        A:1,
        B:1,
    }

    answer, err := client.GetSubtraction(ctx, params)

    if err != nil {
        log.Fatalf("expected: %v, got: %v", nil, err)
    }

    require.Equal(t, answer.Value, int32(0))
}

func TestGetMultiplication(t *testing.T) {
    ctx := context.Background()
    conn := getGRPCConnection(ctx, t)
    defer conn.Close()

    client := pb.NewArithmeticServiceClient(conn)

    params := &pb.OperationParameters{
        A:1,
        B:1,
    }

    answer, err := client.GetMultiplication(ctx, params)

    if err != nil {
        log.Fatalf("expected: %v, got: %v", nil, err)
    }

    require.Equal(t, answer.Value, int32(1))
}

func TestGetDivision(t *testing.T) {
    ctx := context.Background()
    conn := getGRPCConnection(ctx, t)
    defer conn.Close()

    client := pb.NewArithmeticServiceClient(conn)

    params := &pb.OperationParameters{
        A:1,
        B:1,
    }

    answer, err := client.GetDivision(ctx, params)

    if err != nil {
        log.Fatalf("expected: %v, got: %v", nil, err)
    }

    require.Equal(t, answer.Value, int32(1))
}
