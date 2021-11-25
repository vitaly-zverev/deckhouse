// Copyright 2021 Flant JSC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package hooks

import (
	"time"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	"github.com/flant/shell-operator/pkg/kube_events_manager/types"

	"github.com/deckhouse/deckhouse/go_lib/filter"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:       "prometheus_scrape_interval",
			ApiVersion: "v1",
			Kind:       "ConfigMap",
			NameSelector: &types.NameSelector{
				MatchNames: []string{"prometheus-scrape-interval"},
			},
			NamespaceSelector: &types.NamespaceSelector{
				NameSelector: &types.NameSelector{
					MatchNames: []string{"d8-monitoring"},
				},
			},
			FilterFunc: filter.KeyFromConfigMap("scrapeInterval"),
		},
	},
}, discoveryPromscaleScrapeInterval)

// discoveryPromscaleScrapeInterval
// There is CM d8-monitoring/prometheus-scrape-interval with prometheus scrape interval.
// Hook must store it to `global.discovery.prometheusScrapeInterval`.
func discoveryPromscaleScrapeInterval(input *go_hook.HookInput) error {
	intervalScrapSnap := input.Snapshots["prometheus_scrape_interval"]

	intervalInSeconds := 30
	if len(intervalScrapSnap) > 0 {
		interval, err := time.ParseDuration(intervalScrapSnap[0].(string))
		if err != nil {
			input.LogEntry.Warnf("Prometheus scrape interval from ConfigMap was ignored. Use default: %vs. Cannot parse duration: %v", intervalInSeconds, err)
		} else {
			intervalInSeconds = int(interval.Seconds())
		}
	}

	input.Values.Set("global.discovery.prometheusScrapeInterval", intervalInSeconds)

	return nil
}
