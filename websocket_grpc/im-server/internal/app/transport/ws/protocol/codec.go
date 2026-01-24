package protocol

import (
	imv1 "github.com/brianliu-sysu/golang-knowledge/websocket_grpc/im-server/internal/gen/im/v1"

	proto "google.golang.org/protobuf/proto"
)

func DecodeClientMessage(data []byte) (*imv1.ClientEnvelope, error) {
	msg := &imv1.ClientEnvelope{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func EncodeClientMessage(msg *imv1.ClientEnvelope) ([]byte, error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func EncodeServerMessage(msg *imv1.ServerEnvelope) ([]byte, error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return data, nil
}
