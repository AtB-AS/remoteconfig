package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/atb-as/remoteconfig"
)

func main() {
	var credentials string
	var serviceAccountKey []byte
	projectIDFlag := flag.String("project", "", "Google Cloud project ID")
	timeout := flag.Int("timeout", 5, "timeout, in seconds")
	credsFlag := flag.String("credentials", "", "base64 encoded service account key (JSON)")
	flag.Parse()

	credentials = *credsFlag
	if credentials == "" {
		creds, ok := os.LookupEnv("SERVICE_ACCOUNT_KEY")
		if !ok {
			fmt.Fprintf(os.Stderr, "no credentials passed and SERVICE_ACCOUNT_KEY not set\n")
			os.Exit(1)
		}

		credentials = creds
	}

	c, err := base64.StdEncoding.DecodeString(credentials)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	serviceAccountKey = c

	var params map[string]string
	if err := json.NewDecoder(os.Stdin).Decode(&params); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*(time.Second))
	defer cancel()
	client := remoteconfig.NewClient(ctx, *projectIDFlag, serviceAccountKey)

	if err := client.SetDefaultValues(ctx, params); err != nil {
		fmt.Fprintf(os.Stderr, "SetDefaultValues: %v\n", err)
		os.Exit(1)
	}
}