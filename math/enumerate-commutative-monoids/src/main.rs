use enumerate_commutative_monoids::{Monoid, MonoidRep};
use serde::{Deserialize, Serialize};
use std::collections::HashSet;
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

fn dfs1(n: usize, idem: usize, x: usize, table: &mut [Vec<usize>], set: &mut HashSet<Monoid>) {
    if x >= n {
        dfs2(n, 1, 2, table, set);
        return;
    }
    for i in 0..n {
        if x != i {
            table[x][x] = i;
            dfs1(n, idem, x + 1, table, set);
        }
    }
}

fn dfs2(n: usize, x: usize, y: usize, table: &mut [Vec<usize>], set: &mut HashSet<Monoid>) {
    if x >= n {
        if let Some(mon) = Monoid::new(table.to_vec()) {
            set.insert(mon);
        }
        return;
    }
    if y >= n {
        dfs2(n, x + 1, x + 2, table, set);
        return;
    }
    for i in 0..n {
        table[x][y] = i;
        table[y][x] = i;
        dfs2(n, x, y + 1, table, set);
    }
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

    for n in 3..=4 {
        for idem in 1..n + 1 {
            let mut set = HashSet::new();
            let mut table = vec![vec![0; n]; n];
            for i in 0..idem {
                table[i][i] = i;
            }
            for i in 0..n {
                table[i][0] = i;
                table[0][i] = i;
            }
            dfs1(n, idem, idem, &mut table, &mut set);
            for m in set {
                eprintln!("{} {} => {:?} {}", n, idem, m, m.is_normalized());
            }
        }
    }
}
