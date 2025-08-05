package usecases

import (
	"context"
	"fmt"
	"regexp"
	"time"
	"unicode"

	"github.com/InkForge/Blog_Website/domain"
)

type AuthUseCase struct {
	UserRepo            domain.IUserRepository
	PasswordService     domain.IPasswordService
	JWTService          domain.IJWTService
	NotificationService domain.INotificationService
	BaseURL             string
	ContextTimeout      time.Duration
}

func NewAuthUseCase(repo domain.IUserRepository, ps domain.IPasswordService, jw domain.IJWTService, ns domain.INotificationService, bs string, timeout time.Duration) domain.IAuthUsecase {
	return &AuthUseCase{
		UserRepo:            repo,
		PasswordService:     ps,
		JWTService:          jw,
		NotificationService: ns,
		BaseURL:             bs,
		ContextTimeout:      timeout,
	}
}

// register usecase

// Register handles user registration, supporting both traditional and OAuth-based flows
func (uc *AuthUseCase) Register(ctx context.Context,input *domain.User, oauthUser *domain.User) (*domain.User, error) {
	ctx,cancel :=context.WithTimeout(ctx,uc.ContextTimeout)
	defer cancel()

	var email string
	if oauthUser != nil {
		email = oauthUser.Email
	} else {
		email = input.Email

		// check password strength (min 8 chars, at least one number and one letter)
		if !validatePasswordStrength(*input.Password) {
			return nil, fmt.Errorf("%w", domain.ErrWeakPassword)
		}
	}

	// email format validation
	if !validateEmail(email) {
		return nil, fmt.Errorf("%w", domain.ErrInvalidEmailFormat)
	}

	// check if email already exists
	count, err := uc.UserRepo.CountByEmail(ctx,email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrDatabaseOperationFailed, err)
	}
	if count > 0 {
		return nil, fmt.Errorf("%w", domain.ErrEmailAlreadyExists)
	}

	// check if this is the first user
	totalUsers, err := uc.UserRepo.CountAll(ctx)
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
	err = uc.UserRepo.CreateUser(ctx,&newUser)
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
func (uc *AuthUseCase) Login(ctx context.Context, input *domain.User) (string, string, *domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.ContextTimeout)
	defer cancel()

	// find user by email or username
	var user *domain.User
	var err error

	if validateEmail(input.Email) {
		user, err = uc.UserRepo.FindByEmail(ctx, input.Email)
	} else {
		user, err = uc.UserRepo.FindByUserName(ctx, *input.Username)
	}

	if err != nil {
		return "", "", nil, fmt.Errorf("%w: %v", domain.ErrInvalidCredentials, err)
	}

	// reject login if registered via OAuth
	if user.Provider != "" {
		return "", "", nil, fmt.Errorf("%w", domain.ErrOAuthUserCannotLoginWithPassword)
	}

	// check if email is verified
	isVerified, err := uc.UserRepo.IsEmailVerified(ctx, user.UserID)
	if err != nil {
		return "", "", nil, fmt.Errorf("%w", domain.ErrEmailVerficationFailed)
	}
	if !isVerified {
		return "", "", nil, fmt.Errorf("%w", domain.ErrEmailNotVerified)
	}

	// compare passwords
	if user.Password == nil || !uc.PasswordService.ComparePassword(*user.Password, *user.Password) {
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

	user.AccessToken = &accessToken
	user.RefreshToken = &refreshToken

	user.UpdatedAt = time.Now()

	// update the user (save the tokens into database)
	err = uc.UserRepo.UpdateTokens(ctx, user.UserID, accessToken, refreshToken)
	if err != nil {
		return "", "", nil, domain.ErrDatabaseOperationFailed
	}

	return accessToken, refreshToken, user, nil
}


//helper functions
func chooseNonEmpty(field *string, oauthUser *domain.User) *string {
	if field != nil {
		return field
	}
	if oauthUser == nil {
		return nil
	}
	if oauthUser.FirstName != nil && *oauthUser.FirstName != "" {
		return oauthUser.FirstName
	}
	return oauthUser.Name
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

//logout usecase
func (uc *AuthUseCase) Logout(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, uc.ContextTimeout)
	defer cancel()

	//check if empty
	if userID == "" {
		return fmt.Errorf("%w", domain.ErrInvalidUserID)
	}

	//find the users refresh token from db

	user, err := uc.UserRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w", domain.ErrDatabaseOperationFailed)
	}
	refreshToken := user.RefreshToken
	//call the revocation service
	if err := uc.JWTService.RevokeRefreshToken(*refreshToken); err != nil {
		return fmt.Errorf("%w: %v", domain.ErrTokenRevocationFailed, err)
	}
	return nil
}

