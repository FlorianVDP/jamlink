package userUseCase

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	tokenDomain "jamlink-backend/internal/modules/auth/domain/token"
	userDomain "jamlink-backend/internal/modules/auth/domain/user"
	"jamlink-backend/internal/modules/auth/mocks"
	"jamlink-backend/internal/shared/security"
	"testing"
	"time"
)

func TestRefreshToken_Success(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	userRepo := new(mocks.MockUserRepository)
	tokenRepo := new(mocks.MockTokenRepository)

	fakeUser, err := userDomain.CreateUser("test@example.com", "hashedpassword", "fr-FR", "local")
	if err != nil {
		t.Fatal(err)
	}
	fakeUser.Verification.IsVerified = true

	const expiringTimeForRefreshToken = time.Hour * 24 * 7
	refreshToken := "valid_refresh_token"
	newRefreshToken := "new.refresh.token"
	newAccessToken := "new.jwt.token"

	timeNow := time.Now()
	existingToken := &tokenDomain.Token{
		ID:        uuid.New(),
		UserID:    fakeUser.ID,
		Token:     refreshToken,
		CreatedAt: timeNow.Add(-1 * time.Hour),
		ExpiresAt: timeNow.Add(expiringTimeForRefreshToken - time.Hour),
	}

	tokenRepo.On("FindByToken", refreshToken).Return(existingToken, nil)
	mockSecurity.On("GetJWTInfo", refreshToken).Return(fakeUser.ID, nil)
	userRepo.On("FindByID", fakeUser.ID).Return(fakeUser, nil)
	mockSecurity.On("GenerateJWT", &fakeUser.ID, (*string)(nil), time.Minute*15, "login", fakeUser.Verification.IsVerified).Return(newAccessToken, nil)
	mockSecurity.On("GenerateJWT", &fakeUser.ID, (*string)(nil), expiringTimeForRefreshToken, "refresh_token", fakeUser.Verification.IsVerified).Return(newRefreshToken, nil)
	tokenRepo.On("Create", mock.MatchedBy(func(token *tokenDomain.Token) bool {
		return token.Token == newRefreshToken && token.UserID == fakeUser.ID
	})).Return(nil)
	tokenRepo.On("DeleteByID", existingToken.ID).Return(nil)

	usecase := NewRefreshTokenUseCase(mockSecurity, userRepo, tokenRepo)

	// Act
	output, err := usecase.Execute(RefreshTokenInput{RefreshToken: refreshToken})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, newAccessToken, output.Token)
	assert.Equal(t, newRefreshToken, output.RefreshToken)

	mockSecurity.AssertExpectations(t)
	userRepo.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	userRepo := new(mocks.MockUserRepository)
	tokenRepo := new(mocks.MockTokenRepository)

	refreshToken := "invalid_token"

	tokenRepo.On("FindByToken", refreshToken).Return(nil, tokenDomain.ErrTokenExpired)

	usecase := NewRefreshTokenUseCase(mockSecurity, userRepo, tokenRepo)

	// Act
	output, err := usecase.Execute(RefreshTokenInput{RefreshToken: refreshToken})

	// Assert
	assert.ErrorIs(t, err, tokenDomain.ErrTokenExpired)
	assert.Nil(t, output)

	tokenRepo.AssertExpectations(t)
	// Les autres mocks ne devraient pas être appelés
}

func TestRefreshToken_ExpiredToken(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	userRepo := new(mocks.MockUserRepository)
	tokenRepo := new(mocks.MockTokenRepository)

	refreshToken := "expired_token"
	userID := uuid.New()

	// Token expiré (temps d'expiration dans le passé)
	expiredToken := &tokenDomain.Token{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     refreshToken,
		CreatedAt: time.Now().Add(-48 * time.Hour),
		ExpiresAt: time.Now().Add(-1 * time.Hour), // Expiré il y a une heure
	}

	tokenRepo.On("FindByToken", refreshToken).Return(expiredToken, nil)

	usecase := NewRefreshTokenUseCase(mockSecurity, userRepo, tokenRepo)

	// Act
	output, err := usecase.Execute(RefreshTokenInput{RefreshToken: refreshToken})

	// Assert
	assert.ErrorIs(t, err, tokenDomain.ErrTokenExpired)
	assert.Nil(t, output)

	tokenRepo.AssertExpectations(t)
}

