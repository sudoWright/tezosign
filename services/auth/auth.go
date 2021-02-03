package auth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/securecookie"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strings"
	"tezosign/common/apperrors"
	"tezosign/conf"
	"tezosign/models"
	"tezosign/types"
	"time"
)

type Auth struct {
	privateKey   *ecdsa.PrivateKey
	pubKey       *ecdsa.PublicKey
	secureCookie *securecookie.SecureCookie
	network      models.Network
}

const (
	authorizationHeader = "Authorization"
	UserAddressHeader   = "user_address"
	networkHeader       = "network"
)

func NewAuthProvider(authConf conf.Auth, network models.Network) (*Auth, error) {

	bt, err := hex.DecodeString(authConf.AuthKey)
	if err != nil {
		return nil, err
	}

	privKey, err := x509.ParseECPrivateKey(bt)
	if err != nil {
		return nil, err
	}

	// Hash keys should be at least 32 bytes long
	hashKey, err := hex.DecodeString(authConf.SessionHashKey)
	if err != nil {
		return nil, fmt.Errorf("Can not decode hash key: %s", err.Error())
	}

	// Block keys should be 16 bytes (AES-128) or 32 bytes (AES-256) long.
	// Shorter keys may weaken the encryption used.
	blockKey, err := hex.DecodeString(authConf.SessionBlockKey)
	if err != nil {
		return nil, fmt.Errorf("Can not decode hash key: %s", err.Error())
	}

	sc := securecookie.New(hashKey, blockKey)

	return &Auth{privateKey: privKey, pubKey: &privKey.PublicKey, secureCookie: sc, network: network}, nil
}

func (a *Auth) GenerateAuthTokens(address types.Address) (accessToken, refreshToken string, err error) {
	accessToken, err = a.generateAccessToken(address)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = a.generateRefreshToken(address)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (a *Auth) generateAccessToken(address types.Address) (accessToken string, err error) {
	if err = address.Validate(); err != nil {
		return "", err
	}

	// create the jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		UserAddressHeader: address.String(),
		networkHeader:     a.network,
		"exp":             time.Now().Add(time.Second * conf.TtlJWT).Unix(),
	})

	//json.Number(strconv.FormatInt(, 10))
	accessToken, err = token.SignedString(a.privateKey)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (a *Auth) generateRefreshToken(address types.Address) (token string, err error) {
	return uuid.NewV4().String(), nil
}

func (a *Auth) EncodeSessionCookie(data map[string]string) (encodedCookie string, err error) {
	encodedCookie, err = a.secureCookie.Encode("session", data)
	if err != nil {
		return "", err
	}

	return encodedCookie, nil
}

func (a *Auth) DecodeSessionCookie(cookie string) (map[string]string, error) {
	if cookie == "" {
		return nil, apperrors.New(apperrors.ErrBadAuth)
	}

	value := map[string]string{}
	err := a.secureCookie.Decode("session", cookie, &value)
	if err != nil {
		return nil, apperrors.New(apperrors.ErrBadAuth)
	}

	return value, nil
}

func (a *Auth) CheckSignatureAndGetUserAddress(r *http.Request) (string, error) {
	authHeader := strings.SplitN(r.Header.Get(authorizationHeader), " ", 2)
	if len(authHeader) != 2 {
		return "", apperrors.New(apperrors.ErrBadAuth)
	}

	token, claims, err := a.ParseAndCheckToken(authHeader[1])
	if err != nil {
		return "", apperrors.New(apperrors.ErrBadJwt)
	}

	if token == nil {
		return "", apperrors.New(apperrors.ErrBadJwt)
	}

	err = token.Claims.Valid()
	if err != nil {
		return "", apperrors.New(apperrors.ErrBadJwt)
	}

	if network, ok := claims[networkHeader].(string); !ok || network != string(a.network) {
		return "", apperrors.New(apperrors.ErrBadJwt)
	}

	userAddress, ok := claims[UserAddressHeader]
	if !ok || userAddress.(string) == "" {
		return "", apperrors.New(apperrors.ErrBadJwt)
	}

	return userAddress.(string), nil
}

func (a *Auth) ParseAndCheckToken(t string) (*jwt.Token, jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("Bad JWT method")
		}

		return a.pubKey, nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("Can not parse JWT token, %v", err)
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, nil, fmt.Errorf("JWT token is invalid")
	}

	return token, claims, nil
}
