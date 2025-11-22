import { VectorSerialize, ConfigSerialize, VectorParse, ConfigParse } from "./minibuf.js";

const vector = { x: 1.234, y: 5.678, z: 9.012 };
const config = { auto_restart: false, id: 42, user_name: "test", score: 99.5 };

const vectorStr = VectorSerialize(vector); // "[3]1.234;5.678;9.012"
const configStr = ConfigSerialize(config); // "[4]F;42;test;99.500"
console.log("Serialized Vector:", vectorStr);
console.log("Serialized Config:", configStr);

const parsedVector = VectorParse(vectorStr);
const parsedConfig = ConfigParse(configStr);
console.log("Parsed Vector:", parsedVector);
console.log("Parsed Config:", parsedConfig);

const _parsedVector = JSON.stringify(vector);
const _parsedConfig = JSON.stringify(config);
console.log("Parsed Vector:", _parsedVector);
console.log("Parsed Config:", _parsedConfig);