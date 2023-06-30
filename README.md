# Scripty

Generate bash and PowerShell install scripts for binaries hosted on GitHub releases.

## Use

```bash
scripty generate github
```

This will prompt for the following information

- Repo Location: The GitHub repo in the form USERNAME/REPO
- Archive Name: The template for artifact hosted on GitHub. May use variables in the form {VARIABLE_NAME}. Possible variables:
  - {PLATFORM}: The OS of the machine as reported from running `uname -s` on Linux/Darwin machines and `Windows` on Windows machines.
  - {ARCH}: The architecture of the machine as reported from running `uname -m` on Linux/Darwin machines and the environment variable `PROCESSOR_ARCHITECTURE` on Windows.
  - {TAG}: Resolves to the GitHub release.
  - {ARCHIVE_EXT}: Resolves to `.tar.gz` on Linux/Darwin machines and `.zip` on Windows.
  - {BINARY_EXT}: Resolves to blank on Linux/Darwin and `.exe` on Windows.
- Linux/Darwin Platform Mapping: Comma separated list of `KEY=VALUE` pairs representing a map of `uname -s` platform names to what is listed in GitHub releases. Example, `Darwin=mac` can be used to generate install scripts for GitHub projects that use the term `mac` in their release artifacts instead of `Darwin`.
- Linux/Darwin Architecture Mapping: Comma separated list of `KEY=VALUE` paris representing architecture mappings. Example, `x86_64=amd64` can be used to generate install scripts for GitHub projects that use the term `amd64` in their release artifacts instead of `x86_64`.
- Windows Platform Mapping: Map the term `Windows` to a term that is used in the GitHub release artifact. Example, `Windows=win`
- Windows Architecture Mapping: Comma separated list of `KEY=VALUE` pairs to map the environment variable `PROCESSOR_ARCHITECTURE` on Windows to a term used by the GitHub project the script is bein generated for. Example, `AMD64=x86_64` is useful to map the value `AMD64` to `x86_64`.
- Dest: Where to create the `install.sh` and `install.ps1` install scripts.

**Values used to generate the install scripts for this project**

- Repo Location: `dworthen/goscripty`
- Archive Name: `scripty_{PLATFORM}_{ARCH}{ARCHIVE_EXT}`. This matches the format I use to release binary archives: https://github.com/dworthen/goscripty/releases
- Linux/Darwin Platform Mapping: ``. Left blank to use the default of `Linux=Linux,Darwin=Darwin`. Again, this matches the format used in the GitHub releases page.
- Linux/Darwin Architecture Mapping: ``. Left blank to use the default of `x86_64=x86_64`.
- Windows Platform Mapping: ``. Left blank to default to `Windows=Windows`.
- Windows Architecture Mapping: ``. Left blank to default to `AMD64=x86_64`.
- Dest: `./scripts`

This generates install scripts in the [`./scripts`](scripts/).

## Install

### Windows

```powershell
curl -sSfL https://raw.githubusercontent.com/dworthen/goscripty/main/scripts/install.ps1 | pwsh -Command -
```

or with additional flags:

```powershell
curl -sSfL https://raw.githubusercontent.com/dworthen/goscripty/main/scripts/install.ps1 -o install.ps1 &&
pwsh -File install.ps1 -force -tag v0.0.1 -to ~/bin &&
rm install.ps1
```

### Linux/Darwin

```bash
curl -sSfL https://raw.githubusercontent.com/dworthen/goscripty/main/scripts/install.sh | bash
```

or with additional flags:

```bash
curl -sSfL https://raw.githubusercontent.com/dworthen/goscripty/main/scripts/install.sh | bash -s -- --force --tag v0.0.1 --to ~/bin
```
