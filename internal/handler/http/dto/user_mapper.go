package dto

import "github.com/wonjinsin/go-boilerplate/internal/domain"

// ToUserResponse converts domain.User to UserResponse
func ToUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}

// ToUserListResponse converts domain.Users to UserListResponse
func ToUserListResponse(users domain.Users, total, offset, limit int) UserListResponse {
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = ToUserResponse(user)
	}

	return UserListResponse{
		Users:  userResponses,
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}
}
