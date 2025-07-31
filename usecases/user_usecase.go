package usecases

import "github.com/InkForge/Blog_Website/domain"

type UserUseCase struct {
	UserRepo domain.IUserRepository
	PasswordService domain.IPasswordService
	JWTService domain.IJWTService
	NotificationService domain.INotificationService
}

func NewUserUseCase(repo domain.IUserRepository, ps domain.IPasswordService, jw domain.IJWTService ns domain.INotificationService) *UserUseCase {
	return &UserUseCase{
		UserRepo:        repo,
		PasswordService: ps,
		JWTService:      jw,
		NotificationService: ns,
	}
}

//register usecase

func (uc *UserUseCase) Register (input *domain.User)(*domain.User,error){
	//check if email already exits
	count,err:= uc.UserRepo.CountByEmail(*&input.Email)
	if err !=nil{
		return nil,errros
	}
	
}