package cmd

import (
	"fmt"
	"go-news-api/database"

	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with initial data",
	Long:  `This command will seed the database with initial data.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := database.Seed(); err != nil {
			fmt.Printf("Error seeding database: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
}
