use enumerate_commutative_monoids::{Monoid, MonoidRep};
use serde::{Deserialize, Serialize};
use std::fs;

const OUTPUT_FILE: &str = "commutative-monoids.json";

#[derive(Serialize, Deserialize, Debug)]
struct ComMonsRep {
    orderwise: Vec<OrderwiseRep>,
}

#[derive(Serialize, Deserialize, Debug)]
struct OrderwiseRep {
    order: u16,
    idemwise: Vec<IdemwiseRep>,
}

#[derive(Serialize, Deserialize, Debug)]
struct IdemwiseRep {
    idem: u16,
    monoids: Vec<MonoidRep>,
}

fn main() {
    let contents = fs::read_to_string(OUTPUT_FILE).expect("Should have been able to read the file");

    let rep: ComMonsRep = serde_json::from_str(&contents).unwrap();
    eprintln!("rep = {:?}", rep);
    for v in rep.orderwise {
        for v in v.idemwise {
            for v in v.monoids {
                let v: Monoid = v.into();
                eprintln!("{:?} {}", v, v.is_normalized());
            }
        }
    }
}
