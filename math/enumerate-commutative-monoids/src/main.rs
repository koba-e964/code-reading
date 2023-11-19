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

#[derive(Serialize, Deserialize, Debug)]
struct MonoidRep {
    table: Vec<String>,
}

#[derive(Debug, Clone)]
struct Monoid {
    table: Vec<Vec<usize>>,
}

impl From<MonoidRep> for Monoid {
    fn from(a: MonoidRep) -> Self {
        let n = a.table.len();
        let mut table = vec![vec![0; n]; n];
        let s = b"0123456789";
        for i in 0..n {
            let u = a.table[i].as_bytes();
            for j in 0..n {
                table[i][j] = s.iter().cloned().position(|x| x == u[j]).unwrap();
            }
        }
        Self { table }
    }
}

impl From<Monoid> for MonoidRep {
    fn from(a: Monoid) -> Self {
        let n = a.table.len();
        let mut table = vec!["".to_string(); n];
        let s = b"0123456789";
        for i in 0..n {
            for j in 0..n {
                table[i].push(s[a.table[i][j]] as char);
            }
        }
        Self { table }
    }
}

fn main() {
    let contents = fs::read_to_string(OUTPUT_FILE).expect("Should have been able to read the file");

    let rep: ComMonsRep = serde_json::from_str(&contents).unwrap();
    eprintln!("rep = {:?}", rep);
    let n = 3;
    for n in 1..=n {}
}
