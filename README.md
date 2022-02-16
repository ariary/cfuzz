# cfuzz

<div align=center>
<img src= https://github.com/ariary/cfuzz/blob/main/img/CF-logo.png width=300>
</div>
<br>


The same thing as [`wfuzz`](https://github.com/xmendez/wfuzz) **but for command line fuzzing ~> Fuzz your command line execution and filter results.**
<br>*Also a good friend for bruteforcing*


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
