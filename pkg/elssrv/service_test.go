/**
 * (C) Copyright 2012-2016 HP Development Company, L.P.
 * Confidential computer software. Valid license from HP required for possession, use or copying.
 * Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
 * Computer Software Documentation, and Technical Data for Commercial Items are licensed
 * to the U.S. Government under vendor's standard commercial license.
 */

package elssrv

import (
	"testing"
)


func TestElsService_AddKey(t *testing.T) {

	type args struct {
		key string
		srv ServiceInstance
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name:"simple",  args:	args{ key: "1112", srv: ServiceInstance { "127.0.0.1",  "rw"}}, wantErr: false},
		{name:"simple",  args:	args{ key: "2233", srv: ServiceInstance { "127.0.0.1",  "rw"}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.AddKey(tt.args.key, tt.args.srv); (err != nil) != tt.wantErr {
				t.Errorf("ElsService.AddKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
