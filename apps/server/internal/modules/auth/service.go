package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v5"
	"github.com/surajgoraicse/go-next-boilerplate/internal/common/logger"
	"github.com/surajgoraicse/go-next-boilerplate/internal/common/utils"
	"github.com/surajgoraicse/go-next-boilerplate/internal/config"
	db_sqlc "github.com/surajgoraicse/go-next-boilerplate/internal/db/sqlc"
	"go.uber.org/zap"
)

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserExists          = errors.New("user already exists")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrEmailNotVerified    = errors.New("email not verified")
	ErrEmailRequired       = errors.New("email is required")
	ErrInvalidToken        = errors.New("invalid token")
	ErrTokenExpired        = errors.New("verification link has expired, please request a new one")
	ErrTokenAlreadyUsed    = errors.New("this link has already been used")
	ErrUserAlreadyVerified = errors.New("user is already verified")
	ErrUserNotFound        = errors.New("user not found")
	ErrTooManyRequests     = errors.New("too many verification attempts, please try again later")
	ErrOTPInvalid          = errors.New("invalid OTP")
	ErrOTPExpired          = errors.New("OTP has expired")
	ErrOTPAlreadyUsed      = errors.New("OTP has already been used")
	ErrOTPAttemptsExceeded = errors.New("too many incorrect OTP attempts")
	ErrResetTokenInvalid   = errors.New("invalid or expired reset token")
	ErrResetTokenExpired   = errors.New("reset token has expired")
	ErrWeakPassword        = errors.New("password does not meet security requirements")
	ErrRateLimited         = errors.New("too many requests, please try again later")
)

type Service struct {
	queries *db_sqlc.Queries
	db      *pgxpool.Pool
	config  *config.Config
	logger  logger.Logger
}

func NewService(queries *db_sqlc.Queries, db *pgxpool.Pool, config *config.Config, logger logger.Logger) *Service {
	return &Service{
		queries: queries,
		db:      db,
		config:  config,
		logger:  logger,
	}
}

// setAuthCookies sets both auth cookies with the appropriate expiry times.
func (s *Service) setAuthCookies(c *echo.Context, tokens utils.AuthTokens) {
	accessTokenMaxAge := int(s.config.AccessTokenExpiry.Seconds())
	accessCookie := utils.NewSecureCookie("auth_token", tokens.AccessToken, accessTokenMaxAge, "/", s.config.AppEnv)
	c.SetCookie(accessCookie)

	refreshTokenMaxAge := int(s.config.RefreshTokenExpiry.Seconds())
	refreshCookie := utils.NewSecureCookie("refresh_token", tokens.RefreshToken, refreshTokenMaxAge, "/", s.config.AppEnv)
	c.SetCookie(refreshCookie)
}

// clearAuthCookies clears both auth cookies.
func (s *Service) clearAuthCookies(c *echo.Context) {
	c.SetCookie(utils.NewSecureCookie("auth_token", "", -1, "/", s.config.AppEnv))
	c.SetCookie(utils.NewSecureCookie("refresh_token", "", -1, "/", s.config.AppEnv))
}

// setVerificationSessionCookie sets the verification_session cookie with a expiry.
func (s *Service) setVerificationSessionCookie(c *echo.Context, sessionToken string, expiry time.Duration) {
	verificationCookie := utils.NewSecureCookie("verification_session", sessionToken, int(expiry.Seconds()), "/", s.config.AppEnv)
	c.SetCookie(verificationCookie)
}

// clearVerificationSessionCookie clears the verification_session cookie.
func (s *Service) clearVerificationSessionCookie(c *echo.Context) {
	c.SetCookie(utils.NewSecureCookie("verification_session", "", -1, "/", s.config.AppEnv))
}

// generateAndStoreVerificationToken generates a secure token, stores its hash in the DB,
// and returns the raw token string along with the new token's ID for use in exclusion logic.
func (s *Service) generateAndStoreVerificationToken(ctx context.Context, q *db_sqlc.Queries, userID pgtype.UUID) (string, pgtype.UUID, error) {
	return "", pgtype.UUID{}, nil
}

func (s *Service) deliverVerificationEmail(userID pgtype.UUID, email string, rawToken string) error {
	return nil
}

