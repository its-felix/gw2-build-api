//go:build lambda

package main

import "github.com/aws/aws-lambda-go/lambdaurl"

func main() {
	lambdaurl.Start(setupMux())
}
