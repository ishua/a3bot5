package queue

import "context"

func Add(ctx context.Context, queueName, payload string) error {
	return nil
}

func GetLast(ctx context.Context, queueName string) ([]byte, error) {
	return nil, nil
}