// generateAuthTokens creates both access and refresh tokens and persists the session
func (s *Service) generateAuthTokens(ctx context.Context, q *db_sqlc.Queries, userID pgtype.UUID, deviceID pgtype.UUID, payload utils.TokenPayload) (utils.AuthTokens, error) {
	// Generate JWT access token
	accessToken, err := utils.GenerateAccessToken(payload, s.config.AccessTokenExpiry, s.config.JwtSecret)
	if err != nil {
		s.logger.Error("failed to generate access token", zap.String("user_id", userID.String()), zap.Error(err))
		return utils.AuthTokens{}, err
	}

	// Generate refresh token
	rawRefresh, tokenHash, err := utils.GenerateSecureToken()
	if err != nil {
		s.logger.Error("failed to generate refresh token", zap.String("user_id", userID.String()), zap.Error(err))
		return utils.AuthTokens{}, err
	}

	// Upsert session in DB
	expiresAt := time.Now().Add(s.config.RefreshTokenExpiry)
	_, err = q.UpsertRefreshSession(ctx, db_sqlc.UpsertRefreshSessionParams{
		UserID:           userID,
		DeviceID:         deviceID,
		RefreshTokenHash: tokenHash,
		ExpiresAt:        pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		s.logger.Error("failed to upsert refresh session", zap.String("user_id", userID.String()), zap.Error(err))
		return utils.AuthTokens{}, err
	}

	return utils.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
	}, nil
}