//refresh token
func (uc *AuthUseCase) RefreshToken(ctx context.Context, userID string) (*string, *string, time.Duration, error) {
	emptyToken := ""
	//check emptyness
	if userID == "" {
		return &emptyToken, &emptyToken, 0, fmt.Errorf("%w", domain.ErrInvalidInput)
	}
	//find user
	user, err := uc.UserRepo.FindByID(ctx, userID)
	if err != nil {
		return &emptyToken, &emptyToken, 0, domain.ErrDatabaseOperationFailed
	}
	refreshToken := user.RefreshToken

	//validate the refresh token
	userID, role, err := uc.JWTService.ValidateRefreshToken(*refreshToken)
	if err != nil {
		return &emptyToken, &emptyToken, 0, domain.ErrTokenVerificationFailed
	}

	//generate new access token
	newAccessToken, err := uc.JWTService.GenerateAccessToken(userID, role)
	if err != nil {
		return &emptyToken, &emptyToken, 0, domain.ErrTokenGenerationFailed
	}

	//generate new  refresh token
	newRefreshToken, err := uc.JWTService.GenerateRefreshToken(userID, role)
	if err != nil {
		return &emptyToken, &emptyToken, 0, domain.ErrTokenGenerationFailed
	}

	user.AccessToken = &newAccessToken
	user.RefreshToken = &newRefreshToken
	ExpiredTime, err := uc.JWTService.GetAccessTokenRemaining(*user.AccessToken)
	if err != nil {
		return &emptyToken, &emptyToken, 0, domain.ErrGetTokenExpiryFailed
	}
	user.UpdatedAt = time.Now()

	//update the user
	err = uc.UserRepo.UpdateUser(ctx, user)
	if err != nil {
		return &emptyToken, &emptyToken, 0, domain.ErrDatabaseOperationFailed
	}
	//return access refresh duration

	return &newAccessToken, refreshToken, time.Duration(ExpiredTime), nil

}

//verify email
func (uc *AuthUseCase) VerifyEmail(ctx context.Context, token string) error {
	//check emptyness

	if token == "" {
		return fmt.Errorf("%w", domain.ErrInvalidToken)
	}

	//validate token
	userID, err := uc.JWTService.ValidateVerificationToken(token)
	if err != nil {
		return domain.ErrTokenVerificationFailed
	}

	//find  the user
	user, err := uc.UserRepo.FindByID(ctx, userID)
	if err != nil {
		return domain.ErrDatabaseOperationFailed
	}
	user.IsVerified = true
	user.UpdatedAt = time.Now()

	//update the user
	err = uc.UserRepo.UpdateUser(ctx, user)
	if err != nil {
		return domain.ErrDatabaseOperationFailed
	}

	return nil

}

//resend verification email
func (uc *AuthUseCase) ResendVerificationEmail(ctx context.Context, email string) error {
	//check validity
	if !validateEmail(email) {
		return domain.ErrInvalidEmailFormat
	}

	//find id
	user, err := uc.UserRepo.FindByEmail(ctx, email)
	userID := user.UserID
	if err != nil {
		return domain.ErrDatabaseOperationFailed
	}
	//generate verification token
	verificationToken, err := uc.JWTService.GenerateVerificationToken(userID)
	if err != nil {
		return domain.ErrTokenGenerationFailed
	}

	//send the verfication link
	verificationLink := fmt.Sprintf("%s/verify?token=%s", uc.BaseURL, verificationToken)
	emailBody := generateVerificationEmailBody(verificationLink)
	if err = uc.NotificationService.SendEmail(email, "Verify Your Email Address", emailBody); err != nil {
		fmt.Println("email sending failed:", err)
	}
	return nil
}

