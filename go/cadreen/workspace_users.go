package cadreen

import (
	"context"
	"fmt"
)

// ListWorkspaceUsers returns all users in the current workspace.
func (c *Client) ListWorkspaceUsers(ctx context.Context, opts ...RequestOption) (*ListWorkspaceUsersResponse, error) {
	var result ListWorkspaceUsersResponse
	if err := c.do(ctx, "GET", "/api/v1/cadreen/workspace/users", nil, &result, opts...); err != nil {
		return nil, fmt.Errorf("list workspace users: %w", err)
	}
	return &result, nil
}

// InviteUser invites a user to the workspace by email.
func (c *Client) InviteUser(ctx context.Context, req InviteUserRequest, opts ...RequestOption) (*WorkspaceUser, error) {
	var result WorkspaceUser
	if err := c.do(ctx, "POST", "/api/v1/cadreen/workspace/users", req, &result, opts...); err != nil {
		return nil, fmt.Errorf("invite user: %w", err)
	}
	return &result, nil
}

// UpdateUserRole updates a workspace user's role.
func (c *Client) UpdateUserRole(ctx context.Context, id string, req UpdateRoleRequest, opts ...RequestOption) (*WorkspaceUser, error) {
	var result WorkspaceUser
	if err := c.do(ctx, "PATCH", "/api/v1/cadreen/workspace/users/"+id, req, &result, opts...); err != nil {
		return nil, fmt.Errorf("update user role: %w", err)
	}
	return &result, nil
}

// RemoveUser removes a user from the workspace.
func (c *Client) RemoveUser(ctx context.Context, id string, opts ...RequestOption) error {
	if err := c.do(ctx, "DELETE", "/api/v1/cadreen/workspace/users/"+id, nil, nil, opts...); err != nil {
		return fmt.Errorf("remove user: %w", err)
	}
	return nil
}
