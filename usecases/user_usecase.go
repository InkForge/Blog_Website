package usecases

import (
	"fmt"
	"regexp"
	"time"

	"github.com/InkForge/Blog_Website/domain"
)

type UserUseCase struct {
	UserRepo            domain.IUserRepository
	PasswordService     domain.IPasswordService
	JWTService          domain.IJWTService
	NotificationService domain.INotificationService
	BaseURL             string
}

//function to validate email

func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)

}

//function to generate verification email body

func generateVerificationEmailBody(verificationLink string) string {
	return fmt.Sprintf(`
    <html>
      <body style="font-family: Arial, sans-serif; line-height: 1.6;">
        <h2>Welcome!</h2>
        <p>Thanks for signing up. Please verify your email address by clicking the link below.</p>
        <p>This is a one-time link and may expire soon.</p>
        <p>
          <a href="%s" style="display: inline-block; padding: 10px 20px; background-color: #4CAF50;
          color: white; text-decoration: none; border-radius: 4px;">Verify Email</a>
        </p>
        <p>If you didn’t request this, feel free to ignore this email.</p>
        <p>— The Team</p>
      </body>
    </html>
  `, verificationLink)
}
func NewUserUseCase(repo domain.IUserRepository, ps domain.IPasswordService, jw domain.IJWTService, ns domain.INotificationService, bs string) *UserUseCase {
	return &UserUseCase{
		UserRepo:            repo,
		PasswordService:     ps,
		JWTService:          jw,
		NotificationService: ns,
		BaseURL:             bs,
	}
}

//register usecase

// Register handles user registration, supporting both traditional and OAuth-based flows
func (uc *UserUseCase) Register(input *domain.User, oauthUser *domain.User) (*domain.User, error) {

	var email string
	if oauthUser != nil {
		email = oauthUser.Email
	} else {
		email = input.Email
	}

	// email format validation
	if !validateEmail(email) {
		return nil, fmt.Errorf("%w", domain.ErrInvalidEmailFormat)
	}

	// check if email already exists
	count, err := uc.UserRepo.CountByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrDatabaseOperationFailed, err)
	}
	if count > 0 {
		return nil, fmt.Errorf("%w", domain.ErrEmailAlreadyExists)
	}

	// check if this is the first user
	totalUsers, err := uc.UserRepo.CountAll()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrDatabaseOperationFailed, err)
	}

	role := domain.RoleUser
	if totalUsers == 0 {
		role = domain.RoleAdmin
	}

	var hashedPassword *string
	if oauthUser == nil {
		hashed, err := uc.PasswordService.HashPassword(*input.Password)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", domain.ErrPasswordHashingFailed, err)
		}
		hashedPassword = &hashed
	}

	// construct user model
	newUser := domain.User{
		Role:           role,
		Username:       chooseNonEmpty(input.Username, oauthUser),
		FirstName:      chooseNonEmpty(input.FirstName, oauthUser),
		LastName:       chooseNonEmpty(input.LastName, oauthUser),
		Email:          email,
		Password:       hashedPassword,
		ProfilePicture: oauthUserPicture(oauthUser),
		Provider:       oauthUserProvider(oauthUser),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// save user to the database
	err = uc.UserRepo.CreateUser(newUser)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrUserCreationFailed, err)
	}

	// only send verification email if not registered via OAuth
	if oauthUser == nil {
		verificationToken, err := uc.JWTService.GenerateVerificationToken(fmt.Sprint(newUser.UserID))
		if err != nil {
			return nil, fmt.Errorf("%w: %v", domain.ErrTokenGenerationFailed, err)
		}
		verificationLink := fmt.Sprintf("%s/verify?token=%s", uc.BaseURL, verificationToken)
		emailBody := generateVerificationEmailBody(verificationLink)
		if err = uc.NotificationService.SendEmail(newUser.Email, "Verify Your Email Address", emailBody); err != nil {
			fmt.Println("email sending failed:", err)
		}
	}

	return &newUser, nil
}

// login usecase

// Login handles user login usecase
func (uc *UserUseCase) Login(input domain.User) (string, string, *domain.User, error) {
	// validate email format
	if !validateEmail(input.Email) {
		return "", "", nil, fmt.Errorf("%w", domain.ErrInvalidEmailFormat)
	}

	// find user by email
	user, err := uc.UserRepo.FindByEmail(input.Email)
	if err != nil {
		return "", "", nil, fmt.Errorf("%w: %v", domain.ErrInvalidCredentials, err)
	}

	// reject login if user is an OAuth user
	if user.Provider != "" {
		return "", "", nil, fmt.Errorf("%w", domain.ErrOAuthUserCannotLoginWithPassword)
	}

	// ensure email is verified
	if !uc.UserRepo.IsVerified(input.Email) {
		return "", "", nil, fmt.Errorf("%w", domain.ErrEmailNotVerified)
	}

	// compare password
	ok := uc.PasswordService.ComparePassword(*user.Password, *input.Password)
	if !ok {
		return "", "", nil, fmt.Errorf("%w", domain.ErrInvalidCredentials)
	}

	// generate access token
	accessToken, err := uc.JWTService.GenerateAccessToken(user.UserID, string(user.Role))
	if err != nil {
		return "", "", nil, fmt.Errorf("%w: %v", domain.ErrTokenGenerationFailed, err)
	}

	// generate refresh token
	refreshToken, err := uc.JWTService.GenerateRefreshToken(user.UserID, string(user.Role))
	if err != nil {
		return "", "", nil, fmt.Errorf("%w: %v", domain.ErrTokenGenerationFailed, err)
	}

	return accessToken, refreshToken, &user, nil
}

// helper functions
func chooseNonEmpty(field *string, oauthUser *domain.User) *string {
	if field != nil {
		return field
	}
	if oauthUser == nil {
		return nil
	}
	if field == nil && *oauthUser.FirstName != "" {
		return oauthUser.FirstName
	}
	return oauthUser.Name // fallback
}

func oauthUserPicture(oauthUser *domain.User) *string {
	if oauthUser == nil || *oauthUser.ProfilePicture == "" {
		return nil
	}
	return oauthUser.ProfilePicture
}

func oauthUserProvider(oauthUser *domain.User) string {
	if oauthUser == nil {
		return ""
	}
	return oauthUser.Provider
}
