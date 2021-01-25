package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*(time.Second))
	defer cancel()
	client := remoteconfig.NewClient(ctx, *projectIDFlag, serviceAccountKey)

	args := flag.Args()
	keyvals := make([]string, 0, len(args)*2)
	for _, arg := range args {
		keyvals = append(keyvals, strings.SplitN(arg, "=", 2)...)
	}

	if err := client.SetDefaultValues(ctx, keyvals...); err != nil {
		fmt.Fprintf(os.Stderr, "SetDefaultValues: %v\n", err)
		os.Exit(1)
	}
}
