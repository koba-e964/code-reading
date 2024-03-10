set -ex
if [[ ${GITHUB_ACTIONS} -ne 1 ]]; then
    sage phi.sage >phi.log
fi
go run ./cmd/prove_order >prove_order.log
go run ./cmd/encoding >encoding.log
go run ./cmd/phi_surj >phi_surj.log