//Request password reset
func (uc *AuthUseCase) RequestPasswordReset(ctx context.Context, email string) error {
	ctx, cancel := context.WithTimeout(ctx, uc.ContextTimeout)
	defer cancel()

	if !validateEmail(email) {
		return fmt.Errorf("%w", domain.ErrInvalidEmailFormat)
	}

	user, err := uc.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("%w:%v", domain.ErrUserNotFound, err)

	}
	//generate a reset token
	resetToken, err := uc.JWTService.GeneratePasswordResetToken(fmt.Sprint(user.UserID))
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrTokenGenerationFailed, err)

	}
	//reset link
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", uc.BaseURL, resetToken)

	emailBody := fmt.Sprintf(`
    <html>
      <body style="font-family: Arial, sans-serif; line-height: 1.6;">
        <h2>Password Reset Requested</h2>
        <p>We received a request to reset your password. Click the link below to proceed. This link is one-time use and expires soon.</p>
        <p>
          <a href="%s" style="display: inline-block; padding: 10px 20px; background-color: #f39c12;
          color: white; text-decoration: none; border-radius: 4px;">Reset Password</a>
        </p>
        <p>If you didn't request this, you can safely ignore this email.</p>
        <p>— The Team</p>
      </body>
    </html>
    `, resetLink)

	if err := uc.NotificationService.SendEmail(user.Email, "reset your password", emailBody); err != nil {
		return err

	}
	return nil

}

//reset password
func (uc *AuthUseCase) ResetPassword(ctx context.Context, token string, newPassword string) error {
	//check emptyness
	if newPassword == "" || token == "" {
		return domain.ErrInvalidInput
	}

	//validate token
	userID, err := uc.JWTService.ValidatePasswordResetToken(token)
	if err != nil {
		return domain.ErrTokenVerificationFailed
	}

	//find user by userID
	user, err := uc.UserRepo.FindByID(ctx, userID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	//hashpassword
	hashedPassword, err := uc.PasswordService.HashPassword(newPassword)
	if err != nil {
		return domain.ErrPasswordHashingFailed
	}

	user.Password = &hashedPassword
	user.UpdatedAt = time.Now()

	//update user

	err = uc.UserRepo.UpdateUser(ctx, user)
	if err != nil {
		return domain.ErrDatabaseOperationFailed
	}
	return nil
}

//change password
func (uc *AuthUseCase) ChangePassword(ctx context.Context, userID string, oldPassword string, newPassword string) error {
	//check emptyness
	if oldPassword == "" || newPassword == "" {
		return domain.ErrInvalidInput
	}

	// check password strength
	if !validatePasswordStrength(newPassword) {
		return fmt.Errorf("%w", domain.ErrWeakPassword)
	}

	user, err := uc.UserRepo.FindByID(ctx, userID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	//verify old password

	ok := uc.PasswordService.ComparePassword(*user.Password, oldPassword)
	if !ok {
		return domain.ErrPasswordMismatch
	}

	//hash new password
	hashedPassword, err := uc.PasswordService.HashPassword(newPassword)
	if err != nil {
		return domain.ErrPasswordHashingFailed
	}

	//update user
	user.Password = &hashedPassword
	user.UpdatedAt = time.Now()

	err = uc.UserRepo.UpdateUser(ctx, user)
	if err != nil {
		return domain.ErrDatabaseOperationFailed
	}
	return nil

}


//forgot password
func (auc *AuthUseCase) ForgotPassword(ctx context.Context, email string) error {
	if !validateEmail(email) {
		return fmt.Errorf("%w", domain.ErrInvalidEmailFormat)
	}

	user, err := auc.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("%w:%v", domain.ErrUserNotFound, err)

	}
	//generate a reset token
	resetToken, err := auc.JWTService.GeneratePasswordResetToken(fmt.Sprint(user.UserID))
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrTokenGenerationFailed, err)

	}
	//reset link
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", auc.BaseURL, resetToken)

	emailBody := fmt.Sprintf(`
    <html>
      <body style="font-family: Arial, sans-serif; line-height: 1.6;">
        <h2>Password Reset Requested</h2>
        <p>We received a request to reset your password. Click the link below to proceed. This link is one-time use and expires soon.</p>
        <p>
          <a href="%s" style="display: inline-block; padding: 10px 20px; background-color: #f39c12;
          color: white; text-decoration: none; border-radius: 4px;">Reset Password</a>
        </p>
        <p>If you didn't request this, you can safely ignore this email.</p>
        <p>— The Team</p>
      </body>
    </html>
    `, resetLink)

	if err := auc.NotificationService.SendEmail(user.Email, "reset your password", emailBody); err != nil {
		fmt.Println("forgot password email send failed", err)

	}
	return nil

}


//function to validate email

func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)

}

// function to validate password strength 

func validatePasswordStrength(password string) bool {
	if len(password) < 8 {
		return false
	}

	hasLetter := false
	hasNumber := false

	for _, c := range password {
		switch {
		case unicode.IsLetter(c):
			hasLetter = true
		case unicode.IsNumber(c):
			hasNumber = true
		}
	}

	return hasLetter && hasNumber
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
