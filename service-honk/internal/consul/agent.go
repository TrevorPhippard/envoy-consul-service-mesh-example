package consul

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

type Config struct {
	ServiceID   string
	ServiceName string
	Address     string
	Port        int
	Tags        []string
	TTL         time.Duration
	CheckID     string
}

type Agent struct {
	client        *api.Client
	seenInstances map[string]bool
	ttl           time.Duration
	checkID       string
}

func NewAgent(consulConfig *api.Config) *Agent {
	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatal("consul client error:", err)
	}

	return &Agent{
		client:        client,
		seenInstances: make(map[string]bool),
	}
}

/* -------------------------------------------------------------------------- */
/*                           TTL HEALTH CHECK LOOP                            */
/* -------------------------------------------------------------------------- */

func (a *Agent) updateHealthCheck() {
	ticker := time.NewTicker(a.ttl / 2) // update at half TTL

	for range ticker.C {
		err := a.client.Agent().UpdateTTL(a.checkID, "online", api.HealthPassing)
		if err != nil {
			log.Println("ttl update error:", err)
		}
	}
}

/* -------------------------------------------------------------------------- */
/*                             SERVICE REGISTRATION                           */
/* -------------------------------------------------------------------------- */

func (a *Agent) RegisterService(cfg Config) {
	a.ttl = cfg.TTL
	a.checkID = cfg.CheckID

	// TTL health check
	check := &api.AgentServiceCheck{
		DeregisterCriticalServiceAfter: cfg.TTL.String(),
		TLSSkipVerify:                  true,
		TTL:                            cfg.TTL.String(),
		CheckID:                        cfg.CheckID,
	}

	// Service registration
	registration := &api.AgentServiceRegistration{
		ID:      cfg.ServiceID,
		Name:    cfg.ServiceName,
		Tags:    cfg.Tags,
		Address: cfg.Address,
		Port:    cfg.Port,
		Checks:  []*api.AgentServiceCheck{check},
	}

	/* ---------------------------------------------------------------------- */
	/*                         WATCH FOR NEW SERVICE INSTANCES                */
	/* ---------------------------------------------------------------------- */

	query := map[string]any{
		"type":        "service",
		"service":     cfg.ServiceName,
		"passingonly": true,
	}

	plan, err := watch.Parse(query)
	if err != nil {
		log.Fatal("consul watch parse error:", err)
	}

	plan.HybridHandler = func(_ watch.BlockingParamVal, result any) {
		entries, ok := result.([]*api.ServiceEntry)
		if !ok {
			return
		}

		for _, entry := range entries {
			if entry == nil || entry.Service == nil {
				continue
			}

			id := entry.Service.ID
			if !a.seenInstances[id] {
				a.seenInstances[id] = true
				fmt.Printf(
					"🟢 New instance joined: %s (Address=%s Port=%d)\n",
					entry.Service.ID,
					entry.Service.Address,
					entry.Service.Port,
				)
			}
		}
	}

	// Run the watch in the background
	go func() {
		plan.RunWithConfig("", &api.Config{})
	}()

	// Register the service with Consul
	if err := a.client.Agent().ServiceRegister(registration); err != nil {
		log.Fatal("service registration error:", err)
	}

	// Start TTL updater
	go a.updateHealthCheck()
}
