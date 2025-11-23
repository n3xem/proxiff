package plugin

import (
	"github.com/hashicorp/go-plugin"
	"github.com/n3xem/proxiff/comparator"
	"github.com/n3xem/proxiff/plugin/builtin"
)

// StartBuiltinPlugin starts a builtin comparator as an in-process plugin
func StartBuiltinPlugin() (comparator.Comparator, *plugin.Client, error) {
	// Return the builtin comparator directly without gRPC overhead
	// No plugin client is needed since this runs in-process
	comp := builtin.NewSimpleComparator()
	return comp, nil, nil
}
