package useCase

import "jamlink-backend/internal/modules/auth/domain/token"

type DisconnectUserUseCase struct {
	tokenRepo token.TokenRepository
}

type DisconnectUserInput struct {
	RefreshToken string
}

func NewDisconnectUserUseCase(tokenRepo token.TokenRepository) *DisconnectUserUseCase {
	return &DisconnectUserUseCase{tokenRepo}
}

func (uc *DisconnectUserUseCase) Execute(input *DisconnectUserInput) error {
	foundRefreshToken, err := uc.tokenRepo.FindByToken(input.RefreshToken)

	if err != nil {
		return err
	}

	err = uc.tokenRepo.DeleteUserTokens(foundRefreshToken.UserID)

	if err != nil {
		return err
	}

	return nil
}
