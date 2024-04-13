package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetUser(t *testing.T) {
	s := NewSimpleStorage(make(chan StorageData))

	err := s.SetUser("pathe", []AclRule{})
	require.NoError(t, err)

	t.Run("allow user to login with any password", func(t *testing.T) {
		t.Parallel()
		s := NewSimpleStorage(make(chan StorageData))

		err := s.SetUser("pathe", []AclRule{
			AclRule{Type: AclAnyPassword},
		})
		require.NoError(t, err)

		u, err := s.GetUser("pathe")
		require.NoError(t, err)
		assert.Equal(t, true, u.AnyPassword)
	})

	t.Run("can assign password to user", func(t *testing.T) {
		t.Parallel()
		s := NewSimpleStorage(make(chan StorageData))

		err := s.SetUser("pathe", []AclRule{
			AclRule{Type: AclSetUserPassword, Value: "mypwd"},
		})
		require.NoError(t, err)

		u, err := s.GetUser("pathe")
		require.NoError(t, err)
		assert.Equal(t, []string{"mypwd"}, u.Passwords)
	})

	t.Run("can activate a user", func(t *testing.T) {
		t.Parallel()
		s := NewSimpleStorage(make(chan StorageData))

		err := s.SetUser("pathe", []AclRule{
			AclRule{Type: AclActivateUser},
		})
		require.NoError(t, err)

		u, err := s.GetUser("pathe")
		require.NoError(t, err)
		assert.Equal(t, true, u.Active)

		err = s.SetUser("pathe", []AclRule{
			AclRule{Type: AclDisableUser},
		})
		require.NoError(t, err)

		assert.Equal(t, false, u.Active)
	})
}
