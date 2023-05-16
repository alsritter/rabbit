package helloworld_service

import (
	pb "alsritter.icu/rabbit-template/api/helloworld/v1"
	biz "alsritter.icu/rabbit-template/internal/biz/helloworld"
)

type HelloworldService struct {
	pb.UnimplementedHelloServiceServer
	hu *biz.HelloworldUsecase
}

func NewHelloworldService(hu *biz.HelloworldUsecase) *HelloworldService {
	return &HelloworldService{hu: hu}
}
