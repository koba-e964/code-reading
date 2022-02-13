#include "ubc_check.h"

#include <iostream>

using namespace std;

uint32_t rot(uint32_t v, int sh) {
  return sh == 0 ? v : v << sh | v >> (32 - sh);
}

int main(void) {
  for (int i = 0; sha1_dvs[i].dvType; i++) {
    for (int j = 16; j < 80; j++) {
      uint32_t diff = sha1_dvs[i].dm[j];
      diff ^= rot(sha1_dvs[i].dm[j - 3] ^ sha1_dvs[i].dm[j - 8] ^ sha1_dvs[i].dm[j - 14] ^ sha1_dvs[i].dm[j - 16], 1);
      assert (diff == 0);
    }
    int weight = 0;
    for (int j = 0; j < 80; j++) {
      weight += __builtin_popcount(sha1_dvs[i].dm[j]);
    }
    cerr << "verified dvs[" << i << "] = DV(" << sha1_dvs[i].dvType
         << "," << sha1_dvs[i].dvK
         << "," << sha1_dvs[i].dvB
         << "): weight = " << weight << endl;
  }
}
