# minibuf

**minibuf** is a minimal, high-performance serialization library designed for MCU-to-MCU and device communication. It focuses on **simplicity, predictability, and bounded memory usage**, making it ideal for constrained embedded systems. Unlike JSON or protobuf, Buffer avoids keys and hierarchical structures, relying on a **fixed-order, schema-driven line format** with a maximum payload of 256 bytes.

---

## Features

- **Ultra-low overhead**: deterministic parsing, no dynamic allocations.
- **Schema-driven**: order and type of fields are fixed; no runtime reflection needed.
- **Bounded size**: max 256-byte payload for safe MCU stack usage.
- **Default values**: fields can have defaults if missing during schema evolution.
- **Schema evolution friendly**: backward-compatible parsing via count-based approach.
- **Multi-language support**: currently C and TypeScript; other languages planned.

<!-- - **Optional CRC**: protect against transmission errors. -->

---

## Schema Definition

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
  score: float = 0.0;      // default value
}
```

**Notes:**

- `float_precision` specifies serialization rounding.
- Default values are applied when a field is missing (helps with backward compatibility).

---

## Transmission Format

- Each line begins with `[<count>]`, the number of values transmitted.
- Fields are **semicolon-separated**.
- Optional CRC can be appended at the end of a line.

**Example:**

```
[3]3.324;1.444;0.999
[4]T;1239998;Amar Jay;2333.111
```

<!--[4]T;1239998;Manan Junior;2333.111*AF   // CRC appended -->

---

## Usage

### C

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

**Error codes:**

- `MB_OK` — parsing successful
- `MB_ERR_COUNT_MISMATCH` — count prefix does not match schema
- `MB_ERR_PARSE_FLOAT` — invalid float value
- `MB_ERR_PARSE_BOOL` — invalid boolean value
- `MB_ERR_CRC_MISMATCH` — CRC check failed

---

### TypeScript

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

---

## Schema Evolution Rules

1. **Adding fields**: append new fields to the end of the schema. Older parsers will ignore them, and newer parsers will apply default values if missing.
2. **Removing fields**: fields should be deprecated but kept in transmitted data for backward compatibility for minute changes but for massive ones, the schema should be redone.
3. **Changing types**: discouraged; if necessary, treat as a new schema version.

---

## Design Philosophy

- **Predictable memory usage**: fixed buffer, bounded payload.
- **Deterministic parsing**: avoids dynamic key lookups.
- **MCU-first**: minimal footprint, C-friendly API.
- **Human-readable**: easy debugging over UART/Serial.
