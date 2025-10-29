# Test Case Summary for vCore Operators Test Suite

Test case `02_validate_operators.go` is located in `tests/system-tests/vcore/tests/02_validate_operators.go` and is part of the "vCore Operators Test Suite". It calls various verification suites for different operators.

## Goal

The goal of this test case is to validate the proper installation and configuration of several operators within the vCore environment, including NMState, Service Mesh, Helm, Redis, Node Tuning Operator (NTO), SR-IOV, KEDA, and NUMA Resources Operator (NROP).

## Test Setup

No specific manual changes are explicitly mentioned in the test file itself. The test suite is composed of multiple verification functions, each focusing on a different operator:

- `VerifyNMStateSuite()`: This suite verifies the existence of the NMState operator namespace, the successful deployment of the NMState operator, and the existence of an NMState instance.
- `VerifyServiceMeshSuite()`: This suite verifies the existence of required namespaces (DTPO, Kiali, Istio, SMO), and the successful deployment and configuration of Distributed Tracing Platform Operator, Kiali, and Service Mesh.
- `VerifyHelmSuite()`: This suite verifies the Helm deployment procedure by downloading and installing Helm, and then checking its functionality.
- `VerifyRedisSuite()`: This suite verifies the Redis deployment procedure, including checking for existing Redis deployments, mirroring images if in a disconnected environment, installing Redis via Helm, and verifying the statefulset and pod counts.
- `VerifyNTOSuite()`: This suite verifies the Node Tuning Operator by checking for its namespace, successful deployment, creation of performance profiles and node tunings, and configuration of CPU Manager and Huge Pages. It also verifies system reserved memory for user-plane-worker nodes.
- `VerifySRIOVSuite()`: This suite verifies the SR-IOV Operator by checking its namespace, successful deployment, and configuration procedure, including disabling injector/webhook and setting node selectors.
- `VerifyKedaSuite()`: This suite verifies the KEDA operator by checking its namespace, successful deployment, creation of a KedaController instance, and deployment of a ScaleObject instance along with a test application exposing Prometheus metrics.
- `VerifyNROPSuite()`: This suite verifies the NUMA Resources Operator by checking its namespace, successful deployment, custom resource deployment, NUMA-aware secondary pod scheduler configuration, and scheduling workloads with the NUMA-aware scheduler.

It does not require a git config set up.

## Test Steps

**VerifyNMStateSuite() steps:**
1. Verify that `VCoreConfig.NMStateOperatorNamespace` exists.
2. Verify NMState operator deployment succeeded (CSV condition).
3. Verify NMState instance (`vcoreparams.NMStateInstanceName`) exists.

**VerifyServiceMeshSuite() steps:**
1. Verify that the following namespaces exist: `vcoreparams.DTPONamespace`, `vcoreparams.KialiNamespace`, `vcoreparams.IstioNamespace`, `vcoreparams.SMONamespace`.
2. Verify Distributed Tracing Platform Operator deployment succeeded.
3. Verify Kiali deployment succeeded.
4. Verify Service Mesh deployment succeeded.
5. Verify Service Mesh configuration procedure succeeded, including creating namespaces for members, creating a member roll, and creating a service-mesh control plane.

**VerifyHelmSuite() steps:**
1. Download the `get-helm-3` script to the hypervisor.
2. Make the script executable.
3. Install Helm by executing the script.
4. Verify Helm is working properly by checking `helm version` output.

**VerifyRedisSuite() steps:**
1. Check if Redis (`redis-ha-server` statefulset in `redis-ha` namespace) is already installed and ready.
2. If not installed:
    a. Retrieve cluster pull-secret.
    b. Check if deployment is disconnected and mirror Redis images if necessary.
    c. Ensure local directory for configuration files exists.
    d. Create a Redis secret.
    e. Create a custom values file (`redis-custom-values.yaml`) for Helm.
    f. Transfer the custom values file to the hypervisor.
    g. Install Redis using Helm with the custom configuration.
3. Wait for the Redis statefulset (`redis-ha-server`) to be ready.
4. Verify Redis server pods count is 3.

**VerifyNTOSuite() steps:**
1. Verify that `vcoreparams.NTONamespace` exists.
2. Verify Node Tuning Operator deployment succeeded by checking deployment and service.
3. Create a new `performanceprofile`.
4. Create new nodes tuning.
5. Verify CPU Manager config.
6. Verify Node Tuning Operator Huge Pages configuration.
7. Verify System Reserved memory for user-plane-worker nodes configuration.

**VerifySRIOVSuite() steps:**
1. Verify that `vcoreparams.SRIOVNamespace` exists.
2. Verify SR-IOV Operator deployment succeeded.
3. Verify SR-IOV configuration procedure succeeded, including changing `sriovoperatorconfig` to disable Injector and OperatorWebhook and setting a node selector.

**VerifyKedaSuite() steps:**
1. Verify that `vcoreparams.KedaNamespace` exists.
2. Verify Keda operator deployment succeeded.
3. Verify KedaController instance (`vcoreparams.KedaControllerName`) created successfully, with specified admission webhooks, operator, metrics server, and watch namespace.
4. Verify ScaleObject instance (`kedaScaledObjectName`) created successfully, including enabling user workload monitoring, deploying a test application that exposes Prometheus metrics, creating a `ServiceMonitor`, `ServiceAccount`, `Secret`, `ClusterRole`, `ClusterRoleBinding`, `TriggerAuthentication`, and finally the `ScaledObject`.

**VerifyNROPSuite() steps:**
1. Verify that `vcoreparams.NROPNamespace` exists.
2. Verify NUMA Resources Operator deployment succeeded.
3. Verify NUMA Resources Operator Custom Resource deployment (`vcoreparams.NROPInstanceName`).
4. Verify NUMA-aware secondary pod scheduler (`vcoreparams.NumaAwareSecondarySchedulerName`) is configured, including creating the scheduler if it doesn't exist and verifying its deployment.
5. Verify scheduling workloads with the NUMA-aware scheduler by deploying a test workload (`vcoreparams.NumaWorkloadName`) with the `NumaAwareSchedulerName` and verifying that it is scheduled.
