# MiniBuf

**minibuf** is a minimal, high-performance serialization library designed for MCU to MCU and device communication. It focuses on simplicity, predictability, and bounded memory usage, making it ideal for constrained embedded systems. Unlike JSON or protobuf, Buffer avoids keys and hierarchical structures, relying on a fixed-order, schema-driven line format with a maximum payload of 256 bytes.


> [!IMPORTANT]  
> Currently Its best to avoid using strings especially strings with `;` within. 

### Principles

1. deterministic parsing, no dynamic allocations with low overhead.
2. It has a minimal footprint, C-friendly API.
3. the order and type of fields are fixed with deterministic parsing and entirely schema-driven.
4. It has a fixed buffer of 256-byte bounded payload for safe MCU stack usage.
5. fields can have defaults if missing during schema evolution.
6. It is human-readable for easy debugging over UART/Serial.

<!-- - **Optional CRC**: protect against transmission errors. -->

### Schema Definition

Schemas are simple, human-readable, and used to generate parsers:

```buffer
config float_precision = 3;

Vector {
  x: float;
  y: float;
  z: float;
}

Config {
  auto_restart: bool;
  id: number;
  user_name: string;
  score: float = 0.0;      // default value for bachward compatibility
}
```

- Note that `float_precision` specifies serialization rounding.

### Transmission Format

1. Each line begins with `[<count>]`, the number of values transmitted.
2. Fields are **semicolon-separated**.

**Example:**

```
[3]3.324;1.444;0.999
[4]T;1239998;Amar Jay;2333.111
```

<!--[4]T;1239998;Manan Junior;2333.111*AF   // CRC appended -->

## Usage
#### Compilation
To generate the parser code from the schema, use the `minibufc` tool:

```bash
minibufc schema.mb -o types/minibuf --c --ts # generates the types/minibuf.c and types/minibuf.h for C and types/minibuf.ts for TypeScript
```

#### C

```c
#include "minibuf.h"

MiniBuf_t mb = mb_init();
char buf[256];

// parse Vector
Vector_t v;
int err = mb_vector_parse(buf, &v);
if (err != MB_OK) {
    // handle errors
}

// parse Config
Config_t conf;
err = mb_config_parse(buf, &conf);
if (err != MB_OK) {
    // handle errors
}
```

- Note that `minibuf.c` and `minibuf.h` is a single generated c and header file containing everything.

#### TypeScript

```ts
import { VectorParse, ConfigParse } from "./minibuf.ts";

const buf = Buffer.alloc(256);

try {
  const v = VectorParse(buf);
  // parse Config
  const conf = ConfigParse(buf);
} catch (error) {
  console.error(error);
}
```


### Schema Rules
| Rule |                                 |
|-------|--------------------------------------------------|
| Adding fields | append new fields to the end of the schema. Older parsers will ignore them, and newer parsers will apply default values if missing. |
| Removing fields | fields should be deprecated but kept in transmitted data for backward compatibility for minute changes but for massive ones, the schema should be redone. |
| Changing types | discouraged; if necessary, treat as a new schema version. |
