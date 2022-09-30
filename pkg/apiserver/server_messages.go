package apiserver

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/lab5e/l5log/pkg/lg"
	"github.com/lab5e/lospan/pkg/model"
	"github.com/lab5e/lospan/pkg/pb/lospan"
	"github.com/lab5e/lospan/pkg/protocol"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *apiServer) Inbox(ctx context.Context, req *lospan.InboxRequest) (*lospan.InboxResponse, error) {
	eui, err := protocol.EUIFromString(req.Eui)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid EUI")
	}
	list, err := a.store.ListUpstreamMessages(eui, 1000)
	if err != nil {
		return nil, toProtoErr(err)
	}
	ret := &lospan.InboxResponse{
		Messages: make([]*lospan.UpstreamMessage, 0),
	}
	for _, msg := range list {
		ret.Messages = append(ret.Messages, &lospan.UpstreamMessage{
			Eui:        msg.DeviceEUI.String(),
			Timestamp:  msg.Timestamp,
			Payload:    msg.Data[:],
			GatewayEui: msg.GatewayEUI.String(),
			Rssi:       msg.RSSI,
			Snr:        msg.SNR,
			Frequency:  msg.Frequency,
			DataRate:   msg.DataRate,
			DevAddr:    msg.DevAddr.ToUint32(),
		})
	}
	return ret, nil
}

func (a *apiServer) Outbox(ctx context.Context, req *lospan.OutboxRequest) (*lospan.OutboxResponse, error) {
	eui, err := protocol.EUIFromString(req.Eui)
	if err != nil {
		lg.Warning("Invalid EUI: %s: %v", req.Eui, err)
		return nil, status.Error(codes.InvalidArgument, "Invalid EUI")
	}
	list, err := a.store.ListDownstreamMessages(eui)
	if err != nil {
		return nil, toProtoErr(err)
	}

	ret := &lospan.OutboxResponse{
		Messages: make([]*lospan.DownstreamMessage, 0),
	}

	for _, msg := range list {
		payload, err := hex.DecodeString(msg.Data)
		if err != nil {
			lg.Warning("Error decoding payload for downstream message from device %s: %v", msg.DeviceEUI, err)
			continue
		}
		ret.Messages = append(ret.Messages, &lospan.DownstreamMessage{
			Eui:     msg.DeviceEUI.String(),
			Payload: payload,
			Port:    int32(msg.Port),
			Ack:     msg.Ack,
			Created: newPtr(msg.CreatedTime),
			Sent:    newPtr(msg.SentTime),
			AckTime: newPtr(msg.AckTime),
		})
	}
	return ret, nil
}

func (a *apiServer) SendMessage(ctx context.Context, req *lospan.DownstreamMessage) (*lospan.DownstreamMessage, error) {
	eui, err := protocol.EUIFromString(req.Eui)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid EUI")
	}
	if req.Port > 255 || req.Port < 0 {
		return nil, status.Error(codes.InvalidArgument, "Port must be 0-255")
	}
	msg := model.DownstreamMessage{
		DeviceEUI:   eui,
		Data:        hex.EncodeToString(req.Payload),
		Port:        uint8(req.Port),
		Ack:         req.Ack,
		CreatedTime: time.Now().UnixMilli(),
		SentTime:    0,
		AckTime:     0,
	}
	if err := a.store.CreateDownstreamMessage(eui, msg); err != nil {
		return nil, toProtoErr(err)
	}

	return &lospan.DownstreamMessage{
		Eui:     req.Eui,
		Payload: req.Payload[:],
		Port:    int32(msg.Port),
		Ack:     msg.Ack,
		Created: newPtr(msg.CreatedTime),
		Sent:    newPtr(msg.SentTime),
		AckTime: newPtr(msg.AckTime),
	}, nil
}
