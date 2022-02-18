# cfuzz

<div align=center>
<img src= https://github.com/ariary/cfuzz/blob/main/img/CF-logo.png width=300>

<br>


The same thing as [`wfuzz`](https://github.com/xmendez/wfuzz) **but for command line fuzzing. This enables to fuzz command line execution and filter results.**
<br>*Also a good friend for bruteforcing*
</div>

**Why?**<br>
To perform fuzzing or bruteforcing we have plenty of awesome tools ([`fuff`](https://github.com/ffuf/ffuf) and [`wfuzz`](https://github.com/xmendez/wfuzz) for web fuzzing, [`hydra`](https://github.com/vanhauser-thc/thc-hydra) for network bruteforcing, to mention just a few). **`cfuzz`** is a tool that propose a different approach with a step-back. **The aim is to be able to fuzz/bruteforce anything that can be  transcribed in command line**.

Consequently, `cfuzz` can be seen either as an alternative of these tools for simple use case or an extension cause it handles a huge range of use case

*Idea origin: when bruteforcing ipmi service to enumerate users. 3 options: use `msfconsole`, write module for `hydra`, manually or programmaticaly parse `ipmitool` tool output*

## Usage

Indicate the command containing the fuzzing part with the keyword `FUZZ`, the wordlist and let's get it:
```shell
export CFUZZ_CMD="printf FUZZ | sudo -S id" # Example bruteforcing sudo password, I haven't found better
cfuzz -w [wordlist]
```

Or in one line:
```Shell
# example for subdomain enum
cfuzz -w [wordlist] -t 5 ping -c 4 FUZZ.domain.net
```

Also, fuzzing  command `stdin` is possible by adding `--stdin-fuzzing [INPUT_WITH_CFUZZ_KEYWORD]`

### Filter result

Choose an execution element to display and add filters to select specific execution characteristics

#### By command output

Use the flag `-oc`  to display stdout number of character, `-ec` for stderr

Additionnaly you can apply filter:
* *Display only entry with more than n characters*: `--omin n` (Conversely `--omax`)
* *Display only entry with exactly n characters*: `--oeq n`

For stderr flag replace `o` by `e`

#### By command return code

Use the flag `-c` to display result regarding exit code of command execution.

Additionnaly you can apply filter:
* *Display only entry with 0 exit code*: `--success` (Conversely `--failure`)

#### By command execution time

Use the flag `-t` to display result regarding time execution of command.

Additionnaly you can apply filter:
* *Display only entry with more an exectuion time greater than n second*: `--tmin n` (Conversely `--tmax`)
* *Display only entry with exactly n seconds*: `--teq n`

### Configure

* Command input (`-i`, `--input`), to fuzz in stdin use `--stdin-fuzzing` 
* Timeout for command execution process (`-to`, `--timeout`)
* Delay  between each command execution (`-d`, `--delay`)
* Change `cfuzz` Keyword (`-k`, `--keyword`)
