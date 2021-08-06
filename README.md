deployment process of gotti.dev consists of g8s(gotti house k8s) and GitHub Actions.

GitHub Actions pushes built docker image that serve homepage made with hugo, to GitHub Packages.

GitHub Actions also replace image tag in manifests/nginx-static.yml with commit id and pushes to deployment branch.

After pushed, ArgoCD on g8s will detect changes and deploy.
