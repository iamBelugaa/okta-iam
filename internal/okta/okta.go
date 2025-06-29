package okta

import (
	"context"
	"fmt"

	"github.com/iamBelugaa/iam/internal/models"
	oktaSdk "github.com/okta/okta-sdk-golang/v5/okta"
	"go.uber.org/zap"
)

type Service struct {
	client *oktaSdk.APIClient
	log    *zap.SugaredLogger
}

func NewOktaService(client *oktaSdk.APIClient) *Service {
	return &Service{client: client}
}

func (oc *Service) CreateUser(ctx context.Context, userReq *models.CreateUserRequest) (*models.UserResponse, error) {
	var profile oktaSdk.UserProfile

	if userReq.FirstName != "" {
		profile.SetFirstName(userReq.FirstName)
	}

	if userReq.LastName != "" {
		profile.SetLastName(userReq.LastName)
	}

	if userReq.Email != "" {
		profile.SetEmail(userReq.Email)
		profile.SetLogin(userReq.Email)
	}

	if userReq.MobilePhone != "" {
		profile.SetMobilePhone(userReq.MobilePhone)
	}

	createUserRequest := oktaSdk.CreateUserRequest{
		Profile: profile,
		Credentials: &oktaSdk.UserCredentials{
			Password: &oktaSdk.PasswordCredential{
				Value: oktaSdk.PtrString(userReq.Password),
			},
		},
	}

	user, _, err := oc.client.UserAPI.CreateUser(ctx).Body(createUserRequest).Execute()
	if err != nil {
		oc.log.Errorf("Okta SDK: failed to create user: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return models.ConvertSDKUserToUserResponse(user), nil
}

func (oc *Service) GetUser(ctx context.Context, userID string) (*models.UserResponse, error) {
	user, _, err := oc.client.UserAPI.GetUser(ctx, userID).Execute()
	if err != nil {
		oc.log.Errorf("Okta SDK: failed to get user %s: %v", userID, err)
		return nil, fmt.Errorf("failed to get user %s: %w", userID, err)
	}

	return models.ConvertSDKUserToUserResponse(&oktaSdk.User{
		Id:                    user.Id,
		Created:               user.Created,
		Activated:             user.Activated,
		LastLogin:             user.LastLogin,
		Credentials:           user.Credentials,
		LastUpdated:           user.LastUpdated,
		PasswordChanged:       user.PasswordChanged,
		Profile:               user.Profile,
		RealmId:               user.RealmId,
		Status:                user.Status,
		StatusChanged:         user.StatusChanged,
		TransitioningToStatus: user.TransitioningToStatus,
		Type:                  user.Type,
		Links:                 user.Links,
		AdditionalProperties:  user.AdditionalProperties,
	}), nil
}

func (oc *Service) ListUsers(ctx context.Context) ([]*models.UserResponse, error) {

	users, _, err := oc.client.UserAPI.ListUsers(ctx).Execute()
	if err != nil {
		oc.log.Errorf("Okta SDK: failed to list users: %v", err)
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	userResponses := make([]*models.UserResponse, len(users))
	for i := range users {
		userResponses[i] = models.ConvertSDKUserToUserResponse(&users[i])
	}

	return userResponses, nil
}

func (oc *Service) DeactivateUser(ctx context.Context, userID string) error {
	_, err := oc.client.UserAPI.DeactivateUser(ctx, userID).Execute()
	if err != nil {
		oc.log.Errorf("Okta SDK: failed to deactivate user %s: %v", userID, err)
		return fmt.Errorf("failed to deactivate user %s: %w", userID, err)
	}
	return nil
}

func (oc *Service) DeleteUser(ctx context.Context, userID string) error {
	_, err := oc.client.UserAPI.DeleteUser(ctx, userID).Execute()
	if err != nil {
		oc.log.Errorf("Okta SDK: failed to delete user %s: %v", userID, err)
		return fmt.Errorf("failed to delete user %s: %w", userID, err)
	}
	return nil
}
