This is a really time protocol buffer library primarirly built for serial communication across MCUs and device.
It has almost zero overhead and known for its simplicity and avoiding chained JSON-like key-value strucutures.
It is also intentionally designed to have a maximum length of 256 bytes.

### Design

Types are defined as follow

```
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
  score: float;
}
```

And on parsing it, it is transmitted as a single line of max size 256 bytes with each value seperated by `;`, with the line preceded by the count of the values in `[<count>]`
No keys are relevant since it is generated and used. so for instance on the pipe the data looks like;

```
[3]3.324;1.444;0.999

[4]T;1239998;Manan Junior;2333.111
```

upon generation, it assigns the values to the respective structs.

in C:

```
TinyBuffer_t tb = tb_init();
char* buf[256];
// fetch the data however
Vector_t v;

tb_vector_parse(buf, *v);

char* buf[256];
// fetch the data however
tb_config_parse(buf, *v);
```

in typescript:

```
import { VectorParse, ConfigParse } from "./generated.ts"
buf = Buffer.alloc(256)
// fetch the data however
v = VectorParse(buf)

buf = Buffer.alloc(256)
// fetch the data however
conf = ConfigParse(buf)
```

More version will be done later
