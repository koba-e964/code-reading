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

#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Monoid {
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

// Returns the least index of elements that are modified, wrapped with Some.
// If the entire array is reversed, it returns None instead.
// v's elements must be pairwise distinct.
fn next_permutation<T: Ord>(v: &mut [T]) -> Option<usize> {
    let mut tail_dec: usize = 1;
    let n = v.len();
    while tail_dec < n {
        if v[n - tail_dec - 1] > v[n - tail_dec] {
            tail_dec += 1;
        } else {
            break;
        }
    }
    // v[n - tail_dec .. n] is strictly decreasing
    if tail_dec < n {
        let x = n - tail_dec - 1;
        let mut y = n;
        {
            let pivot = &v[x];
            for i in (n - tail_dec..n).rev() {
                if v[i] > *pivot {
                    y = i;
                    break;
                }
            }
            assert!(y < n);
        }
        v.swap(x, y);
    }
    v[n - tail_dec..].reverse();
    if tail_dec < n {
        Some(n - tail_dec - 1)
    } else {
        None
    }
}

impl Monoid {
    pub fn len(&self) -> usize {
        self.table.len()
    }
    /// Does `self` satisfy unit, commutativity and associativity?
    ///
    /// Time complexity: O(n^3)
    pub fn is_valid(&self) -> bool {
        // unit must be the 0th element
        let n = self.len();
        for i in 0..n {
            if self.table[0][i] != i {
                return false;
            }
        }
        for i in 0..n {
            for j in 0..i {
                if self.table[i][j] != self.table[j][i] {
                    return false;
                }
            }
        }
        for i in 1..n {
            for j in 1..n {
                for k in 1..n {
                    let left_assoc = self.table[self.table[i][j]][k];
                    let right_assoc = self.table[i][self.table[j][k]];
                    if left_assoc != right_assoc {
                        return false;
                    }
                }
            }
        }
        true
    }
    pub fn normalize(&self) -> Self {
        let n = self.len();
        if n == 1 {
            return self.clone();
        }
        let n = self.len();
        let mut perm: Vec<_> = (0..n).collect();
        let mut min_table = self.table.clone();
        let mut idem = 0;
        for i in 0..n {
            if self.table[i][i] == i {
                idem += 1;
            }
        }
        // O((n-1)!)
        loop {
            let mut ok = true;
            // idempotents
            let mut table = vec![vec![0; n]; n];
            for i in 0..n {
                for j in 0..n {
                    table[perm[i]][perm[j]] = self.table[i][j];
                }
            }
            for i in 0..idem {
                if self.table[i][i] != i {
                    ok = false;
                }
            }
            if ok {
                min_table = std::cmp::min(min_table.clone(), table);
            }
            if next_permutation(&mut perm[1..]).is_none() {
                break;
            }
        }
        Self { table: min_table }
    }
    /// Is `self` normalized?
    ///
    /// - The identity element (`e`) must be the 0th element.
    /// - All idempotent elements must come first.
    /// - `a.table` should be the least possible.
    pub fn is_normalized(&self) -> bool {
        let cp = self.normalize();
        cp == *self
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
}
