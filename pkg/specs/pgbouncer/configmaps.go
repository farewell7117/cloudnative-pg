package pgbouncer

import (
	"fmt"
	apiv1 "github.com/cloudnative-pg/cloudnative-pg/api/v1"
	"github.com/cloudnative-pg/cloudnative-pg/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConfigMap creates the ConfigMap containing Odyssey configuration
func ConfigMap(pooler *apiv1.Pooler, cluster *apiv1.Cluster) (*corev1.ConfigMap, error) {
	odysseySpec := &apiv1.OdysseySpec{}
	if pooler.Spec.Odyssey != nil {
		odysseySpec = pooler.Spec.Odyssey
	}

	serviceType := pooler.Spec.Type
	clusterHost := fmt.Sprintf("%s-%s.%s.svc.cluster.local", cluster.Name, serviceType, cluster.Namespace)

	logFormat := "%p %t %l [%i] (%c) %m"
	if odysseySpec.LogFormat != nil && *odysseySpec.LogFormat != "" {
		logFormat = *odysseySpec.LogFormat
	}

	workers := intValOrDefault(odysseySpec.Workers, 1)
	resolvers := intValOrDefault(odysseySpec.Resolvers, 1)
	statsInterval := intValOrDefault(odysseySpec.StatsInterval, 3)
	keepalive := intValOrDefault(odysseySpec.Keepalive, 7200)
	listenPort := intValOrDefault(odysseySpec.ListenPort, 6432)

	logDebug := boolYesNo(odysseySpec.LogDebug, true)
	logConfig := boolYesNo(odysseySpec.LogConfig, true)
	logSession := boolYesNo(odysseySpec.LogSession, true)
	logQuery := boolYesNo(odysseySpec.LogQuery, true)
	logStats := boolYesNo(odysseySpec.LogStats, true)

	databaseName := "app"
	if odysseySpec.DatabaseName != "" {
		databaseName = odysseySpec.DatabaseName
	}

	databaseUser := "app"
	if odysseySpec.DatabaseUser != "" {
		decodedUsername, err := decodeBase64(odysseySpec.DatabaseUser)
		if err != nil {
			return nil, fmt.Errorf("failed to decode database user: %w", err)
		}
		databaseUser = decodedUsername
	}

	password := "password"
	if odysseySpec.Password != nil && *odysseySpec.Password != "" {
		decodedPassword, err := decodeBase64(*odysseySpec.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to decode password: %w", err)
		}
		password = decodedPassword
	}

	extraConfig := odysseySpec.Configuration

	// Odyssey configuration
	odysseyConfig := fmt.Sprintf(`daemonize no

pid_file "/var/run/odyssey/odyssey.pid"

log_format "%s"
log_to_stdout yes
log_debug %s
log_config %s
log_session %s
log_query %s
log_stats %s

stats_interval %d

workers %d
resolvers %d

keepalive %d

listen {
  host "0.0.0.0"
  port %d
}

storage "default" {
  type "remote"
  host "%s"
  port 5432
}

database "%s" {
  user "%s" {
    authentication "clear_text"
    password "%s"
    storage "default"
    pool "transaction"
  }
}

%s
`,
		logFormat,
		logDebug, logConfig, logSession, logQuery, logStats,
		statsInterval,
		workers, resolvers,
		keepalive,
		listenPort,
		clusterHost,
		databaseName,
		databaseUser,
		password,
		extraConfig,
	)

	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pooler.Name + "-config",
			Namespace: pooler.Namespace,
			Labels: map[string]string{
				utils.ClusterLabelName:   cluster.Name,
				utils.PgbouncerNameLabel: pooler.Name,
				utils.PodRoleLabelName:   string(utils.PodRolePooler),
				"app":                    "odyssey",
			},
		},
		Data: map[string]string{
			"odyssey.conf": odysseyConfig,
		},
	}, nil
}
