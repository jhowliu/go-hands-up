package localcache

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type localCacheSuite struct {
	suite.Suite
	c Cache
}

func (s *localCacheSuite) SetupTest() {
	s.c = New()
}

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(localCacheSuite))
}

func (s *localCacheSuite) TestGet() {
	tests := []struct {
		Desc      string
		SetupTest func(string)
		Key       string
		Err       error
		Exp       interface{}
	}{
		{
			Desc: "not existed",
			Key:  "not existed",
			Err:  ErrKeyNonExist,
			Exp:  nil,
		},
		{
			Desc: "should get value",
			SetupTest: func(desc string) {
				_ = s.c.Set("value", "value")
			},
			Key: "value",
			Err: nil,
			Exp: "value",
		},
		{
			Desc: "key expired",
			Key:  "expired",
			SetupTest: func(desc string) {
				_ = s.c.Set("expired", "expired")
				s.c.evict("expired")
			},
			Err: ErrKeyNonExist,
			Exp: nil,
		},
	}

	for _, tc := range tests {
		if tc.SetupTest != nil {
			tc.SetupTest(tc.Desc)
		}

		value, err := s.c.Get(tc.Key)
		s.Require().Equal(tc.Err, err, tc.Desc)
		if err == nil {
			s.Require().Equal(tc.Exp, value, tc.Desc)
		}
	}
}

func (s *localCacheSuite) TestSet() {
	tests := []struct {
		Desc      string
		SetupTest func(string)
		Key       string
		Exp       interface{}
	}{
		{
			Desc: "set value if key is not existed",
			Key:  "not existed",
			SetupTest: func(desc string) {
				_ = s.c.Set("not existed", "value")
			},
			Exp: "value",
		},
		{
			Desc: "override if key exits",
			SetupTest: func(desc string) {
				_ = s.c.Set("override", "value")
				_ = s.c.Set("override", "new value")
			},
			Key: "override",
			Exp: "new value",
		},
	}

	for _, tc := range tests {
		if tc.SetupTest != nil {
			tc.SetupTest(tc.Desc)
		}

		value, err := s.c.Get(tc.Key)
		if err == nil {
			s.Require().Equal(tc.Exp, value, tc.Desc)
		}
	}
}
