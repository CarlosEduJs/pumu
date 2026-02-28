use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
struct Example {
    name: String,
    value: i32,
}

fn main() {
    let e = Example { name: "pumu".into(), value: 42 };
    let json = serde_json::to_string_pretty(&e).unwrap();
    println!("{json}");
}
