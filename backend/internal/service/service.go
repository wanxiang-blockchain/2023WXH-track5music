package service

import (
	"backend/pkg/aws"
	"backend/pkg/helper/sid"
	"backend/pkg/jwt"
	"backend/pkg/log"
)

type Service struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	s3     *aws.BucketBasics
}

func NewService(logger *log.Logger, sid *sid.Sid, jwt *jwt.JWT, s3 *aws.BucketBasics) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
		jwt:    jwt,
		s3:     s3,
	}
}
