/**
 * (C) Copyright 2012-2016 HP Development Company, L.P.
 * Confidential computer software. Valid license from HP required for possession, use or copying.
 * Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
 * Computer Software Documentation, and Technical Data for Commercial Items are licensed
 * to the U.S. Government under vendor's standard commercial license.
 */

package elssrv

import (
	"github.com/hpcwp/elsd/pkg/api"
	"golang.org/x/net/context"
)

type HealthGRPCServer struct {
}

func (HealthGRPCServer) Check(context.Context, *api.HealthCheckRequest) (*api.HealthCheckResponse, error) {
	//TODO: check the Els service is actually working
	return &api.HealthCheckResponse{api.HealthCheckResponse_SERVING}, nil
}
