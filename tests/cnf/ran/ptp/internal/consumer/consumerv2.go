// Package consumer provides helpers for deploying the cloud-event-consumer service on a cluster with potentially
// multiple nodes. We want to have a deployment with a single replica per node corresponding to each linuxptp-daemon
// deployment. However, a daemonset is not suitable because each a separate service is needed for each node.
package consumer

import (
	"errors"
	"fmt"

	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/clients"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/deployment"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/nodes"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/pod"
	"github.com/rh-ecosystem-edge/eco-goinfra/pkg/service"
	. "github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/internal/raninittools"
	"github.com/rh-ecosystem-edge/eco-gotests/tests/cnf/ran/ptp/internal/tsparams"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// DeployV2ConsumersOnWorkers deploys the cloud-event-consumer deployment and service on all worker nodes in the
// cluster, assuming the event API version is v2. It accumulates errors during deployment and returns them all at once,
// so one deployment failing does not mean deployments fail for all nodes.
func DeployV2ConsumersOnWorkers(client *clients.Settings) error {
	err := createConsumerNamespace(client)
	if err != nil {
		return fmt.Errorf("failed to deploy consumer namespace: %w", err)
	}

	workers, err := nodes.List(client, metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set(RANConfig.WorkerLabelMap)).String(),
	})
	if err != nil {
		return fmt.Errorf("failed to list worker nodes: %w", err)
	}

	var deployErrors []error

	for _, worker := range workers {
		err = createV2ConsumerServiceOnNode(client, worker.Definition.Name)
		if err != nil {
			deployErrors = append(deployErrors,
				fmt.Errorf("failed to create consumer service on node %s: %w", worker.Definition.Name, err))

			continue
		}

		err = createV2ConsumerDeploymentOnNode(client, worker.Definition.Name)
		if err != nil {
			deployErrors = append(deployErrors,
				fmt.Errorf("failed to create consumer deployment on node %s: %w", worker.Definition.Name, err))

			continue
		}
	}

	return errors.Join(deployErrors...)
}

// CleanupV2ConsumersOnWorkers deletes the cloud-event-consumer deployment and service on all worker nodes in the
// cluster, assuming the event API version is v2. It accumulates errors during deletion and returns them all at once, so
// one deletion failing does not mean deletions fail for all nodes.
func CleanupV2ConsumersOnWorkers(client *clients.Settings) error {
	workers, err := nodes.List(client, metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labels.Set(RANConfig.WorkerLabelMap)).String(),
	})
	if err != nil {
		return fmt.Errorf("failed to list worker nodes: %w", err)
	}

	var cleanupErrors []error

	for _, worker := range workers {
		err = deleteV2ConsumerDeploymentOnNode(client, worker.Definition.Name)
		if err != nil {
			cleanupErrors = append(cleanupErrors,
				fmt.Errorf("failed to delete consumer deployment on node %s: %w", worker.Definition.Name, err))

			continue
		}

		err = deleteV2ConsumerServiceOnNode(client, worker.Definition.Name)
		if err != nil {
			cleanupErrors = append(cleanupErrors,
				fmt.Errorf("failed to delete consumer service on node %s: %w", worker.Definition.Name, err))

			continue
		}
	}

	err = deleteConsumerNamespace(client)
	if err != nil {
		cleanupErrors = append(cleanupErrors,
			fmt.Errorf("failed to delete consumer namespace: %w", err))
	}

	return errors.Join(cleanupErrors...)
}

