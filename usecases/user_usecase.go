package usecases



// 	"github.com/InkForge/Blog_Website/domain"
// )

// type UserUseCase struct {
// 	UserRepo            domain.IUserRepository
// 	PasswordService     domain.IPasswordService
// 	JWTService          domain.IJWTService
// 	NotificationService domain.INotificationService
// 	BaseURL             string
// 	ContextTimeout      time.Duration
// }

// //function to validate email

// func validateEmail(email string) bool {
// 	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
// 	return re.MatchString(email)

// }

// //function to generate verification email body

// func generateVerificationEmailBody(verificationLink string) string {
// 	return fmt.Sprintf(`
//     <html>
//       <body style="font-family: Arial, sans-serif; line-height: 1.6;">
//         <h2>Welcome!</h2>
//         <p>Thanks for signing up. Please verify your email address by clicking the link below.</p>
//         <p>This is a one-time link and may expire soon.</p>
//         <p>
//           <a href="%s" style="display: inline-block; padding: 10px 20px; background-color: #4CAF50;
//           color: white; text-decoration: none; border-radius: 4px;">Verify Email</a>
//         </p>
//         <p>If you didn’t request this, feel free to ignore this email.</p>
//         <p>— The Team</p>
//       </body>
//     </html>
//   `, verificationLink)
// }
// func NewUserUseCase(repo domain.IUserRepository, ps domain.IPasswordService, jw domain.IJWTService, ns domain.INotificationService, bs string, timeout time.Duration) domain.IUserUseCase {
// 	return &UserUseCase{
// 		UserRepo:            repo,
// 		PasswordService:     ps,
// 		JWTService:          jw,
// 		NotificationService: ns,
// 		BaseURL:             bs,
// 		ContextTimeout:      timeout,
// 	}
// }

// //register usecase

// // Register handles user registration, supporting both traditional and OAuth-based flows
// func (uc *UserUseCase) Register(ctx context.Context, input *domain.User, oauthUser *domain.User) (*domain.User, error) {
// 	ctx, cancel := context.WithTimeout(ctx, uc.ContextTimeout)
// 	defer cancel()

// 	var email string
// 	if oauthUser != nil {
// 		email = oauthUser.Email
// 	} else {
// 		email = input.Email
// 	}

// 	// email format validation
// 	if !validateEmail(email) {
// 		return nil, fmt.Errorf("%w", domain.ErrInvalidEmailFormat)
// 	}

// 	// check if email already exists
// 	count, err := uc.UserRepo.CountByEmail(ctx, email)
// 	if err != nil {
// 		return nil, fmt.Errorf("%w: %v", domain.ErrDatabaseOperationFailed, err)
// 	}
// 	if count > 0 {
// 		return nil, fmt.Errorf("%w", domain.ErrEmailAlreadyExists)
// 	}

// 	// check if this is the first user
// 	totalUsers, err := uc.UserRepo.CountAll(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("%w: %v", domain.ErrDatabaseOperationFailed, err)
// 	}

// 	role := domain.RoleUser
// 	if totalUsers == 0 {
// 		role = domain.RoleAdmin
// 	}

// 	var hashedPassword *string
// 	if oauthUser == nil {
// 		hashed, err := uc.PasswordService.HashPassword(*input.Password)
// 		if err != nil {
// 			return nil, fmt.Errorf("%w: %v", domain.ErrPasswordHashingFailed, err)
// 		}
// 		hashedPassword = &hashed
// 	}

// 	// construct user model
// 	newUser := domain.User{
// 		Role:           role,
// 		Username:       chooseNonEmpty(input.Username, oauthUser),
// 		FirstName:      chooseNonEmpty(input.FirstName, oauthUser),
// 		LastName:       chooseNonEmpty(input.LastName, oauthUser),
// 		Email:          email,
// 		Password:       hashedPassword,
// 		ProfilePicture: oauthUserPicture(oauthUser),
// 		Provider:       oauthUserProvider(oauthUser),
// 		CreatedAt:      time.Now(),
// 		UpdatedAt:      time.Now(),
// 	}

