package cmd

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go-v2/service/ecr"

	"github.com/docker/cli/cli/config"
	"github.com/docker/docker/api/types/registry"

	cliTypes "github.com/docker/cli/cli/config/types"
	dockerclient "github.com/docker/docker/client"
	awspkg "github.com/surajincloud/ecrctl/pkg/awsutil"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to ECR",
	Long: `Login to ECR.
	For example:
	ecrctl login
    `,
	RunE: Login,
}

func Login(cmd *cobra.Command, args []string) error {

	ctx := context.Background()

	// read flag values
	region, err := cmd.Flags().GetString("region")
	if err != nil {
		fmt.Println("no region", err)
	}

	// aws config
	cfg, err := awspkg.GetAWSConfig(ctx, region)
	if err != nil {
		log.Fatal(err)
	}

	// Create an ECR client
	client := ecr.NewFromConfig(cfg)

	// Get the login token from ECR
	loginOutput, err := client.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		fmt.Println("Failed to get authorization token:", err)
		return err
	}

	authData := loginOutput.AuthorizationData[0]
	endpoint := *authData.ProxyEndpoint
	token := *authData.AuthorizationToken

	decodedToken, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		fmt.Println("Failed to decode authorization token:", err)
		return err
	}

	password := strings.Split(string(decodedToken), ":")[1]

	// Create a Docker client
	dockerClient, err := dockerclient.NewClientWithOpts()
	if err != nil {
		fmt.Println("Failed to create Docker client:", err)
		return err
	}

	// Authenticate with the ECR repository
	authConfig := registry.AuthConfig{
		Username:      "AWS",
		Password:      password,
		ServerAddress: endpoint,
	}
	_, err = dockerClient.RegistryLogin(ctx, authConfig)
	if err != nil {
		fmt.Println("Failed to log in to ECR:", err)
		return err
	}

	// writing credentials to local file (~/.docker/config.json)
	var configFileErr io.Writer
	localConfigFile := config.LoadDefaultConfigFile(configFileErr)
	localConfigFile.AuthConfigs[endpoint] = cliTypes.AuthConfig{
		Username:      "AWS",
		Password:      password,
		ServerAddress: endpoint,
	}

	err = localConfigFile.Save()
	if err != nil {
		fmt.Println("Failed saving to file")
		return err
	}
	fmt.Println("Logged in to ECR successfully")
	return nil
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.PersistentFlags().String("region", "", "region")
}
