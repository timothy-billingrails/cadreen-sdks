package commands

import (
	"context"
	"fmt"

	"github.com/timothy-billingrails/cadreen-sdks/go/cmd/cadreen/output"

	cadreen "github.com/timothy-billingrails/cadreen-sdks/go/cadreen"

	"github.com/spf13/cobra"
)

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Workspace management",
	Long: `Manage workspace users and roles.

Examples:
  cadreen workspace users list
  cadreen workspace users invite user@example.com --role operator
  cadreen workspace users role <user-id> admin
  cadreen workspace users remove <user-id>`,
}

var workspaceUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage workspace users",
	Long: `List, invite, update roles, and remove users from your workspace.

Examples:
  cadreen workspace users list
  cadreen workspace users invite user@example.com --role operator
  cadreen workspace users role wu_01abc admin
  cadreen workspace users remove wu_01abc`,
}

var workspaceUsersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workspace users",
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		client := newClient()
		resp, err := client.ListWorkspaceUsers(context.Background())
		if err != nil {
			return handleAPIError(err)
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(resp, format)
			return nil
		}

		if len(resp.Users) == 0 {
			fmt.Println("No users in this workspace.")
			return nil
		}

		fmt.Printf("Workspace Users (%d):\n\n", resp.Count)
		for _, u := range resp.Users {
			fmt.Printf("  %-12s  %-10s  %s\n", u.ID, u.Role, u.UserID)
		}
		return nil
	},
}

var workspaceUsersInviteCmd = &cobra.Command{
	Use:   "invite [email]",
	Short: "Invite a user to the workspace",
	Long: `Invite a user by email. Default role is member.

Examples:
  cadreen workspace users invite user@example.com
  cadreen workspace users invite admin@example.com --role admin`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		role, _ := cmd.Flags().GetString("role")
		if role == "" {
			role = "member"
		}

		client := newClient()
		wu, err := client.InviteUser(context.Background(), cadreen.InviteUserRequest{
			Email: args[0],
			Role:  cadreen.WorkspaceRole(role),
		})
		if err != nil {
			return handleAPIError(err)
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(wu, format)
			return nil
		}

		fmt.Printf("User invited: %s\n", wu.UserID)
		fmt.Printf("  Role: %s\n", wu.Role)
		fmt.Printf("  ID:   %s\n", wu.ID)
		return nil
	},
}

var workspaceUsersRoleCmd = &cobra.Command{
	Use:   "role [user-id] [role]",
	Short: "Update a user's role",
	Long: `Change a workspace user's role. Valid roles: admin, operator, member, viewer.

Examples:
  cadreen workspace users role wu_01abc admin
  cadreen workspace users role wu_01abc viewer`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		client := newClient()
		wu, err := client.UpdateUserRole(context.Background(), args[0], cadreen.UpdateRoleRequest{
			Role: cadreen.WorkspaceRole(args[1]),
		})
		if err != nil {
			return handleAPIError(err)
		}

		format := outputFormat()
		if format == output.FormatJSON {
			output.Print(wu, format)
			return nil
		}

		fmt.Printf("Role updated: %s → %s\n", wu.UserID, wu.Role)
		return nil
	},
}

var workspaceUsersRemoveCmd = &cobra.Command{
	Use:   "remove [user-id]",
	Short: "Remove a user from the workspace",
	Long: `Remove a user from the workspace by their workspace user ID.

Examples:
  cadreen workspace users remove wu_01abc`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		requireAuth()

		client := newClient()
		if err := client.RemoveUser(context.Background(), args[0]); err != nil {
			return handleAPIError(err)
		}

		fmt.Println("User removed from workspace.")
		return nil
	},
}

func init() {
	workspaceUsersInviteCmd.Flags().String("role", "member", "role: admin, operator, member, viewer")

	workspaceUsersCmd.AddCommand(workspaceUsersListCmd)
	workspaceUsersCmd.AddCommand(workspaceUsersInviteCmd)
	workspaceUsersCmd.AddCommand(workspaceUsersRoleCmd)
	workspaceUsersCmd.AddCommand(workspaceUsersRemoveCmd)

	workspaceCmd.AddCommand(workspaceUsersCmd)
	rootCmd.AddCommand(workspaceCmd)
}
