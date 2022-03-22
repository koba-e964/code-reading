// Modification of https://yamasa.hatenablog.jp/entry/20101016/1287227939
use std::sync::atomic::{AtomicBool, AtomicU32, Ordering};
use Ordering::*;

const LOOP_COUNT: i32 = 10_000_000;

static mut counter: i32 = 0;

type AB = AtomicBool;

fn lock(my_flag: &AB, other_flag: &AB,
        other_id: u32) {
    my_flag.store(true, Relaxed);
    turn.store(other_id, Relaxed);
    while other_flag.load(Relaxed) && turn.load(Relaxed) == other_id {}
}
fn unlock(my_flag: &AB) {
    my_flag.store(false, SeqCst);
}

fn task0() {
    for _ in 0..LOOP_COUNT {
        lock(&flag0, &flag1, 1);
        unsafe { counter += 1 };
        unlock(&flag0);
    }
}
fn task1() {
    for _ in 0..LOOP_COUNT {
        lock(&flag1, &flag0, 0);
        unsafe { counter += 1 };
        unlock(&flag1);
    }
}

static turn: AtomicU32 = AtomicU32::new(0);
static flag0: AtomicBool =(AtomicBool::new(false));
static flag1: AtomicBool = (AtomicBool::new(false));

fn main() {
    let thd0 = std::thread::spawn(
        move || task0());
    let thd1 = std::thread::spawn(move || task1());
    thd0.join().unwrap();
    thd1.join().unwrap();
    let value = unsafe { counter };
    println!("counter = {}", value);
    println!("diff = {}", 2 * LOOP_COUNT - value);
}
