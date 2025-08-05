package usecases

import (

	"context"

	"time"

	"github.com/InkForge/Blog_Website/domain"
)

type UserUseCase struct {

	UserRepo            domain.IUserRepository	
	ContextTimeout      time.Duration
}


func NewUserUseCase(repo domain.IUserRepository, timeout time.Duration) domain.IUserUseCase {
	return &UserUseCase{
		UserRepo:            repo,
		ContextTimeout:      timeout,
	}
}

//get user by ID
func (uc *UserUseCase) GetUserByID( ctx context.Context ,userID string)(domain.User,error){
	ctx,cancel :=context.WithTimeout(ctx,uc.ContextTimeout)
	defer cancel()

	//check emptyness
	if userID==""{
		return domain.User{},domain.ErrInvalidUserID
	}
	//call the repo
	user,err:=uc.UserRepo.FindByID(ctx,userID)
	if err!=nil{
		return domain.User{},domain.ErrDatabaseOperationFailed
	}

	return *user,nil

}

//get users
func (uc *UserUseCase)GetUsers(ctx context.Context)([]domain.User,error){
	ctx,cancel :=context.WithTimeout(ctx,uc.ContextTimeout)
	defer cancel()
	//declare variable
	var users []domain.User

	//call the repo
	users,err:=uc.UserRepo.GetAllUsers(ctx)
	if err!=nil{
		return []domain.User{},domain.ErrDatabaseOperationFailed
	}
	return users,nil
}
//delete user
func (uc *UserUseCase)DeleteUserByID(ctx context.Context,userID string)(error){
	ctx,cancel :=context.WithTimeout(ctx,uc.ContextTimeout)
	defer cancel()

	if userID==""{
		return domain.ErrInvalidUserID
	}
	return uc.UserRepo.DeleteByID(ctx,userID)
}
//search users
func (uc *UserUseCase)SearchUsers(ctx context.Context,q string)([]domain.User,error){
	ctx,cancel :=context.WithTimeout(ctx,uc.ContextTimeout)
	defer cancel()

	return  uc.UserRepo.SearchUsers(ctx,q)
}
//user/me get my data is the same as getuserbyID
func (uc *UserUseCase) GetMyData(ctx context.Context, userID string) (*domain.User, error) {
	return uc.UserRepo.FindByID(ctx, userID)
}
//update profile
func (uc *UserUseCase)UpdateProfile(ctx context.Context,user *domain.User)(error){
	return uc.UserRepo.UpdateUser(ctx,user)

}

func (uc *UserUseCase) PromoteToAdmin(ctx context.Context, userID string) error {
	if userID == "" {
		return domain.ErrInvalidUserID
	}

	// check if user exists
	_, err := uc.UserRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// promote to admin
	return uc.UserRepo.UpdateRole(ctx, userID, "admin")
}

func (uc *UserUseCase) DemoteFromAdmin(ctx context.Context, userID string) error {
	if userID == "" {
		return domain.ErrInvalidUserID
	}

	// check if user exists
	_, err := uc.UserRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// demote to user
	return uc.UserRepo.UpdateRole(ctx, userID, "user")
}





