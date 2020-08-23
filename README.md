WIP

Done:

- All services deployed, even without all fields specified using default values (constant package)
- fix double route creation.
- use constants across all hardcoded.
- error on deletion, (might even crash, namespace not found)
- allow passing of tolerations(to choose on which nodes to be deployed)
- allow annotations
- add Pods deployed to on status
- add svc type deployed to status
- add if route deployed
- check if openshift, deploy route only when true
- check if openshift, if not deploy Nodeport
TBD:
- allow route TLS 