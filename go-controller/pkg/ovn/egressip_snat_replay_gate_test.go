// SPDX-FileCopyrightText: Copyright The OVN-Kubernetes Contributors
// SPDX-License-Identifier: Apache-2.0
package ovn

import (
	"sync"

	libovsdbclient "github.com/ovn-kubernetes/libovsdb/client"
	"github.com/ovn-kubernetes/libovsdb/model"

	"github.com/ovn-kubernetes/ovn-kubernetes/go-controller/pkg/nbdb"
)

// snatReplayGateClient pauses a gateway NAT lookup in the pod setup client.
// The EgressIP controller keeps its original client, allowing the test to
// verify that assignment waits on the pod lock rather than this client gate.
type snatReplayGateClient struct {
	libovsdbclient.Client
	routerName string
	reached    chan struct{}
	resume     chan struct{}
	once       sync.Once
}

func (c *snatReplayGateClient) Where(models ...model.Model) libovsdbclient.ConditionalAPI {
	for _, m := range models {
		if router, ok := m.(*nbdb.LogicalRouter); ok && router.Name == c.routerName {
			c.once.Do(func() {
				close(c.reached)
				<-c.resume
			})
		}
	}
	return c.Client.Where(models...)
}
