package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/orkarstoft/kscale/pkg/config"
	"github.com/orkarstoft/kscale/pkg/k8s"
	"github.com/orkarstoft/kscale/pkg/logger"
)

func Listen() error {
	ctx := context.Background()

	// Create the client.
	client, err := client(config.Config.ProjectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient error: %v", err)
	}

	// Create subscription name
	subscriptionName := fmt.Sprintf("kscale-%s", config.Config.ClusterName)
	var subscription *pubsub.Subscription
	// Check if subscription exists
	subscription = client.Subscription(subscriptionName)
	exists, err := subscription.Exists(ctx)
	if err != nil {
		return fmt.Errorf("pubsub.Subscription.Exists error: failed to check if subscription exists %v", err)
	}

	// Create subscription if it doesn't exist
	if !exists {
		subscription, err = createSubscription(ctx, client, subscriptionName)
		if err != nil {
			return fmt.Errorf("pubsub.CreateSubscription error: %v", err)
		}
	} else {
		subConf, err := subscription.Config(ctx)
		if err != nil {
			return fmt.Errorf("pubsub.Subscription.Config error: %v", err)
		}

		logger.Log.Info().Msgf("Subscription %s already exists with attribute filter '%s'", subscriptionName, strings.ReplaceAll(subConf.Filter, "\"", ""))
	}

	// Receive messages
	err = subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		logger.Log.Debug().Msgf("Received message: %s", string(msg.Data))

		// Unmarshal message
		var m PubSubMsg
		if err := json.Unmarshal(msg.Data, &m); err != nil {
			logger.Log.Panic().Err(err).Msg("Failed to unmarshal message")
		}

		if m.Action == "kscale_scale_namespace_up" {
			logger.Log.Info().Msgf("Scaling %s namespace %s up", m.Cluster, m.Namespace)
			convertIntToTimeDuration, err := time.ParseDuration(fmt.Sprintf("%dh", m.Duration))
			if err != nil {
				logger.Log.Panic().Err(err).Msg("Failed to parse duration")
			}

			logger.Log.Debug().Msgf("Duration: %d, Duration in time.Duration: %s", m.Duration, convertIntToTimeDuration)

			k8s.ScaleNamespaceUp(m.Namespace, convertIntToTimeDuration)
		}

		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("pubsub.Receive error: %v", err)
	}

	return nil
}

type PubSubMsg struct {
	Action    string `json:"action"`
	Namespace string `json:"namespace"`
	Cluster   string `json:"cluster"`
	Duration  int    `json:"duration"`
}
