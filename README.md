# cfuzz

<div align=center>
<img src= https://github.com/ariary/cfuzz/blob/main/img/CF-logo.png width=300>

<br>


The same thing as [`wfuzz`](https://github.com/xmendez/wfuzz) **but for command line fuzzing. This enables to fuzz any command line execution and filter results.**
<br>*Also a good friend for bruteforcing*
  
<strong><code>{ <a href="#install">Install it</a> ; <a href="#usage">Use it</a> } </code></strong>

</div>

**Why?**<br>
To perform fuzzing or bruteforcing we have plenty of awesome tools ([`fuff`](https://github.com/ffuf/ffuf) and [`wfuzz`](https://github.com/xmendez/wfuzz) for web fuzzing, [`hydra`](https://github.com/vanhauser-thc/thc-hydra) for network bruteforcing, to mention just a few). **`cfuzz`** is a tool that propose a different approach with a step-back. **The aim is to be able to fuzz/bruteforce anything that can be  transcribed in command line**.

Consequently, `cfuzz` can be seen either as an alternative of these tools for simple use case or an extension cause it handles a huge range of use cases

<sub>*Origins of the idea: when bruteforcing ipmi service to enumerate users. 3 options: use `msfconsole`, write module for `hydra`, manually or programmaticaly parse `ipmitool` tool output*</sub>

## Demo
<div align=center>

|user password bruteforcing|
|:---:| 
|![demo](https://github.com/ariary/cfuzz/blob/main/img/cfuzz-user-demo.gif)|

</div>

## Install

From release:
```shell
curl -lO -L -s https://github.com/ariary/cfuzz/releases/latest/download/cfuzz && chmod +x cfuzz
```

With go:
```shell
go install github.com/ariary/cfuzz/cmd/cfuzz@latest
```

## Usage

Indicate:
* the command, with the fuzzing part determined with the keyword `FUZZ`
* the wordlist 

and let's get it!

```shell
export CFUZZ_CMD="printf FUZZ | sudo -S id" # Example bruteforcing user password, I haven't found better
cfuzz -w [wordlist] 
```

Or if you prefer in one line:
```Shell
# example for subdomain enum
cfuzz -w [wordlist] ping -c 4 FUZZ.domain.net
```

Additionnaly it is possible to:
* **[Filter results](#filter-results)**
* **[Custom displayed field](#displayed-field)**
* **[Configure `cfuzz` run](#cfuzz-run-configuration)**
* **[Generate wordlists with AI](#ai-features)**
* **[Use cfuzz as an MCP server](#mcp-server)**

### Filter results

Additionaly, it is possible to filter displayed results:

**stdout filters:**
```shell
  --stdout-min n      show only if stdout character count >= n
  --stdout-max n      show only if stdout character count <= n
  --stdout-eq  n      show only if stdout character count == n
  --stdout-word w     show only if stdout contains word w (repeatable)
```

**stderr filters:**
```shell
  --stderr-min n      show only if stderr character count >= n
  --stderr-max n      show only if stderr character count <= n
  --stderr-eq  n      show only if stderr character count == n
  --stderr-word w     show only if stderr contains word w (repeatable)
```

**execution time filters:**
```shell
  --time-min n        show only if execution time >= n seconds
  --time-max n        show only if execution time <= n seconds
  --time-eq  n        show only if execution time == n seconds
```

**command exit code filters:**
```shell
  --success           show only if execution returns exit code 0
  --failure           show only if execution returns a non-zero exit code
```

To only display results that don't pass the filter use `-H` or `--hide` flag.

### `cfuzz` run configuration
To make cfuzz more flexible and adapt to different constraints, many options are possible:
```shell
  -w, --wordlist        wordlist file(s) for fuzzing (repeatable with --spider)
  -d, --delay           delay in ms between goroutine launches (default: 0)
  -j, --threads         max concurrent workers (default: 50)
  -k, --keyword         keyword to replace in command (default: FUZZ)
  -s, --shell           shell to use for execution (default: /bin/bash)
      --timeout         command execution timeout in seconds (default: 30)
  -i, --input           provide command stdin
      --stdin-fuzzing   fuzz stdin instead of command line
  -m, --spider          fuzz multiple keyword positions (requires multiple -w)
      --stdin-wordlist  read wordlist from cfuzz stdin
```

### Displayed field

It is also possible to choose which result field is displayed in `cfuzz` output (also possible to use several):
```shell
      --stdout-chars    display stdout character count
      --stderr-chars    display stderr character count
  -t, --time            display execution time
  -c, --code            display exit code
      --no-banner       hide banner
  -r, --only-word       print only matched words (no metadata columns)
  -f, --full-output     display full command execution output (can't be combined with other display modes)
```

### AI features

`cfuzz` integrates with Claude (via the Anthropic API) for two AI-powered workflows. Both require the `ANTHROPIC_API_KEY` environment variable to be set.

**AI filter** — describe what an interesting result looks like in plain English; `cfuzz` will ask Claude to evaluate each execution result and only show the ones that match:

```shell
cfuzz -w wordlist.txt --ai-filter "output contains an error about invalid credentials" \
  curl -s http://target/login -d "user=admin&pass=FUZZ"
```

**AI wordlist generation** — generate a context-aware wordlist by describing what you need:

```shell
cfuzz wordlist "default credentials for network switches"
cfuzz wordlist "common web admin paths" -n 50
```

Output is printed to stdout, one entry per line, making it easy to pipe directly into cfuzz:

```shell
cfuzz wordlist "linux privilege escalation binaries" | \
  cfuzz --stdin-wordlist "sudo -l FUZZ 2>/dev/null | grep -v 'not allowed'"
```

### MCP server

`cfuzz` can run as a [Model Context Protocol](https://modelcontextprotocol.io) server, exposing a `fuzz` tool that any MCP-compatible AI assistant can call:

```shell
cfuzz mcp
```

To register with Claude Desktop, add to `~/.claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "cfuzz": { "command": "cfuzz", "args": ["mcp"] }
  }
}
```

The `fuzz` tool accepts: `command` (string), `wordlist` (array of strings), and optional `threads`, `timeout`, `success_only`, `stdout_word`, and `ai_filter` parameters.
