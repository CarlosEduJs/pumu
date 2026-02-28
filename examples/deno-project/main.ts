import { join } from "@std/path";
import { exists } from "@std/fs";

const p = join(".", "hello");
console.log("exists:", await exists(p));
