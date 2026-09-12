# Secrets (SOPS)

Secret values (API keys, DB credentials, tokens) are kept as plain Helm
values files, encrypted at rest with [sops](https://github.com/getsops/sops)
and decrypted on the fly at deploy time — nothing decrypted ever touches
disk or git.

## One-time setup

1. Pick a sops backend (age, GCP KMS, PGP — not decided yet for this repo)
   and create/import the key.
2. Add a `.sops.yaml` at the repo root (or in `deploy/helm/memelo/`) mapping
   `secrets/*.enc.yaml` to that key. Example for age:

   ```yaml
   creation_rules:
     - path_regex: deploy/helm/memelo/secrets/.*\.enc\.yaml$
       age: <age public key>
   ```

## Creating/editing an env's secrets

```bash
cp secrets/test.example.yaml secrets/test.yaml
$EDITOR secrets/test.yaml            # fill in real values
sops -e secrets/test.yaml > secrets/test.enc.yaml
rm secrets/test.yaml                 # don't leave the plaintext around
git add secrets/test.enc.yaml

# to edit an existing encrypted file in place:
sops secrets/test.enc.yaml
```

Repeat for `secrets/production.yaml` -> `secrets/production.enc.yaml`.

Only `*.enc.yaml` files are committed. Plaintext `secrets/*.yaml` (not
`*.example.yaml`) must never be committed — add a `.gitignore` entry for
`secrets/*.yaml` alongside `!secrets/*.example.yaml` once a real file exists
locally.

## Deploying with the decrypted values

```bash
helm upgrade --install memelo deploy/helm/memelo \
  -n memelo-test --create-namespace \
  -f deploy/helm/memelo/values.yaml \
  -f deploy/helm/memelo/values-test.yaml \
  -f <(sops -d deploy/helm/memelo/secrets/test.enc.yaml) \
  --set image.tag=<release tag>
```

Swap `test` for `production` (and the namespace) for the other environment.
