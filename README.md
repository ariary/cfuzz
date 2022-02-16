# cfuzz

<div align=center>
<img src= https://github.com/ariary/cfuzz/blob/main/img/cfuzz-logo.png width=150>
</div>
<br>


The same thing as [`wfuzz`](https://github.com/xmendez/wfuzz) **but for command line fuzzing ~> Fuzz your command line execution and filter results.**
<br>*Also a good friend for bruteforcing*

## Usage

Indicate the command containing the fuzzing part with the kyword `FUZZ`, the wordlist and let's get it:
```shell
export CFUZZ_CMD="echo FUZZ" # Example with echo command, I haven't found better
cfuzz -w [wordlist]
```

### Filter result

#### By command output

#### By command execution time

### Configure

#### Command input

#### `cfuzz` run
