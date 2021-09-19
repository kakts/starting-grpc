package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func FirstRequestInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// リクエストを横取りして処理を加える
		fmt.Println("First Interceptor called!")

		// req = modifyRequest(req)
		// そのリクエストをRPCメソッドに相当するハンドラに渡し、ハンドラで処理された後の返り値をそのまま返却して、次の処理に渡す
		return handler(ctx, req)
	}
}

func SecondRequestInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// リクエストを横取りして処理を加える
		fmt.Println("Second Interceptor called!")

		// req = modifyRequest(req)
		// そのリクエストをRPCメソッドに相当するハンドラに渡し、ハンドラで処理された後の返り値をそのまま返却して、次の処理に渡す
		return handler(ctx, req)
	}
}

// レスポンスを返す前に実行するインターセプタ
// func FirstResponseInterceptor() grpc.UnaryServerInterceptor {
// 	return func(
// 		ctx context.Context,
// 		req interface{},
// 		info *grpc.UnaryServerInfo,
// 		handler grpc.UnaryHandler,
// 	) (interface{}, error) {
// 		// 実行されるRPCハンドラーの戻り値を取得
// 		res, err := handler(ctx, req)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return nil
// 	}
// }
