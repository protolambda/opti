# Opti

Optimistic rollup tech, minimal and generic.

Toy project to explore supporting native EVM rollups,
without the overhead of custom compilers, "L2 operating system" or other magic.
*Keep it stupid simple*.

Goals:
- Full EVM support
- Split EVM in steps, heavy opcodes in multiple steps
- Multi-transaction fraud proofs
- Transaction validation part of fraud-proof
- Interactive search through steps to execute fraud-proof
- Build around Geth EVM tracer, to run EVM test suites and fuzzers


## License

MIT, see [`LICENSE`](./LICENSE) file.