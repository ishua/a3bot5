package queue

import "context"

func AddMessage(ctx context.Context, key string, message []byte) error {
	return nil
}

func GetLastMessage(ctx context.Context, qkey string) ([]byte, error) {
	return nil, nil
}