// 	// save user to the database
// 	err = uc.UserRepo.CreateUser(ctx, &newUser)
// 	if err != nil {
// 		return nil, fmt.Errorf("%w: %v", domain.ErrUserCreationFailed, err)
// 	}

// 	// only send verification email if not registered via OAuth
// 	if oauthUser == nil {
// 		verificationToken, err := uc.JWTService.GenerateVerificationToken(fmt.Sprint(newUser.UserID))
// 		if err != nil {
// 			return nil, fmt.Errorf("%w: %v", domain.ErrTokenGenerationFailed, err)
// 		}
// 		verificationLink := fmt.Sprintf("%s/verify?token=%s", uc.BaseURL, verificationToken)
// 		emailBody := generateVerificationEmailBody(verificationLink)
// 		if err = uc.NotificationService.SendEmail(newUser.Email, "Verify Your Email Address", emailBody); err != nil {
// 			fmt.Println("email sending failed:", err)
// 		}
// 	}

// 	return &newUser, nil
// }

// // login usecase

// // Login handles user login usecase
// func (uc *UserUseCase) Login(ctx context.Context, input domain.User) (string, string, *domain.User, error) {
// 	ctx, cancel := context.WithTimeout(ctx, uc.ContextTimeout)
// 	defer cancel()

// 	// validate email format
// 	if !validateEmail(input.Email) {
// 		return "", "", nil, fmt.Errorf("%w", domain.ErrInvalidEmailFormat)
// 	}

// 	// find user by email
// 	user, err := uc.UserRepo.FindByEmail(ctx, input.Email)
// 	if err != nil {
// 		return "", "", nil, fmt.Errorf("%w: %v", domain.ErrInvalidCredentials, err)
// 	}

// 	// reject login if user is an OAuth user
// 	if user.Provider != "" {
// 		return "", "", nil, fmt.Errorf("%w", domain.ErrOAuthUserCannotLoginWithPassword)
// 	}

// 	// ensure email is verified
// 	ok, err := uc.UserRepo.IsEmailVerified(ctx, input.UserID)
// 	if err != nil {
// 		return "", "", nil, fmt.Errorf("%w", domain.ErrEmailVerficationFailed)
// 	}
// 	if !ok {
// 		return "", "", nil, fmt.Errorf("%w", domain.ErrEmailNotVerified)
// 	}

// 	// compare password
// 	ok = uc.PasswordService.ComparePassword(*user.Password, *input.Password)
// 	if !ok {
// 		return "", "", nil, fmt.Errorf("%w", domain.ErrInvalidCredentials)
// 	}

// 	// generate access token
// 	accessToken, err := uc.JWTService.GenerateAccessToken(user.UserID, string(user.Role))
// 	if err != nil {
// 		return "", "", nil, fmt.Errorf("%w: %v", domain.ErrTokenGenerationFailed, err)
// 	}

// 	// generate refresh token
// 	refreshToken, err := uc.JWTService.GenerateRefreshToken(user.UserID, string(user.Role))
// 	if err != nil {
// 		return "", "", nil, fmt.Errorf("%w: %v", domain.ErrTokenGenerationFailed, err)
// 	}

// 	return accessToken, refreshToken, user, nil
// }

// // helper functions
// func chooseNonEmpty(field *string, oauthUser *domain.User) *string {
// 	if field != nil {
// 		return field
// 	}
// 	if oauthUser == nil {
// 		return nil
// 	}
// 	if oauthUser.FirstName != nil && *oauthUser.FirstName != "" {
// 		return oauthUser.FirstName
// 	}
// 	return oauthUser.Name
// }

// func oauthUserPicture(oauthUser *domain.User) *string {
// 	if oauthUser == nil || *oauthUser.ProfilePicture == "" {
// 		return nil
// 	}
// 	return oauthUser.ProfilePicture
// }

