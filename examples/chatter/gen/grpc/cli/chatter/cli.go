// Code generated by goa v2.0.0-wip, DO NOT EDIT.
//
// chatter gRPC client CLI support package
//
// Command:
// $ goa gen goa.design/goa/examples/chatter/design -o
// $(GOPATH)/src/goa.design/goa/examples/chatter

package cli

import (
	"flag"
	"fmt"
	"os"

	goa "goa.design/goa"
	chattersvcc "goa.design/goa/examples/chatter/gen/grpc/chatter/client"
	grpc "google.golang.org/grpc"
)

// UsageCommands returns the set of commands and sub-commands using the format
//
//    command (subcommand1|subcommand2|...)
//
func UsageCommands() string {
	return `chatter (login|echoer|listener|summary|history)
`
}

// UsageExamples produces an example of a valid invocation of the CLI tool.
func UsageExamples() string {
	return os.Args[0] + ` chatter login --user "username" --password "password"` + "\n" +
		""
}

// ParseEndpoint returns the endpoint and payload as specified on the command
// line.
func ParseEndpoint(cc *grpc.ClientConn, opts ...grpc.CallOption) (goa.Endpoint, interface{}, error) {
	var (
		chatterFlags = flag.NewFlagSet("chatter", flag.ContinueOnError)

		chatterLoginFlags        = flag.NewFlagSet("login", flag.ExitOnError)
		chatterLoginUserFlag     = chatterLoginFlags.String("user", "REQUIRED", "")
		chatterLoginPasswordFlag = chatterLoginFlags.String("password", "REQUIRED", "")

		chatterEchoerFlags     = flag.NewFlagSet("echoer", flag.ExitOnError)
		chatterEchoerTokenFlag = chatterEchoerFlags.String("token", "REQUIRED", "")

		chatterListenerFlags     = flag.NewFlagSet("listener", flag.ExitOnError)
		chatterListenerTokenFlag = chatterListenerFlags.String("token", "REQUIRED", "")

		chatterSummaryFlags     = flag.NewFlagSet("summary", flag.ExitOnError)
		chatterSummaryTokenFlag = chatterSummaryFlags.String("token", "REQUIRED", "")

		chatterHistoryFlags     = flag.NewFlagSet("history", flag.ExitOnError)
		chatterHistoryViewFlag  = chatterHistoryFlags.String("view", "", "")
		chatterHistoryTokenFlag = chatterHistoryFlags.String("token", "REQUIRED", "")
	)
	chatterFlags.Usage = chatterUsage
	chatterLoginFlags.Usage = chatterLoginUsage
	chatterEchoerFlags.Usage = chatterEchoerUsage
	chatterListenerFlags.Usage = chatterListenerUsage
	chatterSummaryFlags.Usage = chatterSummaryUsage
	chatterHistoryFlags.Usage = chatterHistoryUsage

	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return nil, nil, err
	}

	if len(os.Args) < flag.NFlag()+3 {
		return nil, nil, fmt.Errorf("not enough arguments")
	}

	var (
		svcn string
		svcf *flag.FlagSet
	)
	{
		svcn = os.Args[1+flag.NFlag()]
		switch svcn {
		case "chatter":
			svcf = chatterFlags
		default:
			return nil, nil, fmt.Errorf("unknown service %q", svcn)
		}
	}
	if err := svcf.Parse(os.Args[2+flag.NFlag():]); err != nil {
		return nil, nil, err
	}

	var (
		epn string
		epf *flag.FlagSet
	)
	{
		epn = os.Args[2+flag.NFlag()+svcf.NFlag()]
		switch svcn {
		case "chatter":
			switch epn {
			case "login":
				epf = chatterLoginFlags

			case "echoer":
				epf = chatterEchoerFlags

			case "listener":
				epf = chatterListenerFlags

			case "summary":
				epf = chatterSummaryFlags

			case "history":
				epf = chatterHistoryFlags

			}

		}
	}
	if epf == nil {
		return nil, nil, fmt.Errorf("unknown %q endpoint %q", svcn, epn)
	}

	// Parse endpoint flags if any
	if len(os.Args) > 2+flag.NFlag()+svcf.NFlag() {
		if err := epf.Parse(os.Args[3+flag.NFlag()+svcf.NFlag():]); err != nil {
			return nil, nil, err
		}
	}

	var (
		data     interface{}
		endpoint goa.Endpoint
		err      error
	)
	{
		switch svcn {
		case "chatter":
			c := chattersvcc.NewClient(cc, opts...)
			switch epn {
			case "login":
				endpoint = c.Login()
				data, err = chattersvcc.BuildLoginPayload(*chatterLoginUserFlag, *chatterLoginPasswordFlag)
			case "echoer":
				endpoint = c.Echoer()
				data, err = chattersvcc.BuildEchoerPayload(*chatterEchoerTokenFlag)
			case "listener":
				endpoint = c.Listener()
				data, err = chattersvcc.BuildListenerPayload(*chatterListenerTokenFlag)
			case "summary":
				endpoint = c.Summary()
				data, err = chattersvcc.BuildSummaryPayload(*chatterSummaryTokenFlag)
			case "history":
				endpoint = c.History()
				data, err = chattersvcc.BuildHistoryPayload(*chatterHistoryViewFlag, *chatterHistoryTokenFlag)
			}
		}
	}
	if err != nil {
		return nil, nil, err
	}

	return endpoint, data, nil
}

// chatterUsage displays the usage of the chatter command and its subcommands.
func chatterUsage() {
	fmt.Fprintf(os.Stderr, `The chatter service implements a simple client and server chat.
Usage:
    %s [globalflags] chatter COMMAND [flags]

COMMAND:
    login: Creates a valid JWT token for auth to chat.
    echoer: Echoes the message sent by the client.
    listener: Listens to the messages sent by the client.
    summary: Summarizes the chat messages sent by the client.
    history: Returns the chat messages sent to the server.

Additional help:
    %s chatter COMMAND --help
`, os.Args[0], os.Args[0])
}
func chatterLoginUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] chatter login -user STRING -password STRING

Creates a valid JWT token for auth to chat.
    -user STRING: 
    -password STRING: 

Example:
    `+os.Args[0]+` chatter login --user "username" --password "password"
`, os.Args[0])
}

func chatterEchoerUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] chatter echoer -token STRING

Echoes the message sent by the client.
    -token STRING: 

Example:
    `+os.Args[0]+` chatter echoer --token "Laudantium est consectetur ut."
`, os.Args[0])
}

func chatterListenerUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] chatter listener -token STRING

Listens to the messages sent by the client.
    -token STRING: 

Example:
    `+os.Args[0]+` chatter listener --token "Unde ea recusandae."
`, os.Args[0])
}

func chatterSummaryUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] chatter summary -token STRING

Summarizes the chat messages sent by the client.
    -token STRING: 

Example:
    `+os.Args[0]+` chatter summary --token "Ullam quia neque quis qui quia optio."
`, os.Args[0])
}

func chatterHistoryUsage() {
	fmt.Fprintf(os.Stderr, `%s [flags] chatter history -view STRING -token STRING

Returns the chat messages sent to the server.
    -view STRING: 
    -token STRING: 

Example:
    `+os.Args[0]+` chatter history --view "Minima voluptatum." --token "Est quibusdam voluptate nulla autem quisquam dolorum."
`, os.Args[0])
}
