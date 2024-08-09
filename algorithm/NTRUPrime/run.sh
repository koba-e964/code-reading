set -e

URI=https://ntruprime.cr.yp.to/nist/ntruprime-20201007.pdf
curl -sS ${URI} -o ntruprime-20201007.pdf
sha256sum --check sha256sum.txt
