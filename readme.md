# Rattus

```Rattus``` is a lightweight credentials provisioning tool focused on:

- simplicity
- containers support
- repeatable workflow for every credentials provider
- fully configurable through environment variables or flags
- template support

```Rattus``` is designed, to provide the same workflow for credentials provisioning at different environments. 
For example, you have local development environment, that runs under the Kubernetes cluster and you store all your secrets at Vault or at K8S secrets.
But your production environment, deployed at AWS ECS and you can`t use same credential provisioning workflow at both environments.
```Rattus``` fixes that issue, and you can use the same command to retrieve credentials or generating configuration files in every environment.
```Rattus``` is designed for a be configured through environment variables. Because with environment variables you can easily change workflow at different environments, without changing application initialization logic.

# Usage example

Create a shell script, that will be launch at your application startup with followed content:
```bash
#!/bin/sh
/bin/rattus > /app/config.json
```
And that's all! ``Rattus`` will get credentials, render template file, and save the output to application config.
Now you can use this script in every environment, and you will get the same credentials provisioning workflow. All that you need to change - environment variables, that can be easily changed.

See more [examples](https://github.com/rma945/rattus/examples)

# Support credential providers

## Hashicorp Vault

Rattus support [Vault](https://github.com/hashicorp/vault) througt followed auth methods: 

- [kubernetes auth provider](https://www.vaultproject.io/docs/auth/kubernetes/)
- [tokens](https://www.vaultproject.io/docs/concepts/tokens/)

## AWS Secret manager

Rattus support [AWS secret manager](https://aws.amazon.com/secrets-manager/) throught:

- [AWS IAM roles](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles.html)
- [AWS Credential](https://docs.aws.amazon.com/general/latest/gr/aws-security-credentials.html)
