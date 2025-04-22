package gateways

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
)

func GetCognitoUserEmail(tokenString string) (*string, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	cognito := cognitoidentityprovider.NewFromConfig(cfg)

	userOutput, err := cognito.GetUser(context.TODO(), &cognitoidentityprovider.GetUserInput{
		AccessToken: aws.String(tokenString),
	})

	if err != nil {
		return nil, err
	}

	// Para acessar atributos específicos, como o e-mail:
	for _, attr := range userOutput.UserAttributes {
		if *attr.Name == "email" {
			return attr.Value, nil
		}
	}

	return nil, errors.New("email not found")
}
