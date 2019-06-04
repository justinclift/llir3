This is the Go code from Step 3 of:

  https://github.com/llir/llvm/issues/86#issuecomment-498357924

Being worked on to add enough DWARF metadata, so that `-g` on the clang
command line works and embeds decent ".debug_*" sections.

    $ clang --target=x86_64-pc-linux-gnu -g -Wno-override-module -o foo foo.ll

The goal is to have DWARF added to Wasm binaries, so a working debugger for
Wasm can be created.

A working command to generate Wasm from the .ll is:

    $ clang --compile --target=wasm32-unknown-unknown-wasm -g -Wno-override-module -o foo.wasm foo.ll

Note that the `--compile` option there is required.  Without it, clang will
automatically attempt to link the Wasm for the current system (and fail):

```
$ clang --target=wasm32-unknown-unknown-wasm -g -Wno-override-module -o foo.wasm foo.ll
wasm-ld: error: unknown file type: /lib/crt1.o
wasm-ld: error: unable to find library -lc
wasm-ld: error: cannot open /opt/llvm8-wasm/lib/clang/8.0.1/lib/libclang_rt.builtins-wasm32.a: No such file or directory
clang-8: error: lld command failed with exit code 1 (use -v to see invocation)
```
