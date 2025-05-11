package transport

import (
	"context"
	"errors"

	userpb "github.com/eotet/project-protos/gen/go/user"
	errs "github.com/eotet/users-service/internal/errors"
	"github.com/eotet/users-service/internal/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	svc *user.Service
	userpb.UnimplementedUserServiceServer
}

func NewHandler(svc *user.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.CreateUserResponse, error) {
	userToCreate := user.User{
		Email:    req.Email,
		Password: req.Password,
	}

	createdUser, err := h.svc.CreateUser(userToCreate)
	if err != nil {
		if errors.Is(err, errs.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &userpb.CreateUserResponse{
		User: &userpb.User{
			Id:    uint32(createdUser.ID),
			Email: createdUser.Email,
		},
	}, nil
}

func (h *Handler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	user, err := h.svc.GetUserByID(req.Id)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return &userpb.GetUserResponse{
		User: &userpb.User{
			Id:    uint32(user.ID),
			Email: user.Email,
		},
	}, nil
}

func (h *Handler) ListUsers(ctx context.Context, req *emptypb.Empty) (*userpb.ListUsersResponse, error) {
	users, err := h.svc.GetAllUsers()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}

	var response userpb.ListUsersResponse
	for i := range users {
		response.Users = append(response.Users, &userpb.User{
			Id:    uint32(users[i].ID),
			Email: users[i].Email,
		})
	}

	return &response, nil
}

func (h *Handler) UpdateUser(ctx context.Context, req *userpb.UpdateUserRequest) (*userpb.UpdateUserResponse, error) {
	updatedUser, err := h.svc.UpdateUserByID(req.Id, user.UpdateUserRequest{
		Email:    &req.Email,
		Password: &req.Password,
	})
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &userpb.UpdateUserResponse{
		User: &userpb.User{
			Id:    uint32(updatedUser.ID),
			Email: updatedUser.Email,
		},
	}, nil
}

func (h *Handler) DeleteUser(ctx context.Context, req *userpb.DeleteUserRequest) (*emptypb.Empty, error) {
	err := h.svc.DeleteUserByID(req.Id)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &emptypb.Empty{}, nil
}
