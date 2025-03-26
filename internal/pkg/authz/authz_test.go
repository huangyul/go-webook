package authz

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthz_VerifyToken(t *testing.T) {
	authz := &authz{}
	accessToken, _, err := authz.GenerateToken(1, "1")
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	_, err = authz.VerifyToken(accessToken)
	assert.NoError(t, err)
}
