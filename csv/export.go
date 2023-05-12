package csv

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/passbolt/go-passbolt-cli/util"
	"github.com/passbolt/go-passbolt/api"
	"github.com/passbolt/go-passbolt/helper"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"os"
)

var CSVExportCmd = &cobra.Command{
	Use:     "csv",
	Short:   "Exports Passbolt to a CSV File",
	Long:    `Exports Passbolt to a CSV File`,
	Aliases: []string{},
	RunE:    CSVExport,
}

func init() {
	CSVExportCmd.Flags().StringP("file", "f", "passbolt-export.csv", "File name of the CSV File")
}

func CSVExport(cmd *cobra.Command, args []string) error {
	filename, err := cmd.Flags().GetString("file")
	if err != nil {
		return err
	}

	if filename == "" {
		return fmt.Errorf("the Filename cannot be empty")
	}

	ctx := util.GetContext()

	client, err := util.GetClient(ctx)
	if err != nil {
		return err
	}
	defer client.Logout(context.TODO())
	cmd.SilenceUsage = true

	fmt.Println("Getting Resources...")
	ressources, err := client.GetResources(ctx, &api.GetResourcesOptions{
		ContainSecret:       true,
		ContainResourceType: true,
		ContainTags:         true,
	})
	if err != nil {
		return fmt.Errorf("Getting Resources: %w", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Creating File: %w", err)
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	defer csvWriter.Flush()
	err = csvWriter.Write([]string{"Name", "Username", "URI", "Password", "Description"})
	if err != nil {
		return err
	}

	pterm.EnableStyling()
	pterm.DisableColor()
	progressbar, err := pterm.DefaultProgressbar.WithTotal(len(ressources)).WithTitle("Exporting Resources").Start()
	if err != nil {
		return fmt.Errorf("Creating Progressbar: %w", err)
	}

	for _, resource := range ressources {
		_, _, _, _, pass, desc, err := helper.GetResourceFromData(client, resource, resource.Secrets[0], resource.ResourceType)

		err = csvWriter.Write([]string{resource.Name, resource.Username, resource.URI, pass, desc})
		if err != nil {
			return fmt.Errorf("Writing to File: %w", err)
		}
		progressbar.Increment()
	}

	fmt.Println("Done")
	return nil
}
