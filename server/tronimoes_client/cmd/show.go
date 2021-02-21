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

	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"

	spb "github.com/skelterjohn/tronimoes/server/proto"
	"github.com/skelterjohn/tronimoes/server/tronimoes_client/conn"
)

var (
	gameID string
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current board for a game",
	Long:  `Show the current board for a game.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		c, err := conn.GetClient(ctx, serverAddress, useTLS)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating connection: %v", err)
			return
		}

		ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
			"access_token": accessToken,
			"player_id":    playerID,
		}))

		b, err := c.GetBoard(ctx, &spb.GetBoardRequest{
			GameId: gameID,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting board: %v", err)
			return
		}
		fmt.Println(b)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	showCmd.Flags().StringVarP(&gameID, "game_id", "g", "", "Game ID for the board")
	showCmd.MarkFlagRequired("game_id")
}
