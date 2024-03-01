package services

import (
	"server/models"
	"server/repositories"
	"server/utils"
	"server/utils/validations"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(userRepository *repositories.UserRepository) *UserService {
	return &UserService{userRepository}
}

func (us *UserService) CreateUser(user *models.CreateUserRequest) error {
	err := validations.ValidateCreateUserRequest(user)
	if err != nil {
		return err
	}

	existingUser, err := us.userRepository.GetUserByEmail(user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return utils.EmailInUseError
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	userID, err := utils.GenerateUUID()
	if err != nil {
		return err
	}

	validUser := &models.User{
		ID:       userID,
		Name:     utils.CapitalizeText(user.Name),
		Email:    user.Email,
		Password: hashedPassword,
	}

	return us.userRepository.Create(validUser)
}

func (us *UserService) ListUsers() ([]models.User, error) {
	return us.userRepository.List()
}

func (us *UserService) UpdateUser(user *models.UpdateUserRequest, tokenID uuid.UUID) (*models.User, error) {
	if user.ID != tokenID {
		return nil, utils.UnauthorizedActionError
	}

	existingUser, err := us.userRepository.GetUserById(user.ID)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, utils.NotFoundError
	}

	if user.Name != nil {
		capitalizedName := utils.CapitalizeText(*user.Name)
		user.Name = &capitalizedName

		if len(*user.Name) < 2 {
			return nil, utils.InvalidNameError
		}
	}

	if user.Password != nil {
		passwordLength := len(*user.Password)

		if passwordLength < 8 || passwordLength > 128 {
			return nil, utils.PasswordLengthError
		}
	}

	updatedUser, err := us.userRepository.Update(user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (us *UserService) DeleteUser(id, tokenID uuid.UUID) error {
	if id != tokenID {
		return utils.UnauthorizedActionError
	}

	existingUser, err := us.userRepository.GetUserById(id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return utils.NotFoundError
	}

	return us.userRepository.Delete(id)
}

func (us *UserService) LoginUser(loginUserRequest *models.LoginUserResquest) (string, error) {
	existingUser, err := us.userRepository.GetUserByEmail(loginUserRequest.Email)
	if err != nil || existingUser == nil {
		return "", utils.InvalidEmailOrPasswordError
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(loginUserRequest.Password)); err != nil {
		return "", utils.InvalidEmailOrPasswordError
	}

	tokenString, err := utils.GenerateToken(existingUser.ID)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
