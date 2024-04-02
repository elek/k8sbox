# k8sbox

This script prints out all the nodes in a K8s cluster, with the running pods, and summarizes requested/limimt/usage of CPU/MEM.

Can be used to create more reasonable pools...

Example summary:

```
+-------------------------------------------------------+----------------------------------------------+
|                      CPU (CORES)                      |                   MEM (GB)                   |
| ALLOCATABLE | USED        | REQUEST     | LIMIT       | ALLOCATABLE | USED     | REQUEST  | LIMIT    |
+-------------+-------------+-------------+-------------+-------------+----------+----------+----------+
| 55.37       | 7.42        | 19.04       | 22.741      | 43.254      | 16.234   | 21.511   | 22.741   |
+-------------+-------------+-------------+-------------+-------------+----------+----------+----------+
```

But similar information is available for each pod.

## Usage

**It requires kubectl proxy!!!**

```
kubectl proxy
k8sbox --proxy-port=8081
```
