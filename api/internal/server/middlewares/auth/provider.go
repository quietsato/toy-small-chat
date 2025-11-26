package auth

import "github.com/go-chi/jwtauth/v5"

// Provider は chi 用の認証ミドルウェアを提供します
//
// usecase において [*jwtauth.JWTAuth] を直接使うことはないので、
// それを提供する interface は infrastructure 層で定義しています
type Provider interface {
	GetTokenAuthForMiddleware() *jwtauth.JWTAuth
}
