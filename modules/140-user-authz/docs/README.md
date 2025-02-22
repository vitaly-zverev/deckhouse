---
title: "The user-authz module"
---

The module generates RBAC for users and implements the basic multi-tenancy mode with namespace-based access.

Also, it implements the role-based subsystem for end-to-end authorization, thereby extending the functionality of the standard RBAC mechanism.

All the configuration of access rights is performed using [Custom Resources](cr.html).

## Module features

- Manages user and group access control using Kubernetes RBAC;
- Manages access to scaling tools (the `allowScale` parameter of the [`ClusterAuthorizationRule`](cr.html#clusterauthorizationrule) Custom Resource);
- Manages access to port forwarding (the `portForwarding` parameter of the [`ClusterAuthorizationRule`](cr.html#clusterauthorizationrule) Custom Resource);
- Manages the list of allowed namespaces as regular expressions (the `limitNamespaces` parameter of the [`ClusterAuthorizationRule`](cr.html#clusterauthorizationrule) Custom Resource);
- Manages access to system namespaces such as `kube-system`, etc., (the `allowAccessToSystemNamespaces` parameter of the [`ClusterAuthorizationRule`](cr.html#clusterauthorizationrule) Custom Resource);

## Role model

In addition to the RBAC, you can use a set of high-level roles in the module:
- `User` — has access to information about all objects (including viewing pod logs) but cannot exec into containers, read secrets, and perform port-forwarding;
- `PrivilegedUser` — the same as `User` + can exec into containers, read secrets, and delete pods (and thus, restart them);
- `Editor` — is the same as `PrivilegedUser` + can create and edit all objects that are usually required for application tasks.
- `Admin` — the same as `Editor` + can delete service objects (auxiliary resources such as `ReplicaSet`, `certmanager.k8s.io/challenges` and `certmanager.k8s.io/orders`);
- `ClusterEditor` — the same as `Editor` + can manage a limited set of `cluster-wide` objects that can be used in application tasks (`ClusterXXXMetric`, `KeepalivedInstance`, `DaemonSet`, etc.). This role is best suited for cluster operators.
- `ClusterAdmin` — the same as both `ClusterEditor` and `Admin` + can manage `cluster-wide` service objects (e.g.,  `MachineSets`, `Machines`, `OpenstackInstanceClasses`..., as well as `ClusterAuthorizationRule`, `ClusterRoleBindings` and `ClusterRole`). This role is best suited for cluster administrators. **Note** that since `ClusterAdmin` can edit `ClusterRoleBindings`, he can **broaden his privileges within the cluster**;
- `SuperAdmin` — can perform any actions with any objects (note that [`limitNamespaces`](#module-features) restrictions remain valid).

## Implementation nuances

> **Caution!** Currently, the multi-tenancy mode (namespace-based authorization) is implemented according to a temporary scheme and **isn't guaranteed to be entirely safe and secure**!

The `allowAccessToSystemNamespaces` and `limitNamespaces` options in the CR will no longer be applied if the authorization system's webhook is unavailable for some reason. As a result, users will have access to all namespaces. After the webhook availability is restored, the options will become relevant again.

## Default access list for each role:

`verbs` aliases:
<!-- start user-authz roles placeholder -->
* read - `get`, `list`, `watch`
* read-write - `get`, `list`, `watch`, `create`, `delete`, `deletecollection`, `patch`, `update`
* write - `create`, `delete`, `deletecollection`, `patch`, `update`

```yaml
Role `User`:
    read:
        - apiextensions.k8s.io/customresourcedefinitions
        - apps/daemonsets
        - apps/deployments
        - apps/replicasets
        - apps/statefulsets
        - autoscaling.k8s.io/verticalpodautoscalers
        - autoscaling/horizontalpodautoscalers
        - batch/cronjobs
        - batch/jobs
        - configmaps
        - discovery.k8s.io/endpointslices
        - endpoints
        - events
        - events.k8s.io/events
        - extensions/daemonsets
        - extensions/deployments
        - extensions/ingresses
        - extensions/replicasets
        - extensions/replicationcontrollers
        - limitranges
        - metrics.k8s.io/nodes
        - metrics.k8s.io/pods
        - namespaces
        - networking.k8s.io/ingresses
        - networking.k8s.io/networkpolicies
        - nodes
        - persistentvolumeclaims
        - persistentvolumes
        - pods
        - pods/log
        - policy/poddisruptionbudgets
        - rbac.authorization.k8s.io/rolebindings
        - rbac.authorization.k8s.io/roles
        - replicationcontrollers
        - resourcequotas
        - serviceaccounts
        - services
        - storage.k8s.io/storageclasses
    
Role `PrivilegedUser` (and all rules from `User`):
    create,get:
        - pods/attach
        - pods/exec
    delete,deletecollection:
        - pods
    read:
        - secrets
    
Role `Editor` (and all rules from `User`, `PrivilegedUser`):
    read-write:
        - apps/deployments
        - apps/statefulsets
        - autoscaling.k8s.io/verticalpodautoscalers
        - autoscaling/horizontalpodautoscalers
        - batch/cronjobs
        - batch/jobs
        - configmaps
        - discovery.k8s.io/endpointslices
        - endpoints
        - extensions/deployments
        - extensions/ingresses
        - networking.k8s.io/ingresses
        - persistentvolumeclaims
        - policy/poddisruptionbudgets
        - serviceaccounts
        - services
    write:
        - secrets
    
Role `Admin` (and all rules from `User`, `PrivilegedUser`, `Editor`):
    create,patch,update:
        - pods
    delete,deletecollection:
        - apps/replicasets
        - extensions/replicasets
    
Role `ClusterEditor` (and all rules from `User`, `PrivilegedUser`, `Editor`):
    read:
        - rbac.authorization.k8s.io/clusterrolebindings
        - rbac.authorization.k8s.io/clusterroles
    write:
        - apiextensions.k8s.io/customresourcedefinitions
        - apps/daemonsets
        - extensions/daemonsets
        - storage.k8s.io/storageclasses
    
Role `ClusterAdmin` (and all rules from `User`, `PrivilegedUser`, `Editor`, `Admin`, `ClusterEditor`):
    read-write:
        - deckhouse.io/clusterauthorizationrules
    write:
        - limitranges
        - namespaces
        - networking.k8s.io/networkpolicies
        - rbac.authorization.k8s.io/clusterrolebindings
        - rbac.authorization.k8s.io/clusterroles
        - rbac.authorization.k8s.io/rolebindings
        - rbac.authorization.k8s.io/roles
        - resourcequotas
    
```
<!-- end user-authz roles placeholder -->

You can get additional list of access rules for module role from cluster ([existing user defined rules](usage.html#customizing-rights-of-high-level-roles) and non-default rules from other deckhouse modules):
```bash
D8_ROLE_NAME=Editor
kubectl get clusterrole -A -o jsonpath="{range .items[?(@.metadata.annotations.user-authz\.deckhouse\.io/access-level=='$D8_ROLE_NAME')]}{.rules}{'\n'}{end}" | jq -s add
```
