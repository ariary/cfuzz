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
cfuzz -w [wordlist] -t 5 ping -c 4 FUZZ.domain.net
```

Additionnaly it is possible to:
* **[Filter results](#filter-results)**
* **[Custom displayed field](#displayed-field)**
* **[Configure `cfuzz` run](#cfuzz-run-configuration)**

### Filter results

Additionaly, it is possible to filter displayed results:

**stdout filters:**
```shell
  -omin, --stdout-min         filter to only display if stdout characters number is lesser than n
  -omax, --stdout-max         filter to only display if stdout characters number is greater than n
  -oeq,  --stdout-equal       filter to only display if stdout characters number is equal to n
  -ow,   --stdout-word        filter to only display if stdout cointains specific word
```

**stderr filters:**
```shell
  -emin, --stderr-min         filter to only display if stderr characters number is lesser than n
  -emax, --stderr-max         filter to only display if stderr characters number is greater than n
  -eeq,  --stderr-equal       filter to only display if stderr characters number is equal to n
```

**execution time filters:**
```shell
  -tmin, --time-min           filter to only display if exectuion time is shorter than n seconds
  -tmax, --time-max           filter to only display if exectuion time is longer than n seconds
  -teq,  --time-equal         filter to only display if exectuion time is shorter than n seconds
```

**command exit code filters:**
```shell
  --success                  filter to only display if execution return a zero exit code
  --failure                  filter to only display if execution return a non-zero exit code
```

### `cfuzz` run configuration
To make cfuzz more flexible and adapt to different constraints, many options are possible:
```shell
  -w, --wordlist            wordlist used by fuzzer
  -d, --delay               delay in ms between each thread launching. A thread execute one command. (default: 0)
  -k, --keyword             keyword used to determine which zone to fuzz (default: FUZZ)
  -s, --shell               shell to use for execution (default: /bin/bash)
  -to, --timeout            command execution timeout in s. After reaching it the command is killed. (default: 30)
  -i, --input               provide stdin
  -if, --stdin-fuzzing      fuzz sdtin instead of command line
```

### Displayed field

It is also possible to choose which result field is displayed in `cfuzz` output (also possible to use several):
```shell
  -oc, --stdout              display stdout number of characters
  -ec, --stderr              display stderr number of characters
  -t, --time                 display execution time
  -c, --code                 display exit code
```
