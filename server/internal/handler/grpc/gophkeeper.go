// Package grpc implements the gRPC server for the GophKeeper service,
// handling user registration, authentication, and operations related to data items.
package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "gophKeeper/pkg/proto/gophkeeper"
	dataItemsModel "gophKeeper/server/internal/domain/dataItems/model"
	dataItemsU "gophKeeper/server/internal/usecase/dataItems"
	usersU "gophKeeper/server/internal/usecase/users"
)

// St implements the GophKeeperServiceServer interface, providing gRPC handlers
// for user management and data item operations. It uses use cases for both users
// and data items to perform business logic.
type St struct {
	pb.UnsafeGophKeeperServiceServer
	dataItemsUcs *dataItemsU.Usecase
	usersUcs     *usersU.Usecase
}

// New creates a new instance of the St gRPC server with the given use cases for data items and users.
func New(dataItemsUcs *dataItemsU.Usecase, usersUcs *usersU.Usecase) *St {
	return &St{
		dataItemsUcs: dataItemsUcs,
		usersUcs:     usersUcs,
	}
}

// Register handles user registration requests, creating a new user in the system.
func (s *St) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := s.usersUcs.Register(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Canceled, "failed to register user: %v", err)
	}

	return &pb.RegisterResponse{
		Message: "ok",
	}, nil
}

// Login handles user login requests, returning a token if the credentials are valid.
func (s *St) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := s.usersUcs.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{Token: *token}, nil
}

// GetData retrieves a specific data item based on user ID and other provided parameters.
func (s *St) GetData(ctx context.Context, req *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	userID, err := s.usersUcs.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	obj, found, err := s.dataItemsUcs.GetData(ctx, &dataItemsModel.GetPars{
		ID:     req.Id,
		UserID: userID,
		Type:   req.Type,
		URL:    req.URL,
	})
	if err != nil {
		return nil, err
	}
	if !found {
		return &pb.GetDataResponse{
			Data: []*pb.DataItem{},
		}, nil
	}

	dataItem := &pb.DataItem{
		Id:        obj.ID,
		Type:      obj.Type,
		Data:      obj.Data,
		Meta:      obj.Meta,
		CreatedAt: timestamppb.New(obj.CreatedAt),
		UpdatedAt: timestamppb.New(obj.UpdatedAt),
	}

	return &pb.GetDataResponse{
		Data: []*pb.DataItem{dataItem},
	}, nil
}

// ListData retrieves a data items based on user ID
func (s *St) ListData(ctx context.Context, _ *emptypb.Empty) (*pb.ListDataResponse, error) {
	userID, err := s.usersUcs.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	result, _, err := s.dataItemsUcs.ListAll(ctx, &dataItemsModel.ListPars{
		UserID: &userID,
	})
	if err != nil {
		return nil, err
	}

	dataItems := make([]*pb.DataItem, 0, len(result))
	for _, item := range result {
		dataItem := &pb.DataItem{
			Id:        item.ID,
			Type:      item.Type,
			Data:      item.Data,
			Meta:      item.Meta,
			CreatedAt: timestamppb.New(item.CreatedAt),
			UpdatedAt: timestamppb.New(item.UpdatedAt),
		}

		dataItems = append(dataItems, dataItem)
	}

	return &pb.ListDataResponse{
		Data: dataItems,
	}, nil
}

// CreateData handles requests to create new data items for a user.
func (s *St) CreateData(ctx context.Context, req *pb.CreateDataRequest) (*pb.CreateDataResponse, error) {
	userID, err := s.usersUcs.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	data := req.GetData()

	err = s.dataItemsUcs.CreateData(ctx, &dataItemsModel.Edit{
		ID:     data.Id,
		UserID: &userID,
		Type:   &data.Type,
		Data:   &data.Data,
		Meta:   &data.Meta,
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateDataResponse{Message: "Success"}, nil
}

// UpdateData handles requests to update existing data items for a user.
func (s *St) UpdateData(ctx context.Context, req *pb.UpdateDataRequest) (*pb.UpdateDataResponse, error) {
	userID, err := s.usersUcs.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	data := req.GetData()

	editData := &dataItemsModel.Edit{
		UserID: &userID,
		ID:     data.Id,
	}

	if data.Type != "" {
		editData.Type = &data.Type
	}
	if data.Data != nil {
		editData.Data = &data.Data
	}
	if data.Meta != "" {
		editData.Meta = &data.Meta
	}

	err = s.dataItemsUcs.EditData(ctx, editData)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateDataResponse{Message: "Update successful"}, nil
}

// DeleteData handles requests to delete data items based on user ID and item ID.
func (s *St) DeleteData(ctx context.Context, req *pb.DeleteDataRequest) (*pb.DeleteDataResponse, error) {
	userID, err := s.usersUcs.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = s.dataItemsUcs.DeleteData(ctx, &dataItemsModel.GetPars{
		ID:     req.Id,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	return &pb.DeleteDataResponse{Message: "Delete successful"}, nil
}

// SyncData handles requests to synchronize data between the client and the server (currently not implemented).
func (s *St) SyncData(ctx context.Context, req *pb.SyncDataRequest) (*pb.SyncDataResponse, error) {
	return nil, nil
}

// Ping handles requests to show is server available
func (s *St) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
