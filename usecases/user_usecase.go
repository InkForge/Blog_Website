package usecases

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/InkForge/Blog_Website/domain"
)

type UserUseCase struct {
	UserRepo domain.IUserRepository
	PasswordService domain.IPasswordService
	JWTService domain.IJWTService
	NotificationService domain.INotificationService
	BaseURL string
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
func NewUserUseCase(repo domain.IUserRepository, ps domain.IPasswordService, jw domain.IJWTService, ns domain.INotificationService,bs string) *UserUseCase {
	return &UserUseCase{
		UserRepo:        repo,
		PasswordService: ps,
		JWTService:      jw,
		NotificationService: ns,
		BaseURL: bs,
		

	}
}

//register usecase

func (uc *UserUseCase) Register (input *domain.User)(*domain.User,error){
	//email format validation
	if !validateEmail(input.Email){
		return nil,errors.New("invalid email format")
	}


	//check if email already exits
	count,err:= uc.UserRepo.CountByEmail(input.Email)
	if err !=nil{
		return nil,errors.New("error while checking existing user")
	}

	if count >0{
		return nil,errors.New("email already exits")
	}

	//check if this is the first user
	totalUsers,err:=uc.UserRepo.CountAll()
	if err!=nil{
		return nil,errors.New("error while checking total users")

	}
	//hash password
	hashedPassword,err:=uc.PasswordService.HashPassword(*input.Password)
	if err !=nil{
		return nil,errors.New("failed to hash password")

	}
	//assign role 
	role:=domain.RoleUser
	if totalUsers==0{
		role=domain.RoleAdmin
	}

	//create the user model
	newUser := domain.User{
		Role:         role,
		Username:       input.Username,
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		Password: &hashedPassword,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	//save user to the database
	err=uc.UserRepo.CreateUser(newUser)
	if err !=nil{
		return nil,errors.New("failed to create user")

	}
	//email verification

	//generate verification token

	verificationToken,err:=uc.JWTService.GenerateVerificationToken(fmt.Sprint(newUser.UserID))
	if err !=nil{
		return nil,errors.New("failed to generate verification token")
	}
	verificationLink := fmt.Sprintf("%s/verify?token=%s", uc.BaseURL,verificationToken)
	emailBody := generateVerificationEmailBody(verificationLink)
	err= uc.NotificationService.SendEmail(newUser.Email, "Verify Your Email Address", emailBody)

	if err!=nil{
		fmt.Println("email sendign failed:",err)
	}

	return &newUser,nil



	
}

func (uc *UserUseCase) Login(input domain.User)(string,string,*domain.User,error){
	//validate email
	if !validateEmail(input.Email){
		return "","",nil,errors.New("invalid email format")
	}
	//find by email
	user,err:=uc.UserRepo.FindByEmail(input.Email)
	if err!=nil{
		return "","",nil,errors.New("invalid email and password")
	}

	//check if email is verifed
	ok:=uc.UserRepo.IsVerified(input.Email)
	if !ok{
		return "","",nil,errors.New("please verify your email beforre logging")
	}

	//compare password
	ok =uc.PasswordService.ComparePassword(*user.Password,*input.Password)
	if !ok{
		return "","",nil,errors.New("invalid email or password")
	}

	//generate access token
	accessToken,err:=uc.JWTService.GenerateAccessToken(user.UserID,string(input.Role))
	if err !=nil{
		return "","",nil,errors.New("failed to generate access token")

	}
	

	//generate refresh token

	refreshToken,err:= uc.JWTService.GenerateRefreshToken(user.UserID,string(user.Role))

	if err !=nil{
		return "","",nil,errors.New("failed to generate refresh token")
	}

	return accessToken,refreshToken,&user,nil
	
	


}