trait Unsigned {
    fn to_i64(self) -> i64;
}

#[derive(Clone, Copy, Debug, Default)]
struct Z;
impl Unsigned for Z {
    fn to_i64(self) -> i64 {
        0
    }
}

#[derive(Clone, Copy, Debug, Default)]
struct One;
impl Unsigned for One {
    fn to_i64(self) -> i64 {
        1
    }
}

#[derive(Clone, Copy, Debug, Default)]
struct B0;
#[derive(Clone, Copy, Debug, Default)]
struct B1;

trait Bit {
    fn to_i64(self) -> i64;
}
impl Bit for B0 {
    fn to_i64(self) -> i64 { 0 }
}
impl Bit for B1 {
    fn to_i64(self) -> i64 { 1 }
}

// 2 * P + B
trait Pos: Unsigned {}

#[derive(Clone, Copy, Debug, Default)]
struct PX<P, B> {
    p: P,
    b: B,
}
impl<P: Pos, B: Bit> Unsigned for PX<P, B> {
    fn to_i64(self) -> i64 { 2 * self.p.to_i64() + self.b.to_i64() }
}
impl<P: Pos, B: Bit> Pos for PX<P, B> {}
impl Pos for One {}

type U0 = Z;
type U1 = One;
type U2 = PX<One, B0>;
type U3 = PX<One, B1>;
type U4 = PX<U2, B0>;
type U5 = PX<U2, B1>;

trait IncrT: Unsigned {
    type Output;
}
impl IncrT for Z {
    type Output = One;
}
impl IncrT for One {
    type Output = PX<One, B0>;
}
impl<P: Pos> IncrT for PX<P, B0> {
    type Output = PX<P, B1>;
}
impl<P: Pos + IncrT> IncrT for PX<P, B1> {
    type Output = PX<Incr<P>, B0>;
}
type Incr<A> = <A as IncrT>::Output;

// Add

trait AddT<B> {
    type Output;
}
// Add(Z, Y) = Y
// Add(One, Y) = Incr(Y)
// Add(PX(P, B), Z) = PX(P, B)
// Add(PX(P, B), One) = Incr(PX(P, B))
// Add(PX(P, B0), PX(Q, B)) = PX(Add(P, Q), B)
// Add(PX(P, B1), PX(Q, B0)) = PX(Add(P, Q), B1)
// Add(PX(P, B1), PX(Q, B1)) = PX(Incr(Add(P, Q)), B0)
impl<Y: Unsigned> AddT<Y> for Z {
    type Output = Y;
}
impl<Y: IncrT> AddT<Y> for One {
    type Output = Incr<Y>;
}
impl<P: Pos, B: Bit> AddT<Z> for PX<P, B> {
    type Output = PX<P, B>;
}
impl<P: Pos + IncrT, B: Bit> AddT<One> for PX<P, B>
where PX<P, B>: IncrT {
    type Output = Incr<PX<P, B>>;
}
impl<P: Pos + AddT<Q>, Q: Pos, B: Bit> AddT<PX<Q, B>> for PX<P, B0> {
    type Output = PX<Add<P, Q>, B>;
}
impl<P: Pos + AddT<Q>, Q: Pos> AddT<PX<Q, B0>> for PX<P, B1> {
    type Output = PX<Add<P, Q>, B1>;
}
impl<P: Pos + AddT<Q>, Q: Pos> AddT<PX<Q, B1>> for PX<P, B1>
where Add<P, Q>: IncrT {
    type Output = PX<Incr<Add<P, Q>>, B0>;
}
type Add<X, Y> = <X as AddT<Y>>::Output;

// Sub(One, Z) = One
// Sub(PX(P, B0), Z) = PX(P, B0)
// Sub(PX(P, B1), Z) = PX(P, B1)
// Sub(One, One) = Z
// Sub(PX(P, B0), One) = PX(Sub(P, One), B1)
// Sub(PX(P, B1), One) = PX(P, B0)
trait SubT<Y> {
    type Output;
}

type Sub<X, Y> = <X as SubT<Y>>::Output;
macro_rules! sub_0 {
    ($ty:ty) => {
        impl SubT<Z> for $ty {
            type Output = Self;
        }
    }
}
sub_0!(Z);
sub_0!(One);
impl<P: Pos, B: Bit> SubT<Z> for PX<P, B> {
    type Output = Self;
}
impl SubT<One> for One {
    type Output = Z;
}

fn main() {
    let c: Add<U4, U5> = Default::default();
    println!("{}", c.to_i64());
}