// Register registers a new user
func (s *Service) Register(ctx context.Context, req RegisterRequest) (int, error) {
	// hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logger.Error("error hashing password", zap.Error(err))
		return http.StatusInternalServerError, err
	}

	// begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		s.logger.Error("error beginning transaction", zap.Error(err))
		return http.StatusInternalServerError, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	txq := s.queries.WithTx(tx)

	// create user
	insertParams := db_sqlc.InsertUserParams{
		Name:         req.Name,
		Email:        pgtype.Text{String: req.Email, Valid: true},
		PasswordHash: pgtype.Text{String: hashedPassword, Valid: true},
	}
	user, err := txq.InsertUser(ctx, insertParams)
	if err != nil {
		if utils.DbErrIsUniqueViolation(err) {
			return http.StatusConflict, errors.New("user already exists")
		}
		s.logger.Error("failed to create user in the databases", zap.String("email", req.Email), zap.Error(err))
		return http.StatusInternalServerError, err
	}
	// Send verification email
	rawToken, _, err := s.generateAndStoreVerificationToken(ctx, txq, user.ID)
	if err != nil {
		s.logger.Error("failed to initiate verification flow",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		return http.StatusInternalServerError, err
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.Error("failed to commit registration transaction",
			zap.String("email", req.Email),
			zap.Error(err))
		return http.StatusInternalServerError, err
	}

	if err := s.deliverVerificationEmail(user.ID, user.Email.String, rawToken); err != nil {
		s.logger.Error("failed to deliver verification email during registration",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		// We return Created even if email fails, because user exists and they can request resend verification email
	}

	return http.StatusCreated, nil

}

// SignIn logs in a user and returns auth tokens and device ID
func (s *Service) SignIn(ctx context.Context, req LoginRequest, deviceID pgtype.UUID) (utils.AuthTokens, pgtype.UUID, int, error) {
	// 1. check if user exists
	user, err := s.queries.GetUserByEmail(ctx, pgtype.Text{String: req.Email, Valid: true})
	if err != nil {
		if utils.DbErrIsNotFound(err) {
			return utils.AuthTokens{}, pgtype.UUID{}, http.StatusNotFound, ErrUserNotFound
		}
		s.logger.Error("database error during email lookup", zap.String("email", req.Email), zap.Error(err))
		return utils.AuthTokens{}, pgtype.UUID{}, http.StatusInternalServerError, err
	}
	// 2. check if user is verified
	if !user.EmailVerified {
		return utils.AuthTokens{}, pgtype.UUID{}, http.StatusUnauthorized, ErrEmailNotVerified
	}

	// 3. check if password is correct
	if !utils.ComparePassword(user.PasswordHash.String, req.Password) {
		return utils.AuthTokens{}, pgtype.UUID{}, http.StatusUnauthorized, ErrInvalidCredentials
	}

	// Start transaction for session management
	tx, err := s.db.Begin(ctx)
	if err != nil {
		s.logger.Error("error beginning transaction", zap.Error(err))
		return utils.AuthTokens{}, pgtype.UUID{}, http.StatusInternalServerError, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	txq := s.queries.WithTx(tx)

	// 4. Enforce device limit
	activeSessions, err := txq.CountActiveSessionsByUser(ctx, user.ID)
	if err != nil {
		s.logger.Error("failed to count active sessions", zap.Error(err))
		return utils.AuthTokens{}, pgtype.UUID{}, http.StatusInternalServerError, err
	}

	if int(activeSessions) >= s.config.MaxActiveDevices {
		// Revoke oldest session to make room
		oldest, err := txq.GetOldestActiveSessionByUser(ctx, user.ID)
		if err == nil {
			_, err = txq.RevokeSessionByID(ctx, oldest.ID)
			if err != nil {
				s.logger.Warn("failed to revoke oldest session", zap.Error(err))
			}
		}
	}

	// 5. Generate device ID if missing
	if !deviceID.Valid {
		newID, err := uuid.NewV7()
		if err != nil {
			s.logger.Error("failed to generate device id", zap.Error(err))
			return utils.AuthTokens{}, pgtype.UUID{}, http.StatusInternalServerError, err
		}
		deviceID = pgtype.UUID{Bytes: newID, Valid: true}
	}

	// 6. generate tokens and store session
	tokenPayload := utils.TokenPayload{
		UserID:          user.ID.Bytes,
		Email:           user.Email.String,
		Name:            user.Name,
		Role:            user.Role,
		IsEmailVerified: user.EmailVerified,
	}

	tokens, err := s.generateAuthTokens(ctx, txq, user.ID, deviceID, tokenPayload)
	if err != nil {
		s.logger.Error("error generating auth tokens", zap.Error(err))
		return utils.AuthTokens{}, pgtype.UUID{}, http.StatusInternalServerError, err
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.Error("error committing transaction", zap.Error(err))
		return utils.AuthTokens{}, pgtype.UUID{}, http.StatusInternalServerError, err
	}

	return tokens, deviceID, http.StatusOK, nil
}

// RefreshToken rotates the refresh token and returns a new pair of tokens
func (s *Service) RefreshToken(ctx context.Context, refreshToken string) (utils.AuthTokens, int, error) {
	// 1. Hash the incoming token
	tokenHash := utils.HashToken(refreshToken)

	// 2. Begin transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		s.logger.Error("error beginning transaction", zap.Error(err))
		return utils.AuthTokens{}, http.StatusInternalServerError, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	txq := s.queries.WithTx(tx)

	// 3. Find active session for this token hash
	session, err := txq.GetActiveSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		if utils.DbErrIsNotFound(err) {
			return utils.AuthTokens{}, http.StatusUnauthorized, ErrInvalidRefreshToken
		}
		s.logger.Error("failed to get session by token hash", zap.Error(err))
		return utils.AuthTokens{}, http.StatusInternalServerError, err
	}

	// 4. Get user details
	user, err := txq.GetUserByID(ctx, session.UserID)
	if err != nil {
		s.logger.Error("failed to get user for session", zap.String("user_id", session.UserID.String()), zap.Error(err))
		return utils.AuthTokens{}, http.StatusInternalServerError, err
	}

	// 5. Generate new tokens
	tokenPayload := utils.TokenPayload{
		UserID:          user.ID.Bytes,
		Email:           user.Email.String,
		Name:            user.Name,
		Role:            user.Role,
		IsEmailVerified: user.EmailVerified,
	}

	accessToken, err := utils.GenerateAccessToken(tokenPayload, s.config.AccessTokenExpiry, s.config.JwtSecret)
	if err != nil {
		return utils.AuthTokens{}, http.StatusInternalServerError, err
	}

	newRawRefresh, newTokenHash, err := utils.GenerateSecureToken()
	if err != nil {
		return utils.AuthTokens{}, http.StatusInternalServerError, err
	}

	// 6. Rotate session in DB
	newExpiresAt := time.Now().Add(s.config.RefreshTokenExpiry)
	_, err = txq.RotateRefreshToken(ctx, db_sqlc.RotateRefreshTokenParams{
		ID:               session.ID,
		RefreshTokenHash: newTokenHash,
		ExpiresAt:        pgtype.Timestamptz{Time: newExpiresAt, Valid: true},
	})
	if err != nil {
		s.logger.Error("failed to rotate refresh token", zap.String("session_id", session.ID.String()), zap.Error(err))
		return utils.AuthTokens{}, http.StatusInternalServerError, err
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.Error("error committing transaction", zap.Error(err))
		return utils.AuthTokens{}, http.StatusInternalServerError, err
	}

	return utils.AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: newRawRefresh,
	}, http.StatusOK, nil
}