// func oauthUserProvider(oauthUser *domain.User) string {
// 	if oauthUser == nil {
// 		return ""
// 	}
// 	return oauthUser.Provider
// }

// //logout usecase
// func (uc *UserUseCase) Logout(refreshToken string) error {
// 	//check if empty
// 	if refreshToken == "" {
// 		return fmt.Errorf("%w", domain.ErrInvalidToken)
// 	}
// 	//call the revocation service
// 	if err := uc.JWTService.RevokeRefreshToken(refreshToken); err != nil {
// 		return fmt.Errorf("%w: %v", domain.ErrTokenRevocationFailed, err)
// 	}
// 	return nil
// }

// //forgot password
// func (uc *UserUseCase) ForgotPassword(ctx context.Context, email string) error {
// 	ctx, cancel := context.WithTimeout(ctx, uc.ContextTimeout)
// 	defer cancel()

// 	if !validateEmail(email) {
// 		return fmt.Errorf("%w", domain.ErrInvalidEmailFormat)
// 	}

// 	user, err := uc.UserRepo.FindByEmail(ctx, email)
// 	if err != nil {
// 		return fmt.Errorf("%w:%v", domain.ErrUserNotFound, err)

// 	}
// 	//generate a reset token
// 	resetToken, err := uc.JWTService.GeneratePasswordResetToken(fmt.Sprint(user.UserID))
// 	if err != nil {
// 		return fmt.Errorf("%w: %v", domain.ErrTokenGenerationFailed, err)

// 	}
// 	//reset link
// 	resetLink := fmt.Sprintf("%s/reset-password?token=%s", uc.BaseURL, resetToken)

// 	emailBody := fmt.Sprintf(`
//     <html>
//       <body style="font-family: Arial, sans-serif; line-height: 1.6;">
//         <h2>Password Reset Requested</h2>
//         <p>We received a request to reset your password. Click the link below to proceed. This link is one-time use and expires soon.</p>
//         <p>
//           <a href="%s" style="display: inline-block; padding: 10px 20px; background-color: #f39c12;
//           color: white; text-decoration: none; border-radius: 4px;">Reset Password</a>
//         </p>
//         <p>If you didn't request this, you can safely ignore this email.</p>
//         <p>— The Team</p>
//       </body>
//     </html>
//     `, resetLink)

// 	if err := uc.NotificationService.SendEmail(user.Email, "reset your password", emailBody); err != nil {
// 		fmt.Println("forgot password email send failed", err)

// 	}
// 	return nil

// }

// //reset password -consume the reset token and set a new password
// func (uc *UserUseCase) ResetPassword(ctx context.Context, resetToken, newPassword string) error {
// 	ctx, cancel := context.WithTimeout(ctx, uc.ContextTimeout)
// 	defer cancel()

// 	//check emptyness
// 	if resetToken == "" || newPassword == "" {
// 		return fmt.Errorf("%w", domain.ErrInvalidInput)

// 	}

// 	//validate token and extract userID
// 	userIDStr, err := uc.JWTService.ValidatePasswordResetToken(resetToken)
// 	if err != nil {
// 		return fmt.Errorf("%w:%v", domain.ErrInvalidToken, err)

// 	}
// 	//fetch user
// 	user, err := uc.UserRepo.FindByID(ctx, userIDStr)
// 	if err != nil {
// 		return fmt.Errorf("%w:%v", domain.ErrDatabaseOperationFailed, err)

// 	}

// 	//hash new password
// 	hashed, err := uc.PasswordService.HashPassword(newPassword)
// 	if err != nil {
// 		return fmt.Errorf("%w:%v", domain.ErrPasswordHashingFailed, err)
// 	}

// 	user.Password = &hashed
// 	user.UpdatedAt = time.Now()

// 	//persist update
// 	if err := uc.UserRepo.UpdateUser(ctx, user); err != nil {
// 		return fmt.Errorf("%w:%v", domain.ErrUserUpdateFailed, err)

// 	}
// 	return nil

// }

