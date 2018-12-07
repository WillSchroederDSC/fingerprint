package server

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/willschroeder/fingerprint/pkg/models"
	"github.com/willschroeder/fingerprint/pkg/proto"
	"github.com/willschroeder/fingerprint/pkg/session_representations"
)

func BuildSessionRepresentation(user *models.User, sessionUUID string, protoScopeGroupings []*proto.ScopeGrouping) (tokenStr string, json string, err error) {
	tf := session_representations.NewTokenFactory(user.Uuid, sessionUUID)
	for _, sg := range protoScopeGroupings {
		exp, err := ptypes.Timestamp(sg.Expiration)
		if err != nil {
			return "", "", errors.Wrap(err, "couldn't convert timestamp")
		}
		tf.AddScopeGrouping(sg.Scopes, exp)
	}

	sess, err := tf.GenerateSession()
	if err != nil {
		return "", "", err
	}

	return sess.Token, sess.Json, nil
}
