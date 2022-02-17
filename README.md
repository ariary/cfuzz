# cfuzz

<div align=center>
<img src= https://github.com/ariary/cfuzz/blob/main/img/CF-logo.png width=300>
</div>
<br>


The same thing as [`wfuzz`](https://github.com/xmendez/wfuzz) **but for command line fuzzing. This enables to fuzz command line execution and filter results.**
<br>*Also a good friend for bruteforcing*

**Why?**<br>
To perform fuzzing or bruteforcing we have plenty of awesome tools ([`fuff`](https://github.com/ffuf/ffuf) and [`wfuzz`](https://github.com/xmendez/wfuzz) for web fuzzing, [`hydra`](https://github.com/vanhauser-thc/thc-hydra) for network bruteforcing, to mention just a few). **`cfuzz`** is a tool that propose a different approach with a step-back. **The aim is to be able to fuzz/bruteforce anything that can be  transcribed in command line**.

Consequently, `cfuzz` can be seen either as an alternative of these tools for simple use case or an extension cause it handles a huge range of use case

*Idea origin: when bruteforcing ipmi service to enumerate users. 3 options: use `msfconsole`, write module to `hydra`, manually or programmaticaly parse `ipmitool` tool output*

## Usage

Indicate the command containing the fuzzing part with the keyword `FUZZ`, the wordlist and let's get it:
```shell
export CFUZZ_CMD="printf FUZZ | sudo -S id" # Example bruteforcing sudo password, I haven't found better
cfuzz -w [wordlist]
```

### Filter result

#### By command output

#### By command return code

#### By command execution time

### Configure

#### Command input

#### `cfuzz` run
