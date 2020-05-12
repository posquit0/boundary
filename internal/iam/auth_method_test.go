package iam

import (
	"context"
	"testing"

	"github.com/hashicorp/watchtower/internal/db"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func Test_NewAuthMethod(t *testing.T) {
	t.Parallel()
	cleanup, conn := db.TestSetup(t, "postgres")
	defer cleanup()
	assert := assert.New(t)
	defer conn.Close()

	t.Run("valid", func(t *testing.T) {
		w := db.New(conn)
		s, err := NewOrganization()
		assert.Nil(err)
		assert.True(s.Scope != nil)
		err = w.Create(context.Background(), s)
		assert.Nil(err)
		assert.True(s.PublicId != "")

		meth, err := NewAuthMethod(s, AuthUserPass)
		assert.Nil(err)
		assert.True(meth != nil)
		err = w.Create(context.Background(), meth)
		assert.Nil(err)
		assert.True(meth != nil)
		assert.Equal(meth.Type, AuthUserPass.String())
	})
	t.Run("nil-scope", func(t *testing.T) {
		meth, err := NewAuthMethod(nil, AuthUserPass)
		assert.True(err != nil)
		assert.True(meth == nil)
		assert.Equal(err.Error(), "error scope is nil for new auth method")
	})
}

func TestAuthMethod_GetScope(t *testing.T) {
	t.Parallel()
	cleanup, conn := db.TestSetup(t, "postgres")
	defer cleanup()
	assert := assert.New(t)
	defer conn.Close()

	t.Run("valid", func(t *testing.T) {
		w := db.New(conn)
		s, err := NewOrganization()
		assert.Nil(err)
		assert.True(s.Scope != nil)
		err = w.Create(context.Background(), s)
		assert.Nil(err)
		assert.True(s.PublicId != "")

		meth, err := NewAuthMethod(s, AuthUserPass)
		assert.Nil(err)
		assert.True(meth != nil)
		err = w.Create(context.Background(), meth)
		assert.Nil(err)
		assert.True(meth != nil)
		assert.Equal(meth.Type, AuthUserPass.String())

		scope, err := meth.GetScope(context.Background(), w)
		assert.Nil(err)
		assert.Equal(scope.GetPublicId(), s.PublicId)
	})

}

func TestAuthMethod_ResourceType(t *testing.T) {
	t.Parallel()
	cleanup, conn := db.TestSetup(t, "postgres")
	defer cleanup()
	assert := assert.New(t)
	defer conn.Close()

	t.Run("valid", func(t *testing.T) {
		w := db.New(conn)
		s, err := NewOrganization()
		assert.Nil(err)
		assert.True(s.Scope != nil)
		err = w.Create(context.Background(), s)
		assert.Nil(err)
		assert.True(s.PublicId != "")

		meth, err := NewAuthMethod(s, AuthUserPass)
		assert.Nil(err)
		assert.True(meth != nil)
		err = w.Create(context.Background(), meth)
		assert.Nil(err)
		assert.True(meth != nil)
		assert.Equal(meth.Type, AuthUserPass.String())

		ty := meth.ResourceType()
		assert.Equal(ty, ResourceTypeAuthMethod)
	})
}

func TestAuthMethod_Actions(t *testing.T) {
	assert := assert.New(t)
	meth := &AuthMethod{}
	a := meth.Actions()
	assert.Equal(a[ActionCreate.String()], ActionCreate)
	assert.Equal(a[ActionUpdate.String()], ActionUpdate)
	assert.Equal(a[ActionRead.String()], ActionRead)
	assert.Equal(a[ActionDelete.String()], ActionDelete)
}

func TestAuthMethod_Clone(t *testing.T) {
	t.Parallel()
	cleanup, conn := db.TestSetup(t, "postgres")
	defer cleanup()
	assert := assert.New(t)
	defer conn.Close()

	t.Run("valid", func(t *testing.T) {
		w := db.New(conn)
		s, err := NewOrganization()
		assert.Nil(err)
		assert.True(s.Scope != nil)
		err = w.Create(context.Background(), s)
		assert.Nil(err)
		assert.True(s.PublicId != "")

		meth, err := NewAuthMethod(s, AuthUserPass)
		assert.Nil(err)
		assert.True(meth != nil)
		err = w.Create(context.Background(), meth)
		assert.Nil(err)
		assert.True(meth != nil)
		assert.Equal(meth.Type, AuthUserPass.String())

		cp := meth.Clone()
		assert.True(proto.Equal(cp.(*AuthMethod).AuthMethod, meth.AuthMethod))
	})
	t.Run("not-equal", func(t *testing.T) {
		w := db.New(conn)
		s, err := NewOrganization()
		assert.Nil(err)
		assert.True(s.Scope != nil)
		err = w.Create(context.Background(), s)
		assert.Nil(err)
		assert.True(s.PublicId != "")

		meth, err := NewAuthMethod(s, AuthUserPass)
		assert.Nil(err)
		assert.True(meth != nil)
		err = w.Create(context.Background(), meth)
		assert.Nil(err)
		assert.True(meth != nil)
		assert.Equal(meth.Type, AuthUserPass.String())

		meth2, err := NewAuthMethod(s, AuthUserPass)
		assert.Nil(err)
		assert.True(meth2 != nil)
		err = w.Create(context.Background(), meth2)
		assert.Nil(err)

		cp := meth.Clone()
		assert.True(!proto.Equal(cp.(*AuthMethod).AuthMethod, meth2.AuthMethod))
	})
}