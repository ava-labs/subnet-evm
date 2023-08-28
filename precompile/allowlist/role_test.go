// (c) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package allowlist

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsNoRole(t *testing.T) {
	tests := []struct {
		role     Role
		expected bool
	}{
		{
			role:     ManagerRole,
			expected: false,
		},
		{
			role:     AdminRole,
			expected: false,
		},
		{
			role:     EnabledRole,
			expected: false,
		},
		{
			role:     NoRole,
			expected: true,
		},
	}

	for index, test := range tests {
		isNoRole := test.role.IsNoRole()
		require.Equal(t, test.expected, isNoRole, fmt.Sprintf("test index: %d", index))
	}
}

func TestIsEnabled(t *testing.T) {
	tests := []struct {
		role     Role
		expected bool
	}{
		{
			role:     ManagerRole,
			expected: true,
		},
		{
			role:     AdminRole,
			expected: true,
		},
		{
			role:     EnabledRole,
			expected: true,
		},
		{
			role:     NoRole,
			expected: false,
		},
	}

	for index, test := range tests {
		isEnabled := test.role.IsEnabled()
		require.Equal(t, test.expected, isEnabled, fmt.Sprintf("test index: %d", index))
	}
}

func TestCanModify(t *testing.T) {
	tests := []struct {
		role     Role
		expected bool
		from     Role
		target   Role
	}{
		{
			role:     ManagerRole,
			expected: true,
			from:     EnabledRole,
			target:   NoRole,
		},
		{
			role:     ManagerRole,
			expected: true,
			from:     NoRole,
			target:   NoRole,
		},
		{
			role:     ManagerRole,
			expected: true,
			from:     EnabledRole,
			target:   EnabledRole,
		},
		{
			role:     ManagerRole,
			expected: false,
			from:     ManagerRole,
			target:   ManagerRole,
		},
		{
			role:     ManagerRole,
			expected: true,
			from:     NoRole,
			target:   EnabledRole,
		},
		{
			role:     ManagerRole,
			expected: false,
			from:     ManagerRole,
			target:   EnabledRole,
		},
		{
			role:     ManagerRole,
			expected: false,
			from:     AdminRole,
			target:   EnabledRole,
		},
		{
			role:     AdminRole,
			expected: true,
			from:     EnabledRole,
			target:   NoRole,
		},
		{
			role:     AdminRole,
			expected: true,
			from:     AdminRole,
			target:   NoRole,
		},
		{
			role:     EnabledRole,
			expected: false,
			from:     EnabledRole,
			target:   NoRole,
		},
		{
			role:     NoRole,
			expected: false,
			from:     EnabledRole,
			target:   NoRole,
		},
	}
	for index, test := range tests {
		canModify := test.role.CanModify(test.from, test.target)
		require.Equal(t, test.expected, canModify, fmt.Sprintf("test index: %d", index))
	}
}
