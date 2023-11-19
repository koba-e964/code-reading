use serde::{Deserialize, Serialize};

use crate::next_permutation::next_permutation;

#[derive(Serialize, Deserialize, Debug)]
pub struct MonoidRep {
    pub table: Vec<String>,
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

#[derive(Debug, Clone, PartialEq, Eq, Hash, PartialOrd, Ord)]
pub struct Monoid {
    table: Vec<Vec<usize>>,
}

impl Monoid {
    pub fn new(table: Vec<Vec<usize>>) -> Option<Self> {
        let result = Self { table };
        if !result.is_valid() {
            return None;
        }
        Some(result.normalize())
    }
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
        let mut max_table = (vec![], self.table.clone());
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
                    table[perm[i]][perm[j]] = perm[self.table[i][j]];
                }
            }
            let mut tie_break = vec![0; n * n];
            for i in 0..n {
                for j in 0..n {
                    tie_break[i * n + j] = table[n - 1 - i][n - 1 - j];
                }
            }
            for i in 0..idem {
                if table[i][i] != i {
                    ok = false;
                }
            }
            if ok {
                max_table = std::cmp::max(max_table.clone(), (tie_break, table));
            }
            if next_permutation(&mut perm[1..]).is_none() {
                break;
            }
        }
        Self { table: max_table.1 }
    }
    /// Is `self` normalized?
    ///
    /// - The identity element (`e`) must be the 0th element.
    /// - All idempotent elements must come first.
    /// - `rev(1d(a.table))` should be the greatest possible.
    pub fn is_normalized(&self) -> bool {
        let cp = self.normalize();
        cp == *self
    }
}
