# Release Process

## Prerequisites

### GPG Signing Setup (Optional)

To sign releases with GPG:

1. **Install GPG:**
   ```bash
   brew install gnupg  # macOS
   # or
   apt-get install gnupg  # Linux
   ```

2. **Generate a GPG key (if you don't have one):**
   ```bash
   gpg --full-generate-key
   # Choose: RSA and RSA, 4096 bits, no expiration
   # Use your GitHub email
   ```

3. **Get your GPG fingerprint:**
   ```bash
   gpg --list-secret-keys --keyid-format=long
   # Look for the line like: sec   rsa4096/ABCD1234EFGH5678
   # ABCD1234EFGH5678 is your key ID

   # Get full fingerprint:
   gpg --fingerprint ABCD1234EFGH5678
   ```

4. **Export and add to GitHub:**
   ```bash
   gpg --armor --export ABCD1234EFGH5678
   # Copy output and add to: https://github.com/settings/keys
   ```

5. **Configure Git:**
   ```bash
   git config --global user.signingkey ABCD1234EFGH5678
   git config --global commit.gpgsign true
   ```

## Creating a Release

### With Signing

1. **Set environment variable:**
   ```bash
   export GPG_FINGERPRINT=ABCD1234EFGH5678
   export GITHUB_TOKEN=$(gh auth token)
   ```

2. **Create and push tag:**
   ```bash
   git tag v1.x.x
   git push origin v1.x.x
   ```

3. **Run goreleaser:**
   ```bash
   make release
   ```

   This will:
   - Build binaries for all platforms
   - Generate checksums.txt
   - Sign checksums.txt with GPG (creates checksums.txt.sig)
   - Upload all artifacts to GitHub release

### Without Signing

If GPG_FINGERPRINT is not set, goreleaser will skip signing but still create the release with checksums.

```bash
export GITHUB_TOKEN=$(gh auth token)
git tag v1.x.x
git push origin v1.x.x
make release
```

## Verifying Signed Releases

Users can verify signed releases:

```bash
# Download release files
curl -LO https://github.com/tmc/tmpl/releases/download/v1.x.x/checksums.txt
curl -LO https://github.com/tmc/tmpl/releases/download/v1.x.x/checksums.txt.sig
curl -LO https://github.com/tmc/tmpl/releases/download/v1.x.x/tmpl_linux_amd64

# Import your public key
curl https://github.com/yourusername.gpg | gpg --import

# Verify signature
gpg --verify checksums.txt.sig checksums.txt

# Verify checksum
sha256sum -c checksums.txt --ignore-missing
```

## Release Checklist

- [ ] Update version in code if needed
- [ ] Run `go test` to ensure tests pass
- [ ] Create git tag: `git tag v1.x.x`
- [ ] Push tag: `git push origin v1.x.x`
- [ ] Set GPG_FINGERPRINT if signing
- [ ] Set GITHUB_TOKEN: `export GITHUB_TOKEN=$(gh auth token)`
- [ ] Run `make release`
- [ ] Verify release on GitHub
- [ ] Update README.md checksums with `make docs`
- [ ] Commit and push README.md
