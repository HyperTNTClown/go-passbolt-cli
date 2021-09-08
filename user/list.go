package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/speatzle/go-passbolt-cli/util"
	"github.com/speatzle/go-passbolt/api"
	"github.com/spf13/cobra"

	"github.com/pterm/pterm"
)

// UserListCmd Lists a Passbolt User
var UserListCmd = &cobra.Command{
	Use:     "user",
	Short:   "Lists Passbolt Users",
	Long:    `Lists Passbolt Users`,
	Aliases: []string{"groups"},
	RunE:    UserList,
}

func init() {
	UserListCmd.Flags().StringArrayP("groups", "g", []string{}, "Users that are members of groups")
	UserListCmd.Flags().StringArrayP("resources", "r", []string{}, "Users that have access to resources")

	UserListCmd.Flags().StringP("search", "s", "", "Search for Users")
	UserListCmd.Flags().BoolP("admin", "a", false, "Only show Admins")

	UserListCmd.Flags().StringArrayP("columns", "c", []string{"ID", "Username", "FirstName", "LastName", "Role"}, "Columns to return, possible Columns:\nID, Username, FirstName, LastName, Role")
}

func UserList(cmd *cobra.Command, args []string) error {
	groups, err := cmd.Flags().GetStringArray("groups")
	if err != nil {
		return err
	}
	resources, err := cmd.Flags().GetStringArray("resources")
	if err != nil {
		return err
	}
	search, err := cmd.Flags().GetString("search")
	if err != nil {
		return err
	}
	admin, err := cmd.Flags().GetBool("admin")
	if err != nil {
		return err
	}
	columns, err := cmd.Flags().GetStringArray("columns")
	if err != nil {
		return err
	}
	if len(columns) == 0 {
		return fmt.Errorf("You need to specify atleast one column to return")
	}

	ctx := util.GetContext()

	client, err := util.GetClient(ctx)
	if err != nil {
		return err
	}
	defer client.Logout(context.TODO())
	cmd.SilenceUsage = true

	users, err := client.GetUsers(ctx, &api.GetUsersOptions{
		FilterHasGroup:  groups,
		FilterHasAccess: resources,
		FilterSearch:    search,
		FilterIsAdmin:   admin,
	})
	if err != nil {
		return fmt.Errorf("Listing User: %w", err)
	}

	data := pterm.TableData{columns}

	for _, user := range users {
		entry := make([]string, len(columns))
		for i := range columns {
			switch strings.ToLower(columns[i]) {
			case "id":
				entry[i] = user.ID
			case "username":
				entry[i] = user.Username
			case "firstname":
				entry[i] = user.Profile.FirstName
			case "lastname":
				entry[i] = user.Profile.LastName
			case "role":
				entry[i] = user.Role.Name
			default:
				cmd.SilenceUsage = false
				return fmt.Errorf("Unknown Column: %v", columns[i])
			}
		}
		data = append(data, entry)
	}

	pterm.DefaultTable.WithHasHeader().WithData(data).Render()
	return nil
}
