use regex_automata::nfa::thompson::{pikevm::PikeVM, WhichCaptures, NFA};
use regex_syntax::{hir::Hir, parse};

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let hir: Hir = parse("a|bcd*")?;

    // https://docs.rs/regex-automata/0.4.5/regex_automata/nfa/thompson/pikevm/struct.PikeVM.html
    let config = NFA::config()
        .nfa_size_limit(Some(1_000))
        .which_captures(WhichCaptures::None);
    let nfa: NFA = NFA::compiler().configure(config).build_from_hir(&hir)?;

    let re: PikeVM = PikeVM::new_from_nfa(nfa)?;

    eprintln!("{:?}", re);

    Ok(())
}