// Logout revokes a specific device session
func (s *Service) Logout(ctx context.Context, userID pgtype.UUID, deviceID pgtype.UUID) error {
	// We need to find the session ID for this user and device
	sessions, err := s.queries.ListActiveSessionsByUser(ctx, userID)
	if err != nil {
		return err
	}

	for _, session := range sessions {
		if session.DeviceID == deviceID {
			_, err = s.queries.RevokeSessionByID(ctx, session.ID)
			return err
		}
	}

	return nil
}

// LogoutAll revokes all sessions for a user
func (s *Service) LogoutAll(ctx context.Context, userID pgtype.UUID) error {
	_, err := s.queries.RevokeAllSessionsByUser(ctx, userID)
	return err
}

// SendVerificationEmail sends a verification magic link.
// Only allow VerificationEmailRateLimit emails per day.
func (s *Service) SendVerificationEmail(ctx context.Context, email string) (int, error) {
	// 1. check user in db
	user, err := s.queries.GetUserByEmail(ctx, pgtype.Text{String: email, Valid: true})
	if err != nil {
		if utils.DbErrIsNotFound(err) {
			s.logger.Error("user not found", zap.String("email", email), zap.Error(err))
			return http.StatusBadRequest, ErrUserNotFound
		}
		s.logger.Error("error getting user by email", zap.String("email", email), zap.Error(err))
		return http.StatusInternalServerError, err
	}

	// 2. check user already verified - noting to do
	if user.EmailVerified {
		s.logger.Info("user already verified", zap.String("email", email))
		return http.StatusBadRequest, ErrUserAlreadyVerified
	}

	// 3. check rate limit : only allow VerificationEmailRateLimit tokens per day eg. 5 req per day
	count, err := s.queries.CountVerificationTokensSentRecently(ctx, db.CountVerificationTokensSentRecentlyParams{
		UserID:    user.ID,
		CreatedAt: pgtype.Timestamp{Time: time.Now().Add(-24 * time.Hour), Valid: true},
	})
	if err == nil && count >= int64(s.config.VerificationEmailRateLimit) {
		return http.StatusTooManyRequests, ErrTooManyRequests
	}

	// 4. generate and store verification token
	tx, err := s.db.Begin(ctx)
	if err != nil {
		s.logger.Error("failed to begin resend transaction", zap.Error(err))
		return http.StatusInternalServerError, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	txq := s.queries.WithTx(tx)
	rawToken, newTokenID, err := s.generateAndStoreVerificationToken(ctx, txq, user.ID)
	if err != nil {
		s.logger.Error("failed to generate/store new verification token inside transaction",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		return http.StatusInternalServerError, err
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.Error("failed to commit resend transaction", zap.Error(err))
		return http.StatusInternalServerError, err
	}

	// 5. send email
	if err := s.deliverVerificationEmail(user.ID, user.Email.String, rawToken); err != nil {
		return http.StatusInternalServerError, err
	}

	// 6.  Revoke any existing active tokens ONLY after successful delivery of the new one.
	// This prevents the "dead zone" if the new email fails to send.
	// We exclude the newly created token (newTokenID) so it remains active.
	if err := s.queries.RevokeActiveTokensByUserID(ctx, db.RevokeActiveTokensByUserIDParams{
		UserID: user.ID,
		ID:     newTokenID,
	}); err != nil {
		s.logger.Warn("failed to revoke old verification tokens after successful resend",
			zap.String("user_id", user.ID.String()),
			zap.Error(err))
		// We don't return an error here because the new email was successfully sent.
	}

	return http.StatusOK, nil
}
