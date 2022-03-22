/* https://arxiv.org/pdf/1605.05274.pdf */
interface Z {} interface N<x> {} interface L<x> {}
interface Qlr<x> {} interface Qrl<x> {}
interface E<x> extends
		   Qlr<N<?super Qr<?super E<?super E<?super x>>>>>,
		   Qrl<N<?super Ql<?super E<?super E<?super x>>>>> {}
interface Ql<x> extends
		    L<N<?super Ql<?super L<?super N<?super x>>>>>,
		    E<Qlr<?super N<?super x>>> {}
interface Qr<x> extends
		    L<N<?super Qr<?super L<?super N<?super x>>>>>,
		    E<Qrl<?super N<?super x>>> {}
class Main {
    L<?super N<?super L<?super N<?super L<?super N<?super
	E<?super E<?super Z>>>>>>>>
	doit(Qr<? super E<? super E<? super Z>>> v) {return v;}
}