func TestRefreshToken_GetJWTInfoFails(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	userRepo := new(mocks.MockUserRepository)
	tokenRepo := new(mocks.MockTokenRepository)

	refreshToken := "jwt_error_token"
	userID := uuid.New()

	validToken := &tokenDomain.Token{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     refreshToken,
		CreatedAt: time.Now().Add(-1 * time.Hour),
		ExpiresAt: time.Now().Add(23 * time.Hour),
	}

	tokenRepo.On("FindByToken", refreshToken).Return(validToken, nil)
	mockSecurity.On("GetJWTInfo", refreshToken).Return(uuid.Nil, security.ErrInvalidToken)

	usecase := NewRefreshTokenUseCase(mockSecurity, userRepo, tokenRepo)

	// Act
	output, err := usecase.Execute(RefreshTokenInput{RefreshToken: refreshToken})

	// Assert
	assert.ErrorIs(t, err, security.ErrInvalidToken)
	assert.Nil(t, output)

	tokenRepo.AssertExpectations(t)
	mockSecurity.AssertExpectations(t)
}

func TestRefreshToken_UserNotFound(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	userRepo := new(mocks.MockUserRepository)
	tokenRepo := new(mocks.MockTokenRepository)

	refreshToken := "user_not_found_token"
	userID := uuid.New()

	validToken := &tokenDomain.Token{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     refreshToken,
		CreatedAt: time.Now().Add(-1 * time.Hour),
		ExpiresAt: time.Now().Add(23 * time.Hour),
	}

	tokenRepo.On("FindByToken", refreshToken).Return(validToken, nil)
	mockSecurity.On("GetJWTInfo", refreshToken).Return(userID, nil)
	userRepo.On("FindByID", userID).Return(nil, userDomain.ErrUserNotFound)

	usecase := NewRefreshTokenUseCase(mockSecurity, userRepo, tokenRepo)

	// Act
	output, err := usecase.Execute(RefreshTokenInput{RefreshToken: refreshToken})

	// Assert
	assert.ErrorIs(t, err, userDomain.ErrUserNotFound)
	assert.Nil(t, output)

	tokenRepo.AssertExpectations(t)
	mockSecurity.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestRefreshToken_GenerateJWTFails(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	userRepo := new(mocks.MockUserRepository)
	tokenRepo := new(mocks.MockTokenRepository)

	refreshToken := "jwt_gen_error_token"
	userID := uuid.New()

	fakeUser, err := userDomain.CreateUser("test@example.com", "hashedpassword", "fr-FR", "local")
	if err != nil {
		t.Fatal(err)
	}
	fakeUser.ID = userID
	fakeUser.Verification.IsVerified = true

	validToken := &tokenDomain.Token{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     refreshToken,
		CreatedAt: time.Now().Add(-1 * time.Hour),
		ExpiresAt: time.Now().Add(23 * time.Hour),
	}

	tokenRepo.On("FindByToken", refreshToken).Return(validToken, nil)
	mockSecurity.On("GetJWTInfo", refreshToken).Return(userID, nil)
	userRepo.On("FindByID", userID).Return(fakeUser, nil)
	mockSecurity.On("GenerateJWT", &userID, (*string)(nil), time.Minute*15, "login", fakeUser.Verification.IsVerified).Return("", security.ErrJWTGeneration)

	usecase := NewRefreshTokenUseCase(mockSecurity, userRepo, tokenRepo)

	// Act
	output, err := usecase.Execute(RefreshTokenInput{RefreshToken: refreshToken})

	// Assert
	assert.ErrorIs(t, err, security.ErrJWTGeneration)
	assert.Nil(t, output)

	tokenRepo.AssertExpectations(t)
	mockSecurity.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestRefreshToken_GenerateRefreshTokenFails(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	userRepo := new(mocks.MockUserRepository)
	tokenRepo := new(mocks.MockTokenRepository)

	refreshToken := "refresh_gen_error_token"
	userID := uuid.New()

	fakeUser, err := userDomain.CreateUser("test@example.com", "hashedpassword", "fr-FR", "local")
	if err != nil {
		t.Fatal(err)
	}
	fakeUser.ID = userID
	fakeUser.Verification.IsVerified = true

	validToken := &tokenDomain.Token{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     refreshToken,
		CreatedAt: time.Now().Add(-1 * time.Hour),
		ExpiresAt: time.Now().Add(23 * time.Hour),
	}

	tokenRepo.On("FindByToken", refreshToken).Return(validToken, nil)
	mockSecurity.On("GetJWTInfo", refreshToken).Return(userID, nil)
	userRepo.On("FindByID", userID).Return(fakeUser, nil)
	mockSecurity.On("GenerateJWT", &userID, (*string)(nil), time.Minute*15, "login", fakeUser.Verification.IsVerified).Return("new_access_token", nil)
	mockSecurity.On("GenerateJWT", &userID, (*string)(nil), time.Hour*24*7, "refresh_token", fakeUser.Verification.IsVerified).Return("", security.ErrJWTGeneration)

	usecase := NewRefreshTokenUseCase(mockSecurity, userRepo, tokenRepo)

	// Act
	output, err := usecase.Execute(RefreshTokenInput{RefreshToken: refreshToken})

	// Assert
	assert.ErrorIs(t, err, security.ErrJWTGeneration)
	assert.Nil(t, output)

	tokenRepo.AssertExpectations(t)
	mockSecurity.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestRefreshToken_CreateTokenFails(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	userRepo := new(mocks.MockUserRepository)
	tokenRepo := new(mocks.MockTokenRepository)

	refreshToken := "create_token_error"
	userID := uuid.New()
	newRefreshToken := "new.refresh.token"

	fakeUser, err := userDomain.CreateUser("test@example.com", "hashedpassword", "fr-FR", "local")
	if err != nil {
		t.Fatal(err)
	}
	fakeUser.ID = userID
	fakeUser.Verification.IsVerified = true

	validToken := &tokenDomain.Token{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     refreshToken,
		CreatedAt: time.Now().Add(-1 * time.Hour),
		ExpiresAt: time.Now().Add(23 * time.Hour),
	}

	tokenRepo.On("FindByToken", refreshToken).Return(validToken, nil)
	mockSecurity.On("GetJWTInfo", refreshToken).Return(userID, nil)
	userRepo.On("FindByID", userID).Return(fakeUser, nil)
	mockSecurity.On("GenerateJWT", &userID, (*string)(nil), time.Minute*15, "login", fakeUser.Verification.IsVerified).Return("new_access_token", nil)
	mockSecurity.On("GenerateJWT", &userID, (*string)(nil), time.Hour*24*7, "refresh_token", fakeUser.Verification.IsVerified).Return(newRefreshToken, nil)
	tokenRepo.On("Create", mock.MatchedBy(func(token *tokenDomain.Token) bool {
		return token.Token == newRefreshToken && token.UserID == userID
	})).Return(tokenDomain.ErrTokenCreationFailed)

	usecase := NewRefreshTokenUseCase(mockSecurity, userRepo, tokenRepo)

	// Act
	output, err := usecase.Execute(RefreshTokenInput{RefreshToken: refreshToken})

	// Assert
	assert.ErrorIs(t, err, tokenDomain.ErrTokenCreationFailed)
	assert.Nil(t, output)

	tokenRepo.AssertExpectations(t)
	mockSecurity.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestRefreshToken_DeleteTokenFails(t *testing.T) {
	// Arrange
	mockSecurity := new(mocks.MockSecurityService)
	userRepo := new(mocks.MockUserRepository)
	tokenRepo := new(mocks.MockTokenRepository)

	refreshToken := "delete_token_error"
	userID := uuid.New()
	newRefreshToken := "new.refresh.token"

	fakeUser, err := userDomain.CreateUser("test@example.com", "hashedpassword", "fr-FR", "local")
	if err != nil {
		t.Fatal(err)
	}
	fakeUser.ID = userID
	fakeUser.Verification.IsVerified = true

	validToken := &tokenDomain.Token{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     refreshToken,
		CreatedAt: time.Now().Add(-1 * time.Hour),
		ExpiresAt: time.Now().Add(23 * time.Hour),
	}

	tokenRepo.On("FindByToken", refreshToken).Return(validToken, nil)
	mockSecurity.On("GetJWTInfo", refreshToken).Return(userID, nil)
	userRepo.On("FindByID", userID).Return(fakeUser, nil)
	mockSecurity.On("GenerateJWT", &userID, (*string)(nil), time.Minute*15, "login", fakeUser.Verification.IsVerified).Return("new_access_token", nil)
	mockSecurity.On("GenerateJWT", &userID, (*string)(nil), time.Hour*24*7, "refresh_token", fakeUser.Verification.IsVerified).Return(newRefreshToken, nil)
	tokenRepo.On("Create", mock.MatchedBy(func(token *tokenDomain.Token) bool {
		return token.Token == newRefreshToken && token.UserID == userID
	})).Return(nil)
	tokenRepo.On("DeleteByID", validToken.ID).Return(tokenDomain.ErrTokenDeletionFailed)

	usecase := NewRefreshTokenUseCase(mockSecurity, userRepo, tokenRepo)

	// Act
	output, err := usecase.Execute(RefreshTokenInput{RefreshToken: refreshToken})

	// Assert
	assert.ErrorIs(t, err, tokenDomain.ErrTokenDeletionFailed)
	assert.Nil(t, output)

	tokenRepo.AssertExpectations(t)
	mockSecurity.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}
