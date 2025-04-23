package gateways

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type JWKS struct {
	Keys []struct {
		Kty string `json:"kty"`
		Kid string `json:"kid"`
		Use string `json:"use"`
		N   string `json:"n"`
		E   string `json:"e"`
		Alg string `json:"alg"`
	} `json:"keys"`
}

type Auth struct {
	userPoolId     string
	region         string
	publicKeys     JWKS
	publicKeysMap  map[string]*rsa.PublicKey
	foundPublicKey *rsa.PublicKey
}

func NewAuthGateway(userPoolId string, region string) *Auth {
	return &Auth{
		userPoolId: userPoolId,
		region:     region,
	}
}

func (a *Auth) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	if a.foundPublicKey == nil {
		if err := a.getPublicKeys(); err != nil {
			return nil, err
		}
	}

	err := a.convertPublicKeysToSet()

	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("método de assinatura inválido: %v", token.Header["alg"])
		}
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid não encontrado no header")
		}

		key, err := a.getPublicKey(kid)

		if err != nil {
			return nil, err
		}

		return key, err
	})

	if err != nil {
		log.Printf("Erro ao validar token: %v", err)
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return nil, errors.New("Não autorizado")
	}

	return &claims, nil
}

func (a *Auth) getPublicKeys() error {
	jwksURL := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", a.region, a.userPoolId)
	fmt.Println("jwksURL", jwksURL)

	resp, err := http.Get(jwksURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("falha ao buscar JWKS, status: %d", resp.StatusCode)
	}

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return err
	}

	fmt.Println("jwks", jwks)

	a.publicKeys = jwks

	return nil
}

func (a *Auth) convertPublicKeysToSet() error {
	publicKeyMap := make(map[string]*rsa.PublicKey)
	for _, key := range a.publicKeys.Keys {
		if key.Kty == "RSA" {
			nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
			if err != nil {
				return err
			}
			eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
			if err != nil {
				return err
			}

			n := new(big.Int).SetBytes(nBytes)
			e := new(big.Int).SetBytes(eBytes)

			publicKey := &rsa.PublicKey{
				N: n,
				E: int(e.Int64()),
			}
			publicKeyMap[key.Kid] = publicKey
		}
	}

	fmt.Println("publicKeyMap", publicKeyMap)

	a.publicKeysMap = publicKeyMap

	return nil
}

func (a *Auth) getPublicKey(kid string) (*rsa.PublicKey, error) {
	key, ok := a.publicKeysMap[kid]
	if !ok {
		return nil, fmt.Errorf("chave pública não encontrada para kid: %s", kid)
	}
	return key, nil
}
