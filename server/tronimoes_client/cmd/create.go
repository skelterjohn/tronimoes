/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"

	spb "github.com/skelterjohn/tronimoes/server/proto"
	"github.com/skelterjohn/tronimoes/server/tronimoes_client/conn"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new tronimoes game.",
	Long:  `Create a new tronimoes game.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		c, err := conn.GetClient(ctx, serverAddress, useTLS)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating connection: %v", err)
			return
		}

		resp, err := c.CreateGame(ctx, &spb.CreateGameRequest{
			Discoverable: false,
			Private:      false,
			MinPlayers:   0,
			MaxPlayers:   0,
			PlayerId:     playerID,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating game: %v", err)
			return
		}

		fmt.Printf("Creating game with operation %s: ", resp.GetOperationId())
		for !resp.GetDone() {
			fmt.Print(".")
			time.Sleep(1 * time.Second)
			resp, err = c.GetOperation(ctx, &spb.GetOperationRequest{
				OperationId: resp.GetOperationId(),
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching operation game: %v", err)
				return
			}
		}
		fmt.Printf("done.\n")

		if resp.GetStatus() != spb.Operation_SUCCESS {
			fmt.Fprintf(os.Stderr, "Create game operation not successful: %v", err)
			return
		}

		g := &spb.Game{}
		if resp.GetPayload().GetTypeUrl() != "skelterjohn.tronimoes.Game" {
			fmt.Fprintf(os.Stderr, "Unexpected operation payload type %q", resp.GetPayload().GetTypeUrl())
			return
		}
		if err := proto.Unmarshal(resp.GetPayload().GetValue(), g); err != nil {
			fmt.Fprintf(os.Stderr, "Could not unmarshal operation payload: %v", err)
			return
		}
		fmt.Println(g)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