// createV2ConsumerDeploymentOnNode creates a new deployment for the cloud-event-consumer deployment with a specific
// node selected. It uses the definition from
// https://github.com/redhat-cne/cloud-event-proxy/blob/main/examples/manifests/consumer.yaml.
func createV2ConsumerDeploymentOnNode(client *clients.Settings, nodeName string) error {
	v2ConsumerImage := RANConfig.PtpEventConsumerImage + RANConfig.PtpEventConsumerV2Tag
	consumerContainer, err := pod.NewContainerBuilder(
		"cloud-event-consumer", v2ConsumerImage, []string{"./cloud-event-consumer"}).
		WithPorts([]corev1.ContainerPort{{
			Name:          consumerPortName,
			ContainerPort: consumerPort,
		}}).
		WithImagePullPolicy(corev1.PullAlways).
		WithEnvVar("CONSUMER_TYPE", "PTP").
		WithEnvVar("ENABLE_STATUS_CHECK", "true").
		WithEnvVar("NODE_NAME", nodeName).
		GetContainerCfg()

	if err != nil {
		return fmt.Errorf("failed to create consumer container: %w", err)
	}

	apiAddr := fmt.Sprintf("--local-api-addr=%s.%s.svc.cluster.local:9043",
		getConsumerServiceName(nodeName), tsparams.CloudEventsNamespace)
	consumerContainer.Args = []string{
		apiAddr,
		"--api-path=/api/ocloudNotifications/v2/",
		"--http-event-publishers=ptp-event-publisher-service-NODE_NAME.openshift-ptp.svc.cluster.local:9043",
	}
	consumerContainer.SecurityContext = nil

	consumerDeployment := deployment.NewBuilder(
		client,
		getConsumerDeploymentName(nodeName),
		tsparams.CloudEventsNamespace,
		getConsumerSelectorLabels(nodeName),
		*consumerContainer).
		WithReplicas(1).
		WithNodeSelector(RANConfig.WorkerLabelMap).
		WithAffinity(&corev1.Affinity{
			NodeAffinity: &corev1.NodeAffinity{
				RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
					NodeSelectorTerms: []corev1.NodeSelectorTerm{{
						MatchFields: []corev1.NodeSelectorRequirement{{
							Key:      "metadata.name",
							Operator: corev1.NodeSelectorOpIn,
							Values:   []string{nodeName},
						}},
					}},
				},
			},
		})
	consumerDeployment.Definition.Spec.Template.ObjectMeta.Annotations = workloadManagementAnnotation

	_, err = consumerDeployment.CreateAndWaitUntilReady(createDeleteTimeout)
	if err != nil {
		return fmt.Errorf("failed to create consumer deployment: %w", err)
	}

	return nil
}

// deleteV2ConsumerDeploymentOnNode deletes the cloud-event-consumer deployment with a specific node selected. It is
// the inverse of createV2ConsumerDeploymentOnNode.
func deleteV2ConsumerDeploymentOnNode(client *clients.Settings, nodeName string) error {
	consumerDeployment, err := deployment.Pull(client, getConsumerDeploymentName(nodeName), tsparams.CloudEventsNamespace)
	if err != nil {
		return nil
	}

	err = consumerDeployment.DeleteAndWait(createDeleteTimeout)
	if err != nil {
		return fmt.Errorf("failed to delete consumer deployment: %w", err)
	}

	return nil
}

// createV2ConsumerServiceOnNode creates a new service for the cloud-event-consumer deployment with a specific node
// selected. It uses the definition from
// https://github.com/redhat-cne/cloud-event-proxy/blob/main/examples/manifests/service.yaml.
func createV2ConsumerServiceOnNode(client *clients.Settings, nodeName string) error {
	_, err := service.NewBuilder(
		client,
		getConsumerServiceName(nodeName),
		tsparams.CloudEventsNamespace,
		getConsumerSelectorLabels(nodeName),
		corev1.ServicePort{Name: consumerPortName, Port: consumerPort}).
		WithAnnotation(map[string]string{"prometheus.io/scrape": "true"}).
		Create()
	if err != nil {
		return fmt.Errorf("failed to create consumer service: %w", err)
	}

	return nil
}

// deleteV2ConsumerServiceOnNode deletes the cloud-event-consumer service with a specific node selected. It is the
// inverse of createV2ConsumerServiceOnNode.
func deleteV2ConsumerServiceOnNode(client *clients.Settings, nodeName string) error {
	consumerService, err := service.Pull(client, getConsumerServiceName(nodeName), tsparams.CloudEventsNamespace)
	if err != nil {
		return nil
	}

	err = consumerService.Delete()
	if err != nil {
		return fmt.Errorf("failed to delete consumer service: %w", err)
	}

	return nil
}
