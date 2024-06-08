package main

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"github.com/dermaddis/z2m_cli/actions"
	"github.com/dermaddis/z2m_cli/devices"
)

const usage = `
Usage: z2m [device] [actions]
`

func main() {
	// 1: _; 2: device(s); 3: action(s)
	if len(os.Args) < 3 {
		fmt.Print(usage)
		os.Exit(1)
	}
	devicesStr := os.Args[1]
	actionsStr := os.Args[2]

	if len(os.Args) > 3 {
		// The rest of the args should be included in actionsStr.
		// This can happen when you put a space after the comma for multiple actions.
		// `on, [red]` (the [] part is now Args[3])
		for _, arg := range os.Args[3:] {
			actionsStr = actionsStr + " " + arg
		}
	}

	devices, err := devices.Parse(devicesStr)
	if err != nil {
		fmt.Println("❌ | ", err)
		os.Exit(1)
	}
	actions, err := actions.Parse(actionsStr)
	if err != nil {
		fmt.Println("❌ | ", err)
		os.Exit(1)
	}

	client, err := NewFirestoreClient()
	if err != nil {
		fmt.Println("❌ | Could not init firestore", err)
		os.Exit(1)
	}
	deviceCollection := client.Collection("state")

	for _, device := range devices {
		updates := make([]firestore.Update, 0, len(actions))

		for _, action := range actions {
			updates = append(updates, firestore.Update{
				Path:  action.Path,
				Value: action.Value,
			})
		}
		_, err := deviceCollection.Doc(string(device)).Update(context.Background(), updates)
		if err != nil {
			panic(fmt.Errorf("could not update %s: %w", device, err))
		}
	}
	fmt.Println("✅")
}

func NewFirestoreClient() (*firestore.Client, error) {
	bgCtx := context.Background()

	credFile := "/usr/local/bin/z2m/cred.json"
	options := option.WithCredentialsFile(credFile)

	app, err := firebase.NewApp(bgCtx, &firebase.Config{}, options)
	if err != nil {
		return nil, err
	}

	fClient, err := app.Firestore(bgCtx)
	if err != nil {
		return nil, err
	}

	return fClient, nil
}
