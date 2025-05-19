package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewSubscription(t *testing.T) {
	email := "test@example.com"
	city := "Kyiv"
	freq := FrequencyDaily

	s, s_err := NewSubscription(email, city, freq)

	require.NoError(t, s_err)

	require.NotNil(t, s)
	require.Equal(t, email, s.Email)
	require.Equal(t, city, s.City)
	require.Equal(t, freq, s.Frequency)

	// Check ID is valid UUID
	_, err := uuid.Parse(s.ID.String())
	require.NoError(t, err)

	// Check tokens are valid UUIDs and not equal
	confUUID, err := uuid.Parse(s.ConfirmationToken)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, confUUID)

	unsubUUID, err := uuid.Parse(s.UnsubscribeToken)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, unsubUUID)
	require.NotEqual(t, s.ConfirmationToken, s.UnsubscribeToken)

	// Confirmed should be false by default
	require.False(t, s.Confirmed)

	// ConfirmedAt and LastSentAt should be nil
	require.Nil(t, s.ConfirmedAt)
	require.Nil(t, s.LastSentAt)

	// CreatedAt should be recent (within last second)
	now := time.Now().UTC()
	require.WithinDuration(t, now, s.CreatedAt, time.Second)
}
