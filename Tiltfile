docker_build('khos2ow/ratelimiter', '.', dockerfile='Dockerfile')

k8s_yaml([
    'deploy/rbac.yaml',
    'deploy/config.yaml',
    'deploy/secret.yaml',
    'deploy/deployment.yaml',
    'deploy/service.yaml',
    'deploy/ingress.yaml',
])

k8s_resource('rate-limiter', port_forwards=8080)
